package cmd

import (
	"path/filepath"

	"github.com/boykush/livt/internal/recorder"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(recordCmd)
}

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Record stories from USM as markdown files",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := &recorder.Recorder{
			StoriesDir: "stories",
			USMDir:     filepath.Join("discoveries", "usm"),
		}
		return r.Record()
	},
}
