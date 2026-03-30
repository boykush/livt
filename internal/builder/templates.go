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
<title>livt - Example Mappings</title>
<script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100 min-h-screen p-10">
<h1 class="text-2xl font-bold text-gray-800 mb-6">Example Mappings</h1>
<div class="flex flex-wrap gap-4">
{{range .}}
  <a href="{{.Path}}" class="block">
    <div class="bg-yellow-100 border-l-4 border-yellow-400 p-5 rounded shadow hover:-translate-y-0.5 transition-transform min-w-[200px]">
      <span class="font-semibold text-gray-800">{{.StoryName}}</span>
    </div>
  </a>
{{else}}
  <p class="text-gray-500">No example mappings found.</p>
{{end}}
</div>
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
<a href="../index.html" class="text-sm text-gray-500 hover:text-gray-700">← Back</a>

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

type mappingView struct {
	StoryName string
	Mapping   *domain.ExampleMapping
}

func renderIndex(w io.Writer, entries []IndexEntry) error {
	return indexTmpl.Execute(w, entries)
}

func renderMapping(w io.Writer, em *domain.ExampleMapping, storyName string) error {
	return mappingTmpl.Execute(w, mappingView{StoryName: storyName, Mapping: em})
}
