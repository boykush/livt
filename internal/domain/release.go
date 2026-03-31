package domain

type Release struct {
	Name    string
	Stories []StoryKey
}
