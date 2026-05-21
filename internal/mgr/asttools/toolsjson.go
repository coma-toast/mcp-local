package asttools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

type fileEntry struct {
	Enabled     bool   `json:"enabled"`
	Tier        string `json:"tier,omitempty"`
	Description string `json:"description,omitempty"`
}

func DefaultToolsConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".astcache", "tools.json")
}

func ToolsConfigPath(svc config.ServiceConfig) string {
	if p := strings.TrimSpace(svc.ToolsConfigPath); p != "" {
		return config.ExpandPath(p)
	}
	return DefaultToolsConfigPath()
}

func needsToolsConfigEnv(svc config.ServiceConfig) bool {
	return len(svc.Tools) > 0 || strings.TrimSpace(svc.ToolsConfigPath) != ""
}

func OverridesCount(tools []config.ToolConfig) int {
	return len(overridesFromTools(tools))
}

func MarshalOverrides(tools []config.ToolConfig) ([]byte, error) {
	return json.MarshalIndent(overridesFromTools(tools), "", "  ")
}

func overridesFromTools(tools []config.ToolConfig) map[string]fileEntry {
	out := make(map[string]fileEntry)
	for _, t := range tools {
		if t.Name == "" {
			continue
		}
		if !t.Enabled || t.Tier != "" || t.Description != "" {
			e := fileEntry{Enabled: t.Enabled}
			if t.Tier != "" {
				e.Tier = t.Tier
			}
			if t.Description != "" {
				e.Description = t.Description
			}
			out[t.Name] = e
		}
	}
	return out
}

func WriteToolsJSON(path string, tools []config.ToolConfig) error {
	overrides := overridesFromTools(tools)
	if len(overrides) == 0 {
		if _, err := os.Stat(path); err == nil {
			return os.Remove(path)
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(overrides, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}

func ApplyStartEnv(env []string, svc config.ServiceConfig) ([]string, error) {
	if svc.ActiveTier != "" {
		env = append(env, fmt.Sprintf("AST_MCP_TIER=%s", svc.ActiveTier))
	}
	if svc.CodeMode {
		env = append(env, "AST_MCP_CODE_MODE=true")
	}
	if !needsToolsConfigEnv(svc) {
		return env, nil
	}
	path := ToolsConfigPath(svc)
	if len(svc.Tools) > 0 {
		if err := WriteToolsJSON(path, svc.Tools); err != nil {
			return env, fmt.Errorf("write tools config: %w", err)
		}
	}
	env = append(env, fmt.Sprintf("AST_MCP_TOOLS_CONFIG=%s", path))
	return env, nil
}
