package claudedesktop

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/jsonagent"
)

func TestConfigPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	path := ConfigPath()

	var expected string
	switch runtime.GOOS {
	case "darwin":
		expected = filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json")
	case "linux":
		expected = filepath.Join(home, ".config", "Claude", "claude_desktop_config.json")
	case "windows":
		expected = filepath.Join(home, "AppData", "Roaming", "Claude", "claude_desktop_config.json")
	default:
		expected = filepath.Join(home, ".config", "Claude", "claude_desktop_config.json")
	}

	if path != expected {
		t.Errorf("ConfigPath() = %q, want %q", path, expected)
	}
}

func setupHome(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	orig := executable
	executable = func() (string, error) { return "/opt/bin/mcp-local", nil }
	t.Cleanup(func() { executable = orig })
}

func TestEntryToClaude_RemoteIsBridge(t *testing.T) {
	setupHome(t)
	result := entryToClaude(config.AgentEntry{Type: "remote", URL: "http://localhost:8080/mcp"})
	if result["command"] != "/opt/bin/mcp-local" {
		t.Errorf("command = %v, want /opt/bin/mcp-local", result["command"])
	}
	args, _ := result["args"].([]string)
	if len(args) != 2 || args[0] != "bridge" || args[1] != "http://localhost:8080/mcp" {
		t.Errorf("args = %v, want [bridge http://localhost:8080/mcp]", result["args"])
	}
	if _, ok := result["url"]; ok {
		t.Error("bridge entry must not carry url")
	}
}

func TestBridgeCommand_NonMCPLocalExecutable(t *testing.T) {
	orig := executable
	defer func() { executable = orig }()
	executable = func() (string, error) { return "/tmp/go-build/claudedesktop.test", nil }
	if got := BridgeCommand(); got != "mcp-local" {
		t.Errorf("BridgeCommand = %q, want mcp-local", got)
	}
}

func TestRegisterServices_WritesBridgeEntry(t *testing.T) {
	setupHome(t)
	path := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	orig := "{\n  \"globalShortcut\": \"Cmd+Space\",\n  \"mcpServers\": {\n    \"foreign\": {\"command\": \"npx\"}\n  }\n}\n"
	if err := os.WriteFile(path, []byte(orig), 0o644); err != nil {
		t.Fatal(err)
	}
	svcs := []config.ServiceConfig{{Name: "ast", MCPType: "http", Port: 7821, MCPURL: "http://localhost:7821/mcp"}}
	if err := RegisterServices(svcs); err != nil {
		t.Fatal(err)
	}
	top, err := jsonagent.ReadJSON(path)
	if err != nil {
		t.Fatal(err)
	}
	servers := top["mcpServers"].(map[string]interface{})
	entry := servers["ast"].(map[string]interface{})
	if entry["command"] != "/opt/bin/mcp-local" {
		t.Errorf("command = %v", entry["command"])
	}
	if args, _ := entry["args"].([]interface{}); len(args) != 2 || args[0] != "bridge" || args[1] != "http://localhost:7821/mcp" {
		t.Errorf("args = %v", entry["args"])
	}
	if servers["foreign"] == nil || top["globalShortcut"] != "Cmd+Space" {
		t.Errorf("foreign config lost: %v", top)
	}
	found, err := Deregister("ast")
	if err != nil || !found {
		t.Fatalf("Deregister = %v, %v", found, err)
	}
	if after, _ := os.ReadFile(path); string(after) != orig {
		t.Errorf("deregister did not restore original:\n%s", after)
	}
	if found, err := Deregister("ast"); err != nil || found {
		t.Errorf("second Deregister = %v, %v; want false, nil", found, err)
	}
}

func TestDeregister_MissingFile(t *testing.T) {
	setupHome(t)
	found, err := Deregister("nothing")
	if err != nil || found {
		t.Fatalf("Deregister = %v, %v", found, err)
	}
}

func TestEntryToClaude_Stdio(t *testing.T) {
	entry := config.AgentEntry{
		Type:        "local",
		Command:     []string{"/usr/bin/mcp-server", "--config", "cfg.yaml"},
		Environment: map[string]string{"KEY": "VAL"},
	}
	result := entryToClaude(entry)

	if result["type"] != "stdio" {
		t.Errorf("type = %q, want stdio", result["type"])
	}
	if result["command"] != "/usr/bin/mcp-server" {
		t.Errorf("command = %q, want /usr/bin/mcp-server", result["command"])
	}
	// args is []string, not []interface{}
	args, ok := result["args"].([]string)
	if !ok {
		t.Errorf("args should be []string, got %T", result["args"])
	} else if len(args) != 2 || args[0] != "--config" || args[1] != "cfg.yaml" {
		t.Errorf("args = %v, want [--config cfg.yaml]", args)
	}
	// env is map[string]string, not map[string]interface{}
	env, ok := result["env"].(map[string]string)
	if !ok {
		t.Errorf("env should be map[string]string, got %T", result["env"])
	} else if env["KEY"] != "VAL" {
		t.Error("env not set correctly")
	}
}

func TestEntryToClaude_StdioNoArgs(t *testing.T) {
	entry := config.AgentEntry{
		Type:    "local",
		Command: []string{"/usr/bin/mcp-server"},
	}
	result := entryToClaude(entry)

	if result["type"] != "stdio" {
		t.Errorf("type = %q, want stdio", result["type"])
	}
	if _, ok := result["args"]; ok {
		t.Error("args should not be present when empty")
	}
	if _, ok := result["env"]; ok {
		t.Error("env should not be present when empty")
	}
}
