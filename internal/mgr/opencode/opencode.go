package opencode

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/jsonagent"
)

func ConfigPath() string {
	home, _ := os.UserHomeDir()
	j := filepath.Join(home, ".config", "opencode", "opencode.jsonc")
	if _, err := os.Stat(j); err == nil {
		return j
	}
	return filepath.Join(home, ".config", "opencode", "opencode.json")
}

func newAgent() jsonagent.Agent {
	return jsonagent.New(ConfigPath(), "mcp", entryToOpenCode)
}

func entryToOpenCode(entry config.AgentEntry) map[string]interface{} {
	obj := map[string]interface{}{
		"type":    entry.Type,
		"enabled": true,
	}
	if entry.Type == "remote" {
		obj["url"] = entry.URL
		obj["timeout"] = 30000
	} else {
		obj["command"] = entry.Command
		if len(entry.Environment) > 0 {
			obj["environment"] = entry.Environment
		}
	}
	return obj
}

func RegisterRemote(name, url string, timeoutMS int) error {
	if timeoutMS <= 0 {
		timeoutMS = 30000
	}
	return newAgent().RegisterRemote(name, url, map[string]interface{}{
		"type":    "remote",
		"enabled": true,
		"timeout": timeoutMS,
	})
}

func RegisterLocal(name string, command []string, env map[string]string) error {
	if len(command) == 0 {
		return fmt.Errorf("empty command for %q", name)
	}
	entry := map[string]interface{}{
		"type":    "local",
		"command": command,
		"enabled": true,
	}
	if len(env) > 0 {
		entry["environment"] = env
	}
	return newAgent().SetEntries(map[string]map[string]interface{}{name: entry})
}

func Deregister(name string) (bool, error) {
	return newAgent().Deregister(name)
}

func RegisterServices(services []config.ServiceConfig) error {
	return newAgent().RegisterServices(services)
}
