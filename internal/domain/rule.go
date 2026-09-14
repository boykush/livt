package domain

type Rule struct {
	ID       string
	Name     string
	Examples []Example
	// Status is where the rule stands: agreed, awaiting agreement, or closed
	// one of the two ways. The zero value is accepted, as every rule written
	// before the field was; Proposed is what the surfaces branch on.
	Status RuleStatus
	// Issues are the rule's automation Issue URLs on implementation repos.
	// The livt repository records the links; their state lives at the URL target.
	Issues []string
	// Automated records the judgment that the rule is actually automated,
	// which is independent of Issues being filed or closed. Absent means
	// not automated.
	Automated bool
	// SupersededBy names what took a closed rule's place, as livt URIs, so a
	// reference landing on it can go forward. Only the pointer is structured:
	// why the rule closed belongs to the commit that closed it, and a copy of
	// that reasoning here would drift from it. Absent means nothing replaced
	// it.
	SupersededBy []string
}

// Proposed reports whether the rule is still waiting to be agreed — the one
// standing the board, the Tasks page and the sidebar counts all branch on.
func (r Rule) Proposed() bool {
	return r.Status.OrDefault() == RuleProposed
}
