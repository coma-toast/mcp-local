package process

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

// TestMain points HOME at a temp dir so pid and log files never land in the real ~/.mcp-local.
func TestMain(m *testing.M) {
	home, err := os.MkdirTemp("", "process-home-*")
	if err != nil {
		panic(err)
	}
	os.Setenv("HOME", home)
	code := m.Run()
	os.RemoveAll(home)
	os.Exit(code)
}

// spawn starts a service through StartService, reaps it in the background, and kills its
// whole group at cleanup.
func spawn(t *testing.T, cfg config.ServiceConfig) int {
	t.Helper()
	p, err := StartService(cfg)
	if err != nil {
		t.Fatalf("StartService: %v", err)
	}
	reap(t, p.PID)
	return p.PID
}

func reap(t *testing.T, pid int) {
	if proc, err := os.FindProcess(pid); err == nil {
		go func() { _, _ = proc.Wait() }()
	}
	t.Cleanup(func() { _ = syscall.Kill(-pid, syscall.SIGKILL); _ = syscall.Kill(pid, syscall.SIGKILL) })
}

func waitGone(pid int) bool {
	for i := 0; i < 50; i++ {
		if !Alive(pid) {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false
}

func readPID(t *testing.T, path string) int {
	t.Helper()
	for i := 0; i < 50; i++ {
		if data, err := os.ReadFile(path); err == nil {
			if pid, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil && pid > 0 {
				return pid
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("no pid in %s", path)
	return 0
}

func shortStop(t *testing.T) {
	term, kill := TermGrace, KillGrace
	TermGrace, KillGrace = 300*time.Millisecond, 2*time.Second
	t.Cleanup(func() { TermGrace, KillGrace = term, kill })
}

func TestSaveLoadRemovePID(t *testing.T) {
	name := "test-service"
	if err := SavePID(name, 12345); err != nil {
		t.Fatalf("SavePID failed: %v", err)
	}
	if loaded, err := LoadPID(name); err != nil || loaded != 12345 {
		t.Fatalf("LoadPID = %d, %v", loaded, err)
	}
	if err := RemovePID(name); err != nil {
		t.Fatalf("RemovePID failed: %v", err)
	}
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
	if loaded, err := LoadPID(name); err != nil || loaded != 222 {
		t.Errorf("LoadPID = %d, %v; want 222", loaded, err)
	}
}

func TestSavePID_RecordsIdentity(t *testing.T) {
	if err := SavePID("self", os.Getpid()); err != nil {
		t.Fatal(err)
	}
	id, err := LoadIdentity("self")
	if err != nil {
		t.Fatal(err)
	}
	if id.PID != os.Getpid() || id.Start == "" || id.Exe == "" {
		t.Errorf("identity = %+v, want pid, start and exe", id)
	}
	if fi, _ := os.Stat(GetPIDFile("self")); fi.Mode().Perm() != 0o600 {
		t.Errorf("pid file mode = %v, want 0600", fi.Mode().Perm())
	}
}

func TestLoadPID_LegacyTrimsWhitespace(t *testing.T) {
	path := GetPIDFile("legacy")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(" 4321\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	id, err := LoadIdentity("legacy")
	if err != nil || id.PID != 4321 || id.Start != "" {
		t.Errorf("LoadIdentity = %+v, %v; want unverified pid 4321", id, err)
	}
}

func TestGetPIDFile(t *testing.T) {
	path := GetPIDFile("my-service")
	if !strings.HasSuffix(path, filepath.Join(".mcp-local", "pids", "my-service.pid")) {
		t.Errorf("GetPIDFile = %s", path)
	}
}

func TestCheck(t *testing.T) {
	pid := spawn(t, config.ServiceConfig{Name: "check", Command: "/bin/sleep", Args: []string{"30"}})
	if err := SavePID("check", pid); err != nil {
		t.Fatal(err)
	}
	id, _ := LoadIdentity("check")
	if st := id.Check(""); st != Match {
		t.Errorf("live recorded process: state %v, want Match", st)
	}
	legacy := Identity{PID: pid}
	if st := legacy.Check("/bin/sleep"); st != Match {
		t.Errorf("legacy file, same binary: state %v, want Match", st)
	}
	if st := legacy.Check("/usr/bin/other-binary"); st != Unverified {
		t.Errorf("legacy file, other binary: state %v, want Unverified", st)
	}
	orig := ProcInfo
	ProcInfo = func(int) (string, string, error) { return "Mon Jan  1 00:00:00 2001", "/bin/sleep", nil }
	st := id.Check("")
	ProcInfo = orig
	if st != Mismatch {
		t.Errorf("reused PID: state %v, want Mismatch", st)
	}
	_ = syscall.Kill(pid, syscall.SIGKILL)
	if !waitGone(pid) {
		t.Fatal("child did not exit")
	}
	if st := id.Check(""); st != Gone {
		t.Errorf("dead process: state %v, want Gone", st)
	}
}

func TestAlive_EPERMCountsAsAlive(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root can signal pid 1")
	}
	if !Alive(1) {
		t.Error("pid 1 exists but is not ours; EPERM must count as alive")
	}
	if Alive(0) {
		t.Error("pid 0 must not count as alive")
	}
}

func TestTerminate_KillsWrapperChildren(t *testing.T) {
	childFile := filepath.Join(t.TempDir(), "child.pid")
	pid := spawn(t, config.ServiceConfig{Name: "wrapper", Command: "/bin/sh", Args: []string{"-c", "sleep 30 & echo $! > " + childFile + "; wait"}})
	child := readPID(t, childFile)
	target := SignalTarget(pid)
	if target != -pid {
		t.Fatalf("SignalTarget = %d, want group %d", target, -pid)
	}
	killed, err := Terminate(target, nil)
	if err != nil || killed {
		t.Fatalf("Terminate = %v, %v", killed, err)
	}
	if !waitGone(child) {
		t.Errorf("wrapper child %d survived the stop", child)
	}
}

func TestTerminate_EscalatesToSIGKILL(t *testing.T) {
	shortStop(t)
	ready := filepath.Join(t.TempDir(), "ready")
	pid := spawn(t, config.ServiceConfig{Name: "stubborn", Command: "/bin/sh", Args: []string{"-c", "trap '' TERM; touch " + ready + "; sleep 30"}})
	// SIGTERM must not arrive before the trap is installed, or the shell just exits.
	for deadline := time.Now().Add(10 * time.Second); ; time.Sleep(10 * time.Millisecond) {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("stubborn process never installed its TERM trap")
		}
	}
	killed, err := Terminate(SignalTarget(pid), nil)
	if err != nil || !killed {
		t.Fatalf("Terminate = %v, %v; want killed", killed, err)
	}
	if Alive(-pid) {
		t.Error("process group still alive")
	}
}

func TestTerminate_FailsWhileDoneFails(t *testing.T) {
	shortStop(t)
	KillGrace = 200 * time.Millisecond
	pid := spawn(t, config.ServiceConfig{Name: "held", Command: "/bin/sleep", Args: []string{"30"}})
	_, err := Terminate(SignalTarget(pid), func() error { return errors.New("port 1 still in use") })
	if err == nil || !strings.Contains(err.Error(), "still in use") {
		t.Errorf("Terminate error = %v, want port still in use", err)
	}
}

func TestTerminate_RefusesDangerousTargets(t *testing.T) {
	for _, target := range []int{0, 1, -1, os.Getpid(), -syscall.Getpgrp()} {
		if _, err := Terminate(target, nil); err == nil {
			t.Errorf("Terminate(%d) should refuse", target)
		}
	}
}

func TestStartService_NoCommand(t *testing.T) {
	if _, err := StartService(config.ServiceConfig{Name: "test"}); err == nil {
		t.Error("StartService should fail when command is empty")
	}
}

func TestStartService_MissingDep(t *testing.T) {
	cfg := config.ServiceConfig{Name: "test", Command: "echo", Deps: []string{filepath.Join(t.TempDir(), "missing")}}
	_, err := StartService(cfg)
	if err == nil || !strings.Contains(err.Error(), "dependency missing") {
		t.Errorf("StartService error = %v, want dependency missing", err)
	}
}

func TestStartService_DefaultLogIsPrivate(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cfg := config.ServiceConfig{Name: "logged", Command: "/bin/echo", Args: []string{"hello"}}
	spawn(t, cfg)
	path := LogPath(cfg)
	home, _ := os.UserHomeDir()
	if path != filepath.Join(home, ".mcp-local", "logs", "logged.log") {
		t.Errorf("LogPath = %s", path)
	}
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Errorf("log mode = %v, want 0600", fi.Mode().Perm())
	}
}

func TestStartService_DefaultLogRefusesSymlink(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cfg := config.ServiceConfig{Name: "planted", Command: "/bin/echo"}
	target := filepath.Join(t.TempDir(), "victim")
	if err := os.MkdirAll(filepath.Dir(LogPath(cfg)), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, LogPath(cfg)); err != nil {
		t.Fatal(err)
	}
	if _, err := StartService(cfg); err == nil {
		t.Error("StartService followed a symlink planted at the default log path")
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Error("symlink target was created")
	}
}

func TestStartService_ConfiguredLog(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "custom.log")
	spawn(t, config.ServiceConfig{Name: "custom", Command: "/bin/echo", Log: path})
	if _, err := os.Stat(path); err != nil {
		t.Errorf("configured log not created: %v", err)
	}
}

func TestBuild_RunsInCommandDir(t *testing.T) {
	dir, _ := filepath.EvalSymlinks(t.TempDir())
	out := filepath.Join(dir, "pwd.txt")
	if err := Build(config.ServiceConfig{Name: "b", Command: filepath.Join(dir, "server"), BuildCmd: "pwd > " + out}); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(out)
	if strings.TrimSpace(string(got)) != dir {
		t.Errorf("build ran in %q, want %q", strings.TrimSpace(string(got)), dir)
	}
}

func TestBuild_Timeout(t *testing.T) {
	orig := BuildTimeout
	BuildTimeout = 200 * time.Millisecond
	defer func() { BuildTimeout = orig }()
	start := time.Now()
	err := Build(config.ServiceConfig{Name: "slow", BuildCmd: "sleep 5"})
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Errorf("Build error = %v, want timeout", err)
	}
	if time.Since(start) > 3*time.Second {
		t.Errorf("Build took %v", time.Since(start))
	}
}

func TestBuild_Failure(t *testing.T) {
	if err := Build(config.ServiceConfig{Name: "bad", BuildCmd: "exit 3"}); err == nil {
		t.Error("Build should fail")
	}
	if err := Build(config.ServiceConfig{Name: "none"}); err != nil {
		t.Errorf("no build command: %v", err)
	}
}

func TestModeOf(t *testing.T) {
	cases := []struct {
		svc  config.ServiceConfig
		want Mode
	}{
		{config.ServiceConfig{MCPType: "filesystem", Port: 4000}, Embedded},
		{config.ServiceConfig{MCPURL: "https://x/mcp"}, Remote},
		{config.ServiceConfig{MCPType: "http", MCPURL: "https://x/mcp"}, Remote},
		{config.ServiceConfig{Command: "npx", MCPType: "stdio"}, Stdio},
		{config.ServiceConfig{Command: "npx"}, Stdio},
		{config.ServiceConfig{Command: "ast-mcp", Port: 7821}, Daemon},
		{config.ServiceConfig{Command: "srv", MCPURL: "http://localhost:1/mcp"}, Daemon},
	}
	for i, c := range cases {
		if got := ModeOf(c.svc); got != c.want {
			t.Errorf("case %d: ModeOf = %v, want %v", i, got, c.want)
		}
	}
}

func TestProcInfo_Self(t *testing.T) {
	start, exe, err := ProcInfo(os.Getpid())
	if err != nil || start == "" {
		t.Fatalf("ProcInfo(self) = %q, %q, %v", start, exe, err)
	}
	self, _ := os.Executable()
	if !SameBinary(exe, self) {
		t.Errorf("exe %q does not match %q", exe, self)
	}
}
