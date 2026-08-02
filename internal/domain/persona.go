package domain

// Persona is one actor the livt repository's stories are written for. Key is
// derived from the filename, Name is the display name a story's "as a ..." line
// reads, and Body is the Markdown description.
//
// A persona carries no context the way a Term does: a term means different
// things to different parts of the domain, while an actor named twice is the
// same actor, and cutting them by directory would invite the drift the persona
// list exists to stop.
type Persona struct {
	Key  string
	Name string
	Body string
}
