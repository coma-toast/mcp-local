package claudedesktop

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/jsonagent"
)

func ConfigPath() string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json")
	case "linux":
		return filepath.Join(home, ".config", "Claude", "claude_desktop_config.json")
	case "windows":
		return filepath.Join(home, "AppData", "Roaming", "Claude", "claude_desktop_config.json")
	default:
		return filepath.Join(home, ".config", "Claude", "claude_desktop_config.json")
	}
}

var agent = jsonagent.New(ConfigPath(), "mcpServers", jsonagent.ReadJSON, jsonagent.WriteJSON, entryToClaude)

func entryToClaude(entry config.AgentEntry) map[string]interface{} {
	if entry.Type == "remote" {
		return map[string]interface{}{
			"type": "http",
			"url":  entry.URL,
		}
	}
	obj := map[string]interface{}{
		"type":    "stdio",
		"command": entry.Command[0],
	}
	if len(entry.Command) > 1 {
		obj["args"] = entry.Command[1:]
	}
	if len(entry.Environment) > 0 {
		obj["env"] = entry.Environment
	}
	return obj
}

func RegisterServices(services []config.ServiceConfig) error {
	return agent.RegisterServices(services)
}

func Deregister(name string) error {
	_, err := agent.Deregister(name)
	return err
}
