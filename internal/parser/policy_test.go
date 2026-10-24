package parser

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/boykush/livt/internal/domain"
)

// writePolicy creates {dir}/{ctx}/{key}.md, or {dir}/{key}.md when ctx is
// empty.
func writePolicy(t *testing.T, dir, ctx, key, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, ctx), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ctx, key+".md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// livt:automates livt://mapping/define-policies-apart-from-stories/rule/R-01/example/EX-02
// A policy is markdown: frontmatter, then a body.
func TestParsePolicyReadsFrontmatterAndBody(t *testing.T) {
	dir := t.TempDir()
	writePolicy(t, dir, "", "closed-items-are-not-unfinished",
		"---\nname: 閉じたものは未完了に数えない\n---\n\n退役したルールは誰も待っていない。\n")

	p, err := ParsePolicy(dir, "closed-items-are-not-unfinished")
	if err != nil {
		t.Fatal(err)
	}
	if p.Key != "closed-items-are-not-unfinished" {
		t.Errorf("got key %q", p.Key)
	}
	if p.Name != "閉じたものは未完了に数えない" {
		t.Errorf("got name %q", p.Name)
	}
	if p.Body != "退役したルールは誰も待っていない。" {
		t.Errorf("got body %q", p.Body)
	}
}

// livt:automates livt://mapping/define-policies-apart-from-stories/rule/R-03/example/EX-01
// The directory holding the file is the context.
func TestParsePolicyTakesContextFromItsDirectory(t *testing.T) {
	dir := t.TempDir()
	writePolicy(t, dir, "design-system", "lists-filter-by-opportunity", "---\nname: 一覧はオポチュニティで絞り込める\n---\n")

	p, err := ParsePolicy(dir, "design-system/lists-filter-by-opportunity")
	if err != nil {
		t.Fatal(err)
	}
	if p.Ctx != "design-system" || p.Key != "lists-filter-by-opportunity" {
		t.Errorf("got ctx %q key %q, want design-system / lists-filter-by-opportunity", p.Ctx, p.Key)
	}
}

// livt:automates livt://mapping/define-policies-apart-from-stories/rule/R-03/example/EX-02
// A policy at the root of the directory has no context.
func TestParsePolicyAtTheRootHoldsAcrossContexts(t *testing.T) {
	dir := t.TempDir()
	writePolicy(t, dir, "", "closed-items-are-not-unfinished", "---\nname: 閉じたものは未完了に数えない\n---\n")

	p, err := ParsePolicy(dir, "closed-items-are-not-unfinished")
	if err != nil {
		t.Fatal(err)
	}
	if p.Ctx != "" {
		t.Errorf("got ctx %q, want none", p.Ctx)
	}
}

// livt:automates livt://mapping/define-policies-apart-from-stories/rule/R-05/example/EX-01
// The terms a policy takes and yields are read as references, a context-scoped
// one included.
func TestParsePolicyReadsItsInputAndOutputTerms(t *testing.T) {
	dir := t.TempDir()
	writePolicy(t, dir, "design-system", "lists-filter-by-opportunity",
		"---\nname: 一覧はオポチュニティで絞り込める\ninput:\n  - opportunity\noutput:\n  - design-system/list\n---\n")

	p, err := ParsePolicy(dir, "design-system/lists-filter-by-opportunity")
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(p.Input, []string{"opportunity"}) {
		t.Errorf("got input %v", p.Input)
	}
	if !slices.Equal(p.Output, []string{"design-system/list"}) {
		t.Errorf("got output %v", p.Output)
	}
}

// livt:automates livt://mapping/define-policies-apart-from-stories/rule/R-06/example/EX-01
// Each of a rule's statuses reads on a policy too.
func TestParsePolicyReadsEveryRuleStatus(t *testing.T) {
	for _, s := range domain.RuleStatuses {
		dir := t.TempDir()
		writePolicy(t, dir, "", "p", "---\nname: p\nstatus: "+string(s)+"\n---\n")

		p, err := ParsePolicy(dir, "p")
		if err != nil {
			t.Fatalf("status %s: %v", s, err)
		}
		if p.Status != s {
			t.Errorf("got status %q, want %q", p.Status, s)
		}
	}
}

// livt:automates livt://mapping/define-policies-apart-from-stories/rule/R-06/example/EX-02
func TestParsePolicyWithNoStatusIsAccepted(t *testing.T) {
	dir := t.TempDir()
	writePolicy(t, dir, "", "p", "---\nname: p\n---\n")

	p, err := ParsePolicy(dir, "p")
	if err != nil {
		t.Fatal(err)
	}
	if p.Status != domain.RuleAccepted {
		t.Errorf("got status %q, want accepted", p.Status)
	}
}

// A mistyped status is refused rather than read as accepted, which would put a
// withdrawn policy back in force.
func TestParsePolicyRefusesAnUnknownStatus(t *testing.T) {
	dir := t.TempDir()
	writePolicy(t, dir, "", "p", "---\nname: p\nstatus: retierd\n---\n")

	_, err := ParsePolicy(dir, "p")
	if err == nil || !strings.Contains(err.Error(), "retierd") {
		t.Errorf("got %v, want an error naming the unknown status", err)
	}
}

// livt:automates livt://mapping/define-policies-apart-from-stories/rule/R-07/example/EX-02
// A retired policy names what replaced it when something did, and reads as
// retired with nothing when nothing did.
func TestParsePolicyRetiredWithOrWithoutASuccessor(t *testing.T) {
	dir := t.TempDir()
	writePolicy(t, dir, "", "replaced", "---\nname: 置き換えられた\nstatus: retired\nsuperseded_by:\n  - livt://policy/successor\n---\n")
	writePolicy(t, dir, "", "dropped", "---\nname: やめた\nstatus: retired\n---\n")

	replaced, err := ParsePolicy(dir, "replaced")
	if err != nil {
		t.Fatal(err)
	}
	if replaced.Status != domain.RuleRetired || !slices.Equal(replaced.SupersededBy, []string{"livt://policy/successor"}) {
		t.Errorf("got %q %v, want retired pointing at livt://policy/successor", replaced.Status, replaced.SupersededBy)
	}

	dropped, err := ParsePolicy(dir, "dropped")
	if err != nil {
		t.Fatal(err)
	}
	if dropped.Status != domain.RuleRetired || len(dropped.SupersededBy) != 0 {
		t.Errorf("got %q %v, want retired with no successor", dropped.Status, dropped.SupersededBy)
	}
}

// A reference that would climb out of the policies directory is refused before
// anything is read.
func TestParsePolicyRefusesATraversingReference(t *testing.T) {
	if _, err := ParsePolicy(t.TempDir(), "../secret"); err == nil {
		t.Error("got nil, want the reference refused")
	}
}
