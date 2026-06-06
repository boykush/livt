package domain

// Term is a single ubiquitous language entry, stored as ubiquitous/{key}.md.
// The key is derived from the filename, Name is the display term, and Body is
// the Markdown definition.
type Term struct {
	Key  string
	Name string
	Body string
}
