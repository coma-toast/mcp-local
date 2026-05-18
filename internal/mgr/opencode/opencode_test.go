package opencode

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

func setupTestConfig(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	t.Cleanup(func() { os.Setenv("HOME", origHome) })

	cfgDir := filepath.Join(tmpDir, ".config", "opencode")
	os.MkdirAll(cfgDir, 0755)
	return filepath.Join(cfgDir, "opencode.jsonc")
}

func TestRegisterRemote(t *testing.T) {
	setupTestConfig(t)

	if err := RegisterRemote("test-svc", "http://localhost:8080/mcp", 30000); err != nil {
		t.Fatalf("RegisterRemote: %v", err)
	}

	data, _ := os.ReadFile(ConfigPath())
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse: %v", err)
	}
	block := m["mcp"].(map[string]interface{})
	entry := block["test-svc"].(map[string]interface{})
	if entry["type"] != "remote" {
		t.Errorf("type = %v, want remote", entry["type"])
	}
	if entry["url"] != "http://localhost:8080/mcp" {
		t.Errorf("url = %v, want http://localhost:8080/mcp", entry["url"])
	}
}

func TestRegisterLocal(t *testing.T) {
	setupTestConfig(t)

	if err := RegisterLocal("stdio-svc", []string{"/usr/bin/mcp", "--arg"}, map[string]string{"KEY": "VAL"}); err != nil {
		t.Fatalf("RegisterLocal: %v", err)
	}

	data, _ := os.ReadFile(ConfigPath())
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse: %v", err)
	}
	block := m["mcp"].(map[string]interface{})
	entry := block["stdio-svc"].(map[string]interface{})
	if entry["type"] != "local" {
		t.Errorf("type = %v, want local", entry["type"])
	}
}

func TestRegisterServices(t *testing.T) {
	setupTestConfig(t)

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
	block := m["mcp"].(map[string]interface{})
	if _, ok := block["http-svc"]; !ok {
		t.Error("http-svc not registered")
	}
	if _, ok := block["stdio-svc"]; !ok {
		t.Error("stdio-svc not registered")
	}
}

func TestDeregister(t *testing.T) {
	setupTestConfig(t)

	_ = RegisterRemote("to-remove", "http://localhost:9999/mcp", 30000)
	ok, err := Deregister("to-remove")
	if err != nil {
		t.Fatalf("Deregister: %v", err)
	}
	if !ok {
		t.Error("Deregister should return true when entry exists")
	}

	ok, err = Deregister("nonexistent")
	if err != nil {
		t.Fatalf("Deregister nonexistent: %v", err)
	}
	if ok {
		t.Error("Deregister should return false for nonexistent entry")
	}
}

func TestIdempotentRegister(t *testing.T) {
	setupTestConfig(t)

	if err := RegisterRemote("svc", "http://localhost:8080/mcp", 30000); err != nil {
		t.Fatalf("first register: %v", err)
	}
	data1, _ := os.ReadFile(ConfigPath())

	if err := RegisterRemote("svc", "http://localhost:8080/mcp", 30000); err != nil {
		t.Fatalf("second register: %v", err)
	}
	data2, _ := os.ReadFile(ConfigPath())

	if string(data1) != string(data2) {
		t.Error("idempotent register should produce same file")
	}
}
