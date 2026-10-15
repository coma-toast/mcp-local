package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"github.com/coma-toast/mcp-local/internal/mgr/agents"
	"github.com/coma-toast/mcp-local/internal/mgr/claudedesktop"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/cursor"
	"github.com/coma-toast/mcp-local/internal/mgr/opencode"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
	"github.com/coma-toast/mcp-local/internal/mgr/runtime"
	"github.com/spf13/cobra"
)

// lowDiskPercent is the free-space level below which doctor warns about a log volume.
var lowDiskPercent = 5.0

func addDoctor() {
	var fix bool
	doctor := &cobra.Command{
		Use:   "doctor [service]",
		Short: "Diagnose stuck services and optionally apply safe fixes",
		Long: `Checks each service for untracked copies holding its ports, other programs on its
ports, failing health checks, stale PID files, missing dependencies, missing agent
registrations, and low disk space on its log volume.

Without --fix nothing is changed. --fix only stops untracked copies of the service's own
binary, removes stale PID files, restarts degraded services and re-registers. It never
stops another program. Exits non-zero while any problem remains.`,
		Example: `  mcp-local doctor
  mcp-local doctor my-server --fix`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			names := cfg.SortedNames()
			if len(args) == 1 {
				if _, err := cfg.ServiceNamed(args[0]); err != nil {
					return err
				}
				names = args
			}
			remaining := 0
			for _, name := range names {
				svc, _ := cfg.ServiceNamed(name)
				remaining += doctorService(cfg, svc, fix, os.Stdout)
			}
			if remaining > 0 {
				cmd.SilenceUsage = true
				return fmt.Errorf("%d problem(s) remain", remaining)
			}
			fmt.Println("\n✅ No problems found")
			return nil
		},
	}
	doctor.Flags().BoolVarP(&fix, "fix", "f", false, "apply safe fixes")
	rootCmd.AddCommand(doctor)
}

// finding is one problem doctor found, with the fix --fix would apply (nil when there is
// no safe automatic fix).
type finding struct {
	what string
	fix  func() error
}

// doctorService prints svc's report and returns how many problems remain.
func doctorService(cfg *config.ManagerConfig, svc config.ServiceConfig, fix bool, w io.Writer) int {
	st := runtime.CheckStatus(svc)
	fmt.Fprintf(w, "\n🔍 %s: %s", svc.Name, st.State)
	if st.Detail != "" && st.Detail != string(st.State) {
		fmt.Fprintf(w, " (%s)", st.Detail)
	}
	fmt.Fprintln(w)
	var found []finding
	restart := func() error {
		if _, failed := restartServices(cfg, []string{svc.Name}, os.Stderr); failed > 0 {
			return fmt.Errorf("restart failed")
		}
		return nil
	}
	switch {
	case st.Orphan.PID > 0:
		found = append(found, finding{st.Detail, func() error {
			_, err := process.Terminate(process.SignalTarget(st.Orphan.PID), portsFree(svc.OwnedPorts()))
			return err
		}})
	case st.State == config.StateConflict:
		found = append(found, finding{st.Detail + "; stop that program or change the service's ports", nil})
	case st.State == config.StateDegraded:
		found = append(found, finding{st.Detail, restart})
	}
	if st.Stale {
		found = append(found, finding{fmt.Sprintf("stale PID file %s", process.GetPIDFile(svc.Name)), func() error { return process.RemovePID(svc.Name) }})
	}
	for _, dep := range svc.Deps {
		if _, err := os.Stat(config.ExpandPath(dep)); err != nil {
			found = append(found, finding{fmt.Sprintf("dependency missing: %s", dep), nil})
		}
	}
	if missing := missingRegistrations(cfg, svc); len(missing) > 0 {
		found = append(found, finding{fmt.Sprintf("not registered with %v", missing), func() error {
			return agents.RegisterAll([]config.ServiceConfig{svc}, agents.TargetsFromConfig(*cfg))
		}})
	}
	if process.ModeOf(svc) == process.Daemon || process.ModeOf(svc) == process.Embedded {
		logPath := process.LogPath(svc)
		if free, ok := freePercent(logPath); ok && free < lowDiskPercent {
			found = append(found, finding{fmt.Sprintf("only %.1f%% free on the volume holding %s", free, logPath), nil})
		}
	}
	remaining := 0
	for _, f := range found {
		switch {
		case !fix:
			hint := "no automatic fix"
			if f.fix != nil {
				hint = "--fix resolves this"
			}
			fmt.Fprintf(w, "  ✗ %s (%s)\n", f.what, hint)
			remaining++
		case f.fix == nil:
			fmt.Fprintf(w, "  ✗ %s (no automatic fix)\n", f.what)
			remaining++
		default:
			if err := f.fix(); err != nil {
				fmt.Fprintf(w, "  ✗ %s: fix failed: %v\n", f.what, err)
				remaining++
				continue
			}
			fmt.Fprintf(w, "  ✅ fixed: %s\n", f.what)
		}
	}
	if len(found) == 0 {
		fmt.Fprintln(w, "  ✓ no problems")
	}
	return remaining
}

// missingRegistrations lists the enabled agent hosts whose config lacks svc. Services
// mcp-local can't register (no command and no mcp_url) are skipped.
func missingRegistrations(cfg *config.ManagerConfig, svc config.ServiceConfig) []string {
	if svc.Command == "" && svc.MCPURL == "" {
		return nil
	}
	t := agents.TargetsFromConfig(*cfg)
	var missing []string
	for _, h := range []struct {
		on   bool
		name string
		path string
	}{
		{t.OpenCode, "OpenCode", opencode.ConfigPath()},
		{t.Cursor, "Cursor", cursor.ConfigPath()},
		{t.Claude, "Claude Desktop", claudedesktop.ConfigPath()},
	} {
		if h.on && !isRegisteredInFile(h.path, svc.Name) {
			missing = append(missing, h.name)
		}
	}
	return missing
}

// freePercent reports the free space on the volume holding path, using its nearest
// existing ancestor directory.
func freePercent(path string) (float64, bool) {
	dir := filepath.Dir(path)
	for {
		var fs syscall.Statfs_t
		if err := syscall.Statfs(dir, &fs); err == nil {
			if fs.Blocks == 0 {
				return 0, false
			}
			return 100 * float64(fs.Bavail) / float64(fs.Blocks), true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return 0, false
		}
		dir = parent
	}
}
