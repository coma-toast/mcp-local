package cursor

import (
	"os"
	"path/filepath"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/jsonagent"
)

func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cursor", "mcp.json")
}

func newAgent() jsonagent.Agent {
	return jsonagent.New(ConfigPath(), "mcpServers", jsonagent.ReadJSON, jsonagent.WriteJSON, entryToCursor)
}

func entryToCursor(entry config.AgentEntry) map[string]interface{} {
	if entry.Type == "remote" {
		return map[string]interface{}{"url": entry.URL}
	}
	obj := map[string]interface{}{
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

func RegisterRemote(name, url string) error {
	return newAgent().RegisterRemote(name, url, nil)
}

func RegisterLocal(name string, command []string, env map[string]string) error {
	return newAgent().RegisterLocal(name, command, env, "env", "args")
}

func Deregister(name string) error {
	_, err := newAgent().Deregister(name)
	return err
}

func RegisterServices(services []config.ServiceConfig) error {
	return newAgent().RegisterServices(services)
}
