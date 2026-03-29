package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Story struct {
	Key  string
	Name string
}

func ParseStory(path string) (*Story, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := &Story{}
	scanner := bufio.NewScanner(f)
	inFrontmatter := false

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
				s.Key = strings.TrimSpace(strings.TrimPrefix(line, "key:"))
			}
			continue
		}

		if strings.HasPrefix(line, "# ") {
			s.Name = strings.TrimPrefix(line, "# ")
			break
		}
	}

	return s, scanner.Err()
}

func FindStoryByKey(storiesDir string, key string) (*Story, error) {
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

	return &Story{Key: key, Name: key}, nil
}
