package builder

import (
	"fmt"
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
// for it.
func (b *Builder) buildOpportunities(mapsByKey map[string][]storyMapRef) ([]opportunityTile, error) {
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

		outPath := filepath.Join(b.OutDir, uri.OpportunityPage(o.Key.Value))
		if err := b.renderOpportunityPage(outPath, o, canvasPath, mapsByKey[o.Key.Value]); err != nil {
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
		})
	}

	return tiles, nil
}

func (b *Builder) renderOpportunityPage(path string, o *domain.Opportunity, canvasPath string, maps []storyMapRef) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderOpportunity(f, opportunityView{
		Opportunity: o,
		Meta:        metaFieldViews(o.Meta),
		CanvasPath:  canvasPath,
		StoryMaps:   maps,
	})
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
	if err := renderOpportunityCanvas(f, opportunityCanvasView{
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
