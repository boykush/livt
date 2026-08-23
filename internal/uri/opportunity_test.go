package uri

import "testing"

// The opportunity and its canvas sit under two prefixes that share a stem, so
// the risk is one shape parsing as the other. It cannot: the byte after
// "opportunity" is "-" for the canvas and "/" for the opportunity, which is
// what keeps every key on exactly one shape.
func TestOpportunityAndCanvasURIsNeverParseAsEachOther(t *testing.T) {
	opportunity := Opportunity("collaborative-discovery")
	canvas := OpportunityCanvas("collaborative-discovery")

	if _, ok := ParseOpportunity(canvas); ok {
		t.Errorf("ParseOpportunity(%q) parsed a canvas URI as an opportunity", canvas)
	}
	if _, ok := ParseOpportunityCanvas(opportunity); ok {
		t.Errorf("ParseOpportunityCanvas(%q) parsed an opportunity URI as a canvas", opportunity)
	}

	if got, ok := ParseOpportunity(opportunity); !ok || got != "collaborative-discovery" {
		t.Errorf("ParseOpportunity(%q) = (%q, %v)", opportunity, got, ok)
	}
	if got, ok := ParseOpportunityCanvas(canvas); !ok || got != "collaborative-discovery" {
		t.Errorf("ParseOpportunityCanvas(%q) = (%q, %v)", canvas, got, ok)
	}
}

// An externally-supplied key must not walk out of the opportunities directory.
func TestOpportunityURIsRejectTraversalKeys(t *testing.T) {
	for _, s := range []string{
		"livt://opportunity/",
		"livt://opportunity/..",
		"livt://opportunity/../../etc/passwd",
		"livt://opportunity-canvas/",
		"livt://opportunity-canvas/../secrets",
	} {
		if _, ok := Parse(s); ok {
			t.Errorf("Parse(%q) accepted a traversal key", s)
		}
	}
}

// A URI round-trips through Parse and back, and lands on the page the build
// writes for it.
func TestOpportunityURIsRoundTripAndLand(t *testing.T) {
	cases := []struct{ uri, page string }{
		{Opportunity("demo"), "opportunity/demo.html"},
		{OpportunityCanvas("demo"), "opportunity-canvas/demo.html"},
	}
	for _, c := range cases {
		p, ok := Parse(c.uri)
		if !ok {
			t.Fatalf("Parse(%q) failed", c.uri)
		}
		if got := p.String(); got != c.uri {
			t.Errorf("round trip of %q = %q", c.uri, got)
		}
		if got := p.Page(); got != c.page {
			t.Errorf("Page(%q) = %q, want %q", c.uri, got, c.page)
		}
	}
}
