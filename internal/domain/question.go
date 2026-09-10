package domain

type Question struct {
	ID   string
	Text string
	// Resolutions are the Issue or PR URLs where the question is being settled.
	// Named for the role, not the medium as Rule.Issues is: a rule's automation
	// always lands in an Issue (livt://mapping/trace-rule-to-tests/rule/R-01),
	// while a question can be settled in a PR or another tracker entirely. Only
	// the links are recorded; whether they are closed lives at the target.
	Resolutions []string
	// Retired records that the question is no longer asked — answered, or made
	// moot. Kept for the same reason as Rule.Retired: its id stays taken, and
	// the text stays readable. Absent means still open.
	Retired bool
	// SupersededBy as on Rule. A settled question points at the rule that
	// settled it — that rule is where the answer landed, since a Question card
	// never carries one itself.
	SupersededBy []string
}
