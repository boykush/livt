package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/parser"
)

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

		storyName := b.resolveStoryName(em.StoryKey)

		outPath := filepath.Join(b.OutDir, "mapping", em.StoryKey.Value+".html")
		if err := b.buildMapping(outPath, em, storyName); err != nil {
			return err
		}
		fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))
	}

	return nil
}

func (b *Builder) resolveStoryName(key domain.StoryKey) string {
	story, err := parser.FindStoryByKey(b.StoriesDir, key)
	if err != nil {
		return key.Value
	}
	return story.Name
}

func (b *Builder) buildMapping(path string, em *domain.ExampleMapping, storyName string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderMapping(f, em, storyName)
}
