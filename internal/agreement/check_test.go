package agreement

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

const ruleURI = "livt://mapping/checkout/rule/R-01"

// tree writes one livt repository holding the checkout mapping, and its record
// when one is given, and reads it back the way the command does.
func tree(t *testing.T, mapping, record string) (string, map[string]Side) {
	t.Helper()
	root := t.TempDir()
	write(t, filepath.Join(root, filepath.FromSlash(MappingsDir), "checkout.yaml"), mapping)
	if record != "" {
		write(t, filepath.Join(root, Dir, "checkout.json"), record)
	}
	sides, err := LoadTree(root)
	if err != nil {
		t.Fatal(err)
	}
	return root, sides
}

func sides(t *testing.T, mapping, record string) map[string]Side {
	t.Helper()
	_, s := tree(t, mapping, record)
	return s
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// rule writes a one-rule mapping. An empty status leaves the line out, which
// reads as accepted.
func rule(status string, makers ...string) string {
	out := "rules:\n  - id: R-01\n    name: first\n"
	if status != "" {
		out += "    status: " + status + "\n"
	}
	if len(makers) > 0 {
		out += "    decision_makers: [\"" + strings.Join(makers, "\", \"") + "\"]\n"
	}
	return out
}

const noMapping = "rules: []\n"

func record(who, as, in string) string {
	return `{"agreements":[{"uri":"` + ruleURI + `","accepted_by":[{"who":"` + who + `","as":"` + as + `","in":"` + in + `"}]}]}`
}

const pr1, pr2 = "https://forge.example/acme/specs/pull/1", "https://forge.example/acme/specs/pull/2"

// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-01/example/EX-01
// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-01/example/EX-02
// An accepted rule missing even one of the people it names does not pass, and
// passes once all of them stand behind the change.
func TestCheckHoldsAnAcceptedRuleToEveryoneItNames(t *testing.T) {
	base := sides(t, rule("proposed", "@alice", "@bob"), "")
	head := sides(t, rule("accepted", "@alice", "@bob"), "")

	res := Check(base, head, Facts{Change: pr1, Author: "@alice"})
	if res.OK() || len(res.Findings) != 1 || !slices.Equal(res.Findings[0].Missing, []string{"@bob"}) {
		t.Fatalf("findings = %+v, want R-01 waiting for @bob", res.Findings)
	}

	res = Check(base, head, Facts{Change: pr1, Author: "@alice", ApprovedBy: []string{"@bob"}})
	if !res.OK() || res.Checked != 1 {
		t.Fatalf("result = %+v, want the one accepted rule to pass", res)
	}
}

// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-01/example/EX-03
// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-01/example/EX-04
// How the rule came to be accepted makes no difference: turned from a
// proposal, added already accepted, or added with no status at all.
func TestCheckTreatsEveryRouteToAcceptedAlike(t *testing.T) {
	routes := map[string][2]string{
		"proposal accepted":        {rule("proposed", "@bob"), rule("accepted", "@bob")},
		"added accepted":           {noMapping, rule("accepted", "@bob")},
		"added with no status":     {noMapping, rule("", "@bob")},
		"named on an agreed rule":  {rule(""), rule("", "@bob")},
		"another name on the rule": {rule("", "@alice"), rule("", "@alice", "@bob")},
	}
	for name, r := range routes {
		res := Check(sides(t, r[0], ""), sides(t, r[1], ""), Facts{Change: pr1, Author: "@alice"})
		if res.OK() {
			t.Errorf("%s: passed without @bob", name)
		}
		res = Check(sides(t, r[0], ""), sides(t, r[1], ""), Facts{Change: pr1, Author: "@alice", ApprovedBy: []string{"@bob"}})
		if !res.OK() {
			t.Errorf("%s: findings = %+v, want it to pass with @bob's approval", name, res.Findings)
		}
	}
}

// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-01/example/EX-05
// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-01/example/EX-06
// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-01/example/EX-09
// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-01/example/EX-10
// What the check leaves alone: a proposal going on file, a rule naming nobody,
// a rule that was already accepted before the change, and a proposal turned
// down. None of them is a rule being recorded as agreed by this change.
func TestCheckAsksNothingOfChangesThatAcceptNoNamedRule(t *testing.T) {
	changes := map[string][2]string{
		"proposal goes on file":     {noMapping, rule("proposed", "@bob")},
		"rule names nobody":         {noMapping, rule("accepted")},
		"already accepted on main":  {rule("accepted", "@bob"), rule("accepted", "@bob")},
		"proposal turned down":      {rule("proposed", "@bob"), rule("rejected", "@bob")},
		"accepted rule then closed": {rule("accepted", "@bob"), rule("retired", "@bob")},
	}
	for name, c := range changes {
		res := Check(sides(t, c[0], ""), sides(t, c[1], ""), Facts{Change: pr1})
		if !res.OK() || res.Checked != 0 {
			t.Errorf("%s: result = %+v, want nothing checked", name, res)
		}
	}
}

// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-02/example/EX-01
// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-02/example/EX-02
// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-02/example/EX-03
// A person named stands behind the change by having written it or by having
// approved it. When no person wrote it, nobody counts as its author, and the
// person a bot wrote for has to approve like anyone else.
func TestCheckCountsAuthorsAndApprovers(t *testing.T) {
	base := sides(t, noMapping, "")
	head := sides(t, rule("accepted", "@alice"), "")

	if res := Check(base, head, Facts{Change: pr1, Author: "@Alice"}); !res.OK() {
		t.Errorf("findings = %+v, want the author to count", res.Findings)
	}
	if res := Check(base, head, Facts{Change: pr1, Author: "@bob", ApprovedBy: []string{"@alice"}}); !res.OK() {
		t.Errorf("findings = %+v, want an approver to count", res.Findings)
	}
	if res := Check(base, head, Facts{Change: pr1}); res.OK() {
		t.Error("passed with no author and no approval, want @alice still awaited")
	}
}

// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-02/example/EX-05
// A team cannot be matched to a person here, and saying so beats waiting on it
// for ever: even with everyone approving, the entry is named as the problem.
func TestCheckRefusesATeam(t *testing.T) {
	base := sides(t, noMapping, "")
	head := sides(t, rule("accepted", "@acme/design"), "")

	res := Check(base, head, Facts{Change: pr1, Author: "@alice", ApprovedBy: []string{"@acme/design"}})
	if res.OK() || !slices.Equal(res.Findings[0].Teams, []string{"@acme/design"}) {
		t.Fatalf("findings = %+v, want the team entry refused", res.Findings)
	}
}

// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-03/example/EX-02
// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-03/example/EX-03
// Someone on record from an earlier change stands behind the rule in a later
// one: a proposal @bob approved before it was merged is accepted afterwards
// without asking @bob again.
func TestCheckCountsPeopleAlreadyOnRecord(t *testing.T) {
	onMain := record("@bob", AsApprover, pr1)
	base := sides(t, rule("proposed", "@alice", "@bob"), onMain)
	head := sides(t, rule("accepted", "@alice", "@bob"), onMain)

	if res := Check(base, head, Facts{Change: pr2, Author: "@alice"}); !res.OK() {
		t.Fatalf("findings = %+v, want @bob's earlier approval to carry over", res.Findings)
	}
}

// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-03/example/EX-04
// A record the change adds counts only when the facts would have produced it.
// Writing a name into the file by hand is refused, and named as the reason.
func TestCheckRefusesARecordNothingBacks(t *testing.T) {
	base := sides(t, rule("proposed", "@alice", "@bob"), "")
	head := sides(t, rule("accepted", "@alice", "@bob"), record("@bob", AsApprover, pr1))

	res := Check(base, head, Facts{Change: pr1, Author: "@alice"})
	if res.OK() || !slices.Equal(res.Findings[0].Unbacked, []string{"@bob"}) {
		t.Fatalf("findings = %+v, want @bob's record refused as unbacked", res.Findings)
	}
	if res := Check(base, head, Facts{Change: pr1, Author: "@alice", ApprovedBy: []string{"@bob"}}); !res.OK() {
		t.Fatalf("findings = %+v, want the same record to pass once @bob has approved", res.Findings)
	}
}
