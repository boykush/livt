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

	key := strings.TrimSuffix(filepath.Base(path), ".md")

	var name string
	var bodyLines []string
	scanner := bufio.NewScanner(f)
	inFrontmatter := false
	pastTitle := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "---" {
			inFrontmatter = !inFrontmatter
			continue
		}

		if inFrontmatter {
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
