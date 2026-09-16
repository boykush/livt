package builder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/boykush/livt/internal/diff"
	"github.com/boykush/livt/internal/i18n"
	"github.com/boykush/livt/internal/uri"
)

func ruleChange(id string, became diff.Became) diff.Change {
	return diff.Change{
		URI:    "livt://mapping/checkout/rule/" + id,
		Kind:   uri.KindRule,
		Parent: "livt://mapping/checkout",
		Title:  id,
		Became: became,
		Lines:  []diff.Line{{Op: diff.OpAdd, Field: diff.Field{Value: "a rule"}}},
	}
}

func renderedDiff(t *testing.T, result *diff.Result) string {
	t.Helper()
	b := emptyDirsBuilder(t)
	b.diffResult = result
	if err := b.buildDiff(); err != nil {
		t.Fatal(err)
	}
	out, err := os.ReadFile(filepath.Join(b.OutDir, diffPage))
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

// livt://mapping/review-diff-between-revisions/rule/R-01/example/EX-01: a build
// given no revisions renders the site it always did. The page is not merely
// unlinked — it is not there, so an output directory reused from a diff build
// cannot go on serving a comparison nobody asked for.
func TestBuildWithoutRevisionsLeavesNoDiffPage(t *testing.T) {
	b := emptyDirsBuilder(t)
	stale := filepath.Join(b.OutDir, diffPage)
	writeFile(t, stale, "<html>an earlier run's diff</html>")

	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Fatalf("%s survived a build given no revisions (stat error: %v)", diffPage, err)
	}
	index, err := os.ReadFile(filepath.Join(b.OutDir, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(index), diffPage) {
		t.Error("the sidebar links a diff on a build that was given no revisions")
	}
}

// livt://mapping/review-diff-between-revisions/rule/R-01/example/EX-02: the
// page says which revisions it was built between, since the range was given to
// the build and is nowhere else on the site.
func TestBuildDiffNamesTheRevisionsItCompared(t *testing.T) {
	page := renderedDiff(t, &diff.Result{Base: "abc1234", Changes: []diff.Change{ruleChange("R-01", diff.BecameAdded)}, Added: 1})

	if !strings.Contains(page, "abc1234") {
		t.Error("the page does not say which revisions it compared")
	}
}

// livt://mapping/review-diff-between-revisions/rule/R-05/example/EX-04: the
// diff takes no entry in the nav. It is not a kind of thing the livt repository
// holds — it is something said about the things it does — and an entry beside
// the five resource types would read as a sixth.
func TestBuildDiffTakesNoSidebarEntry(t *testing.T) {
	page := renderedDiff(t, &diff.Result{Base: "abc1234", Changes: []diff.Change{ruleChange("R-01", diff.BecameAdded)}, Added: 1})

	nav := page[strings.Index(page, "<nav"):strings.Index(page, "</nav>")]
	if strings.Contains(nav, diffPage) {
		t.Error("the sidebar carries a diff entry")
	}
}

// livt://mapping/review-diff-between-revisions/rule/R-05/example/EX-03: the way
// into the diff is the thing that changed. Reading a board, you want to know
// whether this rule is new — and from there, what it used to say.
func TestAChangedRuleIsMarkedOnItsBoardAndLeadsIntoTheDiff(t *testing.T) {
	b := emptyDirsBuilder(t)
	if err := b.resetGeneratedDirs(); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(b.MappingsDir, "checkout.yaml"), "rules:\n  - id: R-01\n    name: changed\n  - id: R-02\n    name: untouched\n")
	b.diffByURI = map[string]diff.Became{"livt://mapping/checkout/rule/R-01": diff.BecameChanged}

	if _, _, _, err := b.buildMappings(); err != nil {
		t.Fatal(err)
	}
	out, err := os.ReadFile(filepath.Join(b.OutDir, "mapping", "checkout.html"))
	if err != nil {
		t.Fatal(err)
	}
	page := string(out)

	if !strings.Contains(page, `href="../diff.html#mapping/checkout/rule/R-01"`) {
		t.Error("the changed rule's sticky does not lead into the diff")
	}
	if strings.Contains(page, "mapping/checkout/rule/R-02") {
		t.Error("an unchanged rule is marked as though it had changed")
	}
}

// A build given no revisions marks nothing. Every page renders the mark through
// {{with}}, so this is the same absence that keeps the boards unchanged.
func TestABoardCarriesNoMarksWithoutADiff(t *testing.T) {
	b := emptyDirsBuilder(t)
	if err := b.resetGeneratedDirs(); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(b.MappingsDir, "checkout.yaml"), "rules:\n  - id: R-01\n    name: a rule\n")

	if _, _, _, err := b.buildMappings(); err != nil {
		t.Fatal(err)
	}
	out, err := os.ReadFile(filepath.Join(b.OutDir, "mapping", "checkout.html"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), diffPage) {
		t.Error("a board built without revisions points at a diff")
	}
}

// livt://mapping/review-diff-between-revisions/rule/R-02/example/EX-05: rules,
// examples and questions gather under the mapping they hang off. A mapping
// whose own fields did not change is listed once, as the heading for the ones
// that did — not once per child.
func TestBuildDiffGathersMappingChildrenUnderOneMapping(t *testing.T) {
	page := renderedDiff(t, &diff.Result{
		Base: "abc1234",
		Changes: []diff.Change{
			ruleChange("R-01", diff.BecameAdded),
			ruleChange("R-02", diff.BecameChanged),
		},
		Added: 1, Changed: 1,
	})

	if got := strings.Count(page, "livt://mapping/checkout<"); got != 1 {
		t.Errorf("the mapping heads its rules %d times, want once", got)
	}
	for _, id := range []string{"R-01", "R-02"} {
		if !strings.Contains(page, "livt://mapping/checkout/rule/"+id) {
			t.Errorf("rule %s is not listed under its mapping", id)
		}
	}
}

// livt://mapping/review-diff-between-revisions/rule/R-05/example/EX-02: a
// removed URI has nowhere to go. The site holds no page for it, so it is
// rendered as text rather than as a link that would land on a 404.
func TestBuildDiffLeavesARemovedURIUnlinked(t *testing.T) {
	removed := ruleChange("R-01", diff.BecameWithdrawn)
	removed.Lines = []diff.Line{{Op: diff.OpDel, Field: diff.Field{Value: "a rule"}}}
	page := renderedDiff(t, &diff.Result{Base: "abc1234", Changes: []diff.Change{removed}, Withdrawn: 1})

	if strings.Contains(page, `href="mapping/checkout.html#rule-R-01"`) {
		t.Error("a removed URI is linked to a page the site does not hold")
	}
	if !strings.Contains(page, "livt://mapping/checkout/rule/R-01") {
		t.Error("a removed URI is not listed at all")
	}
}

// livt://mapping/review-diff-between-revisions/rule/R-02/example/EX-04: with
// nothing between the revisions, the page says so rather than rendering an
// empty frame a reader would take for a build that went wrong.
func TestBuildDiffSaysWhenNothingChanged(t *testing.T) {
	page := renderedDiff(t, &diff.Result{Base: "abc1234"})

	if !strings.Contains(page, "Nothing changed between these revisions.") {
		t.Error("an empty diff does not say that nothing changed")
	}
}

// livt://mapping/review-diff-between-revisions/rule/R-06/example/EX-02: what
// livt spells in a file — a status, an automation — reaches the page in the
// words the site already uses for it, in the site's own language. "accepted" is
// how the file is written, not how the rule is read.
func TestBuildDiffPutsTheSitesOwnWordsOnLivtsFields(t *testing.T) {
	changed := ruleChange("R-01", diff.BecameChanged)
	changed.Lines = []diff.Line{
		{Op: diff.OpDel, Field: diff.Field{Label: diff.LabelStatus, Translate: true, Value: "diff.status.proposed"}},
		{Op: diff.OpAdd, Field: diff.Field{Label: diff.LabelStatus, Translate: true, Value: "diff.status.accepted"}},
	}
	result := &diff.Result{Base: "abc1234", Changes: []diff.Change{changed}, Changed: 1}

	b := emptyDirsBuilder(t)
	b.Lang = i18n.Ja
	b.diffResult = result
	if err := b.buildDiff(); err != nil {
		t.Fatal(err)
	}
	out, err := os.ReadFile(filepath.Join(b.OutDir, diffPage))
	if err != nil {
		t.Fatal(err)
	}
	page := string(out)

	for _, want := range []string{"状態: 提案中", "状態: 合意済み"} {
		if !strings.Contains(page, want) {
			t.Errorf("the page does not read %q", want)
		}
	}
	if strings.Contains(page, "diff.status.") || strings.Contains(page, "status: accepted") {
		t.Error("the page shows a message key or the file's own spelling")
	}
}

// livt://mapping/review-diff-between-revisions/rule/R-05/example/EX-05: what
// left the board goes back on it, in the language the board already has for a
// sticky that is not spec. A count says something went and nothing about what it
// said, which is the reviewer's actual question.
func TestABoardPutsBackWhatLeftItInThisDiff(t *testing.T) {
	b := emptyDirsBuilder(t)
	if err := b.resetGeneratedDirs(); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(b.MappingsDir, "checkout.yaml"),
		"rules:\n  - id: R-01\n    name: still spec\n  - id: R-02\n    name: no longer spec\n    status: retired\n"+
			"  - id: R-03\n    name: closed long ago\n    status: retired\n")
	// R-02 was retired in this range and R-04 deleted outright; R-03 closed
	// before it and is no business of this diff.
	b.diffResult = &diff.Result{Changes: []diff.Change{
		{URI: "livt://mapping/checkout/rule/R-02", Parent: "livt://mapping/checkout", Kind: uri.KindRule, Status: diff.StatusModified, Became: diff.BecameWithdrawn},
		{URI: "livt://mapping/checkout/rule/R-04", Parent: "livt://mapping/checkout", Kind: uri.KindRule, Status: diff.StatusRemoved, Became: diff.BecameWithdrawn,
			Lines: []diff.Line{{Op: diff.OpDel, Field: diff.Field{Value: "deleted outright"}}}},
	}}
	b.diffByURI = map[string]diff.Became{
		"livt://mapping/checkout/rule/R-02": diff.BecameWithdrawn,
		"livt://mapping/checkout/rule/R-04": diff.BecameWithdrawn,
	}

	if _, _, _, err := b.buildMappings(); err != nil {
		t.Fatal(err)
	}
	out, err := os.ReadFile(filepath.Join(b.OutDir, "mapping", "checkout.html"))
	if err != nil {
		t.Fatal(err)
	}
	page := string(out)

	for _, want := range []string{"no longer spec", "deleted outright"} {
		if !strings.Contains(page, want) {
			t.Errorf("the board does not put back %q", want)
		}
	}
	if strings.Contains(page, "closed long ago") {
		t.Error("a rule closed before this range was put back; it left before the diff did")
	}
	// One sticky, one word. Retiring a rule is a modified record, so the old
	// cut wore both "Changed" and a second stamp saying it was retired — one
	// thing said twice, and the reader left to pick which is the answer.
	if got := strings.Count(page, ">Withdrawn</a>"); got != 2 {
		t.Errorf("%d stickies read as withdrawn, want the two this range lost", got)
	}
	if strings.Contains(page, ">Changed</a>") {
		t.Error("a withdrawn sticky also reads as changed")
	}
}

// A build given no revisions draws the active view and nothing else: the board
// says what the spec asks for today, and a retired rule is not part of that.
func TestABoardWithoutADiffPutsNothingBack(t *testing.T) {
	b := emptyDirsBuilder(t)
	if err := b.resetGeneratedDirs(); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(b.MappingsDir, "checkout.yaml"),
		"rules:\n  - id: R-01\n    name: still spec\n  - id: R-02\n    name: no longer spec\n    status: retired\n")

	if _, _, _, err := b.buildMappings(); err != nil {
		t.Fatal(err)
	}
	out, err := os.ReadFile(filepath.Join(b.OutDir, "mapping", "checkout.html"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "no longer spec") {
		t.Error("a board built without revisions shows a retired rule")
	}
}

// livt://mapping/review-diff-between-revisions/rule/R-03/example/EX-07: the two
// colours answer different questions and are allowed to disagree. Retiring a
// rule adds the line that retires it, so the body is green while the entry is a
// withdrawal — the entry is coloured on its own rail rather than by the lines
// inside it.
func TestAWithdrawnEntryReadsAsOneEvenWhenItsLinesWereAdded(t *testing.T) {
	withdrawn := ruleChange("R-01", diff.BecameWithdrawn)
	withdrawn.Lines = []diff.Line{
		{Op: diff.OpContext, Field: diff.Field{Value: "the statement that stops holding"}},
		{Op: diff.OpAdd, Field: diff.Field{Label: diff.LabelRetired, Translate: true}},
	}
	page := renderedDiff(t, &diff.Result{Base: "abc1234", Changes: []diff.Change{withdrawn}, Withdrawn: 1})

	if !strings.Contains(page, "entry-withdrawn") {
		t.Error("the entry is not coloured by what became of it")
	}
	if !strings.Contains(page, "diff-add") {
		t.Error("the added line lost the colour the record earned it")
	}
	if !strings.Contains(page, `own-text">the statement that stops holding`) {
		t.Error("the withdrawn statement is not the line struck through")
	}
}
