package diff

// livt:automates livt://mapping/automate-from-master-in-impl-repos/rule/R-09
// livt:automates livt://mapping/review-diff-between-revisions/rule/R-01
// livt:automates livt://mapping/review-diff-between-revisions/rule/R-02
// livt:automates livt://mapping/review-diff-between-revisions/rule/R-04
// livt:automates livt://mapping/review-diff-between-revisions/rule/R-06

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/boykush/livt/internal/automation"
	"github.com/boykush/livt/internal/uri"
)

// repoDirs lays a livt repository out under root, in the layout the CLI builds
// one with, so a fixture and a real checkout are read the same way.
func repoDirs(root string) Dirs {
	return Dirs{
		Opportunities: filepath.Join(root, "opportunities"),
		Canvases:      filepath.Join(root, "discoveries", "opportunity-canvases"),
		Mappings:      filepath.Join(root, "discoveries", "example-mappings"),
		Stories:       filepath.Join(root, "stories"),
		USM:           filepath.Join(root, "discoveries", "usm"),
		Ubiquitous:    filepath.Join(root, "ubiquitous"),
		Automations:   filepath.Join(root, "automations"),
	}
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// mapping writes one example mapping and snapshots the repository holding it.
func mapping(t *testing.T, key, yaml string) *Snapshot {
	t.Helper()
	dirs := repoDirs(t.TempDir())
	write(t, filepath.Join(dirs.Mappings, key+".yaml"), yaml)
	snapshot, err := Scan(dirs)
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}

func uris(changes []Change) []string {
	out := make([]string, 0, len(changes))
	for _, c := range changes {
		out = append(out, c.URI)
	}
	return out
}

func only(t *testing.T, changes []Change) Change {
	t.Helper()
	if len(changes) != 1 {
		t.Fatalf("got %d changes %v, want exactly one", len(changes), uris(changes))
	}
	return changes[0]
}

// lineTexts renders a diff for assertion the way the package holds it: label
// keys unresolved, since translating them is the builder's job and not this
// package's to be tested on.
func lineTexts(lines []Line) []string {
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		text := l.Field.Value
		if l.Field.Label != "" {
			text = l.Field.Label + ": " + text
		}
		out = append(out, string(l.Op)+text)
	}
	return out
}

func status(s string) string {
	return LabelStatus + ": " + statusPrefix + s
}

const oneRule = `rules:
  - id: R-01
    name: first
    examples:
      - id: EX-01
        name: an example
`

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-02/example/EX-01
// The reworded rule is the change, and it is the only one. Its mapping and
// its example read the same as before, so neither is anything the reviewer is
// asked to look at.
func TestCompareReportsARewordedRuleOnItsOwnURI(t *testing.T) {
	base := mapping(t, "checkout", oneRule)
	head := mapping(t, "checkout", strings.Replace(oneRule, "name: first", "name: second", 1))

	change := only(t, Compare(base, head))
	if change.URI != "livt://mapping/checkout/rule/R-01" {
		t.Errorf("changed URI = %q, want the rule's", change.URI)
	}
	if change.Status != StatusModified {
		t.Errorf("status = %q, want %q", change.Status, StatusModified)
	}
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-02/example/EX-02
// An added rule is one added URI. Its mapping holds it now and did not
// before, but the mapping's own fields are untouched — counting it again
// there would report one decision twice.
func TestCompareCountsAnAddedRuleOnceRatherThanAlsoOnItsMapping(t *testing.T) {
	base := mapping(t, "checkout", oneRule)
	head := mapping(t, "checkout", oneRule+`  - id: R-02
    name: second
`)

	change := only(t, Compare(base, head))
	if change.URI != "livt://mapping/checkout/rule/R-02" || change.Status != StatusAdded {
		t.Errorf("got %s %q, want the new rule added", change.Status, change.URI)
	}
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-02/example/EX-04
// A revision compared against itself has nothing to show.
func TestCompareLeavesUnchangedURIsOut(t *testing.T) {
	base := mapping(t, "checkout", oneRule)
	head := mapping(t, "checkout", oneRule)

	if changes := Compare(base, head); len(changes) != 0 {
		t.Errorf("unchanged revisions produced %v, want no changes", uris(changes))
	}
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-02/example/EX-03
// Every kind of URI is diffed, not only the ones on an example mapping.
func TestCompareCoversEveryURIKind(t *testing.T) {
	dirs := repoDirs(t.TempDir())
	write(t, filepath.Join(dirs.Opportunities, "growth.md"), "---\nname: Growth\n---\n\nwhy it matters\n")
	write(t, filepath.Join(dirs.Canvases, "growth.yaml"), "problems:\n  - a problem\n")
	write(t, filepath.Join(dirs.USM, "growth.yaml"), "name: Growth map\nactivities: []\n")
	write(t, filepath.Join(dirs.Stories, "checkout.md"), "---\nname: Checkout\n---\n\nbody\n")
	write(t, filepath.Join(dirs.Mappings, "checkout.yaml"), oneRule+"questions:\n  - id: Q-01\n    text: a question\n")
	write(t, filepath.Join(dirs.Ubiquitous, "cart.md"), "---\nname: Cart\n---\n\nwhat a cart is\n")
	head, err := Scan(dirs)
	if err != nil {
		t.Fatal(err)
	}

	got := make(map[uri.Kind]bool)
	for _, c := range Compare(&Snapshot{byURI: map[string]int{}}, head) {
		got[c.Kind] = true
	}
	for _, kind := range []uri.Kind{
		uri.KindOpportunity, uri.KindOpportunityCanvas, uri.KindStoryMap, uri.KindStory,
		uri.KindMapping, uri.KindRule, uri.KindExample, uri.KindQuestion, uri.KindTerm,
	} {
		if !got[kind] {
			t.Errorf("no change reported for %s", kind)
		}
	}
}

// livt:automates livt://mapping/automate-from-master-in-impl-repos/rule/R-09/example/EX-07
// A map is addressed by its key, so renaming it is a change to the one map —
// its name line reworded — rather than one map withdrawn and another added.
func TestCompareReadsARenamedStoryMapAsAChangeToIt(t *testing.T) {
	storyMap := func(name string) *Snapshot {
		t.Helper()
		dirs := repoDirs(t.TempDir())
		write(t, filepath.Join(dirs.USM, "discovery.yaml"), "name: "+name+"\nactivities: []\n")
		snapshot, err := Scan(dirs)
		if err != nil {
			t.Fatal(err)
		}
		return snapshot
	}

	change := only(t, Compare(storyMap("旧い名前"), storyMap("新しい名前")))
	if change.URI != uri.StoryMap("discovery") || change.Status != StatusModified {
		t.Fatalf("got %s %q, want the map modified", change.Status, change.URI)
	}
	if got, want := lineTexts(change.Lines), []string{"-旧い名前", "+新しい名前"}; !equal(got, want) {
		t.Errorf("lines = %v, want %v", got, want)
	}
}

// A canvas stands without an opportunity file, so the diff reads one too: after
// the opportunities, titled by its key. One that shares an opportunity's key is
// read right after it and called by the opportunity's name, as its sheet is.
func TestScanReadsACanvasWithNoOpportunityFile(t *testing.T) {
	dirs := repoDirs(t.TempDir())
	write(t, filepath.Join(dirs.Opportunities, "growth.md"), "---\nname: Growth\n---\n\nwhy it matters\n")
	write(t, filepath.Join(dirs.Canvases, "growth.yaml"), "canvas:\n  problems:\n    - a problem\n")
	write(t, filepath.Join(dirs.Canvases, "unframed.yaml"), "canvas:\n  problems:\n    - another problem\n")
	head, err := Scan(dirs)
	if err != nil {
		t.Fatal(err)
	}

	changes := Compare(&Snapshot{byURI: map[string]int{}}, head)
	want := []string{uri.Opportunity("growth"), uri.OpportunityCanvas("growth"), uri.OpportunityCanvas("unframed")}
	if got := uris(changes); !equal(got, want) {
		t.Fatalf("changes = %v, want %v", got, want)
	}
	if got := changes[1].Title; got != "Growth" {
		t.Errorf("paired canvas title = %q, want its opportunity's name", got)
	}
	if got := changes[2].Title; got != "unframed" {
		t.Errorf("unpaired canvas title = %q, want its key", got)
	}
}

// filedOutOfIDOrder lists a mapping's rules and examples somewhere other than
// where their IDs put them, which is what an author does to keep a new example
// beside the one it speaks to rather than at the end of the rule.
const filedOutOfIDOrder = `rules:
  - id: R-02
    name: second
    examples:
      - id: EX-01
        name: an example of the second
  - id: R-01
    name: first
    examples:
      - id: EX-01
        name: an example
`

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-02/example/EX-06
// The file's order is where an author put things; the ID is what the item is
// called. Read off the file, a retirement lands wherever its replacement was
// inserted, and the pair a review came to read side by side is split.
func TestCompareOrdersAMappingsChangesByTheirIDs(t *testing.T) {
	base := mapping(t, "checkout", filedOutOfIDOrder)
	head := mapping(t, "checkout", strings.NewReplacer(
		"an example of the second", "a reworded example of the second",
		`      - id: EX-01
        name: an example
`, `      - id: EX-02
        name: a replacement
      - id: EX-01
        name: an example
        retired: true
`).Replace(filedOutOfIDOrder))

	want := []string{
		"livt://mapping/checkout/rule/R-01/example/EX-01",
		"livt://mapping/checkout/rule/R-01/example/EX-02",
		"livt://mapping/checkout/rule/R-02/example/EX-01",
	}
	if got := uris(Compare(base, head)); !equal(got, want) {
		t.Errorf("changes = %v, want %v", got, want)
	}
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-03/example/EX-02
// EX-03: what changed is a removal beside its replacement, and what did not
// is still there to read it against.
func TestCompareKeepsUnchangedLinesAsContext(t *testing.T) {
	base := mapping(t, "checkout", oneRule)
	head := mapping(t, "checkout", strings.Replace(oneRule, "name: first", "name: second", 1))

	change := only(t, Compare(base, head))
	want := []string{"-first", "+second", " " + status("accepted")}
	if got := lineTexts(change.Lines); !equal(got, want) {
		t.Errorf("lines = %v, want %v", got, want)
	}
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-04/example/EX-01
// The board drops a retired rule, which is exactly why the diff must not —
// and the record keeps it, which is why the URI goes on resolving. The
// reading is a withdrawal; the record is a modification, and both are true at
// once.
func TestCompareReadsARetiredRuleAsAWithdrawalTheRecordStillHolds(t *testing.T) {
	base := mapping(t, "checkout", oneRule)
	head := mapping(t, "checkout", strings.Replace(oneRule, "name: first", "name: first\n    status: retired", 1))

	change := findURI(t, Compare(base, head), "livt://mapping/checkout/rule/R-01")
	if change.Became != BecameWithdrawn {
		t.Errorf("became %q, want %q", change.Became, BecameWithdrawn)
	}
	if change.Status != StatusModified {
		t.Errorf("status = %q, want %q — the entry is still on file, so its id is not free", change.Status, StatusModified)
	}
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-03/example/EX-07
// EX-08: a withdrawal is drawn as the removal it is. Retiring *adds* the line
// that retires, so drawn off the record it would read as the opposite of what
// happened — and the line that records it says nothing the badge has not.
func TestCompareDrawsAWithdrawalAsARemoval(t *testing.T) {
	base := mapping(t, "checkout", oneRule)
	head := mapping(t, "checkout", strings.Replace(oneRule, "name: first", "name: first\n    status: retired", 1))

	change := findURI(t, Compare(base, head), "livt://mapping/checkout/rule/R-01")
	if got := lineTexts(change.Lines); !equal(got, []string{"-first"}) {
		t.Errorf("lines = %v, want the statement removed and nothing else", got)
	}
}

// What replaced it is what a reader does want from a withdrawal, so the
// successor survives the pruning that drops the rest of the bookkeeping.
func TestCompareKeepsTheSuccessorOfAWithdrawnItem(t *testing.T) {
	const successor = "livt://mapping/checkout/rule/R-02"
	base := mapping(t, "checkout", oneRule)
	head := mapping(t, "checkout", strings.Replace(oneRule, "name: first",
		"name: first\n    status: retired\n    superseded_by:\n      - "+successor, 1))

	change := findURI(t, Compare(base, head), "livt://mapping/checkout/rule/R-01")
	want := []string{"-first", " " + LabelSupersededBy + ": " + successor}
	if got := lineTexts(change.Lines); !equal(got, want) {
		t.Errorf("lines = %v, want %v", got, want)
	}
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-04/example/EX-02
// A proposal being agreed is the one-line change it is on the file, and reads
// as one here. The status is always spelled, so accepting it is not a line
// arriving out of nowhere.
func TestCompareReadsAnAgreedProposalAsAOneLineChange(t *testing.T) {
	proposed := strings.Replace(oneRule, "name: first", "name: first\n    status: proposed", 1)
	base := mapping(t, "checkout", proposed)
	head := mapping(t, "checkout", oneRule)

	change := findURI(t, Compare(base, head), "livt://mapping/checkout/rule/R-01")
	want := []string{" first", "-" + status("proposed"), "+" + status("accepted")}
	if got := lineTexts(change.Lines); !equal(got, want) {
		t.Errorf("lines = %v, want %v", got, want)
	}
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-04/example/EX-03
// A retirement that happened before the base is not news.
func TestCompareLeavesOutAnItemRetiredInBothRevisions(t *testing.T) {
	retired := oneRule + `questions:
  - id: Q-01
    text: settled long ago
    retired: true
`
	base := mapping(t, "checkout", retired)
	head := mapping(t, "checkout", retired)

	if changes := Compare(base, head); len(changes) != 0 {
		t.Errorf("got %v, want nothing — the retirement is in neither revision's news", uris(changes))
	}
}

func findURI(t *testing.T, changes []Change, want string) Change {
	t.Helper()
	for _, c := range changes {
		if c.URI == want {
			return c
		}
	}
	t.Fatalf("no change for %s; got %v", want, uris(changes))
	return Change{}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestParseRangeReadsBothForms(t *testing.T) {
	for _, tc := range []struct {
		arg  string
		want Range
	}{
		{"main..HEAD", Range{Base: "main", Head: "HEAD"}},
		// livt:automates livt://mapping/review-diff-between-revisions/rule/R-01/example/EX-03
		// One revision is the working tree's base.
		{"main", Range{Base: "main"}},
		{"abc123..", Range{Base: "abc123"}},
	} {
		got, err := ParseRange(tc.arg)
		if err != nil {
			t.Errorf("ParseRange(%q) errored: %v", tc.arg, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ParseRange(%q) = %+v, want %+v", tc.arg, got, tc.want)
		}
	}
}

// A three-dot range is git's merge-base spelling and means something livt does
// not do. Cut would read it as a head named ".HEAD" — a ref git itself forbids
// — and the diff would fail deep inside rev-parse instead of at the flag.
func TestParseRangeRefusesWhatItDoesNotMean(t *testing.T) {
	for _, arg := range []string{"", "main...HEAD", "..HEAD"} {
		if got, err := ParseRange(arg); err == nil {
			t.Errorf("ParseRange(%q) = %+v, want an error", arg, got)
		}
	}
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-06/example/EX-01
// A rule is its text, not the value of a `name` key. Spelling the key would
// make the page a prettier YAML diff rather than a different reading of one.
func TestCompareGivesTheItemsOwnTextNoLabel(t *testing.T) {
	base := mapping(t, "checkout", oneRule)
	head := mapping(t, "checkout", strings.Replace(oneRule, "name: first", "name: second", 1))

	change := only(t, Compare(base, head))
	if label := change.Lines[0].Field.Label; label != "" {
		t.Errorf("the rule's own text carries the label %q, want none", label)
	}
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-06/example/EX-03
// A reworded rule shows where it was reworded. Two long sentences side by
// side leave the reader to find the difference by eye, which for a one-phrase
// edit is the whole of the work.
func TestCompareMarksWhatMovedInsideARewordedLine(t *testing.T) {
	const before = "注文を保存されたカードで確定できる"
	const after = "注文を保存されたカードで即時に確定できる"
	base := mapping(t, "checkout", strings.Replace(oneRule, "name: first", "name: "+before, 1))
	head := mapping(t, "checkout", strings.Replace(oneRule, "name: first", "name: "+after, 1))

	change := only(t, Compare(base, head))
	if got := changedParts(change.Lines[1]); got != "即時に" {
		t.Errorf("the added line marks %q as changed, want the phrase that was inserted", got)
	}
	if got := changedParts(change.Lines[0]); got != "" {
		t.Errorf("the removed line marks %q as changed, want nothing — no words were taken out", got)
	}
}

// Two different sentences are not one sentence edited. Marking up the few
// characters they happen to share would scatter highlights through both and
// say they are related when they are not.
func TestCompareLeavesAWhollyDifferentLineWhole(t *testing.T) {
	base := mapping(t, "checkout", strings.Replace(oneRule, "name: first", "name: 注文を確定できる", 1))
	head := mapping(t, "checkout", strings.Replace(oneRule, "name: first", "name: 在庫を引き当てる", 1))

	change := only(t, Compare(base, head))
	for _, l := range change.Lines[:2] {
		if l.Parts != nil {
			t.Errorf("%q was broken down against an unrelated line", l.Field.Value)
		}
	}
}

func changedParts(l Line) string {
	var changed strings.Builder
	for _, p := range l.Parts {
		if p.Changed {
			changed.WriteString(p.Text)
		}
	}
	return changed.String()
}

// livt:automates livt://mapping/name-example-mapping-itself/rule/R-03
// livt:automates livt://mapping/name-example-mapping-itself/rule/R-03/example/EX-01
// A renamed mapping is a change to the mapping itself, read as the rewording it
// is and headed by the name it has now.
func TestCompareReadsARenamedMappingAsAChangeToIt(t *testing.T) {
	base := mapping(t, "fix-login", "name: 旧い名前\n"+oneRule)
	head := mapping(t, "fix-login", "name: 新しい名前\n"+oneRule)

	change := only(t, Compare(base, head))
	if change.URI != uri.Mapping("fix-login") || change.Status != StatusModified {
		t.Fatalf("got %s %q, want the mapping modified", change.Status, change.URI)
	}
	if got, want := lineTexts(change.Lines), []string{"-旧い名前", "+新しい名前"}; !equal(got, want) {
		t.Errorf("lines = %v, want %v", got, want)
	}
	if change.Title != "新しい名前" {
		t.Errorf("title = %q, want the name it has now", change.Title)
	}
}

// scanned writes one example mapping beside the reports implementation
// repositories sent, and snapshots the repository holding them all.
func scanned(t *testing.T, key, yaml string, reports ...automation.Report) *Snapshot {
	t.Helper()
	dirs := repoDirs(t.TempDir())
	write(t, filepath.Join(dirs.Mappings, key+".yaml"), yaml)
	for _, r := range reports {
		writeReport(t, dirs.Automations, r)
	}
	snapshot, err := Scan(dirs)
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}

// writeReport lays a report where the collect station commits it: one file
// per implementation repository, named after it.
func writeReport(t *testing.T, dir string, r automation.Report) {
	t.Helper()
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(dir, r.Repo+".json"), string(data))
}

// report is one scan of an implementation repository, each citation's URL
// pinned to the revision scanned, as `livt automations` pins it.
func report(repo, rev string, cites ...automation.Citation) automation.Report {
	pinned := make([]automation.Citation, 0, len(cites))
	for _, c := range cites {
		c.URL = "https://github.com/" + repo + "/blob/" + rev + "/" + c.File + "#L" + strconv.Itoa(c.Line)
		pinned = append(pinned, c)
	}
	return automation.Report{Repo: repo, Rev: rev, Citations: pinned}
}

func cite(u, file string, line int) automation.Citation {
	return automation.Citation{URI: u, File: file, Line: line}
}

func automatedIn(repo string) string {
	return LabelAutomated + ": " + repo
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-07/example/EX-01
// A rule and an example a test starts citing each gain the repository that
// now automates them, as a line of their own. Neither answers for the other,
// the same way the board answers them.
func TestCompareReadsANewCitationAsTheRepositoryNowAutomatingIt(t *testing.T) {
	const rule = "livt://mapping/checkout/rule/R-01"
	const example = "livt://mapping/checkout/rule/R-01/example/EX-01"
	base := scanned(t, "checkout", oneRule, report("acme/api", "aaaaaaa"))
	head := scanned(t, "checkout", oneRule, report("acme/api", "bbbbbbb",
		cite(rule, "checkout_test.go", 10),
		cite(example, "checkout_test.go", 20),
	))

	changes := Compare(base, head)
	if got, want := uris(changes), []string{rule, example}; !equal(got, want) {
		t.Fatalf("changes = %v, want %v", got, want)
	}
	for _, c := range changes {
		if c.Became != BecameChanged {
			t.Errorf("%s became %q, want %q", c.URI, c.Became, BecameChanged)
		}
	}
	if got, want := lineTexts(changes[0].Lines), []string{" first", " " + status("accepted"), "+" + automatedIn("acme/api")}; !equal(got, want) {
		t.Errorf("rule lines = %v, want %v", got, want)
	}
	if got, want := lineTexts(changes[1].Lines), []string{" an example", "+" + automatedIn("acme/api")}; !equal(got, want) {
		t.Errorf("example lines = %v, want %v", got, want)
	}
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-07/example/EX-02
// A repository whose tests all stopped citing the rule is the line taken away.
// One still citing it stays as context, so the reader sees where the rule is
// still automated as well as where it no longer is.
func TestCompareReadsARepositoryThatStoppedCitingAsItsLineRemoved(t *testing.T) {
	const rule = "livt://mapping/checkout/rule/R-01"
	base := scanned(t, "checkout", oneRule,
		report("acme/api", "aaaaaaa", cite(rule, "checkout_test.go", 10)),
		report("acme/web", "ccccccc", cite(rule, "checkout.spec.ts", 5)),
	)
	head := scanned(t, "checkout", oneRule,
		report("acme/api", "bbbbbbb"),
		report("acme/web", "ccccccc", cite(rule, "checkout.spec.ts", 5)),
	)

	change := only(t, Compare(base, head))
	want := []string{" first", " " + status("accepted"), "-" + automatedIn("acme/api"), " " + automatedIn("acme/web")}
	if got := lineTexts(change.Lines); !equal(got, want) {
		t.Errorf("lines = %v, want %v", got, want)
	}
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-07/example/EX-03
// Every scan names a new revision and time and pins every URL to that
// revision, so the report file is rewritten top to bottom each time. None of
// that is a claim about what the team decided.
func TestCompareLeavesOutWhatMovesOnEveryScan(t *testing.T) {
	const rule = "livt://mapping/checkout/rule/R-01"
	before := report("acme/api", "aaaaaaa", cite(rule, "checkout_test.go", 10))
	before.GeneratedAt = time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	after := report("acme/api", "bbbbbbb", cite(rule, "checkout_test.go", 10))
	after.GeneratedAt = time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)

	if changes := Compare(scanned(t, "checkout", oneRule, before), scanned(t, "checkout", oneRule, after)); len(changes) != 0 {
		t.Errorf("got %v, want nothing — only the scan itself moved", uris(changes))
	}
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-07/example/EX-04
// The repository is the claim. A test that moved, or a second one joining it,
// leaves the rule automated where it was — and a report's pull request would
// not open over it either.
func TestCompareLeavesOutATestThatMovedOrWasJoinedByAnother(t *testing.T) {
	const rule = "livt://mapping/checkout/rule/R-01"
	base := scanned(t, "checkout", oneRule, report("acme/api", "aaaaaaa", cite(rule, "checkout_test.go", 10)))
	head := scanned(t, "checkout", oneRule, report("acme/api", "bbbbbbb",
		cite(rule, "order_test.go", 3),
		cite(rule, "order_test.go", 40),
	))

	if changes := Compare(base, head); len(changes) != 0 {
		t.Errorf("got %v, want nothing — acme/api automates the rule in both", uris(changes))
	}
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-04/example/EX-03
// A rule retired before the range is not this diff's news, whatever its tests
// do: the board counts it nowhere, so a test citing it or ceasing to changes
// nothing a reader can see.
func TestCompareLeavesOutTheCitationsOfAnItemRetiredInBothRevisions(t *testing.T) {
	const rule = "livt://mapping/checkout/rule/R-01"
	retired := strings.Replace(oneRule, "name: first", "name: first\n    status: retired", 1)
	base := scanned(t, "checkout", retired, report("acme/api", "aaaaaaa", cite(rule, "checkout_test.go", 10)))
	head := scanned(t, "checkout", retired, report("acme/api", "bbbbbbb"))

	if changes := Compare(base, head); len(changes) != 0 {
		t.Errorf("got %v, want nothing — the rule was retired in both revisions", uris(changes))
	}
}
