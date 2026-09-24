package cmd

import (
	"encoding/json"
	"fmt"
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
)

func init() {
	automationsCmd.Flags().StringVar(&automationsRepo, "repo", "", "repository the report speaks for, owner/repo (default: read from the checkout's origin)")
	automationsCmd.Flags().StringVar(&automationsRev, "rev", "", "revision scanned (default: read from the checkout)")
	automationsCmd.Flags().StringVar(&automationsForge, "forge", "", "code host to build line URLs for: github or gitlab (default: inferred from origin, otherwise no URLs)")
	automationsCmd.Flags().StringVar(&automationsURLTemplate, "url-template", "", "line URL template for a host livt does not know, e.g. {base}/src/commit/{rev}/{path}#L{line}")
	automationsCmd.Flags().StringVarP(&automationsOut, "out", "o", "", "write the report to this file (default: stdout)")
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
