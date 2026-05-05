package domain

type StoryCard struct {
	Key  StoryKey
	Name string
}

func (s StoryCard) HasKey() bool {
	return s.Key.Value != ""
}
