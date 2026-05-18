package main

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/coma-toast/mcp-local/internal/mgr/agents"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
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
		if svc.Port <= 0 || !portutil.IsRunning(svc.Port) {
			_ = process.RemovePID(name)
			fmt.Printf("  ✅ %s stopped (pid %d)\n", name, pid)
			if deregister {
				_ = agents.DeregisterAll(name, targets)
			}
			return
		}
	}
	_ = proc.Signal(syscall.SIGKILL)
	_ = process.RemovePID(name)
	fmt.Printf("  ✅ %s killed (pid %d)\n", name, pid)
	if deregister {
		_ = agents.DeregisterAll(name, targets)
	}
}

func stopServiceNamed(cfg *config.ManagerConfig, name string, deregister bool) error {
	svc, err := cfg.ServiceNamed(name)
	if err != nil {
		return err
	}
	targets := agents.TargetsFromConfig(*cfg)
	if pid, err := process.LoadPID(name); err == nil {
		stopByPID(name, pid, svc, deregister, targets)
		return nil
	}
	pid := portutil.FindPID(svc.Port)
	if pid == 0 {
		fmt.Printf("  ⏭️  %s: not running\n", name)
		return nil
	}
	stopByPID(name, pid, svc, deregister, targets)
	return nil
}

func startServices(cfg *config.ManagerConfig, names []string, cmd *cobra.Command) ([]config.ServiceConfig, error) {
	var started []config.ServiceConfig
	fmt.Println("🚀 Starting MCP servers...")
	for _, name := range names {
		svc, err := cfg.ServiceNamed(name)
		if err != nil {
			fmt.Printf("  ❌ %v\n", err)
			continue
		}
		if svc.Port > 0 && portutil.IsRunning(svc.Port) {
			fmt.Printf("  ⏭️  %s: already running on port %d\n", name, svc.Port)
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
			continue
		}
		if err := process.SavePID(name, p.PID); err != nil {
			fmt.Printf("  ⚠️  %s: save pid: %v\n", name, err)
		}
		fmt.Printf("  ✅ %s started (PID %d)\n", name, p.PID)
		if svc.Port > 0 && !portutil.WaitForPort(svc.Port, portWait) {
			fmt.Fprintf(cmd.ErrOrStderr(), "  ⚠️  %s: port %d not listening yet\n", name, svc.Port)
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
		fmt.Printf("  ⚠️ Registration failed: %v\n", err)
	} else {
		if targets.OpenCode {
			fmt.Println("  ✅ OpenCode registration complete")
		}
		if targets.Cursor {
			fmt.Println("  ✅ Cursor registration complete")
		}
	}
}

func addLifecycle() {
	var startAll bool
	start := &cobra.Command{
		Use:   "start [service]",
		Short: "Start service(s) and register with agents",
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
