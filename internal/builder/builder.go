package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/parser"
)

type Builder struct {
	MappingsDir   string
	StoriesDir    string
	USMDir        string
	UbiquitousDir string
	OutDir        string
}

type sidebarCounts struct {
	mappings  int
	storyMaps int
	stories   int
	terms     int
}

// computeCounts tallies every resource type so the shared sidebar can show
// per-type counts on each hub page.
func (b *Builder) computeCounts() (sidebarCounts, error) {
	maps, err := parser.ParseAllStoryMaps(b.USMDir)
	if err != nil {
		return sidebarCounts{}, err
	}
	stories, err := parser.ParseAllStories(b.StoriesDir)
	if err != nil {
		return sidebarCounts{}, err
	}
	mappingFiles, err := filepath.Glob(filepath.Join(b.MappingsDir, "*.yaml"))
	if err != nil {
		return sidebarCounts{}, err
	}
	terms, err := parser.ParseAllTerms(b.UbiquitousDir)
	if err != nil {
		return sidebarCounts{}, err
	}
	return sidebarCounts{
		mappings:  len(mappingFiles),
		storyMaps: len(maps),
		stories:   len(stories),
		terms:     len(terms),
	}, nil
}

// sidebar builds the shared navigation data for a hub page. prefix is the
// relative path back to the output root, active marks the current resource type.
func (b *Builder) sidebar(active, prefix string) (Sidebar, error) {
	c, err := b.computeCounts()
	if err != nil {
		return Sidebar{}, err
	}
	return Sidebar{
		Prefix:    prefix,
		Active:    active,
		Mappings:  c.mappings,
		StoryMaps: c.storyMaps,
		Stories:   c.stories,
		Terms:     c.terms,
	}, nil
}

func (b *Builder) Build() error {
	for _, d := range []string{"story", "mapping", "story-map"} {
		if err := os.MkdirAll(filepath.Join(b.OutDir, d), 0o755); err != nil {
			return err
		}
	}

	storyToMaps, storyMapTiles, err := b.buildStoryMaps()
	if err != nil {
		return err
	}

	stories, err := parser.ParseAllStories(b.StoriesDir)
	if err != nil {
		return err
	}

	var storyItems []storyItem
	var storyOpportunitySets [][]opportunityRef
	for _, story := range stories {
		mappingPath := ""
		if b.hasExampleMapping(story.Key) {
			mappingPath = "../mapping/" + story.Key.Value + ".html"
		}

		// Story-page links resolve from the story/ directory ("../story-map/..").
		opportunities := storyToMaps[story.Key.Value]

		storyOutPath := filepath.Join(b.OutDir, "story", story.Key.Value+".html")
		if err := b.buildStory(storyOutPath, story, mappingPath, opportunities); err != nil {
			return err
		}
		fmt.Printf("  %s\n", strings.TrimPrefix(storyOutPath, b.OutDir+"/"))

		// List links resolve from the output root, not the story/ directory.
		listMappingPath := ""
		if mappingPath != "" {
			listMappingPath = "mapping/" + story.Key.Value + ".html"
		}
		listOpportunities := rootRelativeOpportunities(opportunities)
		storyOpportunitySets = append(storyOpportunitySets, listOpportunities)
		storyItems = append(storyItems, storyItem{
			Key:           story.Key.Value,
			Name:          story.Name,
			Opportunities: listOpportunities,
			MappingPath:   listMappingPath,
			Links:         urlMetaFieldViews(story.Meta),
		})
	}

	mappingTiles, err := b.buildMappings()
	if err != nil {
		return err
	}
	// A mapping belongs to the same opportunities as its story (mapping → story
	// → map), so the Example Mappings list filters on the same axis.
	var mappingOpportunitySets [][]opportunityRef
	for i := range mappingTiles {
		mappingTiles[i].Opportunities = rootRelativeOpportunities(storyToMaps[mappingTiles[i].Key])
		mappingOpportunitySets = append(mappingOpportunitySets, mappingTiles[i].Opportunities)
	}

	if err := b.buildGlossary(); err != nil {
		return err
	}

	// Hub pages share the sidebar; index.html is the Example Mappings overview.
	if err := b.buildMappingsIndex(mappingTiles, distinctOpportunityNames(mappingOpportunitySets)); err != nil {
		return err
	}
	fmt.Printf("  index.html\n")

	if err := b.buildStoryMapsIndex(storyMapTiles); err != nil {
		return err
	}
	fmt.Printf("  story-maps.html\n")

	if err := b.buildStoriesIndex(storyItems, distinctOpportunityNames(storyOpportunitySets)); err != nil {
		return err
	}
	fmt.Printf("  stories.html\n")

	return nil
}
