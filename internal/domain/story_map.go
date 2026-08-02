package domain

// StoryMap is one opportunity's board. Personas are the actors whose journey it
// maps, declared here rather than derived from its story cards: the cast is
// settled while the journey is being framed, before any story is committed, so
// a map names an actor it has no story for yet.
type StoryMap struct {
	Name       string
	Activities []Activity
	Releases   []Release
	Personas   []string
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
