package builder

// livt:automates livt://mapping/overview-open-questions/rule/R-02
// livt:automates livt://mapping/trace-test-to-rule/rule/R-03

import (
	"bytes"
	"html"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/i18n"
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

// livt:automates livt://mapping/trace-test-to-rule/rule/R-03/example/EX-01
// livt:automates livt://mapping/trace-test-to-rule/rule/R-03/example/EX-02
// Rule, example and question stickies alike show their own ID, and that badge
// is the trigger that copies the sticky's own URL.
func TestRenderMappingEveryStickyCarriesACopyableIDBadge(t *testing.T) {
	var buf bytes.Buffer
	if err := renderMapping(&buf, i18n.En, board{Mapping: keyedBoard().Active()}, "Story", "", nil, nil, nil); err != nil {
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
	// The badge shows the ID as the livt repository numbers it, so an example reads as the
	// rule-local EX-01 even though the link behind it is qualified by the rule.
	for _, label := range []string{">#R-01</a>", ">#EX-01</a>", ">#Q-01</a>"} {
		if !strings.Contains(html, label) {
			t.Errorf("expected a badge labelled %s", label)
		}
	}
}

// livt:automates livt://mapping/trace-test-to-rule/rule/R-03/example/EX-03
// The badge stays monochrome and tinted to its own sticky, because a dense
// board carries 30-40 of them and an emoji renders full-colour whatever the
// card around it does.
func TestRenderMappingIDBadgesAreTintedNotEmoji(t *testing.T) {
	var buf bytes.Buffer
	if err := renderMapping(&buf, i18n.En, board{Mapping: keyedBoard().Active()}, "Story", "", nil, nil, nil); err != nil {
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

// livt:automates livt://mapping/trace-test-to-rule/rule/R-03/example/EX-01
// An example sticky is a link target like the other two, so arriving at one
// flashes the card instead of leaving the reader to work out which one the
// URL meant.
func TestRenderMappingFlashesEveryLinkableStickyKind(t *testing.T) {
	var buf bytes.Buffer
	if err := renderMapping(&buf, i18n.En, board{Mapping: keyedBoard().Active()}, "Story", "", nil, nil, nil); err != nil {
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

// livt:automates livt://mapping/trace-test-to-rule/rule/R-02/example/EX-02
// livt:automates livt://mapping/trace-test-to-rule/rule/R-02/example/EX-01
// EX-01 recurs under every rule of a board, so an anchor keyed on the example
// ID alone would send both links to whichever card happened to be rendered
// first.
func TestRenderMappingExampleAnchorsCarryTheirRule(t *testing.T) {
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{
			{ID: "R-01", Name: "First rule", Examples: []domain.Example{{ID: "EX-01", Name: "First rule's example"}}},
			{ID: "R-02", Name: "Second rule", Examples: []domain.Example{{ID: "EX-01", Name: "Second rule's example"}}},
		},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, i18n.En, board{Mapping: em.Active()}, "Story", "", nil, nil, nil); err != nil {
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
	if err := renderMapping(&buf, i18n.En, board{Mapping: em.Active()}, "Story", "", nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	if strings.Contains(html, `id="rule--example-EX-01"`) || strings.Contains(html, `href="#rule--example`) {
		t.Fatal("expected no example anchor when the enclosing rule has no ID")
	}
}

// livt:automates livt://mapping/show-automation-status-per-rule/rule/R-01/example/EX-07
// The chip of a citing test is the sticky's mark, and no corner stamp repeats
// it. A rule the deprecated flag alone calls automated has no test to name,
// and keeps a chip so it does not lose the mark it had.
func TestRenderMappingMarksAutomatedRules(t *testing.T) {
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{
			{ID: "R-01", Name: "An automated rule", Automations: []domain.Automation{{Repo: "acme/impl", File: "x_test.go", Line: 1}}},
			{ID: "R-02", Name: "A rule not yet automated"},
			{ID: "R-03", Name: "A rule only the flag calls automated", AutomatedFlag: true},
		},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, i18n.En, board{Mapping: em.Active(), AutomationKnown: true}, "Story", "", nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	if got := strings.Count(html, `class="test-chip `); got != 2 {
		t.Fatalf("test chips rendered %d times, want 2: the cited rule and the flagged one", got)
	}
	if strings.Contains(html, "automated-badge") {
		t.Fatal("a corner stamp repeats what the chips already say")
	}
	if !strings.Contains(html, "Automated by a test") {
		t.Fatal("expected the legend to explain the automated mark")
	}
}

// livt:automates livt://mapping/derive-automation-status/rule/R-01/example/EX-02
// The ✓ legend promises a mark on a board where no sticky can carry one, so it
// goes with the rest of the axis when nothing has been collected.
func TestRenderMappingDrawsNoAutomationLegendWhenNothingWasCollected(t *testing.T) {
	em := &domain.ExampleMapping{Rules: []domain.Rule{{ID: "R-01", Name: "A rule"}}}

	var buf bytes.Buffer
	if err := renderMapping(&buf, i18n.En, board{Mapping: em.Active()}, "Story", "", nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	if html := buf.String(); strings.Contains(html, "Automated by a test") {
		t.Fatal("the legend explains a mark nothing on this board can carry")
	}
}

// livt:automates livt://mapping/derive-automation-status/rule/R-01/example/EX-02
// End to end, because the axis is spread over pages that each decide on their
// own: the board's legend, the Tasks list, and the opportunity figures.
func TestBuildLeavesAutomationOffWhenNothingWasCollected(t *testing.T) {
	b := emptyDirsBuilder(t)
	writeFile(t, filepath.Join(b.MappingsDir, "checkout.yaml"),
		"rules:\n"+
			"  - id: R-01\n"+
			"    name: 締め切り後の注文は受け付けない\n")

	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	if html := readRendered(t, filepath.Join(b.OutDir, "mapping", "checkout.html")); strings.Contains(html, "Automated by a test") {
		t.Fatal("the board explains an automated mark with no report behind it")
	}
	if html := readRendered(t, filepath.Join(b.OutDir, "tasks.html")); strings.Contains(html, "Un-automated Rules") {
		t.Fatal("the Tasks page calls a rule un-automated on a repository nobody has looked at")
	}
}

// livt:automates livt://mapping/show-automation-status-per-rule/rule/R-01/example/EX-04
// An example carries its own mark, and a rule's says nothing about it: a rule
// covered end to end sits over examples nobody has written a test for, and a
// rule with no test of its own sits over examples that all have one.
func TestRenderMappingMarksEachAxisOnItsOwn(t *testing.T) {
	cited := []domain.Automation{{Repo: "acme/impl", File: "x_test.go", Line: 1}}
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{
			{ID: "R-01", Name: "Cited rule, uncited examples", Automations: cited, Examples: []domain.Example{
				{ID: "EX-01", Name: "No test names this one"},
			}},
			{ID: "R-02", Name: "Uncited rule, cited examples", Examples: []domain.Example{
				{ID: "EX-01", Name: "A test names this one", Automations: cited},
			}},
		},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, i18n.En, board{Mapping: em.Active()}, "Story", "", nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	if got := strings.Count(html, `class="test-chip `); got != 2 {
		t.Fatalf("test chips rendered %d times, want 2: the cited rule and the cited example", got)
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
	if err := renderMapping(&buf, i18n.En, board{Mapping: em.Active()}, "Story", "", nil, nil, nil); err != nil {
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

// livt:automates livt://mapping/show-automation-status-per-rule/rule/R-01/example/EX-07
// livt:automates livt://mapping/show-automation-status-per-rule/rule/R-01/example/EX-08
// A sticky links each test that automates it by repository, file and line, on
// a filled chip with the executable-specification icon. An automation issue
// beside it is an outlined, dashed chip with an icon of its own: off GitHub
// an issue's label is nothing but its host, so the two must part by shape.
func TestRenderMappingTellsATestFromAnAutomationIssue(t *testing.T) {
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{{
			ID: "R-01", Name: "A rule with a test and an issue",
			Automations: []domain.Automation{{
				Repo: "acme/impl", File: "internal/x_test.go", Line: 12,
				URL: "https://github.com/acme/impl/blob/abc1234/internal/x_test.go#L12",
			}},
			Issues: []string{"https://tracker.example.com/tickets/9"},
		}},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, i18n.En, board{Mapping: em.Active(), AutomationKnown: true}, "Story", "", nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	page := html.UnescapeString(buf.String())

	type chip struct{ title, class, icon, text string }
	chipOf := func(href string) (chip, bool) {
		anchor := regexp.MustCompile(`<a href="` + regexp.QuoteMeta(href) + `"[^>]*title="([^"]*)" class="([^"]*)">(<svg.*?</svg>)<span[^>]*>([^<]*)</span></a>`)
		m := anchor.FindStringSubmatch(page)
		if m == nil {
			return chip{}, false
		}
		return chip{m[1], m[2], m[3], m[4]}, true
	}

	test, ok := chipOf("https://github.com/acme/impl/blob/abc1234/internal/x_test.go#L12")
	if !ok {
		t.Fatal("no chip links the citing test")
	}
	issue, ok := chipOf("https://tracker.example.com/tickets/9")
	if !ok {
		t.Fatal("no chip links the automation issue")
	}

	if test.text != "impl: x_test.go:12" || test.title != "Automated by tests" {
		t.Errorf("test chip reads %q titled %q", test.text, test.title)
	}
	if issue.text != "tracker.example.com" || issue.title != "Automation issue" {
		t.Errorf("issue chip reads %q titled %q", issue.text, issue.title)
	}
	if test.icon == issue.icon {
		t.Error("a test and an issue wear the same icon")
	}
	if !strings.Contains(issue.class, "border-dashed") || strings.Contains(issue.class, "bg-") {
		t.Errorf("issue chip %q is not an unfilled dashed outline", issue.class)
	}
	if !strings.Contains(test.class, "bg-") || strings.Contains(test.class, "border-dashed") {
		t.Errorf("test chip %q is not a filled chip", test.class)
	}
	if !strings.Contains(page, "Automation issue</span>") {
		t.Error("the legend does not explain the issue chip")
	}
}

// livt:automates livt://mapping/show-automation-status-per-rule/rule/R-01/example/EX-09
// livt:automates livt://mapping/show-automation-status-per-rule/rule/R-01/example/EX-10
// One test chip for both axes, the sticky saying which it sits on, and in none
// of the colours a test runner spends on results: livt knows a test cites the
// sticky, not whether it passes.
func TestRenderMappingPaintsTestChipsOneColourOffTheResults(t *testing.T) {
	cited := []domain.Automation{{Repo: "acme/impl", File: "x_test.go", Line: 1}}
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{{ID: "R-01", Name: "Cited rule", Automations: cited, Examples: []domain.Example{
			{ID: "EX-01", Name: "Cited example", Automations: cited},
		}}},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, i18n.En, board{Mapping: em.Active(), AutomationKnown: true}, "Story", "", nil, nil, nil); err != nil {
		t.Fatal(err)
	}

	classes := regexp.MustCompile(`class="(test-chip [^"]*)"`).FindAllStringSubmatch(buf.String(), -1)
	if len(classes) != 2 {
		t.Fatalf("test chips rendered %d times, want 2: the rule and the example", len(classes))
	}
	if classes[0][1] != classes[1][1] {
		t.Errorf("the rule's chip %q and the example's %q differ", classes[0][1], classes[1][1])
	}
	for _, result := range []string{"green", "red", "yellow", "amber", "emerald", "lime", "cyan"} {
		if strings.Contains(classes[0][1], result) {
			t.Errorf("test chip %q wears %s, a colour runners spend on results", classes[0][1], result)
		}
	}
}

// livt:automates livt://mapping/review-example-mapping-as-list/rule/R-02/example/EX-01
// The list lays a card out side by side, so the chips have to sit in the name's
// column: beside it, a rule cited by several tests squeezes its own name into a
// sliver. The same holds for an example's card.
func TestRenderMappingKeepsChipsUnderTheName(t *testing.T) {
	cited := []domain.Automation{{Repo: "acme/impl", File: "a_test.go", Line: 1}, {Repo: "acme/impl", File: "b_test.go", Line: 2}}
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{{
			ID: "R-01", Name: "A rule with several tests", Automations: cited,
			Issues:   []string{"https://github.com/acme/impl/issues/1"},
			Examples: []domain.Example{{ID: "EX-01", Name: "An example with several tests", Automations: cited}},
		}},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, i18n.En, board{Mapping: em.Active(), AutomationKnown: true}, "Story", "", nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	page := buf.String()

	for _, c := range []struct{ body, name, chips string }{
		{`class="rule-body`, "A rule with several tests", `class="rule-issues`},
		{`class="example-body`, "An example with several tests", `class="example-automations`},
	} {
		start := strings.Index(page, c.body)
		if start < 0 {
			t.Fatalf("no %s column", c.body)
		}
		// The column closes before the card does, so everything up to the
		// next card is its content.
		rest := page[start:]
		name, chips := strings.Index(rest, c.name), strings.Index(rest, c.chips)
		if name < 0 || chips < 0 || chips < name {
			t.Errorf("%s does not hold the name followed by its chips", c.body)
		}
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

// livt:automates livt://mapping/overview-open-questions/rule/R-02/example/EX-02
// A question sticky is linkable, so the home page can send the reader to the
// exact red card, not just the board.
func TestRenderMappingQuestionCarriesIDAnchor(t *testing.T) {
	em := &domain.ExampleMapping{
		Questions: []domain.Question{{ID: "Q-01", Text: "解決した疑問はどう扱うか"}, {Text: "An unkeyed question"}},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, i18n.En, board{Mapping: em.Active()}, "Story", "", nil, nil, nil); err != nil {
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

// livt:automates livt://mapping/overview-open-questions/rule/R-01
// And its mirror in overview-unautomated-rules: a mapping contributes its
// questions and only its un-automated rules; automated rules are finished and
// stay off the list.
func TestCollectTasksSplitsQuestionsFromUnautomatedRules(t *testing.T) {
	em := &domain.ExampleMapping{
		StoryKey: domain.StoryKey{Value: "overview-open-questions"},
		Rules: []domain.Rule{
			{ID: "R-01", Name: "A proven rule", Automations: []domain.Automation{{Repo: "acme/impl", File: "x_test.go", Line: 1}}},
			{ID: "R-02", Name: "A rule not yet proven"},
		},
		Questions: []domain.Question{{ID: "Q-01", Text: "An open question"}},
	}

	out := collectTasks(em, "疑問を見渡す", "story/overview-open-questions.html", true)

	if len(out.Questions) != 1 || out.Questions[0].Text != "An open question" {
		t.Fatalf("questions = %+v, want the one open question", out.Questions)
	}
	if len(out.UnautomatedRules) != 1 || out.UnautomatedRules[0].Text != "A rule not yet proven" {
		t.Fatalf("un-automated rules = %+v, want only the unproven rule", out.UnautomatedRules)
	}
	for _, item := range append(out.Questions, out.UnautomatedRules...) {
		if item.MappingName != "疑問を見渡す" {
			t.Errorf("item %q lost the story it came from "+
				"(livt://mapping/overview-open-questions/rule/R-02/example/EX-01 and its mirror)", item.Text)
		}
	}
}

// livt:automates livt://mapping/overview-open-questions/rule/R-02/example/EX-02
// And its mirror in overview-unautomated-rules: each item deep-links to its
// own sticky, using that sticky's anchor scheme.
func TestCollectTasksLinksItemsToTheirStickies(t *testing.T) {
	em := &domain.ExampleMapping{
		StoryKey:  domain.StoryKey{Value: "checkout"},
		Rules:     []domain.Rule{{ID: "R-02", Name: "A rule"}, {Name: "A rule with no ID"}},
		Questions: []domain.Question{{ID: "Q-01", Text: "A question"}},
	}

	out := collectTasks(em, "Checkout", "story/checkout.html", true)

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

// livt:automates livt://mapping/trace-test-to-rule/rule/R-05/example/EX-02
// livt:automates livt://mapping/propose-rule-before-agreement/rule/R-01/example/EX-04
// Closed stickies are off the board, whichever kind they are — a retired
// rule, a retired example under a live rule, and a retired question alike.
func TestRenderMappingOmitsRetiredStickies(t *testing.T) {
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{
			{ID: "R-01", Name: "現役のルール", Examples: []domain.Example{
				{ID: "EX-01", Name: "現役の実例"},
				{ID: "EX-02", Name: "退役した実例", Retired: true},
			}},
			{ID: "R-02", Name: "退役したルール", Status: domain.RuleRetired},
		},
		Questions: []domain.Question{
			{ID: "Q-01", Text: "現役の疑問"},
			{ID: "Q-02", Text: "退役した疑問", Retired: true},
		},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, i18n.En, board{Mapping: em.Active()}, "Story", "", nil, nil, nil); err != nil {
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
		// Nor their anchors: an id with no sticky is a link that lands nowhere.
		`id="rule-R-02"`, `id="rule-R-01-example-EX-02"`, `id="question-Q-02"`,
	} {
		if strings.Contains(html, retired) {
			t.Errorf("board still shows retired %q", retired)
		}
	}
}

// livt:automates livt://mapping/trace-test-to-rule/rule/R-05/example/EX-02
// A board whose only question is retired carries no Questions column at all —
// an empty one would read as an open question scrolled out of sight.
func TestRenderMappingDropsQuestionsColumnWhenEveryQuestionIsRetired(t *testing.T) {
	em := &domain.ExampleMapping{
		Questions: []domain.Question{{ID: "Q-01", Text: "退役した疑問", Retired: true}},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, i18n.En, board{Mapping: em.Active()}, "Story", "", nil, nil, nil); err != nil {
		t.Fatal(err)
	}

	// "Questions" is the column header; the legend says "Question".
	if strings.Contains(buf.String(), "Questions") {
		t.Fatal("expected no Questions column when every question is retired")
	}
}

// livt:automates livt://mapping/trace-test-to-rule/rule/R-05/example/EX-02
// livt:automates livt://mapping/propose-rule-before-agreement/rule/R-01/example/EX-04
// Closed items are not unfinished work. A retired question left in "open
// questions" could never be closed by a conversation, nor a retired rule by a
// test.
func TestCollectTasksSkipsRetiredItems(t *testing.T) {
	em := &domain.ExampleMapping{
		StoryKey: domain.StoryKey{Value: "trace-test-to-rule"},
		Rules: []domain.Rule{
			{ID: "R-01", Name: "現役の未自動化ルール"},
			{ID: "R-02", Name: "退役したルール", Status: domain.RuleRetired},
		},
		Questions: []domain.Question{
			{ID: "Q-01", Text: "現役の疑問"},
			{ID: "Q-02", Text: "退役した疑問", Retired: true},
		},
	}

	out := collectTasks(em, "テストからルールを辿る", "story/trace-test-to-rule.html", true)

	if len(out.Questions) != 1 || out.Questions[0].ID != "Q-01" {
		t.Errorf("questions = %+v, want only the live Q-01", out.Questions)
	}
	if len(out.UnautomatedRules) != 1 || out.UnautomatedRules[0].ID != "R-01" {
		t.Errorf("un-automated rules = %+v, want only the live R-01", out.UnautomatedRules)
	}
}

// stickyClass returns the class attribute of the sticky anchored at id.
func stickyClass(t *testing.T, html, id string) string {
	t.Helper()
	_, rest, ok := strings.Cut(html, `id="`+id+`" class="`)
	if !ok {
		t.Fatalf("no sticky anchored at %s", id)
	}
	class, _, _ := strings.Cut(rest, `"`)
	return class
}

// livt:automates livt://mapping/propose-rule-before-agreement/rule/R-02
// A proposed rule and its examples are drawn pale beside an agreed pair, and
// the rule is stamped so telling them apart never rests on colour alone. The
// legend explains the look.
func TestRenderMappingSetsProposedRulesApart(t *testing.T) {
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{
			{ID: "R-01", Name: "合意済みのルール", Examples: []domain.Example{{ID: "EX-01", Name: "合意済みの具体例"}}},
			{ID: "R-02", Name: "提案中のルール", Status: domain.RuleProposed, Examples: []domain.Example{{ID: "EX-01", Name: "提案中の具体例"}}},
		},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, i18n.En, board{Mapping: em.Active()}, "Story", "", nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	for _, id := range []string{"rule-R-02", "rule-R-02-example-EX-01"} {
		if class := stickyClass(t, html, id); !strings.Contains(class, "border-dashed") || strings.Contains(class, "shadow") {
			t.Errorf("%s class = %q, want the pale, unpinned look of a proposal", id, class)
		}
	}
	for _, id := range []string{"rule-R-01", "rule-R-01-example-EX-01"} {
		if class := stickyClass(t, html, id); strings.Contains(class, "border-dashed") {
			t.Errorf("%s class = %q, want the pinned look of an agreed sticky", id, class)
		}
	}
	if got := strings.Count(html, ">proposed</span>"); got != 1 {
		t.Errorf("proposed stamp rendered %d times, want once (only on the proposed rule)", got)
	}
	if !strings.Contains(html, "Proposed rule") {
		t.Error("expected the legend to explain the proposed look")
	}
}

// livt:automates livt://mapping/propose-rule-before-agreement/rule/R-03/example/EX-01
// EX-02: agreement is what closes a proposal, so it is listed apart from the
// un-automated rules even once a test covers it. A rejected one has closed on
// the same axis and is off the page (livt://mapping/propose-rule-before-
// agreement/rule/R-05).
func TestCollectTasksListsProposedRulesApart(t *testing.T) {
	em := &domain.ExampleMapping{
		StoryKey: domain.StoryKey{Value: "checkout"},
		Rules: []domain.Rule{
			{ID: "R-01", Name: "未自動化のルール"},
			{ID: "R-02", Name: "提案中のルール", Status: domain.RuleProposed},
			{ID: "R-03", Name: "テストが先に書かれた提案", Status: domain.RuleProposed, Automations: []domain.Automation{{Repo: "acme/impl", File: "x_test.go", Line: 1}}},
			{ID: "R-04", Name: "却下された提案", Status: domain.RuleRejected},
		},
	}

	out := collectTasks(em, "Checkout", "story/checkout.html", true)

	if len(out.UnautomatedRules) != 1 || out.UnautomatedRules[0].ID != "R-01" {
		t.Errorf("un-automated rules = %+v, want only the accepted R-01", out.UnautomatedRules)
	}
	var ids []string
	for _, item := range out.ProposedRules {
		ids = append(ids, item.ID)
		if item.Kind != "proposed-rule" || item.MappingPath != "mapping/checkout.html#rule-"+item.ID {
			t.Errorf("proposed item = %+v, want a proposed-rule linked to its own sticky", item)
		}
	}
	if got := strings.Join(ids, ","); got != "R-02,R-03" {
		t.Errorf("proposed rules = %s, want R-02,R-03 with the retired R-04 left off", got)
	}
}

// The sidebar's Tasks count is what the Tasks page lists, so a proposed rule
// counts while it waits for agreement, automated or not.
func TestTasksCountIncludesProposedRules(t *testing.T) {
	b := emptyDirsBuilder(t)
	writeFile(t, filepath.Join(b.MappingsDir, "checkout.yaml"),
		"rules:\n"+
			"  - id: R-01\n"+
			"    name: 自動化済みのルール\n"+
			"  - id: R-02\n"+
			"    name: テストが先に書かれた提案\n"+
			"    status: proposed\n"+
			"  - id: R-03\n"+
			"    name: 未自動化のルール\n")
	writeAutomations(t, b,
		"livt://mapping/checkout/rule/R-01",
		"livt://mapping/checkout/rule/R-02")

	c, err := b.computeCounts()
	if err != nil {
		t.Fatal(err)
	}
	if c.tasks != 2 {
		t.Errorf("tasks = %d, want 2: the proposed R-02 and the un-automated R-03", c.tasks)
	}
}

func TestRenderMappingRuleWithoutIDOmitsAnchor(t *testing.T) {
	em := &domain.ExampleMapping{
		Rules: []domain.Rule{{Name: "A rule without an ID"}},
	}

	var buf bytes.Buffer
	if err := renderMapping(&buf, i18n.En, board{Mapping: em.Active()}, "Story", "", nil, nil, nil); err != nil {
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

// foldedBoard holds each shape a row takes in the list: a rule with examples to
// fold, a rule with none, and a question.
func foldedBoard() *domain.ExampleMapping {
	return &domain.ExampleMapping{
		Rules: []domain.Rule{
			{
				ID:          "R-01",
				Name:        "A rule with examples",
				Examples:    []domain.Example{{ID: "EX-01", Name: "First example"}, {ID: "EX-02", Name: "Second example"}},
				Issues:      []string{"https://github.com/boykush/livt/issues/25"},
				Automations: []domain.Automation{{Repo: "acme/impl", File: "x_test.go", Line: 1}},
			},
			{ID: "R-02", Name: "A rule without examples"},
		},
		Questions: []domain.Question{{ID: "Q-01", Text: "An open question"}},
	}
}

func renderFoldedBoard(t *testing.T) string {
	t.Helper()
	var buf bytes.Buffer
	if err := renderMapping(&buf, i18n.En, board{Mapping: foldedBoard().Active()}, "Story", "", nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

// livt:automates livt://mapping/review-example-mapping-as-list/rule/R-01/example/EX-01
// A mapping opens on the board, with the list one press away.
func TestRenderMappingOpensOnTheBoardWithAListToSwitchTo(t *testing.T) {
	html := renderFoldedBoard(t)

	if !strings.Contains(html, `<html lang="en">`) {
		t.Fatal("expected the page served with no view chosen, which is the board")
	}
	for _, want := range []string{
		`data-view-toggle="board" aria-pressed="true"`,
		`data-view-toggle="list" aria-pressed="false"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("expected the view switch to carry %s", want)
		}
	}
}

// livt:automates livt://mapping/review-example-mapping-as-list/rule/R-02/example/EX-01
// EX-02: what the list shows on a rule's one row is all on the rule's own
// card, ahead of the examples folded under it.
func TestRenderMappingRuleRowCarriesItsWholeLine(t *testing.T) {
	html := renderFoldedBoard(t)

	start := strings.Index(html, `id="rule-R-01"`)
	end := strings.Index(html, `id="examples-0"`)
	if start < 0 || end < start {
		t.Fatalf("expected R-01's card ahead of its examples (card at %d, examples at %d)", start, end)
	}
	row := html[start:end]
	for _, want := range []string{">#R-01</a>", "A rule with examples", `class="test-chip `, "livt#25", "Examples</span>2</span>"} {
		if !strings.Contains(row, want) {
			t.Errorf("expected R-01's row to carry %q", want)
		}
	}
}

// livt:automates livt://mapping/review-example-mapping-as-list/rule/R-02/example/EX-03
// Questions come after every rule, one row each.
func TestRenderMappingListsQuestionsAfterTheRules(t *testing.T) {
	html := renderFoldedBoard(t)

	if strings.Index(html, `id="question-Q-01"`) < strings.Index(html, `id="rule-R-02"`) {
		t.Fatal("expected the question after the last rule")
	}
}

// livt:automates livt://mapping/review-example-mapping-as-list/rule/R-03/example/EX-01
// EX-03: a rule's examples sit inside the block its fold controls, which
// starts closed; a rule with no examples has nothing to fold.
func TestRenderMappingFoldsExamplesUnderTheirRule(t *testing.T) {
	html := renderFoldedBoard(t)

	if got := strings.Count(html, "data-fold aria-expanded="); got != 1 {
		t.Fatalf("rendered %d folds, want one: R-02 has no examples to fold", got)
	}
	if !strings.Contains(html, `data-fold aria-expanded="false" aria-controls="examples-0"`) {
		t.Fatal("expected R-01's fold to start closed, over its own examples")
	}
	order := []string{`id="examples-0"`, `id="rule-R-01-example-EX-01"`, `id="rule-R-01-example-EX-02"`, `id="rule-R-02"`}
	for i := 1; i < len(order); i++ {
		if strings.Index(html, order[i-1]) > strings.Index(html, order[i]) {
			t.Errorf("expected %s before %s, so the fold holds exactly its rule's examples", order[i-1], order[i])
		}
	}
}

// livt:automates livt://mapping/review-example-mapping-as-list/rule/R-04/example/EX-02
// The list re-lays the board's stickies rather than repeating them, so each
// anchor, and the copy-link aimed at it, exists once whichever view is
// showing.
func TestRenderMappingKeepsOneStickyPerAnchorAcrossViews(t *testing.T) {
	html := renderFoldedBoard(t)

	for _, anchor := range []string{"rule-R-01", "rule-R-01-example-EX-01", "rule-R-01-example-EX-02", "rule-R-02", "question-Q-01"} {
		if got := strings.Count(html, `id="`+anchor+`"`); got != 1 {
			t.Errorf("anchor %s rendered %d times, want once", anchor, got)
		}
		if got := strings.Count(html, `href="#`+anchor+`" data-copy-link`); got != 1 {
			t.Errorf("copy-link to %s rendered %d times, want once", anchor, got)
		}
	}
}

// A story file that carries no name: frontmatter must not leave the surfaces
// naming it blank. The mapping page's title, breadcrumb and heading show only
// this string, as does the Tasks page's story chip, so an empty one renders a
// heading with no text and a link with no label — a worse page than the
// missing-story case, which already falls back to the key.
func TestResolveStoryNameFallsBackToTheKey(t *testing.T) {
	b := emptyDirsBuilder(t)
	writeFile(t, filepath.Join(b.StoriesDir, "named.md"), "---\nname: 名前あり\n---\n\n本文\n")
	writeFile(t, filepath.Join(b.StoriesDir, "unnamed.md"), "本文だけ\n")

	for _, tc := range []struct {
		key, want string
	}{
		{"named", "名前あり"},
		{"unnamed", "unnamed"}, // a story file with no name: frontmatter
		{"missing", "missing"}, // no story file at all
	} {
		if got := b.resolveStoryName(domain.StoryKey{Value: tc.key}); got != tc.want {
			t.Errorf("resolveStoryName(%q) = %q, want %q", tc.key, got, tc.want)
		}
	}
}

// livt:automates livt://mapping/review-example-mapping-as-list/rule/R-02/example/EX-03
// The list heads its rules and its questions the way the Tasks page heads the
// two lists it makes of the same items.
func TestRenderMappingHeadsTheListsSections(t *testing.T) {
	html := renderFoldedBoard(t)

	for _, want := range []string{">Rules</h2>", ">Questions</h2>"} {
		if !strings.Contains(html, want) {
			t.Errorf("expected a section head %q", want)
		}
	}
}

// storySticky is the board's yellow sticky from its opening tag through the
// name it shows — the one place a board names itself besides its title.
func storySticky(t *testing.T, html string) string {
	t.Helper()
	i := strings.Index(html, "bg-yellow-100 border-l-4 border-yellow-400 p-3 rounded shadow font-bold")
	if i < 0 {
		t.Fatal("board has no yellow sticky")
	}
	start := strings.LastIndex(html[:i], "<")
	end := i + strings.Index(html[i:], "</")
	return html[start:end]
}

// livt:automates livt://mapping/name-example-mapping-itself/rule/R-01/example/EX-02
// livt:automates livt://mapping/name-example-mapping-itself/rule/R-02/example/EX-01
// A board is called by its own name, then by its story's, then by its key.
func TestBuildNamesABoardByItsOwnNameThenItsStorysThenItsKey(t *testing.T) {
	b := emptyDirsBuilder(t)
	writeFile(t, filepath.Join(b.MappingsDir, "named.yaml"), "name: 自分の名前\nrules: []\n")
	writeFile(t, filepath.Join(b.MappingsDir, "both.yaml"), "name: 自分の名前\nrules: []\n")
	writeFile(t, filepath.Join(b.StoriesDir, "both.md"), "---\nname: ストーリーの名前\n---\n")
	writeFile(t, filepath.Join(b.MappingsDir, "storied.yaml"), "rules: []\n")
	writeFile(t, filepath.Join(b.StoriesDir, "storied.md"), "---\nname: ストーリーの名前\n---\n")
	writeFile(t, filepath.Join(b.MappingsDir, "bare.yaml"), "rules: []\n")

	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	for key, want := range map[string]string{
		"named":   "自分の名前",
		"both":    "自分の名前",
		"storied": "ストーリーの名前",
		"bare":    "bare",
	} {
		html := readRendered(t, filepath.Join(b.OutDir, "mapping", key+".html"))
		if !strings.Contains(html, "<title>"+want+" - livt</title>") {
			t.Errorf("mapping %s is not titled %q", key, want)
		}
		if sticky := storySticky(t, html); !strings.Contains(sticky, want) {
			t.Errorf("mapping %s: yellow sticky %q does not read %q", key, sticky, want)
		}
	}
}

// livt:automates livt://mapping/name-example-mapping-itself/rule/R-01/example/EX-01
// With no story to lean on, the name is what every surface naming the board
// shows: its sticky, its tile, and the chip its unfinished items wear.
func TestBuildCallsAStorylessMappingByItsNameEverywhere(t *testing.T) {
	const name = "全角スペースでログインできない"
	b := emptyDirsBuilder(t)
	// The Tasks chip this asserts on is an un-automated rule, which is only a
	// task once something has looked.
	writeAutomations(t, b)
	writeFile(t, filepath.Join(b.MappingsDir, "fix-login.yaml"),
		"name: "+name+"\n"+
			"rules:\n"+
			"  - id: R-01\n"+
			"    name: パスワードの前後の空白は取り除かない\n")

	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	board := readRendered(t, filepath.Join(b.OutDir, "mapping", "fix-login.html"))
	if sticky := storySticky(t, board); !strings.HasPrefix(sticky, "<div") || !strings.Contains(sticky, name) {
		t.Errorf("yellow sticky = %q, want an unlinked sticky reading %q", sticky, name)
	}
	if index := readRendered(t, filepath.Join(b.OutDir, "index.html")); !strings.Contains(index, `aria-label="`+name+`"`) {
		t.Errorf("the overview tile is not named %q", name)
	}
	tasks := readRendered(t, filepath.Join(b.OutDir, "tasks.html"))
	if !strings.Contains(tasks, ">"+name+"</span>") || strings.Contains(tasks, ">fix-login<") {
		t.Errorf("the Tasks chip does not name the board %q", name)
	}
}

// livt:automates livt://mapping/name-example-mapping-itself/rule/R-02
// livt:automates livt://mapping/name-example-mapping-itself/rule/R-02/example/EX-02
// livt:automates livt://mapping/name-example-mapping-itself/rule/R-02/example/EX-03
// The two names are read independently: the board keeps its own, the story
// page keeps the story's, and the sticky still leads from one to the other.
func TestBuildKeepsAMappingsNameApartFromItsStorys(t *testing.T) {
	b := emptyDirsBuilder(t)
	writeFile(t, filepath.Join(b.MappingsDir, "checkout.yaml"), "name: 決済の境界\nrules: []\n")
	writeFile(t, filepath.Join(b.StoriesDir, "checkout.md"), "---\nname: 決済する\n---\n")

	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	sticky := storySticky(t, readRendered(t, filepath.Join(b.OutDir, "mapping", "checkout.html")))
	if !strings.HasPrefix(sticky, `<a href="../story/checkout.html"`) || !strings.Contains(sticky, "決済の境界") {
		t.Errorf("yellow sticky = %q, want the board's own name linking to the story page", sticky)
	}
	story := readRendered(t, filepath.Join(b.OutDir, "story", "checkout.html"))
	if !strings.Contains(story, "<title>決済する - livt</title>") || strings.Contains(story, "決済の境界") {
		t.Error("the story page should keep the story's own name")
	}
}
