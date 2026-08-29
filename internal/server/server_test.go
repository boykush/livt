package server

import (
	"bufio"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/boykush/livt/internal/builder"
	"github.com/boykush/livt/internal/config"
	"github.com/boykush/livt/internal/i18n"
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

func TestHTMLInjectorInjectsIntoServedHTML(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html><body>hi</body></html>"), 0o644); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	htmlInjector(dir).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), liveReloadScript) {
		t.Error("expected the live-reload script in a served HTML page")
	}
	// A cached page would keep showing the pre-rebuild output after a reload,
	// which is the one thing serve exists to avoid.
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control = %q, want no-store", got)
	}
}

// Only HTML may be rewritten: a script tag appended to a stylesheet or a JSON
// payload corrupts it, and the browser gives no hint why.
func TestHTMLInjectorLeavesNonHTMLUntouched(t *testing.T) {
	dir := t.TempDir()
	const css = "body { color: red }"
	if err := os.WriteFile(filepath.Join(dir, "app.css"), []byte(css), 0o644); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	htmlInjector(dir).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/app.css", nil))

	if rec.Body.String() != css {
		t.Errorf("body = %q, want the stylesheet served verbatim", rec.Body.String())
	}
}

// A missing page must still 404 rather than be answered with a body holding
// nothing but the injected script.
func TestHTMLInjectorKeeps404ForMissingHTML(t *testing.T) {
	rec := httptest.NewRecorder()
	htmlInjector(t.TempDir()).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/gone.html", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if strings.Contains(rec.Body.String(), liveReloadScript) {
		t.Error("expected no live-reload script on a 404")
	}
}

func TestBroadcastReachesEverySubscriber(t *testing.T) {
	lr := newLiveReload()
	a, b := lr.subscribe(), lr.subscribe()

	lr.broadcast()

	for i, ch := range []chan struct{}{a, b} {
		select {
		case <-ch:
		default:
			t.Errorf("subscriber %d was not notified", i)
		}
	}
}

// broadcast runs on the rebuild goroutine, so a client that has not drained its
// last notification must not stall it: one pending reload already says
// everything a second would.
func TestBroadcastDoesNotBlockOnAPendingClient(t *testing.T) {
	lr := newLiveReload()
	ch := lr.subscribe()

	done := make(chan struct{})
	go func() {
		lr.broadcast()
		lr.broadcast()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("broadcast blocked on a client with a reload already pending")
	}
	if len(ch) != 1 {
		t.Errorf("pending notifications = %d, want 1", len(ch))
	}
}

func TestUnsubscribeStopsDelivery(t *testing.T) {
	lr := newLiveReload()
	ch := lr.subscribe()
	lr.unsubscribe(ch)

	lr.broadcast()

	if len(ch) != 0 {
		t.Error("expected no notification after unsubscribing")
	}
}

func TestHandleSSEStreamsReloadsAndReleasesTheClient(t *testing.T) {
	lr := newLiveReload()
	srv := httptest.NewServer(http.HandlerFunc(lr.handleSSE))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	// The handler subscribes before flushing the headers, so a returned response
	// means the connection is registered and the broadcast below cannot race it.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if got := resp.Header.Get("Content-Type"); got != "text/event-stream" {
		t.Errorf("Content-Type = %q, want text/event-stream", got)
	}

	lr.broadcast()

	line, err := bufio.NewReader(resp.Body).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if line != "data: reload\n" {
		t.Errorf("event = %q, want %q", line, "data: reload\n")
	}

	// A browser that navigates away must not leave its channel behind; every
	// page load opens a new one, so a leak here grows for the whole session.
	resp.Body.Close()
	if !eventually(func() bool { return lr.clientCount() == 0 }) {
		t.Error("expected the client to be unsubscribed once it disconnected")
	}
}

// clientCount reads the subscriber set under the lock, so the assertion above
// does not race the handler goroutine unsubscribing.
func (lr *liveReload) clientCount() int {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	return len(lr.clients)
}

func eventually(cond func() bool) bool {
	for deadline := time.Now().Add(2 * time.Second); time.Now().Before(deadline); {
		if cond() {
			return true
		}
		time.Sleep(5 * time.Millisecond)
	}
	return cond()
}

// watch polls, so the test keeps writing until it reports the change rather than
// timing one write against the tick. The goroutine has no stop signal and
// outlives the test, which is why its callback only ever does a buffered send.
func TestWatchReportsAChangeUnderAWatchedDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "story.md")
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	fired := make(chan struct{}, 1)
	go watch([]string{dir}, time.Millisecond, func() {
		select {
		case fired <- struct{}{}:
		default:
		}
	})

	deadline := time.After(5 * time.Second)
	for i := 2; ; i++ {
		if err := os.WriteFile(path, []byte(strings.Repeat("x", i)), 0o644); err != nil {
			t.Fatal(err)
		}
		select {
		case <-fired:
			return
		case <-deadline:
			t.Fatal("watch did not report a change to a watched file")
		case <-time.After(10 * time.Millisecond):
		}
	}
}

// The config is watched as a bare file, not a directory, so the fingerprint has
// to cover one — a language switch is a single-file edit.
func TestSnapshotFingerprintsAWatchedFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), config.Path)
	if err := os.WriteFile(path, []byte("lang: en\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	before := snapshot([]string{path})
	if len(before) != 1 {
		t.Fatalf("got %d entries for a watched file, want 1", len(before))
	}
	if err := os.WriteFile(path, []byte("lang: ja\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if sameSnapshot(before, snapshot([]string{path})) {
		t.Fatal("expected the snapshot to change after editing the watched file")
	}
}

func TestApplyConfigPicksUpALanguageSwitch(t *testing.T) {
	path := filepath.Join(t.TempDir(), config.Path)
	if err := os.WriteFile(path, []byte("lang: ja\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	b := &builder.Builder{}
	applyConfig(b, path)
	if b.Lang != i18n.Ja {
		t.Fatalf("lang = %q, want %q", b.Lang, i18n.Ja)
	}
}

// A half-typed config is the normal state of a file being edited, so the
// running server keeps the language it last built with instead of tearing the
// preview down over it.
func TestApplyConfigKeepsTheLastGoodLanguage(t *testing.T) {
	path := filepath.Join(t.TempDir(), config.Path)
	if err := os.WriteFile(path, []byte("lang: [ja\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	b := &builder.Builder{Lang: i18n.Ja}
	applyConfig(b, path)
	if b.Lang != i18n.Ja {
		t.Fatalf("lang = %q, want the last good %q", b.Lang, i18n.Ja)
	}
}

// Every input the build reads has to be watched, or an edit to it changes the
// site without the preview ever saying so. Opportunities and their canvases
// were exactly that gap: built into every page's sidebar counts and the
// Opportunities hub, yet absent from the watch list.
func TestWatchTargetsCoverEveryBuildInput(t *testing.T) {
	b := &builder.Builder{
		OpportunitiesDir: "opportunities",
		CanvasesDir:      filepath.Join("discoveries", "opportunity-canvases"),
		MappingsDir:      filepath.Join("discoveries", "example-mappings"),
		StoriesDir:       "stories",
		USMDir:           filepath.Join("discoveries", "usm"),
		UbiquitousDir:    "ubiquitous",
		OutDir:           "dist",
	}

	targets := watchTargets(b, config.Path)

	for _, want := range append(b.InputDirs(), config.Path) {
		if !slices.Contains(targets, want) {
			t.Errorf("watch targets omit %q", want)
		}
	}
	if slices.Contains(targets, b.OutDir) {
		t.Error("watch targets include OutDir; a rebuild would retrigger itself")
	}
}

// watchTargets appends to what the builder hands back, so it must not be able
// to write into a slice the builder still owns.
func TestWatchTargetsLeavesInputDirsIntact(t *testing.T) {
	b := &builder.Builder{OpportunitiesDir: "opportunities", UbiquitousDir: "ubiquitous"}
	before := b.InputDirs()

	watchTargets(b, config.Path)

	if !slices.Equal(before, b.InputDirs()) {
		t.Errorf("InputDirs = %v after watchTargets, want %v", b.InputDirs(), before)
	}
}
