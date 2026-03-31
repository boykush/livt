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

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build static HTML from artifacts",
	RunE: func(cmd *cobra.Command, args []string) error {
		b := &builder.Builder{
			MappingsDir: filepath.Join("discoveries", "example-mappings"),
			StoriesDir:  "stories",
			USMDir:      filepath.Join("discoveries", "usm"),
			OutDir:      outDir,
		}
		fmt.Printf("Building to %s/\n", outDir)
		return b.Build()
	},
}
