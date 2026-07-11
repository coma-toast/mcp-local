package jsonagent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

func TestReadJSON(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.json")

	data := map[string]interface{}{
		"key": "value",
		"nested": map[string]interface{}{
			"inner": 123,
		},
	}
	b, _ := json.Marshal(data)
	if err := os.WriteFile(testFile, b, 0644); err != nil {
		t.Fatal(err)
	}

	got, err := ReadJSON(testFile)
	if err != nil {
		t.Fatal(err)
	}

	if got["key"] != "value" {
		t.Errorf("key = %v, want value", got["key"])
	}
	nested := got["nested"].(map[string]interface{})
	if nested["inner"] != float64(123) {
		t.Errorf("nested.inner = %v, want 123", nested["inner"])
	}
}

func TestReadJSON_NotExist(t *testing.T) {
	_, err := ReadJSON("/nonexistent/path.json")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestWriteJSON(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "output.json")

	data := map[string]interface{}{
		"key": "value",
	}

	if err := WriteJSON(testFile, data); err != nil {
		t.Fatal(err)
	}

	// Read back and verify
	got, err := ReadJSON(testFile)
	if err != nil {
		t.Fatal(err)
	}

	if got["key"] != "value" {
		t.Errorf("key = %v, want value", got["key"])
	}
}

func TestWriteJSON_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "subdir", "output.json")

	data := map[string]interface{}{"key": "value"}
	if err := WriteJSON(testFile, data); err != nil {
		t.Fatal(err)
	}

	got, err := ReadJSON(testFile)
	if err != nil {
		t.Fatal(err)
	}
	if got["key"] != "value" {
		t.Errorf("key = %v, want value", got["key"])
	}
}

func TestEnsureBlock(t *testing.T) {
	a := New("/tmp/test", "mcp", ReadJSON, WriteJSON, func(e config.AgentEntry) map[string]interface{} { return nil })

	// Test creating block when nil
	m := map[string]interface{}{}
	block := a.EnsureBlock(m)
	if block == nil {
		t.Error("EnsureBlock should return a map")
	}
	// Check that the block is stored in the parent map by checking pointer equality
	if _, ok := m["mcp"]; !ok {
		t.Error("EnsureBlock should set the block in the parent map")
	}

	// Test returning existing block
	m2 := map[string]interface{}{"mcp": map[string]interface{}{"existing": "value"}}
	block2 := a.EnsureBlock(m2)
	if block2["existing"] != "value" {
		t.Error("EnsureBlock should return existing block")
	}
}

func TestRegisterServices(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "config.json")

	a := New(testFile, "mcpServers", ReadJSON, WriteJSON, func(e config.AgentEntry) map[string]interface{} {
		if e.Type == "remote" {
			return map[string]interface{}{"type": "http", "url": e.URL}
		}
		return map[string]interface{}{"type": "stdio", "command": e.Command[0]}
	})

	services := []config.ServiceConfig{
		{Name: "test-service", Command: "/bin/echo", Args: []string{"hello"}, MCPType: "stdio"},
		{Name: "remote-service", MCPType: "http", MCPURL: "http://localhost:8080/mcp"},
	}

	if err := a.RegisterServices(services); err != nil {
		t.Fatal(err)
	}

	// Verify written
	got, err := ReadJSON(testFile)
	if err != nil {
		t.Fatal(err)
	}

	block := got["mcpServers"].(map[string]interface{})
	if block["test-service"] == nil {
		t.Error("test-service not registered")
	}
	if block["remote-service"] == nil {
		t.Error("remote-service not registered")
	}
}

func TestDeregister(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "config.json")

	// Pre-populate with a service
	initial := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"to-remove": map[string]interface{}{"type": "stdio", "command": "echo"},
			"to-keep":   map[string]interface{}{"type": "stdio", "command": "cat"},
		},
	}
	b, _ := json.Marshal(initial)
	if err := os.WriteFile(testFile, b, 0644); err != nil {
		t.Fatal(err)
	}

	a := New(testFile, "mcpServers", ReadJSON, WriteJSON, nil)

	// Deregister existing
	removed, err := a.Deregister("to-remove")
	if err != nil {
		t.Fatal(err)
	}
	if !removed {
		t.Error("Deregister should return true for existing service")
	}

	// Try to deregister non-existing
	removed2, err := a.Deregister("nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if removed2 {
		t.Error("Deregister should return false for non-existing service")
	}

	// Verify final state
	got, err := ReadJSON(testFile)
	if err != nil {
		t.Fatal(err)
	}
	block := got["mcpServers"].(map[string]interface{})
	if block["to-remove"] != nil {
		t.Error("to-remove should have been removed")
	}
	if block["to-keep"] == nil {
		t.Error("to-keep should still exist")
	}
}

func TestRegisterRemote(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "config.json")

	a := New(testFile, "mcpServers", ReadJSON, WriteJSON, nil)

	if err := a.RegisterRemote("remote-svc", "http://example.com/mcp", nil); err != nil {
		t.Fatal(err)
	}

	got, err := ReadJSON(testFile)
	if err != nil {
		t.Fatal(err)
	}

	block := got["mcpServers"].(map[string]interface{})
	entry := block["remote-svc"].(map[string]interface{})
	if entry["url"] != "http://example.com/mcp" {
		t.Errorf("url = %v, want http://example.com/mcp", entry["url"])
	}
}

func TestRegisterLocal(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "config.json")

	a := New(testFile, "mcpServers", ReadJSON, WriteJSON, nil)

	if err := a.RegisterLocal("local-svc", []string{"/bin/echo", "hello"}, map[string]string{"KEY": "VAL"}, "env", "args"); err != nil {
		t.Fatal(err)
	}

	got, err := ReadJSON(testFile)
	if err != nil {
		t.Fatal(err)
	}

	block := got["mcpServers"].(map[string]interface{})
	entry := block["local-svc"].(map[string]interface{})
	if entry["command"] != "/bin/echo" {
		t.Errorf("command = %v, want /bin/echo", entry["command"])
	}
	args := entry["args"].([]interface{})
	if len(args) != 1 || args[0] != "hello" {
		t.Errorf("args = %v, want [hello]", args)
	}
	env := entry["env"].(map[string]interface{})
	if env["KEY"] != "VAL" {
		t.Errorf("env = %v, want KEY=VAL", env)
	}
}

func TestRegisterLocal_EmptyCommand(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "config.json")

	a := New(testFile, "mcpServers", ReadJSON, WriteJSON, nil)

	err := a.RegisterLocal("local-svc", []string{}, nil, "env", "args")
	if err == nil {
		t.Error("RegisterLocal should fail with empty command")
	}
}
