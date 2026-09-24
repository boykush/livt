package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/boykush/livt/internal/domain"
)

func TestParseExampleMappingReadsRuleIssues(t *testing.T) {
	path := filepath.Join(t.TempDir(), "story.yaml")
	data := []byte("rules:\n" +
		"  - id: R-01\n" +
		"    name: 記録されたルール\n" +
		"    issues:\n" +
		"      - https://github.com/boykush/livt/issues/25\n" +
		"      - https://github.com/boykush/other/issues/7\n" +
		"    automated: true\n" +
		"  - id: R-02\n" +
		"    name: 素のルール\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	em, err := ParseExampleMapping(path)
	if err != nil {
		t.Fatal(err)
	}

	recorded := em.Rules[0]
	if len(recorded.Issues) != 2 || recorded.Issues[0] != "https://github.com/boykush/livt/issues/25" {
		t.Fatalf("got issues %v, want the two recorded URLs", recorded.Issues)
	}
	if recorded.Automated() {
		t.Fatal("automation is collected from the tests, not read from the mapping")
	}

	bare := em.Rules[1]
	if len(bare.Issues) != 0 {
		t.Fatalf("bare rule should default to unlinked, got issues=%v", bare.Issues)
	}
}

// livt://mapping/trace-test-to-rule/rule/R-05/example/EX-04: retirement is a
// field on the item, so a structural edit of the YAML cannot lose it the way it
// would lose a commented-out block — and the retired body stays readable (EX-03).
func TestParseExampleMappingReadsRetired(t *testing.T) {
	path := filepath.Join(t.TempDir(), "story.yaml")
	data := []byte("rules:\n" +
		"  - id: R-01\n" +
		"    name: 現役のルール\n" +
		"    examples:\n" +
		"      - id: EX-01\n" +
		"        name: 現役の実例\n" +
		"      - id: EX-02\n" +
		"        name: 退役した実例\n" +
		"        retired: true\n" +
		"  - id: R-02\n" +
		"    name: 退役したルール\n" +
		"    status: retired\n" +
		"questions:\n" +
		"  - id: Q-01\n" +
		"    text: 現役の疑問\n" +
		"  - id: Q-02\n" +
		"    text: 退役した疑問\n" +
		"    retired: true\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	em, err := ParseExampleMapping(path)
	if err != nil {
		t.Fatal(err)
	}

	if !em.Rules[0].Status.Active() || em.Rules[0].Examples[0].Retired || em.Questions[0].Retired {
		t.Error("items without the field should default to live")
	}
	if em.Rules[1].Status != domain.RuleRetired || em.Rules[1].Name != "退役したルール" {
		t.Errorf("rule = %+v, want R-02 retired with its text kept", em.Rules[1])
	}
	if !em.Rules[0].Examples[1].Retired {
		t.Errorf("example = %+v, want EX-02 retired", em.Rules[0].Examples[1])
	}
	if !em.Questions[1].Retired || em.Questions[1].Text != "退役した疑問" {
		t.Errorf("question = %+v, want Q-02 retired with its text kept", em.Questions[1])
	}
}

// livt://mapping/propose-rule-before-agreement/rule/R-01: the one axis carries
// all four, and a rule written without it is accepted, as every rule was before
// the field existed (EX-01, EX-02).
func TestParseExampleMappingReadsStatus(t *testing.T) {
	path := filepath.Join(t.TempDir(), "story.yaml")
	data := []byte("rules:\n" +
		"  - id: R-01\n" +
		"    name: 提案中のルール\n" +
		"    status: proposed\n" +
		"  - id: R-02\n" +
		"    name: 合意済みのルール\n" +
		"    status: accepted\n" +
		"  - id: R-03\n" +
		"    name: statusのないルール\n" +
		"  - id: R-04\n" +
		"    name: 却下された提案\n" +
		"    status: rejected\n" +
		"  - id: R-05\n" +
		"    name: 退役したルール\n" +
		"    status: retired\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	em, err := ParseExampleMapping(path)
	if err != nil {
		t.Fatal(err)
	}

	if !em.Rules[0].Proposed() {
		t.Errorf("rule = %+v, want R-01 proposed", em.Rules[0])
	}
	for _, r := range em.Rules[1:3] {
		if r.Proposed() {
			t.Errorf("rule = %+v, want %s accepted", r, r.ID)
		}
	}
	for i, want := range map[int]domain.RuleStatus{3: domain.RuleRejected, 4: domain.RuleRetired} {
		if got := em.Rules[i].Status; got != want {
			t.Errorf("rule %s status = %q, want %q", em.Rules[i].ID, got, want)
		}
	}
}

// livt://mapping/propose-rule-before-agreement/rule/R-01/example/EX-06: the old
// spelling is no longer a rule field, so the line says nothing about where the
// rule stands — status is the whole of that, and a rule still carrying retired
// reads as one written without a status.
func TestParseExampleMappingIgnoresRuleRetired(t *testing.T) {
	path := filepath.Join(t.TempDir(), "story.yaml")
	data := []byte("rules:\n" +
		"  - id: R-01\n" +
		"    name: 旧来の綴りで閉じたルール\n" +
		"    retired: true\n" +
		"  - id: R-02\n" +
		"    name: 旧来の綴りで閉じた提案\n" +
		"    status: proposed\n" +
		"    retired: true\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	em, err := ParseExampleMapping(path)
	if err != nil {
		t.Fatal(err)
	}

	if got := em.Rules[0].Status; got != domain.RuleAccepted {
		t.Errorf("rule R-01 status = %q, want %q", got, domain.RuleAccepted)
	}
	if got := em.Rules[1].Status; got != domain.RuleProposed {
		t.Errorf("rule R-02 status = %q, want %q", got, domain.RuleProposed)
	}
}

// livt://mapping/propose-rule-before-agreement/rule/R-01/example/EX-03: a status
// livt does not know fails the parse. Read as accepted, a mistyped "proposed"
// would put an unagreed rule on the board as spec.
func TestParseExampleMappingRejectsAnUnknownStatus(t *testing.T) {
	path := filepath.Join(t.TempDir(), "story.yaml")
	data := []byte("rules:\n" +
		"  - id: R-01\n" +
		"    name: 打ち間違えたルール\n" +
		"    status: propsed\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := ParseExampleMapping(path)
	if err == nil {
		t.Fatal("expected an unknown status to fail the parse")
	}
	if want := `rule "R-01": unknown status "propsed" (supported: proposed, accepted, rejected, retired)`; err.Error() != want {
		t.Errorf("error = %q, want %q", err, want)
	}
}

func TestParseExampleMappingReadsReferencedTerms(t *testing.T) {
	path := filepath.Join(t.TempDir(), "story.yaml")
	data := []byte("rules: []\nubiquitous:\n  - story-map\n  - story\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	em, err := ParseExampleMapping(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(em.Ubiquitous) != 2 || em.Ubiquitous[0] != "story-map" || em.Ubiquitous[1] != "story" {
		t.Fatalf("got ubiquitous %v, want [story-map story]", em.Ubiquitous)
	}
}

// livt://mapping/trace-test-to-rule/rule/R-09/example/EX-01 and EX-03: a
// retired rule, example, or question names where the spec went, as a list so an
// item that split into two can name both successors.
func TestParseExampleMappingReadsSupersededBy(t *testing.T) {
	path := filepath.Join(t.TempDir(), "story.yaml")
	data := []byte("rules:\n" +
		"  - id: R-01\n" +
		"    name: 分割されて退役したルール\n" +
		"    retired: true\n" +
		"    superseded_by:\n" +
		"      - livt://mapping/story/rule/R-02\n" +
		"      - livt://mapping/other-story/rule/R-07\n" +
		"  - id: R-02\n" +
		"    name: 現役のルール\n" +
		"    examples:\n" +
		"      - id: EX-01\n" +
		"        name: 差し替えられた実例\n" +
		"        retired: true\n" +
		"        superseded_by:\n" +
		"          - livt://mapping/story/rule/R-02/example/EX-02\n" +
		"      - id: EX-02\n" +
		"        name: 現役の実例\n" +
		"questions:\n" +
		"  - id: Q-01\n" +
		"    text: ルール化されて閉じた疑問\n" +
		"    retired: true\n" +
		"    superseded_by:\n" +
		"      - livt://mapping/story/rule/R-02\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	em, err := ParseExampleMapping(path)
	if err != nil {
		t.Fatal(err)
	}

	split := em.Rules[0].SupersededBy
	if len(split) != 2 || split[0] != "livt://mapping/story/rule/R-02" || split[1] != "livt://mapping/other-story/rule/R-07" {
		t.Errorf("rule superseded_by = %v, want both successors, the second in another mapping", split)
	}
	if got := em.Rules[1].Examples[0].SupersededBy; len(got) != 1 || got[0] != "livt://mapping/story/rule/R-02/example/EX-02" {
		t.Errorf("example superseded_by = %v, want the example that replaced it", got)
	}
	if got := em.Questions[0].SupersededBy; len(got) != 1 || got[0] != "livt://mapping/story/rule/R-02" {
		t.Errorf("question superseded_by = %v, want the rule that settled it", got)
	}
	if len(em.Rules[1].SupersededBy) != 0 || len(em.Rules[1].Examples[1].SupersededBy) != 0 {
		t.Error("items without the field should carry no successor")
	}
}
