package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/coma-toast/mcp-local/internal/mgr/agents"
	"github.com/coma-toast/mcp-local/internal/mgr/claudedesktop"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/cursor"
	"github.com/coma-toast/mcp-local/internal/mgr/opencode"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
	"github.com/coma-toast/mcp-local/internal/mgr/runtime"
	"github.com/coma-toast/mcp-local/internal/portutil"
	"github.com/spf13/cobra"
)

func addDoctor() {
	doctor := &cobra.Command{
		Use:   "doctor [service]",
		Short: "Diagnose and optionally fix stuck service issues",
		Example: `  mcp-local doctor
  mcp-local doctor my-server
  mcp-local doctor my-server --fix`,
		RunE: func(cmd *cobra.Command, args []string) error {
			forceFix, _ := cmd.Flags().GetBool("fix")
			name := ""
			if len(args) > 0 {
				name = args[0]
			}
			doctorRun(name, forceFix)
			return nil
		},
	}
	doctor.Flags().BoolP("fix", "f", false, "apply safe fixes automatically")
	rootCmd.AddCommand(doctor)
}

func doctorRun(name string, forceFix bool) {
	cfg := mustConfig()
	targets := agents.TargetsFromConfig(*cfg)

	var svc config.ServiceConfig
	if name != "" {
		var err error
		svc, err = cfg.ServiceNamed(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Service %q not found: %v\n", name, err)
			return
		}
	}

	if name != "" {
		doctorService(cfg, svc, targets, forceFix)
	} else {
		// Run doctor on all services
		for _, sName := range cfg.SortedNames() {
			svc, _ := cfg.ServiceNamed(sName)
			doctorService(cfg, svc, targets, forceFix)
		}
	}
}

func doctorService(cfg *config.ManagerConfig, svc config.ServiceConfig, targets agents.Targets, forceFix bool) {
	fmt.Printf("\n🔍 Doctor examining: %s\n", svc.Name)
	fmt.Printf("   command: %s\n", svc.Command)
	if svc.Port > 0 {
		fmt.Printf("   port: %d\n", svc.Port)
	}
	if svc.HealthURL != "" {
		fmt.Printf("   health_url: %s\n", svc.HealthURL)
	}
	if svc.DashboardURL != "" {
		fmt.Printf("   dashboard_url: %s\n", svc.DashboardURL)
	}
	if svc.MCPURL != "" {
		fmt.Printf("   mcp_url: %s\n", svc.MCPURL)
	}
	if len(svc.ExtraPorts) > 0 {
		fmt.Printf("   extra_ports: %v\n", svc.ExtraPorts)
	}

	ownedPorts := svc.OwnedPorts()

	// 1. Check PID file status
	if pid, err := process.LoadPID(svc.Name); err == nil {
		if runtime.IsPIDAlive(pid) {
			fmt.Printf("   PID file: %s — PID %d (alive)\n", getPIDFilePath(svc.Name), pid)
		} else {
			fmt.Printf("   PID file: %s — PID %d (stale)\n", getPIDFilePath(svc.Name), pid)
			process.RemovePID(svc.Name)
		}
	} else {
		fmt.Printf("   PID file: missing (no PID file configured)\n")
	}

	// 2. Check owned ports for orphans and port conflicts
	fmt.Println("   Owned ports:")
	for _, p := range ownedPorts {
		if portutil.IsRunning(p) {
			pid := portutil.FindPID(p)
			if pid > 0 {
				cmdName := portutil.GetCommandForPID(pid)
				if portutil.IsSameBinary(pid, svc.Command) {
					fmt.Printf("   - port %d: held by %s PID %d (same binary) — orphaned\n", p, cmdName, pid)
					if forceFix {
						// Stop same-binary orphan
						fmt.Printf("   ✅  FIX: Stopping same-binary orphan (PID %d) on port %d\n", pid, p)
						_ = process.RemovePID(svc.Name)
					}
				} else {
					fmt.Printf("   - port %d: held by %s PID %d — port conflict\n", p, cmdName, pid)
					if forceFix {
						fmt.Printf("   ❌  Cannot fix: port %d held by different program %q (PID %d). Not killing.\n", p, cmdName, pid)
					}
				}
			}
		} else {
			fmt.Printf("   - port %d: not listening\n", p)
		}
	}

	// 3. Check health
	if svc.HealthURL != "" {
		fmt.Println("   Health check:")
		healthy, healthDetail := portutil.FetchHealth(svc.HealthURL)
		if healthy {
			fmt.Printf("   ✓ healthy\n")
		} else {
			fmt.Printf("   ✗ %s\n", healthDetail)
			if forceFix {
				// Attempt to restart if degraded
				fmt.Printf("   ℹ️  Attempting restart...\n")
				_ = stopServiceNamed(cfg, svc.Name, true)
				time.Sleep(restartDelay)
				_, _ = startServices(cfg, []string{svc.Name}, &cobra.Command{})
				fmt.Printf("   ✅  Restart attempted\n")
			}
		}
	} else {
		fmt.Println("   No health_url configured")
	}

	// 4. Check missing deps
	fmt.Println("   Dependencies:")
	for _, dep := range svc.Deps {
		expanded := config.ExpandPath(dep)
		if _, err := os.Stat(expanded); os.IsNotExist(err) {
			fmt.Printf("   ✗ missing: %s\n", dep)
		} else {
			fmt.Printf("   ✓ %s\n", dep)
		}
	}

	// 5. Check stale PID file (second check, focused)
	if pid, err := process.LoadPID(svc.Name); err == nil && !runtime.IsPIDAlive(pid) {
		fmt.Printf("   ✗ stale PID file: PID %d (process dead)\n", pid)
		if forceFix {
			_ = process.RemovePID(svc.Name)
			fmt.Printf("   ✅  FIX: Removed stale PID file\n")
		}
	}

	// 6. Check agent registration
	fmt.Println("   Agent registration:")
	if svc.MCPURL != "" && svc.Command == "" {
		fmt.Printf("   → remote only (mcp_url), auto-registered via mcp_url\n")
	} else if svc.Command != "" || svc.MCPURL != "" {
		if targets.OpenCode {
			ocReg := isRegisteredInFile(opencode.ConfigPath(), svc.Name)
			fmt.Printf("   OpenCode: %s\n", statusIcon(ocReg))
		}
		if targets.Cursor {
			curReg := isRegisteredInFile(cursor.ConfigPath(), svc.Name)
			fmt.Printf("   Cursor: %s\n", statusIcon(curReg))
		}
		if targets.Claude {
			claudeReg := isRegisteredInFile(claudedesktop.ConfigPath(), svc.Name)
			fmt.Printf("   Claude: %s\n", statusIcon(claudeReg))
		}
	} else {
		fmt.Printf("   → no mcp_url or command; not registering\n")
	}

	// 7. Check disk space on volume holding log
	if svc.Log != "" {
		logDir := filepath.Dir(svc.Log)
		fs := filepath.Dir(logDir)
		usage := getDiskUsage(fs)
		freePct := usage.FreePercent
		if freePct < 5 {
			fmt.Printf("   ⚠️  low disk space: %.1f%% free on %s\n", freePct, fs)
			if forceFix {
				fmt.Printf("   ⚠️  FIX: Cannot auto-resolve low disk space. Free up space or change log path.\n")
			}
		} else {
			fmt.Printf("   ✓ disk space: %.1f%% free\n", freePct)
		}
	} else {
		fmt.Println("   No log path configured, skipping disk space check")
	}

	fmt.Printf("\n✅ Doctor complete for %s\n", svc.Name)
}

func getPIDFilePath(name string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mcp-local", "pids", name+".pid")
}

func statusIcon(registered bool) string {
	if registered {
		return "✅"
	}
	return "❌"
}

func getDiskUsage(path string) *diskUsageInfo {
	result := &diskUsageInfo{}
	var cmd string
	if os.Getenv("GOOS") == "windows" {
		// On Windows, just report we can't check
		return result
	} else if os.Getenv("GOOS") == "darwin" {
		cmd = "df -P " + path
	} else {
		cmd = "df -B1 " + path
	}
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return result
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return result
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 5 {
		return result
	}
	result.Total, _ = strconv.ParseUint(fields[1], 10, 64)
	result.Used, _ = strconv.ParseUint(fields[2], 10, 64)
	result.Avail, _ = strconv.ParseUint(fields[3], 10, 64)
	result.UsePct, _ = strconv.ParseUint(fields[4], 10, 64)

	if result.Total > 0 {
		result.FreePercent = 100.0 - float64(result.UsePct)
	}
	return result
}

type diskUsageInfo struct {
	Total       uint64
	Used        uint64
	Avail       uint64
	UsePct      uint64
	FreePercent float64
}

func init() {
	// init only once
}
