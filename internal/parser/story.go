package parser

import (
	"path/filepath"

	"github.com/boykush/livt/internal/domain"
)

func ParseStory(path string) (*domain.Story, error) {
	key, name, body, err := parseMarkdownDoc(path)
	if err != nil {
		return nil, err
	}

	return &domain.Story{
		Key:  domain.StoryKey{Value: key},
		Name: name,
		Body: body,
	}, nil
}

func FindStoryByKey(storiesDir string, key domain.StoryKey) (*domain.Story, error) {
	path := filepath.Join(storiesDir, key.Value+".md")
	story, err := ParseStory(path)
	if err != nil {
		return &domain.Story{Key: key, Name: key.Value}, nil
	}
	return story, nil
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
