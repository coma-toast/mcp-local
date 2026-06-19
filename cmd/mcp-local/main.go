package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mcp-local",
	Short: "Local MCP Service Manager",
	Long: `mcp-local is a unified control plane for managing multiple local MCP servers,
lifecycle, OpenCode registration, and tool configuration.`,
	Example: `  mcp-local add github --type stdio --command npx --args "-y" --args "@modelcontextprotocol/server-github" --env GITHUB_TOKEN=ghp_xxxx
  mcp-local add filesystem --type filesystem --path ~/projects --port 4000
  mcp-local add playwright --type stdio --command npx --args "-y" --args "@playwright/mcp"
  mcp-local add fetch --type stdio --command uvx --args "mcp-server-fetch"
  mcp-local add remote-api --type http --mcp-url https://api.example.com/mcp --health https://api.example.com/health
  mcp-local start --all
  mcp-local status --plain
  mcp-local register --dry-run`,
}

func main() {
	addCommands()
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

func mustConfig() *config.ManagerConfig {
	cfg, err := ensureConfig()
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
