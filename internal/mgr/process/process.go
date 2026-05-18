package process

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
	for _, dep := range cfg.Deps {
		path := config.ExpandPath(dep)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, fmt.Errorf("dependency missing: %s", path)
		}
	}

	build := config.EffectiveBuild(cfg)
	if build != "" {
		fmt.Printf("Building %s...\n", cfg.Name)
		cmd := exec.Command("sh", "-c", build)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("build failed: %w", err)
		}
	}

	logPath := cfg.Log
	if logPath == "" {
		logPath = filepath.Join(os.TempDir(), cfg.Name+".log")
	}
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	cmd := exec.Command(config.ExpandPath(cfg.Command), cfg.Args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	env := os.Environ()
	for k, v := range cfg.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	if cfg.ActiveTier != "" {
		env = append(env, fmt.Sprintf("AST_MCP_TIER=%s", cfg.ActiveTier))
	}
	if cfg.CodeMode {
		env = append(env, "AST_MCP_CODE_MODE=true")
	}

	var disabledTools []string
	for _, t := range cfg.Tools {
		if !t.Enabled {
			disabledTools = append(disabledTools, t.Name)
		}
	}
	if len(disabledTools) > 0 {
		env = append(env, fmt.Sprintf("AST_MCP_DISABLED_TOOLS=%s", strings.Join(disabledTools, ",")))
	}

	cmd.Env = env

	if runtime.GOOS != "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	}

	if err := cmd.Start(); err != nil {
		_ = logFile.Close()
		return nil, fmt.Errorf("failed to start process: %w", err)
	}
	_ = logFile.Close()

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

func RemovePID(name string) error {
	return os.Remove(GetPIDFile(name))
}
