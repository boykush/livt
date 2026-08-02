package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func writeStoryFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestParseStoryReadsNameOrderedMetaAndBody(t *testing.T) {
	dir := t.TempDir()
	path := writeStoryFile(t, dir, "preview.md",
		"---\nname: Preview the map\nstatus: draft\nticketUrl: https://example.com/123\n---\n\nAs a maintainer\nI want to preview\n")

	story, err := ParseStory(path)
	if err != nil {
		t.Fatal(err)
	}
	if story.Key.Value != "preview" {
		t.Fatalf("got key %q, want preview", story.Key.Value)
	}
	if story.Name != "Preview the map" {
		t.Fatalf("got name %q", story.Name)
	}
	if story.Body != "As a maintainer\nI want to preview" {
		t.Fatalf("got body %q", story.Body)
	}
	if len(story.Meta) != 2 {
		t.Fatalf("got %d meta fields, want 2", len(story.Meta))
	}
	// Source order is preserved: status before ticketUrl.
	if story.Meta[0].Key != "status" || story.Meta[0].Value != "draft" {
		t.Fatalf("meta[0] = %+v", story.Meta[0])
	}
	if story.Meta[1].Key != "ticketUrl" || story.Meta[1].Value != "https://example.com/123" {
		t.Fatalf("meta[1] = %+v", story.Meta[1])
	}
}

func TestParseStoryJoinsSequenceMeta(t *testing.T) {
	dir := t.TempDir()
	path := writeStoryFile(t, dir, "tagged.md",
		"---\nname: Tagged\ntags:\n  - alpha\n  - beta\n---\nbody\n")

	story, err := ParseStory(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(story.Meta) != 1 || story.Meta[0].Key != "tags" || story.Meta[0].Value != "alpha, beta" {
		t.Fatalf("got meta %+v", story.Meta)
	}
}

func TestParseStoryNoFrontmatter(t *testing.T) {
	dir := t.TempDir()
	path := writeStoryFile(t, dir, "plain.md", "just a body\nsecond line\n")

	story, err := ParseStory(path)
	if err != nil {
		t.Fatal(err)
	}
	if story.Name != "" {
		t.Fatalf("got name %q, want empty", story.Name)
	}
	if len(story.Meta) != 0 {
		t.Fatalf("got %d meta, want 0", len(story.Meta))
	}
	if story.Body != "just a body\nsecond line" {
		t.Fatalf("got body %q", story.Body)
	}
}

// livt://mapping/list-personas-on-story-map/rule/R-03/example/EX-01: persona is
// reserved like name, so it addresses the actor rather than landing among the
// free-form metadata rows.
func TestParseStoryReadsPersonaAsAReservedField(t *testing.T) {
	dir := t.TempDir()
	path := writeStoryFile(t, dir, "automate.md",
		"---\nname: Automate from the spec\npersona: coding-agent\nissue: https://example.com/1\n---\nbody\n")

	story, err := ParseStory(path)
	if err != nil {
		t.Fatal(err)
	}
	if story.Persona != "coding-agent" {
		t.Fatalf("got persona %q, want coding-agent", story.Persona)
	}
	if len(story.Meta) != 1 || story.Meta[0].Key != "issue" {
		t.Fatalf("got meta %+v, want persona kept out of it", story.Meta)
	}
}

// livt://mapping/list-personas-on-story-map/rule/R-03/example/EX-03: naming an actor
// stays optional, so a story written before the personas existed still parses.
func TestParseStoryWithoutPersonaIsEmpty(t *testing.T) {
	dir := t.TempDir()
	path := writeStoryFile(t, dir, "n.md", "---\nname: Only\n---\nbody\n")

	story, err := ParseStory(path)
	if err != nil {
		t.Fatal(err)
	}
	if story.Persona != "" {
		t.Fatalf("got persona %q, want empty", story.Persona)
	}
}

func TestParseStoryNameOnlyHasNoMeta(t *testing.T) {
	dir := t.TempDir()
	path := writeStoryFile(t, dir, "n.md", "---\nname: Only\n---\nbody\n")

	story, err := ParseStory(path)
	if err != nil {
		t.Fatal(err)
	}
	if story.Name != "Only" || len(story.Meta) != 0 {
		t.Fatalf("got name %q meta %+v", story.Name, story.Meta)
	}
}
