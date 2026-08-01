package builder

import (
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"strings"

	"github.com/boykush/livt/internal/domain"
)

//go:embed templates/*.html
var templateFS embed.FS

// All templates are parsed into a single set so pages can share the
// {{define "sidebar"}} partial in _sidebar.html.
var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"issueLabel":       issueLabel,
	"opportunityNames": opportunityNamesJSON,
}).ParseFS(templateFS, "templates/*.html"))

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

// opportunityRef links a card to one story map it belongs to. In livt a story
// map is one opportunity, so Name is the opportunity name (= map name) and Path
// links to that map's board, relative to the page doing the rendering.
type opportunityRef struct {
	Name string
	Path string
}

// opportunityNamesJSON encodes the opportunity names as a JSON array for a
// card's data-opportunities attribute, the hook the client-side filter matches
// against. A story on no map serializes to "[]" so it matches no filter axis.
func opportunityNamesJSON(refs []opportunityRef) string {
	names := make([]string, len(refs))
	for i, r := range refs {
		names[i] = r.Name
	}
	b, err := json.Marshal(names)
	if err != nil {
		return "[]"
	}
	return string(b)
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

// taskItem is one thing the master says is not finished, lifted off its example
// mapping onto the Tasks page: an open question, or a rule with no automation
// recorded. It carries no status, and cannot carry one derived from an issue —
// the master holds automation issue URLs but not their open/closed state, and
// the build may not ask GitHub (show-automation-status-per-rule R-02). Kind
// ("question" or "rule") picks the sticky colour the item carries on its board.
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

// tasksView is the master's two unfinished flanks either side of the example
// mapping: questions close by a conversation, un-automated rules close by a
// test. One FilterOpportunities set serves both lists, since a single filter bar
// drives the page.
type tasksView struct {
	Sidebar             Sidebar
	Questions           []taskItem
	UnautomatedRules    []taskItem
	FilterOpportunities []string
}

type glossaryCard struct {
	Key        string
	Name       string
	Definition string
}

type glossaryView struct {
	Sidebar Sidebar
	Terms   []glossaryCard
}

type storyView struct {
	Story         *domain.Story
	Meta          []metaFieldView
	MappingPath   string
	Opportunities []opportunityRef
}

// termCard is a referenced ubiquitous language term rendered as a pink sticky on
// a board. Href is the link to its glossary row, or empty when the term has no
// matching ubiquitous/{key}.md file (then it renders as a plain card).
type termCard struct {
	Name string
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
