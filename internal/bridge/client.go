// Package bridge speaks the MCP Streamable HTTP transport (spec 2025-06-18) on behalf of a
// stdio client: Client handles one session's HTTP side, and Run pumps newline-delimited
// JSON-RPC between stdin/stdout and a server URL.
package bridge

import (
	"bytes"
	"context"
	"encoding/json"
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
	acceptPost            = "application/json, text/event-stream"
	acceptStream          = "text/event-stream"
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

// Client is a session-aware Streamable HTTP client for one MCP endpoint.
type Client struct {
	URL     string
	Header  http.Header   // extra headers sent on every request
	HTTP    *http.Client  // nil = a client without an overall timeout (streams are long-lived)
	Timeout time.Duration // per POST/DELETE; 0 = none

	mu              sync.RWMutex
	sessionID       string
	protocolVersion string
}

func (c *Client) SessionID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.sessionID
}

func (c *Client) ProtocolVersion() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.protocolVersion
}

func (c *Client) httpClient() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return http.DefaultClient
}

func (c *Client) newRequest(ctx context.Context, method string, body []byte) (*http.Request, error) {
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
	if sid := c.SessionID(); sid != "" {
		req.Header.Set(HeaderSessionID, sid)
	}
	if pv := c.ProtocolVersion(); pv != "" {
		req.Header.Set(HeaderProtocolVersion, pv)
	}
	return req, nil
}

func (c *Client) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if c.Timeout > 0 {
		return context.WithTimeout(ctx, c.Timeout)
	}
	return context.WithCancel(ctx)
}

type envelope struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Result *struct {
		ProtocolVersion string `json:"protocolVersion"`
	} `json:"result"`
}

// Post sends one JSON-RPC message and passes every JSON-RPC message in the response
// (a JSON body or each SSE message event) to emit. 202/204 emit nothing. A successful
// initialize response sets the session id and negotiated protocol version.
func (c *Client) Post(ctx context.Context, msg []byte, emit func([]byte) error) error {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	req, err := c.newRequest(ctx, http.MethodPost, msg)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", acceptPost)
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var in envelope
	_ = json.Unmarshal(msg, &in)
	if sid := resp.Header.Get(HeaderSessionID); sid != "" && in.Method == "initialize" {
		c.mu.Lock()
		c.sessionID = sid
		c.mu.Unlock()
	}
	if resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusNoContent {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return httpError(resp)
	}
	if in.Method == "initialize" {
		emit = c.sniffInitialize(emit)
	}
	switch mediaType(resp) {
	case acceptStream:
		return ParseSSE(resp.Body, func(ev Event) error {
			if !ev.IsMessage() || strings.TrimSpace(ev.Data) == "" {
				return nil
			}
			return emit([]byte(ev.Data))
		})
	case "application/json", "":
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if body = bytes.TrimSpace(body); len(body) == 0 {
			return nil
		}
		if !json.Valid(body) {
			return fmt.Errorf("invalid JSON response body: %s", snippet(body))
		}
		return emit(body)
	default:
		return fmt.Errorf("unexpected response Content-Type %q", resp.Header.Get("Content-Type"))
	}
}

func (c *Client) sniffInitialize(emit func([]byte) error) func([]byte) error {
	return func(b []byte) error {
		var env envelope
		if json.Unmarshal(b, &env) == nil && env.Result != nil && env.Result.ProtocolVersion != "" {
			c.mu.Lock()
			c.protocolVersion = env.Result.ProtocolVersion
			c.mu.Unlock()
		}
		return emit(b)
	}
}

// Stream opens the GET SSE stream for server-initiated messages and blocks until it ends.
// supported is false (with a nil error) when the server does not offer a stream
// (405, 404, or a non-SSE response), as legacy servers do.
func (c *Client) Stream(ctx context.Context, emit func([]byte) error) (supported bool, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Accept", acceptStream)
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return true, err
	}
	defer resp.Body.Close()
	switch {
	case resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusNotFound:
		return false, nil
	case resp.StatusCode < 200 || resp.StatusCode > 299:
		return false, httpError(resp)
	case mediaType(resp) != acceptStream:
		return false, nil
	}
	return true, ParseSSE(resp.Body, func(ev Event) error {
		if !ev.IsMessage() || strings.TrimSpace(ev.Data) == "" {
			return nil
		}
		return emit([]byte(ev.Data))
	})
}

// Delete ends the session (best effort; 405 means the server does not support it).
func (c *Client) Delete(ctx context.Context) error {
	if c.SessionID() == "" {
		return nil
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	req, err := c.newRequest(ctx, http.MethodDelete, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusMethodNotAllowed || (resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil
	}
	return httpError(resp)
}

func mediaType(resp *http.Response) string {
	mt, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	return strings.ToLower(mt)
}

func httpError(resp *http.Response) error {
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
