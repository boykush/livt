package cmd

import (
	"path/filepath"

	"github.com/boykush/livt/internal/opener"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(openCmd)
}

var openCmd = &cobra.Command{
	Use:   "open <story-key>",
	Short: "Open a story card for conversation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		o := &opener.Opener{
			StoriesDir: "stories",
			USMDir:     filepath.Join("discoveries", "usm"),
		}
		return o.Open(args[0])
	},
}
