package builder

import (
	"os"
	"path/filepath"
)

// buildMappingsIndex renders index.html: the Example Mappings overview, where
// each mapping appears as an iframe preview card.
func (b *Builder) buildMappingsIndex(tiles []mappingTile) error {
	sb, err := b.sidebar("example-mapping", "")
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(b.OutDir, "index.html"))
	if err != nil {
		return err
	}
	defer f.Close()
	return renderMappingsIndex(f, mappingsIndexView{Sidebar: sb, Mappings: tiles})
}

// buildStoryMapsIndex renders story-maps.html: the Story Maps overview, where
// each map appears as an iframe preview card.
func (b *Builder) buildStoryMapsIndex(tiles []storyMapTile) error {
	sb, err := b.sidebar("story-map", "")
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(b.OutDir, "story-maps.html"))
	if err != nil {
		return err
	}
	defer f.Close()
	return renderStoryMapsIndex(f, storyMapsIndexView{Sidebar: sb, StoryMaps: tiles})
}

// buildStoriesIndex renders stories.html: the Stories list.
func (b *Builder) buildStoriesIndex(items []storyItem) error {
	sb, err := b.sidebar("story", "")
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(b.OutDir, "stories.html"))
	if err != nil {
		return err
	}
	defer f.Close()
	return renderStoriesIndex(f, storiesIndexView{Sidebar: sb, Stories: items})
}
