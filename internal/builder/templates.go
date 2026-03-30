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

type IndexEntry struct {
	StoryKey  string
	StoryName string
	Path      string
}

type storyView struct {
	Story       *domain.Story
	MappingPath string
}

type mappingView struct {
	StoryName string
	Mapping   *domain.ExampleMapping
}

func renderIndex(w io.Writer, entries []IndexEntry) error {
	return indexTmpl.Execute(w, entries)
}

func renderStory(w io.Writer, story *domain.Story, mappingPath string) error {
	return storyTmpl.Execute(w, storyView{Story: story, MappingPath: mappingPath})
}

func renderMapping(w io.Writer, em *domain.ExampleMapping, storyName string) error {
	return mappingTmpl.Execute(w, mappingView{StoryName: storyName, Mapping: em})
}
