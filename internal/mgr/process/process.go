package process

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/coma-toast/mcp-local/internal/mgr/asttools"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

type Process struct {
	Config config.ServiceConfig
	PID    int
}

// Mode is how mcp-local runs a service.
type Mode int

const (
	Daemon   Mode = iota // spawned and supervised by mcp-local
	Embedded             // filesystem server, spawned as `mcp-local _serve-fs`
	Remote               // mcp_url only: nothing to run, registered only
	Stdio                // launched on demand by each agent: registered only
)

func (m Mode) String() string {
	return [...]string{"daemon", "embedded", "remote", "stdio"}[m]
}

// ModeOf classifies s. stdio servers are never daemonized: with stdin at /dev/null they
// would exit at once, and every agent starts its own copy anyway.
func ModeOf(s config.ServiceConfig) Mode {
	switch {
	case s.MCPType == "filesystem":
		return Embedded
	case strings.TrimSpace(s.Command) == "" && strings.TrimSpace(s.MCPURL) != "":
		return Remote
	case !s.IsHTTP() && strings.TrimSpace(s.MCPURL) == "":
		return Stdio
	}
	return Daemon
}

// BuildTimeout caps a service's build command.
var BuildTimeout = 10 * time.Minute

// Build runs the service's effective build command, if any. It runs from the directory
// of the service's command when that is an absolute path to an existing directory.
func Build(cfg config.ServiceConfig) error {
	build := config.EffectiveBuild(cfg)
	if build == "" {
		return nil
	}
	fmt.Printf("Building %s...\n", cfg.Name)
	ctx, cancel := context.WithTimeout(context.Background(), BuildTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", "-c", build)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if c := config.ExpandPath(cfg.Command); filepath.IsAbs(c) {
		if fi, err := os.Stat(filepath.Dir(c)); err == nil && fi.IsDir() {
			cmd.Dir = filepath.Dir(c)
		}
	}
	if err := cmd.Run(); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("build timed out after %s", BuildTimeout)
		}
		return fmt.Errorf("build failed: %w", err)
	}
	return nil
}

// LogPath is the service's configured log, else ~/.mcp-local/logs/<name>.log.
func LogPath(cfg config.ServiceConfig) string {
	if cfg.Log != "" {
		return config.ExpandPath(cfg.Log)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mcp-local", "logs", cfg.Name+".log")
}

// openLog opens the service log for appending, creating it 0600. The default log is
// opened without following symlinks so a planted link can't redirect it.
func openLog(cfg config.ServiceConfig) (*os.File, string, error) {
	path := LogPath(cfg)
	flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	if cfg.Log == "" {
		flags |= syscall.O_NOFOLLOW
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, path, fmt.Errorf("create log dir: %w", err)
	}
	f, err := os.OpenFile(path, flags, 0o600)
	if err != nil {
		return nil, path, fmt.Errorf("open log file: %w", err)
	}
	if cfg.Log == "" {
		_ = f.Chmod(0o600)
	}
	return f, path, nil
}

// StartService spawns cfg.Command in its own session with output to the service log.
// It does not build; call Build first.
func StartService(cfg config.ServiceConfig) (*Process, error) {
	if strings.TrimSpace(cfg.Command) == "" {
		return nil, fmt.Errorf("no command configured for %s", cfg.Name)
	}
	for _, dep := range cfg.Deps {
		path := config.ExpandPath(dep)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, fmt.Errorf("dependency missing: %s", path)
		}
	}
	logFile, _, err := openLog(cfg)
	if err != nil {
		return nil, err
	}
	defer logFile.Close()
	cmd := exec.Command(config.ExpandPath(cfg.Command), cfg.Args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	env := os.Environ()
	for k, v := range cfg.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	env, err = asttools.ApplyStartEnv(env, cfg)
	if err != nil {
		return nil, err
	}
	if cfg.NoCodeMode {
		env = append(env, "AST_MCP_CODE_MODE=false")
	}
	cmd.Env = env
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start process: %w", err)
	}
	return &Process{Config: cfg, PID: cmd.Process.Pid}, nil
}
