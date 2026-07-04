package builder

import (
	"embed"
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
	"issueLabel": issueLabel,
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

// Sidebar is the shared navigation rendered on every hub page. Prefix is the
// relative path back to the output root ("" for root pages); Active marks the
// current resource type.
type Sidebar struct {
	Prefix    string
	Active    string
	Mappings  int
	StoryMaps int
	Stories   int
	Terms     int
}

type mappingTile struct {
	Key       string
	StoryName string
}

type mappingsIndexView struct {
	Sidebar  Sidebar
	Mappings []mappingTile
}

type storyMapTile struct {
	Name string
}

type storyMapsIndexView struct {
	Sidebar   Sidebar
	StoryMaps []storyMapTile
}

type storyItem struct {
	Key          string
	Name         string
	StoryMapPath string
	MappingPath  string
	Links        []metaFieldView
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
	Sidebar Sidebar
	Stories []storyItem
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
	Story        *domain.Story
	Meta         []metaFieldView
	MappingPath  string
	StoryMapPath string
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

func renderMappingsIndex(w io.Writer, view mappingsIndexView) error {
	return tmpl.ExecuteTemplate(w, "index.html", view)
}

func renderStoryMapsIndex(w io.Writer, view storyMapsIndexView) error {
	return tmpl.ExecuteTemplate(w, "story-maps.html", view)
}

func renderStoriesIndex(w io.Writer, view storiesIndexView) error {
	return tmpl.ExecuteTemplate(w, "stories.html", view)
}

func renderStory(w io.Writer, story *domain.Story, mappingPath, storyMapPath string) error {
	return tmpl.ExecuteTemplate(w, "story.html", storyView{
		Story:        story,
		Meta:         metaFieldViews(story.Meta),
		MappingPath:  mappingPath,
		StoryMapPath: storyMapPath,
	})
}

func renderMapping(w io.Writer, em *domain.ExampleMapping, storyName, storyPath string, ubiquitous []termCard) error {
	return tmpl.ExecuteTemplate(w, "mapping.html", mappingView{StoryName: storyName, StoryPath: storyPath, Mapping: em, Ubiquitous: ubiquitous})
}

func renderStoryMap(w io.Writer, view storyMapView) error {
	return tmpl.ExecuteTemplate(w, "story_map.html", view)
}

func renderGlossary(w io.Writer, view glossaryView) error {
	return tmpl.ExecuteTemplate(w, "glossary.html", view)
}
