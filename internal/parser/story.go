package parser

import (
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"gopkg.in/yaml.v3"
)

func ParseStory(path string) (*domain.Story, error) {
	key, frontmatter, body, err := splitFrontmatter(path)
	if err != nil {
		return nil, err
	}

	name, meta, err := parseStoryFrontmatter(frontmatter)
	if err != nil {
		return nil, err
	}

	return &domain.Story{
		Key:  domain.StoryKey{Value: key},
		Name: name,
		Body: body,
		Meta: meta,
	}, nil
}

// parseStoryFrontmatter pulls out the reserved name field and keeps every other
// key as an ordered MetaField. It decodes into a yaml.Node so the key order from
// the source file is preserved — a map would lose it.
func parseStoryFrontmatter(frontmatter string) (name string, meta []domain.MetaField, err error) {
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

// FindStoryByKey reads the story a mapping names, falling back to a bare keyed
// story when there is none to read: a mapping's board stands on its own, so a
// story page that is missing or unreadable must not stop it rendering. Telling
// those two apart is a different job — mcp.Config.story does it, for a URI that
// names the story itself.
func FindStoryByKey(storiesDir string, key domain.StoryKey) *domain.Story {
	path := filepath.Join(storiesDir, key.Value+".md")
	story, err := ParseStory(path)
	if err != nil {
		return &domain.Story{Key: key}
	}
	return story
}

func ParseAllStories(storiesDir string) ([]*domain.Story, error) {
	files, err := filepath.Glob(filepath.Join(storiesDir, "*.md"))
	if err != nil {
		return nil, err
	}

	var stories []*domain.Story
	for _, f := range files {
		story, err := ParseStory(f)
		if err != nil {
			return nil, err
		}
		stories = append(stories, story)
	}

	return stories, nil
}
