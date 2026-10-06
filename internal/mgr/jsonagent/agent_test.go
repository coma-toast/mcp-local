package jsonagent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

func TestMain(m *testing.M) {
	home, err := os.MkdirTemp("", "jsonagent-home-*")
	if err != nil {
		panic(err)
	}
	os.Setenv("HOME", home)
	code := m.Run()
	os.RemoveAll(home)
	os.Exit(code)
}

func remoteConverter(e config.AgentEntry) map[string]interface{} {
	if e.Type == "remote" {
		return map[string]interface{}{"type": "http", "url": e.URL}
	}
	return map[string]interface{}{"type": "stdio", "command": e.Command[0]}
}

func writeFile(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func block(t *testing.T, path, key string) map[string]interface{} {
	t.Helper()
	m, err := ReadJSON(path)
	if err != nil {
		t.Fatalf("ReadJSON: %v", err)
	}
	b, _ := m[key].(map[string]interface{})
	return b
}

func TestReadJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.json")
	writeFile(t, path, `{"key": "value", "nested": {"inner": 123}}`, 0o644)
	got, err := ReadJSON(path)
	if err != nil {
		t.Fatal(err)
	}
	if got["key"] != "value" || got["nested"].(map[string]interface{})["inner"] != float64(123) {
		t.Errorf("got %v", got)
	}
}

func TestReadJSON_JSONC(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.jsonc")
	writeFile(t, path, "{\n  // line\n  /* block */ \"key\": \"value\",\n}\n", 0o644)
	got, err := ReadJSON(path)
	if err != nil {
		t.Fatal(err)
	}
	if got["key"] != "value" {
		t.Errorf("key = %v", got["key"])
	}
}

func TestReadJSON_NotExist(t *testing.T) {
	if _, err := ReadJSON("/nonexistent/path.json"); err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestRegisterServices_NewFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "config.json")
	a := New(path, "mcpServers", remoteConverter)
	services := []config.ServiceConfig{
		{Name: "test-service", Command: "/bin/echo", Args: []string{"hello"}, MCPType: "stdio"},
		{Name: "remote-service", MCPType: "http", MCPURL: "http://localhost:8080/mcp"},
	}
	if err := a.RegisterServices(services); err != nil {
		t.Fatal(err)
	}
	want := `{
  "mcpServers": {
    "remote-service": {
      "type": "http",
      "url": "http://localhost:8080/mcp"
    },
    "test-service": {
      "command": "/bin/echo",
      "type": "stdio"
    }
  }
}
`
	if got := readFile(t, path); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

const jsoncFixture = `// OpenCode config — user comments must survive
{
	"$schema": "https://opencode.ai/config.json",
	/* theme block
	   spans lines */
	"theme": "dark", // inline after value
	"mcp": {
		// foreign server added by hand
		"foreign": { "type": "local", "command": ["foo"], "enabled": true, }, // keep me
	},
	"tail": [1, 2, 3,],
}
`

func TestSetEntries_PreservesJSONC(t *testing.T) {
	path := filepath.Join(t.TempDir(), "opencode.jsonc")
	writeFile(t, path, jsoncFixture, 0o644)
	a := New(path, "mcp", nil)
	if err := a.SetEntries(map[string]map[string]interface{}{"ours": {"type": "remote", "url": "http://x/mcp"}}); err != nil {
		t.Fatal(err)
	}
	got := readFile(t, path)
	want := `// OpenCode config — user comments must survive
{
	"$schema": "https://opencode.ai/config.json",
	/* theme block
	   spans lines */
	"theme": "dark", // inline after value
	"mcp": {
		// foreign server added by hand
		"foreign": { "type": "local", "command": ["foo"], "enabled": true, }, // keep me
		"ours": {
			"type": "remote",
			"url": "http://x/mcp"
		},
	},
	"tail": [1, 2, 3,],
}
`
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
	if b := block(t, path, "mcp"); b["foreign"] == nil || b["ours"] == nil {
		t.Errorf("block = %v", b)
	}
	// Deregistering our entry restores the original bytes exactly.
	found, err := a.Deregister("ours")
	if err != nil || !found {
		t.Fatalf("Deregister = %v, %v", found, err)
	}
	if got := readFile(t, path); got != jsoncFixture {
		t.Errorf("after deregister got:\n%s\nwant original:\n%s", got, jsoncFixture)
	}
}

func TestSetEntries_ReplaceKeepsSurroundings(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mcp.json")
	orig := "{\n  // top\n  \"mcpServers\": {\n    \"a\": {\"url\": \"http://old\"}, // trailing\n    \"b\": {\"url\": \"http://b\"}\n  }\n}\n"
	writeFile(t, path, orig, 0o644)
	a := New(path, "mcpServers", nil)
	if err := a.SetEntries(map[string]map[string]interface{}{"a": {"url": "http://new"}}); err != nil {
		t.Fatal(err)
	}
	got := readFile(t, path)
	if !strings.Contains(got, "// top") || !strings.Contains(got, "// trailing") || !strings.Contains(got, `"b": {"url": "http://b"}`) {
		t.Errorf("surroundings lost:\n%s", got)
	}
	if block(t, path, "mcpServers")["a"].(map[string]interface{})["url"] != "http://new" {
		t.Errorf("a not replaced:\n%s", got)
	}
}

func TestSetEntries_IdempotentNoWrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mcp.json")
	a := New(path, "mcpServers", nil)
	entry := map[string]map[string]interface{}{"svc": {"url": "http://x", "args": []string{"a"}}}
	if err := a.SetEntries(entry); err != nil {
		t.Fatal(err)
	}
	fi1, _ := os.Stat(path)
	time.Sleep(10 * time.Millisecond)
	if err := a.SetEntries(entry); err != nil {
		t.Fatal(err)
	}
	fi2, _ := os.Stat(path)
	if !fi1.ModTime().Equal(fi2.ModTime()) {
		t.Error("unchanged entry should not rewrite the file")
	}
}

func TestDeregister_ExactNameAndForeignPreserved(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	writeFile(t, path, `{"mcpServers": {"to-remove": {"command": "echo"}, "to-remove-2": {"command": "x"}, "to-keep": {"command": "cat"}}}`, 0o644)
	a := New(path, "mcpServers", nil)
	if found, err := a.Deregister("to-remove"); err != nil || !found {
		t.Fatalf("Deregister = %v, %v", found, err)
	}
	if found, err := a.Deregister("to-"); err != nil || found {
		t.Fatalf("prefix must not match: %v, %v", found, err)
	}
	b := block(t, path, "mcpServers")
	if b["to-remove"] != nil || b["to-remove-2"] == nil || b["to-keep"] == nil {
		t.Errorf("block = %v", b)
	}
}

func TestDeregister_MissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nope.json")
	found, err := New(path, "mcpServers", nil).Deregister("x")
	if err != nil || found {
		t.Fatalf("Deregister = %v, %v", found, err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("deregister must not create the file")
	}
}

func TestDeregister_BlockOwnership(t *testing.T) {
	dir := t.TempDir()
	// Block created by us: removed entirely when emptied.
	created := filepath.Join(dir, "created.json")
	writeFile(t, created, "{\n  \"other\": true\n}\n", 0o644)
	a := New(created, "mcpServers", nil)
	if err := a.RegisterRemote("svc", "http://x", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := a.Deregister("svc"); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, created); got != "{\n  \"other\": true\n}\n" {
		t.Errorf("created block not removed:\n%s", got)
	}
	// Pre-existing block: left as {}.
	pre := filepath.Join(dir, "pre.json")
	writeFile(t, pre, "{\n  \"mcpServers\": {\n    \"svc\": {\"url\": \"http://x\"}\n  }\n}\n", 0o644)
	if _, err := New(pre, "mcpServers", nil).Deregister("svc"); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, pre); got != "{\n  \"mcpServers\": {}\n}\n" {
		t.Errorf("pre-existing block should remain {}:\n%q", got)
	}
}

// A config that doesn't exist yet, in a directory reached through a symlink (as macOS temp and
// home paths can be), must be keyed the same before and after the first write, or the created
// block is never recognized as ours.
func TestDeregister_BlockOwnershipThroughSymlinkedDir(t *testing.T) {
	real := t.TempDir()
	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(real, link); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(link, "mcp.json")
	a := New(path, "mcpServers", nil)
	if err := a.RegisterRemote("svc", "http://x", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := a.Deregister("svc"); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, path); strings.Contains(got, "mcpServers") {
		t.Errorf("block created by mcp-local should be removed once empty:\n%s", got)
	}
}

func TestAtomicWritePreservesModeAndBacksUp(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mcp.json")
	orig := `{"mcpServers": {}}`
	writeFile(t, path, orig, 0o600)
	if err := New(path, "mcpServers", nil).RegisterRemote("svc", "http://x", nil); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Errorf("mode = %v, want 0600", fi.Mode().Perm())
	}
	leftovers, _ := filepath.Glob(filepath.Join(filepath.Dir(path), ".mcp.json.tmp-*"))
	if len(leftovers) != 0 {
		t.Errorf("temp files left behind: %v", leftovers)
	}
	backups, _ := filepath.Glob(filepath.Join(BackupDir(), "*", backupName(path)))
	if len(backups) != 1 {
		t.Fatalf("backups = %v", backups)
	}
	if got := readFile(t, backups[0]); got != orig {
		t.Errorf("backup = %q, want original", got)
	}
}

func TestBackupRetention(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mcp.json")
	writeFile(t, path, `{}`, 0o644)
	base := time.Date(2026, 1, 2, 3, 4, 0, 0, time.Local)
	defer func() { now = time.Now }()
	a := New(path, "mcpServers", nil)
	for i := 0; i < MaxBackups+3; i++ {
		now = func() time.Time { return base.Add(time.Duration(i) * time.Second) }
		if err := a.RegisterRemote("svc", fmt.Sprintf("http://x/%d", i), nil); err != nil {
			t.Fatal(err)
		}
	}
	backups, _ := filepath.Glob(filepath.Join(BackupDir(), "*", backupName(path)))
	if len(backups) != MaxBackups {
		t.Fatalf("kept %d backups, want %d", len(backups), MaxBackups)
	}
	if !strings.Contains(backups[0], base.Add(3*time.Second).Format("20060102-150405")) {
		t.Errorf("oldest kept = %s", backups[0])
	}
}

func TestWriteThroughSymlink(t *testing.T) {
	dir := t.TempDir()
	real := filepath.Join(dir, "real.json")
	link := filepath.Join(dir, "link.json")
	writeFile(t, real, `{}`, 0o644)
	if err := os.Symlink(real, link); err != nil {
		t.Skip(err)
	}
	if err := New(link, "mcpServers", nil).RegisterRemote("svc", "http://x", nil); err != nil {
		t.Fatal(err)
	}
	if fi, _ := os.Lstat(link); fi.Mode()&os.ModeSymlink == 0 {
		t.Error("symlink replaced by regular file")
	}
	if block(t, real, "mcpServers")["svc"] == nil {
		t.Error("target not updated")
	}
}

func TestRegisterRemote_SameURLNoop(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	writeFile(t, path, `{"mcpServers": {"svc": {"url": "http://x", "headers": {"A": "b"}}}}`, 0o644)
	if err := New(path, "mcpServers", nil).RegisterRemote("svc", "http://x", nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(readFile(t, path), "headers") {
		t.Error("same-url register should not rewrite user extras")
	}
}

func TestRegisterLocal(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	a := New(path, "mcpServers", nil)
	if err := a.RegisterLocal("local-svc", []string{"/bin/echo", "hello"}, map[string]string{"KEY": "VAL"}, "env", "args"); err != nil {
		t.Fatal(err)
	}
	entry := block(t, path, "mcpServers")["local-svc"].(map[string]interface{})
	if entry["command"] != "/bin/echo" {
		t.Errorf("command = %v", entry["command"])
	}
	if args := entry["args"].([]interface{}); len(args) != 1 || args[0] != "hello" {
		t.Errorf("args = %v", args)
	}
	if env := entry["env"].(map[string]interface{}); env["KEY"] != "VAL" {
		t.Errorf("env = %v", env)
	}
	if err := a.RegisterLocal("x", nil, nil, "env", "args"); err == nil {
		t.Error("empty command should fail")
	}
}

func TestParseErrorLeavesFileUntouched(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	writeFile(t, path, `{"mcpServers": {`, 0o644)
	if err := New(path, "mcpServers", nil).RegisterRemote("svc", "http://x", nil); err == nil {
		t.Fatal("expected parse error")
	}
	if got := readFile(t, path); got != `{"mcpServers": {` {
		t.Errorf("file modified: %q", got)
	}
}
