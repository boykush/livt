package domain

type StoryCard struct {
	Key     StoryKey
	Name    string
	Release string
}

func (s StoryCard) HasKey() bool {
	return s.Key.Value != ""
}
