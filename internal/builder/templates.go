package builder

import (
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/uri"
)

//go:embed templates/*.html
var templateFS embed.FS

// All templates are parsed into a single set so pages can share the
// {{define "sidebar"}} partial in _sidebar.html.
var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
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
}).ParseFS(templateFS, "templates/*.html"))

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
	// Tasks is what the Tasks page lists: open questions plus rules with no
	// automation recorded.
	Tasks     int
	Mappings  int
	StoryMaps int
	Stories   int
	Terms     int
}

type mappingTile struct {
	Key           string
	StoryName     string
	Opportunities []opportunityRef
}

type mappingsIndexView struct {
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
}

type opportunitiesIndexView struct {
	Sidebar       Sidebar
	Opportunities []opportunityTile
}

// opportunityView is one opportunity's own page. CanvasPath is empty when no
// canvas has been filled in for it, the same way storyView.MappingPath is.
type opportunityView struct {
	Opportunity *domain.Opportunity
	Meta        []metaFieldView
	CanvasPath  string
	StoryMaps   []storyMapRef
}

// opportunityCanvasView is the sheet. Panels are its three columns, since the
// canvas is laid out by zone rather than by the order its boxes are filled in.
type opportunityCanvasView struct {
	OpportunityKey  string
	OpportunityName string
	Panels          domain.CanvasPanels
	Ubiquitous      []termCard
}

type storyMapsIndexView struct {
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
	Sidebar             Sidebar
	Stories             []storyItem
	FilterOpportunities []string
}

// taskItem is one thing the livt repository says is not finished, lifted off its example
// mapping onto the Tasks page: an open question, or a rule with no automation
// recorded. It carries no status, and cannot carry one derived from an issue —
// the livt repository holds automation issue URLs but not their open/closed state, and
// the build may not ask GitHub
// (livt://mapping/show-automation-status-per-rule/rule/R-02). Kind ("question"
// or "rule") picks the sticky colour the item carries on its board.
// StoryPath is empty when the story has no page (mirrors mappingView.StoryPath);
// MappingPath deep-links to the item's own sticky.
type taskItem struct {
	Kind          string
	ID            string
	Text          string
	StoryKey      string
	StoryName     string
	StoryPath     string
	MappingPath   string
	Opportunities []opportunityRef
}

// tasksView is the livt repository's two unfinished flanks either side of the example
// mapping: questions close by a conversation, un-automated rules close by a
// test. One FilterOpportunities set serves both lists, since a single filter bar
// drives the page.
type tasksView struct {
	Sidebar             Sidebar
	Questions           []taskItem
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
}

// glossaryView is the whole ubiquitous language as one table. Contexts are the
// axes it can be narrowed to, omitting the context-free terms — they have no
// axis of their own.
type glossaryView struct {
	Sidebar  Sidebar
	Terms    []glossaryCard
	Contexts []string
}

type storyView struct {
	Story         *domain.Story
	Meta          []metaFieldView
	MappingPath   string
	Opportunities []opportunityRef
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
	StoryName  string
	StoryPath  string
	Mapping    *domain.ExampleMapping
	Ubiquitous []termCard
}

func renderTasks(w io.Writer, view tasksView) error {
	return tmpl.ExecuteTemplate(w, "tasks.html", view)
}

func renderMappingsIndex(w io.Writer, view mappingsIndexView) error {
	return tmpl.ExecuteTemplate(w, "index.html", view)
}

func renderOpportunitiesIndex(w io.Writer, view opportunitiesIndexView) error {
	return tmpl.ExecuteTemplate(w, "opportunities.html", view)
}

func renderOpportunity(w io.Writer, view opportunityView) error {
	return tmpl.ExecuteTemplate(w, "opportunity.html", view)
}

func renderOpportunityCanvas(w io.Writer, view opportunityCanvasView) error {
	return tmpl.ExecuteTemplate(w, "opportunity_canvas.html", view)
}

func renderStoryMapsIndex(w io.Writer, view storyMapsIndexView) error {
	return tmpl.ExecuteTemplate(w, "story-maps.html", view)
}

func renderStoriesIndex(w io.Writer, view storiesIndexView) error {
	return tmpl.ExecuteTemplate(w, "stories.html", view)
}

func renderStory(w io.Writer, story *domain.Story, mappingPath string, opportunities []opportunityRef) error {
	return tmpl.ExecuteTemplate(w, "story.html", storyView{
		Story:         story,
		Meta:          metaFieldViews(story.Meta),
		MappingPath:   mappingPath,
		Opportunities: opportunities,
	})
}

// renderMapping draws the board from the mapping's active view: a retired
// sticky is off the wall, whichever kind it is, so the board shows what the
// spec asks for today.
func renderMapping(w io.Writer, em *domain.ExampleMapping, storyName, storyPath string, ubiquitous []termCard) error {
	return tmpl.ExecuteTemplate(w, "mapping.html", mappingView{StoryName: storyName, StoryPath: storyPath, Mapping: em.Active(), Ubiquitous: ubiquitous})
}

func renderStoryMap(w io.Writer, view storyMapView) error {
	return tmpl.ExecuteTemplate(w, "story_map.html", view)
}

func renderGlossary(w io.Writer, view glossaryView) error {
	return tmpl.ExecuteTemplate(w, "glossary.html", view)
}
