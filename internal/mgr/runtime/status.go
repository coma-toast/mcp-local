package runtime

import (
	"fmt"
	"path/filepath"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
	"github.com/coma-toast/mcp-local/internal/portutil"
)

type ServiceStatus struct {
	Name    string
	State   config.StateKind
	Detail  string
	Running bool
	PID     int
	Port    int
	Uptime  string
	Mode    process.Mode
	// Stale is set when a PID file exists but no longer names the process mcp-local
	// started. Status is read-only, so the file is left for stop or doctor --fix.
	Stale bool
	// Orphan is an untracked copy of the service's own binary holding one of its ports
	// while the main port is down: start replaces it and doctor --fix stops it.
	Orphan Holder
}

// IsPIDAlive reports whether pid exists; EPERM (another user's process) counts as alive.
func IsPIDAlive(pid int) bool {
	return pid > 0 && process.Alive(pid)
}

// Holder is the process listening on one of a service's ports.
type Holder struct {
	Port int
	PID  int
	Exe  string
	Same bool // runs the service's own binary
}

// Lookups used by CheckStatus; variables so tests can fake them.
var (
	portListening = portutil.IsRunning
	portOwner     = portutil.FindPID
	procInfo      = process.ProcInfo
	fetchHealth   = portutil.FetchHealth
)

// HolderOf reports who listens on port, or ok=false when nothing does.
func HolderOf(svc config.ServiceConfig, port int) (Holder, bool) {
	if port <= 0 || !portListening(port) {
		return Holder{}, false
	}
	h := Holder{Port: port, PID: portOwner(port)}
	if h.PID > 0 {
		if _, exe, err := procInfo(h.PID); err == nil {
			h.Exe = exe
			h.Same = process.SameBinary(exe, process.ExpectedExe(svc))
		}
	}
	return h, true
}

func CheckStatus(svc config.ServiceConfig) ServiceStatus {
	st := ServiceStatus{Name: svc.Name, Port: svc.Port, Mode: process.ModeOf(svc), State: config.StateStopped, Detail: "stopped"}
	if h, ok := HolderOf(svc, svc.Port); ok {
		st.setPID(h.PID)
		st.Running, st.State, st.Detail = true, config.StateRunning, "healthy"
		if healthy, why := fetchHealth(svc.HealthURL); !healthy {
			st.State, st.Detail = config.StateDegraded, why
		}
		return st
	}
	if id, err := process.LoadIdentity(svc.Name); err == nil {
		if id.Check(process.ExpectedExe(svc)) == process.Match {
			st.setPID(id.PID)
			if svc.Port <= 0 {
				st.Running, st.State, st.Detail = true, config.StateRunning, "running"
				return st
			}
			// Possibly still starting up, so only restart (explicit) replaces it.
			st.State = config.StateDegraded
			st.Detail = fmt.Sprintf("process alive (pid %d) but port %d not listening", id.PID, svc.Port)
			return st
		}
		st.Stale, st.State = true, config.StateStale
		st.Detail = fmt.Sprintf("PID file names pid %d, which is not the process mcp-local started", id.PID)
	}
	for _, p := range svc.OwnedPorts() {
		if p == svc.Port {
			continue
		}
		h, ok := HolderOf(svc, p)
		if !ok {
			continue
		}
		st.setPID(h.PID)
		if h.Same {
			st.Orphan, st.State = h, config.StateOrphaned
			st.Detail = fmt.Sprintf("port %d held by an untracked %s (pid %d) while port %d is down", p, describeExe(h.Exe), h.PID, svc.Port)
		} else {
			st.State = config.StateConflict
			st.Detail = fmt.Sprintf("port %d held by another program, %s (pid %d)", p, describeExe(h.Exe), h.PID)
		}
		return st
	}
	return st
}

func (st *ServiceStatus) setPID(pid int) {
	st.PID = pid
	if pid > 0 {
		st.Uptime = portutil.GetUptime(pid)
	}
}

func describeExe(exe string) string {
	if exe == "" {
		return "an unknown program"
	}
	return filepath.Base(exe)
}

func CheckAll(cfg *config.ManagerConfig, filter string) []ServiceStatus {
	var results []ServiceStatus
	for _, name := range cfg.SortedNames() {
		if filter != "" && name != filter {
			continue
		}
		svc, err := cfg.ServiceNamed(name)
		if err != nil {
			continue
		}
		results = append(results, CheckStatus(svc))
	}
	return results
}
