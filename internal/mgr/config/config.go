package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type ServiceConfig struct {
	Name        string            `yaml:"name"`
	Command     string            `yaml:"command"`
	Args        []string          `yaml:"args"`
	Env         map[string]string `yaml:"env"`
	Port        int               `yaml:"port"`
	MCPType     string            `yaml:"type"` // "http" or "stdio"
	BuildCmd    string            `yaml:"build_cmd"`
	Deps        []string          `yaml:"deps"` // Files/dirs that must exist
	Description string            `yaml:"description"`
}

type ManagerConfig struct {
	Services []ServiceConfig `yaml:"services"`
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
	return strings.ReplaceAll(path, "${HOME}", home)
}
