package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
	"github.com/spf13/cobra"
)

// lifecycleEnv isolates HOME (pid files, logs, agent configs) and shortens stop timeouts.
func lifecycleEnv(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	term, kill := process.TermGrace, process.KillGrace
	process.TermGrace, process.KillGrace = 500*time.Millisecond, time.Second
	t.Cleanup(func() { process.TermGrace, process.KillGrace = term, kill })
}

// sleeper is a daemon-mode service with no port, so nothing depends on real ports.
func sleeper(name string) config.ServiceConfig {
	return config.ServiceConfig{Name: name, Command: "/bin/sleep", Args: []string{"30"}, MCPType: "http", MCPURL: "http://localhost:1/mcp"}
}

// track reaps pid in the background (it is our child) and kills its group at cleanup.
func track(t *testing.T, pid int) {
	if proc, err := os.FindProcess(pid); err == nil {
		go func() { _, _ = proc.Wait() }()
	}
	t.Cleanup(func() { _ = syscall.Kill(-pid, syscall.SIGKILL); _ = syscall.Kill(pid, syscall.SIGKILL) })
}

func gone(pid int) bool {
	for i := 0; i < 50; i++ {
		if !process.Alive(pid) {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false
}

func startTracked(t *testing.T, cfg *config.ManagerConfig, name string) int {
	t.Helper()
	if _, failed := startServices(cfg, []string{name}, true, io.Discard); failed != 0 {
		t.Fatalf("start %s failed", name)
	}
	pid, err := process.LoadPID(name)
	if err != nil {
		t.Fatal(err)
	}
	track(t, pid)
	return pid
}

func TestStop_ReusedPIDNotSignalled(t *testing.T) {
	lifecycleEnv(t)
	victim := exec.Command("/bin/sleep", "30")
	if err := victim.Start(); err != nil {
		t.Fatal(err)
	}
	track(t, victim.Process.Pid)
	if err := process.SavePID("svc", victim.Process.Pid); err != nil {
		t.Fatal(err)
	}
	orig := process.ProcInfo
	process.ProcInfo = func(int) (string, string, error) { return "Mon Jan  1 00:00:00 2001", "/bin/sleep", nil }
	defer func() { process.ProcInfo = orig }()
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{sleeper("svc")}}
	if err := stopServiceNamed(cfg, "svc", false); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if !process.Alive(victim.Process.Pid) {
		t.Error("stop signalled a process that only reused the recorded PID")
	}
	if _, err := os.Stat(process.GetPIDFile("svc")); !os.IsNotExist(err) {
		t.Error("stale pid file should be removed by stop")
	}
}

func TestStop_KillsWrapperChildren(t *testing.T) {
	lifecycleEnv(t)
	childFile := filepath.Join(t.TempDir(), "child.pid")
	svc := sleeper("wrapped")
	svc.Command, svc.Args = "/bin/sh", []string{"-c", "sleep 30 & echo $! > " + childFile + "; wait"}
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{svc}}
	startTracked(t, cfg, "wrapped")
	var child int
	for i := 0; i < 50 && child == 0; i++ {
		data, _ := os.ReadFile(childFile)
		child, _ = strconv.Atoi(strings.TrimSpace(string(data)))
		time.Sleep(50 * time.Millisecond)
	}
	if child == 0 {
		t.Fatal("wrapper child never started")
	}
	if err := stopServiceNamed(cfg, "wrapped", false); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if !gone(child) {
		t.Errorf("wrapper child %d survived stop", child)
	}
}

func TestStop_NoPortEscalatesToSIGKILL(t *testing.T) {
	lifecycleEnv(t)
	svc := sleeper("stubborn")
	svc.Command, svc.Args = "/bin/sh", []string{"-c", "trap '' TERM; sleep 30"}
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{svc}}
	pid := startTracked(t, cfg, "stubborn")
	time.Sleep(100 * time.Millisecond)
	if err := stopServiceNamed(cfg, "stubborn", false); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if !gone(-pid) {
		t.Error("process group still running after stop reported success")
	}
}

// Success needs the port to be free, not just the process gone.
func TestStop_FailsWhilePortHeld(t *testing.T) {
	lifecycleEnv(t)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	svc := sleeper("held")
	svc.Port = ln.Addr().(*net.TCPAddr).Port
	victim := exec.Command("/bin/sleep", "30")
	victim.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := victim.Start(); err != nil {
		t.Fatal(err)
	}
	track(t, victim.Process.Pid)
	if err := process.SavePID("held", victim.Process.Pid); err != nil {
		t.Fatal(err)
	}
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{svc}}
	err = stopServiceNamed(cfg, "held", false)
	if err == nil || !strings.Contains(err.Error(), "still in use") {
		t.Errorf("stop error = %v, want port still in use", err)
	}
}

// Without a PID file, stop only touches a port owner running the service's binary.
func TestStop_ForeignPortOwnerRefused(t *testing.T) {
	lifecycleEnv(t)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	svc := sleeper("foreign")
	svc.Port = ln.Addr().(*net.TCPAddr).Port
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{svc}}
	if err := stopServiceNamed(cfg, "foreign", false); err == nil {
		t.Error("stop should refuse a port owner it didn't start")
	}
}

func TestStart_RemoteAndStdioRegisterOnly(t *testing.T) {
	lifecycleEnv(t)
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{
		{Name: "remote", MCPURL: "https://example.com/mcp"},
		{Name: "stdio", Command: "/bin/sleep", Args: []string{"30"}, MCPType: "stdio"},
	}}
	started, failed := startServices(cfg, cfg.SortedNames(), true, io.Discard)
	if failed != 0 || len(started) != 2 {
		t.Fatalf("started=%d failed=%d, want both registered", len(started), failed)
	}
	for _, name := range []string{"remote", "stdio"} {
		if _, err := os.Stat(process.GetPIDFile(name)); !os.IsNotExist(err) {
			t.Errorf("%s was spawned", name)
		}
	}
}

func TestStart_RemoteWithIdlePortStillRegisters(t *testing.T) {
	lifecycleEnv(t)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{{Name: "remote", MCPURL: fmt.Sprintf("http://localhost:%d/mcp", port), Port: port}}}
	started, failed := startServices(cfg, cfg.SortedNames(), true, io.Discard)
	if failed != 0 || len(started) != 1 {
		t.Fatalf("started=%d failed=%d, want registered despite the idle port", len(started), failed)
	}
}

func TestStart_NoCommandNoURLFails(t *testing.T) {
	lifecycleEnv(t)
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{{Name: "empty", Port: 1}, {Name: "empty-stdio", MCPType: "stdio"}}}
	started, failed := startServices(cfg, cfg.SortedNames(), true, io.Discard)
	if failed != 2 || len(started) != 0 {
		t.Fatalf("started=%d failed=%d, want both to fail", len(started), failed)
	}
}

func TestStart_FailureCounted(t *testing.T) {
	lifecycleEnv(t)
	svc := sleeper("broken")
	svc.Command = filepath.Join(t.TempDir(), "missing-binary")
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{svc, sleeper("ok")}}
	started, failed := startServices(cfg, []string{"broken"}, true, io.Discard)
	if failed != 1 || len(started) != 0 {
		t.Errorf("started=%d failed=%d, want 0/1", len(started), failed)
	}
}

func TestStart_NoDuplicateWithoutPort(t *testing.T) {
	lifecycleEnv(t)
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{sleeper("once")}}
	first := startTracked(t, cfg, "once")
	if _, failed := startServices(cfg, []string{"once"}, true, io.Discard); failed != 0 {
		t.Fatal("second start failed")
	}
	if second, _ := process.LoadPID("once"); second != first {
		t.Errorf("second start spawned pid %d (first %d)", second, first)
	}
}

func TestRestart_BuildFailureLeavesServiceRunning(t *testing.T) {
	lifecycleEnv(t)
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{sleeper("svc")}}
	pid := startTracked(t, cfg, "svc")
	cfg.Services[0].BuildCmd = "exit 3"
	_, failed := restartServices(cfg, []string{"svc"}, io.Discard)
	if failed != 1 {
		t.Errorf("failed = %d, want 1", failed)
	}
	if !process.Alive(pid) {
		t.Error("a failed build took the running service down")
	}
	if got, _ := process.LoadPID("svc"); got != pid {
		t.Errorf("pid file = %d, want %d", got, pid)
	}
}

func TestRestart_BuildsOnceThenRestarts(t *testing.T) {
	lifecycleEnv(t)
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{sleeper("svc")}}
	old := startTracked(t, cfg, "svc")
	count := filepath.Join(t.TempDir(), "builds")
	cfg.Services[0].BuildCmd = "echo built >> " + count
	if _, failed := restartServices(cfg, []string{"svc"}, io.Discard); failed != 0 {
		t.Fatalf("restart failed")
	}
	pid, _ := process.LoadPID("svc")
	track(t, pid)
	if pid == old || !gone(old) {
		t.Errorf("old pid %d not replaced (new %d)", old, pid)
	}
	data, _ := os.ReadFile(count)
	if n := strings.Count(string(data), "built"); n != 1 {
		t.Errorf("build ran %d times, want 1", n)
	}
}

func TestRunHealth_CountsFailures(t *testing.T) {
	ok := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(`{"ok":true}`)) }))
	defer ok.Close()
	bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusInternalServerError) }))
	defer bad.Close()
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{
		{Name: "bad", HealthURL: bad.URL},
		{Name: "none"},
		{Name: "ok", HealthURL: ok.URL},
	}}
	if n := runHealth(cfg, cfg.SortedNames(), ok.Client(), io.Discard, io.Discard); n != 1 {
		t.Errorf("runHealth = %d failures, want 1", n)
	}
}

func TestFailures(t *testing.T) {
	cmd := &cobra.Command{}
	if failures(cmd, 0, "stop") != nil {
		t.Error("no failures should be nil")
	}
	if err := failures(cmd, 2, "stop"); err == nil || !cmd.SilenceUsage {
		t.Errorf("failures(2) = %v, SilenceUsage=%v", err, cmd.SilenceUsage)
	}
}
