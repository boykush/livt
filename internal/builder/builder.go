package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/automation"
	"github.com/boykush/livt/internal/diff"
	"github.com/boykush/livt/internal/i18n"
	"github.com/boykush/livt/internal/parser"
)

type Builder struct {
	OpportunitiesDir string
	CanvasesDir      string
	MappingsDir      string
	StoriesDir       string
	USMDir           string
	UbiquitousDir    string
	// AutomationsDir holds the collected reports the implementation
	// repositories push. Absent is the ordinary state of a livt repository no
	// implementation points at, and the site then shows no automation at all.
	AutomationsDir string
	OutDir         string
	// Lang is the language of the site chrome, from livt.yaml. The zero value
	// renders in livt's default rather than failing, so a Builder built without
	// a config still produces a site.
	Lang i18n.Lang
	// Diff is the revision range to render a diff for, nil when the build was
	// given none. Root is the git repository those revisions are read from —
	// the directory the input directories are relative to.
	Diff *diff.Range
	Root string

	// diffResult is the range computed for the run in progress. `livt serve`
	// rebuilds on every edit and the working tree is usually the head, so the
	// diff is recomputed per build rather than held across them; the sidebar of
	// every page needs the count, which is why it is on the Builder at all.
	diffResult *diff.Result
	// diffByURI is the same result keyed for the lookup every resource page
	// makes: did this one change, and what became of it.
	diffByURI map[string]diff.Became
	// automations is the run's report index. `livt serve` rebuilds on every
	// edit, so it is reloaded per build rather than held across them.
	automations *automation.Index
}

// automationIndex loads the reports once per build. Every surface that asks
// what a rule is automated by goes through here, so the boards and the counts
// cannot end up reading different answers.
func (b *Builder) automationIndex() (*automation.Index, error) {
	if b.automations == nil {
		idx, err := automation.Load(b.AutomationsDir)
		if err != nil {
			return nil, err
		}
		b.automations = idx
	}
	return b.automations, nil
}

// diffDirs is the input layout the diff reads a revision through, which is this
// Builder's own: the site and the diff must not drift into reading different
// places.
func (b *Builder) diffDirs() diff.Dirs {
	return diff.Dirs{
		Opportunities: b.OpportunitiesDir,
		Canvases:      b.CanvasesDir,
		Mappings:      b.MappingsDir,
		Stories:       b.StoriesDir,
		USM:           b.USMDir,
		Ubiquitous:    b.UbiquitousDir,
	}
}

// root is where git is asked about revisions, defaulting to the directory the
// build already runs in.
func (b *Builder) root() string {
	if b.Root == "" {
		return "."
	}
	return b.Root
}

type sidebarCounts struct {
	opportunities int
	tasks         int
	mappings      int
	storyMaps     int
	stories       int
	terms         int
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
	mappings, err := parser.ParseAllExampleMappings(b.MappingsDir)
	if err != nil {
		return sidebarCounts{}, err
	}
	automations, err := b.automationIndex()
	if err != nil {
		return sidebarCounts{}, err
	}
	for _, em := range mappings {
		automations.Attach(em)
	}
	terms, err := parser.ParseAllTerms(b.UbiquitousDir)
	if err != nil {
		return sidebarCounts{}, err
	}
	opportunities, err := parser.ParseAllOpportunities(b.OpportunitiesDir)
	if err != nil {
		return sidebarCounts{}, err
	}
	tasks := 0
	for _, em := range mappings {
		// Counted off the active view, so the badge matches what the Tasks page
		// lists rather than counting retired items it never shows.
		active := em.Active()
		tasks += len(active.Questions)
		for _, r := range active.Rules {
			// A proposed rule is listed until it is agreed, automated or not.
			if r.Proposed() || !r.Automated() {
				tasks++
			}
		}
	}
	return sidebarCounts{
		opportunities: len(opportunities),
		tasks:         tasks,
		mappings:      len(mappings),
		storyMaps:     len(maps),
		stories:       len(stories),
		terms:         len(terms),
	}, nil
}

// sidebar builds the shared navigation data for a hub page. prefix is the
// relative path back to the output root, active marks the current resource type.
func (b *Builder) sidebar(active, prefix string) (Sidebar, error) {
	c, err := b.computeCounts()
	if err != nil {
		return Sidebar{}, err
	}
	sb := Sidebar{
		Prefix:        prefix,
		Active:        active,
		Opportunities: c.opportunities,
		Tasks:         c.tasks,
		Mappings:      c.mappings,
		StoryMaps:     c.storyMaps,
		Stories:       c.stories,
		Terms:         c.terms,
	}
	return sb, nil
}

// generatedDirs are the output subdirectories holding one page per resource.
// Build owns their contents end to end, so it empties them on every run.
var generatedDirs = []string{"story", "mapping", "story-map", "opportunity", "opportunity-canvas", "opportunity-progress"}

// resetGeneratedDirs empties the per-resource output subdirectories, so a page
// for a renamed or deleted resource cannot outlive its source and keep being
// served by `livt serve`. Only these fixed subdirectories are removed, never
// OutDir itself: --out is user supplied and may point at a directory holding
// files livt did not write. The hub pages at the output root are a fixed set
// that every build rewrites, so they cannot go stale this way.
func (b *Builder) resetGeneratedDirs() error {
	for _, d := range generatedDirs {
		dir := filepath.Join(b.OutDir, d)
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) Build() error {
	if err := b.resetGeneratedDirs(); err != nil {
		return err
	}
	b.automations = nil

	// Computed before any page is written, because every hub page's sidebar
	// carries the count — and because a revision that does not resolve should
	// stop the build rather than leave a site half rewritten.
	if err := b.computeDiff(); err != nil {
		return err
	}

	opportunities, err := b.opportunityIndex()
	if err != nil {
		return err
	}

	if err := b.buildOpportunityCanvases(opportunities); err != nil {
		return err
	}

	maps, err := b.buildStoryMaps(opportunities)
	if err != nil {
		return err
	}
	storyToMaps := maps.StoryOpportunities

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
			Name:          story.DisplayName(),
			Opportunities: listOpportunities,
			MappingPath:   listMappingPath,
			Links:         urlMetaFieldViews(story.Meta),
		})
	}

	mappingTiles, open, tallies, err := b.buildMappings()
	if err != nil {
		return err
	}

	// An opportunity's progress is summed from the mappings its stories have,
	// so its pages are built here rather than beside the maps that named those
	// stories.
	opportunityTiles, err := b.buildOpportunities(maps.MapsByOpportunity, maps.StoriesByOpportunity, tallies)
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

	// Unfinished items inherit their story's opportunities the same way, so the
	// Tasks page filters on that one axis across all of its lists.
	var openOpportunitySets [][]opportunityRef
	for _, items := range [][]taskItem{open.Questions, open.ProposedRules, open.UnautomatedRules} {
		for i := range items {
			items[i].Opportunities = rootRelativeOpportunities(storyToMaps[items[i].StoryKey])
			openOpportunitySets = append(openOpportunitySets, items[i].Opportunities)
		}
	}

	if err := b.buildGlossary(); err != nil {
		return err
	}

	// Hub pages share the sidebar; index.html is the Example Mappings overview
	// and the site's landing page.
	if err := b.buildMappingsIndex(mappingTiles, distinctOpportunityNames(mappingOpportunitySets)); err != nil {
		return err
	}
	fmt.Printf("  index.html\n")

	if err := b.buildTasks(open, distinctOpportunityNames(openOpportunitySets)); err != nil {
		return err
	}
	fmt.Printf("  tasks.html\n")

	if err := b.buildStoryMapsIndex(maps.Tiles); err != nil {
		return err
	}
	fmt.Printf("  story-maps.html\n")

	if err := b.buildOpportunitiesIndex(opportunityTiles); err != nil {
		return err
	}
	fmt.Printf("  opportunities.html\n")

	if err := b.buildStoriesIndex(storyItems, distinctOpportunityNames(storyOpportunitySets)); err != nil {
		return err
	}
	fmt.Printf("  stories.html\n")

	if b.diffResult == nil {
		return b.removeDiff()
	}
	if err := b.buildDiff(); err != nil {
		return err
	}
	fmt.Printf("  %s\n", diffPage)

	return nil
}

// computeDiff reads the two revisions for this run, leaving diffResult nil when
// the build was given no range — which is what keeps the diff off a site nobody
// asked one for.
func (b *Builder) computeDiff() error {
	b.diffResult, b.diffByURI = nil, nil
	if b.Diff == nil {
		return nil
	}
	result, err := b.Diff.Compute(b.root(), b.diffDirs())
	if err != nil {
		return err
	}
	b.diffResult = result
	b.diffByURI = make(map[string]diff.Became, len(result.Changes))
	for _, c := range result.Changes {
		b.diffByURI[c.URI] = c.Became
	}
	return nil
}

// InputDirs is every directory Build reads, in the order the site is composed
// from them. `livt serve` watches this rather than a list of its own, so an
// input the build gained cannot end up unwatched — an edit that changes a page
// and does not reload it is indistinguishable from livt ignoring the edit.
func (b *Builder) InputDirs() []string {
	return []string{
		b.OpportunitiesDir,
		b.CanvasesDir,
		b.MappingsDir,
		b.StoriesDir,
		b.USMDir,
		b.UbiquitousDir,
		b.AutomationsDir,
	}
}
