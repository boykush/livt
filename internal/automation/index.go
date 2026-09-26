package automation

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"github.com/boykush/livt/internal/uri"
)

// Index answers what automates a point of the livt repository. An empty index is the
// ordinary state of a livt repository nobody has pointed an implementation at,
// and it means the site shows no automation at all rather than showing
// everything as un-automated.
type Index struct {
	byURI   map[string][]domain.Automation
	reports []Report
}

// Load reads every report under dir. A missing directory is not an error: a
// livt repository is complete without one.
func Load(dir string) (*Index, error) {
	idx := &Index{byURI: map[string][]domain.Automation{}}
	entries, err := reportFiles(dir)
	if err != nil {
		return nil, err
	}
	for _, path := range entries {
		data, err := os.ReadFile(path) // #nosec G304 -- the path comes from walking the directory the caller named
		if err != nil {
			return nil, err
		}
		var report Report
		if err := json.Unmarshal(data, &report); err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		idx.reports = append(idx.reports, report)
		for _, c := range report.Citations {
			idx.byURI[c.URI] = append(idx.byURI[c.URI], domain.Automation{
				Repo: report.Repo, Rev: report.Rev, File: c.File, Line: c.Line, URL: c.URL,
			})
		}
	}
	sort.Slice(idx.reports, func(i, j int) bool { return idx.reports[i].Repo < idx.reports[j].Repo })
	return idx, nil
}

func reportFiles(dir string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return filepath.SkipAll
			}
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".json") {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)
	return paths, nil
}

// Empty reports whether anything was collected at all, which is what tells
// "nobody has looked" apart from "looked and found nothing".
func (i *Index) Empty() bool { return i == nil || len(i.reports) == 0 }

// Reports are what the index was built from, so a page can say whose answer it
// is showing and how old it is rather than presenting it as timeless.
func (i *Index) Reports() []Report {
	if i == nil {
		return nil
	}
	return i.reports
}

// For returns what automates this exact URI. A rule and its examples are asked
// for separately and neither answers for the other: deriving one from the
// other would be an inference, and the test author already said which they
// meant.
func (i *Index) For(uri string) []domain.Automation {
	if i == nil {
		return nil
	}
	return i.byURI[uri]
}

// Attach hangs the collected citations onto a mapping's rules and examples.
// The mapping on disk says nothing about automation, so this is the only place
// a board learns what the tests cover.
func (i *Index) Attach(em *domain.ExampleMapping) {
	if i == nil || em == nil {
		return
	}
	key := em.StoryKey.Value
	for ri := range em.Rules {
		r := &em.Rules[ri]
		r.Automations = i.For(uri.Rule(key, r.ID))
		for ei := range r.Examples {
			ex := &r.Examples[ei]
			ex.Automations = i.For(uri.Example(key, r.ID, ex.ID))
		}
	}
}
