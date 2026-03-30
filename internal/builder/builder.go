package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/domain"
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
			StoryKey:  story.Key,
			StoryName: story.Name,
			Path:      "story/" + story.Key + ".html",
		})

		mappingPath := ""
		if b.hasExampleMapping(story.Key) {
			mappingPath = "../mapping/" + story.Key + ".html"
		}

		storyOutPath := filepath.Join(b.OutDir, "story", story.Key+".html")
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

func (b *Builder) hasExampleMapping(storyKey string) bool {
	path := filepath.Join(b.MappingsDir, storyKey+".yaml")
	_, err := os.Stat(path)
	return err == nil
}

func (b *Builder) buildMappings() error {
	files, err := filepath.Glob(filepath.Join(b.MappingsDir, "*.yaml"))
	if err != nil {
		return err
	}

	for _, f := range files {
		em, err := parser.ParseExampleMapping(f)
		if err != nil {
			return fmt.Errorf("parse %s: %w", f, err)
		}

		storyName := b.resolveStoryName(em.Story)

		outPath := filepath.Join(b.OutDir, "mapping", em.Story+".html")
		if err := b.buildMapping(outPath, em, storyName); err != nil {
			return err
		}
		fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))
	}

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

func (b *Builder) buildStory(path string, story *domain.Story, mappingPath string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderStory(f, story, mappingPath)
}

func (b *Builder) buildMapping(path string, em *domain.ExampleMapping, storyName string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderMapping(f, em, storyName)
}
