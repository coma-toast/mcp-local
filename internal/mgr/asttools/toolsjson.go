package asttools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/safefile"
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

// toolsScope is the ledger scope for tools.json entries mcp-local wrote.
const toolsScope = "tools"

// WriteToolsJSON merges the config's tool overrides into tools.json. Entries for tools the
// config doesn't override are left alone (e.g. ones set in the ast dashboard), unless
// mcp-local wrote them earlier, in which case they are removed so the tool is back to its
// default. Keys other than enabled/tier/description are kept. The file is backed up and
// written atomically, and never deleted.
func WriteToolsJSON(path string, tools []config.ToolConfig) error {
	overrides := overridesFromTools(tools)
	current, err := os.ReadFile(safefile.ResolvePath(path))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	entries, before := map[string]map[string]interface{}{}, map[string]map[string]interface{}{}
	if len(bytes.TrimSpace(current)) > 0 {
		if err := json.Unmarshal(current, &entries); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		_ = json.Unmarshal(current, &before)
	}
	var add, remove []string
	for _, t := range tools {
		if t.Name == "" {
			continue
		}
		if o, ok := overrides[t.Name]; ok {
			entries[t.Name] = withOverride(entries[t.Name], o)
			add = append(add, t.Name)
		} else if _, ok := entries[t.Name]; ok && safefile.Recorded(path, toolsScope, t.Name) {
			delete(entries, t.Name)
			remove = append(remove, t.Name)
		}
	}
	if !reflect.DeepEqual(entries, before) {
		data, err := json.MarshalIndent(entries, "", "  ")
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := safefile.Replace(path, append(data, '\n')); err != nil {
			return err
		}
	}
	return safefile.Record(path, toolsScope, add, remove)
}

func withOverride(entry map[string]interface{}, o fileEntry) map[string]interface{} {
	if entry == nil {
		entry = map[string]interface{}{}
	}
	entry["enabled"] = o.Enabled
	for k, v := range map[string]string{"tier": o.Tier, "description": o.Description} {
		if v == "" {
			delete(entry, k)
		} else {
			entry[k] = v
		}
	}
	return entry
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
