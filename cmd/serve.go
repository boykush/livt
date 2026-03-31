package cmd

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/boykush/livt/internal/builder"
	"github.com/spf13/cobra"
)

var port int

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 3000, "port to listen on")
	serveCmd.Flags().StringVarP(&outDir, "out", "o", "dist", "output directory")
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Build and start a local server to view artifacts as sticky notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		b := &builder.Builder{
			MappingsDir: filepath.Join("discoveries", "example-mappings"),
			StoriesDir:  "stories",
			USMDir:      filepath.Join("discoveries", "usm"),
			OutDir:      outDir,
		}
		fmt.Printf("Building to %s/\n", outDir)
		if err := b.Build(); err != nil {
			return err
		}

		fmt.Printf("Serving on http://localhost:%d\n", port)
		return http.ListenAndServe(fmt.Sprintf(":%d", port), http.FileServer(http.Dir(outDir)))
	},
}
