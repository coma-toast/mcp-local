package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type AgentsConfig struct {
	OpenCode bool `yaml:"opencode"`
	Cursor   bool `yaml:"cursor"`
	Claude   bool `yaml:"claude"`
}

type ServiceConfig struct {
	Name           string            `yaml:"name"`
	Command        string            `yaml:"command,omitempty"`
	Args           []string          `yaml:"args,omitempty"`
	Env            map[string]string `yaml:"env,omitempty"`
	Port           int               `yaml:"port,omitempty"`
	MCPType        string            `yaml:"type,omitempty"`
	BuildCmd       string            `yaml:"build_cmd,omitempty"`
	BuildCommand   string            `yaml:"build_command,omitempty"`
	NoBuildOnStart bool              `yaml:"no_build_on_start,omitempty"`
	Deps           []string          `yaml:"deps,omitempty"`
	Description    string            `yaml:"description,omitempty"`
	Tools          []ToolConfig      `yaml:"tools,omitempty"`
	HealthURL      string            `yaml:"health_url,omitempty"`
	DashboardURL   string            `yaml:"dashboard_url,omitempty"`
	MCPURL         string            `yaml:"mcp_url,omitempty"`
	// Headers are sent to an HTTP service (written into each host's entry where supported).
	Headers         map[string]string `yaml:"headers,omitempty"`
	Log             string            `yaml:"log,omitempty"`
	ActiveTier      string            `yaml:"active_tier,omitempty"`
	CodeMode        bool              `yaml:"code_mode,omitempty"`
	NoCodeMode      bool              `yaml:"no_code_mode,omitempty"`
	ToolsConfigPath string            `yaml:"tools_config_path,omitempty"`
	Path            string            `yaml:"path,omitempty"`
	// Installer delegates agent registration to the service's own CLI ("ast-mcp"),
	// or "native" to force mcp-local's writers. Empty = auto-detect from command.
	Installer string `yaml:"installer,omitempty"`
}

const (
	InstallerASTMCP = "ast-mcp"
	InstallerNative = "native"
)

// ToolConfig is one tool's settings. In YAML a missing `enabled:` means enabled, and
// `enabled: true` is not written back.
type ToolConfig struct {
	Name        string
	Enabled     bool
	Tier        string
	Description string
	InputSchema map[string]interface{}
}

type toolYAML struct {
	Name        string                 `yaml:"name"`
	Enabled     *bool                  `yaml:"enabled,omitempty"`
	Tier        string                 `yaml:"tier,omitempty"`
	Description string                 `yaml:"description,omitempty"`
	InputSchema map[string]interface{} `yaml:"input_schema,omitempty"`
}

func (t *ToolConfig) UnmarshalYAML(n *yaml.Node) error {
	var y toolYAML
	if err := n.Decode(&y); err != nil {
		return err
	}
	*t = ToolConfig{Name: y.Name, Enabled: y.Enabled == nil || *y.Enabled, Tier: y.Tier, Description: y.Description, InputSchema: y.InputSchema}
	return nil
}

func (t ToolConfig) MarshalYAML() (interface{}, error) {
	y := toolYAML{Name: t.Name, Tier: t.Tier, Description: t.Description, InputSchema: t.InputSchema}
	if !t.Enabled {
		y.Enabled = &t.Enabled
	}
	return y, nil
}

type ManagerConfig struct {
	Services []ServiceConfig `yaml:"services"`
	// Agents is nil when config.yaml has no `agents:` key (legacy: register with every host).
	Agents *AgentsConfig `yaml:"agents,omitempty"`

	doc *yaml.Node               // the file as loaded, so saves keep comments and unknown keys
	raw map[string]ServiceConfig // services as written in the file (unexpanded), by name
}

// IsHTTP reports whether s is reached over HTTP: type http, a port (unless type stdio), or
// only an mcp_url (a remote service).
func (s ServiceConfig) IsHTTP() bool {
	t := strings.ToLower(strings.TrimSpace(s.MCPType))
	if t == "stdio" {
		return false
	}
	return t == "http" || s.Port > 0 || (strings.TrimSpace(s.Command) == "" && strings.TrimSpace(s.MCPURL) != "")
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

// LoadConfig reads config.yaml. Services come back with paths expanded; the file's own
// (unexpanded) form is kept so SaveConfig can write it back unchanged.
func LoadConfig() (*ManagerConfig, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseConfig(data)
}

func parseConfig(data []byte) (*ManagerConfig, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	var cfg, raw ManagerConfig
	if len(doc.Content) > 0 {
		if err := doc.Decode(&cfg); err != nil {
			return nil, err
		}
		_ = doc.Decode(&raw)
		cfg.doc = &doc
	}
	cfg.raw = byName(raw.Services)
	for i := range cfg.Services {
		ExpandService(&cfg.Services[i])
	}
	return &cfg, nil
}

func byName(services []ServiceConfig) map[string]ServiceConfig {
	m := make(map[string]ServiceConfig, len(services))
	for _, s := range services {
		m[s.Name] = s
	}
	return m
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
