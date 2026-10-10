package agreement

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-03
// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-03/example/EX-01
// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-03/example/EX-05
// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-03/example/EX-06
// Writing records who stands behind a rule while it is still proposed: the
// rule's URI, the person, how they stand behind it and the change it was in,
// in the one file kept for that mapping.
func TestWriteRecordsWhoAgreedWhileStillProposed(t *testing.T) {
	base := sides(t, noMapping, "")
	root, head := tree(t, rule("proposed", "@alice", "@bob", "@carol"), "")
	mappings := filepath.Join(root, filepath.FromSlash(MappingsDir))

	written, err := Write(root, mappings, base, head, Facts{Change: pr1, Author: "@alice", ApprovedBy: []string{"@bob"}})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(written.Recorded, []string{ruleURI}) || len(written.Accepted) != 0 {
		t.Fatalf("written = %+v, want R-01 recorded and still proposed", written)
	}

	data, err := os.ReadFile(filepath.Join(root, "agreements", "checkout.json"))
	if err != nil {
		t.Fatal(err)
	}
	var got File
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	want := []Record{{Who: "@alice", As: AsAuthor, In: pr1}, {Who: "@bob", As: AsApprover, In: pr1}}
	if len(got.Agreements) != 1 || got.Agreements[0].URI != ruleURI || !slices.Equal(got.Agreements[0].AcceptedBy, want) {
		t.Errorf("record = %+v, want %s agreed by %+v", got, ruleURI, want)
	}
}

// Recording is for the rules a change touches. An approval of a change that
// leaves a rule alone says nothing about that rule.
func TestWriteLeavesUntouchedRulesOffTheRecord(t *testing.T) {
	proposal := rule("proposed", "@bob")
	base := sides(t, proposal, "")
	root, head := tree(t, proposal, "")

	written, err := Write(root, filepath.Join(root, filepath.FromSlash(MappingsDir)), base, head, Facts{Change: pr1, ApprovedBy: []string{"@bob"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(written.Recorded) != 0 {
		t.Errorf("recorded %v, want nothing: the change did not touch the rule", written.Recorded)
	}
}

// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-04/example/EX-01
// livt:automates livt://mapping/require-decision-makers-approval-to-accept/rule/R-04/example/EX-02
// Once everyone a proposal names stands behind it, writing turns it accepted
// by changing its status line and nothing else — and what is left passes the
// check that would have refused it a moment before.
func TestWriteAcceptsAProposalEveryoneAgreedTo(t *testing.T) {
	base := sides(t, noMapping, "")
	mapping := "# kept as written\n" + rule("proposed", "@alice", "@bob") + "    examples:\n      - id: EX-01\n        name: an example\n"
	root, head := tree(t, mapping, "")
	facts := Facts{Change: pr1, Author: "@alice", ApprovedBy: []string{"@bob"}}

	written, err := Write(root, filepath.Join(root, filepath.FromSlash(MappingsDir)), base, head, facts)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(written.Accepted, []string{ruleURI}) {
		t.Fatalf("accepted = %v, want R-01", written.Accepted)
	}

	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(MappingsDir), "checkout.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if want := strings.Replace(mapping, "status: proposed", "status: accepted", 1); string(data) != want {
		t.Errorf("mapping = %q, want only the status line changed:\n%q", data, want)
	}

	after, err := LoadTree(root)
	if err != nil {
		t.Fatal(err)
	}
	if res := Check(base, after, facts); !res.OK() || res.Checked != 1 {
		t.Errorf("result = %+v, want the rewritten rule to pass its own check", res)
	}
}

// A proposal still short of someone stays a proposal.
func TestWriteLeavesAProposalSomeoneHasNotAgreedTo(t *testing.T) {
	base := sides(t, noMapping, "")
	root, head := tree(t, rule("proposed", "@alice", "@bob"), "")

	written, err := Write(root, filepath.Join(root, filepath.FromSlash(MappingsDir)), base, head, Facts{Change: pr1, Author: "@alice"})
	if err != nil {
		t.Fatal(err)
	}
	if len(written.Accepted) != 0 {
		t.Errorf("accepted %v, want none while @bob has not agreed", written.Accepted)
	}
}
