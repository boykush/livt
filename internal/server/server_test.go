package server

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInjectScriptInsertsBeforeBody(t *testing.T) {
	html := []byte("<html><body><h1>hi</h1></body></html>")
	out := injectScript(html)

	if !bytes.Contains(out, []byte(liveReloadScript)) {
		t.Fatal("expected live-reload script to be injected")
	}
	if got := bytes.Index(out, []byte(liveReloadScript)); got >= bytes.LastIndex(out, []byte("</body>")) {
		t.Fatal("expected script to be injected before </body>")
	}
}

func TestInjectScriptAppendsWhenNoBody(t *testing.T) {
	out := injectScript([]byte("<h1>hi</h1>"))
	if !strings.HasSuffix(string(out), liveReloadScript) {
		t.Fatal("expected script to be appended when no </body> is present")
	}
}

// Index pages embed each mapping and story map as an iframe preview. The
// injected script must bail out inside frames so previews don't each open an
// SSE connection and exhaust the browser's per-origin connection limit.
func TestLiveReloadScriptSkipsIframes(t *testing.T) {
	guard := "if (window.self !== window.top) return;"
	idx := strings.Index(liveReloadScript, guard)
	if idx < 0 {
		t.Fatal("expected live-reload script to skip iframe contexts")
	}
	if src := strings.Index(liveReloadScript, "new EventSource"); src < idx {
		t.Fatal("expected iframe guard to run before opening the SSE connection")
	}
}

func TestSnapshotDetectsAddedFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.md"), []byte("a"), 0o644); err != nil {
		t.Fatal(err)
	}

	before := snapshot([]string{dir})
	if err := os.WriteFile(filepath.Join(dir, "b.md"), []byte("b"), 0o644); err != nil {
		t.Fatal(err)
	}
	after := snapshot([]string{dir})

	if sameSnapshot(before, after) {
		t.Fatal("expected snapshot to change after adding a file")
	}
}

func TestSnapshotDetectsModifiedContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.md")
	if err := os.WriteFile(path, []byte("a"), 0o644); err != nil {
		t.Fatal(err)
	}

	before := snapshot([]string{dir})
	if err := os.WriteFile(path, []byte("aa"), 0o644); err != nil {
		t.Fatal(err)
	}
	after := snapshot([]string{dir})

	if sameSnapshot(before, after) {
		t.Fatal("expected snapshot to change after modifying file content")
	}
}

func TestSnapshotIgnoresMissingDir(t *testing.T) {
	state := snapshot([]string{filepath.Join(t.TempDir(), "does-not-exist")})
	if len(state) != 0 {
		t.Fatalf("got %d entries for missing dir, want 0", len(state))
	}
}
