package bridge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// fakeServer is a minimal Streamable HTTP MCP server.
type fakeServer struct {
	t         *testing.T
	stream    bool // offer a GET SSE stream
	mu        sync.Mutex
	sessions  []string // Mcp-Session-Id seen on each non-initialize POST
	versions  []string // MCP-Protocol-Version seen on each non-initialize POST
	deleted   string
	getHeader http.Header
}

func (f *fakeServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if !f.stream {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		f.mu.Lock()
		f.getHeader = r.Header.Clone()
		f.mu.Unlock()
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, ": keepalive\n\nevent: message\ndata: {\"jsonrpc\":\"2.0\",\"method\":\"notifications/tools/list_changed\"}\n\n")
		w.(http.Flusher).Flush()
		<-r.Context().Done()
		return
	case http.MethodDelete:
		f.mu.Lock()
		f.deleted = r.Header.Get(HeaderSessionID)
		f.mu.Unlock()
		w.WriteHeader(http.StatusOK)
		return
	}
	if got := r.Header.Get("Accept"); got != "application/json, text/event-stream" {
		f.t.Errorf("Accept = %q", got)
	}
	if got := r.Header.Get("X-Test"); got != "yes" {
		f.t.Errorf("custom header X-Test = %q", got)
	}
	body, _ := io.ReadAll(r.Body)
	var req struct {
		ID     json.RawMessage `json:"id"`
		Method string          `json:"method"`
	}
	_ = json.Unmarshal(body, &req)
	if req.Method != "initialize" {
		f.mu.Lock()
		f.sessions = append(f.sessions, r.Header.Get(HeaderSessionID))
		f.versions = append(f.versions, r.Header.Get(HeaderProtocolVersion))
		f.mu.Unlock()
	}
	switch req.Method {
	case "initialize":
		w.Header().Set(HeaderSessionID, "sess-1")
		w.Header().Set("Content-Type", "application/json")
		// Pretty-printed on purpose: the bridge must emit one line.
		fmt.Fprintf(w, "{\n  \"jsonrpc\": \"2.0\",\n  \"id\": %s,\n  \"result\": {\"protocolVersion\": \"2025-06-18\", \"capabilities\": {}}\n}\n", req.ID)
	case "notifications/initialized":
		w.WriteHeader(http.StatusAccepted)
	case "notifications/legacy":
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"jsonrpc":"2.0","id":null,"error":{"code":-32601,"message":"Unknown method"}}`)
	case "tools/list":
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "event: message\ndata: {\"jsonrpc\":\"2.0\",\"method\":\"notifications/progress\",\"params\":{}}\n\n")
		fmt.Fprintf(w, "data: {\"jsonrpc\":\"2.0\",\ndata:  \"id\":%s,\ndata: \"result\":{\"tools\":[]}}\n\n", req.ID)
		fmt.Fprint(w, "event: other\ndata: ignored\n\n")
	case "boom":
		http.Error(w, "kaboom", http.StatusInternalServerError)
	case "notifications/boom":
		http.Error(w, "nope", http.StatusInternalServerError)
	default:
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"result":{}}`, req.ID)
	}
}

type lineBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (l *lineBuffer) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buf.Write(p)
}

func (l *lineBuffer) Lines() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	s := strings.TrimSuffix(l.buf.String(), "\n")
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

func (l *lineBuffer) waitFor(t *testing.T, substr string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		for _, line := range l.Lines() {
			if strings.Contains(line, substr) {
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %q; got %q", substr, l.Lines())
}

const (
	initMsg        = `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"t","version":"0"}}}`
	initializedMsg = `{"jsonrpc":"2.0","method":"notifications/initialized"}`
)

func run(t *testing.T, srv *fakeServer, stdin io.Reader) (*lineBuffer, *bytes.Buffer, error) {
	t.Helper()
	ts := httptest.NewServer(srv)
	defer ts.Close()
	out := &lineBuffer{}
	var logs bytes.Buffer
	err := Run(context.Background(), Options{
		URL:            ts.URL,
		Header:         http.Header{"X-Test": {"yes"}},
		Timeout:        5 * time.Second,
		Stdin:          stdin,
		Stdout:         out,
		Log:            log.New(&logs, "", 0),
		ReconnectDelay: 10 * time.Millisecond,
	})
	return out, &logs, err
}

func TestBridge_JSONAndSSEAndNotification(t *testing.T) {
	srv := &fakeServer{t: t}
	in := strings.Join([]string{initMsg, initializedMsg, `{"jsonrpc":"2.0","id":2,"method":"tools/list"}`}, "\n") + "\n"
	out, _, err := run(t, srv, strings.NewReader(in))
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	lines := out.Lines()
	want := []string{
		`{"jsonrpc":"2.0","id":1,"result":{"protocolVersion":"2025-06-18","capabilities":{}}}`,
		`{"jsonrpc":"2.0","method":"notifications/progress","params":{}}`,
		`{"jsonrpc":"2.0","id":2,"result":{"tools":[]}}`,
	}
	if strings.Join(lines, "\n") != strings.Join(want, "\n") {
		t.Errorf("stdout:\n%s\nwant:\n%s", strings.Join(lines, "\n"), strings.Join(want, "\n"))
	}
}

func TestBridge_NotificationProducesNoOutput(t *testing.T) {
	out, _, err := run(t, &fakeServer{t: t}, strings.NewReader(initializedMsg+"\n"+`{"jsonrpc":"2.0","method":"notifications/legacy"}`+"\n"))
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if lines := out.Lines(); len(lines) != 0 {
		t.Errorf("202 should produce no output, got %q", lines)
	}
}

func TestBridge_SessionHeaderRoundTrip(t *testing.T) {
	srv := &fakeServer{t: t}
	in := initMsg + "\n" + initializedMsg + "\n" + `{"jsonrpc":"2.0","id":2,"method":"ping"}` + "\n"
	if _, _, err := run(t, srv, strings.NewReader(in)); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(srv.sessions) != 2 {
		t.Fatalf("sessions = %v", srv.sessions)
	}
	for i := range srv.sessions {
		if srv.sessions[i] != "sess-1" || srv.versions[i] != "2025-06-18" {
			t.Errorf("request %d: session=%q version=%q", i, srv.sessions[i], srv.versions[i])
		}
	}
	if srv.deleted != "sess-1" {
		t.Errorf("DELETE session = %q, want sess-1", srv.deleted)
	}
}

func TestBridge_GETStreamForwarded(t *testing.T) {
	srv := &fakeServer{t: t, stream: true}
	pr, pw := io.Pipe()
	done := make(chan struct{})
	var out *lineBuffer
	var err error
	ts := httptest.NewServer(srv)
	defer ts.Close()
	out = &lineBuffer{}
	go func() {
		defer close(done)
		err = Run(context.Background(), Options{URL: ts.URL, Header: http.Header{"X-Test": {"yes"}}, Stdin: pr, Stdout: out, ReconnectDelay: 10 * time.Millisecond})
	}()
	fmt.Fprintln(pw, initMsg)
	out.waitFor(t, `"notifications/tools/list_changed"`)
	pw.Close()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not exit after stdin EOF")
	}
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.getHeader.Get(HeaderSessionID) != "sess-1" || srv.getHeader.Get("Accept") != "text/event-stream" {
		t.Errorf("GET headers = %v", srv.getHeader)
	}
}

func TestBridge_HTTPErrorBecomesJSONRPCError(t *testing.T) {
	in := `{"jsonrpc":"2.0","id":"req-7","method":"boom"}` + "\n" + `{"jsonrpc":"2.0","method":"notifications/boom"}` + "\n"
	out, logs, err := run(t, &fakeServer{t: t}, strings.NewReader(in))
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	lines := out.Lines()
	if len(lines) != 1 {
		t.Fatalf("stdout = %q, want exactly one error line", lines)
	}
	var resp struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      json.RawMessage `json:"id"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(lines[0]), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.JSONRPC != "2.0" || string(resp.ID) != `"req-7"` || resp.Error.Code != -32000 || !strings.HasPrefix(resp.Error.Message, "bridge: HTTP 500") || !strings.Contains(resp.Error.Message, "kaboom") {
		t.Errorf("error line = %s", lines[0])
	}
	if !strings.Contains(logs.String(), "notifications/boom") {
		t.Errorf("notification failure should be logged to stderr: %q", logs.String())
	}
}

func TestBridge_StdinEOFExitsCleanly(t *testing.T) {
	out, _, err := run(t, &fakeServer{t: t}, strings.NewReader(""))
	if err != nil || len(out.Lines()) != 0 {
		t.Fatalf("Run = %v, out = %q", err, out.Lines())
	}
}

func TestBridge_UnreachableServer(t *testing.T) {
	out := &lineBuffer{}
	err := Run(context.Background(), Options{URL: "http://127.0.0.1:1/mcp", Stdin: strings.NewReader(`{"jsonrpc":"2.0","id":3,"method":"ping"}` + "\n"), Stdout: out})
	if err != nil {
		t.Fatal(err)
	}
	if lines := out.Lines(); len(lines) != 1 || !strings.Contains(lines[0], `"id":3`) || !strings.Contains(lines[0], "bridge: ") {
		t.Errorf("stdout = %q", lines)
	}
}

func TestParseSSE(t *testing.T) {
	in := "retry: 10\r\n: comment\r\nid: 5\r\nevent: message\r\ndata: a\r\ndata:b\r\n\r\ndata: tail"
	var got []Event
	if err := ParseSSE(strings.NewReader(in), func(e Event) error { got = append(got, e); return nil }); err != nil {
		t.Fatal(err)
	}
	want := []Event{{Type: "message", ID: "5", Data: "a\nb"}, {Data: "tail"}}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("events = %+v, want %+v", got, want)
	}
}
