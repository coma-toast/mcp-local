package runtime

import (
	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
	"github.com/coma-toast/mcp-local/internal/portutil"
)

type ServiceStatus struct {
	Name    string
	Running bool
	PID     int
	Port    int
	Uptime  string
	Mode    process.Mode
	// Stale is set when a PID file exists but no longer names the process mcp-local
	// started. Status is read-only, so the file is left for stop (or a future doctor).
	Stale bool
}

// IsPIDAlive reports whether pid exists; EPERM (another user's process) counts as alive.
func IsPIDAlive(pid int) bool {
	return pid > 0 && process.Alive(pid)
}

func CheckStatus(svc config.ServiceConfig) ServiceStatus {
	status := ServiceStatus{Name: svc.Name, Port: svc.Port, Mode: process.ModeOf(svc)}
	if svc.Port > 0 && portutil.IsRunning(svc.Port) {
		status.Running = true
		if pid := portutil.FindPID(svc.Port); pid > 0 {
			status.PID = pid
			status.Uptime = portutil.GetUptime(pid)
		}
		return status
	}
	id, err := process.LoadIdentity(svc.Name)
	if err != nil {
		return status
	}
	if id.Check(process.ExpectedExe(svc)) != process.Match {
		status.Stale = true
		return status
	}
	status.Running = true
	status.PID = id.PID
	status.Uptime = portutil.GetUptime(id.PID)
	return status
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
