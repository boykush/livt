package domain

type ExampleMapping struct {
	StoryKey   StoryKey
	Rules      []Rule
	Questions  []Question
	Ubiquitous []string
}

// Active returns the mapping without its retired rules, examples, and
// questions — what the spec still asks for. Retired items stay in the livt repository so
// their ids stay taken and their text stays readable, but they are no longer on
// a board and no longer anything the livt repository calls unfinished.
func (em *ExampleMapping) Active() *ExampleMapping {
	out := &ExampleMapping{StoryKey: em.StoryKey, Ubiquitous: em.Ubiquitous}
	for _, r := range em.Rules {
		if r.Retired {
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
