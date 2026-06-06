package domain

type StoryMap struct {
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
