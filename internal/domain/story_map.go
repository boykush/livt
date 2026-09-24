package domain

// StoryMap is one board. OpportunityKey comes from the filename and is what
// joins the map to the opportunity it serves — the same filename join that
// ties an example mapping to its story — and, since an opportunity has one map,
// what addresses it. Name is the board's own, and only what a page shows.
type StoryMap struct {
	OpportunityKey OpportunityKey
	Name           string
	Activities     []Activity
	Releases       []Release
	Ubiquitous     []string
}

// DisplayName is what a page shows for the map, falling back to the key: a map
// need not carry a name now that the key addresses it, and a blank board title
// leaves the reader nothing to read.
func (sm *StoryMap) DisplayName() string {
	if sm.Name != "" {
		return sm.Name
	}
	return sm.OpportunityKey.Value
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
