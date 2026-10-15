package embedfs

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// call runs a tool through the guarded /mcp handler and returns its text or error message.
func call(t *testing.T, s *Server, name string, args map[string]string) (string, string) {
	t.Helper()
	params, _ := json.Marshal(map[string]interface{}{"name": name, "arguments": args})
	rec := post(s.Handler(), `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":`+string(params)+`}`, nil)
	var resp struct {
		Result struct {
			Content []struct{ Text string } `json:"content"`
		} `json:"result"`
		Error *rpcError `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("%s: decode %q: %v", name, rec.Body.String(), err)
	}
	if resp.Error != nil {
		return "", resp.Error.Message
	}
	if len(resp.Result.Content) != 1 {
		t.Fatalf("%s: content = %+v", name, resp.Result.Content)
	}
	return resp.Result.Content[0].Text, ""
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustSymlink(t *testing.T, oldname, newname string) {
	t.Helper()
	if err := os.Symlink(oldname, newname); err != nil {
		t.Fatal(err)
	}
}

func assertAbsent(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		t.Errorf("%s exists (err %v), want absent", path, err)
	}
}

// sandbox builds base/root (served) and base/outside (holding secret.txt) with escape links in root.
func sandbox(t *testing.T) (*Server, string, string) {
	base := t.TempDir()
	root, outside := filepath.Join(base, "root"), filepath.Join(base, "outside")
	for _, d := range []string{root, outside, filepath.Join(root, "data")} {
		if err := os.Mkdir(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	mustWrite(t, filepath.Join(outside, "secret.txt"), "secret")
	mustWrite(t, filepath.Join(root, "data", "a.txt"), "hi")
	mustSymlink(t, outside, filepath.Join(root, "esc"))
	mustSymlink(t, filepath.Join(outside, "secret.txt"), filepath.Join(root, "escfile"))
	mustSymlink(t, filepath.Join(outside, "missing.txt"), filepath.Join(root, "dangling"))
	mustSymlink(t, "data", filepath.Join(root, "in"))
	return New(root), root, outside
}

func TestSymlinkEscapeBlocked(t *testing.T) {
	s, _, outside := sandbox(t)
	tests := []struct {
		tool string
		args map[string]string
	}{
		{"read_file", map[string]string{"path": "esc/secret.txt"}},
		{"read_file", map[string]string{"path": "escfile"}},
		{"read_file", map[string]string{"path": "../outside/secret.txt"}},
		{"write_file", map[string]string{"path": "esc/new.txt", "content": "x"}},
		{"write_file", map[string]string{"path": "esc/sub/new.txt", "content": "x"}},
		{"write_file", map[string]string{"path": "escfile", "content": "overwritten"}},
		{"write_file", map[string]string{"path": "dangling", "content": "x"}},
		{"write_file", map[string]string{"path": "../outside/new.txt", "content": "x"}},
		{"create_directory", map[string]string{"path": "esc/mkdir"}},
		{"create_directory", map[string]string{"path": "dangling/sub"}},
		{"create_directory", map[string]string{"path": "../outside/mkdir"}},
		{"list_directory", map[string]string{"path": "esc"}},
		{"get_file_info", map[string]string{"path": "esc/secret.txt"}},
	}
	for _, tt := range tests {
		t.Run(tt.tool+" "+tt.args["path"], func(t *testing.T) {
			text, errMsg := call(t, s, tt.tool, tt.args)
			if errMsg == "" || strings.Contains(text, "secret") {
				t.Fatalf("got text %q err %q, want an error", text, errMsg)
			}
		})
	}
	for _, p := range []string{"new.txt", "sub", "missing.txt", "mkdir"} {
		assertAbsent(t, filepath.Join(outside, p))
	}
	if b, _ := os.ReadFile(filepath.Join(outside, "secret.txt")); string(b) != "secret" {
		t.Errorf("secret.txt = %q, want unchanged", b)
	}
}

func TestInsideSymlinkAllowed(t *testing.T) {
	s, root, _ := sandbox(t)
	if text, errMsg := call(t, s, "read_file", map[string]string{"path": "in/a.txt"}); text != "hi" {
		t.Fatalf("read via inside link = %q, %q", text, errMsg)
	}
	if _, errMsg := call(t, s, "write_file", map[string]string{"path": "in/new/b.txt", "content": "b"}); errMsg != "" {
		t.Fatal(errMsg)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "data", "new", "b.txt")); string(b) != "b" {
		t.Errorf("data/new/b.txt = %q, want b", b)
	}
	if _, errMsg := call(t, s, "create_directory", map[string]string{"path": "in/c"}); errMsg != "" {
		t.Fatal(errMsg)
	}
	if fi, err := os.Stat(filepath.Join(root, "data", "c")); err != nil || !fi.IsDir() {
		t.Errorf("data/c not created: %v", err)
	}
}

func TestRootBehindSymlink(t *testing.T) {
	_, root, _ := sandbox(t)
	link := filepath.Join(t.TempDir(), "rootlink")
	mustSymlink(t, root, link)
	s := New(link)
	if text, errMsg := call(t, s, "read_file", map[string]string{"path": "data/a.txt"}); text != "hi" {
		t.Fatalf("read under symlinked root = %q, %q", text, errMsg)
	}
	if _, errMsg := call(t, s, "read_file", map[string]string{"path": "esc/secret.txt"}); errMsg == "" {
		t.Fatal("escape under symlinked root was allowed")
	}
}

func TestDotDotPrefixedNames(t *testing.T) {
	root := t.TempDir()
	s := New(root)
	if _, errMsg := call(t, s, "write_file", map[string]string{"path": "..foo", "content": "x"}); errMsg != "" {
		t.Fatalf("write ..foo: %s", errMsg)
	}
	if text, errMsg := call(t, s, "read_file", map[string]string{"path": "..foo"}); text != "x" {
		t.Fatalf("read ..foo = %q, %q", text, errMsg)
	}
	if _, errMsg := call(t, s, "create_directory", map[string]string{"path": "..bar/baz"}); errMsg != "" {
		t.Fatalf("mkdir ..bar/baz: %s", errMsg)
	}
	if _, err := os.Stat(filepath.Join(root, "..bar", "baz")); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"..", "../x", "a/../../x"} {
		if _, errMsg := call(t, s, "write_file", map[string]string{"path": p, "content": "x"}); !strings.Contains(errMsg, "escapes root") {
			t.Errorf("write %q err = %q, want escapes root", p, errMsg)
		}
	}
}

func TestTools(t *testing.T) {
	root := t.TempDir()
	s := New(root)
	steps := []struct {
		tool    string
		args    map[string]string
		want    string
		wantErr string
	}{
		{"write_file", map[string]string{"path": "a/b/c.go", "content": "package c"}, "Written 9 bytes to a/b/c.go", ""},
		{"write_file", map[string]string{"path": "/abs.md", "content": "# hi"}, "Written 4 bytes", ""},
		{"read_file", map[string]string{"path": "a/b/c.go"}, "package c", ""},
		{"read_file", map[string]string{"path": "nope"}, "", "read nope"},
		{"read_file", map[string]string{"path": "a/b/c.go/x"}, "", "not a directory"},
		{"create_directory", map[string]string{"path": "d/e"}, "Created directory d/e", ""},
		{"list_directory", map[string]string{"path": ""}, "a/", ""},
		{"list_directory", map[string]string{"path": "a/b"}, "c.go", ""},
		{"list_directory", map[string]string{"path": "nope"}, "", "list nope"},
		{"search_files", map[string]string{"pattern": "*.go"}, filepath.Join("a", "b", "c.go"), ""},
		{"search_files", map[string]string{"pattern": "*.rs"}, "No matches found", ""},
		{"search_files", map[string]string{"pattern": "["}, "", "invalid pattern"},
		{"get_file_info", map[string]string{"path": "a/b/c.go"}, "Type:  file", ""},
		{"get_file_info", map[string]string{"path": "d"}, "Type:  directory", ""},
		{"get_file_info", map[string]string{"path": "nope"}, "", "stat nope"},
	}
	for _, st := range steps {
		text, errMsg := call(t, s, st.tool, st.args)
		if st.wantErr != "" && !strings.Contains(errMsg, st.wantErr) {
			t.Errorf("%s %v: err = %q, want %q", st.tool, st.args, errMsg, st.wantErr)
		}
		if st.wantErr == "" && (errMsg != "" || !strings.Contains(text, st.want)) {
			t.Errorf("%s %v = %q, %q; want %q", st.tool, st.args, text, errMsg, st.want)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "abs.md")); err != nil {
		t.Errorf("absolute path not mapped under root: %v", err)
	}
}

func TestToolsInvalidArgs(t *testing.T) {
	for _, name := range []string{"read_file", "write_file", "list_directory", "search_files", "get_file_info", "create_directory"} {
		if _, err := toolHandler(name)(t.TempDir(), json.RawMessage(`[]`)); err == nil || !strings.Contains(err.Error(), "invalid args") {
			t.Errorf("%s: err = %v, want invalid args", name, err)
		}
	}
}
