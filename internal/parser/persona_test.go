package parser

import (
	"os"
	"path/filepath"
	"testing"
)

// writePersona creates personas/{key}.md.
func writePersona(t *testing.T, dir, key, name string) {
	t.Helper()
	data := []byte("---\nname: " + name + "\n---\n\nDescription of " + name + ".\n")
	if err := os.WriteFile(filepath.Join(dir, key+".md"), data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// livt://mapping/list-personas-on-story-map/rule/R-01/example/EX-01
func TestParsePersonaReadsNameAndBody(t *testing.T) {
	dir := t.TempDir()
	data := []byte("---\nname: コーディングエージェント\n---\n\nlivtリポジトリをMCP越しに読むAIエージェント。\n")
	if err := os.WriteFile(filepath.Join(dir, "coding-agent.md"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	persona, err := ParsePersona(dir, "coding-agent")
	if err != nil {
		t.Fatal(err)
	}

	if persona.Key != "coding-agent" {
		t.Fatalf("got key %q, want coding-agent", persona.Key)
	}
	if persona.Name != "コーディングエージェント" {
		t.Fatalf("got name %q", persona.Name)
	}
	if persona.Body != "livtリポジトリをMCP越しに読むAIエージェント。" {
		t.Fatalf("got body %q", persona.Body)
	}
}

// A key is externally supplied — it arrives from a story's frontmatter, an MCP
// client, or a CLI argument — so it must not be able to walk out of the personas
// directory. Unlike a term reference it holds no context, so a key with any
// separator in it is malformed rather than scoped.
func TestParsePersonaRejectsKeysThatEscapeTheDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "secret.md"), []byte("---\nname: secret\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{"../secret", "billing/../secret", "a/b", "reader/", "/reader", ""} {
		if _, err := ParsePersona(dir, key); err == nil {
			t.Errorf("ParsePersona(%q) succeeded, want a rejected key", key)
		}
	}
}

// livt://mapping/list-personas-on-story-map/rule/R-01/example/EX-02: the list is
// every persona file, in the order the filesystem lays them out.
func TestParseAllPersonasReturnsEveryFileInOrder(t *testing.T) {
	dir := t.TempDir()
	writePersona(t, dir, "living-document-reader", "リビングドキュメントの閲覧者")
	writePersona(t, dir, "coding-agent", "コーディングエージェント")

	personas, err := ParseAllPersonas(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(personas) != 2 {
		t.Fatalf("got %d personas, want 2", len(personas))
	}
	if personas[0].Key != "coding-agent" || personas[1].Key != "living-document-reader" {
		t.Fatalf("got order %q, %q", personas[0].Key, personas[1].Key)
	}
}

// A persona has no context to cut it by, so a directory under personas/ holds
// nothing the list can address.
func TestParseAllPersonasIgnoresNestedDirectories(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	writePersona(t, dir, filepath.Join("nested", "reader"), "深すぎる")

	personas, err := ParseAllPersonas(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(personas) != 0 {
		t.Fatalf("got %d personas, want none — personas are flat", len(personas))
	}
}

func TestParseAllPersonasMissingDirIsEmpty(t *testing.T) {
	personas, err := ParseAllPersonas(filepath.Join(t.TempDir(), "does-not-exist"))
	if err != nil {
		t.Fatal(err)
	}
	if len(personas) != 0 {
		t.Fatalf("got %d personas, want 0", len(personas))
	}
}
