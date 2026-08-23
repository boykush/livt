package builder

import (
	"os"
	"path/filepath"
)

// buildTasks renders tasks.html: what the livt repository leaves unfinished — open
// questions and un-automated rules. filterOpportunities are the opportunity axes
// both lists can be filtered by.
func (b *Builder) buildTasks(questions, unautomatedRules []taskItem, filterOpportunities []string) error {
	sb, err := b.sidebar("task", "")
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(b.OutDir, "tasks.html"))
	if err != nil {
		return err
	}
	defer f.Close()
	return renderTasks(f, tasksView{Sidebar: sb, Questions: questions, UnautomatedRules: unautomatedRules, FilterOpportunities: filterOpportunities})
}

// buildMappingsIndex renders index.html: the Example Mappings overview and the
// site's landing page, where each mapping appears as an iframe preview card.
// filterOpportunities are the opportunity axes the list can be filtered by.
func (b *Builder) buildMappingsIndex(tiles []mappingTile, filterOpportunities []string) error {
	sb, err := b.sidebar("example-mapping", "")
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(b.OutDir, "index.html"))
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

// buildOpportunitiesIndex renders opportunities.html: the Opportunities hub,
// where each opportunity shows its statement and — once a canvas has been
// filled in for it — a preview of that board.
func (b *Builder) buildOpportunitiesIndex(tiles []opportunityTile) error {
	sb, err := b.sidebar("opportunity", "")
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(b.OutDir, "opportunities.html"))
	if err != nil {
		return err
	}
	defer f.Close()
	return renderOpportunitiesIndex(f, opportunitiesIndexView{Sidebar: sb, Opportunities: tiles})
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
