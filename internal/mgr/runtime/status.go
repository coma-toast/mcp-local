package runtime

import (
	"fmt"
	"os"
	"syscall"

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
}

func IsPIDAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}

func CheckStatus(svc config.ServiceConfig) ServiceStatus {
	status := ServiceStatus{
		Name:   svc.Name,
		Port:   svc.Port,
		State:  config.StateStopped,
		Detail: "",
	}

	ownedPorts := svc.OwnedPorts()
	if len(ownedPorts) == 0 && svc.Command == "" && svc.MCPURL == "" {
		return status
	}

	// 1. Check tracked PID file first
	var trackedPID int
	if pid, err := process.LoadPID(svc.Name); err == nil && IsPIDAlive(pid) {
		trackedPID = pid
	} else if err == nil {
		_ = process.RemovePID(svc.Name)
	}

	// 2. Check main port status
	mainPortRunning := svc.Port > 0 && portutil.IsRunning(svc.Port)
	var mainPortPID int
	if mainPortRunning {
		mainPortPID = portutil.FindPID(svc.Port)
	}

	// If main port is running
	if mainPortRunning {
		pid := mainPortPID
		if pid == 0 {
			pid = trackedPID
		}
		status.PID = pid
		if pid > 0 {
			status.Uptime = portutil.GetUptime(pid)
		}

		// Check health URL if configured
		if svc.HealthURL != "" {
			healthy, healthDetail := portutil.FetchHealth(svc.HealthURL)
			if !healthy {
				status.State = config.StateDegraded
				status.Running = true
				status.Detail = healthDetail
				return status
			}
		}

		status.State = config.StateRunning
		status.Running = true
		status.Detail = "healthy"
		return status
	}

	// If main port is not running, but we have a tracked PID or any owned port is running
	anyPortRunning := false
	var runningPort int
	var runningPortPID int
	for _, p := range ownedPorts {
		if portutil.IsRunning(p) {
			anyPortRunning = true
			runningPort = p
			runningPortPID = portutil.FindPID(p)
			break
		}
	}

	if trackedPID > 0 || anyPortRunning {
		pid := trackedPID
		if pid == 0 {
			pid = runningPortPID
		}
		status.PID = pid
		if pid > 0 {
			status.Uptime = portutil.GetUptime(pid)
		}
		status.Running = true

		// Check if it's same binary or conflict
		if pid > 0 && portutil.IsSameBinary(pid, svc.Command) {
			status.State = config.StateDegraded
			if runningPort > 0 && runningPort != svc.Port {
				status.Detail = fmt.Sprintf("main port %d not listening, but auxiliary port %d is held by PID %d", svc.Port, runningPort, pid)
			} else {
				status.Detail = fmt.Sprintf("process alive (PID %d) but main port %d not listening", pid, svc.Port)
			}
		} else if pid > 0 {
			cmdName := portutil.GetCommandForPID(pid)
			status.State = config.StateConflict
			status.Detail = fmt.Sprintf("port %d held by different process %q (PID %d)", svc.Port, cmdName, pid)
		} else {
			status.State = config.StateDegraded
			status.Detail = "service process is not responding on main port"
		}
		return status
	}

	// Check if any owned port is held by an untracked process
	for _, p := range ownedPorts {
		if portutil.IsRunning(p) {
			pid := portutil.FindPID(p)
			if pid > 0 {
				status.PID = pid
				status.Uptime = portutil.GetUptime(pid)
				if portutil.IsSameBinary(pid, svc.Command) {
					status.State = config.StateOrphaned
					status.Detail = fmt.Sprintf("owned port %d held by untracked instance of same binary (PID %d)", p, pid)
				} else {
					cmdName := portutil.GetCommandForPID(pid)
					status.State = config.StateConflict
					status.Detail = fmt.Sprintf("port %d held by different program %q (PID %d)", p, cmdName, pid)
				}
				return status
			}
		}
	}

	status.State = config.StateStopped
	status.Detail = "stopped"
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
