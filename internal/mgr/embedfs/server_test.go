package embedfs

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// post sends body to /mcp through the guarded handler as a local MCP client would.
func post(h http.Handler, body string, hdr map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
	req.Host = "127.0.0.1:5901"
	req.Header.Set("Content-Type", "application/json")
	for k, v := range hdr {
		if k == "Host" {
			req.Host = v
		} else if v == "" {
			req.Header.Del(k)
		} else {
			req.Header.Set(k, v)
		}
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestHandleMCP(t *testing.T) {
	h := New(t.TempDir()).Handler()
	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantBody   string
	}{
		{"notification is 202", http.MethodPost, `{"jsonrpc":"2.0","method":"notifications/initialized"}`, http.StatusAccepted, ""},
		{"other notification is 202", http.MethodPost, `{"jsonrpc":"2.0","method":"notifications/cancelled","params":{}}`, http.StatusAccepted, ""},
		{"client response is 202", http.MethodPost, `{"jsonrpc":"2.0","id":9,"result":{}}`, http.StatusAccepted, ""},
		{"initialize", http.MethodPost, `{"jsonrpc":"2.0","id":1,"method":"initialize"}`, http.StatusOK, `"protocolVersion"`},
		{"tools/list", http.MethodPost, `{"jsonrpc":"2.0","id":"a","method":"tools/list"}`, http.StatusOK, `"id":"a"`},
		{"unknown method", http.MethodPost, `{"jsonrpc":"2.0","id":2,"method":"nope"}`, http.StatusOK, `-32601`},
		{"parse error", http.MethodPost, `{not json`, http.StatusOK, `-32700`},
		{"invalid params", http.MethodPost, `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":[]}`, http.StatusOK, `-32602`},
		{"unknown tool", http.MethodPost, `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"rm_rf"}}`, http.StatusOK, `Unknown tool`},
		{"GET not allowed", http.MethodGet, "", http.StatusMethodNotAllowed, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/mcp", strings.NewReader(tt.body))
			req.Host = "localhost:5901"
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantBody == "" && rec.Body.Len() != 0 {
				t.Errorf("body = %q, want empty", rec.Body.String())
			}
			if tt.wantBody != "" && !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Errorf("body = %q, want %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestGuard(t *testing.T) {
	init := `{"jsonrpc":"2.0","id":1,"method":"initialize"}`
	tests := []struct {
		name string
		hdr  map[string]string
		want int
	}{
		{"localhost JSON", map[string]string{"Host": "localhost:5901"}, http.StatusOK},
		{"127.0.0.1 no port", map[string]string{"Host": "127.0.0.1"}, http.StatusOK},
		{"ipv6 loopback", map[string]string{"Host": "[::1]:5901"}, http.StatusOK},
		{"ipv6 loopback no port", map[string]string{"Host": "[::1]"}, http.StatusOK},
		{"uppercase localhost", map[string]string{"Host": "LOCALHOST:5901"}, http.StatusOK},
		{"json with charset", map[string]string{"Content-Type": "application/json; charset=utf-8"}, http.StatusOK},
		{"local origin", map[string]string{"Origin": "http://localhost:3000"}, http.StatusOK},
		{"local https origin", map[string]string{"Origin": "https://127.0.0.1"}, http.StatusOK},
		{"ipv6 origin", map[string]string{"Origin": "http://[::1]:8080"}, http.StatusOK},
		{"rebound host", map[string]string{"Host": "evil.example:5901"}, http.StatusForbidden},
		{"rebound host no port", map[string]string{"Host": "evil.example"}, http.StatusForbidden},
		{"lookalike host", map[string]string{"Host": "127.0.0.1.nip.io:5901"}, http.StatusForbidden},
		{"lan host", map[string]string{"Host": "192.168.1.5:5901"}, http.StatusForbidden},
		{"cross-site text/plain", map[string]string{"Origin": "https://evil.example", "Content-Type": "text/plain"}, http.StatusForbidden},
		{"cross-site json", map[string]string{"Origin": "https://evil.example"}, http.StatusForbidden},
		{"lookalike origin", map[string]string{"Origin": "http://localhost.evil.example"}, http.StatusForbidden},
		{"null origin", map[string]string{"Origin": "null"}, http.StatusForbidden},
		{"file origin", map[string]string{"Origin": "file://"}, http.StatusForbidden},
		{"text/plain", map[string]string{"Content-Type": "text/plain"}, http.StatusUnsupportedMediaType},
		{"form", map[string]string{"Content-Type": "application/x-www-form-urlencoded"}, http.StatusUnsupportedMediaType},
		{"no content type", map[string]string{"Content-Type": ""}, http.StatusUnsupportedMediaType},
		{"bad content type", map[string]string{"Content-Type": ";;"}, http.StatusUnsupportedMediaType},
	}
	h := New(t.TempDir()).Handler()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if rec := post(h, init, tt.hdr); rec.Code != tt.want {
				t.Errorf("status = %d, want %d (%s)", rec.Code, tt.want, rec.Body.String())
			}
		})
	}
}

func TestGuardBlocksCrossSiteWrite(t *testing.T) {
	root := t.TempDir()
	h := New(root).Handler()
	body := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"write_file","arguments":{"path":"pwned","content":"x"}}}`
	for _, hdr := range []map[string]string{
		{"Origin": "https://evil.example", "Content-Type": "text/plain"},
		{"Host": "evil.example:5901"},
	} {
		if rec := post(h, body, hdr); rec.Code == http.StatusOK {
			t.Errorf("%v: status 200, want rejection", hdr)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "pwned")); !os.IsNotExist(err) {
		t.Fatalf("cross-site request wrote a file (stat err %v)", err)
	}
	if rec := post(h, body, nil); !strings.Contains(rec.Body.String(), "Written 1 bytes") {
		t.Fatalf("local write failed: %s", rec.Body.String())
	}
}

func TestHealth(t *testing.T) {
	h := New(t.TempDir()).Handler()
	for host, want := range map[string]int{"localhost:5901": http.StatusOK, "evil.example:5901": http.StatusForbidden} {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		req.Host = host
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != want {
			t.Errorf("%s: status = %d, want %d", host, rec.Code, want)
		}
	}
}

func TestMissingRoot(t *testing.T) {
	s := New(filepath.Join(t.TempDir(), "missing"))
	if err := s.Serve(0, ""); err == nil {
		t.Fatal("Serve with a missing root succeeded")
	}
	rec := post(s.Handler(), `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"read_file","arguments":{"path":"x"}}}`, nil)
	var resp jsonRPCResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil || resp.Error == nil {
		t.Fatalf("tools/call with a missing root = %s, want error", rec.Body.String())
	}
}

func TestIsWithin(t *testing.T) {
	root := filepath.FromSlash("/srv/root")
	tests := map[string]bool{
		"/srv/root":           true,
		"/srv/root/a/b":       true,
		"/srv/root/..foo":     true,
		"/srv/root/..foo/bar": true,
		"/srv/root/a..":       true,
		"/srv":                false,
		"/srv/other":          false,
		"/srv/rootx":          false,
		"/":                   false,
	}
	for target, want := range tests {
		if got := isWithin(filepath.FromSlash(target), root); got != want {
			t.Errorf("isWithin(%q) = %v, want %v", target, got, want)
		}
	}
}

func TestLogToPrivateAndNoFollow(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "embed.log")
	LogTo(logPath, "hello")
	fi, err := os.Stat(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Fatalf("mode = %v, want 0600", fi.Mode().Perm())
	}
	victim := filepath.Join(dir, "victim")
	if err := os.WriteFile(victim, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "link.log")
	if err := os.Symlink(victim, link); err != nil {
		t.Fatal(err)
	}
	LogTo(link, "pwned")
	if b, _ := os.ReadFile(victim); string(b) != "keep" {
		t.Fatalf("LogTo followed a symlink: victim = %q", b)
	}
}
