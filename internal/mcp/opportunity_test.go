package mcp

import (
	"testing"

	"github.com/boykush/livt/internal/uri"
)

// The listing is where a consuming agent starts: it should say not just which
// opportunities exist but how far each was taken — a canvas means it was
// thought through, a story map means it was taken on.
func TestListOpportunitiesHandsOutCanvasAndMapURIs(t *testing.T) {
	got, err := newTestServer(t).cfg.opportunities()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d opportunities, want 1", len(got))
	}
	o := got[0]
	if o.Key != "demo-map" || o.Name != "デモ機会" {
		t.Errorf("got key %q name %q", o.Key, o.Name)
	}
	if o.URI != "livt://opportunity/demo-map" {
		t.Errorf("uri = %q", o.URI)
	}
	if o.CanvasURI != "livt://opportunity-canvas/demo-map" {
		t.Errorf("canvas_uri = %q", o.CanvasURI)
	}
	if o.Statement == "" {
		t.Error("statement is empty; the opportunity's own sentence is what a caller reads first")
	}
	// The join is the filename, so the map is found by key while the ref carries
	// the display name the map is addressed by.
	if len(o.StoryMaps) != 1 || o.StoryMaps[0].Name != "デモマップ" {
		t.Fatalf("story_maps = %+v, want the map keyed demo-map", o.StoryMaps)
	}
}

// The canvas resolves as its ten boxes with the questions they ask, so an agent
// reading it cold knows what each box is for. Unanswered boxes stay.
func TestOpportunityCanvasResolvesEveryBoxWithItsPrompt(t *testing.T) {
	s := newTestServer(t)
	canvas, err := s.cfg.opportunityCanvas("demo-map")
	if err != nil {
		t.Fatal(err)
	}
	out := s.cfg.toOpportunityCanvasJSON(canvas)

	if out.OpportunityURI != "livt://opportunity/demo-map" {
		t.Errorf("opportunity_uri = %q, want the opportunity this canvas fills in", out.OpportunityURI)
	}
	if len(out.Boxes) != 10 {
		t.Fatalf("got %d boxes, want the canvas's ten", len(out.Boxes))
	}
	for _, b := range out.Boxes {
		if b.Name == "" || b.Prompt == "" {
			t.Errorf("box %q is missing its name or prompt", b.Key)
		}
		// An unanswered box serializes as [] rather than null, so a consumer
		// reads it as the empty box it is.
		if b.Items == nil {
			t.Errorf("box %q has nil items", b.Key)
		}
	}
	if len(out.UbiquitousTerms) != 1 || out.UbiquitousTerms[0].URI != uri.Term("", "story") {
		t.Errorf("ubiquitous_terms = %+v", out.UbiquitousTerms)
	}
}

// A key that names nothing is not found, and one that tries to walk out of the
// directory is refused before it reaches the filesystem.
func TestOpportunityLookupsRefuseUnknownAndTraversingKeys(t *testing.T) {
	c := newTestServer(t).cfg
	for _, key := range []string{"nope", "..", "../../etc/passwd"} {
		if _, err := c.opportunity(key); err == nil {
			t.Errorf("opportunity(%q) resolved", key)
		}
		if _, err := c.opportunityCanvas(key); err == nil {
			t.Errorf("opportunityCanvas(%q) resolved", key)
		}
	}
}

// A canvas stands on its own: it resolves whether or not an opportunity file
// was committed for it, the way a mapping resolves for an uncommitted story.
func TestOpportunityCanvasResolvesWithoutItsOpportunityFile(t *testing.T) {
	s := newTestServer(t)
	writeFile(t, s.cfg.canvasesDir()+"/orphan.yaml", "canvas:\n  problems:\n    - 課題\n")

	if _, err := s.cfg.opportunityCanvas("orphan"); err != nil {
		t.Fatalf("canvas without an opportunity file: %v", err)
	}
	if _, err := s.cfg.opportunity("orphan"); err == nil {
		t.Fatal("expected the opportunity itself to be not found")
	}
}
