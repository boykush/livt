package builder

import (
	gohtml "html"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/i18n"
)

const demoCanvasYAML = "canvas:\n" +
	"  solution-ideas:\n    - テキストで残す\n" +
	"  problems:\n    - 共通理解が失われる\n" +
	"  budget:\n    - 余暇で開発\n"

// The canvas is one sheet in three zones — the facts, the solution, the value —
// so every box has to land in the panel the method puts it in.
func TestRenderOpportunityCanvasLaysOutTheThreeZones(t *testing.T) {
	b := outDirsBuilder(t)
	writeFile(t, filepath.Join(b.OpportunitiesDir, "demo.md"), "---\nname: デモ\n---\n\n一文\n")
	writeFile(t, filepath.Join(b.CanvasesDir, "demo.yaml"), demoCanvasYAML)

	if err := b.buildOpportunityCanvas(&domain.Opportunity{Key: domain.OpportunityKey{Value: "demo"}, Name: "デモ"}); err != nil {
		t.Fatal(err)
	}
	html := readFile(t, filepath.Join(b.OutDir, "opportunity-canvas", "demo.html"))

	// All ten boxes are on the sheet, unanswered ones included: a blank box is
	// the record of a question the opportunity has not answered.
	// Headings are compared escaped: one of them is Patton's "Customers &
	// Users", and the template is right to render the ampersand as &amp;.
	en := i18n.Of(i18n.En)
	for _, box := range (&domain.OpportunityCanvas{}).Boxes() {
		name := en.Msg("canvas." + box.Key)
		if !strings.Contains(html, gohtml.EscapeString(name)) {
			t.Errorf("box %q is missing from the canvas", name)
		}
	}
	if got := strings.Count(html, en.Msg("canvas.empty-box")); got != 7 {
		t.Errorf("got %d unanswered boxes, want 7 (three of the ten are filled in)", got)
	}
	for _, panel := range []string{"bg-rose-50", "bg-sky-50", "bg-emerald-50"} {
		if !strings.Contains(html, panel) {
			t.Errorf("the %s zone panel is missing; the sheet is not laid out in three zones", panel)
		}
	}
	// Facts sit left of the solution, which sits left of the value: the sheet
	// reads back from the idea to the problem, then forward to the value.
	facts, solution, value := strings.Index(html, "bg-rose-50"), strings.Index(html, "bg-sky-50"), strings.Index(html, "bg-emerald-50")
	if facts >= solution || solution >= value {
		t.Errorf("zones are out of order: facts=%d solution=%d value=%d", facts, solution, value)
	}
}

// livt://opportunity/collaborative-discovery: the Related section is what says
// how far an opportunity has been taken. A canvas means it was thought through,
// a story map means it was taken on.
func TestRenderOpportunityLinksItsCanvasAndStoryMap(t *testing.T) {
	b := outDirsBuilder(t)
	o := &domain.Opportunity{Key: domain.OpportunityKey{Value: "demo"}, Name: "デモ", Body: "一文で言う"}
	maps := []storyMapRef{{Name: "デモマップ", Path: "../story-map/デモマップ.html"}}

	out := filepath.Join(b.OutDir, "opportunity", "demo.html")
	if err := b.renderOpportunityPage(out, o, "../opportunity-canvas/demo.html", maps); err != nil {
		t.Fatal(err)
	}
	html := readFile(t, out)

	for _, want := range []string{"Opportunity Canvas", "デモマップ", "一文で言う"} {
		if !strings.Contains(html, want) {
			t.Errorf("expected %q on the opportunity page", want)
		}
	}
}

// An opportunity nobody has held a canvas session for links to none, so the
// absence stays visible instead of resolving to an empty board.
func TestRenderOpportunityWithoutCanvasLinksNone(t *testing.T) {
	b := outDirsBuilder(t)
	o := &domain.Opportunity{Key: domain.OpportunityKey{Value: "demo"}, Name: "デモ"}

	out := filepath.Join(b.OutDir, "opportunity", "demo.html")
	if err := b.renderOpportunityPage(out, o, "", nil); err != nil {
		t.Fatal(err)
	}
	if html := readFile(t, out); strings.Contains(html, "opportunity-canvas/") {
		t.Fatal("expected no canvas link for an opportunity with no canvas")
	}
}

// The filename is the join: a map whose key names a committed opportunity earns
// its stories a chip pointing at that opportunity's page, not at the map.
func TestStoryChipsPointAtTheOpportunityTheMapServes(t *testing.T) {
	b := outDirsBuilder(t)
	writeFile(t, filepath.Join(b.OpportunitiesDir, "demo.md"), "---\nname: デモ機会\n---\n\n一文\n")
	writeFile(t, filepath.Join(b.USMDir, "demo.yaml"),
		"name: デモマップ\nactivities:\n  - key: a\n    name: A\n    steps:\n      - key: s\n        name: S\n        stories:\n          - key: card\n            name: カード\n")

	opportunities, err := b.opportunityIndex()
	if err != nil {
		t.Fatal(err)
	}
	built, err := b.buildStoryMaps(opportunities)
	if err != nil {
		t.Fatal(err)
	}

	refs := built.StoryOpportunities["card"]
	if len(refs) != 1 {
		t.Fatalf("got %d opportunity refs, want 1", len(refs))
	}
	if refs[0].Name != "デモ機会" {
		t.Errorf("chip reads %q, want the opportunity's name rather than the map's", refs[0].Name)
	}
	if refs[0].Path != "../opportunity/demo.html" {
		t.Errorf("chip links to %q, want the opportunity's page", refs[0].Path)
	}
	if got := built.MapsByOpportunity["demo"]; len(got) != 1 || got[0].Name != "デモマップ" {
		t.Errorf("MapsByOpportunity[demo] = %+v, want the map mapped for it", got)
	}
}

// A livt repository with no opportunities/ at all keeps working exactly as it
// did before opportunities became files: the map stands in as its own
// opportunity, named by the map. Nothing has to be migrated to keep building.
func TestMapWithNoOpportunityFileStandsInAsItsOwn(t *testing.T) {
	b := outDirsBuilder(t)
	writeFile(t, filepath.Join(b.USMDir, "demo.yaml"),
		"name: デモマップ\nactivities:\n  - key: a\n    name: A\n    steps:\n      - key: s\n        name: S\n        stories:\n          - key: card\n            name: カード\n")

	opportunities, err := b.opportunityIndex()
	if err != nil {
		t.Fatal(err)
	}
	built, err := b.buildStoryMaps(opportunities)
	if err != nil {
		t.Fatal(err)
	}

	refs := built.StoryOpportunities["card"]
	if len(refs) != 1 || refs[0].Name != "デモマップ" || refs[0].Path != "../story-map/デモマップ.html" {
		t.Fatalf("got %+v, want the map standing in as its own opportunity", refs)
	}
	if len(built.MapsByOpportunity) != 0 {
		t.Errorf("MapsByOpportunity = %+v, want empty: no opportunity file claims this map", built.MapsByOpportunity)
	}
	// With nothing to link to, the board does not offer an opportunity link.
	if len(built.Tiles) != 1 || built.Tiles[0].Opportunity != nil {
		t.Errorf("tile = %+v, want no opportunity chip", built.Tiles)
	}
}

// A canvas whose opportunity was renamed or deleted must not outlive it, the
// same guarantee TestBuildDropsPagesForRemovedResources makes for the rest.
func TestBuildDropsPagesForRemovedOpportunities(t *testing.T) {
	b := emptyDirsBuilder(t)
	writeFile(t, filepath.Join(b.OpportunitiesDir, "kept.md"), "---\nname: 残る\n---\n\n一文\n")
	writeFile(t, filepath.Join(b.CanvasesDir, "kept.yaml"), demoCanvasYAML)
	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	orphans := []string{
		filepath.Join(b.OutDir, "opportunity", "removed.html"),
		filepath.Join(b.OutDir, "opportunity-canvas", "removed.html"),
	}
	for _, p := range orphans {
		writeFile(t, p, "<html>stale</html>")
	}
	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	for _, p := range orphans {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Errorf("%s survived the rebuild (stat error: %v)", p, err)
		}
	}
	for _, p := range []string{
		filepath.Join(b.OutDir, "opportunity", "kept.html"),
		filepath.Join(b.OutDir, "opportunity-canvas", "kept.html"),
	} {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("rebuild dropped a live page: %v", err)
		}
	}
}

// An opportunity file need not carry a name, and every surface naming one is a
// heading, a filter chip, or a link label. It falls back to the key rather than
// rendering blank — the same call Story.DisplayName makes.
func TestNamelessOpportunityFallsBackToItsKey(t *testing.T) {
	b := outDirsBuilder(t)
	writeFile(t, filepath.Join(b.OpportunitiesDir, "demo.md"), "一文だけの本文\n")
	writeFile(t, filepath.Join(b.USMDir, "demo.yaml"),
		"name: デモマップ\nactivities:\n  - key: a\n    name: A\n    steps:\n      - key: s\n        name: S\n        stories:\n          - key: card\n            name: カード\n")

	opportunities, err := b.opportunityIndex()
	if err != nil {
		t.Fatal(err)
	}
	built, err := b.buildStoryMaps(opportunities)
	if err != nil {
		t.Fatal(err)
	}
	// A blank chip would be an unpickable filter axis, not just an ugly one.
	if refs := built.StoryOpportunities["card"]; len(refs) != 1 || refs[0].Name != "demo" {
		t.Fatalf("chip = %+v, want it to fall back to the key", refs)
	}

	if _, err := b.buildOpportunities(built.MapsByOpportunity); err != nil {
		t.Fatal(err)
	}
	if html := readFile(t, filepath.Join(b.OutDir, "opportunity", "demo.html")); !strings.Contains(html, ">demo</h1>") {
		t.Error("the opportunity page heading is blank")
	}
}

// outDirsBuilder is emptyDirsBuilder with the per-resource output directories
// laid out, for a test that drives one build step rather than a whole Build.
func outDirsBuilder(t *testing.T) Builder {
	t.Helper()
	b := emptyDirsBuilder(t)
	if err := b.resetGeneratedDirs(); err != nil {
		t.Fatal(err)
	}
	return b
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	out, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}
