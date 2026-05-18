package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/coma-toast/mcp-local/internal/mgr/agents"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
	"github.com/coma-toast/mcp-local/internal/portutil"
	"github.com/spf13/cobra"
)

const portWait = 5 * time.Second

func resolveTargets(cfg *config.ManagerConfig, args []string, all bool, verb string) []string {
	if all {
		return cfg.SortedNames()
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: mcp-local %s --all | %s <service>\n", verb, verb)
		os.Exit(1)
	}
	if _, err := cfg.ServiceNamed(args[0]); err != nil {
		log.Fatal(err)
	}
	return []string{args[0]}
}

func stopServiceNamed(cfg *config.ManagerConfig, name string, deregister bool) error {
	svc, err := cfg.ServiceNamed(name)
	if err != nil {
		return err
	}
	targets := agents.TargetsFromConfig(*cfg)
	if pid, err := process.LoadPID(name); err == nil {
		proc, _ := os.FindProcess(pid)
		_ = proc.Signal(syscall.SIGTERM)
		for i := 0; i < 30; i++ {
			time.Sleep(100 * time.Millisecond)
			if svc.Port <= 0 || !portutil.IsRunning(svc.Port) {
				_ = process.RemovePID(name)
				fmt.Printf("  ✅ %s stopped (pid %d)\n", name, pid)
				if deregister {
					_ = agents.DeregisterAll(name, targets)
				}
				return nil
			}
		}
		_ = proc.Signal(syscall.SIGKILL)
		_ = process.RemovePID(name)
		fmt.Printf("  ✅ %s killed (pid %d)\n", name, pid)
		if deregister {
			_ = agents.DeregisterAll(name, targets)
		}
		return nil
	}
	pid := portutil.FindPID(svc.Port)
	if pid == 0 {
		fmt.Printf("  ⏭️  %s: not running\n", name)
		return nil
	}
	proc, _ := os.FindProcess(pid)
	_ = proc.Signal(syscall.SIGTERM)
	for i := 0; i < 30; i++ {
		time.Sleep(100 * time.Millisecond)
		if svc.Port <= 0 || !portutil.IsRunning(svc.Port) {
			fmt.Printf("  ✅ %s stopped (pid %d)\n", name, pid)
			if deregister {
				_ = agents.DeregisterAll(name, targets)
			}
			return nil
		}
	}
	_ = proc.Signal(syscall.SIGKILL)
	fmt.Printf("  ✅ %s killed (pid %d)\n", name, pid)
	if deregister {
		_ = agents.DeregisterAll(name, targets)
	}
	return nil
}

func addLifecycle() {
	var startAll bool
	start := &cobra.Command{
		Use:   "start [service]",
		Short: "Start service(s) and register with agents",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := mustConfig()
			names := resolveTargets(cfg, args, startAll, "start")
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
			if len(started) > 0 {
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
			fmt.Println("\nDone!")
		},
	}
	start.Flags().BoolVar(&startAll, "all", false, "start every configured service")

	var stopAll bool
	var stopDeregister bool
	stop := &cobra.Command{
		Use:   "stop [service]",
		Short: "Stop service(s)",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := mustConfig()
			names := resolveTargets(cfg, args, stopAll, "stop")
			fmt.Println("⏹️  Stopping MCP servers...")
			for _, name := range names {
				if err := stopServiceNamed(cfg, name, stopDeregister); err != nil {
					fmt.Printf("  ❌ %s: %v\n", name, err)
				}
			}
			fmt.Println("\nDone!")
		},
	}
	stop.Flags().BoolVar(&stopAll, "all", false, "stop every configured service")
	stop.Flags().BoolVar(&stopDeregister, "deregister", true, "deregister from agents after stopping")

	var restartAll bool
	restart := &cobra.Command{
		Use:   "restart [service]",
		Short: "Restart service(s)",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := mustConfig()
			names := resolveTargets(cfg, args, restartAll, "restart")
			for _, name := range names {
				_ = stopServiceNamed(cfg, name, false)
			}
			time.Sleep(300 * time.Millisecond)
			var started []config.ServiceConfig
			fmt.Println("🚀 Starting MCP servers...")
			for _, name := range names {
				svc, err := cfg.ServiceNamed(name)
				if err != nil {
					fmt.Printf("  ❌ %v\n", err)
					continue
				}
				if svc.Port > 0 && portutil.IsRunning(svc.Port) {
					fmt.Printf("  ⏭️  %s: already running\n", name)
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
				_ = process.SavePID(name, p.PID)
				fmt.Printf("  ✅ %s started (PID %d)\n", name, p.PID)
				if svc.Port > 0 {
					_ = portutil.WaitForPort(svc.Port, portWait)
				}
				started = append(started, svc)
			}
			if len(started) > 0 {
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
			fmt.Println("\nDone!")
		},
	}
	restart.Flags().BoolVar(&restartAll, "all", false, "restart every configured service")

	rootCmd.AddCommand(start, stop, restart)
}
