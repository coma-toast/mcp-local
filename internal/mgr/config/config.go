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
}

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
