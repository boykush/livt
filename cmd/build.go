package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/boykush/livt/internal/builder"
	"github.com/spf13/cobra"
)

var outDir string

func init() {
	buildCmd.Flags().StringVarP(&outDir, "out", "o", "dist", "output directory")
	rootCmd.AddCommand(buildCmd)
}

// newBuilder lays out the livt repository's input directories. build and serve
// share it so the two cannot drift into reading different places.
func newBuilder(outDir string) *builder.Builder {
	return &builder.Builder{
		OpportunitiesDir: "opportunities",
		CanvasesDir:      filepath.Join("discoveries", "opportunity-canvases"),
		MappingsDir:      filepath.Join("discoveries", "example-mappings"),
		StoriesDir:       "stories",
		USMDir:           filepath.Join("discoveries", "usm"),
		UbiquitousDir:    "ubiquitous",
		OutDir:           outDir,
	}
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build static HTML from artifacts",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Building to %s/\n", outDir)
		return newBuilder(outDir).Build()
	},
}
