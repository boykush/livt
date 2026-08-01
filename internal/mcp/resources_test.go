package mcp

import (
	"context"
	"encoding/json"
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

func readReq(resURI string) *mcpsdk.ReadResourceRequest {
	return &mcpsdk.ReadResourceRequest{Params: &mcpsdk.ReadResourceParams{URI: resURI}}
}
