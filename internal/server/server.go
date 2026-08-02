// Package server builds artifacts and serves them with live reload, so edits
// to discovery and story files are previewed in the browser while editing.
package server

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/boykush/livt/internal/builder"
)

const liveReloadPath = "/livereload"

// liveReloadScript reconnects the browser to the SSE endpoint and reloads the
// page whenever the server signals a rebuild. It is injected into served HTML
// only; the static build output produced by `livt build` is left untouched.
//
// The iframe guard is essential: index pages embed every mapping and story map
// as an iframe preview, so without it each preview would open its own SSE
// connection. A handful of previews exhausts the browser's per-origin HTTP/1.1
// connection limit (typically 6), starving navigation requests and deadlocking
// the page. Only the top-level document keeps a connection, and its reload
// refreshes the embedded previews along with it.
const liveReloadScript = `<script>
(function () {
  if (window.self !== window.top) return;
  var es = new EventSource("` + liveReloadPath + `");
  es.onmessage = function () { location.reload(); };
})();
</script>`

const watchInterval = 500 * time.Millisecond

// Serve builds once, then serves OutDir while watching the builder's input
// directories. A change rebuilds and reloads every connected browser.
func Serve(b *builder.Builder, port int) error {
	fmt.Printf("Building to %s/\n", b.OutDir)
	if err := b.Build(); err != nil {
		return err
	}

	lr := newLiveReload()
	watchDirs := []string{b.MappingsDir, b.StoriesDir, b.USMDir, b.PersonasDir, b.UbiquitousDir}
	go watch(watchDirs, watchInterval, func() {
		fmt.Println("Change detected, rebuilding...")
		if err := b.Build(); err != nil {
			fmt.Printf("  build failed: %v\n", err)
			return
		}
		lr.broadcast()
	})

	mux := http.NewServeMux()
	mux.HandleFunc(liveReloadPath, lr.handleSSE)
	mux.Handle("/", htmlInjector(b.OutDir))

	fmt.Printf("Serving on http://localhost:%d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

// htmlInjector serves files from dir, injecting the live-reload script into
// HTML responses. Non-HTML requests fall through to the standard file server.
func htmlInjector(dir string) http.Handler {
	root := http.Dir(dir)
	fileServer := http.FileServer(root)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := path.Clean(r.URL.Path)
		if strings.HasSuffix(r.URL.Path, "/") {
			name = path.Join(name, "index.html")
		}
		if !strings.HasSuffix(name, ".html") {
			fileServer.ServeHTTP(w, r)
			return
		}

		f, err := root.Open(name)
		if err != nil {
			fileServer.ServeHTTP(w, r) // keep 404 handling consistent
			return
		}
		defer f.Close()

		html, err := io.ReadAll(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write(injectScript(html)) // headers already sent; nothing to recover
	})
}

// injectScript inserts the live-reload script just before the last </body>,
// falling back to appending it when no body tag is present.
func injectScript(html []byte) []byte {
	marker := []byte("</body>")
	idx := bytes.LastIndex(html, marker)
	if idx < 0 {
		return append(html, []byte(liveReloadScript)...)
	}
	out := make([]byte, 0, len(html)+len(liveReloadScript))
	out = append(out, html[:idx]...)
	out = append(out, []byte(liveReloadScript)...)
	out = append(out, html[idx:]...)
	return out
}

// liveReload broadcasts rebuild notifications to connected browsers over SSE.
type liveReload struct {
	mu      sync.Mutex
	clients map[chan struct{}]struct{}
}

func newLiveReload() *liveReload {
	return &liveReload{clients: make(map[chan struct{}]struct{})}
}

func (lr *liveReload) subscribe() chan struct{} {
	ch := make(chan struct{}, 1)
	lr.mu.Lock()
	lr.clients[ch] = struct{}{}
	lr.mu.Unlock()
	return ch
}

func (lr *liveReload) unsubscribe(ch chan struct{}) {
	lr.mu.Lock()
	delete(lr.clients, ch)
	lr.mu.Unlock()
}

func (lr *liveReload) broadcast() {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	for ch := range lr.clients {
		select {
		case ch <- struct{}{}:
		default: // a reload is already pending for this client
		}
	}
}

func (lr *liveReload) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := lr.subscribe()
	defer lr.unsubscribe(ch)
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ch:
			fmt.Fprint(w, "data: reload\n\n")
			flusher.Flush()
		}
	}
}

// watch polls dirs and invokes onChange whenever a watched file is added,
// removed, or modified. It blocks, so run it in a goroutine.
func watch(dirs []string, interval time.Duration, onChange func()) {
	prev := snapshot(dirs)
	for {
		time.Sleep(interval)
		cur := snapshot(dirs)
		if !sameSnapshot(prev, cur) {
			prev = cur
			onChange()
		}
	}
}

// snapshot fingerprints every file under dirs by modtime and size. Missing
// directories simply contribute no entries.
func snapshot(dirs []string) map[string]string {
	state := make(map[string]string)
	for _, dir := range dirs {
		_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			info, err := d.Info()
			if err != nil {
				return nil
			}
			state[p] = fmt.Sprintf("%d-%d", info.ModTime().UnixNano(), info.Size())
			return nil
		})
	}
	return state
}

func sameSnapshot(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
