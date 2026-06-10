package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTermReadsNameAndBody(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "story-map.md")
	data := []byte("---\nname: Story Map\n---\n\nA board to overview activities, steps, and stories.\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	term, err := ParseTerm(path)
	if err != nil {
		t.Fatal(err)
	}

	if term.Key != "story-map" {
		t.Fatalf("got key %q, want story-map", term.Key)
	}
	if term.Name != "Story Map" {
		t.Fatalf("got name %q", term.Name)
	}
	if term.Body != "A board to overview activities, steps, and stories." {
		t.Fatalf("got body %q", term.Body)
	}
}

func TestParseAllTermsReturnsEveryFile(t *testing.T) {
	dir := t.TempDir()
	for _, key := range []string{"story", "discovery"} {
		data := []byte("---\nname: " + key + "\n---\n")
		if err := os.WriteFile(filepath.Join(dir, key+".md"), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	terms, err := ParseAllTerms(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(terms) != 2 {
		t.Fatalf("got %d terms, want 2", len(terms))
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
