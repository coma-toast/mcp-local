package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type AgentsConfig struct {
	OpenCode bool `yaml:"opencode"`
	Cursor   bool `yaml:"cursor"`
	Claude   bool `yaml:"claude"`
}

type ServiceConfig struct {
	Name            string            `yaml:"name"`
	Command         string            `yaml:"command"`
	Args            []string          `yaml:"args"`
	Env             map[string]string `yaml:"env"`
	Port            int               `yaml:"port"`
	MCPType         string            `yaml:"type"`
	BuildCmd        string            `yaml:"build_cmd"`
	BuildCommand    string            `yaml:"build_command,omitempty"`
	NoBuildOnStart  bool              `yaml:"no_build_on_start,omitempty"`
	Deps            []string          `yaml:"deps"`
	Description     string            `yaml:"description"`
	Tools           []ToolConfig      `yaml:"tools"`
	HealthURL       string            `yaml:"health_url,omitempty"`
	DashboardURL    string            `yaml:"dashboard_url,omitempty"`
	MCPURL          string            `yaml:"mcp_url,omitempty"`
	Log             string            `yaml:"log,omitempty"`
	ActiveTier      string            `yaml:"active_tier,omitempty"`
	CodeMode        bool              `yaml:"code_mode,omitempty"`
	NoCodeMode      bool              `yaml:"no_code_mode,omitempty"`
	ToolsConfigPath string            `yaml:"tools_config_path,omitempty"`
	Path            string            `yaml:"path,omitempty"`
	// Installer delegates agent registration to the service's own CLI ("ast-mcp"),
	// or "native" to force mcp-local's writers. Empty = auto-detect from command.
	Installer string `yaml:"installer,omitempty"`
	// ExtraPorts are additional ports the service owns (e.g. dashboard, admin).
	ExtraPorts []int `yaml:"extra_ports,omitempty"`
}

// StateKind classifies a service's live state for status/doctor output.
type StateKind string

const (
	StateRunning  StateKind = "running"
	StateDegraded StateKind = "degraded"
	StateOrphaned StateKind = "orphaned"
	StateConflict StateKind = "port conflict"
	StateStopped  StateKind = "stopped"
	StateAdopting StateKind = "adopting"
)

// PortInfo describes a port the service owns and who is holding it.
type PortInfo struct {
	Port    int
	Running bool
	PID     int
	Command string
	IsSame  bool
}

// OwnedPorts returns every port this service owns: port, the dashboard port,
// the health port (when on localhost), and any extra_ports.
func (s ServiceConfig) OwnedPorts() []int {
	seen := make(map[int]bool)
	var ports []int
	add := func(p int) {
		if p <= 0 || seen[p] {
			return
		}
		seen[p] = true
		ports = append(ports, p)
	}
	add(s.Port)
	if p := portFromURL(s.DashboardURL); p > 0 {
		add(p)
	}
	if p := portFromURL(s.HealthURL); p > 0 {
		add(p)
	}
	if p := portFromURL(s.MCPURL); p > 0 {
		add(p)
	}
	for _, p := range s.ExtraPorts {
		add(p)
	}
	return ports
}

// portFromURL returns the TCP port from a localhost URL, or 0.
func portFromURL(u string) int {
	u = strings.TrimSpace(u)
	if u == "" {
		return 0
	}
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		return 0
	}
	u = strings.TrimPrefix(u, "http://")
	u = strings.TrimPrefix(u, "https://")
	if idx := strings.Index(u, "/"); idx != -1 {
		u = u[:idx]
	}
	if idx := strings.Index(u, ":"); idx != -1 {
		u = u[idx+1:]
	}
	p, err := strconv.Atoi(u)
	if err != nil {
		return 0
	}
	return p
}

const (
	InstallerASTMCP = "ast-mcp"
	InstallerNative = "native"
)

type ToolConfig struct {
	Name        string                 `yaml:"name"`
	Enabled     bool                   `yaml:"enabled"`
	Tier        string                 `yaml:"tier"`
	Description string                 `yaml:"description"`
	InputSchema map[string]interface{} `yaml:"input_schema,omitempty"`
}

type ManagerConfig struct {
	Services []ServiceConfig `yaml:"services"`
	Agents   AgentsConfig    `yaml:"agents,omitempty"`
}

func (s ServiceConfig) IsHTTP() bool {
	t := strings.ToLower(strings.TrimSpace(s.MCPType))
	return t == "http" || (s.Port > 0 && t != "stdio")
}

// EffectiveInstaller returns the installer that owns agent registration for s
// ("" means mcp-local's native writers).
func (s ServiceConfig) EffectiveInstaller() string {
	switch inst := strings.TrimSpace(s.Installer); inst {
	case InstallerNative:
		return ""
	case "":
		if filepath.Base(strings.TrimSpace(s.Command)) == InstallerASTMCP {
			return InstallerASTMCP
		}
		return ""
	default:
		return inst
	}
}

// MCPEndpoint is the HTTP MCP URL for s (mcp_url, else http://localhost:<port>/mcp), or "".
func (s ServiceConfig) MCPEndpoint() string {
	if u := strings.TrimSpace(s.MCPURL); u != "" {
		return u
	}
	if s.Port > 0 {
		return fmt.Sprintf("http://localhost:%d/mcp", s.Port)
	}
	return ""
}

func EffectiveBuild(s ServiceConfig) string {
	if s.BuildCmd != "" {
		return s.BuildCmd
	}
	return s.BuildCommand
}

func ShouldBuildOnStart(s ServiceConfig) bool {
	return !s.NoBuildOnStart
}

func ExpandService(s *ServiceConfig) {
	s.Command = ExpandPath(s.Command)
	s.HealthURL = ExpandPath(s.HealthURL)
	s.DashboardURL = ExpandPath(s.DashboardURL)
	s.MCPURL = ExpandPath(s.MCPURL)
	s.Log = ExpandPath(s.Log)
	s.Path = ExpandPath(s.Path)
	s.ToolsConfigPath = ExpandPath(s.ToolsConfigPath)
	if s.Installer == "" {
		s.Installer = s.EffectiveInstaller()
	}
	for i := range s.Deps {
		s.Deps[i] = ExpandPath(s.Deps[i])
	}
	if s.Env != nil {
		for k, v := range s.Env {
			s.Env[k] = ExpandPath(v)
		}
	}
}

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".mcp-local", "config.yaml"), nil
}

func LoadConfig() (*ManagerConfig, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg ManagerConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	for i := range cfg.Services {
		ExpandService(&cfg.Services[i])
	}
	return &cfg, nil
}

func SaveConfig(cfg *ManagerConfig) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func ExpandPath(path string) string {
	home, _ := os.UserHomeDir()
	path = strings.ReplaceAll(path, "${HOME}", home)
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(home, path[2:])
	}
	return path
}

func (m *ManagerConfig) ServiceNamed(name string) (ServiceConfig, error) {
	for _, s := range m.Services {
		if s.Name == name {
			return s, nil
		}
	}
	return ServiceConfig{}, fmt.Errorf("unknown service %q", name)
}

func (m *ManagerConfig) SortedNames() []string {
	names := make([]string, 0, len(m.Services))
	for _, s := range m.Services {
		names = append(names, s.Name)
	}
	sort.Strings(names)
	return names
}

func (m *ManagerConfig) UpsertService(s ServiceConfig) {
	for i := range m.Services {
		if m.Services[i].Name == s.Name {
			m.Services[i] = s
			return
		}
	}
	m.Services = append(m.Services, s)
}
