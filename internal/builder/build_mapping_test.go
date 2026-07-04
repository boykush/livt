package builder

import (
	"bytes"
	"strings"
	"testing"

	"github.com/boykush/livt/internal/domain"
)

func TestRenderMappingRuleCarriesIDAnchorAndBadge(t *testing.T) {
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{{ID: "R-01", Name: "Activities and steps can be overviewed"}},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, em, "Story", "", nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	if !strings.Contains(html, `id="rule-R-01"`) {
		t.Fatal("expected rule sticky to carry an id anchor")
	}
	if !strings.Contains(html, `href="#rule-R-01" data-copy-link`) {
		t.Fatal("expected a one-click copy-link trigger on the rule sticky")
	}
	if !strings.Contains(html, "🔗") {
		t.Fatal("expected a link emoji to signal the badge is copyable")
	}
	if !strings.Contains(html, "R-01</a>") {
		t.Fatal("expected the rule ID to be displayed as a badge")
	}
}

func TestRenderMappingMarksAutomatedRules(t *testing.T) {
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{
			{ID: "R-01", Name: "An automated rule", Automated: true},
			{ID: "R-02", Name: "A rule not yet automated"},
		},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, em, "Story", "", nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	if got := strings.Count(html, "✓ automated"); got != 1 {
		t.Fatalf("automated badge rendered %d times, want once (only on the automated rule)", got)
	}
	if !strings.Contains(html, "Automated rule") {
		t.Fatal("expected the legend to explain the automated mark")
	}
}

func TestRenderMappingLinksRuleIssues(t *testing.T) {
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{
			{ID: "R-01", Name: "A linked rule", Issues: []string{
				"https://github.com/boykush/livt/issues/25",
				"https://github.com/boykush/other/issues/7",
			}},
			{ID: "R-02", Name: "An unlinked rule"},
		},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, em, "Story", "", nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	if !strings.Contains(html, `href="https://github.com/boykush/livt/issues/25"`) {
		t.Fatal("expected the rule sticky to link out to its recorded issue")
	}
	if !strings.Contains(html, "livt#25") || !strings.Contains(html, "other#7") {
		t.Fatal("expected issue links labelled as repo#number")
	}
	if got := strings.Count(html, `target="_blank"`); got != 2 {
		t.Fatalf("outbound links rendered %d times, want 2 (only on the linked rule)", got)
	}
}

func TestIssueLabelFallsBackToHost(t *testing.T) {
	cases := map[string]string{
		"https://github.com/boykush/livt/issues/25": "livt#25",
		"https://tracker.example.com/tickets/9":     "tracker.example.com",
		"not-a-url":                                 "not-a-url",
	}
	for url, want := range cases {
		if got := issueLabel(url); got != want {
			t.Errorf("issueLabel(%q) = %q, want %q", url, got, want)
		}
	}
}

func TestRenderMappingRuleWithoutIDOmitsAnchor(t *testing.T) {
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{{Name: "A rule without an ID"}},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, em, "Story", "", nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	if strings.Contains(html, `href="#rule-`) {
		t.Fatal("expected no copy-link anchor when a rule has no ID")
	}
	if strings.Contains(html, `id="rule-"`) {
		t.Fatal("expected no empty rule anchor when a rule has no ID")
	}
}
