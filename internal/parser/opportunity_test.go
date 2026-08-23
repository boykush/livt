package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseOpportunityReadsNameStatementAndMeta(t *testing.T) {
	dir := t.TempDir()
	path := writeStoryFile(t, dir, "collaborative-discovery.md",
		"---\nname: 協働ディスカバリー\nrepos:\n  - boykush/livt\n---\n\n共通理解がボードとともに失われる。\nテキストとして残すと、やり直しの時間が減る。\n")

	o, err := ParseOpportunity(path)
	if err != nil {
		t.Fatal(err)
	}
	if o.Key.Value != "collaborative-discovery" {
		t.Fatalf("got key %q", o.Key.Value)
	}
	if o.Name != "協働ディスカバリー" {
		t.Fatalf("got name %q", o.Name)
	}
	// The statement is the body, not the name: an opportunity is a sentence about
	// whose problem it is and what solving it is worth, the way a story is.
	if o.Body != "共通理解がボードとともに失われる。\nテキストとして残すと、やり直しの時間が減る。" {
		t.Fatalf("got body %q", o.Body)
	}
	if len(o.Meta) != 1 || o.Meta[0].Key != "repos" || o.Meta[0].Value != "boykush/livt" {
		t.Fatalf("got meta %+v", o.Meta)
	}
}

func TestParseOpportunityCanvasReadsEveryBox(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "demo.yaml")
	if err := os.WriteFile(path, []byte(
		"canvas:\n"+
			"  solution-ideas:\n    - アイデア\n"+
			"  problems:\n    - 課題1\n    - 課題2\n"+
			"  users-and-customers:\n    - 利用者\n"+
			"  solutions-today:\n    - 今のやり方\n"+
			"  business-challenges:\n    - 事業課題\n"+
			"  user-value:\n    - やること\n"+
			"  user-metrics:\n    - 指標\n"+
			"  adoption-strategy:\n    - 広め方\n"+
			"  business-impact:\n    - 事業影響\n"+
			"  budget:\n    - 予算\n"+
			"ubiquitous:\n  - discovery\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	canvas, err := ParseOpportunityCanvas(path)
	if err != nil {
		t.Fatal(err)
	}
	if canvas.OpportunityKey.Value != "demo" {
		t.Fatalf("got opportunity key %q", canvas.OpportunityKey.Value)
	}
	boxes := canvas.Boxes()
	if len(boxes) != 10 {
		t.Fatalf("got %d boxes, want the canvas's ten", len(boxes))
	}
	for _, b := range boxes {
		if len(b.Items) == 0 {
			t.Errorf("box %q parsed empty; its yaml key does not match what the file authors", b.Key)
		}
	}
	if len(canvas.Ubiquitous) != 1 || canvas.Ubiquitous[0] != "discovery" {
		t.Fatalf("got ubiquitous %v", canvas.Ubiquitous)
	}
}

// A canvas filled in only part way still parses, and the boxes nobody answered
// come back empty rather than missing — the gap is the point.
func TestParseOpportunityCanvasKeepsUnansweredBoxes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "partial.yaml")
	if err := os.WriteFile(path, []byte("canvas:\n  problems:\n    - 課題\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	canvas, err := ParseOpportunityCanvas(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(canvas.Boxes()) != 10 {
		t.Fatalf("got %d boxes, want ten", len(canvas.Boxes()))
	}
	filled := 0
	for _, b := range canvas.Boxes() {
		if len(b.Items) > 0 {
			filled++
		}
	}
	if filled != 1 {
		t.Fatalf("got %d filled boxes, want 1", filled)
	}
}
