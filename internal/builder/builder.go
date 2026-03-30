package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/parser"
)

type Builder struct {
	MappingsDir string
	StoriesDir  string
	OutDir      string
}

func (b *Builder) Build() error {
	if err := os.MkdirAll(filepath.Join(b.OutDir, "story"), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(b.OutDir, "mapping"), 0o755); err != nil {
		return err
	}

	stories, err := parser.ParseAllStories(b.StoriesDir)
	if err != nil {
		return err
	}

	var entries []IndexEntry
	for _, story := range stories {
		entries = append(entries, IndexEntry{
			StoryKey:  story.Key.Value,
			StoryName: story.Name,
			Path:      "story/" + story.Key.Value + ".html",
		})

		mappingPath := ""
		if b.hasExampleMapping(story.Key) {
			mappingPath = "../mapping/" + story.Key.Value + ".html"
		}

		storyOutPath := filepath.Join(b.OutDir, "story", story.Key.Value+".html")
		if err := b.buildStory(storyOutPath, story, mappingPath); err != nil {
			return err
		}
		fmt.Printf("  %s\n", strings.TrimPrefix(storyOutPath, b.OutDir+"/"))
	}

	if err := b.buildMappings(); err != nil {
		return err
	}

	indexPath := filepath.Join(b.OutDir, "index.html")
	if err := b.buildIndex(indexPath, entries); err != nil {
		return err
	}
	fmt.Printf("  index.html\n")

	return nil
}
