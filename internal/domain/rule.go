package domain

type Rule struct {
	ID       string
	Name     string
	Examples []Example
	// Issues are the rule's automation Issue URLs on implementation repos.
	// The master records the links; their state lives at the URL target.
	Issues []string
	// Automated records the judgment that the rule is actually automated,
	// which is independent of Issues being filed or closed. Absent means
	// not automated.
	Automated bool
}
