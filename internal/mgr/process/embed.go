package process

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

func StartEmbedded(cfg config.ServiceConfig) (*Process, error) {
	if cfg.MCPType != "filesystem" {
		return nil, fmt.Errorf("StartEmbedded called for non-filesystem service %s", cfg.Name)
	}
	if cfg.Port <= 0 {
		return nil, fmt.Errorf("filesystem service %s must have a port configured", cfg.Name)
	}
	root := config.ExpandPath(cfg.Path)
	if root == "" {
		return nil, fmt.Errorf("filesystem service %s has no path configured", cfg.Name)
	}

	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("get executable: %w", err)
	}

	logPath := cfg.Log
	if logPath == "" {
		logPath = filepath.Join(os.TempDir(), cfg.Name+".log")
	}
	_ = os.MkdirAll(filepath.Dir(logPath), 0755)

	args := []string{
		"_serve-fs",
		"--name", cfg.Name,
		"--port", strconv.Itoa(cfg.Port),
		"--path", root,
		"--log", logPath,
	}
	cmd := exec.Command(exe, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	if runtime.GOOS != "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start embedded server: %w", err)
	}

	return &Process{Config: cfg, PID: cmd.Process.Pid}, nil
}
