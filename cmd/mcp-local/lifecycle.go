package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/coma-toast/mcp-local/internal/mgr/agents"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
	"github.com/coma-toast/mcp-local/internal/mgr/runtime"
	"github.com/coma-toast/mcp-local/internal/portutil"
	"github.com/spf13/cobra"
)

const (
	portWait         = 5 * time.Second
	maxStopAttempts  = 30
	stopPollInterval = 100 * time.Millisecond
	restartDelay     = 300 * time.Millisecond
)

func resolveTargets(cfg *config.ManagerConfig, args []string, all bool, verb string) ([]string, error) {
	if all {
		return cfg.SortedNames(), nil
	}
	if len(args) != 1 {
		return nil, fmt.Errorf("usage: mcp-local %s --all | %s <service>", verb, verb)
	}
	if _, err := cfg.ServiceNamed(args[0]); err != nil {
		return nil, err
	}
	return []string{args[0]}, nil
}

func stopByPID(name string, pid int, svc config.ServiceConfig, deregister bool, targets agents.Targets) {
	proc, _ := os.FindProcess(pid)
	_ = proc.Signal(syscall.SIGTERM)
	for i := 0; i < maxStopAttempts; i++ {
		time.Sleep(stopPollInterval)
		// Check if any owned port is still running or process is dead
		stillRunning := false
		for _, p := range svc.OwnedPorts() {
			if portutil.IsRunning(p) {
				stillRunning = true
				break
			}
		}
		if !runtime.IsPIDAlive(pid) || !stillRunning {
			_ = process.RemovePID(name)
			fmt.Printf("  ✅ %s stopped (pid %d)\n", name, pid)
			if deregister {
				warnDeregister(name, agents.DeregisterAll(svc, targets))
			}
			return
		}
	}
	_ = proc.Signal(syscall.SIGKILL)
	_ = process.RemovePID(name)
	fmt.Printf("  ✅ %s killed (pid %d)\n", name, pid)
	if deregister {
		warnDeregister(name, agents.DeregisterAll(svc, targets))
	}
}

func warnDeregister(name string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠️  %s: deregister failed:\n%s\n", name, indent(err.Error()))
	}
}

func indent(s string) string {
	return "    " + strings.ReplaceAll(s, "\n", "\n    ")
}

func stopServiceNamed(cfg *config.ManagerConfig, name string, deregister bool) error {
	svc, err := cfg.ServiceNamed(name)
	if err != nil {
		return err
	}
	targets := agents.TargetsFromConfig(*cfg)

	// 1. Check PID file
	if pid, err := process.LoadPID(name); err == nil && runtime.IsPIDAlive(pid) {
		stopByPID(name, pid, svc, deregister, targets)
		return nil
	} else {
		_ = process.RemovePID(name)
	}

	// 2. Check any owned port for same-binary instance
	for _, p := range svc.OwnedPorts() {
		if pid := portutil.FindPID(p); pid > 0 {
			if portutil.IsSameBinary(pid, svc.Command) {
				fmt.Printf("  ℹ️  Found untracked %s instance on port %d (PID %d). Adopting and stopping...\n", name, p, pid)
				stopByPID(name, pid, svc, deregister, targets)
				return nil
			} else {
				cmdName := portutil.GetCommandForPID(pid)
				return fmt.Errorf("port %d held by different process %q (PID %d)", p, cmdName, pid)
			}
		}
	}

	fmt.Printf("  ⏭️  %s: not running\n", name)
	return nil
}

func printLastLogLines(logPath string, n int) {
	f, err := os.Open(logPath)
	if err != nil {
		return
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	start := 0
	if len(lines) > n {
		start = len(lines) - n
	}
	fmt.Fprintf(os.Stderr, "  📄 Last %d lines of log (%s):\n", len(lines)-start, logPath)
	for i := start; i < len(lines); i++ {
		fmt.Fprintf(os.Stderr, "    %s\n", lines[i])
	}
}

func startServices(cfg *config.ManagerConfig, names []string, cmd *cobra.Command) ([]config.ServiceConfig, error) {
	var started []config.ServiceConfig
	fmt.Println("🚀 Starting MCP servers...")
	targets := agents.TargetsFromConfig(*cfg)

	for _, name := range names {
		svc, err := cfg.ServiceNamed(name)
		if err != nil {
			fmt.Printf("  ❌ %v\n", err)
			continue
		}

		// Check status using runtime CheckStatus
		st := runtime.CheckStatus(svc)
		if st.State == config.StateRunning {
			fmt.Printf("  ⏭️  %s: already running on port %d (PID %d)\n", name, svc.Port, st.PID)
			started = append(started, svc)
			continue
		}

		if st.State == config.StateOrphaned || st.State == config.StateDegraded {
			fmt.Printf("  ⚠️  %s: found stuck/orphaned instance (PID %d, state: %s). Stopping before start...\n", name, st.PID, st.State)
			if st.PID > 0 {
				stopByPID(name, st.PID, svc, false, targets)
			} else {
				// Kill by owned ports
				for _, p := range svc.OwnedPorts() {
					if pid := portutil.FindPID(p); pid > 0 && portutil.IsSameBinary(pid, svc.Command) {
						stopByPID(name, pid, svc, false, targets)
					}
				}
			}
		} else if st.State == config.StateConflict {
			fmt.Printf("  ❌ %s: port conflict - %s\n", name, st.Detail)
			continue
		}

		if svc.MCPType == "filesystem" {
			p, err := process.StartEmbedded(svc)
			if err != nil {
				fmt.Printf("  ❌ %s: %v\n", name, err)
				continue
			}
			process.SavePID(name, p.PID)
			fmt.Printf("  ✅ %s started (in-process, PID %d)\n", name, p.PID)
			if svc.Port > 0 && !portutil.WaitForPort(svc.Port, portWait) {
				fmt.Fprintf(cmd.ErrOrStderr(), "  ⚠️  %s: port %d not listening yet\n", name, svc.Port)
				if svc.Log != "" {
					printLastLogLines(svc.Log, 10)
				}
			}
			started = append(started, svc)
			continue
		}

		svcForStart := svc
		if !config.ShouldBuildOnStart(svcForStart) {
			svcForStart.BuildCmd = ""
			svcForStart.BuildCommand = ""
		}
		p, err := process.StartService(svcForStart)
		if err != nil {
			fmt.Printf("  ❌ %s: %v\n", name, err)
			if svc.Log != "" {
				printLastLogLines(svc.Log, 10)
			}
			continue
		}
		if err := process.SavePID(name, p.PID); err != nil {
			fmt.Printf("  ⚠️  %s: save pid: %v\n", name, err)
		}
		fmt.Printf("  ✅ %s started (PID %d)\n", name, p.PID)

		// Wait for main port if configured
		if svc.Port > 0 {
			if !portutil.WaitForPort(svc.Port, portWait) {
				fmt.Fprintf(cmd.ErrOrStderr(), "  ⚠️  %s: port %d not listening yet\n", name, svc.Port)
				if svc.Log != "" {
					printLastLogLines(svc.Log, 10)
				}
			}
		}
		started = append(started, svc)
	}
	return started, nil
}

func registerStarted(cfg *config.ManagerConfig, started []config.ServiceConfig) {
	if len(started) == 0 {
		return
	}
	targets := agents.TargetsFromConfig(*cfg)
	fmt.Println("\n🔄 Registering services with agents...")
	if err := agents.RegisterAll(started, targets); err != nil {
		fmt.Printf("  ⚠️  Registration failed:\n%s\n", indent(err.Error()))
	} else {
		if targets.OpenCode {
			fmt.Println("  ✅ OpenCode registration complete")
		}
		if targets.Cursor {
			fmt.Println("  ✅ Cursor registration complete")
		}
		if targets.Claude {
			fmt.Println("  ✅ Claude registration complete")
		}
	}
}

func addLifecycle() {
	var startAll bool
	start := &cobra.Command{
		Use:   "start [service]",
		Short: "Start service(s) and register with agents",
		Example: `  mcp-local start my-server
  mcp-local start --all`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			names, err := resolveTargets(cfg, args, startAll, "start")
			if err != nil {
				return err
			}
			started, err := startServices(cfg, names, cmd)
			if err != nil {
				return err
			}
			registerStarted(cfg, started)
			fmt.Println("\nDone!")
			return nil
		},
	}
	start.Flags().BoolVar(&startAll, "all", false, "start every configured service")

	var stopAll bool
	var stopDeregister bool
	stop := &cobra.Command{
		Use:   "stop [service]",
		Short: "Stop service(s)",
		Example: `  mcp-local stop my-server
  mcp-local stop --all
  mcp-local stop --deregister=false my-server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			names, err := resolveTargets(cfg, args, stopAll, "stop")
			if err != nil {
				return err
			}
			fmt.Println("⏹️  Stopping MCP servers...")
			for _, name := range names {
				if err := stopServiceNamed(cfg, name, stopDeregister); err != nil {
					fmt.Printf("  ❌ %s: %v\n", name, err)
				}
			}
			fmt.Println("\nDone!")
			return nil
		},
	}
	stop.Flags().BoolVar(&stopAll, "all", false, "stop every configured service")
	stop.Flags().BoolVar(&stopDeregister, "deregister", true, "deregister from agents after stopping")

	var restartAll bool
	restart := &cobra.Command{
		Use:   "restart [service]",
		Short: "Restart service(s)",
		Example: `  mcp-local restart my-server
  mcp-local restart --all`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			names, err := resolveTargets(cfg, args, restartAll, "restart")
			if err != nil {
				return err
			}
			for _, name := range names {
				_ = stopServiceNamed(cfg, name, false)
			}
			time.Sleep(restartDelay)
			started, err := startServices(cfg, names, cmd)
			if err != nil {
				return err
			}
			registerStarted(cfg, started)
			fmt.Println("\nDone!")
			return nil
		},
	}
	restart.Flags().BoolVar(&restartAll, "all", false, "restart every configured service")

	rootCmd.AddCommand(start, stop, restart)
}
