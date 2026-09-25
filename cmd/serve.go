package cmd

import (
	"github.com/boykush/livt/internal/config"
	"github.com/boykush/livt/internal/server"
	"github.com/spf13/cobra"
)

var port int

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 3000, "port to listen on")
	serveCmd.Flags().StringVarP(&outDir, "out", "o", "dist", "output directory")
	serveCmd.Flags().StringVar(&diffRange, "diff", "", diffFlagUsage)
	serveCmd.Flags().StringVar(&reportsDir, "reports", defaultReportsDir, reportsFlagUsage)
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Build and start a local server to view artifacts as sticky notes",
	Long: `Build the site as livt build does, serve it, and rebuild it — reloading the
open page — whenever one of its inputs or livt.yaml changes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := newBuilder(outDir)
		if err != nil {
			return err
		}
		return server.Serve(b, port, config.Path)
	},
}
