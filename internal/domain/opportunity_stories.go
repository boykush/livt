package domain

// OpportunityStory is one story an opportunity took on. OnMap says the story
// is a card on the opportunity's map, which is the only place a release slice
// is drawn: Release is meaningful only then, and a story that names the
// opportunity from its own card alone has no slice to sit in.
type OpportunityStory struct {
	Key     StoryKey
	OnMap   bool
	Release string
}

// OpportunityStories is the join between opportunities and stories, derived
// once so that every surface — chips, filters, progress, MCP — reads the same
// answer however a story reached its opportunity.
type OpportunityStories struct {
	// ByOpportunity is each opportunity's stories: its map's cards in map
	// order, then the stories naming it from their cards in story order.
	ByOpportunity map[OpportunityKey][]OpportunityStory
	// ByStory is the opportunities each story belongs to: one per map it sits
	// on, in map order, then the one its card names if no map gave it already.
	ByStory map[StoryKey][]OpportunityKey
}

// DeriveOpportunityStories joins every keyed card on a map to that map's
// opportunity, and every story naming an opportunity in its card to that one.
// A story reached both ways is the opportunity's once, as its map card: the
// map is where its release was decided.
func DeriveOpportunityStories(maps []*StoryMap, stories []*Story) OpportunityStories {
	out := OpportunityStories{
		ByOpportunity: make(map[OpportunityKey][]OpportunityStory),
		ByStory:       make(map[StoryKey][]OpportunityKey),
	}
	joined := make(map[OpportunityKey]map[StoryKey]bool)
	join := func(o OpportunityKey, s OpportunityStory) {
		if joined[o] == nil {
			joined[o] = make(map[StoryKey]bool)
		}
		if joined[o][s.Key] {
			return
		}
		joined[o][s.Key] = true
		out.ByOpportunity[o] = append(out.ByOpportunity[o], s)
		out.ByStory[s.Key] = append(out.ByStory[s.Key], o)
	}
	for _, sm := range maps {
		for _, a := range sm.Activities {
			for _, st := range a.Steps {
				for _, sc := range st.Stories {
					if sc.HasKey() {
						join(sm.OpportunityKey, OpportunityStory{Key: sc.Key, OnMap: true, Release: sc.Release})
					}
				}
			}
		}
	}
	for _, s := range stories {
		if s.Opportunity.Value != "" {
			join(s.Opportunity, OpportunityStory{Key: s.Key})
		}
	}
	return out
}
