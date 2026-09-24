package domain

type Example struct {
	ID   string
	Name string
	// Retired records that the example no longer illustrates the rule, kept for
	// the same reason as Rule.Retired: its id stays taken. Absent means live.
	Retired bool
	// SupersededBy as on Rule: the examples illustrating the rule in this one's
	// place, as livt URIs.
	SupersededBy []string
	// Automations are the tests citing this example, derived as on Rule.
	Automations []Automation
}

// Automated reports whether any test cites this example.
func (e Example) Automated() bool { return len(e.Automations) > 0 }
