package builder

import (
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/i18n"
	"github.com/boykush/livt/internal/uri"
)

//go:embed templates/*.html
var templateFS embed.FS

// funcs are the template helpers bound to one language. Only "t" varies with
// it; the rest are pure and shared by every set.
func funcs(lang i18n.Lang) template.FuncMap {
	return template.FuncMap{
		"t":                 i18n.Of(lang).T,
		"issueLabel":        issueLabel,
		"filter":            newFilterView,
		"opportunityValues": opportunityValuesJSON,
		"contextValues":     contextValuesJSON,
		"ruleAnchor":        uri.RuleAnchor,
		"exampleAnchor":     uri.ExampleAnchor,
		"questionAnchor":    uri.QuestionAnchor,
		"storyCardAnchor":   uri.StoryCardAnchor,
		"ruleBadge":         ruleBadge,
		"exampleBadge":      exampleBadge,
		"questionBadge":     questionBadge,
		"storyBadge":        storyBadge,
		"percent":           percent,
		"meter":             newMeterView,
		"tasksAnchors":      tasksAnchors,
		"ruleURI":           uri.Rule,
		"exampleURI":        uri.Example,
		"questionURI":       uri.Question,
	}
}

// All templates are parsed into a single set so pages can share the
// {{define "sidebar"}} partial in _sidebar.html.
var baseTmpl = template.Must(template.New("").Funcs(funcs(i18n.Default)).ParseFS(templateFS, "templates/*.html"))

// localized holds one set per language. html/template locks a set the first
// time it renders, so every language's "t" has to be bound into its own clone
// up front — baseTmpl is only ever cloned from, never executed.
var localized = localize()

func localize() map[i18n.Lang]*template.Template {
	sets := make(map[i18n.Lang]*template.Template, len(i18n.Langs))
	for _, lang := range i18n.Langs {
		sets[lang] = template.Must(baseTmpl.Clone()).Funcs(funcs(lang))
	}
	return sets
}

// templates picks the set for a language, falling back to the default for the
// zero value so a Builder left unconfigured still renders.
func templates(lang i18n.Lang) *template.Template {
	if set, ok := localized[lang]; ok {
		return set
	}
	return localized[i18n.Default]
}

// percent is a progress bar's width. A zero total reads as 0 rather than
// dividing: an opportunity that has taken on no story has not finished one.
func percent(part, total int) int {
	if total <= 0 {
		return 0
	}
	return part * 100 / total
}

// meterView is one figure on the opportunity dashboard: a count drawn as its
// share of the whole it belongs to, and the page that count can be acted on.
// Every figure gets its own meter rather than becoming a segment of another's,
// because what closes each of them differs — a test, an agreement, a
// conversation — and a segment inside someone else's bar is not a thing you can
// go and do.
type meterView struct {
	Label string
	Part  int
	Total int
	Fill  string
	// Down marks a burn-down: a count that is done when it reaches zero, as
	// against one done when it reaches Total. The fraction reads the same
	// either way, so the meter says which with a glyph beside it.
	Down bool
	// LinkLabel names the page the count is acted on; empty draws no footer.
	LinkLabel string
	LinkPath  string
}

func newMeterView(label string, part, total int, fill string, down bool, linkLabel, linkPath string) meterView {
	return meterView{Label: label, Part: part, Total: total, Fill: fill, Down: down, LinkLabel: linkLabel, LinkPath: linkPath}
}

// idBadge is a sticky's own ID rendered bottom-right by the id-badge partial.
// Label is the ID as the livt repository numbers it, which for an example is local to
// its rule and so shorter than Anchor. Tint follows the sticky's own
// border-{colour}-400 so the badge recedes into the card it belongs to.
type idBadge struct {
	Anchor string
	Label  string
	Tint   string
}

func ruleBadge(ruleID string) idBadge {
	return idBadge{Anchor: uri.RuleAnchor(ruleID), Label: ruleID, Tint: "text-blue-400/70 hover:text-blue-600"}
}

func exampleBadge(ruleID, exampleID string) idBadge {
	return idBadge{Anchor: uri.ExampleAnchor(ruleID, exampleID), Label: exampleID, Tint: "text-green-400/70 hover:text-green-600"}
}

func questionBadge(questionID string) idBadge {
	return idBadge{Anchor: uri.QuestionAnchor(questionID), Label: questionID, Tint: "text-red-400/70 hover:text-red-600"}
}

// storyBadge labels a story card with its story key. Yellow breaks the
// border-matching shade the other three share: yellow-400 on the card's
// yellow-100 is barely there, so the badge drops to yellow-600.
func storyBadge(storyKey string) idBadge {
	return idBadge{Anchor: uri.StoryCardAnchor(storyKey), Label: storyKey, Tint: "text-yellow-600/70 hover:text-yellow-700"}
}

// issueLabel shortens an automation Issue URL to a sticky-sized link label:
// https://github.com/{owner}/{repo}/issues/{n} becomes "{repo}#{n}". Anything
// else falls back to its host so foreign trackers still read as a chip.
func issueLabel(url string) string {
	rest, ok := strings.CutPrefix(url, "https://")
	if !ok {
		rest, ok = strings.CutPrefix(url, "http://")
	}
	if !ok {
		return url
	}
	parts := strings.Split(rest, "/")
	if len(parts) >= 5 && parts[0] == "github.com" && parts[3] == "issues" {
		return parts[2] + "#" + parts[4]
	}
	return parts[0]
}

// opportunityRef links a card to one opportunity it belongs to. Path leads to
// the opportunity's page, or — when no opportunity file claims the map's key —
// to the story map standing in as its own opportunity. Both are relative to the
// page doing the rendering.
type opportunityRef struct {
	Name string
	Path string
}

// filterView drives the shared filter bar. Param is both the query parameter the
// selection is mirrored in and the name of the axis, Label heads the chip row,
// and Values are the axes a card can be narrowed to. One bar per page, so the
// script finds its param on the bar rather than being told it twice.
type filterView struct {
	Param  string
	Label  string
	Values []string
	// Chip is the class set the bar's chips wear and Pressed the colour a
	// selected one fills with. Both follow the resource the axis narrows on, so a
	// chip in the bar reads as the same thing as the chip on a card below it.
	Chip    string
	Pressed string
}

// filterTints maps a filter axis to its resource's colour. An axis with no
// entry falls back to the neutral one rather than borrowing another resource's.
var filterTints = map[string][2]string{
	"opportunity": {"text-orange-700 bg-orange-50 border border-orange-200 hover:bg-orange-100", "#ea580c"},
	"context":     {"text-pink-700 bg-pink-50 border border-pink-200 hover:bg-pink-100", "#db2777"},
}

var neutralTint = [2]string{"text-gray-700 bg-gray-50 border border-gray-200 hover:bg-gray-100", "#4b5563"}

// opportunityFilterParam is the query parameter the opportunity filter bar
// mirrors its selection in, and the name the templates pass to `filter`. An
// opportunity's page links into those lists pre-narrowed through it, so the two
// have to agree: a link written against a different name lands on the unfiltered
// list and silently shows every opportunity's items as this one's.
const opportunityFilterParam = "opportunity"

// The Tasks page keeps one list per way an item gets closed, and pages that
// count those items link straight at the list they counted — landing on the
// first of three leaves the reader to find the other two. Named here so the
// anchors the template writes and the links other pages aim cannot drift.
const (
	tasksQuestionsAnchor = "open-questions"
	tasksProposedAnchor  = "proposed-rules"
	tasksRulesAnchor     = "unautomated-rules"
)

// tasksAnchors is what the Tasks template writes onto its three sections.
func tasksAnchors() map[string]string {
	return map[string]string{
		"questions": tasksQuestionsAnchor,
		"proposed":  tasksProposedAnchor,
		"rules":     tasksRulesAnchor,
	}
}

func newFilterView(param, label string, values []string) filterView {
	tint, ok := filterTints[param]
	if !ok {
		tint = neutralTint
	}
	return filterView{Param: param, Label: label, Values: values, Chip: tint[0], Pressed: tint[1]}
}

// filterValuesJSON encodes a card's values on the filter axis as a JSON array
// for its data-filter-values attribute, the hook the client-side filter matches
// against. A card on no value serializes to "[]" so it matches no axis.
func filterValuesJSON(values []string) string {
	// A nil slice marshals to null, which the filter would fail to read as the
	// empty list it means.
	if values == nil {
		values = []string{}
	}
	b, err := json.Marshal(values)
	if err != nil {
		return "[]"
	}
	return string(b)
}

// opportunityValuesJSON is the filter hook for a card carrying opportunity refs.
// A story on no map matches no filter axis.
func opportunityValuesJSON(refs []opportunityRef) string {
	names := make([]string, len(refs))
	for i, r := range refs {
		names[i] = r.Name
	}
	return filterValuesJSON(names)
}

// contextValuesJSON is the filter hook for a glossary row. A term holds at most
// one context, and a context-free one matches no axis: it belongs to all of
// them, so narrowing to any single context would misreport it as scoped there.
func contextValuesJSON(ctx string) string {
	if ctx == "" {
		return filterValuesJSON(nil)
	}
	return filterValuesJSON([]string{ctx})
}

// rootRelativeOpportunities rebases refs whose Path is relative to the story/
// directory ("../story-map/x.html") onto the output root ("story-map/x.html"),
// where the Stories and Example Mappings lists render.
func rootRelativeOpportunities(refs []opportunityRef) []opportunityRef {
	out := make([]opportunityRef, len(refs))
	for i, r := range refs {
		out[i] = opportunityRef{Name: r.Name, Path: strings.TrimPrefix(r.Path, "../")}
	}
	return out
}

// distinctOpportunityNames collects the unique opportunity names across a list's
// cards in first-appearance order, so a list only offers filter axes that match
// at least one of its cards.
func distinctOpportunityNames(perCard [][]opportunityRef) []string {
	seen := make(map[string]bool)
	var names []string
	for _, refs := range perCard {
		for _, r := range refs {
			if !seen[r.Name] {
				seen[r.Name] = true
				names = append(names, r.Name)
			}
		}
	}
	return names
}

// Sidebar is the shared navigation rendered on every hub page. Prefix is the
// relative path back to the output root ("" for root pages); Active marks the
// current resource type.
type Sidebar struct {
	Prefix string
	Active string
	// Opportunities counts the committed opportunities, canvas or not: the
	// opportunity is the unit, and a canvas is one thing that may be said of it.
	Opportunities int
	// Tasks is what the Tasks page lists: open questions, proposed rules, and
	// accepted rules with no automation recorded.
	Tasks     int
	Mappings  int
	StoryMaps int
	Stories   int
	Terms     int
}

type mappingTile struct {
	Key           string
	Name          string
	Opportunities []opportunityRef
}

type mappingsIndexView struct {
	DiffGone            *diffGoneView
	Sidebar             Sidebar
	Mappings            []mappingTile
	FilterOpportunities []string
}

type storyMapTile struct {
	Name string
	// Opportunity is the opportunity this map serves, empty when no opportunity
	// file claims the map's key.
	Opportunity *opportunityRef
}

// storyMapRef links an opportunity to one story map mapped for it, the reverse
// of the join opportunityRef makes.
type storyMapRef struct {
	Name string
	Path string
}

// rootRelativeMaps rebases refs written for a page under a subdirectory onto
// the output root, where the hub lists render (mirrors rootRelativeOpportunities).
func rootRelativeMaps(refs []storyMapRef) []storyMapRef {
	out := make([]storyMapRef, len(refs))
	for i, r := range refs {
		out[i] = storyMapRef{Name: r.Name, Path: strings.TrimPrefix(r.Path, "../")}
	}
	return out
}

// opportunityTile is one card on the Opportunities hub. Statement is the
// opportunity's own body — the one thing a reader scanning the list needs —
// and HasCanvas drives the preview, since an opportunity may be committed
// before any canvas session has been held for it.
type opportunityTile struct {
	Key       string
	Name      string
	Statement string
	HasCanvas bool
	StoryMaps []storyMapRef
	Links     []metaFieldView
	// Progress is the same reading the opportunity's own page leads with, so a
	// reader scanning the hub sees which opportunity is moving without opening
	// each one. Its paths are written for that page and are not followed here.
	Progress opportunityProgress
}

type opportunitiesIndexView struct {
	DiffGone      *diffGoneView
	Sidebar       Sidebar
	Opportunities []opportunityTile
}

// opportunityView is one opportunity's own page. CanvasPath is empty when no
// canvas has been filled in for it, the same way storyView.MappingPath is.
type opportunityView struct {
	Diff        *diffMarkView
	Opportunity *domain.Opportunity
	Meta        []metaFieldView
	CanvasPath  string
	// ProgressPath is empty when the opportunity has taken on no story, the
	// same way CanvasPath is empty when no canvas has been filled in.
	ProgressPath string
	StoryMaps    []storyMapRef
	Progress     opportunityProgress
}

// opportunityProgressView is one opportunity's progress on a page of its own.
// The opportunity's page carries the two gauges and leads here for the reading
// story by story.
type opportunityProgressView struct {
	OpportunityKey  string
	OpportunityName string
	Progress        opportunityProgress
}

// opportunityCanvasView is the sheet. Panels are its three columns, since the
// canvas is laid out by zone rather than by the order its boxes are filled in.
type opportunityCanvasView struct {
	Diff            *diffMarkView
	OpportunityKey  string
	OpportunityName string
	// OpportunityPath is empty when no opportunity file shares the key, the same
	// way a mapping's StoryPath is empty when its story has no card.
	OpportunityPath string
	Panels          domain.CanvasPanels
	Ubiquitous      []termCard
}

type storyMapsIndexView struct {
	DiffGone  *diffGoneView
	Sidebar   Sidebar
	StoryMaps []storyMapTile
}

type storyItem struct {
	Key           string
	Name          string
	Opportunities []opportunityRef
	MappingPath   string
	Links         []metaFieldView
}

// metaFieldView renders one Story frontmatter field. Href is non-empty only when
// Value is an http(s) URL, so templates branch on its presence to link the value
// (mirrors termCard.Href) instead of carrying a boolean.
type metaFieldView struct {
	Key   string
	Value string
	Href  string
}

func isHTTPURL(s string) bool {
	return strings.HasPrefix(s, "https://") || strings.HasPrefix(s, "http://")
}

// metaFieldViews maps every frontmatter field to a view, auto-linking URL values.
func metaFieldViews(meta []domain.MetaField) []metaFieldView {
	views := make([]metaFieldView, 0, len(meta))
	for _, m := range meta {
		view := metaFieldView{Key: m.Key, Value: m.Value}
		if isHTTPURL(m.Value) {
			view.Href = m.Value
		}
		views = append(views, view)
	}
	return views
}

// urlMetaFieldViews keeps only URL-valued fields, used as quick-link chips on the
// Stories list.
func urlMetaFieldViews(meta []domain.MetaField) []metaFieldView {
	var views []metaFieldView
	for _, m := range meta {
		if isHTTPURL(m.Value) {
			views = append(views, metaFieldView{Key: m.Key, Value: m.Value, Href: m.Value})
		}
	}
	return views
}

type storiesIndexView struct {
	DiffGone            *diffGoneView
	Sidebar             Sidebar
	Stories             []storyItem
	FilterOpportunities []string
}

// taskItem is one unfinished thing lifted off its example mapping onto the Tasks
// page. It carries no status derived from an issue: the build may not ask GitHub
// (livt://mapping/show-automation-status-per-rule/rule/R-02). Kind ("question",
// "proposed-rule" or "rule") picks the sticky colour it wears on its board;
// MappingPath deep-links to that sticky, and StoryPath is empty with no page.
type taskItem struct {
	Kind     string
	ID       string
	Text     string
	StoryKey string
	// MappingName is what the board's yellow sticky reads, so the chip and the
	// sticky it stands for cannot name the same board two ways.
	MappingName   string
	StoryPath     string
	MappingPath   string
	Opportunities []opportunityRef
}

// tasksView is what the livt repository leaves unfinished, one list per way of
// closing it: questions by a conversation, proposed rules by agreement,
// un-automated rules by a test. One FilterOpportunities set serves every list,
// since a single filter bar drives the page.
type tasksView struct {
	Sidebar             Sidebar
	Questions           []taskItem
	ProposedRules       []taskItem
	UnautomatedRules    []taskItem
	FilterOpportunities []string
}

// glossaryCard is one row of the glossary table. Anchor is the row's id — the
// term's reference, which carries its context — while Key stays the bare key the
// file is named for, so two rows sharing a key still read as the same word meant
// two ways rather than as two unrelated ones.
type glossaryCard struct {
	Anchor     string
	Ctx        string
	Key        string
	Name       string
	Definition string
	// Diff is the note the row carries when the term changed between the
	// revisions a diff build was given, nil otherwise.
	Diff *diffMarkView
}

// glossaryView is the whole ubiquitous language as one table. Contexts are the
// axes it can be narrowed to, omitting the context-free terms — they have no
// axis of their own.
type glossaryView struct {
	DiffGone *diffGoneView
	Sidebar  Sidebar
	Terms    []glossaryCard
	Contexts []string
}

// diffView is what changed between two revisions, gathered by resource type.
// Base and Head are the short hashes git resolved; an empty Head is the working
// tree, which the page says in words because it has no hash to print.
type diffView struct {
	Sidebar   Sidebar
	Base      string
	Head      string
	Added     int
	Changed   int
	Withdrawn int
	Groups    []diffGroupView
}

type diffGroupView struct {
	Label   string
	Entries []diffEntryView
}

// diffEntryView is one livt URI's change. Status is empty on a mapping listed
// only to head its changed rules, and Path is empty when the site holds no page
// for the URI — which is every removed one.
type diffEntryView struct {
	URI      string
	Anchor   string
	Title    string
	Status   string
	Path     string
	Lines    []diffLineView
	Children []diffEntryView
}

// diffLineView is one line of an entry's diff, Op being git diff's own " ", "+"
// or "-" so the page needs no legend to be read. Text is the whole line; Parts
// is the same value broken into what stayed and what moved, and is set only on
// the two halves of a rewording.
type diffLineView struct {
	Op    string
	Label string
	Text  string
	Parts []diffPartView
}

// diffPartView is a run of a reworded line, marked where it changed.
type diffPartView struct {
	Changed bool
	Text    string
}

type storyView struct {
	Story         *domain.Story
	Meta          []metaFieldView
	MappingPath   string
	Opportunities []opportunityRef
	Diff          *diffMarkView
}

// termCard is a referenced ubiquitous language term rendered as a pink sticky on
// a board. Href is the link to its glossary row, or empty when the reference
// resolves to no term file (then it renders as a plain card). Ctx is shown on
// the sticky when the term has one: the display name alone would not say which
// context's meaning a board is reaching for.
type termCard struct {
	Name string
	Ctx  string
	Href string
}

type mappingView struct {
	// Name is the board's, which its yellow sticky reads; StoryPath still leads
	// to the story when there is one, however differently the two are named.
	Name       string
	StoryPath  string
	Mapping    *domain.ExampleMapping
	Ubiquitous []termCard
	Diff       *diffMarkView
	// Ghosts marks the stickies the board is only showing because they were
	// withdrawn in this diff. What became of them is on their own mark; this
	// only says to draw them as no longer spec.
	Ghosts map[string]bool
	// DiffMarks is every changed URI on this board, keyed by URI: the stickies
	// are rendered straight off the domain types, so the mark is looked up
	// beside each rather than carried on it.
	DiffMarks map[string]*diffMarkView
}

func renderTasks(w io.Writer, lang i18n.Lang, view tasksView) error {
	return templates(lang).ExecuteTemplate(w, "tasks.html", view)
}

func renderMappingsIndex(w io.Writer, lang i18n.Lang, view mappingsIndexView) error {
	return templates(lang).ExecuteTemplate(w, "index.html", view)
}

func renderOpportunitiesIndex(w io.Writer, lang i18n.Lang, view opportunitiesIndexView) error {
	return templates(lang).ExecuteTemplate(w, "opportunities.html", view)
}

func renderOpportunity(w io.Writer, lang i18n.Lang, view opportunityView) error {
	return templates(lang).ExecuteTemplate(w, "opportunity.html", view)
}

func renderOpportunityProgress(w io.Writer, lang i18n.Lang, view opportunityProgressView) error {
	return templates(lang).ExecuteTemplate(w, "opportunity_progress.html", view)
}

func renderOpportunityCanvas(w io.Writer, lang i18n.Lang, view opportunityCanvasView) error {
	return templates(lang).ExecuteTemplate(w, "opportunity_canvas.html", view)
}

func renderStoryMapsIndex(w io.Writer, lang i18n.Lang, view storyMapsIndexView) error {
	return templates(lang).ExecuteTemplate(w, "story-maps.html", view)
}

func renderStoriesIndex(w io.Writer, lang i18n.Lang, view storiesIndexView) error {
	return templates(lang).ExecuteTemplate(w, "stories.html", view)
}

func renderStory(w io.Writer, lang i18n.Lang, story *domain.Story, mappingPath string, opportunities []opportunityRef, mark *diffMarkView) error {
	return templates(lang).ExecuteTemplate(w, "story.html", storyView{
		Story:         story,
		Meta:          metaFieldViews(story.Meta),
		MappingPath:   mappingPath,
		Opportunities: opportunities,
		Diff:          mark,
	})
}

// renderMapping draws the board from the mapping's active view: a retired
// sticky is off the wall, whichever kind it is, so the board shows what the
// spec asks for today.
func renderMapping(w io.Writer, lang i18n.Lang, bd board, name, storyPath string, ubiquitous []termCard, diff *diffMarkView, marks map[string]*diffMarkView) error {
	return templates(lang).ExecuteTemplate(w, "mapping.html", mappingView{
		Name: name, StoryPath: storyPath, Mapping: bd.Mapping,
		Ubiquitous: ubiquitous, Diff: diff, Ghosts: bd.Ghosts, DiffMarks: marks,
	})
}

func renderStoryMap(w io.Writer, lang i18n.Lang, view storyMapView) error {
	return templates(lang).ExecuteTemplate(w, "story_map.html", view)
}

func renderGlossary(w io.Writer, lang i18n.Lang, view glossaryView) error {
	return templates(lang).ExecuteTemplate(w, "glossary.html", view)
}

func renderDiff(w io.Writer, lang i18n.Lang, view diffView) error {
	return templates(lang).ExecuteTemplate(w, "diff.html", view)
}
