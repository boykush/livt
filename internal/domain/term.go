package domain

// Term is a single ubiquitous language entry: Key from the filename, Name the
// display term, Body the Markdown definition. Ctx comes from the directory:
// design/ubiquitous/{ctx}/{key}.md is scoped to one context, and
// design/ubiquitous/{key}.md holds across them. The directory keeps the same
// Key at the root and under a context apart without livt having to police it.
type Term struct {
	Ctx  string
	Key  string
	Name string
	Body string
}
