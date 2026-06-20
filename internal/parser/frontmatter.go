package parser

import (
	"os"
	"path/filepath"
	"strings"
)

// splitFrontmatter reads a Markdown file and separates an optional leading YAML
// frontmatter block (delimited by --- fences) from the body. It returns the key
// (filename without extension), the raw frontmatter (the text between the
// fences, empty when there is none), and the trimmed body. Story and term
// parsing share this shape, then each unmarshals the frontmatter into its own
// schema.
func splitFrontmatter(path string) (key, frontmatter, body string, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", "", "", err
	}

	key = strings.TrimSuffix(filepath.Base(path), ".md")
	lines := strings.Split(string(data), "\n")

	// Frontmatter is only recognized when the file opens with a --- fence.
	if len(lines) > 0 && strings.TrimRight(lines[0], "\r") == "---" {
		for i := 1; i < len(lines); i++ {
			if strings.TrimRight(lines[i], "\r") == "---" {
				frontmatter = strings.Join(lines[1:i], "\n")
				body = strings.TrimSpace(strings.Join(lines[i+1:], "\n"))
				return key, frontmatter, body, nil
			}
		}
		// Unterminated fence: treat the remainder as frontmatter, no body.
		return key, strings.Join(lines[1:], "\n"), "", nil
	}

	return key, "", strings.TrimSpace(string(data)), nil
}
