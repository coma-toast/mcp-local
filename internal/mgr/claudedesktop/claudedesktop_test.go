package claudedesktop

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
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

func TestEntryToClaude_Remote(t *testing.T) {
	entry := config.AgentEntry{
		Type: "remote",
		URL:  "http://localhost:8080/mcp",
	}
	result := entryToClaude(entry)

	if result["type"] != "http" {
		t.Errorf("type = %q, want http", result["type"])
	}
	if result["url"] != "http://localhost:8080/mcp" {
		t.Errorf("url = %q, want http://localhost:8080/mcp", result["url"])
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
