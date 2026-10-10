package agreement

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/uri"
)

// Written says what a run changed in the working tree.
type Written struct {
	// Recorded are the rules that gained someone's agreement on record.
	Recorded []string
	// Accepted are the proposed rules rewritten to accepted.
	Accepted []string
}

// stands reports how a person stands behind the change, if they do.
func stands(who string, facts Facts) (as string, ok bool) {
	switch {
	case same(who, facts.Author):
		return AsAuthor, true
	case contains(facts.ApprovedBy, who):
		return AsApprover, true
	}
	return "", false
}

func onRecord(records []Record, who string) bool {
	for _, r := range records {
		if same(r.Who, who) {
			return true
		}
	}
	return false
}

// Write records who stands behind each rule the change touches, and turns a
// proposed rule accepted once everyone it names is on record. Only rules the
// change touches are recorded: an approval of a change says nothing about the
// rules it left alone.
func Write(root, mappingsDir string, base, head map[string]Side, facts Facts) (Written, error) {
	var out Written
	for _, key := range sortedKeys(head) {
		h := head[key]
		b := base[key]
		file := h.File
		changed := false
		var accept []string
		for _, rule := range h.Mapping.Rules {
			if len(rule.DecisionMakers) == 0 || !rule.Status.Active() {
				continue
			}
			baseRule, inBase := ruleIn(b.Mapping, rule.ID)
			if !touched(baseRule, inBase, rule) {
				continue
			}
			ruleURI := uri.Rule(key, rule.ID)
			all := true
			recorded := false
			for _, dm := range rule.DecisionMakers {
				if onRecord(file.records(ruleURI), dm) {
					continue
				}
				as, ok := stands(dm, facts)
				if !ok || isTeam(dm) {
					all = false
					continue
				}
				file.add(ruleURI, Record{Who: dm, As: as, In: facts.Change})
				recorded = true
			}
			if recorded {
				changed = true
				out.Recorded = append(out.Recorded, ruleURI)
			}
			if all && rule.Status.OrDefault() == domain.RuleProposed {
				accept = append(accept, rule.ID)
				out.Accepted = append(out.Accepted, ruleURI)
			}
		}
		if changed {
			if err := writeFile(root, key, file); err != nil {
				return out, err
			}
		}
		if len(accept) > 0 {
			if err := acceptRules(filepath.Join(mappingsDir, key+".yaml"), accept); err != nil {
				return out, err
			}
		}
	}
	return out, nil
}

var (
	ruleStart     = regexp.MustCompile(`^(\s*)- id:\s*(\S+)\s*$`)
	statusLine    = regexp.MustCompile(`^(\s+status:\s*)proposed(\s*)$`)
	topLevelEntry = regexp.MustCompile(`^\S`)
)

// acceptRules changes each rule's status line and nothing else, so the file
// keeps its comments and layout and the diff is the one line the acceptance
// is. Editing lines rather than re-emitting the YAML is what guarantees that.
func acceptRules(path string, ids []string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	inRules := false
	ruleIndent := -1
	current := ""
	done := map[string]bool{}
	for i, line := range lines {
		if topLevelEntry.MatchString(line) {
			inRules = strings.HasPrefix(line, "rules:")
			ruleIndent, current = -1, ""
			continue
		}
		if !inRules {
			continue
		}
		if m := ruleStart.FindStringSubmatch(line); m != nil {
			if ruleIndent == -1 {
				ruleIndent = len(m[1])
			}
			if len(m[1]) == ruleIndent {
				current = m[2]
				continue
			}
		}
		if current == "" || done[current] || !wanted(ids, current) {
			continue
		}
		if m := statusLine.FindStringSubmatch(line); m != nil {
			lines[i] = m[1] + string(domain.RuleAccepted) + m[2]
			done[current] = true
		}
	}
	for _, id := range ids {
		if !done[id] {
			return fmt.Errorf("%s: rule %s has no \"status: proposed\" line to rewrite", path, id)
		}
	}
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644)
}

func wanted(ids []string, id string) bool {
	for _, w := range ids {
		if w == id {
			return true
		}
	}
	return false
}

// add puts a record under its rule, creating the entry on first use.
func (f *File) add(ruleURI string, r Record) {
	for i := range f.Agreements {
		if f.Agreements[i].URI == ruleURI {
			f.Agreements[i].AcceptedBy = append(f.Agreements[i].AcceptedBy, r)
			return
		}
	}
	f.Agreements = append(f.Agreements, Entry{URI: ruleURI, AcceptedBy: []Record{r}})
}

func writeFile(root, storyKey string, f File) error {
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	p := recordPath(root, storyKey)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, append(data, '\n'), 0o644)
}

// touched reports whether the change alters the rule at all, which is what
// makes an approval of the change an approval of this rule.
func touched(base domain.Rule, inBase bool, head domain.Rule) bool {
	if !inBase {
		return true
	}
	if base.Name != head.Name || base.Status.OrDefault() != head.Status.OrDefault() ||
		!slices.Equal(base.DecisionMakers, head.DecisionMakers) {
		return true
	}
	return !slices.EqualFunc(base.Examples, head.Examples, func(a, b domain.Example) bool {
		return a.ID == b.ID && a.Name == b.Name && a.Retired == b.Retired
	})
}
