package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/coma-toast/mcp-local/internal/bridge"
	"github.com/coma-toast/mcp-local/internal/mgr/agents"
	"github.com/coma-toast/mcp-local/internal/mgr/asttools"
	"github.com/coma-toast/mcp-local/internal/mgr/claudedesktop"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/cursor"
	"github.com/coma-toast/mcp-local/internal/mgr/embedfs"
	"github.com/coma-toast/mcp-local/internal/mgr/jsonagent"
	"github.com/coma-toast/mcp-local/internal/mgr/logs"
	"github.com/coma-toast/mcp-local/internal/mgr/opencode"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
	"github.com/coma-toast/mcp-local/internal/portutil"
	"github.com/coma-toast/mcp-local/internal/statusui"
	"github.com/coma-toast/mcp-local/internal/tui"
	"github.com/spf13/cobra"
)

func addCommands() {
	addLifecycle()
	addDoctor()
	addUpdate()

	rootCmd.AddCommand(cmdList())

	var plainStatus bool
	status := &cobra.Command{
		Use:   "status [service]",
		Short: "Live status table (q to quit); use --plain for one-shot",
		Example: `  mcp-local status
  mcp-local status --plain
  mcp-local status my-server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			filter := ""
			if len(args) > 0 {
				filter = args[0]
				if _, err := cfg.ServiceNamed(filter); err != nil {
					return err
				}
			}
			if plainStatus {
				statusui.PrintPlain(cfg, filter)
				return nil
			}
			return statusui.Run(cfg, filter)
		},
	}
	status.Flags().BoolVar(&plainStatus, "plain", false, "print once and exit (no TUI)")
	rootCmd.AddCommand(status)

	health := &cobra.Command{
		Use:   "health [service]",
		Short: "GET health_url for service(s)",
		Example: `  mcp-local health
  mcp-local health my-server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			names := cfg.SortedNames()
			if len(args) > 0 {
				if _, err := cfg.ServiceNamed(args[0]); err != nil {
					return err
				}
				names = []string{args[0]}
			}
			client := &http.Client{Timeout: 3 * time.Second}
			return failures(cmd, runHealth(cfg, names, client, cmd.OutOrStdout(), cmd.ErrOrStderr()), "pass health checks")
		},
	}
	rootCmd.AddCommand(health)

	rootCmd.AddCommand(cmdLog())
	rootCmd.AddCommand(cmdOpen())
	rootCmd.AddCommand(cmdRebuild())
	rootCmd.AddCommand(cmdAdd())
	rootCmd.AddCommand(cmdRemove())
	rootCmd.AddCommand(cmdEdit())
	rootCmd.AddCommand(cmdConfigEdit())
	rootCmd.AddCommand(cmdRegister())
	rootCmd.AddCommand(cmdDeregister())
	rootCmd.AddCommand(cmdRegistered())
	rootCmd.AddCommand(cmdValidate())
	rootCmd.AddCommand(cmdImport())

	rootCmd.AddCommand(cmdLogsUnified())
	rootCmd.AddCommand(cmdTools())
	rootCmd.AddCommand(cmdJSON())
	rootCmd.AddCommand(cmdServeFS())
	rootCmd.AddCommand(cmdBridge())
}

// runHealth GETs each service's health_url and returns how many failed. Services
// without a health_url are skipped, not failed.
func runHealth(cfg *config.ManagerConfig, names []string, client *http.Client, out, errOut io.Writer) int {
	failed := 0
	for _, name := range names {
		svc, _ := cfg.ServiceNamed(name)
		if svc.HealthURL == "" {
			fmt.Fprintf(out, "%-24s  no health_url configured\n", name)
			continue
		}
		resp, err := client.Get(svc.HealthURL)
		if err != nil {
			fmt.Fprintf(errOut, "%-24s  ✗ ERROR: %v\n", name, err)
			failed++
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		var pretty string
		var obj interface{}
		if json.Unmarshal(body, &obj) == nil {
			b, _ := json.MarshalIndent(obj, "                          ", "  ")
			pretty = string(b)
		} else {
			pretty = strings.TrimSpace(string(body))
		}
		status := fmt.Sprintf("HTTP %d", resp.StatusCode)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			fmt.Fprintf(out, "%-24s  ✓ %s  %s\n", name, status, pretty)
		} else {
			fmt.Fprintf(errOut, "%-24s  ✗ %s  %s\n", name, status, pretty)
			failed++
		}
	}
	return failed
}

func cmdList() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List configured services",
		Example: `  mcp-local list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			if len(cfg.Services) == 0 {
				fmt.Println("No services configured. Use: mcp-local add <name> or the wizard.")
				return nil
			}
			for _, name := range cfg.SortedNames() {
				fmt.Println(name)
			}
			return nil
		},
	}
}

func cmdLog() *cobra.Command {
	return &cobra.Command{
		Use:     "log <service>",
		Short:   "Tail the service log file",
		Args:    cobra.ExactArgs(1),
		Example: `  mcp-local log my-server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			svc, err := cfg.ServiceNamed(args[0])
			if err != nil {
				return err
			}
			logPath := process.LogPath(svc)
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			fmt.Printf("==> %s <==\n", logPath)
			return logs.Follower{Path: logPath, Out: os.Stdout, Err: os.Stderr}.Run(ctx)
		},
	}
}

func cmdOpen() *cobra.Command {
	return &cobra.Command{
		Use:     "open <service>",
		Short:   "Open dashboard_url in a browser",
		Args:    cobra.ExactArgs(1),
		Example: `  mcp-local open my-server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			svc, err := cfg.ServiceNamed(args[0])
			if err != nil {
				return err
			}
			if svc.DashboardURL == "" {
				return fmt.Errorf("no dashboard_url configured for %s", args[0])
			}
			opener := "xdg-open"
			if runtime.GOOS == "darwin" {
				opener = "open"
			}
			return exec.Command(opener, svc.DashboardURL).Start()
		},
	}
}

func cmdRebuild() *cobra.Command {
	return &cobra.Command{
		Use:     "rebuild <service>",
		Short:   "Run build_cmd / build_command for a service",
		Args:    cobra.ExactArgs(1),
		Example: `  mcp-local rebuild my-server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			svc, err := cfg.ServiceNamed(args[0])
			if err != nil {
				return err
			}
			b := config.EffectiveBuild(svc)
			if b == "" {
				return fmt.Errorf("no build_cmd configured for %s", args[0])
			}
			fmt.Printf("Rebuilding %s...\n", args[0])
			c := exec.Command("sh", "-c", b)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				return fmt.Errorf("rebuild failed: %w", err)
			}
			fmt.Printf("%s: rebuild complete\n", args[0])
			targets := agents.TargetsFromConfig(*cfg)
			if err := agents.RegisterAll([]config.ServiceConfig{svc}, targets); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "  ⚠️  Registration failed:\n%s\n", indent(err.Error()))
			} else {
				for _, host := range targets.Hosts() {
					fmt.Printf("  ✅ %s: registered %s\n", host, args[0])
				}
			}
			return nil
		},
	}
}

func cmdAdd() *cobra.Command {
	var (
		command      string
		port         int
		svcType      string
		health       string
		dashboard    string
		mcpURL       string
		logPath      string
		buildCommand string
		envPairs     []string
		cmdArgs      []string
		path         string
	)
	c := &cobra.Command{
		Use:   "add <name>",
		Short: "Add or update a service",
		Args:  cobra.ExactArgs(1),
		Example: `  mcp-local add github --type stdio --command npx --args "-y" --args "@modelcontextprotocol/server-github" --env GITHUB_TOKEN=ghp_xxxx
  mcp-local add filesystem --type filesystem --path ~/projects --port 4000
  mcp-local add playwright --type stdio --command npx --args "-y" --args "@playwright/mcp"
  mcp-local add fetch --type stdio --command uvx --args "mcp-server-fetch"
  mcp-local add remote-api --type http --mcp-url https://api.example.com/mcp --health https://api.example.com/health`,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cfg := mustConfig()
			existing, _ := cfg.ServiceNamed(name)
			isNew := existing.Name == ""
			var svc config.ServiceConfig
			if isNew {
				svc.Name = name
			} else {
				svc = existing
			}
			if cmd.Flags().Changed("command") {
				svc.Command = command
			}
			if cmd.Flags().Changed("port") {
				svc.Port = port
			}
			if cmd.Flags().Changed("type") {
				svc.MCPType = svcType
			}
			if cmd.Flags().Changed("health") {
				svc.HealthURL = health
			}
			if cmd.Flags().Changed("dashboard") {
				svc.DashboardURL = dashboard
			}
			if cmd.Flags().Changed("mcp-url") {
				svc.MCPURL = mcpURL
			}
			if cmd.Flags().Changed("log") {
				svc.Log = logPath
			}
			if cmd.Flags().Changed("build-command") {
				svc.BuildCommand = buildCommand
			}
			if cmd.Flags().Changed("path") {
				svc.Path = path
			}
			if cmd.Flags().Changed("args") {
				svc.Args = cmdArgs
			}
			if cmd.Flags().Changed("env") {
				if svc.Env == nil {
					svc.Env = map[string]string{}
				}
				for _, pair := range envPairs {
					parts := strings.SplitN(pair, "=", 2)
					if len(parts) == 2 {
						svc.Env[parts[0]] = parts[1]
					}
				}
			}
			if isNew {
				if svc.MCPType == "filesystem" {
					if svc.Port == 0 {
						portStr := prompt("Port: ")
						svc.Port, _ = strconv.Atoi(strings.TrimSpace(portStr))
					}
				} else if svc.Command == "" && svc.MCPURL == "" {
					svc.Command = prompt("Command (path to binary, or empty for remote-only): ")
				}
				if svc.MCPType == "" {
					svc.MCPType = prompt("Type (http, stdio, filesystem, or empty to auto-detect): ")
				}
				if svc.Command == "" && svc.Port == 0 && svc.MCPURL == "" {
					mcpURL := prompt("MCP URL (for remote-only service): ")
					svc.MCPURL = strings.TrimSpace(mcpURL)
				}
				if svc.Port == 0 && svc.Command != "" {
					portStr := prompt("Port (0 for stdio): ")
					svc.Port, _ = strconv.Atoi(strings.TrimSpace(portStr))
				}
				if svc.HealthURL == "" && svc.Port > 0 {
					suggested := fmt.Sprintf("http://localhost:%d/health", svc.Port)
					val := prompt(fmt.Sprintf("Health URL [%s]: ", suggested))
					if strings.TrimSpace(val) == "" {
						svc.HealthURL = suggested
					} else {
						svc.HealthURL = strings.TrimSpace(val)
					}
				}
				if svc.Log == "" && svc.Command != "" {
					suggested := process.LogPath(config.ServiceConfig{Name: name})
					// Empty keeps the default, resolved per machine, out of the (often shared) config.
					svc.Log = strings.TrimSpace(prompt(fmt.Sprintf("Log path [%s]: ", suggested)))
				}
			}
			if svc.MCPType == "filesystem" && strings.TrimSpace(svc.Path) == "" {
				return fmt.Errorf("when type=filesystem, path must be provided")
			}
			cfg.UpsertService(svc)
			if err := config.SaveConfig(cfg); err != nil {
				return err
			}
			if isNew {
				fmt.Printf("Added service %q\n", name)
			} else {
				fmt.Printf("Updated service %q\n", name)
			}
			return nil
		},
	}
	c.Flags().StringVar(&command, "command", "", "Path to binary")
	c.Flags().IntVar(&port, "port", 0, "Listen port")
	c.Flags().StringVar(&svcType, "type", "", "Service type (http, stdio, filesystem)")
	c.Flags().StringVar(&health, "health", "", "Health check URL")
	c.Flags().StringVar(&dashboard, "dashboard", "", "Dashboard URL")
	c.Flags().StringVar(&mcpURL, "mcp-url", "", "MCP URL for remote registration")
	c.Flags().StringVar(&logPath, "log", "", "Log file path")
	c.Flags().StringVar(&path, "path", "", "Path to filesystem root (used when type=filesystem)")
	c.Flags().StringVar(&buildCommand, "build-command", "", "Rebuild shell command")
	c.Flags().StringArrayVar(&envPairs, "env", nil, "KEY=VAL (repeatable)")
	c.Flags().StringArrayVar(&cmdArgs, "args", nil, "Binary args (repeatable)")
	return c
}

func prompt(label string) string {
	fmt.Print(label)
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimRight(line, "\n\r")
}

func cmdConfigEdit() *cobra.Command {
	return &cobra.Command{
		Use:     "config",
		Short:   "Open ~/.mcp-local/config.yaml in $EDITOR",
		Example: `  mcp-local config`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := config.GetConfigPath()
			if err != nil {
				return err
			}
			if _, err := os.Stat(p); os.IsNotExist(err) {
				if err := config.SaveConfig(&config.ManagerConfig{Services: nil}); err != nil {
					return err
				}
			}
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = os.Getenv("VISUAL")
			}
			if editor == "" {
				editor = "nano"
			}
			c := exec.Command(editor, p)
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			return c.Run()
		},
	}
}

func cmdRemove() *cobra.Command {
	var deregister bool
	c := &cobra.Command{
		Use:   "remove <service>",
		Short: "Remove a service from config and deregister from all agents",
		Args:  cobra.ExactArgs(1),
		Example: `  mcp-local remove my-server
  mcp-local remove my-server --deregister=false`,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cfg := mustConfig()
			if _, err := cfg.ServiceNamed(name); err != nil {
				return err
			}
			if deregister {
				svc, _ := cfg.ServiceNamed(name)
				if err := deregisterHosts(cmd, svc, agents.TargetsFromConfig(*cfg), true); err != nil {
					return fmt.Errorf("%w\n%s is still in config.yaml: fix the agent config file(s) above and retry, or run: mcp-local remove %s --deregister=false", err, name, name)
				}
			}
			var remaining []config.ServiceConfig
			for _, s := range cfg.Services {
				if s.Name != name {
					remaining = append(remaining, s)
				}
			}
			cfg.Services = remaining
			if err := config.SaveConfig(cfg); err != nil {
				return err
			}
			fmt.Printf("Removed service %q\n", name)
			return nil
		},
	}
	c.Flags().BoolVar(&deregister, "deregister", true, "deregister from agents before removing")
	return c
}

func cmdEdit() *cobra.Command {
	return &cobra.Command{
		Use:     "edit <service>",
		Short:   "Edit a service configuration (TUI)",
		Args:    cobra.ExactArgs(1),
		Example: `  mcp-local edit my-server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cfg := mustConfig()
			svc, err := cfg.ServiceNamed(name)
			if err != nil {
				return err
			}
			if err := tui.RunServiceEditor(&svc); err != nil {
				if errors.Is(err, tui.ErrCancelled) {
					fmt.Println("No changes saved.")
					return nil
				}
				return err
			}
			for i := range cfg.Services {
				if cfg.Services[i].Name == name {
					cfg.Services[i] = svc
					break
				}
			}
			if err := config.SaveConfig(cfg); err != nil {
				return err
			}
			fmt.Printf("Updated service %q\n", name)
			return nil
		},
	}
}

func cmdRegistered() *cobra.Command {
	return &cobra.Command{
		Use:   "registered [service]",
		Short: "Show which agents each service is registered with",
		Example: `  mcp-local registered
  mcp-local registered my-server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			targets := agents.TargetsFromConfig(*cfg)
			names := cfg.SortedNames()
			if len(args) > 0 {
				if _, err := cfg.ServiceNamed(args[0]); err != nil {
					return err
				}
				names = []string{args[0]}
			}
			for _, name := range names {
				ocRegistered := false
				curRegistered := false
				claudeRegistered := false
				if targets.OpenCode {
					ocRegistered = isRegisteredInFile(opencode.ConfigPath(), name)
				}
				if targets.Cursor {
					curRegistered = isRegisteredInFile(cursor.ConfigPath(), name)
				}
				if targets.Claude {
					claudeRegistered = isRegisteredInFile(claudedesktop.ConfigPath(), name)
				}
				ocStr := "❌"
				if ocRegistered {
					ocStr = "✅"
				}
				curStr := "❌"
				if curRegistered {
					curStr = "✅"
				}
				claudeStr := "❌"
				if claudeRegistered {
					claudeStr = "✅"
				}
				fmt.Printf("%-24s  OpenCode: %s  Cursor: %s  Claude: %s\n", name, ocStr, curStr, claudeStr)
			}
			return nil
		},
	}
}

func isRegisteredInFile(path, name string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	m, err := jsonagent.ParseJSONC(data)
	if err != nil {
		return false
	}
	for _, key := range []string{"mcp", "mcpServers"} {
		if block, ok := m[key].(map[string]interface{}); ok {
			if _, exists := block[name]; exists {
				return true
			}
		}
	}
	return false
}

func cmdValidate() *cobra.Command {
	return &cobra.Command{
		Use:     "validate",
		Short:   "Validate configuration for common issues",
		Example: `  mcp-local validate`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			errors := 0
			warnings := 0
			for _, svc := range cfg.Services {
				prefix := fmt.Sprintf("[%s]", svc.Name)
				switch {
				case svc.Command == "" && !svc.IsHTTP():
					fmt.Fprintf(cmd.ErrOrStderr(), "%s ERROR: stdio service has no command\n", prefix)
					errors++
				case svc.Command == "" && svc.MCPEndpoint() == "":
					fmt.Fprintf(cmd.ErrOrStderr(), "%s ERROR: no command or mcp_url configured\n", prefix)
					errors++
				case svc.Command != "":
					if err := findCommand(config.ExpandPath(svc.Command)); err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "%s ERROR: command not found: %v\n", prefix, err)
						errors++
					}
				}
				for _, dep := range svc.Deps {
					expanded := config.ExpandPath(dep)
					if _, err := os.Stat(expanded); os.IsNotExist(err) {
						fmt.Fprintf(cmd.ErrOrStderr(), "%s ERROR: dependency missing: %s\n", prefix, expanded)
						errors++
					}
				}
				if svc.Port > 0 && svc.Port < 1024 && svc.Port != 80 && svc.Port != 443 {
					fmt.Fprintf(cmd.ErrOrStderr(), "%s WARNING: port %d may require root\n", prefix, svc.Port)
					warnings++
				}
				if svc.MCPURL != "" && !strings.HasPrefix(svc.MCPURL, "http://") && !strings.HasPrefix(svc.MCPURL, "https://") {
					fmt.Fprintf(cmd.ErrOrStderr(), "%s WARNING: mcp_url should start with http:// or https://\n", prefix)
					warnings++
				}
			}
			if errors == 0 && warnings == 0 {
				fmt.Println("✅ Configuration is valid")
			} else {
				fmt.Printf("\n%d error(s), %d warning(s)\n", errors, warnings)
			}
			if errors > 0 {
				return fmt.Errorf("validation failed")
			}
			return nil
		},
	}
}

// findCommand checks a command: a bare name (npx, uvx) is looked up on PATH, a path must exist.
func findCommand(c string) error {
	if !strings.ContainsRune(c, '/') {
		_, err := exec.LookPath(c)
		return err
	}
	_, err := os.Stat(c)
	return err
}

func cmdRegister() *cobra.Command {
	var dryRun bool
	c := &cobra.Command{
		Use:   "register [service]",
		Short: "Register service(s) with OpenCode, Cursor, and Claude",
		Example: `  mcp-local register
  mcp-local register my-server
  mcp-local register --dry-run`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			names := cfg.SortedNames()
			if len(args) > 0 {
				if _, err := cfg.ServiceNamed(args[0]); err != nil {
					return err
				}
				names = []string{args[0]}
			}
			targets := agents.TargetsFromConfig(*cfg)
			var toRegister []config.ServiceConfig
			for _, name := range names {
				svc, _ := cfg.ServiceNamed(name)
				if svc.MCPURL == "" && svc.Command == "" {
					fmt.Printf("  %-24s skipped (no mcp_url or command)\n", name)
					continue
				}
				toRegister = append(toRegister, svc)
			}
			if len(toRegister) == 0 {
				return nil
			}
			if dryRun {
				fmt.Println("[dry-run] Would register:")
				for _, svc := range toRegister {
					if targets.OpenCode {
						fmt.Printf("  %-24s → opencode\n", svc.Name)
					}
					if targets.Cursor {
						fmt.Printf("  %-24s → cursor\n", svc.Name)
					}
					if targets.Claude {
						fmt.Printf("  %-24s → claude\n", svc.Name)
					}
				}
				return nil
			}
			if err := agents.RegisterAll(toRegister, targets); err != nil {
				return err
			}
			for _, svc := range toRegister {
				if targets.OpenCode {
					fmt.Printf("  %-24s → opencode\n", svc.Name)
				}
				if targets.Cursor {
					fmt.Printf("  %-24s → cursor\n", svc.Name)
				}
				if targets.Claude {
					fmt.Printf("  %-24s → claude\n", svc.Name)
				}
			}
			return nil
		},
	}
	c.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without writing")
	return c
}

func cmdDeregister() *cobra.Command {
	var dryRun bool
	c := &cobra.Command{
		Use:   "deregister [service]",
		Short: "Remove MCP entries from OpenCode, Cursor, and Claude configs",
		Long: `Without a service name: removes the entries of HTTP services with a port that is not in use,
and only entries mcp-local wrote itself (entries it did not write are skipped with a note).
With a name: always removes that service's entries.`,
		Example: `  mcp-local deregister
  mcp-local deregister my-server
  mcp-local deregister --dry-run`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			force := len(args) > 0
			names := cfg.SortedNames()
			if force {
				if _, err := cfg.ServiceNamed(args[0]); err != nil {
					return err
				}
				names = []string{args[0]}
			}
			targets := agents.TargetsFromConfig(*cfg)
			var toDeregister []string
			for _, name := range names {
				svc, _ := cfg.ServiceNamed(name)
				if !force && !(svc.IsHTTP() && svc.Port > 0) {
					continue
				}
				if !force && portIsRunning(svc.Port) {
					fmt.Printf("  %-24s skipped (running)\n", name)
					continue
				}
				toDeregister = append(toDeregister, name)
			}
			if dryRun {
				fmt.Println("[dry-run] Would deregister:")
				for _, name := range toDeregister {
					if targets.OpenCode {
						fmt.Printf("  %-24s ← opencode\n", name)
					}
					if targets.Cursor {
						fmt.Printf("  %-24s ← cursor\n", name)
					}
					if targets.Claude {
						fmt.Printf("  %-24s ← claude\n", name)
					}
				}
				return nil
			}
			var errs []error
			for _, name := range toDeregister {
				svc, _ := cfg.ServiceNamed(name)
				errs = append(errs, deregisterHosts(cmd, svc, targets, force))
			}
			return errors.Join(errs...)
		},
	}
	c.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without writing")
	return c
}

// portIsRunning is swappable in tests.
var portIsRunning = portutil.IsRunning

// deregisterHosts removes svc from each enabled host, reporting entries removed or skipped,
// and returns every host's error joined. explicit also removes entries mcp-local didn't write.
func deregisterHosts(cmd *cobra.Command, svc config.ServiceConfig, targets agents.Targets, explicit bool) error {
	var errs []error
	for _, host := range targets.Hosts() {
		res, err := agents.Deregister(svc, host, explicit)
		switch {
		case err != nil:
			fmt.Fprintf(cmd.ErrOrStderr(), "  %-24s %s ERROR: %v\n", svc.Name, host, err)
			errs = append(errs, fmt.Errorf("%s: %s: %w", svc.Name, host, err))
		case res == agents.Removed:
			fmt.Fprintf(cmd.OutOrStdout(), "  %-24s %s removed\n", svc.Name, host)
		case res == agents.NotOwned:
			fmt.Fprintf(cmd.OutOrStdout(), "  %-24s %s skipped (entry not written by mcp-local; name the service to remove it)\n", svc.Name, host)
		}
	}
	return errors.Join(errs...)
}

func cmdLogsUnified() *cobra.Command {
	return &cobra.Command{
		Use:     "logs",
		Short:   "Tail logs for all services (Ctrl+C to stop)",
		Example: `  mcp-local logs`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := mustConfig()
			var entries []logs.Entry
			for _, name := range cfg.SortedNames() {
				svc, err := cfg.ServiceNamed(name)
				if err != nil {
					continue
				}
				entries = append(entries, logs.Entry{Name: name, Path: process.LogPath(svc)})
			}
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			fmt.Println("Tailing logs (Press Ctrl+C to stop)...")
			logs.StreamLogs(ctx, entries, os.Stdout, os.Stderr)
		},
	}
}

func cmdTools() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Manage tool tiers (TUI)",
		Example: `  mcp-local tools sync my-server
  mcp-local tools apply my-server
  mcp-local json tools my-server`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := mustConfig()
			if err := tui.RunToolManager(cfg); err != nil {
				if errors.Is(err, tui.ErrCancelled) {
					fmt.Println("No changes saved.")
					return
				}
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
	cmd.AddCommand(cmdToolsSync())
	cmd.AddCommand(cmdToolsApply())
	return cmd
}

func cmdToolsApply() *cobra.Command {
	return &cobra.Command{
		Use:     "apply <service>",
		Short:   "Write ~/.astcache/tools.json from config (restart service to load)",
		Args:    cobra.ExactArgs(1),
		Example: `  mcp-local tools apply my-server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			svc, err := cfg.ServiceNamed(args[0])
			if err != nil {
				return err
			}
			if len(svc.Tools) == 0 {
				return fmt.Errorf("no tools configured for %s (run tools sync first)", args[0])
			}
			path := asttools.ToolsConfigPath(svc)
			if err := asttools.WriteToolsJSON(path, svc.Tools); err != nil {
				return err
			}
			fmt.Printf("Wrote %s (%d tool overrides)\n", path, asttools.OverridesCount(svc.Tools))
			fmt.Printf("Restart %s to apply: mcp-local restart %s\n", args[0], args[0])
			return nil
		},
	}
}

func cmdJSON() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "json",
		Short:   "Emit JSON artifacts from config",
		Example: `  mcp-local json tools my-server`,
	}
	cmd.AddCommand(&cobra.Command{
		Use:     "tools <service>",
		Short:   "Print tools.json overrides for a service",
		Args:    cobra.ExactArgs(1),
		Example: `  mcp-local json tools my-server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			svc, err := cfg.ServiceNamed(args[0])
			if err != nil {
				return err
			}
			if len(svc.Tools) == 0 {
				return fmt.Errorf("no tools configured for %s", args[0])
			}
			b, err := asttools.MarshalOverrides(svc.Tools)
			if err != nil {
				return err
			}
			fmt.Println(string(b))
			return nil
		},
	})
	return cmd
}

func cmdToolsSync() *cobra.Command {
	return &cobra.Command{
		Use:     "sync <service>",
		Short:   "Sync tools list from a running MCP server",
		Args:    cobra.ExactArgs(1),
		Example: `  mcp-local tools sync my-server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			svc, err := cfg.ServiceNamed(args[0])
			if err != nil {
				return err
			}
			mcpURL := strings.TrimSpace(svc.MCPURL)
			if mcpURL == "" && svc.Port > 0 {
				mcpURL = fmt.Sprintf("http://localhost:%d/mcp", svc.Port)
			}
			if mcpURL == "" {
				return fmt.Errorf("no mcp_url or port for %s", args[0])
			}
			tools, err := fetchToolsList(mcpURL)
			if err != nil {
				return fmt.Errorf("failed to fetch tools: %w", err)
			}
			existing := make(map[string]config.ToolConfig)
			for _, t := range svc.Tools {
				existing[t.Name] = t
			}
			merged := mergeSyncedTools(tools, existing)
			for i := range cfg.Services {
				if cfg.Services[i].Name == args[0] {
					cfg.Services[i].Tools = merged
					break
				}
			}
			if err := config.SaveConfig(cfg); err != nil {
				return err
			}
			fmt.Printf("Synced %d tools for %s\n", len(merged), args[0])
			for _, t := range merged {
				en := "enabled"
				if !t.Enabled {
					en = "disabled"
				}
				fmt.Printf("  %-30s  %s  tier=%s\n", t.Name, en, t.Tier)
			}
			return nil
		},
	}
}

// mergeSyncedTools keeps the user's settings for each tool the server lists. The server's
// own description is not stored: only a description the user changed is an override.
func mergeSyncedTools(server []config.ToolConfig, existing map[string]config.ToolConfig) []config.ToolConfig {
	var merged []config.ToolConfig
	for _, t := range server {
		serverDesc := t.Description
		t.Description = ""
		if prev, ok := existing[t.Name]; ok {
			t.Enabled, t.Tier = prev.Enabled, prev.Tier
			if prev.Description != serverDesc {
				t.Description = prev.Description
			}
		}
		merged = append(merged, t)
	}
	return merged
}

type jsonRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      *int        `json:"id,omitempty"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  struct {
		Tools []struct {
			Name        string                 `json:"name"`
			Description string                 `json:"description"`
			InputSchema map[string]interface{} `json:"inputSchema"`
		} `json:"tools"`
	} `json:"result,omitempty"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

const mcpClientProtocolVersion = "2025-06-18"

func intPtr(i int) *int { return &i }

// rpcCall POSTs req and returns the response whose id matches (JSON or SSE body).
// Notifications (nil ID) return a nil response.
func rpcCall(ctx context.Context, c *bridge.Client, req jsonRPCRequest) (*jsonRPCResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var resp *jsonRPCResponse
	err = c.Post(ctx, body, func(msg []byte) error {
		var r jsonRPCResponse
		if json.Unmarshal(msg, &r) == nil && req.ID != nil && string(r.ID) == strconv.Itoa(*req.ID) {
			resp = &r
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if req.ID != nil && resp == nil {
		return nil, fmt.Errorf("%s: no response", req.Method)
	}
	if resp != nil && resp.Error != nil {
		return nil, fmt.Errorf("%s error: %s", req.Method, resp.Error.Message)
	}
	return resp, nil
}

func fetchToolsList(mcpURL string) ([]config.ToolConfig, error) {
	ctx := context.Background()
	c := &bridge.Client{URL: mcpURL, Timeout: 10 * time.Second}
	defer c.Delete(ctx)
	_, err := rpcCall(ctx, c, jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      intPtr(1),
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": mcpClientProtocolVersion,
			"capabilities":    map[string]interface{}{},
			"clientInfo":      map[string]string{"name": "mcp-local", "version": "0.1.0"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("initialize: %w", err)
	}
	if _, err := rpcCall(ctx, c, jsonRPCRequest{JSONRPC: "2.0", Method: "notifications/initialized"}); err != nil {
		return nil, fmt.Errorf("send initialized notification: %w", err)
	}
	rpcResp, err := rpcCall(ctx, c, jsonRPCRequest{JSONRPC: "2.0", ID: intPtr(2), Method: "tools/list"})
	if err != nil {
		return nil, err
	}
	var tools []config.ToolConfig
	for _, t := range rpcResp.Result.Tools {
		tools = append(tools, config.ToolConfig{
			Name:        t.Name,
			Description: t.Description,
			Enabled:     true,
			Tier:        "",
			InputSchema: t.InputSchema,
		})
	}
	return tools, nil
}

func cmdImport() *cobra.Command {
	var source string
	c := &cobra.Command{
		Use:   "import",
		Short: "Import MCP configs from agents into config.yaml",
		Long:  "Reads existing MCP configurations from OpenCode, Cursor, or Claude Desktop and imports them as services in config.yaml.",
		Example: `  mcp-local import
  mcp-local import --from opencode
  mcp-local import --from cursor
  mcp-local import --from claude`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			existing := make(map[string]bool)
			for _, s := range cfg.Services {
				existing[s.Name] = true
			}
			var imported int
			sources := []string{"opencode", "cursor", "claude"}
			if source != "" {
				sources = []string{source}
			}
			for _, src := range sources {
				var path string
				var blockKey string
				switch src {
				case "opencode":
					path = opencode.ConfigPath()
					blockKey = "mcp"
				case "cursor":
					path = cursor.ConfigPath()
					blockKey = "mcpServers"
				case "claude":
					path = claudedesktop.ConfigPath()
					blockKey = "mcpServers"
				}
				data, err := os.ReadFile(path)
				if err != nil {
					continue
				}
				m, err := jsonagent.ParseJSONC(data)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "  ⚠️  %s: %v\n", path, err)
					continue
				}
				block, _ := m[blockKey].(map[string]interface{})
				if block == nil {
					continue
				}
				for name, raw := range block {
					if existing[name] {
						continue
					}
					entry, _ := raw.(map[string]interface{})
					if entry == nil {
						continue
					}
					svc, ok := serviceFromEntry(name, entry)
					if ok {
						cfg.Services = append(cfg.Services, svc)
						existing[name] = true
						imported++
						fmt.Printf("  Imported %s from %s\n", name, src)
					}
				}
			}
			if imported == 0 {
				fmt.Println("No new services to import")
				return nil
			}
			if err := config.SaveConfig(cfg); err != nil {
				return err
			}
			fmt.Printf("\nImported %d service(s)\n", imported)
			return nil
		},
	}
	c.Flags().StringVar(&source, "from", "", "Source agent: opencode, cursor, or claude (default: all)")
	return c
}

// serviceFromEntry converts an OpenCode, Cursor or Claude Desktop entry into a service,
// keeping everything register writes back (args, env, headers) so a later register doesn't
// clobber the entry. ok is false when the entry has neither a command nor a URL.
func serviceFromEntry(name string, entry map[string]interface{}) (config.ServiceConfig, bool) {
	svc := config.ServiceConfig{Name: name}
	entryType, _ := entry["type"].(string)
	url, _ := entry["url"].(string)
	if bridge, headers, ok := bridgeEntry(entry); ok {
		svc.MCPType, svc.MCPURL, svc.Headers = "http", bridge, headers
	} else if url != "" || entryType == "remote" || entryType == "http" || entryType == "sse" || entryType == "streamable-http" {
		svc.MCPType, svc.MCPURL, svc.Headers = "http", url, stringMap(entry["headers"])
	} else {
		svc.MCPType = "stdio"
		switch c := entry["command"].(type) {
		case []interface{}:
			if cmd := stringSlice(c); len(cmd) > 0 {
				svc.Command, svc.Args = cmd[0], cmd[1:]
			}
		case string:
			svc.Command = c
		}
		svc.Args = append(svc.Args, stringSlice(entry["args"])...)
		svc.Env = stringMap(entry["env"])
		for k, v := range stringMap(entry["environment"]) {
			if svc.Env == nil {
				svc.Env = map[string]string{}
			}
			svc.Env[k] = v
		}
	}
	if len(svc.Args) == 0 {
		svc.Args = nil
	}
	return svc, svc.Command != "" || svc.MCPURL != ""
}

func stringSlice(v interface{}) []string {
	arr, _ := v.([]interface{})
	var out []string
	for _, a := range arr {
		if s, ok := a.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func stringMap(v interface{}) map[string]string {
	m, _ := v.(map[string]interface{})
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, val := range m {
		out[k], _ = val.(string)
	}
	return out
}

// bridgeEntry recognizes a Claude Desktop entry written by mcp-local:
// {"command": ".../mcp-local", "args": ["bridge", url, "--header", "K=V", ...]}. Entries with
// any other bridge flags are not recognized, so importing them as stdio keeps those flags.
func bridgeEntry(entry map[string]interface{}) (string, map[string]string, bool) {
	cmd, _ := entry["command"].(string)
	args := stringSlice(entry["args"])
	if filepath.Base(cmd) != "mcp-local" || len(args) < 2 || args[0] != "bridge" {
		return "", nil, false
	}
	var headers map[string]string
	for rest := args[2:]; len(rest) > 0; rest = rest[2:] {
		if rest[0] != "--header" || len(rest) < 2 {
			return "", nil, false
		}
		k, v, ok := strings.Cut(rest[1], "=")
		if !ok {
			return "", nil, false
		}
		if headers == nil {
			headers = map[string]string{}
		}
		headers[k] = v
	}
	return args[1], headers, true
}

func cmdServeFS() *cobra.Command {
	var name, logPath, root string
	var port int
	c := &cobra.Command{
		Use:    "_serve-fs",
		Short:  "Start an embedded filesystem MCP server (internal)",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" || port == 0 || root == "" {
				return fmt.Errorf("usage: _serve-fs --name <name> --port <port> --path <path> [--log <log>]")
			}
			embedfs.LogTo(logPath, fmt.Sprintf("%s: starting embedfs on port %d (root: %s)", name, port, root))
			if err := embedfs.New(root).Serve(port, logPath); err != nil {
				embedfs.LogTo(logPath, fmt.Sprintf("%s: server error: %v", name, err))
				return err
			}
			return nil
		},
	}
	c.Flags().StringVar(&name, "name", "", "Service name")
	c.Flags().IntVar(&port, "port", 0, "Listen port")
	c.Flags().StringVar(&root, "path", "", "Filesystem root path")
	c.Flags().StringVar(&logPath, "log", "", "Log file path")
	return c
}
