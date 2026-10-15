// Package bridge speaks the MCP Streamable HTTP transport (spec 2025-06-18) on behalf of a
// stdio client: Client handles one session's HTTP side, and Run pumps newline-delimited
// JSON-RPC between stdin/stdout and a server URL.
package bridge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	HeaderSessionID       = "Mcp-Session-Id"
	HeaderProtocolVersion = "MCP-Protocol-Version"
	HeaderLastEventID     = "Last-Event-ID"
	acceptPost            = "application/json, text/event-stream"
	acceptStream          = "text/event-stream"
)

var (
	// ErrTimeout means the server sent no response headers, or went silent on a response
	// body, for longer than Client.Timeout.
	ErrTimeout = errors.New("timed out waiting for the server")
	// ErrNoStream means the server does not offer a GET stream (405, or a legacy server).
	ErrNoStream = errors.New("server does not offer a GET stream")
)

// HTTPError is a non-2xx response from the server.
type HTTPError struct {
	Status int
	Text   string
}

func (e *HTTPError) Error() string {
	msg := fmt.Sprintf("HTTP %d %s", e.Status, http.StatusText(e.Status))
	if e.Text != "" {
		msg += ": " + e.Text
	}
	return msg
}

// SessionExpiredError is a 404 to a request that carried a session id: the server no longer
// knows the session (it restarted, or expired it).
type SessionExpiredError struct{ Err *HTTPError }

func (e *SessionExpiredError) Error() string { return "session expired: " + e.Err.Error() }
func (e *SessionExpiredError) Unwrap() error { return e.Err }

// Client is a session-aware Streamable HTTP client for one MCP endpoint. It caches the
// client's initialize and notifications/initialized messages, and when the server forgets
// the session (404) it replays them to open a new session and retries the message once.
type Client struct {
	URL     string
	Header  http.Header            // extra headers sent on every request
	HTTP    *http.Client           // nil = a client without an overall timeout (streams are long-lived)
	Timeout time.Duration          // max wait for POST response headers, and max silence on a POST response body; 0 = none
	Logf    func(string, ...any)   // diagnostics; nil discards
	Notify  func(msg []byte) error // server requests/notifications seen while re-initializing; nil drops

	initMu         sync.Mutex // serializes re-initialization
	mu             sync.Mutex // guards the fields below
	sess           session
	initMsg        []byte // the client's initialize, as sent
	initializedMsg []byte // the client's notifications/initialized, as sent
}

type session struct {
	id, version string
	gen         uint64 // bumped on every session change
	lost        bool   // dropped after a 404 and not yet re-established
}

func (c *Client) SessionID() string       { return c.state().id }
func (c *Client) ProtocolVersion() string { return c.state().version }

func (c *Client) state() session {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sess
}

func (c *Client) setSession(s session) {
	c.mu.Lock()
	s.gen = c.sess.gen + 1
	c.sess = s
	c.mu.Unlock()
}

func (c *Client) logf(format string, args ...any) {
	if c.Logf != nil {
		c.Logf(format, args...)
	}
}

func (c *Client) httpClient() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return http.DefaultClient
}

// newRequest builds a request; st is nil for initialize, which never carries session headers.
func (c *Client) newRequest(ctx context.Context, method string, body []byte, st *session) (*http.Request, error) {
	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.URL, r)
	if err != nil {
		return nil, err
	}
	for k, vs := range c.Header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	if st != nil && st.id != "" {
		req.Header.Set(HeaderSessionID, st.id)
	}
	if st != nil && st.version != "" {
		req.Header.Set(HeaderProtocolVersion, st.version)
	}
	return req, nil
}

type envelope struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Result *struct {
		ProtocolVersion string `json:"protocolVersion"`
	} `json:"result"`
	Error json.RawMessage `json:"error"`
}

// Post sends one JSON-RPC message and passes every message in the response (a JSON body or
// each SSE message event) to emit. 202/204 emit nothing. A successful initialize response
// sets the session id and negotiated protocol version. A 404 for a message sent with a
// session re-initializes the session and retries the message once.
func (c *Client) Post(ctx context.Context, msg []byte, emit func([]byte) error) error {
	var in envelope
	_ = json.Unmarshal(msg, &in)
	switch in.Method {
	case "initialize":
		c.mu.Lock()
		c.initMsg = msg
		c.mu.Unlock()
		return c.post(ctx, msg, nil, c.sniffInitialize(emit))
	case "notifications/initialized":
		c.mu.Lock()
		c.initializedMsg = msg
		c.mu.Unlock()
	}
	retried := false
	for range 4 {
		st := c.state()
		if st.lost {
			if err := c.reinit(ctx, st.gen); err != nil {
				return err
			}
			continue
		}
		err := c.post(ctx, msg, &st, emit)
		var se *SessionExpiredError
		if retried || !errors.As(err, &se) {
			return err
		}
		retried = true
		if rerr := c.reinit(ctx, st.gen); rerr != nil {
			return fmt.Errorf("%w; re-initialize: %v", err, rerr)
		}
	}
	return errors.New("session keeps expiring")
}

func (c *Client) post(ctx context.Context, msg []byte, st *session, emit func([]byte) error) error {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	wd := c.watchdog(cancel)
	defer wd.stop()
	req, err := c.newRequest(ctx, http.MethodPost, msg, st)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", acceptPost)
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return causeOr(ctx, err)
	}
	defer resp.Body.Close()
	wd.kick()
	body := &kickReader{r: resp.Body, wd: wd}
	if resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusNoContent {
		_, _ = io.Copy(io.Discard, body)
		return nil
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		herr := httpError(resp)
		if resp.StatusCode == http.StatusNotFound && st != nil && st.id != "" {
			return &SessionExpiredError{Err: herr}
		}
		return herr
	}
	if st == nil {
		c.setSession(session{id: resp.Header.Get(HeaderSessionID)})
	}
	switch mediaType(resp) {
	case acceptStream:
		err = ParseSSE(body, func(ev Event) error {
			if !ev.IsMessage() || strings.TrimSpace(ev.Data) == "" {
				return nil
			}
			return emit([]byte(ev.Data))
		})
	case "application/json", "":
		var b []byte
		if b, err = io.ReadAll(body); err == nil {
			if b = bytes.TrimSpace(b); len(b) > 0 {
				err = emit(b)
			}
		}
	default:
		return fmt.Errorf("unexpected response Content-Type %q", resp.Header.Get("Content-Type"))
	}
	return causeOr(ctx, err)
}

func (c *Client) sniffInitialize(emit func([]byte) error) func([]byte) error {
	return func(b []byte) error {
		var env envelope
		if json.Unmarshal(b, &env) == nil && env.Result != nil && env.Result.ProtocolVersion != "" {
			c.mu.Lock()
			c.sess.version = env.Result.ProtocolVersion
			c.mu.Unlock()
		}
		return emit(b)
	}
}

// reinit replaces the session that was current at generation gen by replaying the client's
// initialize (its response is consumed here) and notifications/initialized. Concurrent
// callers for the same generation wait for the first and then reuse its session.
func (c *Client) reinit(ctx context.Context, gen uint64) error {
	c.initMu.Lock()
	defer c.initMu.Unlock()
	c.mu.Lock()
	if c.sess.gen != gen {
		c.mu.Unlock()
		return nil
	}
	stale, initMsg, initializedMsg := c.sess.id, c.initMsg, c.initializedMsg
	c.sess = session{gen: gen + 1, lost: true}
	c.mu.Unlock()
	if initMsg == nil {
		return errors.New("no initialize to replay")
	}
	c.logf("session %q expired; re-initializing", stale)
	err := c.replayInit(ctx, initMsg, initializedMsg)
	if err != nil {
		c.setSession(session{lost: true})
		return err
	}
	c.logf("re-initialized as session %q", c.SessionID())
	return nil
}

func (c *Client) replayInit(ctx context.Context, initMsg, initializedMsg []byte) error {
	var in envelope
	_ = json.Unmarshal(initMsg, &in)
	responded, rpcErr := false, json.RawMessage(nil)
	err := c.post(ctx, initMsg, nil, c.sniffInitialize(func(m []byte) error {
		var env envelope
		switch {
		case json.Unmarshal(m, &env) != nil:
		case env.Method != "":
			if c.Notify != nil {
				return c.Notify(m)
			}
		case sameID(env.ID, in.ID):
			responded, rpcErr = true, env.Error
		}
		return nil
	}))
	switch {
	case err != nil:
		return fmt.Errorf("initialize: %w", err)
	case !responded:
		return errors.New("initialize: stream ended without response")
	case len(rpcErr) > 0 && string(rpcErr) != "null":
		return fmt.Errorf("initialize: %s", snippet(rpcErr))
	case initializedMsg == nil:
		return nil
	}
	st := c.state()
	return c.post(ctx, initializedMsg, &st, func([]byte) error { return nil })
}

// StreamCursor carries GET stream state across reconnects.
type StreamCursor struct {
	LastEventID string        // sent as Last-Event-ID on reconnect
	Retry       time.Duration // the server's latest SSE retry: hint
}

// Stream opens the GET SSE stream for server-initiated messages and blocks until it ends.
// connected reports whether the server accepted the stream. It returns ErrNoStream when the
// server does not offer one (405, a 404 without a session, or a non-SSE response). A 404 for
// a session re-initializes it and reopens the stream once.
func (c *Client) Stream(ctx context.Context, cur *StreamCursor, emit func([]byte) error) (connected bool, err error) {
	for attempt := 0; ; attempt++ {
		st := c.state()
		if st.lost {
			if err := c.reinit(ctx, st.gen); err != nil {
				return false, err
			}
			st = c.state()
		}
		connected, err = c.stream(ctx, &st, cur, emit)
		var se *SessionExpiredError
		if attempt > 0 || !errors.As(err, &se) {
			return connected, err
		}
		if rerr := c.reinit(ctx, st.gen); rerr != nil {
			return false, fmt.Errorf("%w; re-initialize: %v", err, rerr)
		}
	}
}

func (c *Client) stream(ctx context.Context, st *session, cur *StreamCursor, emit func([]byte) error) (bool, error) {
	req, err := c.newRequest(ctx, http.MethodGet, nil, st)
	if err != nil {
		return false, err
	}
	req.Header.Set("Accept", acceptStream)
	if cur.LastEventID != "" {
		req.Header.Set(HeaderLastEventID, cur.LastEventID)
	}
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	switch {
	case resp.StatusCode == http.StatusMethodNotAllowed:
		return false, ErrNoStream
	case resp.StatusCode == http.StatusNotFound && st.id != "":
		return false, &SessionExpiredError{Err: httpError(resp)}
	case resp.StatusCode == http.StatusNotFound:
		return false, ErrNoStream
	case resp.StatusCode < 200 || resp.StatusCode > 299:
		return false, httpError(resp)
	case mediaType(resp) != acceptStream:
		return false, ErrNoStream
	}
	return true, ParseSSE(resp.Body, func(ev Event) error {
		if ev.ID != "" {
			cur.LastEventID = ev.ID
		}
		if ev.Retry > 0 {
			cur.Retry = ev.Retry
		}
		if !ev.IsMessage() || strings.TrimSpace(ev.Data) == "" {
			return nil
		}
		return emit([]byte(ev.Data))
	})
}

// Delete ends the session (best effort; 404/405 mean there is nothing to end). It is bounded
// by ctx and by Timeout.
func (c *Client) Delete(ctx context.Context) error {
	st := c.state()
	if st.id == "" {
		return nil
	}
	if c.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.Timeout)
		defer cancel()
	}
	req, err := c.newRequest(ctx, http.MethodDelete, nil, &st)
	if err != nil {
		return err
	}
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusNotFound || (resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil
	}
	return httpError(resp)
}

// watchdog cancels a request with ErrTimeout unless it is kicked within d.
type watchdog struct {
	t *time.Timer
	d time.Duration
}

func (c *Client) watchdog(cancel context.CancelCauseFunc) *watchdog {
	w := &watchdog{d: c.Timeout}
	if w.d > 0 {
		w.t = time.AfterFunc(w.d, func() { cancel(fmt.Errorf("%w (%s)", ErrTimeout, w.d)) })
	}
	return w
}

func (w *watchdog) kick() {
	if w.t != nil {
		w.t.Reset(w.d)
	}
}

func (w *watchdog) stop() {
	if w.t != nil {
		w.t.Stop()
	}
}

// kickReader kicks the watchdog whenever response bytes arrive, so a stream that keeps
// sending (progress, keepalives) is never cut off, but a silent one is.
type kickReader struct {
	r  io.Reader
	wd *watchdog
}

func (k *kickReader) Read(p []byte) (int, error) {
	n, err := k.r.Read(p)
	if n > 0 {
		k.wd.kick()
	}
	return n, err
}

// causeOr prefers the context's cancellation cause (timeout, shutdown) over the transport
// error it produced.
func causeOr(ctx context.Context, err error) error {
	if err != nil && ctx.Err() != nil {
		return context.Cause(ctx)
	}
	return err
}

func mediaType(resp *http.Response) string {
	mt, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	return strings.ToLower(mt)
}

func httpError(resp *http.Response) *HTTPError {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return &HTTPError{Status: resp.StatusCode, Text: snippet(bytes.TrimSpace(body))}
}

func snippet(b []byte) string {
	s := strings.Join(strings.Fields(string(b)), " ")
	if len(s) > 200 {
		s = s[:200] + "…"
	}
	return s
}
