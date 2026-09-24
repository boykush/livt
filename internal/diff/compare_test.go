package diff

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

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

// Where the spec went is what a reader does want from a withdrawal, so the
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
