package builder

import (
	"bytes"
	"strings"
	"testing"

	"github.com/boykush/livt/internal/domain"
)

// keyedBoard is a board holding one of every sticky kind that carries an ID.
func keyedBoard() *domain.ExampleMapping {
	return &domain.ExampleMapping{
		Rules: []domain.Rule{{
			ID:       "R-01",
			Name:     "Activities and steps can be overviewed",
			Examples: []domain.Example{{ID: "EX-01", Name: "An example"}},
		}},
		Questions: []domain.Question{{ID: "Q-01", Text: "An open question"}},
	}
}

// livt://mapping/trace-test-to-rule/rule/R-03/example/EX-01 and
// livt://mapping/trace-test-to-rule/rule/R-03/example/EX-02: rule, example and
// question stickies alike show their own ID, and that badge is the trigger that
// copies the sticky's own URL.
func TestRenderMappingEveryStickyCarriesACopyableIDBadge(t *testing.T) {
	var buf bytes.Buffer
	if err := renderMapping(&buf, keyedBoard(), "Story", "", nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	for _, anchor := range []string{"rule-R-01", "rule-R-01-example-EX-01", "question-Q-01"} {
		if !strings.Contains(html, `id="`+anchor+`"`) {
			t.Errorf("expected a sticky anchored at %s", anchor)
		}
		if !strings.Contains(html, `href="#`+anchor+`" data-copy-link`) {
			t.Errorf("expected a one-click copy-link trigger aimed at %s", anchor)
		}
	}
	// The badge shows the ID as the master numbers it, so an example reads as the
	// rule-local EX-01 even though the link behind it is qualified by the rule.
	for _, label := range []string{">#R-01</a>", ">#EX-01</a>", ">#Q-01</a>"} {
		if !strings.Contains(html, label) {
			t.Errorf("expected a badge labelled %s", label)
		}
	}
}

// livt://mapping/trace-test-to-rule/rule/R-03/example/EX-03: the badge stays
// monochrome and tinted to its own sticky, because a dense board carries 30-40
// of them and an emoji renders full-colour whatever the card around it does.
func TestRenderMappingIDBadgesAreTintedNotEmoji(t *testing.T) {
	var buf bytes.Buffer
	if err := renderMapping(&buf, keyedBoard(), "Story", "", nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	if strings.Contains(html, "🔗") {
		t.Error("expected the ID badges to carry no emoji")
	}
	for _, tint := range []string{"text-blue-400/70", "text-green-400/70", "text-red-400/70"} {
		if !strings.Contains(html, tint) {
			t.Errorf("expected a badge tinted %s to follow its own sticky", tint)
		}
	}
}

// livt://mapping/trace-test-to-rule/rule/R-03/example/EX-01: an example sticky
// is a link target like the other two, so arriving at one flashes the card
// instead of leaving the reader to work out which one the URL meant.
func TestRenderMappingFlashesEveryLinkableStickyKind(t *testing.T) {
	var buf bytes.Buffer
	if err := renderMapping(&buf, keyedBoard(), "Story", "", nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	for _, kind := range []string{"rule-card", "example-card", "question-card"} {
		if !strings.Contains(html, "."+kind+":target") {
			t.Errorf("expected a %s to flash when a deep link lands on it", kind)
		}
		if !strings.Contains(html, `class="`+kind+` relative`) {
			t.Errorf("expected the %s markup to carry the class the flash keys on", kind)
		}
	}
}

// livt://mapping/trace-test-to-rule/rule/R-02/example/EX-02: EX-01 recurs under
// every rule of a board, so an anchor keyed on the example ID alone would send
// both links to whichever card happened to be rendered first.
func TestRenderMappingExampleAnchorsCarryTheirRule(t *testing.T) {
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{
			{ID: "R-01", Name: "First rule", Examples: []domain.Example{{ID: "EX-01", Name: "First rule's example"}}},
			{ID: "R-02", Name: "Second rule", Examples: []domain.Example{{ID: "EX-01", Name: "Second rule's example"}}},
		},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, em, "Story", "", nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	for _, anchor := range []string{"rule-R-01-example-EX-01", "rule-R-02-example-EX-01"} {
		if !strings.Contains(html, `id="`+anchor+`"`) {
			t.Errorf("expected the two EX-01s to be told apart by %s", anchor)
		}
	}
	if strings.Contains(html, `id="example-EX-01"`) {
		t.Error("expected no rule-blind example anchor, which both EX-01s would answer to")
	}
}

// An example under a rule with no ID has nothing to qualify its anchor with, so
// it stays unlinkable rather than claiming an ambiguous one.
func TestRenderMappingExampleUnderUnkeyedRuleOmitsAnchor(t *testing.T) {
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{{Name: "A rule without an ID", Examples: []domain.Example{{ID: "EX-01", Name: "An example"}}}},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, em, "Story", "", nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	if strings.Contains(html, `id="rule--example-EX-01"`) || strings.Contains(html, `href="#rule--example`) {
		t.Fatal("expected no example anchor when the enclosing rule has no ID")
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
func TestCollectTasksSplitsQuestionsFromUnautomatedRules(t *testing.T) {
	em := &domain.ExampleMapping{
		StoryKey: domain.StoryKey{Value: "overview-open-questions"},
		Rules: []domain.Rule{
			{ID: "R-01", Name: "A proven rule", Automated: true},
			{ID: "R-02", Name: "A rule not yet proven"},
		},
		Questions: []domain.Question{{ID: "Q-01", Text: "An open question"}},
	}

	out := collectTasks(em, "疑問を見渡す", "story/overview-open-questions.html")

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
func TestCollectTasksLinksItemsToTheirStickies(t *testing.T) {
	em := &domain.ExampleMapping{
		StoryKey:  domain.StoryKey{Value: "checkout"},
		Rules:     []domain.Rule{{ID: "R-02", Name: "A rule"}, {Name: "A rule with no ID"}},
		Questions: []domain.Question{{ID: "Q-01", Text: "A question"}},
	}

	out := collectTasks(em, "Checkout", "story/checkout.html")

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

// livt://mapping/trace-test-to-rule/rule/R-05/example/EX-02: retired stickies are
// off the board, whichever kind they are — a retired rule, a retired example
// under a live rule, and a retired question alike.
func TestRenderMappingOmitsRetiredStickies(t *testing.T) {
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{
			{ID: "R-01", Name: "現役のルール", Examples: []domain.Example{
				{ID: "EX-01", Name: "現役の実例"},
				{ID: "EX-02", Name: "退役した実例", Retired: true},
			}},
			{ID: "R-02", Name: "退役したルール", Retired: true},
		},
		Questions: []domain.Question{
			{ID: "Q-01", Text: "現役の疑問"},
			{ID: "Q-02", Text: "退役した疑問", Retired: true},
		},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, em, "Story", "", nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	for _, live := range []string{"現役のルール", "現役の実例", "現役の疑問"} {
		if !strings.Contains(html, live) {
			t.Errorf("board dropped a live sticky: %q", live)
		}
	}
	for _, retired := range []string{
		"退役したルール", "退役した実例", "退役した疑問",
		`id="rule-R-02"`, `id="question-Q-02"`, // nor their anchors
	} {
		if strings.Contains(html, retired) {
			t.Errorf("board still shows retired %q", retired)
		}
	}
}

// livt://mapping/trace-test-to-rule/rule/R-05/example/EX-02: a board whose only
// question is retired carries no Questions column at all — an empty one would
// read as an open question scrolled out of sight.
func TestRenderMappingDropsQuestionsColumnWhenEveryQuestionIsRetired(t *testing.T) {
	em := &domain.ExampleMapping{
		Questions: []domain.Question{{ID: "Q-01", Text: "退役した疑問", Retired: true}},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, em, "Story", "", nil); err != nil {
		t.Fatal(err)
	}

	// "Questions" is the column header; the legend says "Question".
	if strings.Contains(buf.String(), "Questions") {
		t.Fatal("expected no Questions column when every question is retired")
	}
}

// livt://mapping/trace-test-to-rule/rule/R-05/example/EX-02: retired items are
// not unfinished work. A retired question left in "open questions" could never
// be closed by a conversation, nor a retired rule by a test.
func TestCollectTasksSkipsRetiredItems(t *testing.T) {
	em := &domain.ExampleMapping{
		StoryKey: domain.StoryKey{Value: "trace-test-to-rule"},
		Rules: []domain.Rule{
			{ID: "R-01", Name: "現役の未自動化ルール"},
			{ID: "R-02", Name: "退役したルール", Retired: true},
		},
		Questions: []domain.Question{
			{ID: "Q-01", Text: "現役の疑問"},
			{ID: "Q-02", Text: "退役した疑問", Retired: true},
		},
	}

	out := collectTasks(em, "テストからルールを辿る", "story/trace-test-to-rule.html")

	if len(out.Questions) != 1 || out.Questions[0].ID != "Q-01" {
		t.Errorf("questions = %+v, want only the live Q-01", out.Questions)
	}
	if len(out.UnautomatedRules) != 1 || out.UnautomatedRules[0].ID != "R-01" {
		t.Errorf("un-automated rules = %+v, want only the live R-01", out.UnautomatedRules)
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
