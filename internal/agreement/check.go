package agreement

import (
	"slices"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/uri"
)

// Side is one revision of a mapping and its record.
type Side struct {
	Mapping *domain.ExampleMapping
	File    File
}

// Finding is one rule that may not be accepted as the change stands.
type Finding struct {
	URI string
	// Missing are the people named who have not agreed.
	Missing []string
	// Teams are entries naming a team, which cannot be matched to a person.
	Teams []string
	// Unbacked are records in the file that neither the base nor the facts
	// stand behind.
	Unbacked []string
}

// Result is what the check found across every mapping.
type Result struct {
	// Checked counts the rules the change makes accepted that name someone.
	Checked  int
	Findings []Finding
}

// OK reports whether the change may go in.
func (r Result) OK() bool { return len(r.Findings) == 0 }

// isTeam reads the forge spelling of a team, owner/name. Matching one needs
// the forge's membership, which this package does not have.
func isTeam(who string) bool { return strings.Contains(who, "/") }

// ruleIn finds a rule by id, since a mapping may be absent on one side.
func ruleIn(m *domain.ExampleMapping, id string) (domain.Rule, bool) {
	if m == nil {
		return domain.Rule{}, false
	}
	for _, r := range m.Rules {
		if r.ID == id {
			return r, true
		}
	}
	return domain.Rule{}, false
}

func accepted(r domain.Rule) bool { return r.Status.OrDefault() == domain.RuleAccepted }

// makesAccepted reports whether the change is what puts this rule, with these
// people named on it, in the accepted state. A rule that was already there in
// the base with the same people is not the change's doing.
func makesAccepted(base domain.Rule, inBase bool, head domain.Rule) bool {
	if !accepted(head) || len(head.DecisionMakers) == 0 {
		return false
	}
	return !inBase || !accepted(base) || !slices.Equal(base.DecisionMakers, head.DecisionMakers)
}

// backed reports whether a record is one the facts would have produced.
func backed(r Record, facts Facts) bool {
	if r.In != facts.Change {
		return false
	}
	switch r.As {
	case AsAuthor:
		return same(r.Who, facts.Author)
	case AsApprover:
		return contains(facts.ApprovedBy, r.Who)
	}
	return false
}

func hasRecord(records []Record, r Record) bool {
	return slices.ContainsFunc(records, func(o Record) bool { return o == r })
}

// agreed lists who stands behind the rule: everyone on record that the base
// or the facts back, the author, and the approvers. A record only counts when
// something other than the change's own diff vouches for it.
func agreed(ruleURI string, base, head File, facts Facts) (who []string, unbacked []string) {
	baseRecords := base.records(ruleURI)
	for _, r := range head.records(ruleURI) {
		if hasRecord(baseRecords, r) || backed(r, facts) {
			who = append(who, r.Who)
		} else {
			unbacked = append(unbacked, r.Who)
		}
	}
	if facts.Author != "" {
		who = append(who, facts.Author)
	}
	return append(who, facts.ApprovedBy...), unbacked
}

// Check holds every rule the change makes accepted to the people it names.
func Check(base, head map[string]Side, facts Facts) Result {
	var res Result
	for _, key := range sortedKeys(head) {
		h := head[key]
		b := base[key]
		for _, rule := range h.Mapping.Rules {
			baseRule, inBase := ruleIn(b.Mapping, rule.ID)
			if !makesAccepted(baseRule, inBase, rule) {
				continue
			}
			res.Checked++
			ruleURI := uri.Rule(key, rule.ID)
			who, unbacked := agreed(ruleURI, b.File, h.File, facts)
			f := Finding{URI: ruleURI, Unbacked: unbacked}
			for _, dm := range rule.DecisionMakers {
				switch {
				case isTeam(dm):
					f.Teams = append(f.Teams, dm)
				case !contains(who, dm):
					f.Missing = append(f.Missing, dm)
				}
			}
			if len(f.Missing)+len(f.Teams)+len(f.Unbacked) > 0 {
				res.Findings = append(res.Findings, f)
			}
		}
	}
	return res
}

func sortedKeys(m map[string]Side) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
