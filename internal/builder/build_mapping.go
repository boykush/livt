package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/parser"
)

// buildMappings builds example mapping HTML pages and returns a preview tile per
// mapping for the Example Mappings overview page (index.html).
func (b *Builder) buildMappings() ([]mappingTile, error) {
	files, err := filepath.Glob(filepath.Join(b.MappingsDir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	var tiles []mappingTile
	for _, f := range files {
		em, err := parser.ParseExampleMapping(f)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", f, err)
		}

		storyName := b.resolveStoryName(em.StoryKey)
		storyPath := ""
		if b.hasStoryPage(em.StoryKey) {
			storyPath = "../story/" + em.StoryKey.Value + ".html"
		}

		ubiquitous := b.resolveTermCards(em.Ubiquitous)

		outPath := filepath.Join(b.OutDir, "mapping", em.StoryKey.Value+".html")
		if err := b.buildMapping(outPath, em, storyName, storyPath, ubiquitous); err != nil {
			return nil, err
		}
		fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))

		tiles = append(tiles, mappingTile{Key: em.StoryKey.Value, StoryName: storyName})
	}

	return tiles, nil
}

func (b *Builder) resolveStoryName(key domain.StoryKey) string {
	story, err := parser.FindStoryByKey(b.StoriesDir, key)
	if err != nil {
		return key.Value
	}
	return story.Name
}

func (b *Builder) buildMapping(path string, em *domain.ExampleMapping, storyName, storyPath string, ubiquitous []termCard) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderMapping(f, em, storyName, storyPath, ubiquitous)
}
