package opencode

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/jsonagent"
)

var stripLineComments = regexp.MustCompile(`(?m)^\s*//.*$`)

func ConfigPath() string {
	home, _ := os.UserHomeDir()
	j := filepath.Join(home, ".config", "opencode", "opencode.jsonc")
	if _, err := os.Stat(j); err == nil {
		return j
	}
	return filepath.Join(home, ".config", "opencode", "opencode.json")
}

func readFn(path string) (map[string]interface{}, error) {
	return jsonagent.ReadJSONC(path, func(raw []byte) []byte {
		return stripLineComments.ReplaceAll(raw, []byte{})
	})
}

func newAgent() jsonagent.Agent {
	return jsonagent.New(ConfigPath(), "mcp", readFn, jsonagent.WriteJSON, entryToOpenCode)
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
	m, err := readFn(ConfigPath())
	if os.IsNotExist(err) {
		m = map[string]interface{}{}
	} else if err != nil {
		return err
	}
	block := newAgent().EnsureBlock(m)
	entry := map[string]interface{}{
		"type":    "local",
		"command": command,
		"enabled": true,
	}
	if len(env) > 0 {
		entry["environment"] = env
	}
	block[name] = entry
	m["mcp"] = block
	return jsonagent.WriteJSON(ConfigPath(), m)
}

func Deregister(name string) (bool, error) {
	return newAgent().Deregister(name)
}

func RegisterServices(services []config.ServiceConfig) error {
	return newAgent().RegisterServices(services)
}
