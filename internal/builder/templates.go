package builder

import (
	"html/template"
	"io"

	"github.com/boykush/livt/internal/domain"
)

var indexTmpl = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>livt - Stories</title>
<script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100 min-h-screen p-10">
<h1 class="text-2xl font-bold text-gray-800 mb-6">Stories</h1>
<ul class="space-y-2">
{{range .}}
  <li><a href="{{.Path}}" class="text-blue-600 hover:underline">{{.StoryName}}</a></li>
{{else}}
  <p class="text-gray-500">No stories found.</p>
{{end}}
</ul>
</body>
</html>`))

var storyTmpl = template.Must(template.New("story").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Story.Name}} - livt</title>
<script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100 min-h-screen p-10">
<a href="../index.html" class="text-sm text-gray-500 hover:text-gray-700">&larr; Back</a>

<h1 class="text-2xl font-bold text-gray-800 mt-4 mb-4">{{.Story.Name}}</h1>

<pre class="text-gray-700 whitespace-pre-wrap mb-6">{{.Story.Body}}</pre>

{{if .MappingPath}}
  <a href="{{.MappingPath}}" class="text-blue-600 hover:underline">Example Mapping</a>
{{end}}
</body>
</html>`))

var mappingTmpl = template.Must(template.New("mapping").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.StoryName}} - livt</title>
<script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100 min-h-screen p-10">
<a href="../index.html" class="text-sm text-gray-500 hover:text-gray-700">&larr; Back</a>

<div class="bg-yellow-100 border-l-4 border-yellow-400 p-4 rounded shadow text-center text-lg font-bold text-gray-800 mt-4 mb-6">
  {{.StoryName}}
</div>

<div class="flex gap-6 items-start flex-wrap">
{{range .Mapping.Rules}}
  <div class="flex flex-col gap-2 min-w-[200px] max-w-[260px]">
    <div class="bg-blue-100 border-l-4 border-blue-400 p-4 rounded shadow font-semibold text-gray-800">
      {{.Name}}
    </div>
    {{range .Examples}}
    <div class="bg-green-100 border-l-4 border-green-400 p-4 rounded shadow text-gray-700 text-sm">
      {{.Name}}
    </div>
    {{end}}
  </div>
{{end}}
{{if .Mapping.Questions}}
  <div class="flex flex-col gap-2 min-w-[200px] max-w-[260px]">
    <div class="bg-red-50 border-l-4 border-red-300 p-3 rounded text-xs font-semibold text-red-400 uppercase tracking-wide">
      Questions
    </div>
    {{range .Mapping.Questions}}
    <div class="bg-red-100 border-l-4 border-red-400 p-4 rounded shadow text-gray-700 text-sm">
      {{.Text}}
    </div>
    {{end}}
  </div>
{{end}}
</div>
</body>
</html>`))

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
