package agents

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/coma-toast/mcp-local/internal/mgr/claudedesktop"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/cursor"
	"github.com/coma-toast/mcp-local/internal/mgr/jsonagent"
	"github.com/coma-toast/mcp-local/internal/mgr/opencode"
)

const fakeASTMCP = `#!/bin/sh
echo "$*" >> "$FAKE_AST_LOG"
if [ "$1" = "--version" ]; then
  echo "ast-mcp version ${FAKE_AST_VERSION:-v4.0.0} (abc123)"
  exit 0
fi
echo '{"changes":[{"path":"/tmp/x.json","kind":"'"$1"'","diff":"","skipped":false,"reason":"conflict here"}],"status":[],"warnings":["heads up"]}'
exit ${FAKE_AST_EXIT:-0}
`

// setup redirects HOME, installs a fake ast-mcp, and returns the service plus its arg log path.
func setup(t *testing.T, version string, exit string) (config.ServiceConfig, string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	bin := filepath.Join(t.TempDir(), "ast-mcp")
	if err := os.WriteFile(bin, []byte(fakeASTMCP), 0o755); err != nil {
		t.Fatal(err)
	}
	log := filepath.Join(t.TempDir(), "args.log")
	t.Setenv("FAKE_AST_LOG", log)
	t.Setenv("FAKE_AST_VERSION", version)
	t.Setenv("FAKE_AST_EXIT", exit)
	var warn bytes.Buffer
	orig := Warn
	Warn = &warn
	t.Cleanup(func() { Warn = orig })
	return config.ServiceConfig{Name: "ast-context-cache", Command: bin, MCPType: "http", Port: 7821}, log
}

func logLines(t *testing.T, path string) []string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(strings.TrimSpace(string(b)), "\n")
}

func registeredIn(path, block, name string) bool {
	m, err := jsonagent.ReadJSON(path)
	if err != nil {
		return false
	}
	b, _ := m[block].(map[string]interface{})
	_, ok := b[name]
	return ok
}

func TestRegisterAll_DelegatesToASTMCP(t *testing.T) {
	svc, log := setup(t, "v4.0.0", "0")
	if err := RegisterAll([]config.ServiceConfig{svc}, DefaultTargets()); err != nil {
		t.Fatalf("RegisterAll: %v", err)
	}
	want := []string{
		"--version",
		"install --target opencode --component mcp --yes --json --mcp-url http://localhost:7821/mcp",
		"install --target cursor --component mcp --yes --json --mcp-url http://localhost:7821/mcp",
		"install --target claude_desktop --component mcp --yes --json --mcp-url http://localhost:7821/mcp",
	}
	if got := logLines(t, log); strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("args:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if _, err := os.Stat(cursor.ConfigPath()); !os.IsNotExist(err) {
		t.Error("native cursor writer must not run when delegated")
	}
	if !strings.Contains(Warn.(*bytes.Buffer).String(), "heads up") {
		t.Errorf("installer warnings not surfaced: %q", Warn.(*bytes.Buffer).String())
	}
}

func TestDeregister_DelegatesUninstall(t *testing.T) {
	svc, log := setup(t, "4.2.1", "0")
	found, err := Deregister(svc, HostClaude)
	if err != nil || !found {
		t.Fatalf("Deregister = %v, %v", found, err)
	}
	got := logLines(t, log)
	if got[len(got)-1] != "uninstall --target claude_desktop --component mcp --yes --json --mcp-url http://localhost:7821/mcp" {
		t.Errorf("args = %v", got)
	}
}

func TestRegisterAll_OldVersionFallsBackToNative(t *testing.T) {
	svc, log := setup(t, "3.0.40", "0")
	if err := RegisterAll([]config.ServiceConfig{svc}, DefaultTargets()); err != nil {
		t.Fatalf("RegisterAll: %v", err)
	}
	if got := logLines(t, log); len(got) != 1 || got[0] != "--version" {
		t.Errorf("only --version should run (cached), got %v", got)
	}
	if !registeredIn(opencode.ConfigPath(), "mcp", svc.Name) || !registeredIn(cursor.ConfigPath(), "mcpServers", svc.Name) || !registeredIn(claudedesktop.ConfigPath(), "mcpServers", svc.Name) {
		t.Error("native writers should register on fallback")
	}
	if !strings.Contains(Warn.(*bytes.Buffer).String(), "3.0.40") {
		t.Errorf("fallback warning missing: %q", Warn.(*bytes.Buffer).String())
	}
}

func TestRegisterAll_InstallerErrorsAggregated(t *testing.T) {
	svc, _ := setup(t, "4.0.0", "3")
	other := config.ServiceConfig{Name: "plain", MCPType: "http", MCPURL: "http://localhost:9/mcp"}
	err := RegisterAll([]config.ServiceConfig{svc, other}, DefaultTargets())
	if err == nil {
		t.Fatal("expected aggregated error")
	}
	for _, want := range []string{"opencode: ast-context-cache", "cursor: ast-context-cache", "claude: ast-context-cache", "conflict or parse error", "conflict here"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error missing %q: %v", want, err)
		}
	}
	if !registeredIn(cursor.ConfigPath(), "mcpServers", "plain") {
		t.Error("non-delegated service should still register natively")
	}
}

func TestRegisterAll_UnsupportedFallsBack(t *testing.T) {
	svc, _ := setup(t, "4.0.0", "4")
	if err := RegisterAll([]config.ServiceConfig{svc}, Targets{Cursor: true}); err != nil {
		t.Fatalf("RegisterAll: %v", err)
	}
	if !registeredIn(cursor.ConfigPath(), "mcpServers", svc.Name) {
		t.Error("exit 4 should fall back to native writer")
	}
}

func TestRegisterAll_ContinuesPastHostError(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	if err := os.MkdirAll(filepath.Dir(cursor.ConfigPath()), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cursor.ConfigPath(), []byte(`{"mcpServers": {`), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := config.ServiceConfig{Name: "svc", MCPType: "http", MCPURL: "http://localhost:9/mcp"}
	err := RegisterAll([]config.ServiceConfig{svc}, DefaultTargets())
	if err == nil || !strings.Contains(err.Error(), "cursor:") {
		t.Fatalf("expected cursor error, got %v", err)
	}
	if !registeredIn(opencode.ConfigPath(), "mcp", "svc") || !registeredIn(claudedesktop.ConfigPath(), "mcpServers", "svc") {
		t.Error("other hosts should still be registered")
	}
	err = DeregisterAll(svc, DefaultTargets())
	if err == nil || !strings.Contains(err.Error(), "cursor:") {
		t.Fatalf("expected cursor deregister error, got %v", err)
	}
	if registeredIn(opencode.ConfigPath(), "mcp", "svc") || registeredIn(claudedesktop.ConfigPath(), "mcpServers", "svc") {
		t.Error("other hosts should still be deregistered")
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		in    string
		major int
		ok    bool
	}{
		{"ast-mcp v4.0.0", 4, true},
		{"3.0.40\n", 3, true},
		{"ast-mcp 10.1.2-rc1 (deadbeef)", 10, true},
		{"dev", 0, false},
	}
	for _, tt := range tests {
		major, _, ok := parseVersion(tt.in)
		if major != tt.major || ok != tt.ok {
			t.Errorf("parseVersion(%q) = %d, %v", tt.in, major, ok)
		}
	}
}
