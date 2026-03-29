package server

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/parser"
)

type Server struct {
	mappingsDir string
}

func Start(port int, mappingsDir string) error {
	s := &Server{mappingsDir: mappingsDir}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/mapping/", s.handleMapping)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files, err := filepath.Glob(filepath.Join(s.mappingsDir, "*.yaml"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var mappings []indexEntry
	for _, f := range files {
		em, err := parser.ParseExampleMapping(f)
		if err != nil {
			continue
		}
		mappings = append(mappings, indexEntry{
			Story: em.Story,
			Path:  "/mapping/" + em.Story,
		})
	}

	renderIndex(w, mappings)
}

func (s *Server) handleMapping(w http.ResponseWriter, r *http.Request) {
	storyKey := strings.TrimPrefix(r.URL.Path, "/mapping/")
	if storyKey == "" {
		http.NotFound(w, r)
		return
	}

	path := filepath.Join(s.mappingsDir, storyKey+".yaml")
	em, err := parser.ParseExampleMapping(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load mapping: %v", err), http.StatusNotFound)
		return
	}

	renderMapping(w, em)
}

type indexEntry struct {
	Story string
	Path  string
}
