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
	// SupersededBy names what took a retired rule's place, as livt URIs, so a
	// reference landing on it can go forward. Only the pointer is structured:
	// why the rule was retired belongs to the commit that retired it, and a
	// copy of that reasoning here would drift from it. Absent means nothing
	// replaced it.
	SupersededBy []string
}
