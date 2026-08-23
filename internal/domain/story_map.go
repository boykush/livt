package domain

// StoryMap is one board. Key comes from the filename and is what joins the map
// to the opportunity it serves — the same filename join that ties an example
// mapping to its story. Name is the board's own display name, and is what
// addresses it.
type StoryMap struct {
	Key        string
	Name       string
	Activities []Activity
	Releases   []Release
	Ubiquitous []string
}

type Activity struct {
	Key   string
	Name  string
	Steps []Step
}

type Step struct {
	Key     string
	Name    string
	Stories []StoryCard
}
