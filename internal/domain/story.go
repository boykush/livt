package domain

type Story struct {
	Key  StoryKey
	Name string
	Body string
	Meta []MetaField
}

// MetaField is one frontmatter entry beyond the reserved name field, kept in
// source order. Value holds the scalar value; sequence values are comma-joined.
type MetaField struct {
	Key   string
	Value string
}
