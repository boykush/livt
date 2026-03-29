package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of livt",
	Run: func(cmd *cobra.Command, args []string) {
		version := "(unknown)"
		if info, ok := debug.ReadBuildInfo(); ok {
			if info.Main.Version != "" {
				version = info.Main.Version
			}
		}
		fmt.Printf("livt version %s\n", version)
	},
}
