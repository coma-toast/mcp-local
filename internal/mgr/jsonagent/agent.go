package jsonagent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

type Reader func(path string) (map[string]interface{}, error)
type Writer func(path string, m map[string]interface{}) error
type EntryConverter func(entry config.AgentEntry) map[string]interface{}

type Agent struct {
	configPath   string
	blockKey     string
	read         Reader
	write        Writer
	convertEntry EntryConverter
}

func New(configPath string, blockKey string, read Reader, write Writer, convert EntryConverter) Agent {
	return Agent{
		configPath:   configPath,
		blockKey:     blockKey,
		read:         read,
		write:        write,
		convertEntry: convert,
	}
}

func ReadJSON(path string) (map[string]interface{}, error) {
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

func ReadJSONC(path string, stripComments func([]byte) []byte) (map[string]interface{}, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	stripped := stripComments(raw)
	var m map[string]interface{}
	if err := json.Unmarshal(stripped, &m); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return m, nil
}

func WriteJSON(path string, m map[string]interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0644)
}

func (a Agent) EnsureBlock(m map[string]interface{}) map[string]interface{} {
	block, _ := m[a.blockKey].(map[string]interface{})
	if block == nil {
		block = map[string]interface{}{}
		m[a.blockKey] = block
	}
	return block
}

func (a Agent) RegisterServices(services []config.ServiceConfig) error {
	m, err := a.read(a.configPath)
	if os.IsNotExist(err) {
		m = map[string]interface{}{}
	} else if err != nil {
		return err
	}
	block := a.EnsureBlock(m)
	for _, s := range services {
		entry, err := config.ServiceToEntry(s)
		if err != nil {
			return fmt.Errorf("%s: %w", s.Name, err)
		}
		block[s.Name] = a.convertEntry(entry)
	}
	m[a.blockKey] = block
	return a.write(a.configPath, m)
}

func (a Agent) Deregister(name string) (bool, error) {
	m, err := a.read(a.configPath)
	if err != nil {
		return false, err
	}
	block := a.EnsureBlock(m)
	if _, ok := block[name]; !ok {
		return false, nil
	}
	delete(block, name)
	m[a.blockKey] = block
	return true, a.write(a.configPath, m)
}

func (a Agent) RegisterRemote(name, url string, extra map[string]interface{}) error {
	m, err := a.read(a.configPath)
	if os.IsNotExist(err) {
		m = map[string]interface{}{}
	} else if err != nil {
		return err
	}
	block := a.EnsureBlock(m)
	entry := map[string]interface{}{"url": url}
	for k, v := range extra {
		entry[k] = v
	}
	if existing, ok := block[name].(map[string]interface{}); ok {
		if existing["url"] == url {
			return nil
		}
	}
	block[name] = entry
	m[a.blockKey] = block
	return a.write(a.configPath, m)
}

func (a Agent) RegisterLocal(name string, command []string, env map[string]string, envKey, argsKey string) error {
	if len(command) == 0 {
		return fmt.Errorf("empty command for %q", name)
	}
	m, err := a.read(a.configPath)
	if os.IsNotExist(err) {
		m = map[string]interface{}{}
	} else if err != nil {
		return err
	}
	block := a.EnsureBlock(m)
	entry := map[string]interface{}{"command": command[0]}
	if len(command) > 1 {
		entry[argsKey] = command[1:]
	}
	if len(env) > 0 {
		entry[envKey] = env
	}
	block[name] = entry
	m[a.blockKey] = block
	return a.write(a.configPath, m)
}
