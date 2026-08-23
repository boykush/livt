package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/boykush/livt/internal/builder"
	"github.com/boykush/livt/internal/config"
	"github.com/spf13/cobra"
)

var outDir string

func init() {
	buildCmd.Flags().StringVarP(&outDir, "out", "o", "dist", "output directory")
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
	return &builder.Builder{
		MappingsDir:   filepath.Join("discoveries", "example-mappings"),
		StoriesDir:    "stories",
		USMDir:        filepath.Join("discoveries", "usm"),
		UbiquitousDir: "ubiquitous",
		OutDir:        outDir,
		Lang:          cfg.Lang,
	}, nil
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
