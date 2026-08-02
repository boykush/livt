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

// writeGlossaryTerm creates ubiquitous/{ctx}/{key}.md, or ubiquitous/{key}.md
// when ctx is empty.
func writeGlossaryTerm(t *testing.T, dir, ctx, key, name string) {
	t.Helper()
	if ctx != "" {
		if err := os.MkdirAll(filepath.Join(dir, ctx), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	body := []byte("---\nname: " + name + "\n---\n\nDefinition of " + name + ".\n")
	if err := os.WriteFile(filepath.Join(dir, ctx, key+".md"), body, 0o644); err != nil {
		t.Fatal(err)
	}
}

func buildGlossaryHTML(t *testing.T, ubiquitousDir string) string {
	t.Helper()
	outDir := t.TempDir()
	b := Builder{UbiquitousDir: ubiquitousDir, OutDir: outDir}
	if err := b.buildGlossary(); err != nil {
		t.Fatal(err)
	}
	out, err := os.ReadFile(filepath.Join(outDir, "ubiquitous.html"))
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

// livt://mapping/scope-terms-by-context/rule/R-03/example/EX-01 and EX-02: the
// table shows a scoped term's context and leaves a context-free one without one.
func TestBuildGlossaryShowsContextOnScopedTermsOnly(t *testing.T) {
	dir := t.TempDir()
	writeGlossaryTerm(t, dir, "", "story", "ストーリー")
	writeGlossaryTerm(t, dir, "billing", "invoice", "請求書")
	html := buildGlossaryHTML(t, dir)

	if !strings.Contains(html, `id="billing/invoice"`) {
		t.Fatal("expected the scoped term's row to be anchored by its full reference")
	}
	if !strings.Contains(html, `id="story"`) {
		t.Fatal("expected the context-free term to keep its bare-key anchor")
	}
	if !strings.Contains(html, `data-filter-values="[&#34;billing&#34;]"`) {
		t.Fatal("expected the scoped row to carry its context as the filter hook")
	}
	if !strings.Contains(html, `data-filter-values="[]"`) {
		t.Fatal("expected the context-free row to match no context filter")
	}
}

// livt://mapping/scope-terms-by-context/rule/R-03/example/EX-03: one key under
// two contexts gives two rows, told apart by their anchors and their contexts.
func TestBuildGlossaryKeepsTheSameKeyUnderTwoContextsApart(t *testing.T) {
	dir := t.TempDir()
	writeGlossaryTerm(t, dir, "billing", "invoice", "請求の請求書")
	writeGlossaryTerm(t, dir, "shipping", "invoice", "出荷の請求書")
	html := buildGlossaryHTML(t, dir)

	for _, want := range []string{`id="billing/invoice"`, `id="shipping/invoice"`, "請求の請求書", "出荷の請求書"} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected both rows to stand on their own, missing %q", want)
		}
	}
}

// livt://mapping/scope-terms-by-context/rule/R-04/example/EX-01 and EX-03: the
// table offers one chip per context and mirrors the selection in ?context=.
func TestBuildGlossaryRendersContextFilterControls(t *testing.T) {
	dir := t.TempDir()
	writeGlossaryTerm(t, dir, "", "story", "ストーリー")
	writeGlossaryTerm(t, dir, "billing", "invoice", "請求書")
	html := buildGlossaryHTML(t, dir)

	if !strings.Contains(html, `data-filter-param="context"`) {
		t.Fatal("expected the glossary filter to ride in the context query param")
	}
	if !strings.Contains(html, `data-filter-value="billing"`) {
		t.Fatal("expected a chip for the billing context")
	}
	if !strings.Contains(html, `data-filter-value=""`) || !strings.Contains(html, ">All<") {
		t.Fatal("expected an All chip that clears the filter")
	}
	for _, network := range []string{"fetch(", "XMLHttpRequest"} {
		if strings.Contains(html, network) {
			t.Fatalf("filter must stay client-side, found %q", network)
		}
	}
}

// A master that cuts no contexts gains no filter bar: every chip would be "All".
func TestBuildGlossaryWithoutContextsRendersNoFilterBar(t *testing.T) {
	dir := t.TempDir()
	writeGlossaryTerm(t, dir, "", "story", "ストーリー")
	html := buildGlossaryHTML(t, dir)

	// The bare attribute also appears in the script's selector, so match the
	// element by the param only the bar itself carries.
	if strings.Contains(html, `data-filter-param="context"`) {
		t.Fatal("expected no filter bar when the master cuts no contexts")
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
