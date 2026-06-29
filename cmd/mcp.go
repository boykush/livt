package cmd

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/boykush/livt/internal/mcp"
	"github.com/spf13/cobra"
)

var (
	mcpRoot string
	mcpHTTP string
)

func init() {
	mcpCmd.Flags().StringVar(&mcpRoot, "root", "", "path to the livt project root holding the discovery master (default: $LIVT_ROOT, then the current directory)")
	mcpCmd.Flags().StringVar(&mcpHTTP, "http", "", "serve over Streamable HTTP at this address (e.g. localhost:5488) instead of stdio; the MCP endpoint is <addr>/mcp")
	rootCmd.AddCommand(mcpCmd)
}

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Run an MCP server exposing the discovery master over stdio or HTTP",
	Long: `Run a Model Context Protocol server that exposes the discovery master
(stories and example mappings), so an implementation repo's agent can fetch the
spec for a story or rule without reading livt's source.

By default it serves over stdio, spawned per consumer. Pass --http to instead
serve over Streamable HTTP from one long-running process, so several local repos
can share a single server (each points its MCP client at <addr>/mcp) without a
per-repo checkout of the master.

The master usually lives in a separate checkout from the consumer, so point at
it with --root or the LIVT_ROOT environment variable (the flag takes precedence;
both default to the current directory).`,
	// A server failure shouldn't print CLI usage.
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		root := resolveRoot(mcpRoot)
		srv := mcp.NewServer(mcp.Config{Root: root}, livtVersion())
		if mcpHTTP != "" {
			return srv.RunHTTP(cmd.Context(), mcpHTTP)
		}
		if err := srv.Run(cmd.Context()); !isCleanShutdown(err) {
			return err
		}
		return nil
	},
}

// isCleanShutdown reports whether err is the normal end of an stdio session: the
// client closed the stream (EOF) or cancelled the context. The SDK wraps these
// in an internal "server is closing" error that isn't comparable with errors.Is,
// so the message is matched as a fallback.
func isCleanShutdown(err error) bool {
	if err == nil {
		return true
	}
	if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "EOF") || strings.Contains(msg, "server is closing")
}

// resolveRoot selects the master location: the --root flag wins, then the
// LIVT_ROOT environment variable, then the current directory.
func resolveRoot(flagRoot string) string {
	if flagRoot != "" {
		return flagRoot
	}
	if env := os.Getenv("LIVT_ROOT"); env != "" {
		return env
	}
	return "."
}
