// Package claudedesktop registers services in Claude Desktop's claude_desktop_config.json.
//
// Claude Desktop only launches stdio servers, so HTTP services are written as a stdio
// entry that runs `mcp-local bridge <mcp_url>`.
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
	case "windows":
		return filepath.Join(home, "AppData", "Roaming", "Claude", "claude_desktop_config.json")
	default:
		return filepath.Join(home, ".config", "Claude", "claude_desktop_config.json")
	}
}

// executable is swappable in tests.
var executable = os.Executable

// BridgeCommand is the command Claude Desktop runs for HTTP services: the absolute path of
// the running mcp-local binary when it is one, else "mcp-local" from PATH.
func BridgeCommand() string {
	exe, err := executable()
	if err != nil {
		return "mcp-local"
	}
	if real, err := filepath.EvalSymlinks(exe); err == nil {
		exe = real
	}
	if filepath.Base(exe) != "mcp-local" {
		return "mcp-local"
	}
	return exe
}

func newAgent() jsonagent.Agent {
	return jsonagent.New(ConfigPath(), "mcpServers", entryToClaude)
}

func entryToClaude(entry config.AgentEntry) map[string]interface{} {
	if entry.Type == "remote" {
		return map[string]interface{}{
			"command": BridgeCommand(),
			"args":    []string{"bridge", entry.URL},
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
	return newAgent().RegisterServices(services)
}

// Deregister removes the named entry; found is false when nothing was there.
func Deregister(name string) (bool, error) {
	return newAgent().Deregister(name)
}
