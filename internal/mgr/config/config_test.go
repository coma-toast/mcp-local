package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPath(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		input    string
		expected string
	}{
		{"${HOME}/foo", home + "/foo"},
		{"~/bar", home + "/bar"},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
	}

	for _, tt := range tests {
		got := ExpandPath(tt.input)
		if got != tt.expected {
			t.Errorf("ExpandPath(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestExpandService(t *testing.T) {
	home, _ := os.UserHomeDir()
	svc := ServiceConfig{
		Command:      "~/bin/mcp",
		HealthURL:    "${HOME}/health",
		DashboardURL: "~/dash",
		MCPURL:       "${HOME}/mcp",
		Log:          "~/logs/mcp.log",
		Deps:         []string{"${HOME}/dep1", "~/dep2"},
		Env:          map[string]string{"PATH": "${HOME}/bin"},
	}
	ExpandService(&svc)

	if svc.Command != home+"/bin/mcp" {
		t.Errorf("Command = %q, want %q", svc.Command, home+"/bin/mcp")
	}
	if svc.HealthURL != home+"/health" {
		t.Errorf("HealthURL = %q, want %q", svc.HealthURL, home+"/health")
	}
	if svc.Env["PATH"] != home+"/bin" {
		t.Errorf("Env[PATH] = %q, want %q", svc.Env["PATH"], home+"/bin")
	}
}

func TestServiceToEntry_HTTP(t *testing.T) {
	svc := ServiceConfig{
		Name:    "test-svc",
		Port:    8080,
		MCPType: "http",
	}
	entry, err := ServiceToEntry(svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Type != "remote" {
		t.Errorf("Type = %q, want remote", entry.Type)
	}
	if entry.URL != "http://localhost:8080/mcp" {
		t.Errorf("URL = %q, want http://localhost:8080/mcp", entry.URL)
	}
}

func TestServiceToEntry_HTTPWithMCPURL(t *testing.T) {
	svc := ServiceConfig{
		Name:   "test-svc",
		Port:   8080,
		MCPURL: "http://example.com/custom",
	}
	entry, err := ServiceToEntry(svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.URL != "http://example.com/custom" {
		t.Errorf("URL = %q, want http://example.com/custom", entry.URL)
	}
}

func TestServiceToEntry_Stdio(t *testing.T) {
	svc := ServiceConfig{
		Name:    "test-svc",
		Command: "/usr/local/bin/mcp",
		Args:    []string{"--config", "cfg.yaml"},
		MCPType: "stdio",
		Env:     map[string]string{"KEY": "VAL"},
	}
	entry, err := ServiceToEntry(svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Type != "local" {
		t.Errorf("Type = %q, want local", entry.Type)
	}
	if len(entry.Command) != 3 || entry.Command[0] != "/usr/local/bin/mcp" {
		t.Errorf("Command = %v, want [/usr/local/bin/mcp --config cfg.yaml]", entry.Command)
	}
	if entry.Environment["KEY"] != "VAL" {
		t.Errorf("Environment[KEY] = %q, want VAL", entry.Environment["KEY"])
	}
}

func TestShouldBuildOnStart(t *testing.T) {
	svcDefault := ServiceConfig{Name: "test"}
	if !ShouldBuildOnStart(svcDefault) {
		t.Error("ShouldBuildOnStart should be true by default")
	}

	svcSkip := ServiceConfig{Name: "test", NoBuildOnStart: true}
	if ShouldBuildOnStart(svcSkip) {
		t.Error("ShouldBuildOnStart should be false when NoBuildOnStart=true")
	}
}

func TestUpsertService(t *testing.T) {
	cfg := &ManagerConfig{
		Services: []ServiceConfig{{Name: "a"}, {Name: "b"}},
	}
	cfg.UpsertService(ServiceConfig{Name: "b", Port: 9999})
	if cfg.Services[1].Port != 9999 {
		t.Error("UpsertService should update existing")
	}

	cfg.UpsertService(ServiceConfig{Name: "c"})
	if len(cfg.Services) != 3 {
		t.Error("UpsertService should append new")
	}
}

func TestLoadSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfgDir := filepath.Join(tmpDir, ".mcp-local")
	os.MkdirAll(cfgDir, 0755)

	orig := &ManagerConfig{
		Services: []ServiceConfig{{Name: "test", Port: 8080}},
		Agents:   AgentsConfig{OpenCode: true, Cursor: true},
	}
	if err := SaveConfig(orig); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}

	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if len(loaded.Services) != 1 || loaded.Services[0].Name != "test" {
		t.Errorf("loaded config mismatch: %+v", loaded)
	}
	if !loaded.Agents.OpenCode || !loaded.Agents.Cursor {
		t.Error("agents config not preserved")
	}
}
