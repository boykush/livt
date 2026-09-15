package builder

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/parser"
	"github.com/boykush/livt/internal/uri"
)

func (b *Builder) hasOpportunityCanvas(key domain.OpportunityKey) bool {
	_, err := os.Stat(filepath.Join(b.CanvasesDir, key.Value+".yaml"))
	return err == nil
}

// opportunityIndex reads the committed opportunities into a lookup by key. A
// story map resolves the opportunity it serves through this, so the directory
// is read once per build rather than once per map.
func (b *Builder) opportunityIndex() (map[string]*domain.Opportunity, error) {
	all, err := parser.ParseAllOpportunities(b.OpportunitiesDir)
	if err != nil {
		return nil, err
	}
	index := make(map[string]*domain.Opportunity, len(all))
	for _, o := range all {
		index[o.Key.Value] = o
	}
	return index, nil
}

// mapOpportunity is the opportunity a story map serves, as a ref relative to the
// story/ directory (the deepest page that renders one; the lists rebase it onto
// the output root). A map whose key names no committed opportunity falls back to
// standing in as its own opportunity, named by the map — which is what every
// opportunity chip meant before opportunities became files of their own.
func mapOpportunity(sm *domain.StoryMap, index map[string]*domain.Opportunity) opportunityRef {
	if o, ok := index[sm.Key]; ok {
		return opportunityRef{Name: o.DisplayName(), Path: "../" + uri.OpportunityPage(o.Key.Value)}
	}
	return opportunityRef{Name: sm.Name, Path: "../" + uri.StoryMapPage(sm.Name)}
}

// buildOpportunities builds a page per opportunity and per canvas, and returns a
// preview tile for each on the Opportunities hub. mapsByKey names the story maps
// that share an opportunity's key, so the opportunity links the journey mapped
// for it; storiesByKey and tallies are what its progress is summed from, which
// is why this runs after the mappings rather than beside the maps.
func (b *Builder) buildOpportunities(mapsByKey map[string][]storyMapRef, storiesByKey map[string][]string, tallies map[string]mappingTally) ([]opportunityTile, error) {
	opportunities, err := parser.ParseAllOpportunities(b.OpportunitiesDir)
	if err != nil {
		return nil, err
	}

	var tiles []opportunityTile
	for _, o := range opportunities {
		canvasPath := ""
		if b.hasOpportunityCanvas(o.Key) {
			canvasPath = "../" + uri.OpportunityCanvasPage(o.Key.Value)
			if err := b.buildOpportunityCanvas(o); err != nil {
				return nil, err
			}
		}

		progress := b.progressOf(o, storiesByKey[o.Key.Value], tallies)

		outPath := filepath.Join(b.OutDir, uri.OpportunityPage(o.Key.Value))
		if err := b.renderOpportunityPage(outPath, o, canvasPath, mapsByKey[o.Key.Value], progress); err != nil {
			return nil, err
		}
		fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))

		tiles = append(tiles, opportunityTile{
			Key:       o.Key.Value,
			Name:      o.DisplayName(),
			Statement: o.Body,
			HasCanvas: canvasPath != "",
			StoryMaps: rootRelativeMaps(mapsByKey[o.Key.Value]),
			Links:     urlMetaFieldViews(o.Meta),
			Progress:  progress,
		})
	}

	return tiles, nil
}

func (b *Builder) renderOpportunityPage(path string, o *domain.Opportunity, canvasPath string, maps []storyMapRef, progress opportunityProgress) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderOpportunity(f, b.Lang, opportunityView{
		Opportunity: o,
		Meta:        metaFieldViews(o.Meta),
		CanvasPath:  canvasPath,
		StoryMaps:   maps,
		Progress:    progress,
	})
}

// opportunityProgress is how far an opportunity has been taken, on the two axes
// that move independently: how many of the stories it took on have been through
// an example mapping, and how many of the rules that produced are held by
// tests. Folding the two into one number would let a high automation ratio over
// a third of the stories read as nearly done.
type opportunityProgress struct {
	Stories       []opportunityStoryRow
	MappedStories int
	Rules         int
	Automated     int
	Proposed      int
	Questions     int
	// StoriesPath and TasksPath narrow the lists that already render these
	// items down to this opportunity, through the filter bar's own query
	// parameter — so a figure here is a way into the page that owns it rather
	// than a second place the same list is kept.
	StoriesPath string
	TasksPath   string
}

// opportunityStoryRow is one story the opportunity took on. Path is the story's
// example mapping once it has one and its story page until then, because that
// is where a reader can act on the row in either state.
type opportunityStoryRow struct {
	Name      string
	Path      string
	Mapped    bool
	Rules     int
	Automated int
}

// progressOf sums one opportunity's stories. storyKeys is the opportunity's own
// keyed stories in map order, and tallies holds the counts of every mapping the
// build read; a key missing from it is a story whose conversation has not been
// held yet, which is the measurement rather than an absence of data.
func (b *Builder) progressOf(o *domain.Opportunity, storyKeys []string, tallies map[string]mappingTally) opportunityProgress {
	narrow := "?" + opportunityFilterParam + "=" + url.QueryEscape(o.DisplayName())
	p := opportunityProgress{
		StoriesPath: "../stories.html" + narrow,
		TasksPath:   "../tasks.html" + narrow,
	}
	for _, key := range storyKeys {
		storyKey := domain.StoryKey{Value: key}
		t, mapped := tallies[key]
		row := opportunityStoryRow{
			Name:      b.resolveStoryName(storyKey),
			Mapped:    mapped,
			Rules:     t.Rules,
			Automated: t.Automated,
		}
		switch {
		case mapped:
			row.Path = "../" + uri.MappingPage(key)
			p.MappedStories++
			p.Rules += t.Rules
			p.Automated += t.Automated
			p.Proposed += t.Proposed
			p.Questions += t.Questions
		case b.hasStoryPage(storyKey):
			row.Path = "../" + uri.StoryPage(key)
		}
		p.Stories = append(p.Stories, row)
	}
	return p
}

// buildOpportunityCanvas renders the ten-box board for one opportunity.
func (b *Builder) buildOpportunityCanvas(o *domain.Opportunity) error {
	canvas, err := parser.ParseOpportunityCanvas(filepath.Join(b.CanvasesDir, o.Key.Value+".yaml"))
	if err != nil {
		return fmt.Errorf("parse opportunity canvas %q: %w", o.Key.Value, err)
	}

	outPath := filepath.Join(b.OutDir, uri.OpportunityCanvasPage(o.Key.Value))
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := renderOpportunityCanvas(f, b.Lang, opportunityCanvasView{
		OpportunityKey:  o.Key.Value,
		OpportunityName: o.DisplayName(),
		Panels:          canvas.Panels(),
		Ubiquitous:      b.resolveTermCards(canvas.Ubiquitous),
	}); err != nil {
		return err
	}
	fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))
	return nil
}
