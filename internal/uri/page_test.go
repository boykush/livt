package uri

import (
	"strings"
	"testing"
)

// livt://mapping/trace-test-to-rule/rule/R-02/example/EX-02: the derivation from
// a livt URI to its anchor on the page is written down in
// docs/src/reference/uri.md, and this is the one implementation the SSG and the
// CLI both read it through, so neither can drift off the table.
func TestPagesFollowTheDocumentedScheme(t *testing.T) {
	cases := []struct {
		uri  string
		got  string
		want string
	}{
		{Mapping("checkout"), MappingPage("checkout"), "mapping/checkout.html"},
		{Rule("checkout", "R-01"), RulePage("checkout", "R-01"), "mapping/checkout.html#rule-R-01"},
		{Example("checkout", "R-01", "EX-01"), ExamplePage("checkout", "R-01", "EX-01"), "mapping/checkout.html#rule-R-01-example-EX-01"},
		{Question("checkout", "Q-01"), QuestionPage("checkout", "Q-01"), "mapping/checkout.html#question-Q-01"},
		{Story("checkout"), StoryPage("checkout"), "story/checkout.html"},
		{StoryMap("discovery"), StoryMapPage("discovery"), "story-map/discovery.html"},
		{Term("livt-uri"), TermPage("livt-uri"), "ubiquitous.html#livt-uri"},
	}

	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s lands on %q, want %q", c.uri, c.got, c.want)
		}
	}
}

// livt://mapping/trace-test-to-rule/rule/R-01/example/EX-02: an example id is
// numbered within its rule, so one board holds an EX-01 under every rule. The
// anchor carries the rule for the same reason the URI does.
func TestExampleAnchorsSeparateTheSameIDUnderDifferentRules(t *testing.T) {
	if ExampleAnchor("R-01", "EX-01") == ExampleAnchor("R-02", "EX-01") {
		t.Fatal("EX-01 under two rules collapsed onto a single anchor")
	}
}

// livt://mapping/trace-test-to-rule/rule/R-02/example/EX-01: what the master
// stores is the livt URI, and the deployment URL is prefixed onto it only at
// render time — so no page here may resolve to a host or an absolute path on
// its own, or the master would be carrying a deployment it cannot know.
func TestPagesCarryNoDeployment(t *testing.T) {
	for _, got := range []string{
		MappingPage("checkout"),
		RulePage("checkout", "R-01"),
		ExamplePage("checkout", "R-01", "EX-01"),
		QuestionPage("checkout", "Q-01"),
		StoryPage("checkout"),
		StoryMapPage("discovery"),
		TermPage("livt-uri"),
	} {
		if strings.Contains(got, "://") || strings.HasPrefix(got, "/") {
			t.Errorf("%q resolves on its own; prefixing the deployment is the renderer's job", got)
		}
	}
}
