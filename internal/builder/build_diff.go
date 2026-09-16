package builder

import (
	"os"
	"path/filepath"

	"github.com/boykush/livt/internal/diff"
	"github.com/boykush/livt/internal/uri"
)

// diffPage is where the diff lands. It sits at the output root beside the hub
// pages, and is written only by a build that was given revisions — a build
// given none must leave the site exactly as it was.
const diffPage = "diff.html"

// diffGroups are the resource types the diff is gathered under, in the order
// the sidebar lists them, so a reader looks for a change where they already
// look for the thing that changed. A rule, example, or question is not a group
// of its own: it gathers under the mapping it hangs off.
var diffGroups = []struct {
	kind  uri.Kind
	label string
}{
	{uri.KindMapping, "nav.example-mappings"},
	{uri.KindOpportunity, "nav.opportunities"},
	{uri.KindStoryMap, "nav.story-maps"},
	{uri.KindStory, "nav.stories"},
	{uri.KindTerm, "nav.ubiquitous"},
}

// diffKindGroup says which group a kind is listed under. A canvas goes with its
// opportunity rather than standing alone, the same way a rule goes with its
// mapping — both are one more thing said about the resource above them.
var diffKindGroup = map[uri.Kind]uri.Kind{
	uri.KindMapping:           uri.KindMapping,
	uri.KindRule:              uri.KindMapping,
	uri.KindExample:           uri.KindMapping,
	uri.KindQuestion:          uri.KindMapping,
	uri.KindOpportunity:       uri.KindOpportunity,
	uri.KindOpportunityCanvas: uri.KindOpportunity,
	uri.KindStoryMap:          uri.KindStoryMap,
	uri.KindStory:             uri.KindStory,
	uri.KindTerm:              uri.KindTerm,
}

// buildDiff renders diff.html from a computed result.
func (b *Builder) buildDiff(result *diff.Result) error {
	sb, err := b.sidebar("diff", "")
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(b.OutDir, diffPage))
	if err != nil {
		return err
	}
	defer f.Close()
	return renderDiff(f, b.Lang, diffView{
		Sidebar:  sb,
		Base:     result.Base,
		Head:     result.Head,
		Added:    result.Added,
		Removed:  result.Removed,
		Modified: result.Modified,
		Groups:   diffGroupViews(result.Changes),
	})
}

// removeDiff takes the page away again, so a build given no revisions leaves
// none of the last one's diff behind in an output directory it is reusing.
func (b *Builder) removeDiff() error {
	if err := os.Remove(filepath.Join(b.OutDir, diffPage)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// diffGroupViews gathers the changes under their resource type, and the ones
// belonging to a mapping under that mapping. A mapping whose own fields did not
// change still heads its rules: the group needs something to hang them off, so
// it is listed with no status of its own.
func diffGroupViews(changes []diff.Change) []diffGroupView {
	entries := make(map[uri.Kind][]*diffEntryView)
	parents := make(map[string]*diffEntryView)
	for _, c := range changes {
		entry := &diffEntryView{
			URI:    c.URI,
			Title:  c.Title,
			Status: string(c.Status),
			Path:   c.Page,
			Lines:  diffLineViews(c.Lines),
		}
		group := diffKindGroup[c.Kind]
		if c.Parent == "" {
			if held, ok := parents[c.URI]; ok {
				// The mapping was listed as a bare header by a child that came
				// first; fill in its own change rather than listing it twice.
				held.Status, held.Path, held.Lines = entry.Status, entry.Path, entry.Lines
				continue
			}
			parents[c.URI] = entry
			entries[group] = append(entries[group], entry)
			continue
		}
		parent, ok := parents[c.Parent]
		if !ok {
			parent = &diffEntryView{URI: c.Parent, Title: c.Parent}
			parents[c.Parent] = parent
			entries[group] = append(entries[group], parent)
		}
		parent.Children = append(parent.Children, *entry)
	}
	var groups []diffGroupView
	for _, g := range diffGroups {
		listed := entries[g.kind]
		if len(listed) == 0 {
			continue
		}
		view := diffGroupView{Label: g.label}
		for _, e := range listed {
			view.Entries = append(view.Entries, *e)
		}
		groups = append(groups, view)
	}
	return groups
}

func diffLineViews(lines []diff.Line) []diffLineView {
	views := make([]diffLineView, 0, len(lines))
	for _, l := range lines {
		views = append(views, diffLineView{Op: string(l.Op), Text: l.Text})
	}
	return views
}
