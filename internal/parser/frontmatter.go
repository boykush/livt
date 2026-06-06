package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// parseMarkdownDoc reads a Markdown file with optional YAML frontmatter and
// returns the key (filename without extension), the frontmatter name field, and
// the trimmed body. It backs both story and term parsing, which share this
// name + body shape.
func parseMarkdownDoc(path string) (key, name, body string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", "", "", err
	}
	defer f.Close()

	key = strings.TrimSuffix(filepath.Base(path), ".md")

	var bodyLines []string
	scanner := bufio.NewScanner(f)
	inFrontmatter := false
	frontmatterDone := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "---" {
			if !frontmatterDone {
				inFrontmatter = !inFrontmatter
				if !inFrontmatter {
					frontmatterDone = true
				}
				continue
			}
		}

		if inFrontmatter {
			if strings.HasPrefix(line, "name:") {
				name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
			}
			continue
		}

		bodyLines = append(bodyLines, line)
	}

	if err := scanner.Err(); err != nil {
		return "", "", "", err
	}

	body = strings.TrimSpace(strings.Join(bodyLines, "\n"))
	return key, name, body, nil
}
