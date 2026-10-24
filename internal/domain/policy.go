package domain

// Policy is a decision that holds across stories, kept apart from the example
// mappings that log what each story decided. Ctx comes from the directory
// holding it, the way a term's does: design/policies/{ctx}/{key}.md is scoped
// to one context, design/policies/{key}.md holds across them.
type Policy struct {
	Ctx  string
	Key  string
	Name string
	// Status is agreed the way a rule's is, so it reads the same vocabulary and
	// the same default: a policy written with no status is accepted.
	Status RuleStatus
	// SupersededBy names what took a retired policy's place, as livt URIs.
	// Absent means nothing replaced it: a policy can be retired on its own.
	SupersededBy []string
	Body         string
}
