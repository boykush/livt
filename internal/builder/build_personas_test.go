package builder

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/boykush/livt/internal/domain"
)

func writePersona(t *testing.T, dir, key, name string) {
	t.Helper()
	body := []byte("---\nname: " + name + "\n---\n\nDescription of " + name + ".\n")
	if err := os.WriteFile(filepath.Join(dir, key+".md"), body, 0o644); err != nil {
		t.Fatal(err)
	}
}

func buildPersonasHTML(t *testing.T, personasDir string, storyIndex map[string][]personaStoryRef) string {
	t.Helper()
	outDir := t.TempDir()
	b := Builder{PersonasDir: personasDir, OutDir: outDir}
	if err := b.buildPersonas(storyIndex); err != nil {
		t.Fatal(err)
	}
	out, err := os.ReadFile(filepath.Join(outDir, "personas.html"))
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

// livt://mapping/look-up-story-personas/rule/R-01/example/EX-03: each persona is
// a row addressable by its key, carrying the name and the description.
func TestBuildPersonasRendersPersonaAsAnchoredRow(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "coding-agent.md"), []byte("---\nname: コーディングエージェント\n---\n\nMCP越しに仕様を読むAIエージェント。\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	html := buildPersonasHTML(t, dir, nil)

	for _, want := range []string{`id="coding-agent"`, "コーディングエージェント", "MCP越しに仕様を読むAIエージェント。"} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected the persona row to carry %q", want)
		}
	}
}

// livt://mapping/look-up-story-personas/rule/R-03/example/EX-01: the row lists
// the stories written for the persona, which is the reverse of the link the
// story's frontmatter holds.
func TestBuildPersonasListsTheStoriesWrittenForEachPersona(t *testing.T) {
	dir := t.TempDir()
	writePersona(t, dir, "coding-agent", "コーディングエージェント")
	writePersona(t, dir, "living-document-reader", "リビングドキュメントの閲覧者")

	html := buildPersonasHTML(t, dir, map[string][]personaStoryRef{
		"coding-agent": {{Name: "仕様を見て自動化する", Path: "story/automate.html"}},
	})

	if !strings.Contains(html, `href="story/automate.html"`) || !strings.Contains(html, "仕様を見て自動化する") {
		t.Fatal("expected the persona's story to be listed and linked")
	}
	// livt://mapping/look-up-story-personas/rule/R-03/example/EX-02
	if !strings.Contains(html, "No stories yet") {
		t.Fatal("expected a persona no story names to say so rather than render blank")
	}
}

func TestBuildPersonasWithoutPersonasRendersEmptyState(t *testing.T) {
	html := buildPersonasHTML(t, filepath.Join(t.TempDir(), "missing"), nil)
	if !strings.Contains(html, "No personas found.") {
		t.Fatal("expected empty-state message when no personas exist")
	}
}

// livt://mapping/look-up-story-personas/rule/R-02/example/EX-02: a committed
// persona resolves to its display name and links to its row.
func TestResolvePersonaCardLinksACommittedPersona(t *testing.T) {
	dir := t.TempDir()
	writePersona(t, dir, "coding-agent", "コーディングエージェント")
	b := Builder{PersonasDir: dir}

	card := b.resolvePersonaCard("coding-agent", "../")
	if card == nil {
		t.Fatal("expected a card for a committed persona")
	}
	if card.Name != "コーディングエージェント" {
		t.Fatalf("got name %q", card.Name)
	}
	if card.Href != "../personas.html#coding-agent" {
		t.Fatalf("got href %q", card.Href)
	}
}

// A key with no persona file behind it still reads on the page, so a story can
// name its actor before the persona is committed.
func TestResolvePersonaCardDegradesWhenThePersonaIsMissing(t *testing.T) {
	b := Builder{PersonasDir: t.TempDir()}

	card := b.resolvePersonaCard("not-committed", "")
	if card == nil || card.Name != "not-committed" {
		t.Fatalf("got card %+v, want the bare key kept", card)
	}
	if card.Href != "" {
		t.Fatalf("got href %q, want no link to a page that was never built", card.Href)
	}
}

// A story's persona key reaches this from frontmatter, so it must not resolve
// to a file outside the personas directory.
func TestResolvePersonaCardRejectsKeysThatEscapeTheDirectory(t *testing.T) {
	dir := t.TempDir()
	writePersona(t, dir, "secret", "秘密")
	b := Builder{PersonasDir: filepath.Join(dir, "nested")}

	card := b.resolvePersonaCard("../secret", "")
	if card == nil || card.Href != "" {
		t.Fatalf("got card %+v, want no link out of the personas directory", card)
	}
}

func TestResolvePersonaCardWithoutAPersonaIsNil(t *testing.T) {
	b := Builder{PersonasDir: t.TempDir()}
	if card := b.resolvePersonaCard("", ""); card != nil {
		t.Fatalf("got card %+v, want none for a story that names no actor", card)
	}
}

// livt://mapping/look-up-story-personas/rule/R-02/example/EX-02: the story page
// shows the actor beside the story key, linked to its row.
func TestRenderStoryShowsThePersonaChip(t *testing.T) {
	story := &domain.Story{Key: domain.StoryKey{Value: "automate"}, Name: "Automate", Persona: "coding-agent"}
	card := &personaCard{Name: "コーディングエージェント", Href: "../personas.html#coding-agent"}

	var buf bytes.Buffer
	if err := renderStory(&buf, story, card, "", nil); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	if !strings.Contains(html, `href="../personas.html#coding-agent"`) {
		t.Fatal("expected the persona chip to link to its row")
	}
	if !strings.Contains(html, "コーディングエージェント") {
		t.Fatal("expected the persona display name on the story page")
	}
}
