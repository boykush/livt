package builder

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/diff"
	"github.com/boykush/livt/internal/i18n"
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

// diffMarkView is the note a resource's own page carries when that resource
// changed between the revisions: what happened to it, and where in the diff to
// read what it used to say. It is how the diff is reached at all — the diff is
// not a kind of thing the livt repository holds, but something said about the
// things it does, so the way in is the thing that changed.
type diffMarkView struct {
	Status string
	Label  string
	Href   string
}

// diffAnchor is a URI's id on the diff page. The livt:// scheme is dropped and
// the rest kept as it stands: a livt URI's path is already unique across every
// kind, so nothing is lost and the fragment stays readable in a link someone
// pastes into a review.
func diffAnchor(u string) string {
	return strings.TrimPrefix(u, "livt://")
}

// diffMark is what one URI's page shows, nil when that URI did not change — or
// when the build was given no revisions, which is what keeps the mark off every
// page without a single template having to ask.
func (b *Builder) diffMark(prefix, u string) *diffMarkView {
	became, changed := b.diffByURI[u]
	if !changed {
		return nil
	}
	return &diffMarkView{
		Status: string(became),
		Label:  i18n.Of(b.Lang).Msg("diff." + string(became)),
		Href:   prefix + diffPage + "#" + diffAnchor(u),
	}
}

// diffMarks is every changed URI at once, for a page rendering many of them.
func (b *Builder) diffMarks(prefix string) map[string]*diffMarkView {
	if len(b.diffByURI) == 0 {
		return nil
	}
	marks := make(map[string]*diffMarkView, len(b.diffByURI))
	for u := range b.diffByURI {
		marks[u] = b.diffMark(prefix, u)
	}
	return marks
}

// diffGoneView is what changed under a page and is no longer on it: an item
// deleted, or retired and so off the board. Nothing there can carry its mark,
// so the page carries the count — an absence announced nowhere is one the
// reader takes for something that never existed.
type diffGoneView struct {
	Count int
	Label string
	Href  string
}

// diffGoneFromList is what a hub list lost: resources of its kind the site no
// longer holds at all, which is what an empty Page records.
func (b *Builder) diffGoneFromList(kind uri.Kind) *diffGoneView {
	if b.diffResult == nil {
		return nil
	}
	count := 0
	for _, c := range b.diffResult.Changes {
		if diffKindGroup[c.Kind] == kind && c.Page == "" {
			count++
		}
	}
	return b.gone(count, diffPage)
}

func (b *Builder) gone(count int, href string) *diffGoneView {
	if count == 0 {
		return nil
	}
	return &diffGoneView{Count: count, Label: i18n.Of(b.Lang).Msg("diff.gone"), Href: href}
}

// buildDiff renders diff.html from the range computed for this run. It reads
// diffResult rather than taking the result as an argument, because the sidebar
// on every page already reads it: two sources for one thing is how a page comes
// to link a diff the build did not write.
func (b *Builder) buildDiff() error {
	result := b.diffResult
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
		Sidebar:   sb,
		Base:      result.Base,
		Head:      result.Head,
		Added:     result.Added,
		Changed:   result.Changed,
		Withdrawn: result.Withdrawn,
		Groups:    b.diffGroupViews(result.Changes),
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
func (b *Builder) diffGroupViews(changes []diff.Change) []diffGroupView {
	entries := make(map[uri.Kind][]*diffEntryView)
	parents := make(map[string]*diffEntryView)
	for _, c := range changes {
		entry := &diffEntryView{
			URI:    c.URI,
			Anchor: diffAnchor(c.URI),
			Title:  c.Title,
			Status: string(c.Became),
			Path:   c.Page,
			Lines:  b.diffLineViews(c.Lines),
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
			// Named off the child, because the mapping has no change of its own
			// to be named by: nothing about it moved except the rules hanging
			// off it.
			parent = &diffEntryView{URI: c.Parent, Anchor: diffAnchor(c.Parent), Title: c.ParentTitle}
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

func (b *Builder) diffLineViews(lines []diff.Line) []diffLineView {
	views := make([]diffLineView, 0, len(lines))
	for _, l := range lines {
		label := b.diffLabel(l.Field)
		views = append(views, diffLineView{
			Op:    string(l.Op),
			Label: label,
			Text:  diffText(label, b.diffValue(l.Field)),
			Parts: diffPartViews(l.Parts),
		})
	}
	return views
}

// diffLabel and diffValue put the site's own words on a field. What livt names
// it translates; what the livt repository named — a frontmatter key — is the
// repository's word and is carried through as written.
func (b *Builder) diffLabel(f diff.Field) string {
	if f.Translate {
		return i18n.Of(b.Lang).Msg(f.Label)
	}
	return f.Label
}

// diffValue translates a value that is itself a message key, which is how a
// rule's status reaches the page as a word rather than as the spelling in the
// file.
func (b *Builder) diffValue(f diff.Field) string {
	if f.Label == diff.LabelStatus {
		return i18n.Of(b.Lang).Msg(f.Value)
	}
	return f.Value
}

// diffText is one line as it reads. A field with no label is the item speaking
// for itself; one with no value says something by being there at all, so the
// word stands alone rather than trailing an empty colon.
func diffText(label, value string) string {
	switch {
	case label == "":
		return value
	case value == "":
		return label
	}
	return label + ": " + value
}

// diffPartViews carry a reworded line's breakdown. A line with none renders
// whole — there was nothing to pair it with, or the two were different
// sentences rather than one edited.
func diffPartViews(parts []diff.Part) []diffPartView {
	if len(parts) == 0 {
		return nil
	}
	views := make([]diffPartView, 0, len(parts))
	for _, p := range parts {
		views = append(views, diffPartView{Changed: p.Changed, Text: p.Text})
	}
	return views
}
