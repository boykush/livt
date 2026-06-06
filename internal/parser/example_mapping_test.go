package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseExampleMappingReadsReferencedTerms(t *testing.T) {
	path := filepath.Join(t.TempDir(), "story.yaml")
	data := []byte("rules: []\nterms:\n  - story-map\n  - story\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	em, err := ParseExampleMapping(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(em.Terms) != 2 || em.Terms[0] != "story-map" || em.Terms[1] != "story" {
		t.Fatalf("got terms %v, want [story-map story]", em.Terms)
	}
}
