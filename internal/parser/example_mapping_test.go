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
