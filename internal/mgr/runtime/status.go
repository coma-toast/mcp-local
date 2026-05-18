package runtime

import (
	"os"
	"syscall"

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
		Name: svc.Name,
		Port: svc.Port,
	}

	if svc.Port > 0 {
		if portutil.IsRunning(svc.Port) {
			status.Running = true
			pid := portutil.FindPID(svc.Port)
			if pid > 0 {
				status.PID = pid
				status.Uptime = portutil.GetUptime(pid)
			}
			return status
		}
	}

	if pid, err := process.LoadPID(svc.Name); err == nil {
		if IsPIDAlive(pid) {
			status.Running = true
			status.PID = pid
			status.Uptime = portutil.GetUptime(pid)
			return status
		}
		_ = process.RemovePID(svc.Name)
	}

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
