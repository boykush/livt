package domain

type ExampleMapping struct {
	StoryKey StoryKey
	// Name is the mapping's own, independent of its story's. A mapping with no
	// story file — a fix driven straight into tests — has nothing else to be
	// called by but its key.
	Name       string
	Rules      []Rule
	Questions  []Question
	Ubiquitous []string
}

// DisplayName is what every surface naming the board calls it: its own name,
// else its story's, which falls back to the key when there is no story either.
func (em *ExampleMapping) DisplayName(story *Story) string {
	if em.Name != "" {
		return em.Name
	}
	return story.DisplayName()
}

// Active returns the mapping without its closed rules and its retired examples
// and questions — what the spec still asks for. They stay in the livt repository so
// their ids stay taken and their text stays readable, but they are no longer on
// a board and no longer anything the livt repository calls unfinished.
func (em *ExampleMapping) Active() *ExampleMapping {
	out := &ExampleMapping{StoryKey: em.StoryKey, Name: em.Name, Ubiquitous: em.Ubiquitous}
	for _, r := range em.Rules {
		if !r.Status.Active() {
			continue
		}
		var examples []Example
		for _, ex := range r.Examples {
			if !ex.Retired {
				examples = append(examples, ex)
			}
		}
		r.Examples = examples
		out.Rules = append(out.Rules, r)
	}
	for _, q := range em.Questions {
		if q.Retired {
			continue
		}
		out.Questions = append(out.Questions, q)
	}
	return out
}
