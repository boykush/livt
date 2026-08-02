package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/boykush/livt/internal/mcp"
	"github.com/boykush/livt/internal/uri"
	"github.com/spf13/cobra"
)

// Output forms. json is the default because it is what a livt URI means
// everywhere else: the payload an MCP client reads.
const (
	formatJSON = "json"
	formatURL  = "url"
)

var (
	resolveRootDir string
	resolveFormat  string
	resolveBaseURL string
)

func init() {
	resolveCmd.Flags().StringVar(&resolveRootDir, "root", "", "path to the root of the livt repository (default: $LIVT_ROOT, then the current directory)")
	resolveCmd.Flags().StringVar(&resolveFormat, "format", formatJSON, "output form: json (the resource payload an MCP client reads) or url (its page on the deployed site)")
	resolveCmd.Flags().StringVar(&resolveBaseURL, "base-url", "", "root of the deployed site, required by --format url (e.g. https://boykush.github.io/livt)")
	rootCmd.AddCommand(resolveCmd)
}

var resolveCmd = &cobra.Command{
	Use:   "resolve <uri>",
	Short: "Resolve a livt URI against the livt repository",
	Long: `Resolve a livt URI -- the deployment-independent name for one point of the
spec -- without running an MCP client.

A rule, example, or question cited in a test comment, an Issue, or a commit
message can then be read back by whoever needs it: CI, an editor, or an agent
that is not connected to livt over MCP.

The livt repository usually lives in a separate checkout from the consumer, so point at
it with --root or the LIVT_ROOT environment variable (the flag takes precedence;
both default to the current directory).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Past argument parsing, a failure is about the URI or the livt repository, so
		// usage would only bury the reason.
		cmd.SilenceUsage = true
		return resolveURI(cmd.OutOrStdout(), resolveRoot(resolveRootDir), args[0], resolveFormat, resolveBaseURL)
	},
}

// resolveURI writes rawURI in the requested form. It is split from the command
// so its behaviour can be tested without driving cobra's global flag state.
func resolveURI(out io.Writer, root, rawURI, format, baseURL string) error {
	p, ok := uri.Parse(rawURI)
	if !ok {
		return malformedURIError(rawURI)
	}

	switch format {
	case formatJSON:
		payload, err := mcp.NewServer(mcp.Config{Root: root}, livtVersion()).Resolve(p)
		if err != nil {
			return err
		}
		body, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(out, string(body))
		return err
	case formatURL:
		if baseURL == "" {
			return fmt.Errorf("--format %s needs --base-url, the root of the deployed site (e.g. https://boykush.github.io/livt)", formatURL)
		}
		// The page path derives from the URI alone, so the livt repository is consulted
		// only to refuse a confident link to a page that is not there.
		if err := (mcp.Config{Root: root}).Verify(p); err != nil {
			return err
		}
		_, err := fmt.Fprintln(out, strings.TrimSuffix(baseURL, "/")+"/"+p.Page())
		return err
	}
	return fmt.Errorf("unknown --format %q: want %s or %s", format, formatJSON, formatURL)
}

// malformedURIError separates a URI that is not a livt URI at all from one that
// is well formed but names something the livt repository does not hold -- the two need
// different fixes, so they must not read alike.
func malformedURIError(rawURI string) error {
	var b strings.Builder
	fmt.Fprintf(&b, "malformed livt URI %q; expected one of:", rawURI)
	for _, t := range uri.Templates {
		fmt.Fprintf(&b, "\n  %s", t)
	}
	return errors.New(b.String())
}
