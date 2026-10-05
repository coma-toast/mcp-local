package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/coma-toast/mcp-local/internal/bridge"
	"github.com/spf13/cobra"
)

func cmdBridge() *cobra.Command {
	var headers []string
	var timeout time.Duration
	c := &cobra.Command{
		Use:   "bridge <url>",
		Short: "Bridge stdio JSON-RPC to a Streamable HTTP MCP server",
		Long: `Reads newline-delimited JSON-RPC from stdin, POSTs each message to <url> using the
MCP Streamable HTTP transport, and writes every server message to stdout as one line.
Used for hosts that only launch stdio servers (e.g. Claude Desktop). Logs go to stderr.`,
		Args: cobra.ExactArgs(1),
		Example: `  mcp-local bridge http://localhost:7821/mcp
  mcp-local bridge https://api.example.com/mcp --header "Authorization=Bearer $TOKEN"`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			h := http.Header{}
			for _, kv := range headers {
				k, v, ok := strings.Cut(kv, "=")
				if !ok || strings.TrimSpace(k) == "" {
					return fmt.Errorf("--header %q: want key=value", kv)
				}
				h.Add(strings.TrimSpace(k), strings.TrimSpace(v))
			}
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			return bridge.Run(ctx, bridge.Options{
				URL:     args[0],
				Header:  h,
				Timeout: timeout,
				Stdin:   os.Stdin,
				Stdout:  os.Stdout,
				Log:     log.New(os.Stderr, "mcp-local bridge: ", log.LstdFlags),
			})
		},
	}
	c.Flags().StringArrayVar(&headers, "header", nil, "extra HTTP header key=value (repeatable)")
	c.Flags().DurationVar(&timeout, "timeout", 60*time.Second, "per-request timeout")
	return c
}
