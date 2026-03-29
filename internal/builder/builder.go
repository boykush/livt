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
	if err := os.MkdirAll(filepath.Join(b.OutDir, "mapping"), 0o755); err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(b.MappingsDir, "*.yaml"))
	if err != nil {
		return err
	}

	var entries []IndexEntry
	for _, f := range files {
		em, err := parser.ParseExampleMapping(f)
		if err != nil {
			return fmt.Errorf("parse %s: %w", f, err)
		}

		storyName := b.resolveStoryName(em.Story)
		entries = append(entries, IndexEntry{
			StoryKey:  em.Story,
			StoryName: storyName,
			Path:      "mapping/" + em.Story + ".html",
		})

		outPath := filepath.Join(b.OutDir, "mapping", em.Story+".html")
		if err := b.buildMapping(outPath, em, storyName); err != nil {
			return err
		}

		fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))
	}

	indexPath := filepath.Join(b.OutDir, "index.html")
	if err := b.buildIndex(indexPath, entries); err != nil {
		return err
	}
	fmt.Printf("  index.html\n")

	return nil
}

func (b *Builder) resolveStoryName(key string) string {
	story, err := parser.FindStoryByKey(b.StoriesDir, key)
	if err != nil {
		return key
	}
	return story.Name
}

func (b *Builder) buildIndex(path string, entries []IndexEntry) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderIndex(f, entries)
}

func (b *Builder) buildMapping(path string, em *parser.ExampleMapping, storyName string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderMapping(f, em, storyName)
}
