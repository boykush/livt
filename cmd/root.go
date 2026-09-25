package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "livt",
	Short: "Living Text - Collaborate on board. Make it living in text.",
	Long: `Living Text - Collaborate on board. Make it living in text.

Each command's --help says what it takes. Every other detail lives where it
cannot drift from what it describes: https://boykush.github.io/livt/reference.html
says where.`,
}

func Execute() error {
	return rootCmd.Execute()
}
