package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/coma-toast/mcp-local/internal/mcpserver"
	"github.com/coma-toast/mcp-local/internal/mgr/agents"
	"github.com/coma-toast/mcp-local/internal/mgr/asttools"
	"github.com/coma-toast/mcp-local/internal/mgr/claudedesktop"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/cursor"
	"github.com/coma-toast/mcp-local/internal/mgr/opencode"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
	"github.com/coma-toast/mcp-local/internal/mgr/runtime"
	"github.com/spf13/cobra"
)

func addMCP() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "mcp",
		Short: "Run mcp-local as a stdio MCP server",
		Long:  "Exposes mcp-local service, agent registration, and tool-tier management over stdio JSON-RPC 2.0.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Keep stdout exclusively for JSON-RPC frames; redirect stray fmt.Print* to stderr.
			mcpOut := os.Stdout
			os.Stdout = os.Stderr
			agents.Warn = os.Stderr
			defer func() { os.Stdout = mcpOut }()

			ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer cancel()
			return newControlMCPServer().Serve(ctx, os.Stdin, mcpOut)
		},
	})
}

func newControlMCPServer() *mcpserver.Server {
	s := mcpserver.New("mcp-local", version)
	registerControlTools(s)
	return s
}

func registerControlTools(s *mcpserver.Server) {
	s.Register(mcpserver.ToolDef{
		Name:        "list_services",
		Description: "List all configured MCP services with mode, port, mcp_url, active_tier, and live runtime status.",
		InputSchema: map[string]any{"type": "object", "properties": map[string]any{}},
	}, toolListServices)

	s.Register(mcpserver.ToolDef{
		Name:        "service_status",
		Description: "Get runtime status and per-agent registration status for a configured service.",
		InputSchema: nameSchema("Service name"),
	}, toolServiceStatus)

	s.Register(mcpserver.ToolDef{
		Name:        "start_service",
		Description: "Start a service and register it with enabled agents.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":  map[string]any{"type": "string", "description": "Service name"},
				"build": map[string]any{"type": "boolean", "description": "Run build_command before starting (default true)"},
			},
			"required": []string{"name"},
		},
	}, toolStartService)

	s.Register(mcpserver.ToolDef{
		Name:        "stop_service",
		Description: "Stop a running service and optionally deregister it from agents.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":       map[string]any{"type": "string", "description": "Service name"},
				"deregister": map[string]any{"type": "boolean", "description": "Remove agent registrations (default true)"},
			},
			"required": []string{"name"},
		},
	}, toolStopService)

	s.Register(mcpserver.ToolDef{
		Name:        "restart_service",
		Description: "Restart a service (rebuilds first unless no_build_on_start is set) and re-registers with agents.",
		InputSchema: nameSchema("Service name"),
	}, toolRestartService)

	s.Register(mcpserver.ToolDef{
		Name:        "register_service",
		Description: "Register a service with enabled agents (OpenCode, Cursor, Claude Desktop).",
		InputSchema: nameSchema("Service name"),
	}, toolRegisterService)

	s.Register(mcpserver.ToolDef{
		Name:        "deregister_service",
		Description: "Remove a service's entries from OpenCode, Cursor, and Claude Desktop.",
		InputSchema: nameSchema("Service name"),
	}, toolDeregisterService)

	s.Register(mcpserver.ToolDef{
		Name:        "list_tools",
		Description: "List configured tools and effective tiers for a service, or fetch live tools/list from a running service.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{"type": "string", "description": "Service name"},
				"live": map[string]any{"type": "boolean", "description": "Query live tools/list from the running MCP endpoint"},
			},
			"required": []string{"name"},
		},
	}, toolListTools)

	s.Register(mcpserver.ToolDef{
		Name:        "sync_tools",
		Description: "Fetch tools/list from a running service and merge into config.yaml, keeping existing user overrides.",
		InputSchema: nameSchema("Service name"),
	}, toolSyncTools)

	s.Register(mcpserver.ToolDef{
		Name:        "configure_tools",
		Description: "Update a service's active_tier, code_mode, or per-tool enabled/tier/description override and write tools.json.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":        map[string]any{"type": "string", "description": "Service name"},
				"active_tier": map[string]any{"type": "string", "description": "core, extended, or complete"},
				"code_mode":   map[string]any{"type": "boolean", "description": "Set code_mode (true) or no_code_mode (false)"},
				"tool":        map[string]any{"type": "string", "description": "Tool name to add/update"},
				"enabled":     map[string]any{"type": "boolean", "description": "Enable or disable the specified tool"},
				"tier":        map[string]any{"type": "string", "description": "Tier override for the specified tool"},
				"description": map[string]any{"type": "string", "description": "Description override for the specified tool"},
				"apply":       map[string]any{"type": "boolean", "description": "Write tools.json overrides (default true)"},
			},
			"required": []string{"name"},
		},
	}, toolConfigureTools)

	s.Register(mcpserver.ToolDef{
		Name:        "add_service",
		Description: "Add or update a service in config.yaml.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":          map[string]any{"type": "string", "description": "Service name"},
				"type":          map[string]any{"type": "string", "description": "http, stdio, or filesystem"},
				"command":       map[string]any{"type": "string", "description": "Binary path"},
				"args":          map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Command arguments"},
				"port":          map[string]any{"type": "integer", "description": "TCP port"},
				"mcp_url":       map[string]any{"type": "string", "description": "MCP endpoint URL"},
				"health_url":    map[string]any{"type": "string", "description": "Health check URL"},
				"dashboard_url": map[string]any{"type": "string", "description": "Dashboard URL"},
				"path":          map[string]any{"type": "string", "description": "Filesystem root path (for type=filesystem)"},
				"build_command": map[string]any{"type": "string", "description": "Shell build command"},
				"env":           map[string]any{"type": "object", "additionalProperties": map[string]any{"type": "string"}, "description": "Environment variables"},
			},
			"required": []string{"name"},
		},
	}, toolAddService)

	s.Register(mcpserver.ToolDef{
		Name:        "remove_service",
		Description: "Deregister a service from all agents and remove it from config.yaml.",
		InputSchema: nameSchema("Service name"),
	}, toolRemoveService)
}

func nameSchema(desc string) map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{"type": "string", "description": desc},
		},
		"required": []string{"name"},
	}
}

func decodeArgs(raw json.RawMessage, dst any) error {
	if len(raw) == 0 {
		return nil
	}
	return json.Unmarshal(raw, dst)
}

func jsonText(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func toolListServices(_ context.Context, _ json.RawMessage) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	type item struct {
		Name       string `json:"name"`
		Mode       string `json:"mode"`
		State      string `json:"state"`
		Detail     string `json:"detail,omitempty"`
		PID        int    `json:"pid,omitempty"`
		Port       int    `json:"port,omitempty"`
		MCPURL     string `json:"mcp_url,omitempty"`
		ActiveTier string `json:"active_tier,omitempty"`
		ToolsCount int    `json:"tools_count"`
	}
	var out []item
	for _, name := range cfg.SortedNames() {
		svc, _ := cfg.ServiceNamed(name)
		st := runtime.CheckStatus(svc)
		out = append(out, item{
			Name:       svc.Name,
			Mode:       process.ModeOf(svc).String(),
			State:      string(st.State),
			Detail:     st.Detail,
			PID:        st.PID,
			Port:       svc.Port,
			MCPURL:     svc.MCPEndpoint(),
			ActiveTier: svc.ActiveTier,
			ToolsCount: len(svc.Tools),
		})
	}
	return jsonText(map[string]any{"services": out})
}

func toolServiceStatus(_ context.Context, raw json.RawMessage) (string, error) {
	var a struct {
		Name string `json:"name"`
	}
	if err := decodeArgs(raw, &a); err != nil {
		return "", err
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	svc, err := cfg.ServiceNamed(a.Name)
	if err != nil {
		return "", err
	}
	st := runtime.CheckStatus(svc)
	return jsonText(map[string]any{
		"name":   svc.Name,
		"mode":   process.ModeOf(svc).String(),
		"state":  string(st.State),
		"detail": st.Detail,
		"pid":    st.PID,
		"port":   svc.Port,
		"stale":  st.Stale,
		"registered": map[string]bool{
			"opencode": isRegisteredInFile(opencode.ConfigPath(), svc.Name),
			"cursor":   isRegisteredInFile(cursor.ConfigPath(), svc.Name),
			"claude":   isRegisteredInFile(claudedesktop.ConfigPath(), svc.Name),
		},
	})
}

func toolStartService(_ context.Context, raw json.RawMessage) (string, error) {
	var a struct {
		Name  string `json:"name"`
		Build *bool  `json:"build"`
	}
	if err := decodeArgs(raw, &a); err != nil {
		return "", err
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	if _, err := cfg.ServiceNamed(a.Name); err != nil {
		return "", err
	}
	build := true
	if a.Build != nil {
		build = *a.Build
	}
	started, failed := startServices(cfg, []string{a.Name}, build, os.Stderr)
	if failed > 0 {
		return "", fmt.Errorf("failed to start %s", a.Name)
	}
	registerStarted(cfg, started)
	return jsonText(map[string]any{"started": a.Name, "status": "ok"})
}

func toolStopService(_ context.Context, raw json.RawMessage) (string, error) {
	var a struct {
		Name       string `json:"name"`
		Deregister *bool  `json:"deregister"`
	}
	if err := decodeArgs(raw, &a); err != nil {
		return "", err
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	deregister := true
	if a.Deregister != nil {
		deregister = *a.Deregister
	}
	if err := stopServiceNamed(cfg, a.Name, deregister); err != nil {
		return "", err
	}
	return jsonText(map[string]any{"stopped": a.Name, "deregistered": deregister})
}

func toolRestartService(_ context.Context, raw json.RawMessage) (string, error) {
	var a struct {
		Name string `json:"name"`
	}
	if err := decodeArgs(raw, &a); err != nil {
		return "", err
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	if _, err := cfg.ServiceNamed(a.Name); err != nil {
		return "", err
	}
	started, failed := restartServices(cfg, []string{a.Name}, os.Stderr)
	if failed > 0 {
		return "", fmt.Errorf("failed to restart %s", a.Name)
	}
	registerStarted(cfg, started)
	return jsonText(map[string]any{"restarted": a.Name, "status": "ok"})
}

func toolRegisterService(_ context.Context, raw json.RawMessage) (string, error) {
	var a struct {
		Name string `json:"name"`
	}
	if err := decodeArgs(raw, &a); err != nil {
		return "", err
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	svc, err := cfg.ServiceNamed(a.Name)
	if err != nil {
		return "", err
	}
	targets := agents.TargetsFromConfig(*cfg)
	if err := agents.RegisterAll([]config.ServiceConfig{svc}, targets); err != nil {
		return "", err
	}
	return jsonText(map[string]any{"registered": svc.Name, "hosts": targets.Hosts()})
}

func toolDeregisterService(_ context.Context, raw json.RawMessage) (string, error) {
	var a struct {
		Name string `json:"name"`
	}
	if err := decodeArgs(raw, &a); err != nil {
		return "", err
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	svc, err := cfg.ServiceNamed(a.Name)
	if err != nil {
		return "", err
	}
	targets := agents.TargetsFromConfig(*cfg)
	if err := agents.DeregisterAll(svc, targets); err != nil {
		return "", err
	}
	return jsonText(map[string]any{"deregistered": svc.Name, "hosts": targets.Hosts()})
}

func toolListTools(_ context.Context, raw json.RawMessage) (string, error) {
	var a struct {
		Name string `json:"name"`
		Live bool   `json:"live"`
	}
	if err := decodeArgs(raw, &a); err != nil {
		return "", err
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	svc, err := cfg.ServiceNamed(a.Name)
	if err != nil {
		return "", err
	}
	if a.Live {
		mcpURL := strings.TrimSpace(svc.MCPEndpoint())
		if mcpURL == "" {
			return "", fmt.Errorf("no mcp_url or port for %s", a.Name)
		}
		liveTools, err := fetchToolsList(mcpURL)
		if err != nil {
			return "", err
		}
		return jsonText(map[string]any{"name": svc.Name, "source": "live", "tools": liveTools})
	}
	type toolOut struct {
		Name          string `json:"name"`
		Enabled       bool   `json:"enabled"`
		Tier          string `json:"tier,omitempty"`
		EffectiveTier string `json:"effective_tier,omitempty"`
		Description   string `json:"description,omitempty"`
	}
	var out []toolOut
	for _, t := range svc.Tools {
		out = append(out, toolOut{
			Name:          t.Name,
			Enabled:       t.Enabled,
			Tier:          t.Tier,
			EffectiveTier: asttools.EffectiveTierForDisplay(t.Name, t.Tier),
			Description:   t.Description,
		})
	}
	return jsonText(map[string]any{
		"name":              svc.Name,
		"source":            "config",
		"active_tier":       svc.ActiveTier,
		"code_mode":         svc.CodeMode,
		"no_code_mode":      svc.NoCodeMode,
		"tools_config_path": asttools.ToolsConfigPath(svc),
		"tools":             out,
	})
}

func toolSyncTools(_ context.Context, raw json.RawMessage) (string, error) {
	var a struct {
		Name string `json:"name"`
	}
	if err := decodeArgs(raw, &a); err != nil {
		return "", err
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	svc, err := cfg.ServiceNamed(a.Name)
	if err != nil {
		return "", err
	}
	mcpURL := strings.TrimSpace(svc.MCPEndpoint())
	if mcpURL == "" {
		return "", fmt.Errorf("no mcp_url or port for %s", a.Name)
	}
	tools, err := fetchToolsList(mcpURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch tools: %w", err)
	}
	existing := make(map[string]config.ToolConfig)
	for _, t := range svc.Tools {
		existing[t.Name] = t
	}
	merged := mergeSyncedTools(tools, existing)
	for i := range cfg.Services {
		if cfg.Services[i].Name == a.Name {
			cfg.Services[i].Tools = merged
			break
		}
	}
	if err := config.SaveConfig(cfg); err != nil {
		return "", err
	}
	return jsonText(map[string]any{"name": a.Name, "synced": len(merged), "tools": merged})
}

func toolConfigureTools(_ context.Context, raw json.RawMessage) (string, error) {
	var a struct {
		Name        string  `json:"name"`
		ActiveTier  *string `json:"active_tier"`
		CodeMode    *bool   `json:"code_mode"`
		Tool        string  `json:"tool"`
		Enabled     *bool   `json:"enabled"`
		Tier        *string `json:"tier"`
		Description *string `json:"description"`
		Apply       *bool   `json:"apply"`
	}
	if err := decodeArgs(raw, &a); err != nil {
		return "", err
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	idx := -1
	for i := range cfg.Services {
		if cfg.Services[i].Name == a.Name {
			idx = i
			break
		}
	}
	if idx < 0 {
		return "", fmt.Errorf("service %q not found", a.Name)
	}
	svc := &cfg.Services[idx]
	if a.ActiveTier != nil {
		svc.ActiveTier = *a.ActiveTier
	}
	if a.CodeMode != nil {
		svc.CodeMode = *a.CodeMode
		svc.NoCodeMode = !*a.CodeMode
	}
	if a.Tool != "" {
		found := false
		for i := range svc.Tools {
			if svc.Tools[i].Name == a.Tool {
				if a.Enabled != nil {
					svc.Tools[i].Enabled = *a.Enabled
				}
				if a.Tier != nil {
					svc.Tools[i].Tier = *a.Tier
				}
				if a.Description != nil {
					svc.Tools[i].Description = *a.Description
				}
				found = true
				break
			}
		}
		if !found {
			en := true
			if a.Enabled != nil {
				en = *a.Enabled
			}
			tc := config.ToolConfig{Name: a.Tool, Enabled: en}
			if a.Tier != nil {
				tc.Tier = *a.Tier
			}
			if a.Description != nil {
				tc.Description = *a.Description
			}
			svc.Tools = append(svc.Tools, tc)
		}
	}
	if err := config.SaveConfig(cfg); err != nil {
		return "", err
	}
	apply := true
	if a.Apply != nil {
		apply = *a.Apply
	}
	var toolsJSONPath string
	if apply {
		toolsJSONPath = asttools.ToolsConfigPath(*svc)
		if err := asttools.WriteToolsJSON(toolsJSONPath, svc.Tools); err != nil {
			return "", err
		}
	}
	return jsonText(map[string]any{
		"name":              svc.Name,
		"active_tier":       svc.ActiveTier,
		"code_mode":         svc.CodeMode,
		"tools_json_path":   toolsJSONPath,
		"overrides_applied": asttools.OverridesCount(svc.Tools),
	})
}

func toolAddService(_ context.Context, raw json.RawMessage) (string, error) {
	var a struct {
		Name         string            `json:"name"`
		Type         string            `json:"type"`
		Command      string            `json:"command"`
		Args         []string          `json:"args"`
		Port         int               `json:"port"`
		MCPURL       string            `json:"mcp_url"`
		HealthURL    string            `json:"health_url"`
		DashboardURL string            `json:"dashboard_url"`
		Path         string            `json:"path"`
		BuildCommand string            `json:"build_command"`
		Env          map[string]string `json:"env"`
	}
	if err := decodeArgs(raw, &a); err != nil {
		return "", err
	}
	if strings.TrimSpace(a.Name) == "" {
		return "", fmt.Errorf("name is required")
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		cfg = &config.ManagerConfig{}
	}
	svc := config.ServiceConfig{
		Name:         a.Name,
		MCPType:      a.Type,
		Command:      a.Command,
		Args:         a.Args,
		Port:         a.Port,
		MCPURL:       a.MCPURL,
		HealthURL:    a.HealthURL,
		DashboardURL: a.DashboardURL,
		Path:         a.Path,
		BuildCommand: a.BuildCommand,
		Env:          a.Env,
	}
	cfg.UpsertService(svc)
	if err := config.SaveConfig(cfg); err != nil {
		return "", err
	}
	return jsonText(map[string]any{"saved": svc.Name, "mode": process.ModeOf(svc).String()})
}

func toolRemoveService(_ context.Context, raw json.RawMessage) (string, error) {
	var a struct {
		Name string `json:"name"`
	}
	if err := decodeArgs(raw, &a); err != nil {
		return "", err
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	svc, err := cfg.ServiceNamed(a.Name)
	if err != nil {
		return "", err
	}
	targets := agents.TargetsFromConfig(*cfg)
	if err := agents.DeregisterAll(svc, targets); err != nil {
		return "", err
	}
	filtered := make([]config.ServiceConfig, 0, len(cfg.Services))
	for _, s := range cfg.Services {
		if s.Name != a.Name {
			filtered = append(filtered, s)
		}
	}
	cfg.Services = filtered
	if err := config.SaveConfig(cfg); err != nil {
		return "", err
	}
	return jsonText(map[string]any{"removed": a.Name})
}
