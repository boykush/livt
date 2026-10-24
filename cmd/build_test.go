package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// livt:automates livt://mapping/define-policies-apart-from-stories/rule/R-02/example/EX-01
// The site reads terms from design/ubiquitous, beside design/policies.
func TestBuildReadsTermsFromDesign(t *testing.T) {
	b, err := newBuilder(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join("design", "ubiquitous"); b.UbiquitousDir != want {
		t.Errorf("UbiquitousDir = %q, want %q", b.UbiquitousDir, want)
	}
}

// livt:automates livt://mapping/define-policies-apart-from-stories/rule/R-02/example/EX-02
// livt:automates livt://mapping/define-policies-apart-from-stories/rule/R-04/example/EX-03
// A term filed under design/ubiquitous keeps the URI it had under ubiquitous/,
// and the old directory is no longer read.
func TestTermsMovedUnderDesignKeepTheirURIs(t *testing.T) {
	root := t.TempDir()
	for _, dir := range []string{filepath.Join("design", "ubiquitous"), "ubiquitous"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "design", "ubiquitous", "story.md"), []byte("---\nname: ストーリー\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "ubiquitous", "stale.md"), []byte("---\nname: 古い置き場\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := resolveURI(&bytes.Buffer{}, root, "livt://ubiquitous/story", formatJSON, ""); err != nil {
		t.Errorf("livt://ubiquitous/story: %v, want it to resolve from design/ubiquitous", err)
	}
	if err := resolveURI(&bytes.Buffer{}, root, "livt://ubiquitous/stale", formatJSON, ""); err == nil {
		t.Error("livt://ubiquitous/stale resolved from the old ubiquitous/ directory")
	}
}
