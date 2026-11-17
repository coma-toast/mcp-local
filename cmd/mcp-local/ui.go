package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	goruntime "runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/logs"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
	"github.com/coma-toast/mcp-local/internal/mgr/runtime"
	"github.com/coma-toast/mcp-local/internal/portutil"
	"github.com/coma-toast/mcp-local/internal/webui"
	"github.com/spf13/cobra"
)

const (
	uiServiceName = "ui"
	defaultUIHost = "127.0.0.1"
	defaultUIPort = 7831
)

// uiConfigState persists the last-used Web UI bind settings in ~/.mcp-local/ui.json
// so `mcp-local ui start`, `mcp-local start ui`, and `mcp-local ui restart` remember
// custom --host / --port / --allow-remote flags across restarts.
type uiConfigState struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	AllowRemote bool   `json:"allow_remote"`
}

func defaultUIConfigState() uiConfigState {
	return uiConfigState{
		Host:        defaultUIHost,
		Port:        defaultUIPort,
		AllowRemote: false,
	}
}

func uiStateFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mcp-local", "ui.json")
}

func uiLogFilePath() string {
	return process.LogPath(config.ServiceConfig{Name: uiServiceName})
}

func loadUIConfigState() uiConfigState {
	st := defaultUIConfigState()
	data, err := os.ReadFile(uiStateFilePath())
	if err != nil {
		return st
	}
	var parsed uiConfigState
	if json.Unmarshal(data, &parsed) == nil {
		if strings.TrimSpace(parsed.Host) != "" {
			st.Host = strings.TrimSpace(parsed.Host)
		}
		if parsed.Port > 0 {
			st.Port = parsed.Port
		}
		st.AllowRemote = parsed.AllowRemote
	}
	return st
}

func saveUIConfigState(st uiConfigState) {
	path := uiStateFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return
	}
	if data, err := json.MarshalIndent(st, "", "  "); err == nil {
		_ = os.WriteFile(path, append(data, '\n'), 0o600)
	}
}

func resolveUIConfigFromFlags(cmd *cobra.Command, host string, port int, allowRemote bool) uiConfigState {
	st := loadUIConfigState()
	if cmd != nil {
		if cmd.Flags().Changed("host") {
			st.Host = host
		}
		if cmd.Flags().Changed("port") {
			st.Port = port
		}
		if cmd.Flags().Changed("allow-remote") {
			st.AllowRemote = allowRemote
		}
	}
	if strings.TrimSpace(st.Host) == "" {
		st.Host = defaultUIHost
	}
	if st.Port <= 0 {
		st.Port = defaultUIPort
	}
	return st
}

func uiDisplayHost(host string) string {
	h := strings.TrimSpace(host)
	if h == "" || h == "0.0.0.0" || h == "::" {
		return "127.0.0.1"
	}
	return h
}

func uiDisplayURL(st uiConfigState) string {
	return fmt.Sprintf("http://%s", net.JoinHostPort(uiDisplayHost(st.Host), strconv.Itoa(st.Port)))
}

func isUIServiceAlias(name string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "ui", "web", "webui", "web-ui":
		return true
	default:
		return false
	}
}

// isUIServiceTarget returns true when name is a Web UI alias ("ui", "web", "webui")
// and the user has not defined a custom service with that exact name in config.yaml.
func isUIServiceTarget(name string) bool {
	if !isUIServiceAlias(name) {
		return false
	}
	if cfg, err := config.LoadConfig(); err == nil {
		if _, err := cfg.ServiceNamed(name); err == nil {
			return false
		}
	}
	return true
}

func isUIAddrListening(st uiConfigState) bool {
	if st.Port <= 0 {
		return false
	}
	addr := net.JoinHostPort(uiDisplayHost(st.Host), strconv.Itoa(st.Port))
	if conn, err := net.DialTimeout("tcp", addr, 400*time.Millisecond); err == nil {
		_ = conn.Close()
		return true
	}
	return portutil.IsRunning(st.Port)
}

func waitForUIAddr(st uiConfigState, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if isUIAddrListening(st) {
			return true
		}
		time.Sleep(150 * time.Millisecond)
	}
	return isUIAddrListening(st)
}

func uiExpectedExe() string {
	if exe, err := os.Executable(); err == nil && exe != "" {
		return exe
	}
	return "mcp-local"
}

func runningUIPID(st uiConfigState) (int, bool) {
	want := uiExpectedExe()
	if id, err := process.LoadIdentity(uiServiceName); err == nil && id.Check(want) == process.Match {
		return id.PID, true
	}
	if isUIAddrListening(st) {
		if owner := portutil.FindPID(st.Port); owner > 0 {
			if _, exe, err := process.ProcInfo(owner); err == nil && process.SameBinary(exe, want) {
				return owner, true
			}
		}
	}
	return 0, false
}

func uiStopTarget(st uiConfigState) (int, error) {
	want := uiExpectedExe()
	hint := 0
	if id, err := process.LoadIdentity(uiServiceName); err == nil {
		state := id.Check(want)
		if state == process.Match {
			return process.SignalTarget(id.PID), nil
		}
		_ = process.RemovePID(uiServiceName)
		if state != process.Mismatch {
			hint = id.PID
		}
	} else if !os.IsNotExist(err) {
		_ = process.RemovePID(uiServiceName)
	}

	if !isUIAddrListening(st) {
		return 0, nil
	}
	owner := portutil.FindPID(st.Port)
	if owner <= 0 {
		return 0, fmt.Errorf("port %d is in use but its owner can't be found", st.Port)
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
	return 0, fmt.Errorf("port %d is held by pid %d, which mcp-local didn't start; not stopping it", st.Port, owner)
}

func uiPortFree(st uiConfigState) func() error {
	return func() error {
		if isUIAddrListening(st) {
			return fmt.Errorf("port %d still in use (pid %d)", st.Port, portutil.FindPID(st.Port))
		}
		return nil
	}
}

func cleanupOwnUIPID(selfPID int) {
	if id, err := process.LoadIdentity(uiServiceName); err == nil && id.PID == selfPID {
		_ = process.RemovePID(uiServiceName)
	}
}

func runUIForeground(st uiConfigState, openBrowser bool, out io.Writer) error {
	if pid, running := runningUIPID(st); running && pid != os.Getpid() {
		return fmt.Errorf("ui is already running (PID %d at %s); use 'mcp-local ui stop' or 'mcp-local ui restart'", pid, uiDisplayURL(st))
	}

	saveUIConfigState(st)
	_ = process.SavePID(uiServiceName, os.Getpid())
	defer cleanupOwnUIPID(os.Getpid())

	addr := net.JoinHostPort(st.Host, strconv.Itoa(st.Port))
	url := uiDisplayURL(st)
	fmt.Fprintf(out, "🖥️  mcp-local Web UI (React95) running at %s\n", url)
	fmt.Fprintln(out, "   Press Ctrl+C to stop.")

	srv := webui.New(webuiActions(), webui.Options{
		AllowRemote: st.AllowRemote,
	})

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if openBrowser {
		go func() {
			time.Sleep(250 * time.Millisecond)
			_ = openURLInBrowser(url)
		}()
	}

	return srv.ListenAndServe(ctx, addr)
}

func startUIService(st uiConfigState, openBrowser bool, out, errOut io.Writer) error {
	if pid, running := runningUIPID(st); running {
		fmt.Fprintf(out, "  ⏭️  ui: already running (PID %d, %s)\n", pid, uiDisplayURL(st))
		if openBrowser {
			_ = openURLInBrowser(uiDisplayURL(st))
		}
		return nil
	}
	if isUIAddrListening(st) {
		owner := portutil.FindPID(st.Port)
		return fmt.Errorf("port %d is already in use by pid %d", st.Port, owner)
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("determine executable: %w", err)
	}

	saveUIConfigState(st)

	logPath := uiLogFilePath()
	if err := os.MkdirAll(filepath.Dir(logPath), 0o700); err != nil {
		return fmt.Errorf("create log dir: %w", err)
	}
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND|syscall.O_NOFOLLOW, 0o600)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}
	_ = logFile.Chmod(0o600)
	defer logFile.Close()

	args := []string{
		"ui", "run",
		"--host", st.Host,
		"--port", strconv.Itoa(st.Port),
	}
	if st.AllowRemote {
		args = append(args, "--allow-remote")
	}

	cmd := exec.Command(exe, args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ui service: %w", err)
	}

	pid := cmd.Process.Pid
	if err := process.SavePID(uiServiceName, pid); err != nil {
		fmt.Fprintf(errOut, "  ⚠️  ui: save pid: %v\n", err)
	}

	if !waitForUIAddr(st, portWait) {
		printLogTail(errOut, logPath, logTailLines)
		return fmt.Errorf("ui started (PID %d) but %s is not listening yet", pid, net.JoinHostPort(st.Host, strconv.Itoa(st.Port)))
	}

	url := uiDisplayURL(st)
	fmt.Fprintf(out, "  ✅ ui started (PID %d) — %s\n", pid, url)
	fmt.Fprintf(out, "  📄 log: %s\n", logPath)

	if openBrowser {
		_ = openURLInBrowser(url)
	}
	return nil
}

func stopUIService(out, errOut io.Writer) error {
	st := loadUIConfigState()
	target, err := uiStopTarget(st)
	if err != nil {
		return err
	}
	if target == 0 {
		fmt.Fprintln(out, "  ⏭️  ui: not running")
		return nil
	}
	killed, err := process.Terminate(target, uiPortFree(st))
	if err != nil {
		return err
	}
	_ = process.RemovePID(uiServiceName)
	pid, verb := target, "stopped"
	if pid < 0 {
		pid = -pid
	}
	if killed {
		verb = "killed"
	}
	fmt.Fprintf(out, "  ✅ ui %s (pid %d)\n", verb, pid)
	_ = errOut
	return nil
}

func restartUIService(st uiConfigState, openBrowser bool, out, errOut io.Writer) error {
	if err := stopUIService(out, errOut); err != nil {
		return err
	}
	time.Sleep(restartDelay)
	return startUIService(st, openBrowser, out, errOut)
}

func printUIStatus(out io.Writer) error {
	st := loadUIConfigState()
	url := uiDisplayURL(st)
	logPath := uiLogFilePath()
	if pid, running := runningUIPID(st); running {
		uptime := portutil.GetUptime(pid)
		fmt.Fprintf(out, "🖥️  ui: running (pid %d, uptime %s) — %s\n", pid, uptime, url)
		fmt.Fprintf(out, "   log: %s\n", logPath)
		return nil
	}
	fmt.Fprintf(out, "⏹️  ui: stopped (configured for %s, log: %s)\n", url, logPath)
	return nil
}

func tailUILog(out, errOut io.Writer) error {
	logPath := uiLogFilePath()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	fmt.Fprintf(out, "==> %s <==\n", logPath)
	return logs.Follower{Path: logPath, Out: out, Err: errOut}.Run(ctx)
}

func addUIFlags(cmd *cobra.Command, host *string, port *int, openBrowser, allowRemote *bool) {
	cmd.Flags().StringVar(host, "host", defaultUIHost, "bind address (remembered in ~/.mcp-local/ui.json)")
	cmd.Flags().IntVar(port, "port", defaultUIPort, "HTTP listen port (remembered in ~/.mcp-local/ui.json)")
	cmd.Flags().BoolVar(openBrowser, "open", false, "open Web UI in the default browser")
	cmd.Flags().BoolVar(allowRemote, "allow-remote", false, "disable localhost Host/Origin checks (for Tailscale/LAN)")
}

func addUI() {
	var (
		host        string
		port        int
		openBrowser bool
		allowRemote bool
		daemon      bool
	)
	cmd := &cobra.Command{
		Use:     "ui",
		Aliases: []string{"web"},
		Short:   "Run or manage the React95 Web UI (foreground or background service)",
		Long: `Launch the React95 Web UI in the foreground (stops on Ctrl+C / SIGTERM), or
manage it as a background daemon service with start, stop, restart, status, and log.`,
		Example: `  mcp-local ui                                            # Run in foreground (stops when killed)
  mcp-local ui --host 0.0.0.0 --port 7831 --allow-remote
  mcp-local ui start                                      # Start as background service
  mcp-local ui stop                                       # Stop background service
  mcp-local ui restart                                    # Restart background service
  mcp-local ui status                                     # Check background service status
  mcp-local start ui                                      # Also works via top-level start/stop/restart`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			st := resolveUIConfigFromFlags(cmd, host, port, allowRemote)
			if daemon {
				fmt.Fprintln(cmd.OutOrStdout(), "🚀 Starting Web UI service...")
				return startUIService(st, openBrowser, cmd.OutOrStdout(), cmd.ErrOrStderr())
			}
			return runUIForeground(st, openBrowser, cmd.OutOrStdout())
		},
	}
	addUIFlags(cmd, &host, &port, &openBrowser, &allowRemote)
	cmd.Flags().BoolVarP(&daemon, "daemon", "d", false, "run Web UI in the background as a daemon service")

	// `mcp-local ui run` (explicit foreground subcommand)
	var runHost string
	var runPort int
	var runOpen, runAllowRemote bool
	runSub := &cobra.Command{
		Use:     "run",
		Aliases: []string{"serve", "fg"},
		Short:   "Run the Web UI server in the foreground (stops when killed)",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			st := resolveUIConfigFromFlags(cmd, runHost, runPort, runAllowRemote)
			return runUIForeground(st, runOpen, cmd.OutOrStdout())
		},
	}
	addUIFlags(runSub, &runHost, &runPort, &runOpen, &runAllowRemote)

	// `mcp-local ui start`
	var startHost string
	var startPort int
	var startOpen, startAllowRemote bool
	startSub := &cobra.Command{
		Use:   "start",
		Short: "Start the Web UI as a background daemon service",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			st := resolveUIConfigFromFlags(cmd, startHost, startPort, startAllowRemote)
			fmt.Fprintln(cmd.OutOrStdout(), "🚀 Starting Web UI service...")
			return startUIService(st, startOpen, cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}
	addUIFlags(startSub, &startHost, &startPort, &startOpen, &startAllowRemote)

	// `mcp-local ui stop`
	stopSub := &cobra.Command{
		Use:   "stop",
		Short: "Stop the background Web UI service",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "⏹️  Stopping Web UI service...")
			return stopUIService(cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}

	// `mcp-local ui restart`
	var restartHost string
	var restartPort int
	var restartOpen, restartAllowRemote bool
	restartSub := &cobra.Command{
		Use:   "restart",
		Short: "Restart the background Web UI service",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			st := resolveUIConfigFromFlags(cmd, restartHost, restartPort, restartAllowRemote)
			fmt.Fprintln(cmd.OutOrStdout(), "🔁 Restarting Web UI service...")
			return restartUIService(st, restartOpen, cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}
	addUIFlags(restartSub, &restartHost, &restartPort, &restartOpen, &restartAllowRemote)

	// `mcp-local ui status`
	statusSub := &cobra.Command{
		Use:   "status",
		Short: "Show the status of the Web UI service",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return printUIStatus(cmd.OutOrStdout())
		},
	}

	// `mcp-local ui log`
	logSub := &cobra.Command{
		Use:     "log",
		Aliases: []string{"logs"},
		Short:   "Tail the Web UI service log file",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return tailUILog(cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}

	// `mcp-local ui open`
	openSub := &cobra.Command{
		Use:   "open",
		Short: "Open the Web UI in the default browser",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			st := loadUIConfigState()
			return openURLInBrowser(uiDisplayURL(st))
		},
	}

	cmd.AddCommand(runSub, startSub, stopSub, restartSub, statusSub, logSub, openSub)
	rootCmd.AddCommand(cmd)
}

func openURLInBrowser(url string) error {
	opener := "xdg-open"
	if goruntime.GOOS == "darwin" {
		opener = "open"
	}
	return exec.Command(opener, url).Start()
}

func webuiActions() webui.Actions {
	return webui.Actions{
		Version:          version,
		StartServices:    startServicesTo,
		StopService:      stopServiceNamedTo,
		RestartServices:  restartServicesTo,
		RegisterStarted:  registerStartedTo,
		FetchToolsList:   fetchToolsList,
		MergeSyncedTools: mergeSyncedTools,
		FindCommand:      findCommand,
		ServiceFromEntry: serviceFromEntry,
		DoctorService: func(cfg *config.ManagerConfig, svc config.ServiceConfig, fix bool, out, errOut io.Writer) webui.DoctorServiceResult {
			st := runtime.CheckStatus(svc)
			fmt.Fprintf(out, "\n🔍 %s: %s", svc.Name, st.State)
			if st.Detail != "" && st.Detail != string(st.State) {
				fmt.Fprintf(out, " (%s)", st.Detail)
			}
			fmt.Fprintln(out)
			found := collectDoctorFindings(cfg, svc, st, out, errOut)
			res := webui.DoctorServiceResult{
				Name:     svc.Name,
				State:    string(st.State),
				Detail:   st.Detail,
				Findings: make([]webui.DoctorFindingDTO, 0, len(found)),
			}
			for _, f := range found {
				dto := webui.DoctorFindingDTO{
					Message: f.what,
					Fixable: f.fix != nil,
				}
				switch {
				case !fix:
					hint := "no automatic fix"
					if f.fix != nil {
						hint = "--fix resolves this"
					}
					fmt.Fprintf(out, "  ✗ %s (%s)\n", f.what, hint)
					res.Remaining++
				case f.fix == nil:
					fmt.Fprintf(out, "  ✗ %s (no automatic fix)\n", f.what)
					res.Remaining++
				default:
					if err := f.fix(); err != nil {
						fmt.Fprintf(out, "  ✗ %s: fix failed: %v\n", f.what, err)
						dto.Error = err.Error()
						res.Remaining++
					} else {
						fmt.Fprintf(out, "  ✅ fixed: %s\n", f.what)
						dto.Fixed = true
					}
				}
				res.Findings = append(res.Findings, dto)
			}
			if len(found) == 0 {
				fmt.Fprintln(out, "  ✓ no problems")
			}
			return res
		},
	}
}
