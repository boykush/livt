package domain

import (
	"slices"
	"strings"
)

// RuleStatus is where a rule stands in being agreed: an ADR's status, narrowed
// to the two a rule moves between. Rejection is Retired's business, so a
// rejected proposal keeps RuleProposed beside it and reads as a proposal that
// never became spec.
type RuleStatus string

const (
	RuleProposed RuleStatus = "proposed"
	RuleAccepted RuleStatus = "accepted"
)

// RuleStatusDefault is what a rule with no status written means: every rule on
// file before the field existed was agreed when it went on the board.
const RuleStatusDefault = RuleAccepted

// RuleStatuses are the statuses a mapping may write, in the order an error
// message lists them.
var RuleStatuses = []RuleStatus{RuleProposed, RuleAccepted}

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

// RuleStatusList names every status a mapping may write, for an error telling
// the author what to pick instead.
func RuleStatusList() string {
	names := make([]string, len(RuleStatuses))
	for i, s := range RuleStatuses {
		names[i] = string(s)
	}
	return strings.Join(names, ", ")
}
