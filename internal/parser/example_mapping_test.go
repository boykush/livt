package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseExampleMappingReadsRuleIssuesAndAutomated(t *testing.T) {
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
	if !recorded.Automated {
		t.Fatal("recorded rule should be automated")
	}

	bare := em.Rules[1]
	if len(bare.Issues) != 0 || bare.Automated {
		t.Fatalf("bare rule should default to unlinked and not automated, got issues=%v automated=%v", bare.Issues, bare.Automated)
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
		"    retired: true\n" +
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

	if em.Rules[0].Retired || em.Rules[0].Examples[0].Retired || em.Questions[0].Retired {
		t.Error("items without the field should default to live")
	}
	if !em.Rules[1].Retired || em.Rules[1].Name != "退役したルール" {
		t.Errorf("rule = %+v, want R-02 retired with its text kept", em.Rules[1])
	}
	if !em.Rules[0].Examples[1].Retired {
		t.Errorf("example = %+v, want EX-02 retired", em.Rules[0].Examples[1])
	}
	if !em.Questions[1].Retired || em.Questions[1].Text != "退役した疑問" {
		t.Errorf("question = %+v, want Q-02 retired with its text kept", em.Questions[1])
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
