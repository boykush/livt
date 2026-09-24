package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/boykush/livt/internal/automation"
	"github.com/spf13/cobra"
)

var (
	automationsRepo        string
	automationsRev         string
	automationsForge       string
	automationsURLTemplate string
	automationsOut         string
	automationsChangedPath string
)

func init() {
	automationsCmd.Flags().StringVar(&automationsRepo, "repo", "", "repository the report speaks for, owner/repo (default: read from the checkout's origin)")
	automationsCmd.Flags().StringVar(&automationsRev, "rev", "", "revision scanned (default: read from the checkout)")
	automationsCmd.Flags().StringVar(&automationsForge, "forge", "", "code host to build line URLs for: github or gitlab (default: inferred from origin, otherwise no URLs)")
	automationsCmd.Flags().StringVar(&automationsURLTemplate, "url-template", "", "line URL template for a host livt does not know, e.g. {base}/src/commit/{rev}/{path}#L{line}")
	automationsCmd.Flags().StringVarP(&automationsOut, "out", "o", "", "write the report to this file (default: stdout)")
	automationsChangedCmd.Flags().StringVar(&automationsChangedPath, "path", ".", "the checkout to read the diff in")
	automationsCmd.AddCommand(automationsChangedCmd)
	rootCmd.AddCommand(automationsCmd)
}

var automationsCmd = &cobra.Command{
	Use:   "automations [path]",
	Short: "Collect what a repository's tests say they automate",
	Long: `Collect the automations an implementation repository's tests declare.

A test declares one by writing the marker and a livt URI on a comment line
above itself:

    // livt:automates livt://mapping/{story-key}/rule/{rule-id}/example/{example-id}

The comment syntax is the language's; livt only looks for the marker and the
URI, and never reads the structure around them. A livt URI without the marker
is a reference, not a claim, and is not collected; neither is one written with
placeholders, so the form can be documented without inventing a citation.

The report carries the citations and nothing else — whether a rule counts as
automated is derived when the site is built, not decided here.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		root := "."
		if len(args) == 1 {
			root = args[0]
		}

		// The report usually lands inside the tree being scanned, and a run
		// must still be able to name the revision it read.
		ignore := []string{}
		if rel, err := filepath.Rel(root, automationsOut); err == nil && automationsOut != "" {
			ignore = append(ignore, filepath.ToSlash(rel))
		}
		origin := automation.ReadOrigin(root, ignore...)
		template, err := automation.URLTemplate(automationsForge, automationsURLTemplate, origin)
		if err != nil {
			return err
		}
		opts := automation.Options{
			Repo:        firstNonEmpty(automationsRepo, origin.Repo),
			Rev:         firstNonEmpty(automationsRev, origin.Rev),
			Base:        origin.Base,
			URLTemplate: template,
		}

		report, warnings, err := automation.Scan(root, opts)
		if err != nil {
			return err
		}
		if err := writeReport(cmd.OutOrStdout(), report); err != nil {
			return err
		}
		for _, w := range warnings {
			fmt.Fprintln(os.Stderr, "warning: "+w)
		}
		if len(warnings) > 0 {
			return fmt.Errorf("%d citation(s) could not be read; the report was written without them", len(warnings))
		}
		return nil
	},
}

var automationsChangedCmd = &cobra.Command{
	Use:   "changed <base> [head]",
	Short: "Report whether a range could have changed what the tests claim",
	Long: `Report whether the diff between two revisions adds or removes a marker line.

Collecting walks every file, which is worth avoiding on a large repository,
and only an added or removed marker line can change what a report claims. A
rename, or a line that merely moved, produces neither and answers false: a
report is a dated snapshot, so a later commit that shifts a cited line leaves
the link pinned to the revision the report names. A deleted test's markers
arrive as removals and answer true, because a claim that outlives its test
leaves the board naming a rule nothing covers any more.

Prints true or false. A range that cannot be read — a new branch, a
force-push, a revision since gone — is not an answer of false: it warns and
answers true, since the cost of the extra walk is a walk, and the cost of the
wrong skip is a board that lies.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		head := "HEAD"
		if len(args) == 2 {
			head = args[1]
		}
		reportChanged(cmd.OutOrStdout(), cmd.ErrOrStderr(), automationsChangedPath, args[0], head)
		return nil
	},
}

// reportChanged answers on stdout and explains itself on stderr, so CI reads
// the answer with a plain command substitution and still sees why.
func reportChanged(out, warn io.Writer, root, base, head string) {
	changed, err := automation.Changed(root, base, head)
	if err != nil {
		fmt.Fprintln(warn, "warning: "+err.Error()+"; answering true rather than guessing")
		changed = true
	}
	fmt.Fprintln(out, changed)
}

// writeReport puts the report where the caller asked. --out exists so CI can
// land it on the path the livt repository expects without a shell redirect
// swallowing the exit code.
func writeReport(stdout interface{ Write([]byte) (int, error) }, report automation.Report) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if automationsOut == "" {
		_, err = stdout.Write(data)
		return err
	}
	if err := os.MkdirAll(filepath.Dir(automationsOut), 0o755); err != nil {
		return err
	}
	return os.WriteFile(automationsOut, data, 0o644) // #nosec G306 -- a generated report read by everyone who can read the repository
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
