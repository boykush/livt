package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "livt",
	Short: "Living Text - Collaborate on board. Make it living in text.",
}

func Execute() error {
	return rootCmd.Execute()
}
