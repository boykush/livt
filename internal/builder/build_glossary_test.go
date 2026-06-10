package builder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildGlossaryRendersTermAsAnchoredRow(t *testing.T) {
	ubiquitousDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(ubiquitousDir, "story-map.md"), []byte("---\nname: Story Map\n---\n\nA board to overview the whole story.\n"), 0o644); err != nil {
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
	if !strings.Contains(html, "Story Map") {
		t.Fatal("expected term name to be rendered")
	}
	if !strings.Contains(html, "A board to overview the whole story.") {
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
