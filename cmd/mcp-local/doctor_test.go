package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
	"github.com/coma-toast/mcp-local/internal/portutil"
)

// listenEnv makes the test binary act as a service: it listens on the comma-separated
// ports and blocks, so tests get a real process running a known executable.
const listenEnv = "MCPL_TEST_LISTEN"

func TestMain(m *testing.M) {
	if ports := os.Getenv(listenEnv); ports != "" {
		for _, p := range strings.Split(ports, ",") {
			ln, err := net.Listen("tcp", "127.0.0.1:"+p)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}
			defer ln.Close()
		}
		time.Sleep(time.Minute)
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port
}

// selfService is a service whose binary is this test binary, listening on its main port.
func selfService(t *testing.T, name string, main, dash int) config.ServiceConfig {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	return config.ServiceConfig{Name: name, Command: exe, Args: []string{"-test.run=^$"}, Port: main, MCPType: "http",
		MCPURL: fmt.Sprintf("http://localhost:%d/mcp", main), DashboardURL: fmt.Sprintf("http://localhost:%d", dash),
		Env: map[string]string{listenEnv: strconv.Itoa(main)}}
}

// startOrphan runs an untracked copy of the test binary holding port, as a process left
// behind by an earlier run would.
func startOrphan(t *testing.T, port int) *exec.Cmd {
	t.Helper()
	exe, _ := os.Executable()
	c := exec.Command(exe, "-test.run=^$")
	c.Env = append(os.Environ(), listenEnv+"="+strconv.Itoa(port))
	if err := c.Start(); err != nil {
		t.Fatal(err)
	}
	go func() { _ = c.Wait() }()
	t.Cleanup(func() { _ = c.Process.Kill() })
	if !portutil.WaitForPort(port, 5*time.Second) {
		t.Fatal("orphan never listened")
	}
	return c
}

func TestStart_ReplacesSameBinaryOrphan(t *testing.T) {
	lifecycleEnv(t)
	main, dash := freePort(t), freePort(t)
	orphan := startOrphan(t, dash)
	svc := selfService(t, "svc", main, dash)
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{svc}}
	started, failed := startServices(cfg, []string{"svc"}, false, io.Discard)
	if failed != 0 || len(started) != 1 {
		t.Fatalf("started=%d failed=%d", len(started), failed)
	}
	t.Cleanup(func() { _ = stopServiceNamed(cfg, "svc", false) })
	if process.Alive(orphan.Process.Pid) {
		t.Error("orphan still running")
	}
	if !portutil.IsRunning(main) {
		t.Error("service not listening on its main port")
	}
}

func TestStart_RefusesForeignProgramOnOwnedPort(t *testing.T) {
	lifecycleEnv(t)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	dash := ln.Addr().(*net.TCPAddr).Port
	svc := sleeper("svc")
	svc.Port, svc.DashboardURL = freePort(t), fmt.Sprintf("http://localhost:%d", dash)
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{svc}}
	started, failed := startServices(cfg, []string{"svc"}, false, io.Discard)
	if failed != 1 || len(started) != 0 {
		t.Fatalf("started=%d failed=%d, want refused", len(started), failed)
	}
	if !portutil.IsRunning(dash) {
		t.Error("the other program's port was taken away")
	}
}

func TestStop_StopsSameBinaryOrphan(t *testing.T) {
	lifecycleEnv(t)
	main, dash := freePort(t), freePort(t)
	orphan := startOrphan(t, dash)
	cfg := &config.ManagerConfig{Agents: &config.AgentsConfig{}, Services: []config.ServiceConfig{selfService(t, "svc", main, dash)}}
	if err := stopServiceNamed(cfg, "svc", false); err != nil {
		t.Fatal(err)
	}
	if process.Alive(orphan.Process.Pid) {
		t.Error("orphan still running")
	}
}

func TestTailLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x.log")
	var b strings.Builder
	for i := 1; i <= 30; i++ {
		fmt.Fprintf(&b, "line %d\n", i)
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := tailLines(path, 3)
	if err != nil || strings.Join(got, "|") != "line 28|line 29|line 30" {
		t.Fatalf("tailLines = %q, %v", got, err)
	}
	big := strings.Repeat("x", 70<<10) + "\nlast\n"
	if err := os.WriteFile(path, []byte(big), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, _ := tailLines(path, 10); len(got) != 1 || got[0] != "last" {
		t.Fatalf("tailLines on a big file = %d lines, want just the complete last line", len(got))
	}
}

// doctorEnv is lifecycleEnv plus a disk threshold of zero, so the host's real free space
// never adds a finding.
func doctorEnv(t *testing.T) {
	lifecycleEnv(t)
	old := lowDiskPercent
	lowDiskPercent = 0
	t.Cleanup(func() { lowDiskPercent = old })
}

func TestDoctor_ReadOnlyWithoutFix(t *testing.T) {
	doctorEnv(t)
	svc := sleeper("svc")
	cfg := &config.ManagerConfig{Agents: &config.AgentsConfig{}, Services: []config.ServiceConfig{svc}}
	pidFile := process.GetPIDFile("svc")
	if err := os.MkdirAll(filepath.Dir(pidFile), 0o700); err != nil {
		t.Fatal(err)
	}
	// Our own PID with a start time that can't match: a reused PID, so the file is stale.
	stale, _ := json.Marshal(map[string]any{"pid": os.Getpid(), "start": "Thu Jan  1 00:00:00 1970"})
	if err := os.WriteFile(pidFile, stale, 0o600); err != nil {
		t.Fatal(err)
	}
	if n := doctorService(cfg, svc, false, io.Discard); n != 1 {
		t.Fatalf("problems = %d, want 1 (stale PID file)", n)
	}
	if _, err := os.Stat(pidFile); err != nil {
		t.Fatal("doctor without --fix removed the PID file")
	}
	if n := doctorService(cfg, svc, true, io.Discard); n != 0 {
		t.Fatalf("problems after --fix = %d, want 0", n)
	}
	if _, err := os.Stat(pidFile); !os.IsNotExist(err) {
		t.Error("--fix left the stale PID file")
	}
}

func TestDoctor_ForeignProgramIsNeverFixed(t *testing.T) {
	doctorEnv(t)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	svc := sleeper("svc")
	svc.Port, svc.ExtraPorts = freePort(t), []int{ln.Addr().(*net.TCPAddr).Port}
	cfg := &config.ManagerConfig{Agents: &config.AgentsConfig{}, Services: []config.ServiceConfig{svc}}
	if n := doctorService(cfg, svc, true, io.Discard); n != 1 {
		t.Fatalf("problems = %d, want the conflict to remain", n)
	}
}

func TestDoctor_LowDiskReported(t *testing.T) {
	doctorEnv(t)
	lowDiskPercent = 101
	svc := sleeper("svc")
	cfg := &config.ManagerConfig{Agents: &config.AgentsConfig{}, Services: []config.ServiceConfig{svc}}
	if n := doctorService(cfg, svc, true, io.Discard); n != 1 {
		t.Fatalf("problems = %d, want the low-disk warning (no automatic fix)", n)
	}
}
