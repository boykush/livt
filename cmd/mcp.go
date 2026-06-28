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

var mcpRoot string

func init() {
	mcpCmd.Flags().StringVar(&mcpRoot, "root", "", "path to the livt project root holding the discovery master (default: $LIVT_ROOT, then the current directory)")
	rootCmd.AddCommand(mcpCmd)
}

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Run an MCP server exposing the discovery master over stdio",
	Long: `Run a Model Context Protocol server over stdio that exposes the discovery
master (stories and example mappings), so an implementation repo's agent can
fetch the spec for a story or rule without reading livt's source.

The master usually lives in a separate checkout from the consumer, so point at
it with --root or the LIVT_ROOT environment variable (the flag takes precedence;
both default to the current directory).`,
	// A server failure shouldn't print CLI usage.
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		root := resolveRoot(mcpRoot)
		err := mcp.NewServer(mcp.Config{Root: root}, livtVersion()).Run(cmd.Context())
		if isCleanShutdown(err) {
			return nil
		}
		return err
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
