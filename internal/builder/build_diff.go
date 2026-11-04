package builder

import (
	"cmp"
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
		Sections:  b.diffSectionViews(result),
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

// diffSectionViews gathers the changes under their resource type. Under a
// mapping, an example joins the item of the rule it illustrates, since a rule
// is weighed together with what illustrates it; a rule that did not change
// itself still heads that item, in the words result.Context carries for it.
func (b *Builder) diffSectionViews(result *diff.Result) []diffSectionView {
	type group struct {
		view  diffGroupView
		items []*diffItemView
		byURI map[string]*diffItemView
	}
	groups := make(map[string]*group)
	var groupOrder []string
	flat := make(map[uri.Kind][]diffItemView)
	for _, c := range result.Changes {
		kind := diffKindGroup[c.Kind]
		if kind != uri.KindMapping {
			flat[kind] = append(flat[kind], b.diffItem(c))
			continue
		}
		mapping := c.Parent
		if c.Kind == uri.KindMapping {
			mapping = c.URI
		}
		g, ok := groups[mapping]
		if !ok {
			g = &group{view: diffGroupView{Title: cmp.Or(c.ParentTitle, mapping)}, byURI: make(map[string]*diffItemView)}
			groups[mapping] = g
			groupOrder = append(groupOrder, mapping)
		}
		g.view.Count++
		switch c.Kind {
		case uri.KindMapping:
			g.view.Title, g.view.Status = cmp.Or(c.Title, g.view.Title), string(c.Became)
			item := b.diffItem(c)
			item.Title = i18n.Of(b.Lang).Msg("diff.mapping-own")
			// The mapping's name is already its group's heading.
			item.Before[0].Title, item.After[0].Title = "", ""
			g.items = append([]*diffItemView{&item}, g.items...)
		case uri.KindExample:
			parsed, _ := uri.Parse(c.URI)
			ruleURI := uri.Rule(parsed.StoryKey, parsed.RuleID)
			item, held := g.byURI[ruleURI]
			if !held {
				item = b.diffContextItem(ruleURI, result.Context[ruleURI])
				g.byURI[ruleURI] = item
				g.items = append(g.items, item)
			}
			before, after := b.diffCards(c, uri.KindExample)
			before.Nested, after.Nested = true, true
			item.Before = append(item.Before, before)
			item.After = append(item.After, after)
			item.Examples++
		default:
			item := b.diffItem(c)
			if held, ok := g.byURI[c.URI]; ok {
				// The rule's examples came first and gave it a stand-in head;
				// the rule's own change takes that place.
				*held = mergeHead(item, *held)
				continue
			}
			g.byURI[c.URI] = &item
			g.items = append(g.items, &item)
		}
	}
	var sections []diffSectionView
	for _, dg := range diffGroups {
		section := diffSectionView{Label: dg.label}
		if dg.kind == uri.KindMapping {
			for _, m := range groupOrder {
				g := groups[m]
				for _, item := range g.items {
					g.view.Items = append(g.view.Items, *item)
				}
				section.Groups = append(section.Groups, g.view)
			}
		} else {
			section.Items = flat[dg.kind]
		}
		if len(section.Groups) == 0 && len(section.Items) == 0 {
			continue
		}
		sections = append(sections, section)
	}
	return sections
}

// mergeHead puts a rule's own change in front of the examples already held
// under it.
func mergeHead(head, held diffItemView) diffItemView {
	head.Before = append(head.Before, held.Before[1:]...)
	head.After = append(head.After, held.After[1:]...)
	head.Examples = held.Examples
	return head
}

// diffItem is one change read on its own.
func (b *Builder) diffItem(c diff.Change) diffItemView {
	before, after := b.diffCards(c, c.Kind)
	return diffItemView{
		Anchor:  diffAnchor(c.URI),
		Title:   c.Title,
		Summary: diffSummary(c.Lines),
		Status:  string(c.Became),
		URI:     c.URI,
		Path:    c.Page,
		Before:  []diffCardView{before},
		After:   []diffCardView{after},
	}
}

// diffContextItem heads an item with a rule that did not change, so the same
// sticky stands on both sides and only its examples differ.
func (b *Builder) diffContextItem(ruleURI string, rule diff.Entry) *diffItemView {
	lines := make([]diffLineView, 0, len(rule.Fields))
	for _, f := range rule.Fields {
		lines = append(lines, b.diffLineView(diff.Line{Op: diff.OpContext, Field: f}))
	}
	page := ""
	if parsed, ok := uri.Parse(ruleURI); ok {
		page = parsed.Page()
	}
	card := diffCardView{Title: rule.Title, URI: ruleURI, Path: page, Tone: diffTone[uri.KindRule], Lines: lines}
	after := card
	after.Anchor = diffAnchor(ruleURI)
	return &diffItemView{
		Anchor:  diffAnchor(ruleURI),
		Title:   rule.Title,
		Summary: diffSummary(contextLines(rule.Fields)),
		URI:     ruleURI,
		Path:    page,
		Before:  []diffCardView{card},
		After:   []diffCardView{after},
	}
}

func contextLines(fields []diff.Field) []diff.Line {
	lines := make([]diff.Line, 0, len(fields))
	for _, f := range fields {
		lines = append(lines, diff.Line{Op: diff.OpContext, Field: f})
	}
	return lines
}

// diffCards draws one change as the sticky it was and the sticky it became.
// An addition has nothing before it. A withdrawal's after side is its last
// words struck through, with what replaced it beside them; its before side is
// those words as they stood.
func (b *Builder) diffCards(c diff.Change, kind uri.Kind) (before, after diffCardView) {
	base := diffCardView{Title: c.Title, Status: string(c.Became), URI: c.URI, Path: c.Page, Tone: diffTone[kind]}
	before, after = base, base
	after.Anchor = diffAnchor(c.URI)
	for _, l := range c.Lines {
		switch {
		case c.Became == diff.BecameWithdrawn:
			if l.Op == diff.OpDel {
				before.Lines = append(before.Lines, b.diffLineView(diff.Line{Op: diff.OpContext, Field: l.Field}))
			}
			after.Lines = append(after.Lines, b.diffLineView(l))
		case c.Became == diff.BecameAdded:
			after.Lines = append(after.Lines, b.diffLineView(diff.Line{Op: diff.OpContext, Field: l.Field}))
		case l.Op == diff.OpDel:
			before.Lines = append(before.Lines, b.diffLineView(l))
		case l.Op == diff.OpAdd:
			after.Lines = append(after.Lines, b.diffLineView(l))
		default:
			before.Lines = append(before.Lines, b.diffLineView(l))
			after.Lines = append(after.Lines, b.diffLineView(l))
		}
	}
	before.Empty = len(before.Lines) == 0
	after.Ghost = c.Became == diff.BecameWithdrawn
	return before, after
}

// diffSummary is the item's own sentence, for the list: the line with no field
// name in front of it, as the head has it where it still holds.
func diffSummary(lines []diff.Line) string {
	summary := ""
	for _, l := range lines {
		if l.Field.Label != "" {
			continue
		}
		if l.Op != diff.OpDel {
			return l.Field.Value
		}
		summary = cmp.Or(summary, l.Field.Value)
	}
	return summary
}

// diffTone is each kind's sticky, the colours the boards already use for it.
var diffTone = map[uri.Kind]string{
	uri.KindMapping:           "bg-yellow-100 border-yellow-400",
	uri.KindRule:              "bg-blue-100 border-blue-400 font-semibold",
	uri.KindExample:           "bg-green-100 border-green-400",
	uri.KindQuestion:          "bg-red-100 border-red-400",
	uri.KindOpportunity:       "bg-orange-100 border-orange-400",
	uri.KindOpportunityCanvas: "bg-orange-100 border-orange-400",
	uri.KindStoryMap:          "bg-purple-100 border-purple-400",
	uri.KindStory:             "bg-yellow-100 border-yellow-400",
	uri.KindTerm:              "bg-pink-100 border-pink-400",
}

func (b *Builder) diffLineView(l diff.Line) diffLineView {
	label := b.diffLabel(l.Field)
	return diffLineView{
		Op:    string(l.Op),
		Label: label,
		Text:  diffText(label, b.diffValue(l.Field)),
		Parts: diffPartViews(l.Parts),
	}
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
