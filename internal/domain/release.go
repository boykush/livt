package domain

import "fmt"

type Release struct {
	Name    string
	Stories []StoryCard
}

func (r Release) DisplayName(index int) string {
	if r.Name != "" {
		return r.Name
	}
	return fmt.Sprintf("Release %d", index+1)
}
