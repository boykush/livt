package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/boykush/livt/internal/domain"
)

func TestResolveTermCardLinksResolvedTerm(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "story-map.md"), []byte("---\nname: Story Map\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	b := Builder{UbiquitousDir: dir}
	card := b.resolveTermCard("story-map")

	if card.Name != "Story Map" {
		t.Fatalf("got name %q, want Story Map", card.Name)
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

// livt://mapping/scope-terms-by-context/rule/R-02/example/EX-03: a board writes
// "{ctx}/{term-key}" to reach a term scoped to one context, and the sticky keeps
// the context so it says which meaning is being reached for.
func TestResolveTermCardResolvesScopedReference(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "billing"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "billing", "invoice.md"), []byte("---\nname: 請求書\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	b := Builder{UbiquitousDir: dir}
	card := b.resolveTermCard("billing/invoice")

	if card.Name != "請求書" {
		t.Fatalf("got name %q, want 請求書", card.Name)
	}
	if card.Ctx != "billing" {
		t.Fatalf("got ctx %q, want billing", card.Ctx)
	}
	if card.Href != "../ubiquitous.html#billing/invoice" {
		t.Fatalf("got href %q", card.Href)
	}
}

// A bare key names the context-free term, so it must not fall back to a scoped
// one of the same key — that would hand a board the wrong meaning silently.
func TestResolveTermCardBareKeyDoesNotFallBackToAScopedTerm(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "billing"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "billing", "invoice.md"), []byte("---\nname: 請求書\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	b := Builder{UbiquitousDir: dir}
	card := b.resolveTermCard("invoice")

	if card.Href != "" {
		t.Fatalf("a bare key resolved to the scoped term at %q", card.Href)
	}
	if card.Ctx != "" {
		t.Fatalf("got ctx %q, want none for a bare reference", card.Ctx)
	}
}

// A reference arrives from a file the build does not control, so one that walks
// out of the ubiquitous directory degrades to a plain card rather than reading
// whatever it points at.
func TestResolveTermCardRejectsReferenceEscapingTheDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(filepath.Dir(dir), "secret.md"), []byte("---\nname: secret\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	b := Builder{UbiquitousDir: dir}
	card := b.resolveTermCard("../secret")

	if card.Href != "" || card.Name != "../secret" {
		t.Fatalf("got (%q, %q), want the reference rendered as a plain card", card.Name, card.Href)
	}
}

func TestStoryMapViewResolvesReferencedTerms(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "story-map.md"), []byte("---\nname: Story Map\n---\n"), 0o644); err != nil {
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
