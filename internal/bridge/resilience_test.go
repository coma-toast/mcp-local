package bridge

import (
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

// mcpServer is a restartable Streamable HTTP MCP server: restart forgets every session and
// ends open streams, as a real server restart does.
type mcpServer struct {
	mu          sync.Mutex
	sessions    map[string]bool
	created     int
	initHeaders []http.Header
	initBodies  []string
	initialized int
	posts       []rpcPost // non-initialize POSTs that reached a live session
	gets        []getRec
	deletes     []string
	stop        chan struct{}
	hold404     int // when > 0, 404s wait until this many are pending (forces concurrency)
	pending404  int
	release404  chan struct{}
	getFn       func(w http.ResponseWriter, r *http.Request, n int, stop <-chan struct{}) // nil = 405
	methods     map[string]func(w http.ResponseWriter, r *http.Request, id json.RawMessage)
}

type rpcPost struct {
	Method, Session, Body string
	ID                    json.RawMessage
}

type getRec struct {
	Header     http.Header
	Start, End time.Time
}

func newMCPServer() *mcpServer {
	return &mcpServer{sessions: map[string]bool{}, stop: make(chan struct{}), release404: make(chan struct{}), methods: map[string]func(http.ResponseWriter, *http.Request, json.RawMessage){}}
}

func (s *mcpServer) restart() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions = map[string]bool{}
	close(s.stop)
	s.stop = make(chan struct{})
}

func (s *mcpServer) snapshot() (inits, initialized int, posts []rpcPost, gets []getRec) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.initBodies), s.initialized, append([]rpcPost(nil), s.posts...), append([]getRec(nil), s.gets...)
}

func (s *mcpServer) notFound(w http.ResponseWriter) {
	s.mu.Lock()
	if s.hold404 > 0 {
		if s.pending404++; s.pending404 == s.hold404 {
			close(s.release404)
		}
		ch := s.release404
		s.mu.Unlock()
		select {
		case <-ch:
		case <-time.After(2 * time.Second):
		}
	} else {
		s.mu.Unlock()
	}
	http.Error(w, "session not found", http.StatusNotFound)
}

func (s *mcpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sid := r.Header.Get(HeaderSessionID)
	s.mu.Lock()
	live, stop := s.sessions[sid], s.stop
	s.mu.Unlock()
	switch r.Method {
	case http.MethodDelete:
		s.mu.Lock()
		s.deletes = append(s.deletes, sid)
		s.mu.Unlock()
		return
	case http.MethodGet:
		if !live {
			s.notFound(w)
			return
		}
		s.mu.Lock()
		n := len(s.gets)
		s.gets = append(s.gets, getRec{Header: r.Header.Clone(), Start: time.Now()})
		s.mu.Unlock()
		defer func() { s.mu.Lock(); s.gets[n].End = time.Now(); s.mu.Unlock() }()
		if s.getFn == nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		s.getFn(w, r, n, stop)
		return
	}
	body, _ := io.ReadAll(r.Body)
	var req struct {
		ID     json.RawMessage `json:"id"`
		Method string          `json:"method"`
	}
	_ = json.Unmarshal(body, &req)
	if req.Method == "initialize" {
		s.mu.Lock()
		s.created++
		sid = fmt.Sprintf("sess-%d", s.created)
		s.sessions[sid] = true
		s.initHeaders = append(s.initHeaders, r.Header.Clone())
		s.initBodies = append(s.initBodies, string(body))
		s.mu.Unlock()
		w.Header().Set(HeaderSessionID, sid)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"result":{"protocolVersion":"2025-06-18","capabilities":{}}}`, req.ID)
		return
	}
	if !live {
		s.notFound(w)
		return
	}
	s.mu.Lock()
	s.posts = append(s.posts, rpcPost{Method: req.Method, Session: sid, Body: string(body), ID: req.ID})
	if req.Method == "notifications/initialized" {
		s.initialized++
	}
	fn := s.methods[req.Method]
	s.mu.Unlock()
	switch {
	case fn != nil:
		fn(w, r, req.ID)
	case req.Method == "" || strings.HasPrefix(req.Method, "notifications/"):
		w.WriteHeader(http.StatusAccepted)
	default:
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"result":{}}`, req.ID)
	}
}

type harness struct {
	t     *testing.T
	in    *io.PipeWriter
	out   *lineBuffer
	logs  *lineBuffer
	done  chan error
	once  sync.Once
	err   error
	start time.Time
}

func startBridge(t *testing.T, h http.Handler, mod func(*Options)) *harness {
	t.Helper()
	ts := httptest.NewServer(h)
	t.Cleanup(ts.Close)
	pr, pw := io.Pipe()
	hs := &harness{t: t, in: pw, out: &lineBuffer{}, logs: &lineBuffer{}, done: make(chan error, 1), start: time.Now()}
	o := Options{URL: ts.URL, Stdin: pr, Stdout: hs.out, Log: log.New(hs.logs, "", 0), Timeout: 5 * time.Second,
		ReconnectDelay: 10 * time.Millisecond, MaxReconnectDelay: 100 * time.Millisecond, ShutdownGrace: 200 * time.Millisecond}
	if mod != nil {
		mod(&o)
	}
	go func() { hs.done <- Run(context.Background(), o) }()
	t.Cleanup(func() { _ = hs.close() })
	return hs
}

func (h *harness) send(lines ...string) {
	h.t.Helper()
	for _, l := range lines {
		if _, err := fmt.Fprintln(h.in, l); err != nil {
			h.t.Fatalf("write stdin: %v", err)
		}
	}
}

// close ends stdin and waits for Run to return.
func (h *harness) close() error {
	h.once.Do(func() {
		h.in.Close()
		select {
		case h.err = <-h.done:
		case <-time.After(10 * time.Second):
			h.err = fmt.Errorf("Run did not exit after stdin EOF")
		}
	})
	return h.err
}

func (h *harness) count(substr string) int {
	n := 0
	for _, l := range h.out.Lines() {
		if strings.Contains(l, substr) {
			n++
		}
	}
	return n
}

type rpcResp struct {
	ID     json.RawMessage `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (h *harness) response(id string) rpcResp {
	h.t.Helper()
	h.out.waitFor(h.t, `"id":`+id)
	for _, l := range h.out.Lines() {
		var r rpcResp
		if json.Unmarshal([]byte(l), &r) == nil && string(r.ID) == id {
			return r
		}
	}
	h.t.Fatalf("no response for id %s in %q", id, h.out.Lines())
	return rpcResp{}
}

func waitUntil(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func ping(id int) string { return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"method":"ping"}`, id) }

func TestBridge_SessionExpiryReinitializesTransparently(t *testing.T) {
	srv := newMCPServer()
	h := startBridge(t, srv, nil)
	h.send(initMsg, initializedMsg, ping(2))
	h.response("2")
	srv.restart()
	h.send(ping(3))
	if r := h.response("3"); r.Error != nil {
		t.Fatalf("ping after restart: %+v", r.Error)
	}
	if err := h.close(); err != nil {
		t.Fatal(err)
	}
	if n := h.count(`"id":1`); n != 1 {
		t.Errorf("the replayed initialize response must not reach stdout; id 1 lines = %d: %q", n, h.out.Lines())
	}
	inits, initialized, posts, _ := srv.snapshot()
	if inits != 2 || initialized != 2 {
		t.Fatalf("inits=%d initialized=%d, want 2 and 2", inits, initialized)
	}
	if srv.initBodies[1] != initMsg {
		t.Errorf("replayed initialize = %s, want the client's own %s", srv.initBodies[1], initMsg)
	}
	for i, hd := range srv.initHeaders {
		if hd.Get(HeaderSessionID) != "" || hd.Get(HeaderProtocolVersion) != "" {
			t.Errorf("initialize %d carried session headers: %v", i, hd)
		}
	}
	if last := posts[len(posts)-1]; last.Method != "ping" || last.Session != "sess-2" {
		t.Errorf("retried ping = %+v, want it on sess-2", last)
	}
}

func TestBridge_ConcurrentSessionExpiryReinitializesOnce(t *testing.T) {
	srv := newMCPServer()
	h := startBridge(t, srv, nil)
	h.send(initMsg, initializedMsg)
	waitUntil(t, "initialized", func() bool { _, n, _, _ := srv.snapshot(); return n == 1 })
	srv.restart()
	const n = 8
	srv.mu.Lock()
	srv.hold404 = n
	srv.mu.Unlock()
	for i := range n {
		h.send(ping(100 + i))
	}
	for i := range n {
		if r := h.response(fmt.Sprint(100 + i)); r.Error != nil {
			t.Errorf("ping %d: %+v", 100+i, r.Error)
		}
	}
	if inits, _, _, _ := srv.snapshot(); inits != 2 {
		t.Errorf("initialize sent %d times, want 2 (one replay for %d concurrent 404s)", inits, n)
	}
}

func TestBridge_ReinitFailureRetriesOnNextRequest(t *testing.T) {
	srv := newMCPServer()
	var mu sync.Mutex
	initDown := false
	h := startBridge(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(strings.NewReader(string(body)))
		mu.Lock()
		down := initDown && strings.Contains(string(body), `"initialize"`)
		mu.Unlock()
		if down {
			http.Error(w, "starting up", http.StatusServiceUnavailable)
			return
		}
		srv.ServeHTTP(w, r)
	}), nil)
	h.send(initMsg, initializedMsg)
	waitUntil(t, "initialized", func() bool { _, n, _, _ := srv.snapshot(); return n == 1 })
	srv.restart()
	mu.Lock()
	initDown = true
	mu.Unlock()
	h.send(`{"jsonrpc":"2.0","id":"a","method":"ping"}`)
	if r := h.response(`"a"`); r.Error == nil || !strings.Contains(r.Error.Message, "re-initialize") || !strings.Contains(r.Error.Message, "503") {
		t.Fatalf("ping while the server cannot initialize = %+v", r)
	}
	mu.Lock()
	initDown = false
	mu.Unlock()
	h.send(`{"jsonrpc":"2.0","id":"b","method":"ping"}`)
	if r := h.response(`"b"`); r.Error != nil {
		t.Fatalf("ping after the server recovered: %+v", r.Error)
	}
	if _, _, posts, _ := srv.snapshot(); posts[len(posts)-1].Session != "sess-2" {
		t.Errorf("last POST = %+v, want sess-2", posts[len(posts)-1])
	}
}

func TestBridge_StreamEndedWithoutResponse(t *testing.T) {
	srv := newMCPServer()
	srv.methods["sse-silent"] = func(w http.ResponseWriter, r *http.Request, id json.RawMessage) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "data: {\"jsonrpc\":\"2.0\",\"method\":\"notifications/progress\",\"params\":{}}\n\n")
	}
	srv.methods["json-empty"] = func(w http.ResponseWriter, r *http.Request, id json.RawMessage) {
		w.Header().Set("Content-Type", "application/json")
	}
	h := startBridge(t, srv, nil)
	h.send(initMsg, initializedMsg, `{"jsonrpc":"2.0","id":5,"method":"sse-silent"}`, `{"jsonrpc":"2.0","id":6,"method":"json-empty"}`)
	for _, id := range []string{"5", "6"} {
		r := h.response(id)
		if r.Error == nil || r.Error.Code != -32000 || r.Error.Message != "bridge: stream ended without response" {
			t.Errorf("id %s: %+v", id, r)
		}
	}
}

func TestBridge_TimeoutIsIdleNotTotal(t *testing.T) {
	srv := newMCPServer()
	srv.methods["slow"] = func(w http.ResponseWriter, r *http.Request, id json.RawMessage) {
		w.Header().Set("Content-Type", "text/event-stream")
		for range 8 {
			fmt.Fprint(w, "data: {\"jsonrpc\":\"2.0\",\"method\":\"notifications/progress\",\"params\":{}}\n\n")
			w.(http.Flusher).Flush()
			time.Sleep(50 * time.Millisecond)
		}
		fmt.Fprintf(w, "data: {\"jsonrpc\":\"2.0\",\"id\":%s,\"result\":{\"done\":true}}\n\n", id)
	}
	srv.methods["stuck"] = func(w http.ResponseWriter, r *http.Request, id json.RawMessage) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}
	srv.methods["no-headers"] = func(w http.ResponseWriter, r *http.Request, id json.RawMessage) { <-r.Context().Done() }
	h := startBridge(t, srv, func(o *Options) { o.Timeout = 150 * time.Millisecond })
	h.send(initMsg, initializedMsg, `{"jsonrpc":"2.0","id":7,"method":"slow"}`, `{"jsonrpc":"2.0","id":8,"method":"stuck"}`, `{"jsonrpc":"2.0","id":9,"method":"no-headers"}`)
	if r := h.response("7"); r.Error != nil || string(r.Result) != `{"done":true}` {
		t.Errorf("slow call streaming progress for 400ms with a 150ms timeout: %+v", r)
	}
	for _, id := range []string{"8", "9"} {
		if r := h.response(id); r.Error == nil || !strings.Contains(r.Error.Message, "timed out") {
			t.Errorf("id %s: %+v", id, r)
		}
	}
	waitUntil(t, "notifications/cancelled for 8 and 9", func() bool {
		_, _, posts, _ := srv.snapshot()
		got := map[string]bool{}
		for _, p := range posts {
			var n struct {
				Params struct {
					RequestID json.RawMessage `json:"requestId"`
					Reason    string          `json:"reason"`
				} `json:"params"`
			}
			if p.Method == "notifications/cancelled" && json.Unmarshal([]byte(p.Body), &n) == nil && strings.Contains(n.Params.Reason, "timed out") {
				got[string(n.Params.RequestID)] = true
			}
		}
		return got["8"] && got["9"]
	})
}

func TestBridge_ClientCancelledRequestGetsNoSyntheticError(t *testing.T) {
	srv := newMCPServer()
	cancelled := make(chan struct{})
	srv.methods["hang"] = func(w http.ResponseWriter, r *http.Request, id json.RawMessage) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.(http.Flusher).Flush()
		<-cancelled // server honours the cancel by closing the stream without a response
	}
	srv.methods["notifications/cancelled"] = func(w http.ResponseWriter, r *http.Request, id json.RawMessage) {
		close(cancelled)
		w.WriteHeader(http.StatusAccepted)
	}
	h := startBridge(t, srv, nil)
	h.send(initMsg, initializedMsg, `{"jsonrpc":"2.0","id":11,"method":"hang"}`)
	waitUntil(t, "hang to reach the server", func() bool { _, _, p, _ := srv.snapshot(); return len(p) == 2 })
	h.send(`{"jsonrpc":"2.0","method":"notifications/cancelled","params":{"requestId":11}}`, ping(12))
	h.response("12")
	if err := h.close(); err != nil {
		t.Fatal(err)
	}
	if n := h.count(`"id":11`); n != 0 {
		t.Errorf("client-cancelled request got %d output lines: %q", n, h.out.Lines())
	}
}

func TestBridge_GETStreamReconnects(t *testing.T) {
	srv := newMCPServer()
	srv.getFn = func(w http.ResponseWriter, r *http.Request, n int, stop <-chan struct{}) {
		switch {
		case n < 4:
			http.Error(w, "busy", http.StatusServiceUnavailable)
		case n == 4:
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close()
		case n == 5:
			w.Header().Set("Content-Type", "text/event-stream")
			fmt.Fprint(w, "retry: 300\nid: ev-1\ndata: {\"jsonrpc\":\"2.0\",\"method\":\"notifications/tools/list_changed\"}\n\n")
		default:
			w.Header().Set("Content-Type", "text/event-stream")
			fmt.Fprint(w, "id: ev-2\ndata: {\"jsonrpc\":\"2.0\",\"method\":\"notifications/resources/list_changed\"}\n\n")
			w.(http.Flusher).Flush()
			<-r.Context().Done()
		}
	}
	h := startBridge(t, srv, func(o *Options) { o.ReconnectDelay, o.MaxReconnectDelay = 10*time.Millisecond, 40*time.Millisecond })
	h.send(initMsg)
	h.out.waitFor(t, "notifications/tools/list_changed")
	h.out.waitFor(t, "notifications/resources/list_changed")
	if err := h.close(); err != nil {
		t.Fatal(err)
	}
	_, _, _, gets := srv.snapshot()
	if got := gets[6].Header.Get(HeaderLastEventID); got != "ev-1" {
		t.Errorf("Last-Event-ID on reconnect = %q, want ev-1", got)
	}
	if gap := gets[6].Start.Sub(gets[5].End); gap < 250*time.Millisecond {
		t.Errorf("reconnect after retry: 300 took %s", gap)
	}
	if gap := gets[3].Start.Sub(gets[2].End); gap < 15*time.Millisecond {
		t.Errorf("third retry after %s; backoff should have grown past the 10ms base", gap)
	}
	logs := 0
	for _, l := range h.logs.Lines() {
		if strings.Contains(l, "reconnecting") {
			logs++
		}
	}
	if logs != 1 {
		t.Errorf("reconnect logging should be rate-limited to 1 line, got %d: %q", logs, h.logs.Lines())
	}
}

func TestBridge_GETStreamNotOfferedStopsAfterOneTry(t *testing.T) {
	srv := newMCPServer()
	h := startBridge(t, srv, nil)
	h.send(initMsg, ping(2))
	h.response("2")
	time.Sleep(100 * time.Millisecond)
	if _, _, _, gets := srv.snapshot(); len(gets) != 1 {
		t.Errorf("GET attempts after 405 = %d, want 1", len(gets))
	}
}

func TestBridge_GETStreamSessionExpiryReinitializes(t *testing.T) {
	srv := newMCPServer()
	srv.getFn = func(w http.ResponseWriter, r *http.Request, n int, stop <-chan struct{}) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprintf(w, "data: {\"jsonrpc\":\"2.0\",\"method\":\"notifications/message\",\"params\":{\"n\":%d}}\n\n", n)
		w.(http.Flusher).Flush()
		select {
		case <-stop:
		case <-r.Context().Done():
		}
	}
	h := startBridge(t, srv, nil)
	h.send(initMsg, initializedMsg)
	h.out.waitFor(t, `"n":0`)
	srv.restart()
	h.out.waitFor(t, `"n":1`)
	if err := h.close(); err != nil {
		t.Fatal(err)
	}
	inits, initialized, _, gets := srv.snapshot()
	if inits != 2 || initialized != 2 || gets[1].Header.Get(HeaderSessionID) != "sess-2" {
		t.Errorf("inits=%d initialized=%d second GET session=%q", inits, initialized, gets[1].Header.Get(HeaderSessionID))
	}
	if n := h.count(`"id":1`); n != 1 {
		t.Errorf("id 1 lines = %d, want only the client's own initialize response", n)
	}
}

func TestBridge_ShutdownCancelsInFlightAndDeletes(t *testing.T) {
	srv := newMCPServer()
	srv.methods["hang"] = func(w http.ResponseWriter, r *http.Request, id json.RawMessage) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}
	h := startBridge(t, srv, func(o *Options) { o.Timeout = time.Minute })
	h.send(initMsg, initializedMsg, `{"jsonrpc":"2.0","id":5,"method":"hang"}`)
	waitUntil(t, "hang to reach the server", func() bool { _, _, p, _ := srv.snapshot(); return len(p) == 2 })
	start := time.Now()
	if err := h.close(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if d := time.Since(start); d > 3*time.Second {
		t.Errorf("shutdown took %s", d)
	}
	if r := h.response("5"); r.Error == nil || !strings.Contains(r.Error.Message, "shutting down") {
		t.Errorf("in-flight request at shutdown: %+v", r)
	}
	_, _, posts, _ := srv.snapshot()
	if last := posts[len(posts)-1]; last.Method != "notifications/cancelled" || !strings.Contains(last.Body, `"requestId":5`) {
		t.Errorf("last POST = %+v, want notifications/cancelled for 5", last)
	}
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if len(srv.deletes) != 1 || srv.deletes[0] != "sess-1" {
		t.Errorf("DELETEs = %v", srv.deletes)
	}
}

func TestBridge_MalformedInput(t *testing.T) {
	srv := newMCPServer()
	h := startBridge(t, srv, nil)
	h.send(`{"jsonrpc":"2.0",`,
		`[{"jsonrpc":"2.0","id":7,"method":"ping"},{"jsonrpc":"2.0","method":"notifications/x"},{"jsonrpc":"2.0","id":"b","method":"ping"}]`,
		`[]`, `42`)
	if err := h.close(); err != nil {
		t.Fatal(err)
	}
	want := []string{
		`{"error":{"code":-32700,"message":"Parse error"},"id":null,"jsonrpc":"2.0"}`,
		`{"error":{"code":-32600,"message":"Invalid Request: JSON-RPC batches are not supported"},"id":7,"jsonrpc":"2.0"}`,
		`{"error":{"code":-32600,"message":"Invalid Request: JSON-RPC batches are not supported"},"id":"b","jsonrpc":"2.0"}`,
		`{"error":{"code":-32600,"message":"Invalid Request: empty batch"},"id":null,"jsonrpc":"2.0"}`,
		`{"error":{"code":-32600,"message":"Invalid Request"},"id":null,"jsonrpc":"2.0"}`,
	}
	if got := h.out.Lines(); strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("stdout:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if inits, _, posts, _ := srv.snapshot(); inits+len(posts) != 0 {
		t.Errorf("malformed input reached the server: %d inits, %v", inits, posts)
	}
}

func TestBridge_NonJSONServerDataDropped(t *testing.T) {
	srv := newMCPServer()
	srv.methods["garbage"] = func(w http.ResponseWriter, r *http.Request, id json.RawMessage) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprintf(w, "data: not json\n\ndata: {\"jsonrpc\":\"2.0\",\"id\":%s,\"result\":{}}\n\n", id)
	}
	srv.methods["garbage-json"] = func(w http.ResponseWriter, r *http.Request, id json.RawMessage) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "<html>oops</html>")
	}
	h := startBridge(t, srv, nil)
	h.send(initMsg, `{"jsonrpc":"2.0","id":4,"method":"garbage"}`, `{"jsonrpc":"2.0","id":5,"method":"garbage-json"}`)
	h.response("5")
	if err := h.close(); err != nil {
		t.Fatal(err)
	}
	for _, l := range h.out.Lines() {
		if !json.Valid([]byte(l)) {
			t.Errorf("non-JSON line on stdout: %q", l)
		}
	}
	if r := h.response("4"); r.Error != nil {
		t.Errorf("id 4: %+v", r)
	}
	if !strings.Contains(strings.Join(h.logs.Lines(), "\n"), "dropped non-JSON data from server: not json") {
		t.Errorf("logs = %q", h.logs.Lines())
	}
}

func TestBridge_ServerToClientRequestRoundTrip(t *testing.T) {
	srv := newMCPServer()
	srv.getFn = func(w http.ResponseWriter, r *http.Request, n int, stop <-chan struct{}) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "data: {\"jsonrpc\":\"2.0\",\"id\":\"srv-1\",\"method\":\"roots/list\"}\n\n")
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}
	h := startBridge(t, srv, nil)
	h.send(initMsg, initializedMsg)
	h.out.waitFor(t, `"method":"roots/list"`)
	reply := `{"jsonrpc":"2.0","id":"srv-1","result":{"roots":[{"uri":"file:///tmp"}]}}`
	h.send(reply)
	waitUntil(t, "the client's reply to reach the server", func() bool {
		_, _, posts, _ := srv.snapshot()
		for _, p := range posts {
			if p.Body == reply && p.Session == "sess-1" {
				return true
			}
		}
		return false
	})
	if err := h.close(); err != nil {
		t.Fatal(err)
	}
	if n := len(h.out.Lines()); n != 2 {
		t.Errorf("stdout = %q, want the init response and the server request only", h.out.Lines())
	}
}

func TestBridge_MultiMegabyteMessages(t *testing.T) {
	srv := newMCPServer()
	srv.methods["echo"] = func(w http.ResponseWriter, r *http.Request, id json.RawMessage) {
		w.Header().Set("Content-Type", "text/event-stream")
		srv.mu.Lock()
		body := srv.posts[len(srv.posts)-1].Body
		srv.mu.Unlock()
		var in struct {
			Params json.RawMessage `json:"params"`
		}
		_ = json.Unmarshal([]byte(body), &in)
		fmt.Fprintf(w, "data: {\"jsonrpc\":\"2.0\",\"id\":%s,\"result\":%s}\n\n", id, in.Params)
	}
	h := startBridge(t, srv, nil)
	big := strings.Repeat("x", 6<<20)
	h.send(initMsg, fmt.Sprintf(`{"jsonrpc":"2.0","id":3,"method":"echo","params":{"blob":%q}}`, big))
	r := h.response("3")
	var res struct{ Blob string }
	if err := json.Unmarshal(r.Result, &res); err != nil || res.Blob != big {
		t.Errorf("echoed blob: len %d, err %v", len(res.Blob), err)
	}
}
