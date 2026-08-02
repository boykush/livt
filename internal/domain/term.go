package domain

// Term is a single ubiquitous language entry. Key is derived from the filename,
// Name is the display term, and Body is the Markdown definition.
//
// Ctx is the context the term belongs to, taken from the directory holding it:
// ubiquitous/{ctx}/{key}.md is scoped to one context, ubiquitous/{key}.md holds
// across them and leaves Ctx empty. Cutting the context by directory is what
// makes the pair unique — the same Key can mean one thing at the root and
// another under a context, and the filesystem keeps the two files apart without
// livt having to police it.
type Term struct {
	Ctx  string
	Key  string
	Name string
	Body string
}
