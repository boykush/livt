package domain

type Question struct {
	ID   string
	Text string
	// Retired records that the question is no longer asked — answered, or made
	// moot. Kept for the same reason as Rule.Retired: its id stays taken, and
	// the text stays readable. Absent means still open.
	Retired bool
}
