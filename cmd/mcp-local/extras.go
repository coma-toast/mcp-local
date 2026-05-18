package main

import (
	"bufio"
	"bytes"
	"encoding/json"
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

	"github.com/coma-toast/mcp-local/internal/mgr/agents"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/cursor"
	"github.com/coma-toast/mcp-local/internal/mgr/logs"
	"github.com/coma-toast/mcp-local/internal/mgr/opencode"
	"github.com/coma-toast/mcp-local/internal/portutil"
	"github.com/coma-toast/mcp-local/internal/statusui"
	"github.com/coma-toast/mcp-local/internal/tui"
	"github.com/spf13/cobra"
)

func addCommands() {
	addLifecycle()

	rootCmd.AddCommand(cmdList())

	var plainStatus bool
	status := &cobra.Command{
		Use:   "status [service]",
		Short: "Live status table (q to quit); use --plain for one-shot",
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
			for _, name := range names {
				svc, _ := cfg.ServiceNamed(name)
				if svc.HealthURL == "" {
					fmt.Printf("%-24s  no health_url configured\n", name)
					continue
				}
				resp, err := client.Get(svc.HealthURL)
				if err != nil {
					fmt.Printf("%-24s  ERROR: %v\n", name, err)
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
					fmt.Printf("%-24s  ✓ %s  %s\n", name, status, pretty)
				} else {
					fmt.Fprintf(cmd.ErrOrStderr(), "%-24s  ✗ %s  %s\n", name, status, pretty)
				}
			}
			return nil
		},
	}
	rootCmd.AddCommand(health)

	rootCmd.AddCommand(cmdLog())
	rootCmd.AddCommand(cmdOpen())
	rootCmd.AddCommand(cmdRebuild())
	rootCmd.AddCommand(cmdAdd())
	rootCmd.AddCommand(cmdConfigEdit())
	rootCmd.AddCommand(cmdRegister())
	rootCmd.AddCommand(cmdDeregister())

	rootCmd.AddCommand(cmdLogsUnified())
	rootCmd.AddCommand(cmdTools())
}

func cmdList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured services",
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
		Use:   "log <service>",
		Short: "Tail the service log file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig()
			svc, err := cfg.ServiceNamed(args[0])
			if err != nil {
				return err
			}
			logPath := svc.Log
			if logPath == "" {
				logPath = filepath.Join(os.TempDir(), args[0]+".log")
			}
			f, err := os.Open(logPath)
			if err != nil {
				return fmt.Errorf("open log %s: %w", logPath, err)
			}
			defer f.Close()
			_, _ = f.Seek(0, io.SeekEnd)
			reader := bufio.NewReader(f)
			fmt.Printf("==> %s <==\n", logPath)
			for {
				line, err := reader.ReadString('\n')
				if line != "" {
					fmt.Print(line)
				}
				if err != nil {
					time.Sleep(200 * time.Millisecond)
				}
			}
		},
	}
}

func cmdOpen() *cobra.Command {
	return &cobra.Command{
		Use:   "open <service>",
		Short: "Open dashboard_url in a browser",
		Args:  cobra.ExactArgs(1),
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
		Use:   "rebuild <service>",
		Short: "Run build_cmd / build_command for a service",
		Args:  cobra.ExactArgs(1),
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
			if targets.OpenCode && svc.MCPURL != "" {
				if err := opencode.RegisterRemote(args[0], svc.MCPURL, 30000); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "  opencode: register failed: %v\n", err)
				} else {
					fmt.Printf("  opencode: registered %s → %s\n", args[0], svc.MCPURL)
				}
			}
			if targets.Cursor && svc.MCPURL != "" {
				if err := cursor.RegisterRemote(args[0], svc.MCPURL); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "  cursor: register failed: %v\n", err)
				} else {
					fmt.Printf("  cursor: registered %s → %s\n", args[0], svc.MCPURL)
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
		health       string
		dashboard    string
		mcpURL       string
		logPath      string
		buildCommand string
		envPairs     []string
		cmdArgs      []string
	)
	c := &cobra.Command{
		Use:   "add <name>",
		Short: "Add or update a service",
		Args:  cobra.ExactArgs(1),
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
				if svc.Command == "" {
					svc.Command = prompt("Command (path to binary): ")
				}
				if svc.Port == 0 {
					portStr := prompt("Port: ")
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
				if svc.Log == "" {
					suggested := filepath.Join(os.TempDir(), name+".log")
					val := prompt(fmt.Sprintf("Log path [%s]: ", suggested))
					if strings.TrimSpace(val) == "" {
						svc.Log = suggested
					} else {
						svc.Log = strings.TrimSpace(val)
					}
				}
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
	c.Flags().StringVar(&health, "health", "", "Health check URL")
	c.Flags().StringVar(&dashboard, "dashboard", "", "Dashboard URL")
	c.Flags().StringVar(&mcpURL, "mcp-url", "", "MCP URL for OpenCode remote entry")
	c.Flags().StringVar(&logPath, "log", "", "Log file path")
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
		Use:   "config",
		Short: "Open ~/.mcp-local/config.yaml in $EDITOR",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := config.GetConfigPath()
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

func cmdRegister() *cobra.Command {
	return &cobra.Command{
		Use:   "register [service]",
		Short: "Register service(s) with OpenCode and Cursor",
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
			if len(toRegister) > 0 {
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
				}
			}
			return nil
		},
	}
}

func cmdDeregister() *cobra.Command {
	return &cobra.Command{
		Use:   "deregister [service]",
		Short: "Remove MCP entries from OpenCode and Cursor configs",
		Long:  "Without a service name: removes entries for services that have mcp_url and are not running. With a name: always removes that entry.",
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
			for _, name := range names {
				svc, _ := cfg.ServiceNamed(name)
				if svc.MCPURL == "" && svc.Command == "" {
					continue
				}
				if !force && svc.Port > 0 && portutil.IsRunning(svc.Port) {
					fmt.Printf("  %-24s skipped (running)\n", name)
					continue
				}
				if targets.OpenCode {
					ok, err := opencode.Deregister(name)
					if err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "  %-24s opencode ERROR: %v\n", name, err)
					} else if ok {
						fmt.Printf("  %-24s opencode removed\n", name)
					}
				}
				if targets.Cursor {
					err := cursor.Deregister(name)
					if err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "  %-24s cursor ERROR: %v\n", name, err)
					} else {
						fmt.Printf("  %-24s cursor removed\n", name)
					}
				}
			}
			return nil
		},
	}
}

func cmdLogsUnified() *cobra.Command {
	return &cobra.Command{
		Use:   "logs",
		Short: "Tail logs for all services (Ctrl+C to stop)",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := mustConfig()
			var entries []logs.Entry
			for _, name := range cfg.SortedNames() {
				svc, err := cfg.ServiceNamed(name)
				if err != nil {
					continue
				}
				entries = append(entries, logs.Entry{Name: name, Path: svc.Log})
			}
			stop := make(chan struct{})
			go func() {
				sigChan := make(chan os.Signal, 1)
				signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
				<-sigChan
				close(stop)
			}()
			fmt.Println("Tailing logs (Press Ctrl+C to stop)...")
			logs.StreamLogs(entries, stop)
		},
	}
}

func cmdTools() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Manage tool tiers (TUI)",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := mustConfig()
			if err := tui.RunToolManager(cfg); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
	cmd.AddCommand(cmdToolsSync())
	return cmd
}

func cmdToolsSync() *cobra.Command {
	return &cobra.Command{
		Use:   "sync <service>",
		Short: "Sync tools list from a running MCP server",
		Args:  cobra.ExactArgs(1),
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
			var merged []config.ToolConfig
			for _, t := range tools {
				if prev, ok := existing[t.Name]; ok {
					t.Enabled = prev.Enabled
					if prev.Tier != "" {
						t.Tier = prev.Tier
					}
					if prev.Description != "" {
						t.Description = prev.Description
					}
				}
				merged = append(merged, t)
			}
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

type jsonRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type jsonRPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
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

func fetchToolsList(mcpURL string) ([]config.ToolConfig, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	initReq := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"clientInfo":      map[string]string{"name": "mcp-local", "version": "0.1.0"},
		},
	}
	initBody, _ := json.Marshal(initReq)
	resp, err := client.Post(mcpURL, "application/json", bytes.NewReader(initBody))
	if err != nil {
		return nil, err
	}
	var initResp jsonRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&initResp); err != nil {
		resp.Body.Close()
		return nil, fmt.Errorf("initialize response: %w", err)
	}
	resp.Body.Close()
	if initResp.Error != nil {
		return nil, fmt.Errorf("initialize error: %s", initResp.Error.Message)
	}

	initializedReq := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  "notifications/initialized",
	}
	initializedBody, _ := json.Marshal(initializedReq)
	resp, err = client.Post(mcpURL, "application/json", bytes.NewReader(initializedBody))
	if err != nil {
		return nil, fmt.Errorf("send initialized notification: %w", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	toolsReq := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/list",
	}
	toolsBody, _ := json.Marshal(toolsReq)
	resp, err = client.Post(mcpURL, "application/json", bytes.NewReader(toolsBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rpcResp jsonRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return nil, err
	}
	if rpcResp.Error != nil {
		return nil, fmt.Errorf("MCP error: %s", rpcResp.Error.Message)
	}
	var tools []config.ToolConfig
	for _, t := range rpcResp.Result.Tools {
		tools = append(tools, config.ToolConfig{
			Name:        t.Name,
			Description: t.Description,
			Enabled:     true,
			Tier:        "core",
			InputSchema: t.InputSchema,
		})
	}
	return tools, nil
}
