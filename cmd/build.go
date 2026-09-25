package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/boykush/livt/internal/builder"
	"github.com/boykush/livt/internal/config"
	"github.com/boykush/livt/internal/diff"
	"github.com/boykush/livt/internal/gitrev"
	"github.com/boykush/livt/internal/i18n"
	"github.com/spf13/cobra"
)

var outDir string
var diffRange string
var reportsDir string

// diffFlagUsage is worded once and set on both commands, so `livt serve --diff`
// cannot come to mean something `livt build --diff` does not.
const diffFlagUsage = "render a diff between two revisions: <base>..<head>, or <base> alone against the working tree"

// The reports flag is worded and defaulted once for the same reason as the
// diff flag: build and serve must read the same place.
const (
	reportsFlagUsage  = "directory of collected automation reports"
	defaultReportsDir = "automations"
)

func init() {
	buildCmd.Flags().StringVarP(&outDir, "out", "o", "dist", "output directory")
	buildCmd.Flags().StringVar(&diffRange, "diff", "", diffFlagUsage)
	buildCmd.Flags().StringVar(&reportsDir, "reports", defaultReportsDir, reportsFlagUsage)
	rootCmd.AddCommand(buildCmd)
}

// newBuilder lays out the livt repository's input directories and reads
// livt.yaml beside them. build and serve share it so the two cannot drift into
// reading different places, or into building the site in different languages.
func newBuilder(outDir string) (*builder.Builder, error) {
	cfg, err := config.Load(config.Path)
	if err != nil {
		return nil, err
	}
	revisions, err := diffRangeFlag()
	if err != nil {
		return nil, err
	}
	return &builder.Builder{
		OpportunitiesDir: "opportunities",
		CanvasesDir:      filepath.Join("discoveries", "opportunity-canvases"),
		MappingsDir:      filepath.Join("discoveries", "example-mappings"),
		StoriesDir:       "stories",
		USMDir:           filepath.Join("discoveries", "usm"),
		UbiquitousDir:    "ubiquitous",
		AutomationsDir:   reportsDir,
		OutDir:           outDir,
		Lang:             cfg.Lang,
		Diff:             revisions,
		// What the site foots with. The binary knows its own version, and the
		// livt repository is the working directory the paths above are
		// relative to.
		LivtVersion: livtVersion(),
		SpecVersion: gitrev.Short("."),
	}, nil
}

// diffRangeFlag reads --diff, which is absent far more often than not: a build
// given no revisions renders the site it always did, and returning nil is what
// says so.
func diffRangeFlag() (*diff.Range, error) {
	if diffRange == "" {
		return nil, nil
	}
	parsed, err := diff.ParseRange(diffRange)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build static HTML from artifacts",
	// The languages are read from the catalogs, so the help cannot offer one
	// that livt.yaml would refuse.
	Long: fmt.Sprintf(`Build the livt repository in the current directory as a static site.

An optional livt.yaml beside it sets lang, the language of the labels livt
adds to the site. What the livt repository itself says is rendered as
written. lang is %s unless set, and one of: %s.`, i18n.Default, i18n.List()),
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := newBuilder(outDir)
		if err != nil {
			return err
		}
		fmt.Printf("Building to %s/\n", outDir)
		return b.Build()
	},
}
