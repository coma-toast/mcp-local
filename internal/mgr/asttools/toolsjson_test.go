package asttools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/safefile"
)

func TestWriteToolsJSON_disabledTool(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tools.json")
	tools := []config.ToolConfig{
		{Name: "execute_code", Enabled: false, Tier: "complete"},
		{Name: "index_files", Enabled: true},
	}
	if err := WriteToolsJSON(path, tools); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]fileEntry
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	if m["execute_code"].Enabled {
		t.Fatal("execute_code should be disabled")
	}
	if _, ok := m["index_files"]; ok {
		t.Fatal("enabled tool with no overrides should be omitted")
	}
}

func TestWriteToolsJSON_tierOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tools.json")
	tools := []config.ToolConfig{
		{Name: "index_files", Enabled: true, Tier: "core", Description: "promoted"},
	}
	if err := WriteToolsJSON(path, tools); err != nil {
		t.Fatal(err)
	}
	var m map[string]fileEntry
	data, _ := os.ReadFile(path)
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	if m["index_files"].Tier != "core" || m["index_files"].Description != "promoted" {
		t.Fatalf("got %+v", m["index_files"])
	}
}

func TestApplyStartEnv_setsTierAndToolsConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tools.json")
	svc := config.ServiceConfig{
		Name:       "ast-context-cache",
		ActiveTier: "extended",
		Tools: []config.ToolConfig{
			{Name: "execute_code", Enabled: false},
		},
	}
	_ = path
	svc.ToolsConfigPath = path
	env, err := ApplyStartEnv([]string{}, svc)
	if err != nil {
		t.Fatal(err)
	}
	var hasTier, hasConfig bool
	for _, e := range env {
		if e == "AST_MCP_TIER=extended" {
			hasTier = true
		}
		if e == "AST_MCP_TOOLS_CONFIG="+path {
			hasConfig = true
		}
	}
	if !hasTier || !hasConfig {
		t.Fatalf("env: %v", env)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("tools.json not written")
	}
}

func TestMain(m *testing.M) {
	home, err := os.MkdirTemp("", "asttools-home-*")
	if err != nil {
		panic(err)
	}
	os.Setenv("HOME", home)
	code := m.Run()
	os.RemoveAll(home)
	os.Exit(code)
}

func readEntries(t *testing.T, path string) map[string]map[string]interface{} {
	t.Helper()
	var m map[string]map[string]interface{}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	return m
}

func TestWriteToolsJSON_mergesWithDashboardEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tools.json")
	orig := `{"dash_tool": {"enabled": false, "note": "x"}, "execute_code": {"enabled": true, "extra": 1}}` + "\n"
	if err := os.WriteFile(path, []byte(orig), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := WriteToolsJSON(path, []config.ToolConfig{{Name: "execute_code", Enabled: false}, {Name: "index_files", Enabled: true}}); err != nil {
		t.Fatal(err)
	}
	m := readEntries(t, path)
	if m["dash_tool"]["note"] != "x" || m["dash_tool"]["enabled"] != false {
		t.Errorf("dashboard entry lost: %v", m)
	}
	if m["execute_code"]["enabled"] != false || m["execute_code"]["extra"] != float64(1) {
		t.Errorf("override not merged: %v", m["execute_code"])
	}
	if fi, _ := os.Stat(path); fi.Mode().Perm() != 0o600 {
		t.Errorf("mode = %v", fi.Mode().Perm())
	}
	backups, _ := filepath.Glob(filepath.Join(safefile.BackupDir(), "*", safefile.BackupName(path)))
	if len(backups) != 1 {
		t.Errorf("backups = %v", backups)
	}
}

func TestWriteToolsJSON_noOverridesKeepsFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tools.json")
	orig := `{"dash_tool": {"enabled": false}}` + "\n"
	if err := os.WriteFile(path, []byte(orig), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := WriteToolsJSON(path, []config.ToolConfig{{Name: "dash_tool", Enabled: true}}); err != nil {
		t.Fatal(err)
	}
	if b, err := os.ReadFile(path); err != nil || string(b) != orig {
		t.Errorf("file changed or deleted: %q, %v", b, err)
	}
}

func TestWriteToolsJSON_droppedOverrideRemovesOwnEntry(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tools.json")
	if err := WriteToolsJSON(path, []config.ToolConfig{{Name: "execute_code", Enabled: false}}); err != nil {
		t.Fatal(err)
	}
	if err := WriteToolsJSON(path, []config.ToolConfig{{Name: "execute_code", Enabled: true}}); err != nil {
		t.Fatal(err)
	}
	if m := readEntries(t, path); len(m) != 0 {
		t.Errorf("re-enabled tool still overridden: %v", m)
	}
}
