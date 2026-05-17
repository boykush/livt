package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var storyKey string
var storyName string
var storiesDir string
var usmPath string
var candidateName string
var storyAs string
var storyWant string
var storySoThat string

var storyKeyPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

func init() {
	storyCommitCmd.Flags().StringVar(&storyKey, "key", "", "story key")
	storyCommitCmd.Flags().StringVar(&storyName, "name", "", "story display name")
	storyCommitCmd.Flags().StringVar(&storiesDir, "stories-dir", "stories", "stories directory")
	storyCommitCmd.Flags().StringVar(&usmPath, "usm", "", "story map YAML file")
	storyCommitCmd.Flags().StringVar(&candidateName, "candidate", "", "story candidate name in the story map")
	storyCommitCmd.Flags().StringVar(&storyAs, "as", "", "story persona")
	storyCommitCmd.Flags().StringVar(&storyWant, "want", "", "story goal")
	storyCommitCmd.Flags().StringVar(&storySoThat, "so-that", "", "story benefit")
	storyCmd.AddCommand(storyCommitCmd)
	rootCmd.AddCommand(storyCmd)
}

var storyCmd = &cobra.Command{
	Use:   "story",
	Short: "Work with story files",
}

var storyCommitCmd = &cobra.Command{
	Use:   "commit [key]",
	Short: "Commit a story candidate to detailed discovery",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := strings.TrimSpace(storyKey)
		if len(args) == 1 {
			if key != "" && key != args[0] {
				return fmt.Errorf("story key specified both as argument and --key")
			}
			key = strings.TrimSpace(args[0])
		}
		if key == "" {
			return fmt.Errorf("story key is required")
		}
		if !storyKeyPattern.MatchString(key) {
			return fmt.Errorf("story key must be kebab-case: %s", key)
		}

		if err := validateStoryBodyFlags(); err != nil {
			return err
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

		name := strings.TrimSpace(storyName)
		var updatedStoryMap []byte
		if usmPath != "" {
			if candidateName == "" {
				return fmt.Errorf("--candidate is required when --usm is set")
			}
			if name != "" {
				return fmt.Errorf("--name cannot be used with --usm; the story name is read from the story map")
			}

			var err error
			name, updatedStoryMap, err = prepareCandidateCommit(usmPath, candidateName, key)
			if err != nil {
				return err
			}
		}

		if name == "" {
			return fmt.Errorf("--name is required when --usm is not set")
		}

		content := storyMarkdown(name)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
		if updatedStoryMap != nil {
			if err := os.WriteFile(usmPath, updatedStoryMap, 0o644); err != nil {
				return err
			}
		}

		fmt.Printf("Committed %s\n", path)
		return nil
	},
}

func validateStoryBodyFlags() error {
	count := 0
	for _, value := range []string{storyAs, storyWant, storySoThat} {
		if value != "" {
			count++
		}
	}
	if count != 0 && count != 3 {
		return fmt.Errorf("--as, --want, and --so-that must be specified together")
	}
	return nil
}

func storyMarkdown(name string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "---\nname: %s\n---\n", name)
	if storyAs != "" {
		fmt.Fprintf(&b, "\nAs a %s\nI want %s\nSo that %s\n", storyAs, storyWant, storySoThat)
	}
	return b.String()
}

func prepareCandidateCommit(path, candidate, key string) (string, []byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil, err
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return "", nil, err
	}
	if len(doc.Content) == 0 {
		return "", nil, fmt.Errorf("empty story map: %s", path)
	}

	var matches []*yaml.Node
	activities := mappingValue(doc.Content[0], "activities")
	if activities != nil {
		for _, activity := range activities.Content {
			steps := mappingValue(activity, "steps")
			if steps == nil {
				continue
			}
			for _, step := range steps.Content {
				stories := mappingValue(step, "stories")
				if stories == nil {
					continue
				}
				for _, story := range stories.Content {
					nameNode := mappingValue(story, "name")
					if nameNode != nil && nameNode.Value == candidate {
						matches = append(matches, story)
					}
				}
			}
		}
	}

	if len(matches) == 0 {
		return "", nil, fmt.Errorf("story candidate %q not found in %s", candidate, path)
	}
	if len(matches) > 1 {
		return "", nil, fmt.Errorf("story candidate %q is not unique in %s", candidate, path)
	}

	story := matches[0]
	if existing := mappingValue(story, "key"); existing != nil && existing.Value != "" {
		return "", nil, fmt.Errorf("story candidate %q already has key %q", candidate, existing.Value)
	}
	setMappingValue(story, "key", key)

	var out bytes.Buffer
	encoder := yaml.NewEncoder(&out)
	encoder.SetIndent(2)
	if err := encoder.Encode(&doc); err != nil {
		return "", nil, err
	}
	if err := encoder.Close(); err != nil {
		return "", nil, err
	}

	return candidate, out.Bytes(), nil
}

func mappingValue(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(node.Content)-1; i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}

func setMappingValue(node *yaml.Node, key, value string) {
	if existing := mappingValue(node, key); existing != nil {
		existing.Value = value
		return
	}
	node.Content = append(node.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key},
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value},
	)
}
