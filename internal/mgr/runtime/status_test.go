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

func TestCheckStatus_States(t *testing.T) {
	// Test stopped state
	svc := config.ServiceConfig{
		Name: "stopped-svc",
		Port: 99997,
	}
	st := CheckStatus(svc)
	if st.State != config.StateStopped {
		t.Errorf("state = %s, want stopped", st.State)
	}
}

// fakeHost stubs the port and process lookups CheckStatus uses.
func fakeHost(t *testing.T, listening map[int]int, exes map[int]string, health func(string) (bool, string)) {
	t.Helper()
	l, o, p, h := portListening, portOwner, procInfo, fetchHealth
	t.Cleanup(func() { portListening, portOwner, procInfo, fetchHealth = l, o, p, h })
	portListening = func(port int) bool { _, ok := listening[port]; return ok }
	portOwner = func(port int) int { return listening[port] }
	procInfo = func(pid int) (string, string, error) { return "", exes[pid], nil }
	fetchHealth = func(u string) (bool, string) {
		if health == nil || u == "" {
			return true, "healthy"
		}
		return health(u)
	}
}

func TestCheckStatus_OwnedPortStates(t *testing.T) {
	svc := config.ServiceConfig{Name: "ast", Command: "/opt/ast/ast-mcp", Port: 7821, DashboardURL: "http://localhost:7830",
		HealthURL: "http://localhost:7821/health"}
	cases := []struct {
		name      string
		listening map[int]int
		exes      map[int]string
		health    func(string) (bool, string)
		state     config.StateKind
		orphan    int
	}{
		{"orphan holds the dashboard port", map[int]int{7830: 4242}, map[int]string{4242: "/opt/ast/ast-mcp"}, nil, config.StateOrphaned, 4242},
		{"another program holds it", map[int]int{7830: 5151}, map[int]string{5151: "/usr/bin/python3"}, nil, config.StateConflict, 0},
		{"main port up and healthy", map[int]int{7821: 4242}, map[int]string{4242: "/opt/ast/ast-mcp"}, nil, config.StateRunning, 0},
		{"main port up, health failing", map[int]int{7821: 4242}, map[int]string{4242: "/opt/ast/ast-mcp"},
			func(string) (bool, string) { return false, "failed: data directory unavailable" }, config.StateDegraded, 0},
		{"nothing listening", nil, nil, nil, config.StateStopped, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fakeHost(t, c.listening, c.exes, c.health)
			st := CheckStatus(svc)
			if st.State != c.state || st.Orphan.PID != c.orphan {
				t.Errorf("state=%s orphan=%d detail=%q, want %s/%d", st.State, st.Orphan.PID, st.Detail, c.state, c.orphan)
			}
		})
	}
}
