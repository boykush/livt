package parser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"gopkg.in/yaml.v3"
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

// parseNamedFrontmatter pulls out the reserved name field and keeps every other
// key as an ordered MetaField. It decodes into a yaml.Node so the key order from
// the source file is preserved — a map would lose it. Stories and opportunities
// share the shape, so they share this.
func parseNamedFrontmatter(frontmatter string) (name string, meta []domain.MetaField, err error) {
	if strings.TrimSpace(frontmatter) == "" {
		return "", nil, nil
	}

	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(frontmatter), &doc); err != nil {
		return "", nil, err
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return "", nil, nil
	}

	root := doc.Content[0]
	for i := 0; i+1 < len(root.Content); i += 2 {
		key := root.Content[i].Value
		value := root.Content[i+1]
		if key == "name" {
			name = value.Value
			continue
		}
		meta = append(meta, domain.MetaField{Key: key, Value: scalarOrJoin(value)})
	}
	return name, meta, nil
}

// scalarOrJoin renders a frontmatter value as a display string: scalars pass
// through, sequences of scalars are comma-joined, and anything else falls back
// to its scalar Value (empty for nested mappings).
func scalarOrJoin(n *yaml.Node) string {
	if n.Kind == yaml.SequenceNode {
		parts := make([]string, 0, len(n.Content))
		for _, item := range n.Content {
			parts = append(parts, item.Value)
		}
		return strings.Join(parts, ", ")
	}
	return n.Value
}
