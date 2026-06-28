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

func TestGetExampleMappingReturnsRulesExamplesQuestions(t *testing.T) {
	s := newTestServer(t)

	_, out, err := s.getExampleMapping(context.Background(), nil, getExampleMappingInput{StoryKey: "demo"})
	if err != nil {
		t.Fatalf("getExampleMapping: %v", err)
	}

	if out.Mapping.StoryKey != "demo" {
		t.Errorf("story_key = %q, want demo", out.Mapping.StoryKey)
	}
	if len(out.Mapping.Rules) != 1 || out.Mapping.Rules[0].ID != "R-01" {
		t.Fatalf("rules = %+v, want one R-01", out.Mapping.Rules)
	}
	if len(out.Mapping.Rules[0].Examples) != 1 || out.Mapping.Rules[0].Examples[0].ID != "EX-01" {
		t.Errorf("examples = %+v, want one EX-01", out.Mapping.Rules[0].Examples)
	}
	if len(out.Mapping.Questions) != 1 || out.Mapping.Questions[0].ID != "Q-01" {
		t.Errorf("questions = %+v, want one Q-01", out.Mapping.Questions)
	}
}

func TestGetExampleMappingUnknownStoryErrors(t *testing.T) {
	s := newTestServer(t)
	if _, _, err := s.getExampleMapping(context.Background(), nil, getExampleMappingInput{StoryKey: "nope"}); err == nil {
		t.Fatal("expected error for unknown story key")
	}
}

func TestGetRuleReturnsSingleRule(t *testing.T) {
	s := newTestServer(t)

	_, out, err := s.getRule(context.Background(), nil, getRuleInput{StoryKey: "demo", RuleID: "R-01"})
	if err != nil {
		t.Fatalf("getRule: %v", err)
	}
	if out.Rule.ID != "R-01" || out.Rule.Name != "ルール1" {
		t.Errorf("rule = %+v, want R-01 ルール1", out.Rule)
	}
}

func TestGetRuleUnknownRuleErrors(t *testing.T) {
	s := newTestServer(t)
	if _, _, err := s.getRule(context.Background(), nil, getRuleInput{StoryKey: "demo", RuleID: "R-99"}); err == nil {
		t.Fatal("expected error for unknown rule id")
	}
}

func TestListStoriesFlagsExampleMapping(t *testing.T) {
	s := newTestServer(t)

	_, out, err := s.listStories(context.Background(), nil, listStoriesInput{})
	if err != nil {
		t.Fatalf("listStories: %v", err)
	}

	flags := map[string]bool{}
	for _, st := range out.Stories {
		flags[st.Key] = st.HasExampleMapping
	}
	if got, ok := flags["demo"]; !ok || !got {
		t.Errorf("demo has_example_mapping = %v (present=%v), want true", got, ok)
	}
	if got, ok := flags["other"]; !ok || got {
		t.Errorf("other has_example_mapping = %v (present=%v), want false", got, ok)
	}
}

func TestServerRegistersToolsWithoutPanic(t *testing.T) {
	s := newTestServer(t)
	if s.mcpServer() == nil {
		t.Fatal("mcpServer returned nil")
	}
}
