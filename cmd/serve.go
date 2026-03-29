package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/boykush/livt/internal/server"
	"github.com/spf13/cobra"
)

var port int

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 3000, "port to listen on")
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a local server to view artifacts as sticky notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		mappingsDir := filepath.Join("discoveries", "example-mappings")
		fmt.Printf("Starting livt server on http://localhost:%d\n", port)
		return server.Start(port, mappingsDir)
	},
}
