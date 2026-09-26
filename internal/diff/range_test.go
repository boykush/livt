package diff

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-01

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// relDirs names the input directories as the repository holds them. Compute
// reads them under the root it is given, which is what a real build does from
// the directory it runs in.
var relDirs = Dirs{
	Opportunities: "opportunities",
	Canvases:      filepath.Join("discoveries", "opportunity-canvases"),
	Mappings:      filepath.Join("discoveries", "example-mappings"),
	Stories:       "stories",
	USM:           filepath.Join("discoveries", "usm"),
	Ubiquitous:    "ubiquitous",
	Automations:   "automations",
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("git %v failed (%v): %s", args, err, out)
	}
}

// commitedRepo is a git repository holding one example mapping, committed.
func committedRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	root := t.TempDir()
	git(t, root, "init")
	git(t, root, "config", "user.email", "t@example.com")
	git(t, root, "config", "user.name", "t")
	write(t, filepath.Join(root, relDirs.Mappings, "checkout.yaml"), oneRule)
	git(t, root, "add", ".")
	git(t, root, "commit", "-m", "spec")
	return root
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-01/example/EX-03
// One revision compares against the working tree, which is the shape a review
// takes — you are on the branch that makes the change. The head has no hash
// to print because it is not a revision.
func TestComputeDiffsARevisionAgainstTheWorkingTree(t *testing.T) {
	root := committedRepo(t)
	write(t, filepath.Join(root, relDirs.Mappings, "checkout.yaml"), strings.Replace(oneRule, "name: first", "name: second", 1))

	result, err := Range{Base: "HEAD"}.Compute(root, relDirs)
	if err != nil {
		t.Fatal(err)
	}
	if result.Head != "" {
		t.Errorf("head = %q, want empty for the working tree", result.Head)
	}
	if result.Base == "" {
		t.Error("base = empty, want the short hash git resolved")
	}
	change := findURI(t, result.Changes, "livt://mapping/checkout/rule/R-01")
	if change.Became != BecameChanged || result.Changed != 1 {
		t.Errorf("got %s with %d changed, want one changed rule", change.Became, result.Changed)
	}
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-05/example/EX-01
// EX-02: a URI the site holds is linked to the page it lands on, and one it
// does not hold is left unlinked rather than pointed at a page nobody built.
func TestComputeLinksOnlyWhatTheSiteHolds(t *testing.T) {
	root := committedRepo(t)
	write(t, filepath.Join(root, relDirs.Mappings, "checkout.yaml"), `rules:
  - id: R-02
    name: replaced the first
`)

	result, err := Range{Base: "HEAD"}.Compute(root, relDirs)
	if err != nil {
		t.Fatal(err)
	}
	added := findURI(t, result.Changes, "livt://mapping/checkout/rule/R-02")
	if added.Page != "mapping/checkout.html#rule-R-02" {
		t.Errorf("added rule page = %q, want the sticky it lands on", added.Page)
	}
	removed := findURI(t, result.Changes, "livt://mapping/checkout/rule/R-01")
	if removed.Status != StatusRemoved || removed.Page != "" {
		t.Errorf("removed rule = %s with page %q, want removed and unlinked", removed.Status, removed.Page)
	}
}

// A revision from before a directory existed is a legitimate side of a diff:
// showing that it held none of a resource is the whole job, so the export must
// not fail on the directories it cannot find.
func TestComputeReadsARevisionMissingWholeDirectories(t *testing.T) {
	root := committedRepo(t)
	write(t, filepath.Join(root, relDirs.Ubiquitous, "cart.md"), "---\nname: Cart\n---\n\nwhat a cart is\n")

	result, err := Range{Base: "HEAD"}.Compute(root, relDirs)
	if err != nil {
		t.Fatal(err)
	}
	if term := findURI(t, result.Changes, "livt://ubiquitous/cart"); term.Status != StatusAdded {
		t.Errorf("term = %s, want added — the base revision had no ubiquitous directory", term.Status)
	}
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-01/example/EX-04
// Both ways a diff can be asked for something unreadable stop the build, and
// each says the thing that fixes it.
func TestComputeFailsWhenThereIsNothingToRead(t *testing.T) {
	root := committedRepo(t)
	for _, tc := range []struct {
		name  string
		root  string
		rev   string
		wants string
	}{
		{"outside a repository", t.TempDir(), "HEAD", "not a git repository"},
		{"revision git cannot find", root, "no-such-rev", `revision "no-such-rev" not found`},
	} {
		_, err := Range{Base: tc.rev}.Compute(tc.root, relDirs)
		if err == nil {
			t.Errorf("%s: no error, want one", tc.name)
			continue
		}
		if !strings.Contains(err.Error(), tc.wants) {
			t.Errorf("%s: error = %q, want it to mention %q", tc.name, err, tc.wants)
		}
	}
}

// livt:automates livt://mapping/review-diff-between-revisions/rule/R-07
// Each revision is read with the reports it held, so the collect station's
// pull request — a report committed and nothing else — reads as the rules it
// automates rather than as a file rewritten top to bottom.
func TestComputeReadsTheReportsEachRevisionHeld(t *testing.T) {
	const rule = "livt://mapping/checkout/rule/R-01"
	root := committedRepo(t)
	reports := filepath.Join(root, relDirs.Automations)
	writeReport(t, reports, report("acme/api", "aaaaaaa"))
	git(t, root, "add", ".")
	git(t, root, "commit", "-m", "collect")
	writeReport(t, reports, report("acme/api", "bbbbbbb", cite(rule, "checkout_test.go", 10)))
	git(t, root, "add", ".")
	git(t, root, "commit", "-m", "collect again")

	result, err := Range{Base: "HEAD~1", Head: "HEAD"}.Compute(root, relDirs)
	if err != nil {
		t.Fatal(err)
	}
	change := only(t, result.Changes)
	if change.URI != rule || result.Changed != 1 {
		t.Errorf("got %s with %d changed, want the rule the report now cites", change.URI, result.Changed)
	}
}
