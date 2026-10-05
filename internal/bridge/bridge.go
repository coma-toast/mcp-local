package bridge

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Options configures Run.
type Options struct {
	URL            string
	Header         http.Header
	Timeout        time.Duration // per POST request (default 60s)
	Stdin          io.Reader
	Stdout         io.Writer
	Log            *log.Logger   // stderr logger; nil discards
	HTTPClient     *http.Client  // nil = default client (no overall timeout)
	ReconnectDelay time.Duration // GET stream reconnect delay (default 1s)
}

type bridge struct {
	c      *Client
	out    io.Writer
	outMu  sync.Mutex
	log    *log.Logger
	delay  time.Duration
	stream sync.Once
	wg     sync.WaitGroup // in-flight requests
	swg    sync.WaitGroup // GET stream goroutine
}

// Run bridges newline-delimited JSON-RPC on Stdin/Stdout to the Streamable HTTP server at URL.
// It returns nil on stdin EOF after in-flight requests finish and the session is deleted.
func Run(ctx context.Context, o Options) error {
	if o.Timeout <= 0 {
		o.Timeout = 60 * time.Second
	}
	if o.ReconnectDelay <= 0 {
		o.ReconnectDelay = time.Second
	}
	if o.Log == nil {
		o.Log = log.New(io.Discard, "", 0)
	}
	b := &bridge{
		c:     &Client{URL: o.URL, Header: o.Header, HTTP: o.HTTPClient, Timeout: o.Timeout},
		out:   o.Stdout,
		log:   o.Log,
		delay: o.ReconnectDelay,
	}
	streamCtx, stopStream := context.WithCancel(ctx)
	defer stopStream()
	lines := make(chan []byte)
	readErr := make(chan error, 1)
	go readLines(o.Stdin, lines, readErr)
	var err error
loop:
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			break loop
		case line, ok := <-lines:
			if !ok {
				if rerr := <-readErr; rerr != nil {
					err = fmt.Errorf("read stdin: %w", rerr)
				}
				break loop
			}
			b.handle(ctx, streamCtx, line)
		}
	}
	b.wg.Wait()
	stopStream()
	b.swg.Wait()
	dctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if derr := b.c.Delete(dctx); derr != nil {
		b.log.Printf("delete session: %v", derr)
	}
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}

func readLines(r io.Reader, lines chan<- []byte, errc chan<- error) {
	defer close(lines)
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadBytes('\n')
		if msg := bytes.TrimSpace(line); len(msg) > 0 {
			lines <- msg
		}
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			errc <- err
			return
		}
	}
}

// handle POSTs one message. initialize, notifications, and client responses are sent in order;
// other requests run concurrently so a slow tool call does not block cancellations.
func (b *bridge) handle(ctx, streamCtx context.Context, msg []byte) {
	var env envelope
	_ = json.Unmarshal(msg, &env)
	isRequest := env.Method != "" && len(env.ID) > 0 && string(env.ID) != "null"
	switch {
	case env.Method == "initialize":
		b.post(ctx, msg, env.ID, isRequest)
		if b.c.ProtocolVersion() != "" {
			b.stream.Do(func() { b.startStream(streamCtx) })
		}
	case isRequest:
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			b.post(ctx, msg, env.ID, true)
		}()
	default:
		b.post(ctx, msg, env.ID, false)
	}
}

func (b *bridge) post(ctx context.Context, msg []byte, id json.RawMessage, isRequest bool) {
	responded := false
	err := b.c.Post(ctx, msg, func(resp []byte) error {
		var env envelope
		isResponse := json.Unmarshal(resp, &env) == nil && env.Method == ""
		if !isRequest && isResponse {
			// Notifications have no response; some legacy servers send one anyway.
			b.log.Printf("POST %s: dropped unexpected response %s", describe(msg), snippet(resp))
			return nil
		}
		if isRequest && isResponse && sameID(env.ID, id) {
			responded = true
		}
		return b.write(resp)
	})
	if err == nil {
		return
	}
	if !isRequest {
		b.log.Printf("POST %s: %v", describe(msg), err)
		return
	}
	if responded {
		b.log.Printf("POST %s after response: %v", describe(msg), err)
		return
	}
	b.log.Printf("POST %s: %v", describe(msg), err)
	_ = b.write(errorResponse(id, err))
}

func (b *bridge) startStream(ctx context.Context) {
	b.swg.Add(1)
	go func() {
		defer b.swg.Done()
		for {
			supported, err := b.c.Stream(ctx, b.write)
			if ctx.Err() != nil {
				return
			}
			if !supported {
				if err != nil {
					b.log.Printf("GET stream: %v", err)
				}
				return
			}
			if err != nil {
				b.log.Printf("GET stream ended: %v; reconnecting", err)
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(b.delay):
			}
		}
	}()
}

// write emits one message as a single newline-terminated stdout line.
func (b *bridge) write(msg []byte) error {
	var buf bytes.Buffer
	if err := json.Compact(&buf, msg); err != nil {
		buf.Reset()
		buf.WriteString(strings.NewReplacer("\r", " ", "\n", " ").Replace(string(msg)))
	}
	buf.WriteByte('\n')
	b.outMu.Lock()
	defer b.outMu.Unlock()
	_, err := b.out.Write(buf.Bytes())
	return err
}

func errorResponse(id json.RawMessage, err error) []byte {
	resp, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"error":   map[string]interface{}{"code": -32000, "message": "bridge: " + err.Error()},
	})
	return resp
}

func sameID(a, b json.RawMessage) bool {
	var x, y interface{}
	return json.Unmarshal(a, &x) == nil && json.Unmarshal(b, &y) == nil && fmt.Sprintf("%T:%v", x, x) == fmt.Sprintf("%T:%v", y, y)
}

func describe(msg []byte) string {
	var env envelope
	if json.Unmarshal(msg, &env) == nil && env.Method != "" {
		if len(env.ID) > 0 {
			return fmt.Sprintf("%s (id %s)", env.Method, env.ID)
		}
		return env.Method
	}
	return snippet(msg)
}
