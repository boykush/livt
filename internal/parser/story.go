package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/domain"
)

func ParseStory(path string) (*domain.Story, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var key, name string
	var bodyLines []string
	scanner := bufio.NewScanner(f)
	inFrontmatter := false
	pastTitle := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			inFrontmatter = false
			continue
		}

		if inFrontmatter {
			if strings.HasPrefix(line, "key:") {
				key = strings.TrimSpace(strings.TrimPrefix(line, "key:"))
			}
			continue
		}

		if !pastTitle {
			if strings.HasPrefix(line, "# ") {
				name = strings.TrimPrefix(line, "# ")
				pastTitle = true
			}
			continue
		}

		bodyLines = append(bodyLines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	body := strings.TrimSpace(strings.Join(bodyLines, "\n"))

	return &domain.Story{
		Key:  domain.StoryKey{Value: key},
		Name: name,
		Body: body,
	}, nil
}

func FindStoryByKey(storiesDir string, key domain.StoryKey) (*domain.Story, error) {
	files, err := filepath.Glob(filepath.Join(storiesDir, "*.md"))
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		story, err := ParseStory(f)
		if err != nil {
			continue
		}
		if story.Key == key {
			return story, nil
		}
	}

	return &domain.Story{Key: key, Name: key.Value}, nil
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
