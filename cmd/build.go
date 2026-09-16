package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/boykush/livt/internal/builder"
	"github.com/boykush/livt/internal/config"
	"github.com/boykush/livt/internal/diff"
	"github.com/spf13/cobra"
)

var outDir string
var diffRange string

// diffFlagUsage is worded once and set on both commands, so `livt serve --diff`
// cannot come to mean something `livt build --diff` does not.
const diffFlagUsage = "render a diff between two revisions: <base>..<head>, or <base> alone against the working tree"

func init() {
	buildCmd.Flags().StringVarP(&outDir, "out", "o", "dist", "output directory")
	buildCmd.Flags().StringVar(&diffRange, "diff", "", diffFlagUsage)
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
		OutDir:           outDir,
		Lang:             cfg.Lang,
		Diff:             revisions,
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
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := newBuilder(outDir)
		if err != nil {
			return err
		}
		fmt.Printf("Building to %s/\n", outDir)
		return b.Build()
	},
}
