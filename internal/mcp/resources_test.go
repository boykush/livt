package mcp

import "testing"

func TestParseMappingURI(t *testing.T) {
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
		key, ok := parseMappingURI(c.uri)
		if ok != c.ok || key != c.key {
			t.Errorf("parseMappingURI(%q) = (%q, %v), want (%q, %v)", c.uri, key, ok, c.key, c.ok)
		}
	}
}

func TestParseRuleURI(t *testing.T) {
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
	}
	for _, c := range cases {
		key, rule, ok := parseRuleURI(c.uri)
		if ok != c.ok || key != c.key || rule != c.rule {
			t.Errorf("parseRuleURI(%q) = (%q, %q, %v), want (%q, %q, %v)", c.uri, key, rule, ok, c.key, c.rule, c.ok)
		}
	}
}

func TestParseStoryMapURI(t *testing.T) {
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
		name, ok := parseStoryMapURI(c.uri)
		if ok != c.ok || name != c.name {
			t.Errorf("parseStoryMapURI(%q) = (%q, %v), want (%q, %v)", c.uri, name, ok, c.name, c.ok)
		}
	}
}

func TestParseStoryURI(t *testing.T) {
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
		key, ok := parseStoryURI(c.uri)
		if ok != c.ok || key != c.key {
			t.Errorf("parseStoryURI(%q) = (%q, %v), want (%q, %v)", c.uri, key, ok, c.key, c.ok)
		}
	}
}

func TestParseTermURI(t *testing.T) {
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
		key, ok := parseTermURI(c.uri)
		if ok != c.ok || key != c.key {
			t.Errorf("parseTermURI(%q) = (%q, %v), want (%q, %v)", c.uri, key, ok, c.key, c.ok)
		}
	}
}

func TestURIsRoundTrip(t *testing.T) {
	key, ok := parseMappingURI(mappingURI("demo"))
	if !ok || key != "demo" {
		t.Fatalf("mapping round trip = (%q, %v), want (demo, true)", key, ok)
	}
	k, id, ok := parseRuleURI(ruleURI("demo", "R-01"))
	if !ok || k != "demo" || id != "R-01" {
		t.Fatalf("rule round trip = (%q, %q, %v), want (demo, R-01, true)", k, id, ok)
	}
	name, ok := parseStoryMapURI(storyMapURI("デモマップ"))
	if !ok || name != "デモマップ" {
		t.Fatalf("story map round trip = (%q, %v), want (デモマップ, true)", name, ok)
	}
	sk, ok := parseStoryURI(storyURI("demo"))
	if !ok || sk != "demo" {
		t.Fatalf("story round trip = (%q, %v), want (demo, true)", sk, ok)
	}
	tk, ok := parseTermURI(termURI("story"))
	if !ok || tk != "story" {
		t.Fatalf("term round trip = (%q, %v), want (story, true)", tk, ok)
	}
}
