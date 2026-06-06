package builder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildGlossaryRendersTermAsAnchoredRow(t *testing.T) {
	ubiquitousDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(ubiquitousDir, "story-map.md"), []byte("---\nname: ストーリーマップ\n---\n\n全体ストーリーを俯瞰するボード。\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	outDir := t.TempDir()
	b := Builder{UbiquitousDir: ubiquitousDir, OutDir: outDir}
	if err := b.buildGlossary(); err != nil {
		t.Fatal(err)
	}

	out, err := os.ReadFile(filepath.Join(outDir, "ubiquitous.html"))
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)

	if !strings.Contains(html, `id="story-map"`) {
		t.Fatal("expected term card to carry an id anchor")
	}
	if !strings.Contains(html, "ストーリーマップ") {
		t.Fatal("expected term name to be rendered")
	}
	if !strings.Contains(html, "全体ストーリーを俯瞰するボード。") {
		t.Fatal("expected definition body to be rendered")
	}
}

func TestBuildGlossaryWithoutTermsRendersEmptyState(t *testing.T) {
	outDir := t.TempDir()
	b := Builder{UbiquitousDir: filepath.Join(t.TempDir(), "missing"), OutDir: outDir}
	if err := b.buildGlossary(); err != nil {
		t.Fatal(err)
	}

	out, err := os.ReadFile(filepath.Join(outDir, "ubiquitous.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "No terms found.") {
		t.Fatal("expected empty-state message when no terms exist")
	}
}
