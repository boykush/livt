package parser

import (
	"path/filepath"

	"github.com/boykush/livt/internal/domain"
)

func ParseStory(path string) (*domain.Story, error) {
	key, frontmatter, body, err := splitFrontmatter(path)
	if err != nil {
		return nil, err
	}

	name, meta, err := parseNamedFrontmatter(frontmatter)
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
