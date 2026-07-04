package mcp

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// newTestServer lays out a master under a temp root with one mapped story
// (demo, which has an example mapping) and one unmapped story (other), one
// story map (デモマップ, holding a committed card and a bare candidate card),
// and one committed ubiquitous term (story; missing-term stays uncommitted).
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
			"  - story\n"+
			"  - missing-term\n")
	writeFile(t, filepath.Join(root, "stories", "demo.md"),
		"---\nname: デモストーリー\nissue: https://example.com/issues/1\n---\n\n本文\n")
	writeFile(t, filepath.Join(root, "stories", "other.md"), "---\nname: 別ストーリー\n---\n\n本文\n")
	writeFile(t, filepath.Join(root, "discoveries", "usm", "demo-map.yaml"),
		"name: デモマップ\n"+
			"ubiquitous:\n"+
			"  - story\n"+
			"releases:\n"+
			"  - id: mvp\n"+
			"    name: MVP\n"+
			"activities:\n"+
			"  - key: activity-1\n"+
			"    name: アクティビティ1\n"+
			"    steps:\n"+
			"      - key: step-1\n"+
			"        name: ステップ1\n"+
			"        stories:\n"+
			"          - name: デモストーリー\n"+
			"            key: demo\n"+
			"            release: mvp\n"+
			"          - name: 候補ストーリー\n")
	writeFile(t, filepath.Join(root, "ubiquitous", "story.md"), "---\nname: ストーリー\n---\n\nストーリーの定義\n")

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

	mappingURIs := map[string]string{}
	storyURIs := map[string]string{}
	for _, st := range out.Stories {
		mappingURIs[st.Key] = st.ExampleMappingURI
		storyURIs[st.Key] = st.URI
	}
	if got := mappingURIs["demo"]; got != "livt://mapping/demo" {
		t.Errorf("demo example_mapping_uri = %q, want livt://mapping/demo", got)
	}
	if _, ok := storyURIs["other"]; !ok {
		t.Fatal("story other missing from list")
	}
	if got := mappingURIs["other"]; got != "" {
		t.Errorf("other example_mapping_uri = %q, want empty (no mapping)", got)
	}
	// Every listed story links to its own story resource.
	for _, key := range []string{"demo", "other"} {
		if got := storyURIs[key]; got != "livt://story/"+key {
			t.Errorf("%s uri = %q, want livt://story/%s", key, got, key)
		}
	}
}

func TestListStoryMapsLinksStoryMapResource(t *testing.T) {
	s := newTestServer(t)

	_, out, err := s.listStoryMaps(context.Background(), nil, listStoryMapsInput{})
	if err != nil {
		t.Fatalf("listStoryMaps: %v", err)
	}
	if len(out.StoryMaps) != 1 || out.StoryMaps[0].Name != "デモマップ" {
		t.Fatalf("story maps = %+v, want one デモマップ", out.StoryMaps)
	}
	// The URI carries the display name percent-encoded, so it round-trips
	// through RFC 6570 template matching.
	want := "livt://story-map/%E3%83%87%E3%83%A2%E3%83%9E%E3%83%83%E3%83%97"
	if got := out.StoryMaps[0].URI; got != want {
		t.Errorf("story map uri = %q, want %q", got, want)
	}
}

func TestStoryMapReturnsActivitiesStepsCardsReleases(t *testing.T) {
	cfg := newTestServer(t).cfg
	sm, err := cfg.storyMap("デモマップ")
	if err != nil {
		t.Fatalf("storyMap: %v", err)
	}

	got := cfg.toStoryMapJSON(sm)
	if got.Name != "デモマップ" {
		t.Errorf("name = %q, want デモマップ", got.Name)
	}
	if len(got.Releases) != 1 || got.Releases[0].ID != "mvp" || got.Releases[0].Name != "MVP" {
		t.Fatalf("releases = %+v, want one mvp/MVP", got.Releases)
	}
	if len(got.Activities) != 1 || len(got.Activities[0].Steps) != 1 {
		t.Fatalf("activities = %+v, want one with one step", got.Activities)
	}
	cards := got.Activities[0].Steps[0].Stories
	if len(cards) != 2 {
		t.Fatalf("cards = %+v, want two", cards)
	}
	// A committed card links to its story resource; a candidate card stays bare.
	if cards[0].Key != "demo" || cards[0].URI != "livt://story/demo" || cards[0].Release != "mvp" {
		t.Errorf("committed card = %+v, want key demo with uri livt://story/demo and release mvp", cards[0])
	}
	if cards[1].Name != "候補ストーリー" || cards[1].Key != "" || cards[1].URI != "" {
		t.Errorf("candidate card = %+v, want bare 候補ストーリー", cards[1])
	}
}

func TestStoryMapUnknownNameErrors(t *testing.T) {
	if _, err := newTestServer(t).cfg.storyMap("なし"); err == nil {
		t.Fatal("expected error for unknown story map name")
	}
}

func TestStoryReturnsNameBodyMeta(t *testing.T) {
	cfg := newTestServer(t).cfg
	story, err := cfg.story("demo")
	if err != nil {
		t.Fatalf("story: %v", err)
	}

	got := cfg.toStoryJSON(story)
	if got.Key != "demo" || got.Name != "デモストーリー" || got.Body != "本文" {
		t.Errorf("story = %+v, want demo/デモストーリー/本文", got)
	}
	if len(got.Meta) != 1 || got.Meta[0].Key != "issue" || got.Meta[0].Value != "https://example.com/issues/1" {
		t.Errorf("meta = %+v, want one issue field", got.Meta)
	}
	if got.ExampleMappingURI != "livt://mapping/demo" {
		t.Errorf("example_mapping_uri = %q, want livt://mapping/demo", got.ExampleMappingURI)
	}
}

func TestStoryWithoutMappingHasNoMappingURI(t *testing.T) {
	cfg := newTestServer(t).cfg
	story, err := cfg.story("other")
	if err != nil {
		t.Fatalf("story: %v", err)
	}
	if got := cfg.toStoryJSON(story).ExampleMappingURI; got != "" {
		t.Errorf("example_mapping_uri = %q, want empty (no mapping)", got)
	}
}

func TestStoryUnknownKeyErrors(t *testing.T) {
	if _, err := newTestServer(t).cfg.story("nope"); err == nil {
		t.Fatal("expected error for unknown story key")
	}
}

func TestStoryRejectsTraversalKey(t *testing.T) {
	if _, err := newTestServer(t).cfg.story("../../etc/passwd"); err == nil {
		t.Fatal("expected error for traversal key")
	}
}

func TestTermReturnsNameAndBody(t *testing.T) {
	term, err := newTestServer(t).cfg.term("story")
	if err != nil {
		t.Fatalf("term: %v", err)
	}
	if term.Key != "story" || term.Name != "ストーリー" || term.Body != "ストーリーの定義" {
		t.Errorf("term = %+v, want story/ストーリー/ストーリーの定義", term)
	}
}

func TestTermUnknownKeyErrors(t *testing.T) {
	if _, err := newTestServer(t).cfg.term("nope"); err == nil {
		t.Fatal("expected error for unknown term key")
	}
}

func TestExampleMappingResolvesTermRefs(t *testing.T) {
	cfg := newTestServer(t).cfg
	em, err := cfg.exampleMapping("demo")
	if err != nil {
		t.Fatalf("exampleMapping: %v", err)
	}

	got := cfg.toExampleMappingJSON(em)
	// The raw key list is kept as-is for compatibility...
	if len(got.Ubiquitous) != 2 || got.Ubiquitous[0] != "story" || got.Ubiquitous[1] != "missing-term" {
		t.Fatalf("ubiquitous = %+v, want [story missing-term]", got.Ubiquitous)
	}
	// ...while ubiquitous_terms resolves each key: committed terms get a name
	// and resource URI, uncommitted ones stay bare keys.
	if len(got.UbiquitousTerms) != 2 {
		t.Fatalf("ubiquitous_terms = %+v, want two", got.UbiquitousTerms)
	}
	want := termRefJSON{Key: "story", Name: "ストーリー", URI: "livt://ubiquitous/story"}
	if got.UbiquitousTerms[0] != want {
		t.Errorf("resolved term ref = %+v, want %+v", got.UbiquitousTerms[0], want)
	}
	if got.UbiquitousTerms[1] != (termRefJSON{Key: "missing-term"}) {
		t.Errorf("unresolved term ref = %+v, want bare missing-term", got.UbiquitousTerms[1])
	}
}

func TestServerRegistersToolsWithoutPanic(t *testing.T) {
	if newTestServer(t).mcpServer() == nil {
		t.Fatal("mcpServer returned nil")
	}
}
