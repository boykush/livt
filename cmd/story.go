package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var storyName string
var storiesDir string

func init() {
	storyInitCmd.Flags().StringVar(&storyName, "name", "", "story display name")
	storyInitCmd.Flags().StringVar(&storiesDir, "stories-dir", "stories", "stories directory")
	storyCmd.AddCommand(storyInitCmd)
	rootCmd.AddCommand(storyCmd)
}

var storyCmd = &cobra.Command{
	Use:   "story",
	Short: "Work with story files",
}

var storyInitCmd = &cobra.Command{
	Use:   "init <key>",
	Short: "Initialize a story file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := strings.TrimSpace(args[0])
		if key == "" {
			return fmt.Errorf("story key is required")
		}
		if storyName == "" {
			return fmt.Errorf("--name is required")
		}

		if err := os.MkdirAll(storiesDir, 0o755); err != nil {
			return err
		}

		path := filepath.Join(storiesDir, key+".md")
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("story file already exists: %s", path)
		} else if !os.IsNotExist(err) {
			return err
		}

		content := fmt.Sprintf("---\nname: %s\n---\n", storyName)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}

		fmt.Printf("Initialized %s\n", path)
		return nil
	},
}
