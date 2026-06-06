package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/boykush/livt/internal/domain"
)

func TestResolveTermCardLinksResolvedTerm(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "story-map.md"), []byte("---\nname: ストーリーマップ\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	b := Builder{UbiquitousDir: dir}
	card := b.resolveTermCard("story-map")

	if card.Name != "ストーリーマップ" {
		t.Fatalf("got name %q, want ストーリーマップ", card.Name)
	}
	if card.Href != "../ubiquitous.html#story-map" {
		t.Fatalf("got href %q", card.Href)
	}
}

func TestResolveTermCardUnresolvedRendersPlainCard(t *testing.T) {
	b := Builder{UbiquitousDir: t.TempDir()}
	card := b.resolveTermCard("missing")

	if card.Name != "missing" {
		t.Fatalf("got name %q, want the key as fallback", card.Name)
	}
	if card.Href != "" {
		t.Fatalf("expected no href for an unresolved term, got %q", card.Href)
	}
}

func TestStoryMapViewResolvesReferencedTerms(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "story-map.md"), []byte("---\nname: ストーリーマップ\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	b := Builder{UbiquitousDir: dir}
	view := b.toStoryMapView(&domain.StoryMap{Name: "discovery", Ubiquitous: []string{"story-map", "missing"}})

	if len(view.StoryMap.Ubiquitous) != 2 {
		t.Fatalf("got %d term cards, want 2", len(view.StoryMap.Ubiquitous))
	}
	if view.StoryMap.Ubiquitous[0].Href == "" {
		t.Fatal("expected the resolved term to link to the glossary")
	}
	if view.StoryMap.Ubiquitous[1].Href != "" {
		t.Fatal("expected the unresolved term to render without a link")
	}
}
