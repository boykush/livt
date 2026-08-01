package uri

import "testing"

// livt://mapping/trace-test-to-rule/rule/R-04/example/EX-01: a rule, example,
// question, story, story map, or term all resolve — so a URI arriving as free
// text is recognised whichever point of the master it names.
func TestParseRecognisesEveryKind(t *testing.T) {
	cases := []struct {
		uri  string
		want Parsed
	}{
		{"livt://mapping/demo", Parsed{Kind: KindMapping, StoryKey: "demo"}},
		{"livt://mapping/demo/rule/R-01", Parsed{Kind: KindRule, StoryKey: "demo", RuleID: "R-01"}},
		{
			"livt://mapping/demo/rule/R-01/example/EX-02",
			Parsed{Kind: KindExample, StoryKey: "demo", RuleID: "R-01", ExampleID: "EX-02"},
		},
		{"livt://mapping/demo/question/Q-01", Parsed{Kind: KindQuestion, StoryKey: "demo", QuestionID: "Q-01"}},
		{"livt://story/demo", Parsed{Kind: KindStory, StoryKey: "demo"}},
		{"livt://story-map/デモマップ", Parsed{Kind: KindStoryMap, MapName: "デモマップ"}},
		{"livt://story-map/%E3%83%87%E3%83%A2%E3%83%9E%E3%83%83%E3%83%97", Parsed{Kind: KindStoryMap, MapName: "デモマップ"}},
		{"livt://ubiquitous/livt-uri", Parsed{Kind: KindTerm, TermKey: "livt-uri"}},
	}
	for _, c := range cases {
		got, ok := Parse(c.uri)
		if !ok || got != c.want {
			t.Errorf("Parse(%q) = (%+v, %v), want (%+v, true)", c.uri, got, ok, c.want)
		}
	}
}

func TestParseRejectsMalformed(t *testing.T) {
	for _, s := range []string{
		"",
		"livt://",
		"livt://mapping/",
		"livt://nope/demo",
		"https://boykush.github.io/livt/mapping/demo.html", // a site URL, not a livt URI
		"mapping/demo",              // no scheme
		"livt://mapping/../secret",  // traversal
		"livt://mapping/demo/rule/", // empty rule id
		"livt://story/a/b",          // nested path
	} {
		if got, ok := Parse(s); ok {
			t.Errorf("Parse(%q) = (%+v, true), want false", s, got)
		}
	}
}

// A longer shape never falls back to the shorter one it is built on: an example
// URI must not resolve as the rule that carries it.
func TestParseKeepsNestedShapesApart(t *testing.T) {
	got, ok := Parse("livt://mapping/demo/rule/R-01/example/EX-02")
	if !ok || got.Kind != KindExample {
		t.Fatalf("Parse of an example URI = (%+v, %v), want kind %q", got, ok, KindExample)
	}
	got, ok = Parse("livt://mapping/demo/question/Q-01")
	if !ok || got.Kind != KindQuestion {
		t.Fatalf("Parse of a question URI = (%+v, %v), want kind %q", got, ok, KindQuestion)
	}
}

func TestParsedStringRoundTrips(t *testing.T) {
	for _, s := range []string{
		"livt://mapping/demo",
		"livt://mapping/demo/rule/R-01",
		"livt://mapping/demo/rule/R-01/example/EX-02",
		"livt://mapping/demo/question/Q-01",
		"livt://story/demo",
		"livt://ubiquitous/livt-uri",
		"livt://story-map/%E3%83%87%E3%83%A2%E3%83%9E%E3%83%83%E3%83%97",
	} {
		p, ok := Parse(s)
		if !ok {
			t.Fatalf("Parse(%q) failed", s)
		}
		if got := p.String(); got != s {
			t.Errorf("Parse(%q).String() = %q, want the same URI", s, got)
		}
	}

	// A plainly written map name comes back percent-encoded.
	p, _ := Parse("livt://story-map/デモマップ")
	if got := p.String(); got != "livt://story-map/%E3%83%87%E3%83%A2%E3%83%9E%E3%83%83%E3%83%97" {
		t.Errorf("String() = %q, want the percent-encoded form", got)
	}
}
