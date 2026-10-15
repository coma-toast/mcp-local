package main

import (
	"fmt"
	"io"
	"os"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/tui"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

// version is set at release build time (-ldflags "-X main.version=...").
var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "mcp-local",
	Version: version,
	// main prints the error; without this cobra prints it too ("Error: X" then "X").
	SilenceErrors: true,
	Short:         "Local MCP Service Manager",
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
	os.Exit(execute(os.Stderr))
}

// execute runs the root command and prints a failing command's error once, to stderr.
func execute(stderr io.Writer) int {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

// Swappable in tests.
var (
	runWizard       = tui.RunConfigWizard
	stdinIsTerminal = func() bool { return isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd()) }
)

// ensureConfig loads config.yaml. The setup wizard runs only when there is no file at all
// (not even a dangling symlink) and stdin is a terminal; any other load error is returned,
// so a typo or a sync conflict never gets the file overwritten.
func ensureConfig() (*config.ManagerConfig, error) {
	cfg, err := config.LoadConfig()
	if err == nil {
		return cfg, nil
	}
	path, perr := config.GetConfigPath()
	if perr != nil {
		return nil, perr
	}
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("cannot load %s: %w", path, err)
	}
	if fi, lerr := os.Lstat(path); lerr == nil && fi.Mode()&os.ModeSymlink != 0 {
		target, _ := os.Readlink(path)
		return nil, fmt.Errorf("%s is a symlink to %s, which does not exist", path, target)
	}
	if !stdinIsTerminal() {
		return nil, fmt.Errorf("%s does not exist; copy examples/config.starter.yaml there, or run mcp-local in a terminal for the setup wizard", path)
	}
	fmt.Fprintf(os.Stderr, "⚠️  No configuration at %s. Launching setup wizard...\n", path)
	wizardCfg, err := runWizard()
	if err != nil {
		return nil, fmt.Errorf("config wizard failed: %v", err)
	}
	if err := config.SaveConfig(wizardCfg); err != nil {
		return nil, fmt.Errorf("failed to save wizard config: %v", err)
	}
	fmt.Fprintf(os.Stderr, "✅ Configuration saved to %s\n", path)
	return wizardCfg, nil
}

func mustConfig() *config.ManagerConfig {
	cfg, err := ensureConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌", err)
		os.Exit(1)
	}
	return cfg
}
