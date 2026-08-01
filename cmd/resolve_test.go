package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// newTestMaster writes a master small enough to assert against but holding one
// of every addressable kind.
func newTestMaster(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	write := func(path, content string) {
		full := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write(filepath.Join("discoveries", "example-mappings", "demo.yaml"),
		"rules:\n"+
			"  - id: R-01\n"+
			"    name: ルール1\n"+
			"    examples:\n"+
			"      - id: EX-01\n"+
			"        name: 実例1\n"+
			"  - id: R-02\n"+
			"    name: 退役したルール\n"+
			"    retired: true\n"+
			"questions:\n"+
			"  - id: Q-01\n"+
			"    text: 質問1\n")
	write(filepath.Join("stories", "demo.md"), "---\nname: デモストーリー\n---\n\n本文\n")
	write(filepath.Join("discoveries", "usm", "demo-map.yaml"),
		"name: デモマップ\nactivities: []\n")
	write(filepath.Join("ubiquitous", "story.md"), "---\nname: ストーリー\n---\n\n定義\n")

	return root
}

// livt://mapping/trace-test-to-rule/rule/R-04/example/EX-01: a rule, example,
// question, story, story map, or term all resolve from the command line.
func TestResolveURIResolvesEveryKind(t *testing.T) {
	root := newTestMaster(t)
	cases := []struct {
		uri   string
		field string
	}{
		{"livt://mapping/demo", "mapping"},
		{"livt://mapping/demo/rule/R-01", "rule"},
		{"livt://mapping/demo/rule/R-01/example/EX-01", "example"},
		{"livt://mapping/demo/question/Q-01", "question"},
		{"livt://story-map/デモマップ", "story_map"},
		{"livt://story/demo", "story"},
		{"livt://ubiquitous/story", "term"},
	}
	for _, c := range cases {
		var out bytes.Buffer
		if err := resolveURI(&out, root, c.uri, formatJSON, ""); err != nil {
			t.Errorf("resolve %s: %v", c.uri, err)
			continue
		}
		var payload map[string]any
		if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
			t.Errorf("decode %s: %v", c.uri, err)
			continue
		}
		if _, ok := payload[c.field]; !ok {
			t.Errorf("%s produced %v, want a %q field", c.uri, keysOf(payload), c.field)
		}
		if _, ok := payload["spec_version"]; !ok {
			t.Errorf("%s produced no spec_version, so a consumer cannot tell which master it read", c.uri)
		}
	}
}

// livt://mapping/trace-test-to-rule/rule/R-04/example/EX-03: the URL form is
// the item's page on the deployed site. The paths and anchors come from the
// same derivation the site build anchors its stickies to, so a link resolved
// here lands where the board says it does.
func TestResolveURIWritesTheURLForm(t *testing.T) {
	root := newTestMaster(t)
	const base = "https://boykush.github.io/livt"
	cases := []struct{ uri, want string }{
		{"livt://mapping/demo", base + "/mapping/demo.html"},
		{"livt://mapping/demo/rule/R-01", base + "/mapping/demo.html#rule-R-01"},
		{"livt://mapping/demo/rule/R-01/example/EX-01", base + "/mapping/demo.html#rule-R-01-example-EX-01"},
		{"livt://mapping/demo/question/Q-01", base + "/mapping/demo.html#question-Q-01"},
		{"livt://story/demo", base + "/story/demo.html"},
		{"livt://ubiquitous/story", base + "/ubiquitous.html#story"},
		{"livt://story-map/デモマップ", base + "/story-map/デモマップ.html"},
	}
	for _, c := range cases {
		var out bytes.Buffer
		if err := resolveURI(&out, root, c.uri, formatURL, base); err != nil {
			t.Errorf("resolve %s: %v", c.uri, err)
			continue
		}
		if got := strings.TrimSpace(out.String()); got != c.want {
			t.Errorf("resolve %s --format url = %q, want %q", c.uri, got, c.want)
		}
	}

	// A site root pasted with its trailing slash must not double it.
	var out bytes.Buffer
	if err := resolveURI(&out, root, "livt://story/demo", formatURL, base+"/"); err != nil {
		t.Fatalf("resolve with a trailing slash: %v", err)
	}
	if got := strings.TrimSpace(out.String()); got != base+"/story/demo.html" {
		t.Errorf("trailing-slash base gave %q, want %q", got, base+"/story/demo.html")
	}
}

// A URL is derivable from the URI without reading anything, which is exactly
// why the master still has to be checked: a link to a page that was never built
// is worse than an error.
func TestResolveURIRefusesAURLItCannotBack(t *testing.T) {
	root := newTestMaster(t)

	err := resolveURI(&bytes.Buffer{}, root, "livt://mapping/demo/rule/R-99", formatURL, "https://example.com")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("url for a missing rule = %v, want a not-found error", err)
	}

	err = resolveURI(&bytes.Buffer{}, root, "livt://story/demo", formatURL, "")
	if err == nil || !strings.Contains(err.Error(), "--base-url") {
		t.Errorf("url without a base = %v, want it to name --base-url", err)
	}
}

// The two ways a URI fails need different fixes -- correct the URI, or add the
// item to the master -- so they must not read alike.
func TestResolveURIDistinguishesMalformedFromMissing(t *testing.T) {
	root := newTestMaster(t)

	for _, raw := range []string{"", "R-01", "livt://nope/x", "livt://mapping/demo/rule/"} {
		err := resolveURI(&bytes.Buffer{}, root, raw, formatJSON, "")
		if err == nil {
			t.Errorf("resolve %q succeeded, want a malformed-URI error", raw)
			continue
		}
		if !strings.Contains(err.Error(), "malformed livt URI") {
			t.Errorf("resolve %q said %q, want it to name the URI as malformed", raw, err)
		}
	}

	for _, raw := range []string{"livt://mapping/demo/rule/R-99", "livt://story/nope"} {
		err := resolveURI(&bytes.Buffer{}, root, raw, formatJSON, "")
		if err == nil {
			t.Errorf("resolve %q succeeded, want a not-found error", raw)
			continue
		}
		if strings.Contains(err.Error(), "malformed") || !strings.Contains(err.Error(), "not found") {
			t.Errorf("resolve %q said %q, want it to report a well-formed URI with no item", raw, err)
		}
	}
}

// livt://mapping/trace-test-to-rule/rule/R-05/example/EX-01: a retired rule
// still resolves from the command line, flagged rather than refused. A
// reference filed against it outlives the decision to retire it.
func TestResolveURIResolvesARetiredRule(t *testing.T) {
	var out bytes.Buffer
	if err := resolveURI(&out, newTestMaster(t), "livt://mapping/demo/rule/R-02", formatJSON, ""); err != nil {
		t.Fatalf("resolving a retired rule failed: %v", err)
	}

	var payload struct {
		Rule struct {
			Name    string `json:"name"`
			Retired bool   `json:"retired"`
		} `json:"rule"`
	}
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if !payload.Rule.Retired {
		t.Error("a retired rule should resolve carrying retired")
	}
	if payload.Rule.Name != "退役したルール" {
		t.Errorf("name = %q, want the retired rule's text kept", payload.Rule.Name)
	}
}

func TestResolveURIRejectsUnknownFormat(t *testing.T) {
	err := resolveURI(&bytes.Buffer{}, newTestMaster(t), "livt://story/demo", "yaml", "")
	if err == nil || !strings.Contains(err.Error(), "unknown --format") {
		t.Errorf("unknown format error = %v, want it to name the flag", err)
	}
}

func keysOf(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
