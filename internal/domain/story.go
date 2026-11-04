package domain

type Story struct {
	Key  StoryKey
	Name string
	// Opportunity is the one the card names in `opportunity:`, empty when it
	// names none. A story on a map belongs to that map's opportunity without
	// it, so this is for a story no map places.
	Opportunity OpportunityKey
	Body        string
	Meta        []MetaField
}

// DisplayName is what a page shows for the story. Name is free text and a story
// file need not carry one, so it falls back to the key: every surface naming a
// story is a heading, a title, or a link label, and a blank one leaves the
// reader nothing to read and a link nothing to announce.
func (s *Story) DisplayName() string {
	if s.Name != "" {
		return s.Name
	}
	return s.Key.Value
}

// MetaField is one frontmatter entry beyond the reserved name field, kept in
// source order. Value holds the scalar value; sequence values are comma-joined.
type MetaField struct {
	Key   string
	Value string
}
