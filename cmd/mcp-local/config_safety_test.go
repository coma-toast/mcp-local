package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/cursor"
	"github.com/coma-toast/mcp-local/internal/mgr/jsonagent"
	"github.com/spf13/cobra"
)

// tempHome points HOME at a temp dir and writes config.yaml there when content is non-empty.
func tempHome(t *testing.T, content string) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	if content != "" {
		path, _ := config.GetConfigPath()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return home
}

func runCmd(t *testing.T, c *cobra.Command, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	c.SetArgs(args)
	c.SetOut(&out)
	c.SetErr(&out)
	err := c.Execute()
	return out.String(), err
}

func TestExecutePrintsErrorOnce(t *testing.T) {
	fail := &cobra.Command{Use: "fail-once", SilenceUsage: true, RunE: func(*cobra.Command, []string) error { return errors.New("boom") }}
	rootCmd.AddCommand(fail)
	defer rootCmd.RemoveCommand(fail)
	var cobraOut, stderr bytes.Buffer
	rootCmd.SetErr(&cobraOut)
	rootCmd.SetArgs([]string{"fail-once"})
	defer func() { rootCmd.SetErr(nil); rootCmd.SetArgs(nil) }()
	if code := execute(&stderr); code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
	if n := strings.Count(cobraOut.String()+stderr.String(), "boom"); n != 1 {
		t.Errorf("error printed %d times: cobra %q, main %q", n, cobraOut.String(), stderr.String())
	}
}

func stubWizard(t *testing.T, tty bool) *bool {
	t.Helper()
	called := false
	origWizard, origTTY := runWizard, stdinIsTerminal
	runWizard = func() (*config.ManagerConfig, error) {
		called = true
		return &config.ManagerConfig{Services: []config.ServiceConfig{{Name: "w", Command: "/bin/w"}}}, nil
	}
	stdinIsTerminal = func() bool { return tty }
	t.Cleanup(func() { runWizard, stdinIsTerminal = origWizard, origTTY })
	return &called
}

func TestEnsureConfig_CorruptFileIsAnErrorNotTheWizard(t *testing.T) {
	tempHome(t, "services:\n  - name: [unclosed\n")
	called := stubWizard(t, true)
	path, _ := config.GetConfigPath()
	_, err := ensureConfig()
	if err == nil || !strings.Contains(err.Error(), path) {
		t.Fatalf("err = %v; want it to name %s", err, path)
	}
	if *called {
		t.Error("wizard ran for a corrupt config")
	}
	if b, _ := os.ReadFile(path); string(b) != "services:\n  - name: [unclosed\n" {
		t.Error("corrupt config was overwritten")
	}
}

func TestEnsureConfig_MissingWithoutTTY(t *testing.T) {
	tempHome(t, "")
	called := stubWizard(t, false)
	if _, err := ensureConfig(); err == nil || *called {
		t.Fatalf("err = %v, wizard ran = %v", err, *called)
	}
}

func TestEnsureConfig_DanglingSymlinkIsAnError(t *testing.T) {
	home := tempHome(t, "")
	called := stubWizard(t, true)
	path, _ := config.GetConfigPath()
	os.MkdirAll(filepath.Dir(path), 0o755)
	if err := os.Symlink(filepath.Join(home, "configSync", "mcp-local.yaml"), path); err != nil {
		t.Skip(err)
	}
	if _, err := ensureConfig(); err == nil || *called {
		t.Fatalf("err = %v, wizard ran = %v", err, *called)
	}
}

func TestEnsureConfig_MissingWithTTYRunsWizard(t *testing.T) {
	tempHome(t, "")
	called := stubWizard(t, true)
	cfg, err := ensureConfig()
	if err != nil || !*called || len(cfg.Services) != 1 {
		t.Fatalf("cfg = %+v, err = %v, wizard ran = %v", cfg, err, *called)
	}
	if _, err := config.LoadConfig(); err != nil {
		t.Errorf("wizard config not saved: %v", err)
	}
}

func TestServiceFromEntry(t *testing.T) {
	tests := []struct {
		name  string
		entry map[string]interface{}
		want  config.ServiceConfig
	}{
		{"cursor stdio", map[string]interface{}{"command": "npx", "args": []interface{}{"-y", "@playwright/mcp"}, "env": map[string]interface{}{"K": "v"}},
			config.ServiceConfig{MCPType: "stdio", Command: "npx", Args: []string{"-y", "@playwright/mcp"}, Env: map[string]string{"K": "v"}}},
		{"opencode local", map[string]interface{}{"type": "local", "command": []interface{}{"uvx", "fetch"}, "environment": map[string]interface{}{"K": "v"}},
			config.ServiceConfig{MCPType: "stdio", Command: "uvx", Args: []string{"fetch"}, Env: map[string]string{"K": "v"}}},
		{"cursor url with headers", map[string]interface{}{"url": "https://x/mcp", "headers": map[string]interface{}{"Authorization": "Bearer t"}},
			config.ServiceConfig{MCPType: "http", MCPURL: "https://x/mcp", Headers: map[string]string{"Authorization": "Bearer t"}}},
		{"bridge with header", map[string]interface{}{"command": "/usr/local/bin/mcp-local", "args": []interface{}{"bridge", "https://x/mcp", "--header", "A=b=c"}},
			config.ServiceConfig{MCPType: "http", MCPURL: "https://x/mcp", Headers: map[string]string{"A": "b=c"}}},
		{"bridge with other flags stays stdio", map[string]interface{}{"command": "mcp-local", "args": []interface{}{"bridge", "https://x/mcp", "--timeout", "2m"}},
			config.ServiceConfig{MCPType: "stdio", Command: "mcp-local", Args: []string{"bridge", "https://x/mcp", "--timeout", "2m"}}},
	}
	for _, tt := range tests {
		got, ok := serviceFromEntry("svc", tt.entry)
		tt.want.Name = "svc"
		if !ok || !equalService(got, tt.want) {
			t.Errorf("%s: got %+v, want %+v", tt.name, got, tt.want)
		}
	}
	if _, ok := serviceFromEntry("svc", map[string]interface{}{"type": "stdio"}); ok {
		t.Error("an entry without command or url should not import")
	}
}

func equalService(a, b config.ServiceConfig) bool {
	return a.Name == b.Name && a.MCPType == b.MCPType && a.Command == b.Command && a.MCPURL == b.MCPURL &&
		strings.Join(a.Args, "\x00") == strings.Join(b.Args, "\x00") && mapsEqual(a.Env, b.Env) && mapsEqual(a.Headers, b.Headers)
}

func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func TestMergeSyncedToolsKeepsOnlyUserOverrides(t *testing.T) {
	server := []config.ToolConfig{
		{Name: "a", Description: "server a", Enabled: true},
		{Name: "b", Description: "server b", Enabled: true},
		{Name: "c", Description: "server c", Enabled: true},
	}
	existing := map[string]config.ToolConfig{
		"a": {Name: "a", Description: "server a", Enabled: true}, // frozen by an older sync
		"b": {Name: "b", Description: "mine", Enabled: false, Tier: "core"},
	}
	got := mergeSyncedTools(server, existing)
	if got[0].Description != "" || got[2].Description != "" {
		t.Errorf("server descriptions stored as overrides: %+v", got)
	}
	if got[1].Description != "mine" || got[1].Enabled || got[1].Tier != "core" || !got[2].Enabled {
		t.Errorf("user settings lost: %+v", got)
	}
}

const twoServices = `agents:
  cursor: true
services:
  - name: stdio1
    command: /bin/echo
  - name: http1
    type: http
    port: 1
`

func writeCursor(t *testing.T, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(cursor.ConfigPath()), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cursor.ConfigPath(), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func cursorHas(name string) bool {
	return jsonagent.New(cursor.ConfigPath(), "mcpServers", nil).Has(name)
}

func TestRemove_KeepsServiceWhenDeregisterFails(t *testing.T) {
	tempHome(t, twoServices)
	writeCursor(t, `{"mcpServers": {`)
	if _, err := runCmd(t, cmdRemove(), "http1"); err == nil || !strings.Contains(err.Error(), "--deregister=false") {
		t.Fatalf("err = %v", err)
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := cfg.ServiceNamed("http1"); err != nil {
		t.Error("service removed from config although deregistration failed")
	}
}

func TestDeregister_BareTargetsOnlyStoppedHTTPEntriesMCPLocalWrote(t *testing.T) {
	tempHome(t, twoServices)
	orig := portIsRunning
	portIsRunning = func(int) bool { return false }
	defer func() { portIsRunning = orig }()
	writeCursor(t, `{"mcpServers": {"stdio1": {"command": "/bin/echo"}, "http1": {"url": "http://user"}}}`)
	out, err := runCmd(t, cmdDeregister())
	if err != nil {
		t.Fatal(err)
	}
	if !cursorHas("stdio1") || !cursorHas("http1") || !strings.Contains(out, "skipped (entry not written by mcp-local") {
		t.Errorf("user entries touched or no note: %q", out)
	}
	if err := cursor.RegisterServices([]config.ServiceConfig{{Name: "http1", Port: 1, MCPType: "http"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(t, cmdDeregister()); err != nil {
		t.Fatal(err)
	}
	if cursorHas("http1") || !cursorHas("stdio1") {
		t.Error("bare deregister should remove only the HTTP entry mcp-local wrote")
	}
	writeCursor(t, `{"mcpServers": {`)
	if _, err := runCmd(t, cmdDeregister(), "http1"); err == nil {
		t.Error("a host failure must fail the command")
	}
}

func TestValidate_CommandChecks(t *testing.T) {
	tempHome(t, `services:
  - name: nocmd
    type: stdio
  - name: onpath
    command: sh
  - name: missing
    command: definitely-not-a-command-xyz
  - name: remote
    mcp_url: https://x/mcp
`)
	out, err := runCmd(t, cmdValidate())
	if err == nil {
		t.Fatal("expected validation failure")
	}
	if !strings.Contains(out, "[nocmd] ERROR: stdio service has no command") || !strings.Contains(out, "[missing] ERROR: command not found") {
		t.Errorf("output: %s", out)
	}
	if strings.Contains(out, "[onpath]") || strings.Contains(out, "[remote]") {
		t.Errorf("false positives: %s", out)
	}
}
