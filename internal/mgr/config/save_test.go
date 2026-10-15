package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/coma-toast/mcp-local/internal/mgr/safefile"
)

const sharedConfig = `# mcp-local config shared by both Macs
agents:
  opencode: true
  cursor: false
  claude: false

services:
  # the ast server
  - name: ast-context-cache
    command: ${HOME}/git/ast-context-cache/ast-mcp # built by make
    port: 7821
    type: http
    mcp_url: http://localhost:7821/mcp
    log: ~/logs/ast.log
    x_future_key: keep me
    env:
      AST_HOME: ${HOME}/.astcache
    tools:
      - name: execute_code
        enabled: false
      - name: index_files
        tier: core
  - name: fetch
    command: uvx
    args: ["mcp-server-fetch"]
`

// setupShared makes ~/.mcp-local/config.yaml a symlink to a file elsewhere, like the
// Syncthing-shared config, and returns the target path.
func setupShared(t *testing.T, content string) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	target := filepath.Join(t.TempDir(), "mcp-local.yaml")
	if err := os.WriteFile(target, []byte(content), 0o640); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(home, ".mcp-local"), 0o755); err != nil {
		t.Fatal(err)
	}
	link, _ := GetConfigPath()
	if err := os.Symlink(target, link); err != nil {
		t.Skip(err)
	}
	return target
}

func TestSaveConfigKeepsSymlinkCommentsAndUnexpandedPaths(t *testing.T) {
	target := setupShared(t, sharedConfig)
	before, _ := os.Stat(target)
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	home, _ := os.UserHomeDir()
	ast, _ := cfg.ServiceNamed("ast-context-cache")
	if ast.Command != home+"/git/ast-context-cache/ast-mcp" || ast.Env["AST_HOME"] != home+"/.astcache" {
		t.Fatalf("callers must see expanded values: %q %q", ast.Command, ast.Env["AST_HOME"])
	}
	if !ast.Tools[1].Enabled {
		t.Error("a tool without enabled: must be enabled")
	}
	ast.Port = 7822
	cfg.UpsertService(ast)
	cfg.UpsertService(ServiceConfig{Name: "new", Command: "/opt/new"})
	if err := SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	link, _ := GetConfigPath()
	if fi, _ := os.Lstat(link); fi.Mode()&os.ModeSymlink == 0 {
		t.Fatal("config.yaml symlink replaced by a regular file")
	}
	b, _ := os.ReadFile(target)
	got := string(b)
	for _, want := range []string{
		"# mcp-local config shared by both Macs", "# the ast server", "# built by make",
		"command: ${HOME}/git/ast-context-cache/ast-mcp", "log: ~/logs/ast.log", "AST_HOME: ${HOME}/.astcache",
		"x_future_key: keep me", "port: 7822", "name: new", "cursor: false", "args: [\"mcp-server-fetch\"]",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("saved config lacks %q:\n%s", want, got)
		}
	}
	for _, bad := range []string{home, "installer:", "enabled: true", "deps:", "build_cmd:"} {
		if strings.Contains(got, bad) {
			t.Errorf("saved config contains %q:\n%s", bad, got)
		}
	}
	after, _ := os.Stat(target)
	if os.SameFile(before, after) {
		t.Error("target was rewritten in place; want temp file + rename")
	}
	if after.Mode().Perm() != 0o640 {
		t.Errorf("mode = %v, want 0640", after.Mode().Perm())
	}
	backups, _ := filepath.Glob(filepath.Join(safefile.BackupDir(), "*", safefile.BackupName(link)))
	if len(backups) != 1 {
		t.Fatalf("backups = %v", backups)
	}
	if b, _ := os.ReadFile(backups[0]); string(b) != sharedConfig {
		t.Error("backup is not the previous file")
	}
	reloaded, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if s, _ := reloaded.ServiceNamed("ast-context-cache"); s.Port != 7822 || s.Command != ast.Command || s.Tools[0].Enabled {
		t.Errorf("reloaded = %+v", s)
	}
}

func TestSaveConfigUnchangedDoesNotWrite(t *testing.T) {
	target := setupShared(t, sharedConfig)
	before, _ := os.Stat(target)
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if err := SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	after, _ := os.Stat(target)
	if !os.SameFile(before, after) {
		t.Error("an unchanged config was rewritten")
	}
}

func TestSaveConfigChangedPathAndRemovedKeys(t *testing.T) {
	target := setupShared(t, sharedConfig)
	cfg, _ := LoadConfig()
	ast, _ := cfg.ServiceNamed("ast-context-cache")
	ast.Command = "/usr/local/bin/ast-mcp"
	ast.Log = ""
	ast.Tools = []ToolConfig{{Name: "execute_code", Enabled: true}}
	cfg.UpsertService(ast)
	cfg.Services = cfg.Services[:1]
	if err := SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(target)
	got := string(b)
	if !strings.Contains(got, "command: /usr/local/bin/ast-mcp") || strings.Contains(got, "log:") ||
		strings.Contains(got, "name: fetch") || strings.Contains(got, "enabled: false") || strings.Contains(got, "index_files") {
		t.Errorf("saved:\n%s", got)
	}
}

func TestSaveConfigWithoutLoadMergesIntoFile(t *testing.T) {
	target := setupShared(t, sharedConfig)
	cfg := &ManagerConfig{Services: []ServiceConfig{{Name: "fetch", Command: "uvx", Args: []string{"mcp-server-fetch"}}}}
	if err := SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(target)
	if !strings.Contains(string(b), "# mcp-local config shared by both Macs") || strings.Contains(string(b), "ast-context-cache") {
		t.Errorf("saved:\n%s", b)
	}
}

func TestLoadConfigAgentsPresence(t *testing.T) {
	cfg, err := parseConfig([]byte("services: []\n"))
	if err != nil || cfg.Agents != nil {
		t.Fatalf("no agents key: %+v, %v", cfg.Agents, err)
	}
	cfg, err = parseConfig([]byte("agents:\n  opencode: false\n"))
	if err != nil || cfg.Agents == nil || cfg.Agents.OpenCode {
		t.Fatalf("agents key present: %+v, %v", cfg.Agents, err)
	}
}

func TestIsHTTPRemoteOnly(t *testing.T) {
	if !(ServiceConfig{MCPURL: "https://x/mcp"}).IsHTTP() {
		t.Error("mcp_url only should be HTTP")
	}
	if (ServiceConfig{Command: "npx", MCPURL: "https://x/mcp"}).IsHTTP() {
		t.Error("a command without port/type is stdio")
	}
	if (ServiceConfig{MCPType: "stdio", MCPURL: "https://x/mcp"}).IsHTTP() {
		t.Error("explicit stdio wins")
	}
	e, err := ServiceToEntry(ServiceConfig{MCPURL: "https://x/mcp", Headers: map[string]string{"A": "b"}})
	if err != nil || e.Type != "remote" || e.URL != "https://x/mcp" || e.Headers["A"] != "b" {
		t.Errorf("entry = %+v, %v", e, err)
	}
	if _, err := ServiceToEntry(ServiceConfig{MCPType: "stdio"}); err == nil {
		t.Error("stdio without command should not produce an entry")
	}
}
