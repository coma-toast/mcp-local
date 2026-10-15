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
	"math/rand/v2"
	"net/http"
	"sync"
	"time"
)

const (
	cancelTimeout = time.Second     // bound on a notifications/cancelled POST
	deleteTimeout = 2 * time.Second // bound on the session DELETE at exit
	logEvery      = 30 * time.Second
)

var (
	errNoResponse = errors.New("stream ended without response")
	errShutdown   = errors.New("bridge shutting down")
)

// Options configures Run.
type Options struct {
	URL               string
	Header            http.Header
	Timeout           time.Duration // max wait for response headers, and max silence on a response stream (default 60s)
	Stdin             io.Reader
	Stdout            io.Writer
	Log               *log.Logger   // stderr logger; nil discards
	HTTPClient        *http.Client  // nil = default client (no overall timeout)
	ReconnectDelay    time.Duration // first GET stream reconnect delay, doubled per failure (default 500ms)
	MaxReconnectDelay time.Duration // cap on the GET stream reconnect delay (default 30s)
	ShutdownGrace     time.Duration // how long in-flight requests may finish after stdin EOF (default 2s)
}

type bridge struct {
	c        *Client
	out      io.Writer
	outMu    sync.Mutex
	log      *log.Logger
	minDelay time.Duration
	maxDelay time.Duration
	stream   sync.Once
	wg       sync.WaitGroup // in-flight requests
	swg      sync.WaitGroup // GET stream goroutine
	mu       sync.Mutex
	calls    map[string]*call // in-flight requests by id
}

type call struct{ clientCancelled bool } // guarded by bridge.mu

// Run bridges newline-delimited JSON-RPC on Stdin/Stdout to the Streamable HTTP server at URL.
// On stdin EOF it gives in-flight requests ShutdownGrace to finish, cancels the rest, deletes
// the session and returns nil.
func Run(ctx context.Context, o Options) error {
	if o.Timeout <= 0 {
		o.Timeout = 60 * time.Second
	}
	if o.ReconnectDelay <= 0 {
		o.ReconnectDelay = 500 * time.Millisecond
	}
	if o.MaxReconnectDelay <= 0 {
		o.MaxReconnectDelay = 30 * time.Second
	}
	if o.ShutdownGrace <= 0 {
		o.ShutdownGrace = 2 * time.Second
	}
	if o.Log == nil {
		o.Log = log.New(io.Discard, "", 0)
	}
	b := &bridge{out: o.Stdout, log: o.Log, minDelay: o.ReconnectDelay, maxDelay: o.MaxReconnectDelay, calls: map[string]*call{}}
	b.c = &Client{URL: o.URL, Header: o.Header, HTTP: o.HTTPClient, Timeout: o.Timeout, Logf: b.log.Printf, Notify: b.forward}
	// Requests outlive stdin EOF (for the grace period) but not a signal.
	reqCtx, cancelReqs := context.WithCancelCause(context.WithoutCancel(ctx))
	defer cancelReqs(nil)
	defer context.AfterFunc(ctx, func() { cancelReqs(errShutdown) })()
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
			b.handle(reqCtx, streamCtx, line)
		}
	}
	b.drain(o.ShutdownGrace, func() { cancelReqs(errShutdown) })
	stopStream()
	b.swg.Wait()
	dctx, cancel := context.WithTimeout(context.Background(), deleteTimeout)
	defer cancel()
	if derr := b.c.Delete(dctx); derr != nil {
		b.log.Printf("delete session: %v", derr)
	}
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}

// drain waits up to grace for in-flight requests, then cancels the rest and waits for them
// to report their errors.
func (b *bridge) drain(grace time.Duration, cancel func()) {
	done := make(chan struct{})
	go func() { b.wg.Wait(); close(done) }()
	select {
	case <-done:
		return
	case <-time.After(grace):
	}
	cancel()
	<-done
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
	switch {
	case !json.Valid(msg):
		b.log.Printf("stdin: invalid JSON: %s", snippet(msg))
		_ = b.write(errorResponse(nil, -32700, "Parse error"))
		return
	case msg[0] == '[':
		b.rejectBatch(msg)
		return
	case msg[0] != '{':
		b.log.Printf("stdin: not a JSON-RPC message: %s", snippet(msg))
		_ = b.write(errorResponse(nil, -32600, "Invalid Request"))
		return
	}
	var env envelope
	_ = json.Unmarshal(msg, &env)
	isRequest := env.Method != "" && hasID(env.ID)
	switch {
	case env.Method == "initialize" && isRequest:
		b.wg.Add(1)
		b.request(ctx, msg, env.ID)
		if b.c.ProtocolVersion() != "" {
			b.stream.Do(func() { b.startStream(streamCtx) })
		}
	case isRequest:
		b.wg.Add(1)
		go b.request(ctx, msg, env.ID)
	default:
		if env.Method == "notifications/cancelled" {
			b.markCancelled(msg)
		}
		b.notify(ctx, msg)
	}
}

// rejectBatch answers each request in a JSON array with -32600: MCP 2025-06-18 removed
// JSON-RPC batching. Notifications and responses in the batch are dropped.
func (b *bridge) rejectBatch(msg []byte) {
	var items []json.RawMessage
	_ = json.Unmarshal(msg, &items)
	b.log.Printf("stdin: rejected a JSON-RPC batch of %d messages (batching is not supported)", len(items))
	if len(items) == 0 {
		_ = b.write(errorResponse(nil, -32600, "Invalid Request: empty batch"))
		return
	}
	for _, it := range items {
		var env envelope
		if json.Unmarshal(it, &env) == nil && env.Method != "" && hasID(env.ID) {
			_ = b.write(errorResponse(env.ID, -32600, "Invalid Request: JSON-RPC batches are not supported"))
		}
	}
}

// request POSTs a client request and guarantees the client exactly one response for it,
// unless the client cancelled it. The caller has done b.wg.Add(1).
func (b *bridge) request(ctx context.Context, msg []byte, id json.RawMessage) {
	defer b.wg.Done()
	cl, key := &call{}, idKey(id)
	b.mu.Lock()
	b.calls[key] = cl
	b.mu.Unlock()
	defer func() {
		b.mu.Lock()
		if b.calls[key] == cl {
			delete(b.calls, key)
		}
		b.mu.Unlock()
	}()
	responded := false
	err := b.c.Post(ctx, msg, func(resp []byte) error {
		var env envelope
		if json.Unmarshal(resp, &env) == nil && env.Method == "" && sameID(env.ID, id) {
			responded = true
		}
		return b.forward(resp)
	})
	b.mu.Lock()
	clientCancelled := cl.clientCancelled
	b.mu.Unlock()
	switch {
	case responded && err != nil:
		b.log.Printf("POST %s after response: %v", describe(msg), err)
		return
	case responded:
		return
	case clientCancelled:
		b.log.Printf("POST %s: cancelled by client (%v)", describe(msg), err)
		return
	case err == nil:
		err = errNoResponse
	case errors.Is(err, ErrTimeout) || errors.Is(err, errShutdown):
		b.cancelRemote(id, err)
	}
	b.log.Printf("POST %s: %v", describe(msg), err)
	_ = b.write(errorResponse(id, -32000, "bridge: "+err.Error()))
}

func (b *bridge) notify(ctx context.Context, msg []byte) {
	err := b.c.Post(ctx, msg, func(resp []byte) error {
		var env envelope
		if json.Unmarshal(resp, &env) == nil && env.Method == "" {
			// Notifications have no response; some legacy servers send one anyway.
			b.log.Printf("POST %s: dropped unexpected response %s", describe(msg), snippet(resp))
			return nil
		}
		return b.forward(resp)
	})
	if err != nil {
		b.log.Printf("POST %s: %v", describe(msg), err)
	}
}

// markCancelled notes a client notifications/cancelled so the request's missing response
// is expected rather than reported as an error.
func (b *bridge) markCancelled(msg []byte) {
	var n struct {
		Params struct {
			RequestID json.RawMessage `json:"requestId"`
		} `json:"params"`
	}
	if json.Unmarshal(msg, &n) != nil || !hasID(n.Params.RequestID) {
		return
	}
	b.mu.Lock()
	if cl := b.calls[idKey(n.Params.RequestID)]; cl != nil {
		cl.clientCancelled = true
	}
	b.mu.Unlock()
}

// cancelRemote tells the server the bridge gave up on request id.
func (b *bridge) cancelRemote(id json.RawMessage, reason error) {
	note, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"method":  "notifications/cancelled",
		"params":  map[string]any{"requestId": id, "reason": "bridge: " + reason.Error()},
	})
	ctx, cancel := context.WithTimeout(context.Background(), cancelTimeout)
	defer cancel()
	if err := b.c.Post(ctx, note, func([]byte) error { return nil }); err != nil {
		b.log.Printf("POST notifications/cancelled (id %s): %v", id, err)
	}
}

// startStream keeps the GET stream open, reconnecting with capped exponential backoff and
// jitter until the server says it has no stream (405) or ctx ends.
func (b *bridge) startStream(ctx context.Context) {
	b.swg.Add(1)
	go func() {
		defer b.swg.Done()
		cur, delay, lim := &StreamCursor{}, b.minDelay, &limiter{every: logEvery}
		for {
			connected, err := b.c.Stream(ctx, cur, b.forward)
			if ctx.Err() != nil {
				return
			}
			if errors.Is(err, ErrNoStream) {
				b.log.Printf("GET stream: %v", err)
				return
			}
			if connected {
				delay = b.minDelay
			}
			wait := delay/2 + rand.N(delay/2+1)
			wait = max(wait, cur.Retry)
			if err == nil {
				err = errors.New("closed by server")
			}
			lim.printf(b.log, "GET stream: %v; reconnecting in %s", err, wait.Round(time.Millisecond))
			if !connected {
				delay = min(delay*2, b.maxDelay)
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(wait):
			}
		}
	}()
}

// limiter rate-limits a repeating log line, counting what it suppressed.
type limiter struct {
	every   time.Duration
	last    time.Time
	skipped int
}

func (l *limiter) printf(lg *log.Logger, format string, args ...any) {
	if !l.last.IsZero() && time.Since(l.last) < l.every {
		l.skipped++
		return
	}
	if l.skipped > 0 {
		format += fmt.Sprintf(" (%d similar suppressed)", l.skipped)
	}
	l.last, l.skipped = time.Now(), 0
	lg.Printf(format, args...)
}

// forward writes one server message to stdout; anything that is not JSON is dropped.
func (b *bridge) forward(msg []byte) error {
	if !json.Valid(msg) {
		b.log.Printf("dropped non-JSON data from server: %s", snippet(msg))
		return nil
	}
	return b.write(msg)
}

// write emits one JSON message as a single newline-terminated stdout line.
func (b *bridge) write(msg []byte) error {
	var buf bytes.Buffer
	if err := json.Compact(&buf, msg); err != nil {
		return err
	}
	buf.WriteByte('\n')
	b.outMu.Lock()
	defer b.outMu.Unlock()
	_, err := b.out.Write(buf.Bytes())
	return err
}

func errorResponse(id json.RawMessage, code int, msg string) []byte {
	if id == nil {
		id = json.RawMessage("null")
	}
	resp, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
		"error":   map[string]any{"code": code, "message": msg},
	})
	return resp
}

func hasID(id json.RawMessage) bool { return len(id) > 0 && string(id) != "null" }

// idKey normalizes a JSON-RPC id, so 1 and 1.0 match but 1 and "1" do not.
func idKey(id json.RawMessage) string {
	var x any
	_ = json.Unmarshal(id, &x)
	return fmt.Sprintf("%T:%v", x, x)
}

func sameID(a, b json.RawMessage) bool { return json.Valid(a) && json.Valid(b) && idKey(a) == idKey(b) }

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
