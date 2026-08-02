package domain

type Story struct {
	Key  StoryKey
	Name string
	// Persona is the key of the actor the story is written for, empty when the
	// story names none. It is a field of its own rather than a MetaField because
	// it addresses a persona the livt repository holds — the body's "as a ..."
	// line is prose, and prose is what drifts into a second name for one actor.
	Persona string
	Body    string
	Meta    []MetaField
}

// MetaField is one frontmatter entry beyond the reserved name field, kept in
// source order. Value holds the scalar value; sequence values are comma-joined.
type MetaField struct {
	Key   string
	Value string
}
