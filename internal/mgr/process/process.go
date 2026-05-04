package process

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

type Process struct {
	Config config.ServiceConfig
	PID    int
}

func StartService(cfg config.ServiceConfig) (*Process, error) {
	// Dependency Guard
	for _, dep := range cfg.Deps {
		path := config.ExpandPath(dep)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, fmt.Errorf("dependency missing: %s", path)
		}
	}

	// Build if needed
	if cfg.BuildCmd != "" {
		fmt.Printf("Building %s...\n", cfg.Name)
		cmd := exec.Command("sh", "-c", cfg.BuildCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("build failed: %w", err)
		}
	}

	// Start process
	cmd := exec.Command("sh", "-c", config.ExpandPath(cfg.Command))
	cmd.Args = append(cmd.Args, cfg.Args...)
	
	env := os.Environ()
	for k, v := range cfg.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = env

	// Redirect logs to /tmp/
	logFile, err := os.OpenFile(fmt.Sprintf("/tmp/%s.log", cfg.Name), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start process: %w", err)
	}

	return &Process{
		Config: cfg,
		PID:    cmd.Process.Pid,
	}, nil
}

func StopService(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Signal(syscall.SIGTERM)
}

func GetPIDFile(name string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mcp-local", "pids", name+".pid")
}

func SavePID(name string, pid int) error {
	path := GetPIDFile(name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(strconv.Itoa(pid)), 0644)
}

func LoadPID(name string) (int, error) {
	data, err := os.ReadFile(GetPIDFile(name))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data))
}
