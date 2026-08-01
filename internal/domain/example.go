package domain

type Example struct {
	ID   string
	Name string
	// Retired records that the example no longer illustrates the rule, kept for
	// the same reason as Rule.Retired: its id stays taken. Absent means live.
	Retired bool
}
