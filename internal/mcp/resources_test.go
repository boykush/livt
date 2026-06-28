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

func TestURIsRoundTrip(t *testing.T) {
	key, ok := parseMappingURI(mappingURI("demo"))
	if !ok || key != "demo" {
		t.Fatalf("mapping round trip = (%q, %v), want (demo, true)", key, ok)
	}
	k, id, ok := parseRuleURI(ruleURI("demo", "R-01"))
	if !ok || k != "demo" || id != "R-01" {
		t.Fatalf("rule round trip = (%q, %q, %v), want (demo, R-01, true)", k, id, ok)
	}
}
