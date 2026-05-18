package opencode

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
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

func readJSONC(path string) (map[string]interface{}, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	stripped := stripLineComments.ReplaceAll(raw, []byte{})
	var m map[string]interface{}
	if err := json.Unmarshal(stripped, &m); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return m, nil
}

func writeJSONC(path string, m map[string]interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0644)
}

func ensureMCPBlock(m map[string]interface{}) map[string]interface{} {
	block, _ := m["mcp"].(map[string]interface{})
	if block == nil {
		block = map[string]interface{}{}
		m["mcp"] = block
	}
	return block
}

func entryToOpenCode(entry config.AgentEntry) map[string]interface{} {
	obj := map[string]interface{}{
		"type":    entry.Type,
		"enabled": true,
	}
	if entry.Type == "remote" {
		obj["url"] = entry.URL
		obj["timeout"] = entry.Timeout
	} else {
		obj["command"] = entry.Command
		if len(entry.Environment) > 0 {
			obj["environment"] = entry.Environment
		}
	}
	return obj
}

func RegisterRemote(name, url string, timeoutMS int) error {
	path := ConfigPath()
	m, err := readJSONC(path)
	if os.IsNotExist(err) {
		m = map[string]interface{}{}
	} else if err != nil {
		return err
	}
	block := ensureMCPBlock(m)
	entry := map[string]interface{}{
		"type":    "remote",
		"url":     url,
		"enabled": true,
	}
	if timeoutMS <= 0 {
		timeoutMS = 30000
	}
	entry["timeout"] = timeoutMS

	if existing, ok := block[name].(map[string]interface{}); ok {
		if existing["url"] == url {
			return nil
		}
	}
	block[name] = entry
	m["mcp"] = block
	return writeJSONC(path, m)
}

func RegisterLocal(name string, command []string, env map[string]string) error {
	if len(command) == 0 {
		return fmt.Errorf("empty command for %q", name)
	}
	path := ConfigPath()
	m, err := readJSONC(path)
	if os.IsNotExist(err) {
		m = map[string]interface{}{}
	} else if err != nil {
		return err
	}
	block := ensureMCPBlock(m)
	entry := map[string]interface{}{
		"type":    "local",
		"command": command,
		"enabled": true,
	}
	if len(env) > 0 {
		envObj := map[string]interface{}{}
		for k, v := range env {
			envObj[k] = v
		}
		entry["environment"] = envObj
	}
	block[name] = entry
	m["mcp"] = block
	return writeJSONC(path, m)
}

func Deregister(name string) (bool, error) {
	path := ConfigPath()
	m, err := readJSONC(path)
	if err != nil {
		return false, err
	}
	block, _ := m["mcp"].(map[string]interface{})
	if block == nil {
		return false, nil
	}
	if _, ok := block[name]; !ok {
		return false, nil
	}
	delete(block, name)
	m["mcp"] = block
	return true, writeJSONC(path, m)
}

func RegisterServices(services []config.ServiceConfig) error {
	path := ConfigPath()
	m, err := readJSONC(path)
	if os.IsNotExist(err) {
		m = map[string]interface{}{}
	} else if err != nil {
		return err
	}
	block := ensureMCPBlock(m)
	for _, s := range services {
		entry, err := config.ServiceToEntry(s)
		if err != nil {
			return fmt.Errorf("%s: %w", s.Name, err)
		}
		block[s.Name] = entryToOpenCode(entry)
	}
	m["mcp"] = block
	return writeJSONC(path, m)
}
