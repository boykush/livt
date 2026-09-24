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
	if o, ok := index[sm.OpportunityKey.Value]; ok {
		return opportunityRef{Name: o.DisplayName(), Path: "../" + uri.OpportunityPage(o.Key.Value)}
	}
	return opportunityRef{Name: sm.DisplayName(), Path: "../" + uri.StoryMapPage(sm.OpportunityKey.Value)}
}

// buildOpportunities builds a page per opportunity, and returns a preview tile
// for each on the Opportunities hub. mapByKey names the story map filed under an
// opportunity's key, so the opportunity links the journey drawn for it;
// storiesByKey and tallies are what its progress is summed from, which is why
// this runs after the mappings rather than beside the maps.
func (b *Builder) buildOpportunities(mapByKey map[string]storyMapRef, storiesByKey map[string][]opportunityReleaseStories, tallies map[string]mappingTally) ([]opportunityTile, error) {
	opportunities, err := parser.ParseAllOpportunities(b.OpportunitiesDir)
	if err != nil {
		return nil, err
	}

	var tiles []opportunityTile
	for _, o := range opportunities {
		canvasPath := ""
		if b.hasOpportunityCanvas(o.Key) {
			canvasPath = "../" + uri.OpportunityCanvasPage(o.Key.Value)
		}

		progress := b.progressOf(o, storiesByKey[o.Key.Value], tallies)
		progressPath := ""
		// An opportunity nobody has mapped a journey for has no progress to
		// read, so it gets no page — the same way one with no canvas links to
		// none, and for the same reason: the absence is the record.
		if progress.TotalStories > 0 {
			progressPath = "../" + uri.OpportunityProgressPage(o.Key.Value)
			if err := b.buildOpportunityProgress(o, progress); err != nil {
				return nil, err
			}
		}

		var storyMap *storyMapRef
		if ref, ok := mapByKey[o.Key.Value]; ok {
			storyMap = &ref
		}
		outPath := filepath.Join(b.OutDir, uri.OpportunityPage(o.Key.Value))
		if err := b.renderOpportunityPage(outPath, o, canvasPath, progressPath, storyMap, progress); err != nil {
			return nil, err
		}
		fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))

		tiles = append(tiles, opportunityTile{
			Key:       o.Key.Value,
			Name:      o.DisplayName(),
			Statement: o.Body,
			HasCanvas: canvasPath != "",
			StoryMap:  rootRelativeMap(storyMap),
			Links:     urlMetaFieldViews(o.Meta),
			Progress:  progress,
		})
	}

	return tiles, nil
}

func (b *Builder) renderOpportunityPage(path string, o *domain.Opportunity, canvasPath, progressPath string, storyMap *storyMapRef, progress opportunityProgress) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderOpportunity(f, b.Lang, opportunityView{
		Diff:         b.diffMark("../", uri.Opportunity(o.Key.Value)),
		Opportunity:  o,
		Meta:         metaFieldViews(o.Meta),
		CanvasPath:   canvasPath,
		ProgressPath: progressPath,
		StoryMap:     storyMap,
		Progress:     progress,
	})
}

// buildOpportunityProgress renders the story-by-story reading of how far an
// opportunity has been taken. The opportunity's own page carries the two gauges
// and leads here: what an opportunity *is* does not change between two visits,
// and this is the part that does.
func (b *Builder) buildOpportunityProgress(o *domain.Opportunity, progress opportunityProgress) error {
	outPath := filepath.Join(b.OutDir, uri.OpportunityProgressPage(o.Key.Value))
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := renderOpportunityProgress(f, b.Lang, opportunityProgressView{
		OpportunityKey:  o.Key.Value,
		OpportunityName: o.DisplayName(),
		Progress:        progress,
	}); err != nil {
		return err
	}
	fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))
	return nil
}

// opportunityProgress is how far an opportunity has been taken, on the two axes
// that move independently: how many of the stories it took on have been through
// an example mapping, and how many of the rules that produced are held by
// tests. Folding the two into one number would let a high automation ratio over
// a third of the stories read as nearly done.
type opportunityProgress struct {
	// Releases are the opportunity's stories in its map's own release slices,
	// which is the order a team heading for a release reads them in. A map that
	// declares no release leaves one unnamed slice holding every story.
	Releases      []opportunityReleaseRow
	TotalStories  int
	MappedStories int
	Rules         int
	Automated     int
	Proposed      int
	// Unautomated is what the Tasks page lists as waiting for a test: neither
	// automated nor proposed. Not derived from the three above, since a
	// proposal that already carries a test is in two of them.
	Unautomated    int
	Questions      int
	QuestionsAsked int
	// StoriesPath and TasksPath narrow the lists that already render these
	// items down to this opportunity, through the filter bar's own query
	// parameter — so a figure here is a way into the page that owns it rather
	// than a second place the same list is kept. Each Tasks link carries the
	// anchor of the list it counts, since that page keeps three and landing on
	// the first leaves the reader to find the other two.
	StoriesPath string
	// MappingsPath is the opportunity's boards, and belongs to no meter: a
	// reader who came to look at the mappings themselves is not answering a
	// figure, and hanging it off one would make it a claim about that count.
	MappingsPath    string
	QuestionsPath   string
	ProposedPath    string
	UnautomatedPath string
}

// opportunityReleaseRow is one release slice on the progress page, carrying its
// own coverage: a team heading for a release wants that slice's number, not the
// opportunity's.
type opportunityReleaseRow struct {
	ID            string
	Name          string
	Stories       []opportunityStoryRow
	MappedStories int
	Rules         int
	Automated     int
}

// Unscoped reports whether this is the remainder rather than a stage of the
// plan. The page names it only when there is a named slice to tell it from.
func (r opportunityReleaseRow) Unscoped() bool { return r.ID == "" }

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
func (b *Builder) progressOf(o *domain.Opportunity, slices []opportunityReleaseStories, tallies map[string]mappingTally) opportunityProgress {
	narrow := "?" + opportunityFilterParam + "=" + url.QueryEscape(o.DisplayName())
	tasks := "../tasks.html" + narrow + "#"
	p := opportunityProgress{
		StoriesPath:     "../stories.html" + narrow,
		MappingsPath:    "../index.html" + narrow,
		QuestionsPath:   tasks + tasksQuestionsAnchor,
		ProposedPath:    tasks + tasksProposedAnchor,
		UnautomatedPath: tasks + tasksRulesAnchor,
	}
	for _, slice := range slices {
		row := opportunityReleaseRow{ID: slice.ID, Name: slice.Name}
		for _, key := range slice.Keys {
			storyKey := domain.StoryKey{Value: key}
			t, mapped := tallies[key]
			story := opportunityStoryRow{
				Name:      b.resolveStoryName(storyKey),
				Mapped:    mapped,
				Rules:     t.Rules,
				Automated: t.Automated,
			}
			switch {
			case mapped:
				story.Path = "../" + uri.MappingPage(key)
				row.MappedStories++
				row.Rules += t.Rules
				row.Automated += t.Automated
				p.Proposed += t.Proposed
				p.Unautomated += t.Unautomated
				p.Questions += t.Questions
				p.QuestionsAsked += t.QuestionsAsked
			case b.hasStoryPage(storyKey):
				story.Path = "../" + uri.StoryPage(key)
			}
			row.Stories = append(row.Stories, story)
		}
		p.TotalStories += len(row.Stories)
		p.MappedStories += row.MappedStories
		p.Rules += row.Rules
		p.Automated += row.Automated
		p.Releases = append(p.Releases, row)
	}
	return p
}

// buildOpportunityCanvases renders a sheet for every canvas file, whether or not
// an opportunity file shares its key: the canvas sits beside its opportunity the
// way a mapping sits beside its story, and either stands without the other.
func (b *Builder) buildOpportunityCanvases(opportunities map[string]*domain.Opportunity) error {
	files, err := filepath.Glob(filepath.Join(b.CanvasesDir, "*.yaml"))
	if err != nil {
		return err
	}
	for _, f := range files {
		canvas, err := parser.ParseOpportunityCanvas(f)
		if err != nil {
			return fmt.Errorf("parse %s: %w", f, err)
		}
		if err := b.buildOpportunityCanvas(canvas, opportunities[canvas.OpportunityKey.Value]); err != nil {
			return err
		}
	}
	return nil
}

// buildOpportunityCanvas renders the ten-box board for one canvas. o is nil
// when no opportunity file shares the canvas's key; the key still names the
// opportunity, so the sheet is called by it, as a mapping with neither a name
// nor a story is called by its key.
func (b *Builder) buildOpportunityCanvas(canvas *domain.OpportunityCanvas, o *domain.Opportunity) error {
	key := canvas.OpportunityKey.Value
	view := opportunityCanvasView{
		Diff:            b.diffMark("../", uri.OpportunityCanvas(key)),
		OpportunityKey:  key,
		OpportunityName: key,
		Panels:          canvas.Panels(),
		Ubiquitous:      b.resolveTermCards(canvas.Ubiquitous),
	}
	if o != nil {
		view.OpportunityName = o.DisplayName()
		view.OpportunityPath = "../" + uri.OpportunityPage(key)
	}

	outPath := filepath.Join(b.OutDir, uri.OpportunityCanvasPage(key))
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := renderOpportunityCanvas(f, b.Lang, view); err != nil {
		return err
	}
	fmt.Printf("  %s\n", strings.TrimPrefix(outPath, b.OutDir+"/"))
	return nil
}
