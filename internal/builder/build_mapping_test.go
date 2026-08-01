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

// overview-open-questions R-02 EX-02: a question sticky is linkable, so the home
// page can send the reader to the exact red card, not just the board.
func TestRenderMappingQuestionCarriesIDAnchor(t *testing.T) {
	em := &domain.ExampleMapping{
		Questions: []domain.Question{{ID: "Q-01", Text: "解決した疑問はどう扱うか"}, {Text: "An unkeyed question"}},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, em, "Story", "", nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	if !strings.Contains(html, `id="question-Q-01"`) {
		t.Fatal("expected the question sticky to carry an id anchor")
	}
	if strings.Contains(html, `id="question-"`) {
		t.Fatal("expected no empty anchor when a question has no ID")
	}
}

// overview-open-questions R-01 & overview-unautomated-rules R-01: a mapping
// contributes its questions and only its un-automated rules; automated rules are
// finished and stay off the list.
func TestCollectOutstandingSplitsQuestionsFromUnautomatedRules(t *testing.T) {
	em := &domain.ExampleMapping{
		StoryKey: domain.StoryKey{Value: "overview-open-questions"},
		Rules: []domain.Rule{
			{ID: "R-01", Name: "A proven rule", Automated: true},
			{ID: "R-02", Name: "A rule not yet proven"},
		},
		Questions: []domain.Question{{ID: "Q-01", Text: "An open question"}},
	}

	out := collectOutstanding(em, "疑問を見渡す", "story/overview-open-questions.html")

	if len(out.Questions) != 1 || out.Questions[0].Text != "An open question" {
		t.Fatalf("questions = %+v, want the one open question", out.Questions)
	}
	if len(out.UnautomatedRules) != 1 || out.UnautomatedRules[0].Text != "A rule not yet proven" {
		t.Fatalf("un-automated rules = %+v, want only the unproven rule", out.UnautomatedRules)
	}
	for _, item := range append(out.Questions, out.UnautomatedRules...) {
		if item.StoryName != "疑問を見渡す" {
			t.Errorf("item %q lost the story it came from (R-02 EX-01)", item.Text)
		}
	}
}

// overview-open-questions R-02 EX-02 & overview-unautomated-rules R-02 EX-02:
// each item deep-links to its own sticky, using that sticky's anchor scheme.
func TestCollectOutstandingLinksItemsToTheirStickies(t *testing.T) {
	em := &domain.ExampleMapping{
		StoryKey:  domain.StoryKey{Value: "checkout"},
		Rules:     []domain.Rule{{ID: "R-02", Name: "A rule"}, {Name: "A rule with no ID"}},
		Questions: []domain.Question{{ID: "Q-01", Text: "A question"}},
	}

	out := collectOutstanding(em, "Checkout", "story/checkout.html")

	if got := out.Questions[0].MappingPath; got != "mapping/checkout.html#question-Q-01" {
		t.Errorf("question link = %q, want the question sticky's anchor", got)
	}
	if got := out.UnautomatedRules[0].MappingPath; got != "mapping/checkout.html#rule-R-02" {
		t.Errorf("rule link = %q, want the rule sticky's anchor", got)
	}
	// An unkeyed sticky renders no anchor, so its card aims at the board itself.
	if got := out.UnautomatedRules[1].MappingPath; got != "mapping/checkout.html" {
		t.Errorf("unkeyed rule link = %q, want the board with no fragment", got)
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
