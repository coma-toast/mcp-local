package process

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

// Identity is what a PID file records about a process mcp-local started. Start and Exe
// are empty for legacy plain-PID files.
type Identity struct {
	PID   int    `json:"pid"`
	Start string `json:"start,omitempty"`
	Exe   string `json:"exe,omitempty"`
}

// State is how a recorded Identity relates to whatever now has its PID.
type State int

const (
	Gone       State = iota // no process has the PID
	Match                   // the process mcp-local recorded
	Mismatch                // the PID was reused by another process
	Unverified              // alive, but a legacy PID file can't prove it is ours
)

// ProcInfo returns pid's start time and executable, or an error when there is no such
// process. It is a variable so tests can fake PID reuse.
var ProcInfo = procInfo

func procInfo(pid int) (start, exe string, err error) {
	if pid <= 0 {
		return "", "", fmt.Errorf("invalid pid %d", pid)
	}
	if runtime.GOOS == "linux" {
		return procInfoLinux(pid)
	}
	return procInfoPS(pid)
}

func procInfoPS(pid int) (string, string, error) {
	cmd := exec.Command("ps", "-o", "lstart=,comm=", "-p", strconv.Itoa(pid))
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	out, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("pid %d: no such process", pid)
	}
	f := strings.Fields(string(out))
	if len(f) < 6 {
		return "", "", fmt.Errorf("pid %d: unexpected ps output %q", pid, out)
	}
	return strings.Join(f[:5], " "), strings.Join(f[5:], " "), nil
}

func procInfoLinux(pid int) (string, string, error) {
	dir := fmt.Sprintf("/proc/%d", pid)
	stat, err := os.ReadFile(dir + "/stat")
	if err != nil {
		return "", "", fmt.Errorf("pid %d: no such process", pid)
	}
	// comm (field 2) may contain spaces; starttime is field 22, the 20th after the ")".
	i := bytes.LastIndexByte(stat, ')')
	if i < 0 {
		return "", "", fmt.Errorf("pid %d: unexpected stat %q", pid, stat)
	}
	f := strings.Fields(string(stat[i+1:]))
	if len(f) < 20 {
		return "", "", fmt.Errorf("pid %d: unexpected stat %q", pid, stat)
	}
	boot, _ := os.ReadFile("/proc/sys/kernel/random/boot_id")
	start := strings.TrimSpace(string(boot)) + ":" + f[19]
	exe, err := os.Readlink(dir + "/exe")
	if err != nil {
		comm, _ := os.ReadFile(dir + "/comm")
		exe = strings.TrimSpace(string(comm))
	}
	// A rebuilt binary replaces the file the running process was started from.
	return start, strings.TrimSuffix(exe, " (deleted)"), nil
}

// Check reports whether id still names the process it was recorded for. The start time
// decides; a legacy plain-PID file is accepted only if the process runs wantExe.
func (id Identity) Check(wantExe string) State {
	if id.PID <= 1 {
		return Mismatch
	}
	start, exe, err := ProcInfo(id.PID)
	if err != nil {
		if Alive(id.PID) {
			return Unverified
		}
		return Gone
	}
	if id.Start != "" {
		if start == id.Start {
			return Match
		}
		return Mismatch
	}
	if SameBinary(exe, wantExe) {
		return Match
	}
	return Unverified
}

// SameBinary reports whether exe (as reported by the OS) is the program want names.
func SameBinary(exe, want string) bool {
	if exe == "" || want == "" {
		return false
	}
	b := filepath.Base(exe)
	if b == filepath.Base(want) {
		return true
	}
	r, err := filepath.EvalSymlinks(want)
	return err == nil && filepath.Base(r) == b
}

// ExpectedExe is the executable a running copy of s should show.
func ExpectedExe(s config.ServiceConfig) string {
	if ModeOf(s) == Embedded {
		exe, _ := os.Executable()
		return exe
	}
	c := config.ExpandPath(strings.TrimSpace(s.Command))
	if p, err := exec.LookPath(c); err == nil {
		return p
	}
	return c
}

func GetPIDFile(name string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mcp-local", "pids", name+".pid")
}

// SavePID records pid with its start time and executable so later commands can tell it
// apart from an unrelated process that reuses the PID.
func SavePID(name string, pid int) error {
	id := Identity{PID: pid}
	if start, exe, err := ProcInfo(pid); err == nil {
		id.Start, id.Exe = start, exe
	}
	path := GetPIDFile(name)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.Marshal(id)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// LoadIdentity reads a PID file in either the JSON or the legacy plain-PID format.
func LoadIdentity(name string) (Identity, error) {
	data, err := os.ReadFile(GetPIDFile(name))
	if err != nil {
		return Identity{}, err
	}
	s := strings.TrimSpace(string(data))
	var id Identity
	if strings.HasPrefix(s, "{") {
		if err := json.Unmarshal([]byte(s), &id); err != nil {
			return Identity{}, fmt.Errorf("pid file %s: %w", GetPIDFile(name), err)
		}
	} else if id.PID, err = strconv.Atoi(s); err != nil {
		return Identity{}, fmt.Errorf("pid file %s: %w", GetPIDFile(name), err)
	}
	if id.PID <= 0 {
		return Identity{}, fmt.Errorf("pid file %s: invalid pid %d", GetPIDFile(name), id.PID)
	}
	return id, nil
}

func LoadPID(name string) (int, error) {
	id, err := LoadIdentity(name)
	return id.PID, err
}

func RemovePID(name string) error {
	return os.Remove(GetPIDFile(name))
}
