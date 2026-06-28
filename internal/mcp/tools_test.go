package mcp

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// newTestServer lays out a master under a temp root with one mapped story
// (demo, which has an example mapping) and one unmapped story (other).
func newTestServer(t *testing.T) *Server {
	t.Helper()
	root := t.TempDir()

	writeFile(t, filepath.Join(root, "discoveries", "example-mappings", "demo.yaml"),
		"rules:\n"+
			"  - id: R-01\n"+
			"    name: ルール1\n"+
			"    examples:\n"+
			"      - id: EX-01\n"+
			"        name: 実例1\n"+
			"questions:\n"+
			"  - id: Q-01\n"+
			"    text: 質問1\n"+
			"ubiquitous:\n"+
			"  - story\n")
	writeFile(t, filepath.Join(root, "stories", "demo.md"), "---\nname: デモストーリー\n---\n\n本文\n")
	writeFile(t, filepath.Join(root, "stories", "other.md"), "---\nname: 別ストーリー\n---\n\n本文\n")

	return NewServer(Config{Root: root}, "test")
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestExampleMappingReturnsRulesExamplesQuestions(t *testing.T) {
	em, err := newTestServer(t).cfg.exampleMapping("demo")
	if err != nil {
		t.Fatalf("exampleMapping: %v", err)
	}
	if em.StoryKey.Value != "demo" {
		t.Errorf("story key = %q, want demo", em.StoryKey.Value)
	}
	if len(em.Rules) != 1 || em.Rules[0].ID != "R-01" {
		t.Fatalf("rules = %+v, want one R-01", em.Rules)
	}
	if len(em.Rules[0].Examples) != 1 || em.Rules[0].Examples[0].ID != "EX-01" {
		t.Errorf("examples = %+v, want one EX-01", em.Rules[0].Examples)
	}
	if len(em.Questions) != 1 || em.Questions[0].ID != "Q-01" {
		t.Errorf("questions = %+v, want one Q-01", em.Questions)
	}
}

func TestExampleMappingUnknownStoryErrors(t *testing.T) {
	if _, err := newTestServer(t).cfg.exampleMapping("nope"); err == nil {
		t.Fatal("expected error for unknown story key")
	}
}

func TestExampleMappingRejectsTraversalKey(t *testing.T) {
	if _, err := newTestServer(t).cfg.exampleMapping("../../etc/passwd"); err == nil {
		t.Fatal("expected error for traversal key")
	}
}

func TestRuleReturnsSingleRule(t *testing.T) {
	rule, err := newTestServer(t).cfg.rule("demo", "R-01")
	if err != nil {
		t.Fatalf("rule: %v", err)
	}
	if rule.ID != "R-01" || rule.Name != "ルール1" {
		t.Errorf("rule = %+v, want R-01 ルール1", rule)
	}
}

func TestRuleUnknownRuleErrors(t *testing.T) {
	if _, err := newTestServer(t).cfg.rule("demo", "R-99"); err == nil {
		t.Fatal("expected error for unknown rule id")
	}
}

func TestListStoriesLinksExampleMapping(t *testing.T) {
	s := newTestServer(t)

	_, out, err := s.listStories(context.Background(), nil, listStoriesInput{})
	if err != nil {
		t.Fatalf("listStories: %v", err)
	}

	uris := map[string]string{}
	present := map[string]bool{}
	for _, st := range out.Stories {
		uris[st.Key] = st.ExampleMappingURI
		present[st.Key] = true
	}
	if got := uris["demo"]; got != "livt://mapping/demo" {
		t.Errorf("demo example_mapping_uri = %q, want livt://mapping/demo", got)
	}
	if !present["other"] {
		t.Fatal("story other missing from list")
	}
	if got := uris["other"]; got != "" {
		t.Errorf("other example_mapping_uri = %q, want empty (no mapping)", got)
	}
}

func TestServerRegistersToolsWithoutPanic(t *testing.T) {
	if newTestServer(t).mcpServer() == nil {
		t.Fatal("mcpServer returned nil")
	}
}
