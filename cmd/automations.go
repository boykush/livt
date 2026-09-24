package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/boykush/livt/internal/automation"
	"github.com/boykush/livt/internal/mcp"
	"github.com/boykush/livt/internal/uri"
	"github.com/spf13/cobra"
)

var (
	automationsRepo          string
	automationsRev           string
	automationsForge         string
	automationsURLTemplate   string
	automationsOut           string
	automationsChangedPath   string
	automationsVerifyRoot    string
	automationsVerifyReports string
)

func init() {
	automationsCmd.Flags().StringVar(&automationsRepo, "repo", "", "repository the report speaks for, owner/repo (default: read from the checkout's origin)")
	automationsCmd.Flags().StringVar(&automationsRev, "rev", "", "revision scanned (default: read from the checkout)")
	automationsCmd.Flags().StringVar(&automationsForge, "forge", "", "code host to build line URLs for: github or gitlab (default: inferred from origin, otherwise no URLs)")
	automationsCmd.Flags().StringVar(&automationsURLTemplate, "url-template", "", "line URL template for a host livt does not know, e.g. {base}/src/commit/{rev}/{path}#L{line}")
	automationsCmd.Flags().StringVarP(&automationsOut, "out", "o", "", "write the report to this file (default: stdout)")
	automationsChangedCmd.Flags().StringVar(&automationsChangedPath, "path", ".", "the checkout to read the diff in")
	automationsVerifyCmd.Flags().StringVar(&automationsVerifyRoot, "root", "", "path to the root of the livt repository (default: $LIVT_ROOT, then the current directory)")
	automationsVerifyCmd.Flags().StringVar(&automationsVerifyReports, "reports", "", reportsFlagUsage+" (default: "+defaultReportsDir+" under --root)")
	automationsCmd.AddCommand(automationsChangedCmd, automationsVerifyCmd)
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

var automationsVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Check that every livt URI the collected reports cite resolves",
	Long: `Resolve every livt URI the collected reports cite against the livt repository
holding them.

A citation is written by hand in a test comment, and collecting only reads it:
a mistyped marker, or one naming an id that moved, lands in the report as a
claim about a point of the spec that is not there. Nothing downstream says so
-- an unresolvable URI simply matches no rule, so the board shows a rule as
un-automated and the test goes on passing.

Run it where the reports land, on the pull request that commits them, and the
typo fails there instead of on main.

A URI that resolves to a retired rule is not a failure. The rule closed on the
spec's side while the test naming it still runs and passes; that gap is the
board's to show, not this command's to break a build over.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		root := resolveRoot(automationsVerifyRoot)
		// The reports live in the livt repository being resolved against, so a
		// --root pointing elsewhere must take them with it: the default read
		// from the wrong checkout would find no reports and pass on nothing.
		reports := automationsVerifyReports
		if reports == "" {
			reports = filepath.Join(root, defaultReportsDir)
		}
		return verifyCitations(cmd.OutOrStdout(), root, reports)
	},
}

// verifyCitations resolves every citation the reports carry, naming the ones
// that resolve to nothing. Reports are read through the same loader the build
// uses, so a missing directory is no error here either: a livt repository
// nobody has pointed an implementation at has no citations to check.
func verifyCitations(out io.Writer, root, reportsDir string) error {
	idx, err := automation.Load(reportsDir)
	if err != nil {
		return err
	}
	cfg := mcp.Config{Root: root}
	var checked, unresolved int
	for _, report := range idx.Reports() {
		for _, c := range report.Citations {
			checked++
			err := resolveCitation(cfg, c.URI)
			if err == nil {
				continue
			}
			unresolved++
			fmt.Fprintf(out, "%s%s:%d cites %s: %v\n", repoPrefix(report.Repo), c.File, c.Line, c.URI, err)
		}
	}
	if unresolved > 0 {
		return fmt.Errorf("%d of %d cited livt URI(s) do not resolve against %s", unresolved, checked, root)
	}
	fmt.Fprintf(out, "%d cited livt URI(s) resolve\n", checked)
	return nil
}

// resolveCitation answers whether the livt repository holds what the URI names.
// A retired item does resolve and so passes: retiring is a decision taken on
// the spec's side, and failing a build over it would reach further than the
// decision does.
func resolveCitation(cfg mcp.Config, rawURI string) error {
	p, ok := uri.Parse(rawURI)
	if !ok {
		return errors.New("not a livt URI")
	}
	return cfg.Verify(p)
}

// repoPrefix names whose report a failing citation came from. A report
// generated outside a checkout names no repository, and a bare space before
// the file path would be the only trace left of it.
func repoPrefix(repo string) string {
	if repo == "" {
		return ""
	}
	return repo + " "
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
