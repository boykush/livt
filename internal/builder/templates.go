package builder

import (
	"embed"
	"html/template"
	"io"

	"github.com/boykush/livt/internal/domain"
)

//go:embed templates/*.html
var templateFS embed.FS

// All templates are parsed into a single set so pages can share the
// {{define "sidebar"}} partial in _sidebar.html.
var tmpl = template.Must(template.ParseFS(templateFS, "templates/*.html"))

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
	Key  string
	Name string
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
	MappingPath  string
	StoryMapPath string
}

type mappingView struct {
	StoryName string
	StoryPath string
	Mapping   *domain.ExampleMapping
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
	return tmpl.ExecuteTemplate(w, "story.html", storyView{Story: story, MappingPath: mappingPath, StoryMapPath: storyMapPath})
}

func renderMapping(w io.Writer, em *domain.ExampleMapping, storyName, storyPath string) error {
	return tmpl.ExecuteTemplate(w, "mapping.html", mappingView{StoryName: storyName, StoryPath: storyPath, Mapping: em})
}

func renderStoryMap(w io.Writer, view storyMapView) error {
	return tmpl.ExecuteTemplate(w, "story_map.html", view)
}

func renderGlossary(w io.Writer, view glossaryView) error {
	return tmpl.ExecuteTemplate(w, "glossary.html", view)
}
