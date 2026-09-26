package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "livt",
	Short: "Living Text - Your team's product decisions, living in text — context for AI.",
	Long: `Living Text - Your team's product decisions, living in text — context for AI.

Each command's --help says what it takes. Every other detail lives where it
cannot drift from what it describes: https://boykush.github.io/livt/reference.html
says where.`,
}

func Execute() error {
	return rootCmd.Execute()
}
