package builder

import (
	"embed"
	"html/template"
	"io"

	"github.com/boykush/livt/internal/domain"
)

//go:embed templates/*.html
var templateFS embed.FS

var indexTmpl = template.Must(template.ParseFS(templateFS, "templates/index.html"))
var storyTmpl = template.Must(template.ParseFS(templateFS, "templates/story.html"))
var mappingTmpl = template.Must(template.ParseFS(templateFS, "templates/mapping.html"))
var storyMapTmpl = template.Must(template.ParseFS(templateFS, "templates/story_map.html"))
var glossaryTmpl = template.Must(template.ParseFS(templateFS, "templates/glossary.html"))

type IndexEntry struct {
	StoryKey  string
	StoryName string
	Path      string
}

type glossaryCard struct {
	Key        string
	Name       string
	Definition string
}

type glossaryView struct {
	Terms []glossaryCard
}

type storyView struct {
	Story        *domain.Story
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
	StoryName string
	StoryPath string
	Mapping   *domain.ExampleMapping
	Terms     []termCard
}

func renderIndex(w io.Writer, entries []IndexEntry) error {
	return indexTmpl.Execute(w, entries)
}

func renderStory(w io.Writer, story *domain.Story, mappingPath, storyMapPath string) error {
	return storyTmpl.Execute(w, storyView{Story: story, MappingPath: mappingPath, StoryMapPath: storyMapPath})
}

func renderMapping(w io.Writer, em *domain.ExampleMapping, storyName, storyPath string, terms []termCard) error {
	return mappingTmpl.Execute(w, mappingView{StoryName: storyName, StoryPath: storyPath, Mapping: em, Terms: terms})
}

func renderStoryMap(w io.Writer, view storyMapView) error {
	return storyMapTmpl.Execute(w, view)
}

func renderGlossary(w io.Writer, view glossaryView) error {
	return glossaryTmpl.Execute(w, view)
}
