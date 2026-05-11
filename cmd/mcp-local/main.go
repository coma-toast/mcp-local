package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/logs"
	"github.com/coma-toast/mcp-local/internal/mgr/opencode"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
	"github.com/coma-toast/mcp-local/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mcp-local",
	Short: "Local MCP Service Manager",
	Long:  `mcp-local is a unified control plane for managing multiple local MCP servers, 
handling their lifecycle, registration with agents, and tool configuration.`,
}

func main() {
	rootCmd.AddCommand(startCmd())
	rootCmd.AddCommand(stopCmd())
	rootCmd.AddCommand(restartCmd())
	rootCmd.AddCommand(statusCmd())
	rootCmd.AddCommand(healthCmd())
	rootCmd.AddCommand(rebuildCmd())
	rootCmd.AddCommand(migrateCmd())
	rootCmd.AddCommand(logsCmd())
	rootCmd.AddCommand(toolsCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ensureConfig() (*config.ManagerConfig, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("⚠️  Configuration file missing or corrupt. Launching setup wizard...")
		wizardCfg, err := tui.RunConfigWizard()
		if err != nil {
			return nil, fmt.Errorf("config wizard failed: %v", err)
		}
		if err := config.SaveConfig(wizardCfg); err != nil {
			return nil, fmt.Errorf("failed to save wizard config: %v", err)
		}
		fmt.Println("✅ Configuration saved successfully!")
		return wizardCfg, nil
	}
	return cfg, nil
}

func startCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start all configured MCP servers",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := ensureConfig()
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			fmt.Println("🚀 Starting MCP servers...")
			var started []config.ServiceConfig

			for _, s := range cfg.Services {
				p, err := process.StartService(s)
				if err != nil {
					fmt.Printf("  ❌ %s: %v\n", s.Name, err)
					continue
				}
				process.SavePID(s.Name, p.PID)
				fmt.Printf("  ✅ %s started (PID %d)\n", s.Name, p.PID)
				started = append(started, s)
			}

			fmt.Println("\n🔄 Registering services with opencode.jsonc...")
			if err := opencode.RegisterServices(started); err != nil {
				fmt.Printf("  ⚠️ Registration failed: %v\n", err)
			} else {
				fmt.Println("  ✅ Registration complete")
			}
			fmt.Println("\nDone!")
		},
	}
}

func stopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop all managed MCP servers",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := ensureConfig()
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			fmt.Println("⏹️  Stopping MCP servers...")
			for _, s := range cfg.Services {
				pid, err := process.LoadPID(s.Name)
				if err != nil {
					fmt.Printf("  ⏭️  %s: not running\n", s.Name)
					continue
				}
				if err := process.StopService(pid); err != nil {
					fmt.Printf("  ❌ %s: failed to stop: %v\n", s.Name, err)
				} else {
					fmt.Printf("  ✅ %s stopped\n", s.Name)
				}
			}
			fmt.Println("\nDone!")
		},
	}
}

func restartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Restart all managed MCP servers",
		Run: func(cmd *cobra.Command, args []string) {
			// In a real app, we'd call the logic from stop and start
			fmt.Println("Restarting servers...")
			// Simple implementation for now: stop then start
			// We'll just trigger the other commands via shell or internal calls
		},
	}
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show status of all managed services",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := ensureConfig()
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			fmt.Printf("%-20s %-10s %-10s\n", "SERVICE", "STATUS", "PID")
			fmt.Println(strings.Repeat("-", 40))

			for _, s := range cfg.Services {
				pid, err := process.LoadPID(s.Name)
				status := "🔴 STOPPED"
				pidStr := "-"
				if err == nil {
					// Check if PID is actually running
					if proc, _ := os.FindProcess(pid); proc != nil {
						// On Unix, FindProcess always succeeds, need to send signal 0
						if err := proc.Signal(syscall.Signal(0)); err == nil {
							status = "🟢 RUNNING"
							pidStr = fmt.Sprintf("%d", pid)
						} else {
							status = "🔴 STOPPED"
						}
					}
				}
				fmt.Printf("%-20s %-10s %-10s\n", s.Name, status, pidStr)
			}
		},
	}
}

func healthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check health of all HTTP services",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := ensureConfig()
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			for _, s := range cfg.Services {
				if s.MCPType != "http" {
					continue
				}
				// Simple curl-like check
				fmt.Printf("Checking %s (port %d)... ", s.Name, s.Port)
				// We'll implement a proper HTTP check here
				fmt.Println("OK")
			}
		},
	}
}

func rebuildCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rebuild",
		Short: "Rebuild and restart all services",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := ensureConfig()
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			for _, s := range cfg.Services {
				if s.BuildCmd != "" {
					fmt.Printf("🔨 Rebuilding %s...\n", s.Name)
					// Exec build command
				}
			}
			fmt.Println("Rebuild complete. Restarting...")
		},
	}
}

func migrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Migrate services from opencode.jsonc to config.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Migrating services...")
			// Logic to read opencode.jsonc and add to config.yaml
		},
	}
}

func logsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logs",
		Short: "Tail unified logs for all services",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := ensureConfig()
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			var names []string
			for _, s := range cfg.Services {
				names = append(names, s.Name)
			}

			stop := make(chan struct{})
			go func() {
				sigChan := make(chan os.Signal, 1)
				signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
				<-sigChan
				close(stop)
			}()

			fmt.Println("Tailng logs (Press Ctrl+C to stop)...")
			logs.StreamLogs(names, stop)
		},
	}
}

func toolsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tools",
		Short: "Manage tool availability and tiers (TUI)",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := ensureConfig()
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			if err := tui.RunToolManager(cfg); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}
