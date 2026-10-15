package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/coma-toast/mcp-local/internal/mgr/agents"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
	"github.com/coma-toast/mcp-local/internal/mgr/runtime"
	"github.com/coma-toast/mcp-local/internal/portutil"
	"github.com/spf13/cobra"
)

const (
	portWait     = 5 * time.Second
	restartDelay = 300 * time.Millisecond
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

// failures turns a failure count into the command's error, so the exit status is non-zero.
func failures(cmd *cobra.Command, n int, verb string) error {
	if n == 0 {
		return nil
	}
	cmd.SilenceUsage = true
	return fmt.Errorf("%d service(s) failed to %s", n, verb)
}

func warnDeregister(name string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠️  %s: deregister failed:\n%s\n", name, indent(err.Error()))
	}
}

func indent(s string) string {
	return "    " + strings.ReplaceAll(s, "\n", "\n    ")
}

// stopTarget finds what to signal for a service: its PID file entry when that still names
// the process mcp-local started, else the port owner when it runs the service's binary or
// belongs to the process group the PID file recorded. 0 means nothing is running; a
// negative target is a process group. Stale PID files are removed.
func stopTarget(name string, svc config.ServiceConfig) (int, error) {
	want := process.ExpectedExe(svc)
	hint := 0
	if id, err := process.LoadIdentity(name); err == nil {
		st := id.Check(want)
		if st == process.Match {
			return process.SignalTarget(id.PID), nil
		}
		_ = process.RemovePID(name)
		if st == process.Unverified {
			fmt.Fprintf(os.Stderr, "  ⚠️  %s: pid %d from %s can't be verified; ignoring it\n", name, id.PID, process.GetPIDFile(name))
		}
		// A reused PID may lead an unrelated group, so only a dead or unverified one is a hint.
		if st != process.Mismatch {
			hint = id.PID
		}
	} else if !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "  ⚠️  %s: %v\n", name, err)
		_ = process.RemovePID(name)
	}
	if svc.Port <= 0 || !portutil.IsRunning(svc.Port) {
		if st := runtime.CheckStatus(svc); st.Orphan.PID > 0 {
			return process.SignalTarget(st.Orphan.PID), nil
		}
		return 0, nil
	}
	owner := portutil.FindPID(svc.Port)
	if owner <= 0 {
		return 0, fmt.Errorf("port %d is in use but its owner can't be found", svc.Port)
	}
	if hint > 0 && owner == hint {
		return process.SignalTarget(owner), nil
	}
	if hint > 1 && process.GroupOf(owner) == hint {
		return -hint, nil
	}
	if _, exe, err := process.ProcInfo(owner); err == nil && process.SameBinary(exe, want) {
		return process.SignalTarget(owner), nil
	}
	return 0, fmt.Errorf("port %d is held by pid %d, which mcp-local didn't start; not stopping it", svc.Port, owner)
}

func portsFree(ports []int) func() error {
	return func() error {
		for _, port := range ports {
			if port > 0 && portutil.IsRunning(port) {
				return fmt.Errorf("port %d still in use (pid %d)", port, portutil.FindPID(port))
			}
		}
		return nil
	}
}

// clearOrphan stops an untracked copy of the service's own binary that holds one of its
// ports while the main port is down (a process left behind by an earlier run), and
// refuses to start when another program holds one of the ports. It never touches a
// process that isn't the service's binary.
func clearOrphan(name string, svc config.ServiceConfig) error {
	st := runtime.CheckStatus(svc)
	switch {
	case st.State == config.StateConflict:
		return fmt.Errorf("port conflict: %s", st.Detail)
	case st.Orphan.PID <= 0:
		return nil
	}
	fmt.Printf("  ⚠️  %s: stopping %s\n", name, st.Detail)
	if _, err := process.Terminate(process.SignalTarget(st.Orphan.PID), portsFree(svc.OwnedPorts())); err != nil {
		return fmt.Errorf("could not stop the untracked copy: %w", err)
	}
	return nil
}

const logTailLines = 10

// printLogTail prints the last n lines of the service log, so a failed start shows why.
func printLogTail(w io.Writer, path string, n int) {
	lines, err := tailLines(path, n)
	if err != nil || len(lines) == 0 {
		return
	}
	fmt.Fprintf(w, "  📄 last %d lines of %s:\n", len(lines), path)
	for _, l := range lines {
		fmt.Fprintf(w, "    %s\n", l)
	}
}

// tailLines returns up to n trailing lines, reading at most the last 64 KiB of the file.
func tailLines(path string, n int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	off := fi.Size() - 64<<10
	if off < 0 {
		off = 0
	}
	buf := make([]byte, fi.Size()-off)
	if _, err := f.ReadAt(buf, off); err != nil && err != io.EOF {
		return nil, err
	}
	lines := strings.Split(strings.TrimRight(string(buf), "\n"), "\n")
	if off > 0 && len(lines) > 1 {
		lines = lines[1:] // the first line is probably cut off
	}
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return lines, nil
}

// stopServiceNamed stops one service and prints its ✅/⏭️ line. It returns an error, for
// the caller to print, unless the process is gone and its port is free.
func stopServiceNamed(cfg *config.ManagerConfig, name string, deregister bool) error {
	svc, err := cfg.ServiceNamed(name)
	if err != nil {
		return err
	}
	target, err := stopTarget(name, svc)
	if err != nil {
		return err
	}
	mode := process.ModeOf(svc)
	switch {
	case target != 0:
		killed, err := process.Terminate(target, portsFree(svc.OwnedPorts()))
		if err != nil {
			return err
		}
		_ = process.RemovePID(name)
		pid, verb := target, "stopped"
		if pid < 0 {
			pid = -pid
		}
		if killed {
			verb = "killed"
		}
		fmt.Printf("  ✅ %s %s (pid %d)\n", name, verb, pid)
	case mode == process.Remote || mode == process.Stdio:
		fmt.Printf("  ⏭️  %s: %s service, no process to stop\n", name, mode)
	default:
		fmt.Printf("  ⏭️  %s: not running\n", name)
		return nil
	}
	if deregister {
		warnDeregister(name, agents.DeregisterAll(svc, agents.TargetsFromConfig(*cfg)))
	}
	return nil
}

// runningPID reports whether a copy of svc is already up: its port is listening or its
// PID file still names the process mcp-local started.
func runningPID(name string, svc config.ServiceConfig) (int, bool) {
	if svc.Port > 0 && portutil.IsRunning(svc.Port) {
		return portutil.FindPID(svc.Port), true
	}
	if id, err := process.LoadIdentity(name); err == nil && id.Check(process.ExpectedExe(svc)) == process.Match {
		return id.PID, true
	}
	return 0, false
}

// startOne starts a service and prints its status line. Remote and stdio services are
// only registered: there is nothing to run, or the agents launch them on demand.
func startOne(name string, svc config.ServiceConfig, build bool, errOut io.Writer) error {
	mode := process.ModeOf(svc)
	if mode != process.Embedded && strings.TrimSpace(svc.Command) == "" && strings.TrimSpace(svc.MCPURL) == "" {
		return fmt.Errorf("no command or mcp_url configured")
	}
	if mode == process.Remote {
		if svc.Port > 0 && !portutil.IsRunning(svc.Port) {
			fmt.Printf("  ⚠️  %s: remote; nothing listening on port %d (start the external server), registering only\n", name, svc.Port)
			return nil
		}
		fmt.Printf("  🌐 %s: remote (%s), registering only\n", name, svc.MCPURL)
		return nil
	}
	if mode != process.Stdio {
		if err := clearOrphan(name, svc); err != nil {
			return err
		}
		if pid, ok := runningPID(name, svc); ok {
			if svc.Port > 0 && !portutil.IsRunning(svc.Port) {
				fmt.Printf("  ⚠️  %s: already running (pid %d) but port %d isn't listening yet; `mcp-local restart %s` replaces it\n", name, pid, svc.Port, name)
				return nil
			}
			where := fmt.Sprintf("pid %d", pid)
			if svc.Port > 0 {
				where = fmt.Sprintf("port %d", svc.Port)
			}
			fmt.Printf("  ⏭️  %s: already running (%s)\n", name, where)
			return nil
		}
	}
	if build && config.ShouldBuildOnStart(svc) {
		if err := process.Build(svc); err != nil {
			return err
		}
	}
	if mode == process.Stdio {
		fmt.Printf("  📎 %s: stdio, launched by agents on demand; registering only\n", name)
		return nil
	}
	start, how := process.StartService, ""
	if mode == process.Embedded {
		start, how = process.StartEmbedded, "in-process, "
	}
	p, err := start(svc)
	if err != nil {
		printLogTail(errOut, process.LogPath(svc), logTailLines)
		return err
	}
	if err := process.SavePID(name, p.PID); err != nil {
		fmt.Fprintf(errOut, "  ⚠️  %s: save pid: %v\n", name, err)
	}
	fmt.Printf("  ✅ %s started (%sPID %d)\n", name, how, p.PID)
	if svc.Port > 0 && !portutil.WaitForPort(svc.Port, portWait) {
		fmt.Fprintf(errOut, "  ⚠️  %s: port %d not listening yet\n", name, svc.Port)
		printLogTail(errOut, process.LogPath(svc), logTailLines)
	}
	return nil
}

// startServices starts names, returning the services to register and how many failed.
func startServices(cfg *config.ManagerConfig, names []string, build bool, errOut io.Writer) ([]config.ServiceConfig, int) {
	var started []config.ServiceConfig
	failed := 0
	fmt.Println("🚀 Starting MCP servers...")
	for _, name := range names {
		svc, err := cfg.ServiceNamed(name)
		if err == nil {
			err = startOne(name, svc, build, errOut)
		}
		if err != nil {
			fmt.Printf("  ❌ %s: %v\n", name, err)
			failed++
			continue
		}
		started = append(started, svc)
	}
	return started, failed
}

// restartServices builds each service first, so a failed build leaves the running copy
// alone, then stops it and starts it again without a second build.
func restartServices(cfg *config.ManagerConfig, names []string, errOut io.Writer) ([]config.ServiceConfig, int) {
	failed := 0
	var ready []string
	fmt.Println("🔁 Restarting MCP servers...")
	for _, name := range names {
		svc, err := cfg.ServiceNamed(name)
		if err == nil && config.ShouldBuildOnStart(svc) && process.ModeOf(svc) != process.Remote {
			if err = process.Build(svc); err != nil {
				err = fmt.Errorf("%w; left as it was", err)
			}
		}
		if err == nil {
			err = stopServiceNamed(cfg, name, false)
		}
		if err != nil {
			fmt.Printf("  ❌ %s: %v\n", name, err)
			failed++
			continue
		}
		ready = append(ready, name)
	}
	if len(ready) == 0 {
		return nil, failed
	}
	time.Sleep(restartDelay)
	started, n := startServices(cfg, ready, false, errOut)
	return started, failed + n
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
			started, failed := startServices(cfg, names, true, cmd.ErrOrStderr())
			registerStarted(cfg, started)
			fmt.Println("\nDone!")
			return failures(cmd, failed, "start")
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
			failed := 0
			for _, name := range names {
				if err := stopServiceNamed(cfg, name, stopDeregister); err != nil {
					fmt.Printf("  ❌ %s: %v\n", name, err)
					failed++
				}
			}
			fmt.Println("\nDone!")
			return failures(cmd, failed, "stop")
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
			started, failed := restartServices(cfg, names, cmd.ErrOrStderr())
			registerStarted(cfg, started)
			fmt.Println("\nDone!")
			return failures(cmd, failed, "restart")
		},
	}
	restart.Flags().BoolVar(&restartAll, "all", false, "restart every configured service")

	rootCmd.AddCommand(start, stop, restart)
}
