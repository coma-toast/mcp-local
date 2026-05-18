package cursor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cursor", "mcp.json")
}

func readConfig(path string) (map[string]interface{}, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return m, nil
}

func writeConfig(path string, m map[string]interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0644)
}

func ensureMCPServers(m map[string]interface{}) map[string]interface{} {
	block, _ := m["mcpServers"].(map[string]interface{})
	if block == nil {
		block = map[string]interface{}{}
		m["mcpServers"] = block
	}
	return block
}

func entryToCursor(entry config.AgentEntry) map[string]interface{} {
	if entry.Type == "remote" {
		obj := map[string]interface{}{
			"url": entry.URL,
		}
		return obj
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

func RegisterServices(services []config.ServiceConfig) error {
	path := ConfigPath()
	m, err := readConfig(path)
	if os.IsNotExist(err) {
		m = map[string]interface{}{}
	} else if err != nil {
		return err
	}
	block := ensureMCPServers(m)
	for _, s := range services {
		entry, err := config.ServiceToEntry(s)
		if err != nil {
			return fmt.Errorf("%s: %w", s.Name, err)
		}
		block[s.Name] = entryToCursor(entry)
	}
	m["mcpServers"] = block
	return writeConfig(path, m)
}

func Deregister(name string) error {
	path := ConfigPath()
	m, err := readConfig(path)
	if err != nil {
		return err
	}
	block, _ := m["mcpServers"].(map[string]interface{})
	if block == nil {
		return nil
	}
	if _, ok := block[name]; !ok {
		return nil
	}
	delete(block, name)
	m["mcpServers"] = block
	return writeConfig(path, m)
}

func RegisterRemote(name, url string) error {
	path := ConfigPath()
	m, err := readConfig(path)
	if os.IsNotExist(err) {
		m = map[string]interface{}{}
	} else if err != nil {
		return err
	}
	block := ensureMCPServers(m)
	block[name] = map[string]interface{}{"url": url}
	m["mcpServers"] = block
	return writeConfig(path, m)
}

func RegisterLocal(name string, command []string, env map[string]string) error {
	if len(command) == 0 {
		return fmt.Errorf("empty command for %q", name)
	}
	path := ConfigPath()
	m, err := readConfig(path)
	if os.IsNotExist(err) {
		m = map[string]interface{}{}
	} else if err != nil {
		return err
	}
	block := ensureMCPServers(m)
	entry := map[string]interface{}{
		"command": command[0],
	}
	if len(command) > 1 {
		entry["args"] = command[1:]
	}
	if len(env) > 0 {
		entry["env"] = env
	}
	block[name] = entry
	m["mcpServers"] = block
	return writeConfig(path, m)
}

func IsHTTPService(s config.ServiceConfig) bool {
	t := strings.ToLower(strings.TrimSpace(s.MCPType))
	return t == "http" || (s.Port > 0 && t != "stdio")
}
