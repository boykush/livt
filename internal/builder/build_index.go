package builder

import (
	"os"
	"path/filepath"
)

// buildHome renders index.html: the landing page, listing what the master
// leaves unfinished — open questions and un-automated rules.
// filterOpportunities are the opportunity axes both lists can be filtered by.
func (b *Builder) buildHome(questions, unautomatedRules []outstandingItem, filterOpportunities []string) error {
	sb, err := b.sidebar("home", "")
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(b.OutDir, "index.html"))
	if err != nil {
		return err
	}
	defer f.Close()
	return renderHome(f, homeView{Sidebar: sb, Questions: questions, UnautomatedRules: unautomatedRules, FilterOpportunities: filterOpportunities})
}

// buildMappingsIndex renders example-mappings.html: the Example Mappings
// overview, where each mapping appears as an iframe preview card.
// filterOpportunities are the opportunity axes the list can be filtered by.
func (b *Builder) buildMappingsIndex(tiles []mappingTile, filterOpportunities []string) error {
	sb, err := b.sidebar("example-mapping", "")
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(b.OutDir, "example-mappings.html"))
	if err != nil {
		return err
	}
	defer f.Close()
	return renderMappingsIndex(f, mappingsIndexView{Sidebar: sb, Mappings: tiles, FilterOpportunities: filterOpportunities})
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

// buildStoriesIndex renders stories.html: the Stories list. filterOpportunities
// are the opportunity axes the list can be filtered by.
func (b *Builder) buildStoriesIndex(items []storyItem, filterOpportunities []string) error {
	sb, err := b.sidebar("story", "")
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(b.OutDir, "stories.html"))
	if err != nil {
		return err
	}
	defer f.Close()
	return renderStoriesIndex(f, storiesIndexView{Sidebar: sb, Stories: items, FilterOpportunities: filterOpportunities})
}
