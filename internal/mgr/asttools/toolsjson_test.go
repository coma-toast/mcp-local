package asttools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
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
