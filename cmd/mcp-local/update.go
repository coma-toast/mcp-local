package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/coma-toast/mcp-local/internal/selfupdate"
	"github.com/spf13/cobra"
)

// Swappable in tests.
var (
	newUpdater = func() *selfupdate.Updater {
		token := os.Getenv("GH_TOKEN")
		if token == "" {
			token = os.Getenv("GITHUB_TOKEN")
		}
		return &selfupdate.Updater{Token: token}
	}
	updateTarget = selfupdate.Executable
)

func addUpdate() {
	var check, force bool
	var tag string
	update := &cobra.Command{
		Use:   "update",
		Short: "Update mcp-local to the latest GitHub release",
		Long: `Downloads the release build for this OS and architecture from
github.com/coma-toast/mcp-local, checks it against the release's checksums.txt,
and replaces this binary (symlinks are followed to the real file).

Development builds (version "dev") and reinstalls or downgrades need --force.
Running "mcp-local bridge" processes keep the old version until their host restarts them.
Set GH_TOKEN or GITHUB_TOKEN to avoid GitHub's anonymous API rate limit.`,
		Example: `  mcp-local update --check
  mcp-local update
  mcp-local update --version v0.1.8 --force`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runUpdate(cmd, os.Stdout, version, tag, check, force)
		},
	}
	update.Flags().BoolVar(&check, "check", false, "only report whether an update is available")
	update.Flags().StringVar(&tag, "version", "", "install this release tag instead of the latest")
	update.Flags().BoolVar(&force, "force", false, "install even if not newer (dev builds, reinstalls, downgrades)")
	rootCmd.AddCommand(update)
}

func runUpdate(cmd *cobra.Command, out io.Writer, current, tag string, check, force bool) error {
	u := newUpdater()
	ctx := cmd.Context()
	var rel *selfupdate.Release
	var err error
	if tag != "" {
		rel, err = u.ByTag(ctx, tag)
	} else {
		rel, err = u.Latest(ctx)
	}
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "Current: %s\nRelease: %s\n", current, rel.Version)

	newer, cmpErr := selfupdate.Compare(current, rel.Version)
	switch {
	case errors.Is(cmpErr, selfupdate.ErrDevBuild):
		if check {
			fmt.Fprintln(out, "This is a development build; run with --force to replace it with the release.")
			return nil
		}
		if !force {
			return fmt.Errorf("this is a development build (version %q); re-run with --force to replace it with %s", current, rel.Tag)
		}
	case cmpErr != nil:
		return cmpErr
	case check:
		if newer {
			fmt.Fprintf(out, "Update available: run mcp-local update\n")
		} else {
			fmt.Fprintln(out, "✅ Up to date")
		}
		return nil
	case !newer && !force:
		fmt.Fprintln(out, "✅ Up to date (use --force to reinstall or downgrade)")
		return nil
	}

	exe, err := updateTarget()
	if err != nil {
		return fmt.Errorf("cannot find this binary: %w", err)
	}
	fmt.Fprintf(out, "Installing %s to %s...\n", u.ArchiveName(rel.Version), exe)
	if err := u.Install(ctx, rel, exe); err != nil {
		return err
	}
	fmt.Fprintf(out, "✅ Updated mcp-local %s → %s\n", current, rel.Version)
	return nil
}
