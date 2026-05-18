package cursor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

func setupTestCursor(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	t.Cleanup(func() { os.Setenv("HOME", origHome) })

	cfgDir := filepath.Join(tmpDir, ".cursor")
	os.MkdirAll(cfgDir, 0755)
	return filepath.Join(cfgDir, "mcp.json")
}

func TestRegisterRemote(t *testing.T) {
	setupTestCursor(t)

	if err := RegisterRemote("test-svc", "http://localhost:8080/mcp"); err != nil {
		t.Fatalf("RegisterRemote: %v", err)
	}

	data, _ := os.ReadFile(ConfigPath())
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse: %v", err)
	}
	block := m["mcpServers"].(map[string]interface{})
	entry := block["test-svc"].(map[string]interface{})
	if entry["url"] != "http://localhost:8080/mcp" {
		t.Errorf("url = %v, want http://localhost:8080/mcp", entry["url"])
	}
}

func TestRegisterLocal(t *testing.T) {
	setupTestCursor(t)

	if err := RegisterLocal("stdio-svc", []string{"/usr/bin/mcp", "--arg"}, map[string]string{"KEY": "VAL"}); err != nil {
		t.Fatalf("RegisterLocal: %v", err)
	}

	data, _ := os.ReadFile(ConfigPath())
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse: %v", err)
	}
	block := m["mcpServers"].(map[string]interface{})
	entry := block["stdio-svc"].(map[string]interface{})
	if entry["command"] != "/usr/bin/mcp" {
		t.Errorf("command = %v, want /usr/bin/mcp", entry["command"])
	}
}

func TestRegisterServices(t *testing.T) {
	setupTestCursor(t)

	services := []config.ServiceConfig{
		{Name: "http-svc", Port: 8080, MCPType: "http"},
		{Name: "stdio-svc", Command: "/usr/bin/mcp", MCPType: "stdio"},
	}
	if err := RegisterServices(services); err != nil {
		t.Fatalf("RegisterServices: %v", err)
	}

	data, _ := os.ReadFile(ConfigPath())
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse: %v", err)
	}
	block := m["mcpServers"].(map[string]interface{})
	if _, ok := block["http-svc"]; !ok {
		t.Error("http-svc not registered")
	}
	if _, ok := block["stdio-svc"]; !ok {
		t.Error("stdio-svc not registered")
	}
}

func TestDeregister(t *testing.T) {
	setupTestCursor(t)

	_ = RegisterRemote("to-remove", "http://localhost:9999/mcp")
	if err := Deregister("to-remove"); err != nil {
		t.Fatalf("Deregister: %v", err)
	}

	data, _ := os.ReadFile(ConfigPath())
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse: %v", err)
	}
	block := m["mcpServers"].(map[string]interface{})
	if _, ok := block["to-remove"]; ok {
		t.Error("to-remove should be deregistered")
	}
}
