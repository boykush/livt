package mcp

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/boykush/livt/internal/uri"
)

// newTestServer lays out a livt repository under a temp root with one mapped story
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
			"    issues:\n"+
			"      - https://github.com/boykush/livt/issues/25\n"+
			"  - id: R-02\n"+
			"    name: ルール2\n"+
			"questions:\n"+
			"  - id: Q-01\n"+
			"    text: 質問1\n"+
			"ubiquitous:\n"+
			"  - story\n"+
			"  - missing-term\n")
	writeFile(t, filepath.Join(root, "automations", "acme", "impl.json"),
		"{\n"+
			"  \"repo\": \"acme/impl\",\n"+
			"  \"rev\": \"abc123\",\n"+
			"  \"citations\": [\n"+
			"    {\"livt_uri\": \"livt://mapping/demo/rule/R-01\", \"file\": \"x_test.go\", \"line\": 7}\n"+
			"  ]\n"+
			"}\n")
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
	// Keyed demo-map so it joins the story map by filename, the way an example
	// mapping joins its story.
	writeFile(t, filepath.Join(root, "opportunities", "demo-map.md"),
		"---\nname: デモ機会\nrepos:\n  - boykush/livt\n---\n\n誰かの課題を解くと、事業の便益が立つ\n")
	writeFile(t, filepath.Join(root, "discoveries", "opportunity-canvases", "demo-map.yaml"),
		"canvas:\n  problems:\n    - 課題\n  budget:\n    - 余暇\nubiquitous:\n  - story\n")

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
	if len(em.Rules) != 2 || em.Rules[0].ID != "R-01" {
		t.Fatalf("rules = %+v, want R-01 first of two", em.Rules)
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

func TestRuleJSONCarriesIssuesAndAutomated(t *testing.T) {
	cfg := newTestServer(t).cfg

	recorded, err := cfg.rule("demo", "R-01")
	if err != nil {
		t.Fatalf("rule: %v", err)
	}
	j := toRuleJSON("demo", recorded)
	if len(j.Issues) != 1 || j.Issues[0] != "https://github.com/boykush/livt/issues/25" {
		t.Errorf("issues = %v, want the recorded URL", j.Issues)
	}
	if !j.Automated {
		t.Error("R-01 should be automated")
	}

	bare, err := cfg.rule("demo", "R-02")
	if err != nil {
		t.Fatalf("rule: %v", err)
	}
	bj := toRuleJSON("demo", bare)
	if len(bj.Issues) != 0 || bj.Automated {
		t.Errorf("R-02 should be unlinked and not automated, got issues=%v automated=%v", bj.Issues, bj.Automated)
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

// demoOpportunityRef is the ref a story on the test repository's デモマップ
// carries: the map is keyed demo-map, and an opportunity file describes that
// key, so the ref is the opportunity's rather than the map's.
var demoOpportunityRef = opportunityRefJSON{Key: "demo-map", Name: "デモ機会", URI: "livt://opportunity/demo-map"}

// livt:automates livt://mapping/automate-from-master-in-impl-repos/rule/R-13/example/EX-05
// Each entry shows its opportunities by key, name, and resource URI, and a
// story on no map shows none.
func TestListStoriesCarriesOpportunities(t *testing.T) {
	s := newTestServer(t)

	_, out, err := s.listStories(context.Background(), nil, listStoriesInput{})
	if err != nil {
		t.Fatalf("listStories: %v", err)
	}
	byKey := storiesByKey(out.Stories)
	if got := byKey["demo"].Opportunities; len(got) != 1 || got[0] != demoOpportunityRef {
		t.Errorf("demo opportunities = %+v, want [%+v]", got, demoOpportunityRef)
	}
	if got := byKey["other"].Opportunities; len(got) != 0 {
		t.Errorf("other opportunities = %+v, want none (on no map)", got)
	}
}

// livt:automates livt://mapping/automate-from-master-in-impl-repos/rule/R-13/example/EX-06
// A map with no opportunity file stands in as its own opportunity: the key is
// still the map's, and the name and URI are the map's too.
func TestMapWithNoOpportunityFileStandsInAsItsOwn(t *testing.T) {
	s := newTestServer(t)
	writeFile(t, filepath.Join(s.cfg.Root, "discoveries", "usm", "second-map.yaml"),
		"name: 第二マップ\n"+
			"activities:\n"+
			"  - name: アクティビティ\n"+
			"    steps:\n"+
			"      - name: ステップ\n"+
			"        stories:\n"+
			"          - name: 別ストーリー\n"+
			"            key: other\n")

	_, out, err := s.listStories(context.Background(), nil, listStoriesInput{})
	if err != nil {
		t.Fatalf("listStories: %v", err)
	}
	want := opportunityRefJSON{Key: "second-map", Name: "第二マップ", URI: uri.StoryMap("第二マップ")}
	if got := storiesByKey(out.Stories)["other"].Opportunities; len(got) != 1 || got[0] != want {
		t.Errorf("other opportunities = %+v, want [%+v]", got, want)
	}
}

// livt:automates livt://mapping/automate-from-master-in-impl-repos/rule/R-13/example/EX-07
// The opportunity parameter takes the key list_opportunities hands out and
// keeps only the stories on its map. A name is not a key, so it matches none.
func TestListStoriesFiltersByOpportunity(t *testing.T) {
	s := newTestServer(t)

	_, out, err := s.listStories(context.Background(), nil, listStoriesInput{Opportunity: "demo-map"})
	if err != nil {
		t.Fatalf("listStories: %v", err)
	}
	if len(out.Stories) != 1 || out.Stories[0].Key != "demo" {
		t.Fatalf("stories = %+v, want only demo (the one story on demo-map)", out.Stories)
	}

	for _, name := range []string{"デモ機会", "デモマップ"} {
		_, byName, err := s.listStories(context.Background(), nil, listStoriesInput{Opportunity: name})
		if err != nil {
			t.Fatalf("listStories(%q): %v", name, err)
		}
		if len(byName.Stories) != 0 {
			t.Errorf("filtering by the name %q kept %+v, want none", name, byName.Stories)
		}
	}
}

// livt:automates livt://mapping/automate-from-master-in-impl-repos/rule/R-13/example/EX-08
// An unknown opportunity key yields an empty list, not an error.
func TestListStoriesUnknownOpportunityYieldsEmptyList(t *testing.T) {
	s := newTestServer(t)

	_, out, err := s.listStories(context.Background(), nil, listStoriesInput{Opportunity: "no-such-opportunity"})
	if err != nil {
		t.Fatalf("listStories: %v", err)
	}
	if len(out.Stories) != 0 {
		t.Errorf("stories = %+v, want empty for an unknown opportunity", out.Stories)
	}
}

// A story on several maps carries one ref per map in map file order, a key
// recurring across steps within one map still gets a single ref, and the
// filter matches the story through either of its maps.
func TestListStoriesStoryOnSeveralMaps(t *testing.T) {
	s := newTestServer(t)
	writeFile(t, filepath.Join(s.cfg.Root, "discoveries", "usm", "second-map.yaml"),
		"name: 第二マップ\n"+
			"activities:\n"+
			"  - name: アクティビティ\n"+
			"    steps:\n"+
			"      - name: ステップ1\n"+
			"        stories:\n"+
			"          - name: デモストーリー\n"+
			"            key: demo\n"+
			"      - name: ステップ2\n"+
			"        stories:\n"+
			"          - name: デモストーリー（再掲）\n"+
			"            key: demo\n")

	_, out, err := s.listStories(context.Background(), nil, listStoriesInput{})
	if err != nil {
		t.Fatalf("listStories: %v", err)
	}
	got := storiesByKey(out.Stories)["demo"].Opportunities
	if len(got) != 2 || got[0].Key != "demo-map" || got[1].Key != "second-map" {
		t.Fatalf("demo opportunities = %+v, want [demo-map second-map] in map file order", got)
	}

	_, filtered, err := s.listStories(context.Background(), nil, listStoriesInput{Opportunity: "second-map"})
	if err != nil {
		t.Fatalf("listStories filtered: %v", err)
	}
	if len(filtered.Stories) != 1 || filtered.Stories[0].Key != "demo" {
		t.Fatalf("filtered stories = %+v, want only demo (on second-map)", filtered.Stories)
	}
}

func storiesByKey(stories []storySummaryJSON) map[string]storySummaryJSON {
	byKey := make(map[string]storySummaryJSON, len(stories))
	for _, st := range stories {
		byKey[st.Key] = st
	}
	return byKey
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

// Automates livt://mapping/automate-from-master-in-impl-repos/rule/R-16.
func TestListTermsEnumeratesGlossaryWithResourceURIs(t *testing.T) {
	s := newTestServer(t)
	// A term under a context directory, and one no board references — the
	// glossary is the livt repository's vocabulary, not what the boards cite.
	writeFile(t, filepath.Join(s.cfg.Root, "ubiquitous", "unreferenced.md"), "---\nname: 未参照の用語\n---\n\n定義\n")
	writeFile(t, filepath.Join(s.cfg.Root, "ubiquitous", "billing", "story.md"), "---\nname: ストーリー（請求）\n---\n\n請求文脈の定義\n")

	_, out, err := s.listTerms(context.Background(), nil, listTermsInput{})
	if err != nil {
		t.Fatalf("listTerms: %v", err)
	}

	byURI := make(map[string]termSummaryJSON, len(out.Terms))
	for _, term := range out.Terms {
		byURI[term.URI] = term
	}
	if len(out.Terms) != 3 || len(byURI) != 3 {
		t.Fatalf("terms = %+v, want three distinct (story, unreferenced, billing/story)", out.Terms)
	}
	if got := byURI["livt://ubiquitous/story"]; got.Key != "story" || got.Ctx != "" || got.Name != "ストーリー" {
		t.Errorf("context-free term = %+v, want key story, no ctx, name ストーリー", got)
	}
	// The same key at the root and under a context are two terms, so the ctx
	// has to come through — a listing carrying key alone shows two rows a
	// consumer cannot tell apart.
	if got := byURI["livt://ubiquitous/billing/story"]; got.Key != "story" || got.Ctx != "billing" || got.Name != "ストーリー（請求）" {
		t.Errorf("scoped term = %+v, want key story, ctx billing, name ストーリー（請求）", got)
	}
	if _, ok := byURI["livt://ubiquitous/unreferenced"]; !ok {
		t.Errorf("terms = %+v, want the unreferenced term listed too", out.Terms)
	}
	// missing-term is referenced by the demo mapping but has no file, so it is
	// not a term — the listing enumerates the glossary, not the references.
	if _, ok := byURI["livt://ubiquitous/missing-term"]; ok {
		t.Errorf("terms = %+v, want no entry for the uncommitted missing-term", out.Terms)
	}
}

func TestListTermsOnMissingUbiquitousDirIsEmpty(t *testing.T) {
	s := NewServer(Config{Root: t.TempDir()}, "test")

	_, out, err := s.listTerms(context.Background(), nil, listTermsInput{})
	if err != nil {
		t.Fatalf("listTerms: %v", err)
	}
	if len(out.Terms) != 0 {
		t.Errorf("terms = %+v, want empty", out.Terms)
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

	got, err := cfg.toStoryJSON(story)
	if err != nil {
		t.Fatalf("toStoryJSON: %v", err)
	}
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
	got, err := cfg.toStoryJSON(story)
	if err != nil {
		t.Fatalf("toStoryJSON: %v", err)
	}
	if got.ExampleMappingURI != "" {
		t.Errorf("example_mapping_uri = %q, want empty (no mapping)", got.ExampleMappingURI)
	}
}

// livt:automates livt://mapping/automate-from-master-in-impl-repos/rule/R-13/example/EX-04
// The story resource shows the same opportunities as the list — the maps the
// story sits on, empty for a story on no map.
func TestStoryJSONCarriesOpportunities(t *testing.T) {
	cfg := newTestServer(t).cfg

	story, err := cfg.story("demo")
	if err != nil {
		t.Fatalf("story: %v", err)
	}
	got, err := cfg.toStoryJSON(story)
	if err != nil {
		t.Fatalf("toStoryJSON: %v", err)
	}
	if len(got.Opportunities) != 1 || got.Opportunities[0] != demoOpportunityRef {
		t.Errorf("demo opportunities = %+v, want [%+v]", got.Opportunities, demoOpportunityRef)
	}

	unmapped, err := cfg.story("other")
	if err != nil {
		t.Fatalf("story: %v", err)
	}
	oj, err := cfg.toStoryJSON(unmapped)
	if err != nil {
		t.Fatalf("toStoryJSON: %v", err)
	}
	if len(oj.Opportunities) != 0 {
		t.Errorf("other opportunities = %+v, want none (on no map)", oj.Opportunities)
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

// livt:automates livt://mapping/name-example-mapping-itself/rule/R-06
// livt:automates livt://mapping/name-example-mapping-itself/rule/R-06/example/EX-01
// livt:automates livt://mapping/name-example-mapping-itself/rule/R-06/example/EX-02
// A mapping with no story is listed beside the rest with its own name, while a
// story's name stays where it lives, behind story_uri.
func TestListExampleMappingsReachesAMappingWithoutAStory(t *testing.T) {
	s := newTestServer(t)
	writeFile(t, filepath.Join(s.cfg.Root, "discoveries", "example-mappings", "fix-login-with-full-width-space.yaml"),
		"name: 全角スペースでログインできない\nrules: []\n")

	_, out, err := s.listExampleMappings(context.Background(), nil, listExampleMappingsInput{})
	if err != nil {
		t.Fatalf("listExampleMappings: %v", err)
	}
	byKey := make(map[string]exampleMappingSummaryJSON, len(out.ExampleMappings))
	for _, m := range out.ExampleMappings {
		byKey[m.StoryKey] = m
	}

	fix := byKey["fix-login-with-full-width-space"]
	if fix.URI != "livt://mapping/fix-login-with-full-width-space" || fix.Name != "全角スペースでログインできない" || fix.StoryURI != "" {
		t.Errorf("story-less mapping = %+v, want its uri and own name and no story_uri", fix)
	}
	demo := byKey["demo"]
	if demo.URI != "livt://mapping/demo" || demo.StoryURI != "livt://story/demo" || demo.Name != "" {
		t.Errorf("story-backed mapping = %+v, want its uri and story_uri and no name of its own", demo)
	}
}

// livt:automates livt://mapping/name-example-mapping-itself/rule/R-06/example/EX-03
// Narrowed to an opportunity, the list keeps the mappings on that map, and one
// with no story sits on none.
func TestListExampleMappingsFiltersByOpportunity(t *testing.T) {
	s := newTestServer(t)
	writeFile(t, filepath.Join(s.cfg.Root, "discoveries", "example-mappings", "fix-login-with-full-width-space.yaml"),
		"name: 全角スペースでログインできない\nrules: []\n")

	_, out, err := s.listExampleMappings(context.Background(), nil, listExampleMappingsInput{Opportunity: "demo-map"})
	if err != nil {
		t.Fatalf("listExampleMappings: %v", err)
	}
	if len(out.ExampleMappings) != 1 || out.ExampleMappings[0].StoryKey != "demo" {
		t.Fatalf("mappings = %+v, want only demo (the one mapping on demo-map)", out.ExampleMappings)
	}
	if got := out.ExampleMappings[0].Opportunities; len(got) != 1 || got[0] != demoOpportunityRef {
		t.Errorf("demo opportunities = %+v, want [%+v]", got, demoOpportunityRef)
	}
}
