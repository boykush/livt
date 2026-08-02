package domain

type Rule struct {
	ID       string
	Name     string
	Examples []Example
	// Issues are the rule's automation Issue URLs on implementation repos.
	// The livt repository records the links; their state lives at the URL target.
	Issues []string
	// Automated records the judgment that the rule is actually automated,
	// which is independent of Issues being filed or closed. Absent means
	// not automated.
	Automated bool
	// Retired records that the rule is no longer part of the spec. It stays in
	// the livt repository rather than being deleted so its id is never handed to another
	// rule — references to it keep resolving, as retired, instead of quietly
	// pointing at whatever took the id. Absent means live.
	Retired bool
}
