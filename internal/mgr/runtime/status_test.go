package runtime

import (
	"os"
	"testing"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
)

func TestIsPIDAlive(t *testing.T) {
	if IsPIDAlive(0) {
		t.Error("PID 0 should not be alive")
	}
	if IsPIDAlive(-1) {
		t.Error("PID -1 should not be alive")
	}
	if !IsPIDAlive(os.Getpid()) {
		t.Error("current process should be alive")
	}
}

func TestCheckStatus_StdioWithPIDFile(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	svc := config.ServiceConfig{
		Name:    "stdio-svc",
		Command: "/usr/bin/true",
		MCPType: "stdio",
	}

	_ = process.SavePID("stdio-svc", os.Getpid())
	status := CheckStatus(svc)

	if !status.Running {
		t.Error("stdio service with valid PID file should be running")
	}
	if status.PID != os.Getpid() {
		t.Errorf("PID = %d, want %d", status.PID, os.Getpid())
	}
}

func TestCheckStatus_HTTPRunning(t *testing.T) {
	svc := config.ServiceConfig{
		Name: "http-svc",
		Port: 99999,
	}
	status := CheckStatus(svc)
	if status.Running {
		t.Error("service on invalid port should not be running")
	}
}

func TestCheckStatus_MissingPIDFile(t *testing.T) {
	svc := config.ServiceConfig{
		Name:    "missing-svc",
		Command: "/usr/bin/true",
		MCPType: "stdio",
	}
	status := CheckStatus(svc)
	if status.Running {
		t.Error("service with no PID file should not be running")
	}
}

func TestCheckAll(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg := &config.ManagerConfig{
		Services: []config.ServiceConfig{
			{Name: "svc-a", Port: 99999},
			{Name: "svc-b", Command: "/usr/bin/true", MCPType: "stdio"},
		},
	}

	results := CheckAll(cfg, "")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestCheckAll_Filter(t *testing.T) {
	cfg := &config.ManagerConfig{
		Services: []config.ServiceConfig{
			{Name: "svc-a", Port: 99999},
			{Name: "svc-b", Port: 99998},
		},
	}

	results := CheckAll(cfg, "svc-a")
	if len(results) != 1 {
		t.Fatalf("expected 1 result with filter, got %d", len(results))
	}
	if results[0].Name != "svc-a" {
		t.Errorf("Name = %q, want svc-a", results[0].Name)
	}
}
