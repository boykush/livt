package mcp

import (
	"context"
	"encoding/json"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/boykush/livt/internal/uri"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// readResource drives one resource handler the way the SDK does, returning the
// decoded body. URI parsing and building live in internal/uri and are tested
// there; what these tests cover is that each handler serves its own shape.
func readResource[T any](t *testing.T, h func(context.Context, *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error), resURI string) T {
	t.Helper()
	res, err := h(context.Background(), &mcpsdk.ReadResourceRequest{Params: &mcpsdk.ReadResourceParams{URI: resURI}})
	if err != nil {
		t.Fatalf("read %q: %v", resURI, err)
	}
	var out T
	if err := json.Unmarshal([]byte(resourceText(t, res)), &out); err != nil {
		t.Fatalf("decode %q: %v", resURI, err)
	}
	return out
}

// livt://mapping/trace-test-to-rule/rule/R-01/example/EX-02: an example resolves
// through the rule that numbers it.
func TestReadExampleReturnsSingleExample(t *testing.T) {
	s := newTestServer(t)
	resURI := uri.Example("demo", "R-01", "EX-01")

	got := readResource[exampleResult](t, s.readExample, resURI).Example
	if got.ID != "EX-01" || got.Name != "実例1" {
		t.Errorf("example = %+v, want EX-01 実例1", got)
	}
	if got.URI != resURI {
		t.Errorf("example uri = %q, want %q", got.URI, resURI)
	}
}

// livt://mapping/trace-test-to-rule/rule/R-01/example/EX-03: a question resolves
// off the mapping, without a rule in the address.
func TestReadQuestionReturnsSingleQuestion(t *testing.T) {
	s := newTestServer(t)
	resURI := uri.Question("demo", "Q-01")

	got := readResource[questionResult](t, s.readQuestion, resURI).Question
	if got.ID != "Q-01" || got.Text != "質問1" {
		t.Errorf("question = %+v, want Q-01 質問1", got)
	}
	if got.URI != resURI {
		t.Errorf("question uri = %q, want %q", got.URI, resURI)
	}
}

// livt://mapping/trace-test-to-rule/rule/R-01/example/EX-02: a rule URI and an
// example URI address different things, so neither handler answers for the
// other's shape.
func TestRuleAndExampleHandlersRejectEachOthersURIs(t *testing.T) {
	s := newTestServer(t)

	if _, err := s.readExample(context.Background(), readReq(uri.Rule("demo", "R-01"))); err == nil {
		t.Error("readExample answered a rule URI")
	}
	if _, err := s.readRule(context.Background(), readReq(uri.Example("demo", "R-01", "EX-01"))); err == nil {
		t.Error("readRule answered an example URI")
	}
	if _, err := s.readMapping(context.Background(), readReq(uri.Example("demo", "R-01", "EX-01"))); err == nil {
		t.Error("readMapping answered an example URI")
	}
	if _, err := s.readMapping(context.Background(), readReq(uri.Question("demo", "Q-01"))); err == nil {
		t.Error("readMapping answered a question URI")
	}
}

func TestReadExampleAndQuestionNotFound(t *testing.T) {
	s := newTestServer(t)
	for _, resURI := range []string{
		uri.Example("demo", "R-01", "EX-99"), // unknown example
		uri.Example("demo", "R-99", "EX-01"), // unknown rule
		uri.Example("nope", "R-01", "EX-01"), // unknown story
		"livt://mapping/demo/rule/R-01/example/../secret",
		"livt://mapping/demo/rule/R-01/example/",
	} {
		if _, err := s.readExample(context.Background(), readReq(resURI)); err == nil {
			t.Errorf("expected error reading %q", resURI)
		}
	}
	for _, resURI := range []string{
		uri.Question("demo", "Q-99"), // unknown question
		uri.Question("nope", "Q-01"), // unknown story
		"livt://mapping/demo/question/../secret",
		"livt://mapping/demo/question/",
	} {
		if _, err := s.readQuestion(context.Background(), readReq(resURI)); err == nil {
			t.Errorf("expected error reading %q", resURI)
		}
	}
}

// livt://mapping/trace-test-to-rule/rule/R-01/example/EX-01 and EX-03: the
// examples and questions listed inside a mapping carry the URIs that address
// them on their own.
func TestMappingLinksExamplesAndQuestions(t *testing.T) {
	s := newTestServer(t)

	em := readResource[exampleMappingResult](t, s.readMapping, uri.Mapping("demo")).Mapping
	if got := em.Rules[0].Examples[0].URI; got != "livt://mapping/demo/rule/R-01/example/EX-01" {
		t.Errorf("example uri = %q, want livt://mapping/demo/rule/R-01/example/EX-01", got)
	}
	if got := em.Questions[0].URI; got != "livt://mapping/demo/question/Q-01" {
		t.Errorf("question uri = %q, want livt://mapping/demo/question/Q-01", got)
	}
}

// newRetiredTestServer lays out a livt repository whose demo mapping holds retired items
// beside live ones: rule R-02, example EX-02 of the live R-01, and question Q-02
// are retired, so their ids stay taken and their text stays on file.
func newRetiredTestServer(t *testing.T) *Server {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "discoveries", "example-mappings", "demo.yaml"),
		"rules:\n"+
			"  - id: R-01\n"+
			"    name: 現役のルール\n"+
			"    examples:\n"+
			"      - id: EX-01\n"+
			"        name: 現役の実例\n"+
			"      - id: EX-02\n"+
			"        name: 退役した実例\n"+
			"        retired: true\n"+
			"  - id: R-02\n"+
			"    name: 退役したルール\n"+
			"    retired: true\n"+
			"questions:\n"+
			"  - id: Q-01\n"+
			"    text: 現役の疑問\n"+
			"  - id: Q-02\n"+
			"    text: 退役した疑問\n"+
			"    retired: true\n")
	return NewServer(Config{Root: root}, "test")
}

// livt://mapping/trace-test-to-rule/rule/R-05/example/EX-01 and EX-03: a retired
// rule, example, and question each resolve by URI — flagged retired, text intact
// — rather than 404ing, which is what a reference embedded elsewhere follows.
func TestReadRetiredItemsResolveAsRetired(t *testing.T) {
	s := newRetiredTestServer(t)

	rule := readResource[ruleResult](t, s.readRule, uri.Rule("demo", "R-02")).Rule
	if !rule.Retired || rule.Name != "退役したルール" {
		t.Errorf("rule = %+v, want R-02 retired with its text kept", rule)
	}
	example := readResource[exampleResult](t, s.readExample, uri.Example("demo", "R-01", "EX-02")).Example
	if !example.Retired || example.Name != "退役した実例" {
		t.Errorf("example = %+v, want EX-02 retired with its text kept", example)
	}
	question := readResource[questionResult](t, s.readQuestion, uri.Question("demo", "Q-02")).Question
	if !question.Retired || question.Text != "退役した疑問" {
		t.Errorf("question = %+v, want Q-02 retired with its text kept", question)
	}

	live := readResource[ruleResult](t, s.readRule, uri.Rule("demo", "R-01")).Rule
	if live.Retired {
		t.Errorf("rule = %+v, want R-01 live", live)
	}
}

// livt://mapping/trace-test-to-rule/rule/R-05: the mapping resource is the
// structural record its ids are numbered from, so it keeps retired entries —
// flagged, so a consumer can drop them, but never making a taken id look free.
func TestMappingKeepsRetiredEntriesFlagged(t *testing.T) {
	s := newRetiredTestServer(t)

	em := readResource[exampleMappingResult](t, s.readMapping, uri.Mapping("demo")).Mapping
	if len(em.Rules) != 2 || !em.Rules[1].Retired {
		t.Fatalf("rules = %+v, want both, the second flagged retired", em.Rules)
	}
	if len(em.Rules[0].Examples) != 2 || !em.Rules[0].Examples[1].Retired {
		t.Errorf("examples = %+v, want both, the second flagged retired", em.Rules[0].Examples)
	}
	if len(em.Questions) != 2 || !em.Questions[1].Retired {
		t.Errorf("questions = %+v, want both, the second flagged retired", em.Questions)
	}
}

// A live item carries no retired field at all, so the flag reads as the
// exception it marks rather than as noise on every item.
func TestLiveItemsOmitRetiredFromJSON(t *testing.T) {
	s := newRetiredTestServer(t)

	// Leaf resources only: a rule's payload nests its examples, one of which is
	// retired here.
	leaves := []struct {
		read   func(context.Context, *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error)
		resURI string
	}{
		{s.readExample, uri.Example("demo", "R-01", "EX-01")},
		{s.readQuestion, uri.Question("demo", "Q-01")},
	}
	for _, leaf := range leaves {
		res, err := leaf.read(context.Background(), readReq(leaf.resURI))
		if err != nil {
			t.Fatalf("read %q: %v", leaf.resURI, err)
		}
		if body := resourceText(t, res); strings.Contains(body, "retired") {
			t.Errorf("live payload mentions retired: %s", body)
		}
	}
}

func readReq(resURI string) *mcpsdk.ReadResourceRequest {
	return &mcpsdk.ReadResourceRequest{Params: &mcpsdk.ReadResourceParams{URI: resURI}}
}

// newSupersededTestServer lays out a livt repository whose retired items name where the
// spec went: R-01 split into R-02 here and a rule in another mapping, EX-01 was
// replaced by EX-02 under the same rule, and Q-01 was settled by R-02.
func newSupersededTestServer(t *testing.T) *Server {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "discoveries", "example-mappings", "demo.yaml"),
		"rules:\n"+
			"  - id: R-01\n"+
			"    name: 分割されて退役したルール\n"+
			"    retired: true\n"+
			"    superseded_by:\n"+
			"      - livt://mapping/demo/rule/R-02\n"+
			"      - livt://mapping/other/rule/R-07\n"+
			"  - id: R-02\n"+
			"    name: 現役のルール\n"+
			"    examples:\n"+
			"      - id: EX-01\n"+
			"        name: 差し替えられた実例\n"+
			"        retired: true\n"+
			"        superseded_by:\n"+
			"          - livt://mapping/demo/rule/R-02/example/EX-02\n"+
			"      - id: EX-02\n"+
			"        name: 現役の実例\n"+
			"questions:\n"+
			"  - id: Q-01\n"+
			"    text: ルール化されて閉じた疑問\n"+
			"    retired: true\n"+
			"    superseded_by:\n"+
			"      - livt://mapping/demo/rule/R-02\n")
	return NewServer(Config{Root: root}, "test")
}

// livt://mapping/trace-test-to-rule/rule/R-09/example/EX-01, EX-02 and EX-03: a
// reference that lands on a retired item reads on to its successors, which are
// livt URIs — so one of them can live in another mapping, and a split can name
// both.
func TestReadRetiredItemsCarrySuccessors(t *testing.T) {
	s := newSupersededTestServer(t)

	rule := readResource[ruleResult](t, s.readRule, uri.Rule("demo", "R-01")).Rule
	want := []string{"livt://mapping/demo/rule/R-02", "livt://mapping/other/rule/R-07"}
	if !slices.Equal(rule.SupersededBy, want) {
		t.Errorf("rule superseded_by = %v, want %v", rule.SupersededBy, want)
	}
	example := readResource[exampleResult](t, s.readExample, uri.Example("demo", "R-02", "EX-01")).Example
	if !slices.Equal(example.SupersededBy, []string{"livt://mapping/demo/rule/R-02/example/EX-02"}) {
		t.Errorf("example superseded_by = %v, want the example that replaced it", example.SupersededBy)
	}
	question := readResource[questionResult](t, s.readQuestion, uri.Question("demo", "Q-01")).Question
	if !slices.Equal(question.SupersededBy, []string{"livt://mapping/demo/rule/R-02"}) {
		t.Errorf("question superseded_by = %v, want the rule that settled it", question.SupersededBy)
	}
}

// livt://mapping/trace-test-to-rule/rule/R-09/example/EX-06: the successor
// travels as a URI and nothing more. Inlining what it says would spend the
// caller's context on a hop most of them never take.
func TestSuccessorsTravelAsURIsNotText(t *testing.T) {
	s := newSupersededTestServer(t)

	res, err := s.readRule(context.Background(), readReq(uri.Rule("demo", "R-01")))
	if err != nil {
		t.Fatalf("read retired rule: %v", err)
	}
	if body := resourceText(t, res); strings.Contains(body, "現役のルール") {
		t.Errorf("payload inlines the successor's text: %s", body)
	}
}

// livt://mapping/trace-test-to-rule/rule/R-09/example/EX-04: a rule the spec
// simply stopped asking for is retired with nothing to point at, so the field
// is absent rather than empty.
func TestRetiredWithoutSuccessorOmitsSupersededBy(t *testing.T) {
	s := newRetiredTestServer(t)

	res, err := s.readRule(context.Background(), readReq(uri.Rule("demo", "R-02")))
	if err != nil {
		t.Fatalf("read retired rule: %v", err)
	}
	body := resourceText(t, res)
	if !strings.Contains(body, "retired") {
		t.Fatalf("payload should still be flagged retired: %s", body)
	}
	if strings.Contains(body, "superseded_by") {
		t.Errorf("payload mentions superseded_by with no successor: %s", body)
	}
}
