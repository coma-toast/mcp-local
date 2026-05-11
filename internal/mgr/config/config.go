package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type ServiceConfig struct {
	Name         string            `yaml:"name"`
	Command      string            `yaml:"command"`
	Args         []string          `yaml:"args"`
	Env          map[string]string `yaml:"env"`
	Port         int               `yaml:"port"`
	MCPType      string            `yaml:"type"` // "http" or "stdio"
	BuildCmd     string            `yaml:"build_cmd"`
	BuildCommand string            `yaml:"build_command,omitempty"`
	Deps         []string          `yaml:"deps"`
	Description  string            `yaml:"description"`
	Tools        []ToolConfig      `yaml:"tools"`
	HealthURL    string            `yaml:"health_url,omitempty"`
	DashboardURL string            `yaml:"dashboard_url,omitempty"`
	MCPURL       string            `yaml:"mcp_url,omitempty"`
	Log          string            `yaml:"log,omitempty"`
}

type ToolConfig struct {
	Name        string `yaml:"name"`
	Enabled     bool   `yaml:"enabled"`
	Tier        string `yaml:"tier"`
	Description string `yaml:"description"`
}

type ManagerConfig struct {
	Services []ServiceConfig `yaml:"services"`
}

func EffectiveBuild(s ServiceConfig) string {
	if s.BuildCmd != "" {
		return s.BuildCmd
	}
	return s.BuildCommand
}

func ExpandService(s *ServiceConfig) {
	s.Command = ExpandPath(s.Command)
	s.HealthURL = ExpandPath(s.HealthURL)
	s.DashboardURL = ExpandPath(s.DashboardURL)
	s.MCPURL = ExpandPath(s.MCPURL)
	s.Log = ExpandPath(s.Log)
	for i := range s.Deps {
		s.Deps[i] = ExpandPath(s.Deps[i])
	}
	if s.Env != nil {
		for k, v := range s.Env {
			s.Env[k] = ExpandPath(v)
		}
	}
}

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mcp-local", "config.yaml")
}

func LoadConfig() (*ManagerConfig, error) {
	path := GetConfigPath()
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
	path := GetConfigPath()
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
