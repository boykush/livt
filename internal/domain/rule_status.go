package domain

import (
	"slices"
	"strings"
)

// RuleStatus is where a rule stands on the one axis an ADR's status also
// carries: proposed while it waits for agreement, accepted once it has it, and
// the two ways it closes — rejected when the proposal was turned down, retired
// when spec it once was stopped holding. A closed rule stays in the livt
// repository so its id is never handed to another rule.
type RuleStatus string

const (
	RuleProposed RuleStatus = "proposed"
	RuleAccepted RuleStatus = "accepted"
	RuleRejected RuleStatus = "rejected"
	RuleRetired  RuleStatus = "retired"
)

// RuleStatusDefault is what a rule with no status written means: every rule on
// file before the field existed was agreed when it went on the board.
const RuleStatusDefault = RuleAccepted

// RuleStatuses are the statuses a mapping may write, in the order an error
// message lists them — the order a rule moves through them.
var RuleStatuses = []RuleStatus{RuleProposed, RuleAccepted, RuleRejected, RuleRetired}

// Valid reports whether the status is one livt knows, so a mapping naming
// another is rejected at parse time rather than read as the default.
func (s RuleStatus) Valid() bool {
	return slices.Contains(RuleStatuses, s)
}

// OrDefault resolves the zero value to the status it means, so a surface
// reporting a rule's status never spells that default itself.
func (s RuleStatus) OrDefault() RuleStatus {
	if s == "" {
		return RuleStatusDefault
	}
	return s
}

// Active reports whether the rule is one the livt repository still asks about —
// a proposal awaiting agreement, or the agreement itself. It names the live
// statuses rather than the closed ones, so a status added later stays off the
// boards until it is deliberately let on.
func (s RuleStatus) Active() bool {
	switch s.OrDefault() {
	case RuleProposed, RuleAccepted:
		return true
	}
	return false
}

// RuleStatusList names every status a mapping may write, for an error telling
// the author what to pick instead.
func RuleStatusList() string {
	names := make([]string, len(RuleStatuses))
	for i, s := range RuleStatuses {
		names[i] = string(s)
	}
	return strings.Join(names, ", ")
}
