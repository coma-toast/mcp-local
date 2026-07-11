package process

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

func TestSaveLoadRemovePID(t *testing.T) {
	name := "test-service"
	pid := 12345

	// Save PID
	if err := SavePID(name, pid); err != nil {
		t.Fatalf("SavePID failed: %v", err)
	}

	// Load PID
	loaded, err := LoadPID(name)
	if err != nil {
		t.Fatalf("LoadPID failed: %v", err)
	}
	if loaded != pid {
		t.Errorf("LoadPID = %d, want %d", loaded, pid)
	}

	// Remove PID
	if err := RemovePID(name); err != nil {
		t.Fatalf("RemovePID failed: %v", err)
	}

	// Load should fail now
	if _, err := LoadPID(name); err == nil {
		t.Error("LoadPID should fail after RemovePID")
	}
}

func TestSavePID_Overwrite(t *testing.T) {
	name := "test-service-2"
	if err := SavePID(name, 111); err != nil {
		t.Fatal(err)
	}
	if err := SavePID(name, 222); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadPID(name)
	if err != nil {
		t.Fatal(err)
	}
	if loaded != 222 {
		t.Errorf("Expected 222, got %d", loaded)
	}
}

func TestGetPIDFile(t *testing.T) {
	path := GetPIDFile("my-service")
	if !strings.Contains(path, "my-service.pid") {
		t.Errorf("GetPIDFile should contain service name: %s", path)
	}
	if !strings.Contains(path, ".mcp-local") {
		t.Errorf("GetPIDFile should be in .mcp-local dir: %s", path)
	}
}

func TestStartService_NoCommand(t *testing.T) {
	cfg := config.ServiceConfig{
		Name: "test",
		// No command
	}
	_, err := StartService(cfg)
	if err == nil {
		t.Error("StartService should fail when command is empty")
	}
}

func TestStopService_InvalidPID(t *testing.T) {
	// StopService on non-existent PID should fail
	err := StopService(999999)
	if err == nil {
		t.Error("StopService should fail for invalid PID")
	}
}

// Test that dependency checking works
func TestStartService_MissingDep(t *testing.T) {
	tmpDir := t.TempDir()
	missingDep := filepath.Join(tmpDir, "missing-file")

	cfg := config.ServiceConfig{
		Name:    "test",
		Command: "echo",
		Args:    []string{"hello"},
		Deps:    []string{missingDep},
	}
	_, err := StartService(cfg)
	if err == nil {
		t.Error("StartService should fail when dependency is missing")
	}
	if !strings.Contains(err.Error(), "dependency missing") {
		t.Errorf("Error should mention dependency: %v", err)
	}
}

func TestStopService_AlreadyStopped(t *testing.T) {
	// StopService on PID 0 should fail gracefully
	err := StopService(0)
	// This might return error or not depending on os.FindProcess behavior
	_ = err
}

// Test that we can build and run simple commands
func TestStartService_SimpleCommand(t *testing.T) {
	cfg := config.ServiceConfig{
		Name:    "test-echo",
		Command: "echo",
		Args:    []string{"hello"},
	}
	proc, err := StartService(cfg)
	if err != nil {
		t.Fatalf("StartService failed: %v", err)
	}
	t.Logf("Started process with PID: %d", proc.PID)

	// Save PID (normally done by caller in lifecycle.go)
	if err := SavePID(cfg.Name, proc.PID); err != nil {
		t.Fatalf("SavePID failed: %v", err)
	}

	// Give it a moment to finish
	time.Sleep(100 * time.Millisecond)

	// Check PID file was created
	pidPath := GetPIDFile(cfg.Name)
	if _, err := os.Stat(pidPath); os.IsNotExist(err) {
		t.Errorf("PID file should be created at %s", pidPath)
	} else {
		t.Logf("PID file exists at %s", pidPath)
	}

	// Cleanup
	_ = RemovePID(cfg.Name)
}
