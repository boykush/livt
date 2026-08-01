package uri

import "testing"

// livt://mapping/trace-test-to-rule/rule/R-01/example/EX-01: a rule is
// livt://mapping/{story-key}/rule/{rule-id}.
func TestParseRule(t *testing.T) {
	cases := []struct {
		uri  string
		key  string
		rule string
		ok   bool
	}{
		{"livt://mapping/demo/rule/R-01", "demo", "R-01", true},
		{"livt://mapping/demo", "", "", false},           // mapping URI, no rule
		{"livt://mapping/demo/rule/", "", "", false},     // empty rule id
		{"livt://mapping//rule/R-01", "", "", false},     // empty key
		{"livt://mapping/demo/rule/../x", "", "", false}, // traversal in rule id
		{"livt://mapping/demo/question/Q-01", "", "", false},
	}
	for _, c := range cases {
		key, rule, ok := ParseRule(c.uri)
		if ok != c.ok || key != c.key || rule != c.rule {
			t.Errorf("ParseRule(%q) = (%q, %q, %v), want (%q, %q, %v)", c.uri, key, rule, ok, c.key, c.rule, c.ok)
		}
	}
}

// livt://mapping/trace-test-to-rule/rule/R-01/example/EX-02: an example is
// livt://mapping/{story-key}/rule/{rule-id}/example/{example-id} — the rule is
// part of the address because example ids are numbered within their rule.
func TestParseExample(t *testing.T) {
	cases := []struct {
		uri     string
		key     string
		rule    string
		example string
		ok      bool
	}{
		{"livt://mapping/demo/rule/R-01/example/EX-01", "demo", "R-01", "EX-01", true},
		{"livt://mapping/demo/rule/R-01", "", "", "", false},               // rule URI, no example
		{"livt://mapping/demo/rule/R-01/example/", "", "", "", false},      // empty example id
		{"livt://mapping/demo/rule//example/EX-01", "", "", "", false},     // empty rule id
		{"livt://mapping//rule/R-01/example/EX-01", "", "", "", false},     // empty key
		{"livt://mapping/demo/rule/R-01/example/../x", "", "", "", false},  // traversal in example id
		{"livt://mapping/demo/example/EX-01/rule/R-01", "", "", "", false}, // segments out of order
		{"livt://mapping/demo/rule/R-01/example/EX-01/example/EX-02", "", "", "", false},
	}
	for _, c := range cases {
		key, rule, example, ok := ParseExample(c.uri)
		if ok != c.ok || key != c.key || rule != c.rule || example != c.example {
			t.Errorf("ParseExample(%q) = (%q, %q, %q, %v), want (%q, %q, %q, %v)",
				c.uri, key, rule, example, ok, c.key, c.rule, c.example, c.ok)
		}
	}
}

// livt://mapping/trace-test-to-rule/rule/R-01/example/EX-02: an example URI and
// the rule URI it hangs off never resolve to each other, so following one
// cannot silently serve the other.
func TestRuleAndExampleURIsDoNotCrossParse(t *testing.T) {
	ruleURI := Rule("demo", "R-01")
	exampleURI := Example("demo", "R-01", "EX-01")

	if _, _, _, ok := ParseExample(ruleURI); ok {
		t.Errorf("ParseExample(%q) parsed a rule URI as an example", ruleURI)
	}
	if _, _, ok := ParseRule(exampleURI); ok {
		t.Errorf("ParseRule(%q) parsed an example URI as a rule", exampleURI)
	}
	if _, ok := ParseMapping(ruleURI); ok {
		t.Errorf("ParseMapping(%q) parsed a rule URI as a mapping", ruleURI)
	}
	if _, ok := ParseMapping(exampleURI); ok {
		t.Errorf("ParseMapping(%q) parsed an example URI as a mapping", exampleURI)
	}
}

// livt://mapping/trace-test-to-rule/rule/R-01/example/EX-03: a question is
// livt://mapping/{story-key}/question/{question-id} — it hangs off the mapping,
// not off a rule.
func TestParseQuestion(t *testing.T) {
	cases := []struct {
		uri      string
		key      string
		question string
		ok       bool
	}{
		{"livt://mapping/demo/question/Q-01", "demo", "Q-01", true},
		{"livt://mapping/demo", "", "", false},               // mapping URI, no question
		{"livt://mapping/demo/question/", "", "", false},     // empty question id
		{"livt://mapping//question/Q-01", "", "", false},     // empty key
		{"livt://mapping/demo/question/../x", "", "", false}, // traversal in question id
		{"livt://mapping/demo/rule/R-01", "", "", false},     // a rule URI, not a question URI
		{"livt://mapping/demo/rule/question", "", "", false}, // "question" as a rule id
	}
	for _, c := range cases {
		key, question, ok := ParseQuestion(c.uri)
		if ok != c.ok || key != c.key || question != c.question {
			t.Errorf("ParseQuestion(%q) = (%q, %q, %v), want (%q, %q, %v)", c.uri, key, question, ok, c.key, c.question, c.ok)
		}
	}
}

// livt://mapping/trace-test-to-rule/rule/R-01/example/EX-04: a bare id is not an
// address — the same rule and example ids live in every mapping, so only the
// story key tells them apart.
func TestBareIDsNeedTheStoryKey(t *testing.T) {
	if a, b := Rule("demo", "R-02"), Rule("other", "R-02"); a == b {
		t.Errorf("R-02 in two stories built the same URI %q", a)
	}
	if a, b := Example("demo", "R-01", "EX-01"), Example("demo", "R-02", "EX-01"); a == b {
		t.Errorf("EX-01 under two rules built the same URI %q", a)
	}
}

func TestParseMapping(t *testing.T) {
	cases := []struct {
		uri string
		key string
		ok  bool
	}{
		{"livt://mapping/demo", "demo", true},
		{"livt://mapping/automate-from-master-in-impl-repos", "automate-from-master-in-impl-repos", true},
		{"livt://mapping/", "", false},               // empty key
		{"livt://story/demo", "", false},             // wrong path
		{"livt://mapping/demo/rule/R-01", "", false}, // a rule URI, not a mapping URI
		{"livt://mapping/../secret", "", false},      // traversal
		{"livt://mapping/a/b", "", false},            // nested path
	}
	for _, c := range cases {
		key, ok := ParseMapping(c.uri)
		if ok != c.ok || key != c.key {
			t.Errorf("ParseMapping(%q) = (%q, %v), want (%q, %v)", c.uri, key, ok, c.key, c.ok)
		}
	}
}

func TestParseStoryMap(t *testing.T) {
	cases := []struct {
		uri  string
		name string
		ok   bool
	}{
		{"livt://story-map/plain-name", "plain-name", true},
		// Display names travel percent-encoded and are decoded on read.
		{"livt://story-map/%E3%83%87%E3%83%A2%E3%83%9E%E3%83%83%E3%83%97", "デモマップ", true},
		{"livt://story-map/", "", false},    // empty name
		{"livt://story/demo", "", false},    // a story URI, not a story map URI
		{"livt://story-map/%zz", "", false}, // invalid percent-encoding
	}
	for _, c := range cases {
		name, ok := ParseStoryMap(c.uri)
		if ok != c.ok || name != c.name {
			t.Errorf("ParseStoryMap(%q) = (%q, %v), want (%q, %v)", c.uri, name, ok, c.name, c.ok)
		}
	}
}

func TestParseStory(t *testing.T) {
	cases := []struct {
		uri string
		key string
		ok  bool
	}{
		{"livt://story/demo", "demo", true},
		{"livt://story/", "", false},          // empty key
		{"livt://story-map/demo", "", false},  // a story map URI, not a story URI
		{"livt://mapping/demo", "", false},    // wrong path
		{"livt://story/../secret", "", false}, // traversal
		{"livt://story/a/b", "", false},       // nested path
	}
	for _, c := range cases {
		key, ok := ParseStory(c.uri)
		if ok != c.ok || key != c.key {
			t.Errorf("ParseStory(%q) = (%q, %v), want (%q, %v)", c.uri, key, ok, c.key, c.ok)
		}
	}
}

func TestParseTerm(t *testing.T) {
	cases := []struct {
		uri string
		key string
		ok  bool
	}{
		{"livt://ubiquitous/story", "story", true},
		{"livt://ubiquitous/", "", false},          // empty key
		{"livt://story/story", "", false},          // wrong path
		{"livt://ubiquitous/../secret", "", false}, // traversal
	}
	for _, c := range cases {
		key, ok := ParseTerm(c.uri)
		if ok != c.ok || key != c.key {
			t.Errorf("ParseTerm(%q) = (%q, %v), want (%q, %v)", c.uri, key, ok, c.key, c.ok)
		}
	}
}

func TestValidSegment(t *testing.T) {
	for _, s := range []string{"demo", "R-01", "EX-01", "Q-01", "a.b"} {
		if !ValidSegment(s) {
			t.Errorf("ValidSegment(%q) = false, want true", s)
		}
	}
	for _, s := range []string{"", ".", "..", "../x", "a/b", `a\b`, "a..b"} {
		if ValidSegment(s) {
			t.Errorf("ValidSegment(%q) = true, want false", s)
		}
	}
}

func TestURIsRoundTrip(t *testing.T) {
	key, ok := ParseMapping(Mapping("demo"))
	if !ok || key != "demo" {
		t.Fatalf("mapping round trip = (%q, %v), want (demo, true)", key, ok)
	}
	k, id, ok := ParseRule(Rule("demo", "R-01"))
	if !ok || k != "demo" || id != "R-01" {
		t.Fatalf("rule round trip = (%q, %q, %v), want (demo, R-01, true)", k, id, ok)
	}
	ek, eid, exID, ok := ParseExample(Example("demo", "R-01", "EX-01"))
	if !ok || ek != "demo" || eid != "R-01" || exID != "EX-01" {
		t.Fatalf("example round trip = (%q, %q, %q, %v), want (demo, R-01, EX-01, true)", ek, eid, exID, ok)
	}
	qk, qid, ok := ParseQuestion(Question("demo", "Q-01"))
	if !ok || qk != "demo" || qid != "Q-01" {
		t.Fatalf("question round trip = (%q, %q, %v), want (demo, Q-01, true)", qk, qid, ok)
	}
	name, ok := ParseStoryMap(StoryMap("デモマップ"))
	if !ok || name != "デモマップ" {
		t.Fatalf("story map round trip = (%q, %v), want (デモマップ, true)", name, ok)
	}
	sk, ok := ParseStory(Story("demo"))
	if !ok || sk != "demo" {
		t.Fatalf("story round trip = (%q, %v), want (demo, true)", sk, ok)
	}
	tk, ok := ParseTerm(Term("story"))
	if !ok || tk != "story" {
		t.Fatalf("term round trip = (%q, %v), want (story, true)", tk, ok)
	}
}
