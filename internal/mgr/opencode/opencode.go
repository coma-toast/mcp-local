package opencode

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

type OpenCodeConfig struct {
	MCPServers map[string]MCPServer `json:"mcpServers"`
}

type MCPServer struct {
	Type     string `json:"type"`
	URL      string `json:"url,omitempty"`
	Command  string `json:"command,omitempty"`
	Args     []string `json:"args,omitempty"`
	Enabled  bool   `json:"enabled"`
	Note     string `json:"note,omitempty"`
}

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "opencode", "opencode.jsonc")
}

func RegisterServices(services []config.ServiceConfig) error {
	path := GetConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read opencode config: %w", err)
	}

	// opencode.jsonc might have comments, we should handle it. 
	// For simplicity, we'll treat it as JSON but a production version would use a JSONC parser.
	var cfg OpenCodeConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse opencode config: %w", err)
	}

	if cfg.MCPServers == nil {
		cfg.MCPServers = make(map[string]MCPServer)
	}

	for _, s := range services {
		server := MCPServer{
			Enabled: true,
		}
		if s.MCPType == "http" {
			server.Type = "remote"
			server.URL = fmt.Sprintf("http://localhost:%d/mcp", s.Port)
		} else {
			server.Type = "stdio"
			server.Command = config.ExpandPath(s.Command)
			server.Args = s.Args
		}
		cfg.MCPServers[s.Name] = server
	}

	newData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal opencode config: %w", err)
	}

	return os.WriteFile(path, newData, 0644)
}
