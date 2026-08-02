package parser

import (
	"os"
	"path/filepath"
	"testing"
)

// writeTerm creates ubiquitous/{ctx}/{key}.md, or ubiquitous/{key}.md when ctx
// is empty.
func writeTerm(t *testing.T, dir, ctx, key, name string) {
	t.Helper()
	if ctx != "" {
		if err := os.MkdirAll(filepath.Join(dir, ctx), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	data := []byte("---\nname: " + name + "\n---\n\nDefinition of " + name + ".\n")
	if err := os.WriteFile(filepath.Join(dir, ctx, key+".md"), data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestParseTermReadsNameAndBody(t *testing.T) {
	dir := t.TempDir()
	data := []byte("---\nname: Story Map\n---\n\nA board to overview activities, steps, and stories.\n")
	if err := os.WriteFile(filepath.Join(dir, "story-map.md"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	term, err := ParseTerm(dir, "story-map")
	if err != nil {
		t.Fatal(err)
	}

	if term.Key != "story-map" {
		t.Fatalf("got key %q, want story-map", term.Key)
	}
	if term.Ctx != "" {
		t.Fatalf("got ctx %q, want a term at the root to hold across contexts", term.Ctx)
	}
	if term.Name != "Story Map" {
		t.Fatalf("got name %q", term.Name)
	}
	if term.Body != "A board to overview activities, steps, and stories." {
		t.Fatalf("got body %q", term.Body)
	}
}

// livt://mapping/scope-terms-by-context/rule/R-01/example/EX-01: the directory
// holding the file is the context, so a term under one reads as scoped to it.
func TestParseTermTakesContextFromItsDirectory(t *testing.T) {
	dir := t.TempDir()
	writeTerm(t, dir, "billing", "invoice", "請求書")

	term, err := ParseTerm(dir, "billing/invoice")
	if err != nil {
		t.Fatal(err)
	}
	if term.Ctx != "billing" || term.Key != "invoice" {
		t.Fatalf("got (%q, %q), want (billing, invoice)", term.Ctx, term.Key)
	}
	if term.Name != "請求書" {
		t.Fatalf("got name %q", term.Name)
	}
}

// livt://mapping/scope-terms-by-context/rule/R-01/example/EX-03: the same key at
// the root and under a context are two terms, and the filesystem is what keeps
// them apart. Resolving one must never hand back the other.
func TestSameKeyAtRootAndUnderContextAreDistinctTerms(t *testing.T) {
	dir := t.TempDir()
	writeTerm(t, dir, "", "invoice", "共通の請求書")
	writeTerm(t, dir, "billing", "invoice", "請求コンテキストの請求書")

	shared, err := ParseTerm(dir, "invoice")
	if err != nil {
		t.Fatal(err)
	}
	scoped, err := ParseTerm(dir, "billing/invoice")
	if err != nil {
		t.Fatal(err)
	}

	if shared.Name == scoped.Name {
		t.Fatalf("both resolved to %q; the context is part of what identifies a term", shared.Name)
	}
	if shared.Ctx != "" || scoped.Ctx != "billing" {
		t.Fatalf("got contexts (%q, %q), want (\"\", billing)", shared.Ctx, scoped.Ctx)
	}
}

// A reference is externally supplied — it arrives from a board, an MCP client,
// or a CLI argument — so it must not be able to walk out of the ubiquitous
// directory or nest a context inside another.
func TestParseTermRejectsReferencesThatEscapeTheDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "secret.md"), []byte("---\nname: secret\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, ref := range []string{"../secret", "billing/../secret", "a/b/c", "billing/", "/invoice", ""} {
		if _, err := ParseTerm(dir, ref); err == nil {
			t.Errorf("ParseTerm(%q) succeeded, want a rejected reference", ref)
		}
	}
}

func TestParseAllTermsReturnsEveryFile(t *testing.T) {
	dir := t.TempDir()
	for _, key := range []string{"story", "discovery"} {
		writeTerm(t, dir, "", key, key)
	}

	terms, err := ParseAllTerms(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(terms) != 2 {
		t.Fatalf("got %d terms, want 2", len(terms))
	}
}

// The context-free terms come first, then each context's, so a build renders the
// table in the order the filesystem lays it out rather than an incidental one.
func TestParseAllTermsWalksContextDirectories(t *testing.T) {
	dir := t.TempDir()
	writeTerm(t, dir, "", "story", "ストーリー")
	writeTerm(t, dir, "shipping", "invoice", "出荷の請求書")
	writeTerm(t, dir, "billing", "invoice", "請求の請求書")

	terms, err := ParseAllTerms(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(terms) != 3 {
		t.Fatalf("got %d terms, want 3", len(terms))
	}

	type ref struct{ ctx, key string }
	got := make([]ref, len(terms))
	for i, term := range terms {
		got[i] = ref{term.Ctx, term.Key}
	}
	want := []ref{{"", "story"}, {"billing", "invoice"}, {"shipping", "invoice"}}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("term %d is %v, want %v (full order %v)", i, got[i], want[i], got)
		}
	}
}

// Contexts are one path segment deep, the same depth a term URI carries, so a
// directory nested below one holds no terms the glossary can address.
func TestParseAllTermsIgnoresTermsBelowOneContextLevel(t *testing.T) {
	dir := t.TempDir()
	writeTerm(t, dir, filepath.Join("billing", "deeper"), "invoice", "深すぎる")

	terms, err := ParseAllTerms(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(terms) != 0 {
		t.Fatalf("got %d terms, want none — a context is one level deep", len(terms))
	}
}

func TestParseAllTermsMissingDirIsEmpty(t *testing.T) {
	terms, err := ParseAllTerms(filepath.Join(t.TempDir(), "does-not-exist"))
	if err != nil {
		t.Fatal(err)
	}
	if len(terms) != 0 {
		t.Fatalf("got %d terms, want 0", len(terms))
	}
}
