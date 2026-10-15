package runtime

import (
	"os"
	"os/exec"
	"testing"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
)

// TestMain points HOME at a temp dir so pid files are never read from the real ~/.mcp-local.
func TestMain(m *testing.M) {
	home, err := os.MkdirTemp("", "runtime-home-*")
	if err != nil {
		panic(err)
	}
	os.Setenv("HOME", home)
	code := m.Run()
	os.RemoveAll(home)
	os.Exit(code)
}

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
	if os.Geteuid() != 0 && !IsPIDAlive(1) {
		t.Error("EPERM (pid 1, not ours) should count as alive")
	}
}

// A read-only status must not delete PID files, even stale ones.
func TestCheckStatus_StalePIDFileKept(t *testing.T) {
	cmd := exec.Command("/bin/sleep", "30")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	pid := cmd.Process.Pid
	if err := process.SavePID("stale-svc", pid); err != nil {
		t.Fatal(err)
	}
	_ = cmd.Process.Kill()
	_ = cmd.Wait()
	svc := config.ServiceConfig{Name: "stale-svc", Command: "/bin/sleep", MCPType: "http", MCPURL: "http://localhost:1/mcp"}
	st := CheckStatus(svc)
	if st.Running || !st.Stale {
		t.Errorf("status = %+v, want stopped and stale", st)
	}
	if _, err := os.Stat(process.GetPIDFile("stale-svc")); err != nil {
		t.Errorf("status removed the pid file: %v", err)
	}
}

// A PID reused by another process must not read as running.
func TestCheckStatus_ReusedPIDNotRunning(t *testing.T) {
	if err := process.SavePID("reused-svc", os.Getpid()); err != nil {
		t.Fatal(err)
	}
	orig := process.ProcInfo
	process.ProcInfo = func(int) (string, string, error) { return "Mon Jan  1 00:00:00 2001", "/bin/other", nil }
	defer func() { process.ProcInfo = orig }()
	st := CheckStatus(config.ServiceConfig{Name: "reused-svc", Command: "/bin/sleep", MCPType: "stdio"})
	if st.Running {
		t.Errorf("reused PID reported as running: %+v", st)
	}
}

func TestCheckStatus_StdioWithPIDFile(t *testing.T) {
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
