package automation

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/boykush/livt/internal/uri"
)

// commentLead is what may precede the marker: whitespace and the punctuation
// every language starts a comment with. Word characters are excluded so the
// marker written inside code — a string literal, this file — is not read as a
// claim.
const commentLead = `^[\s/#*;%!<>'()\-]*`

var (
	// citation requires the line to hold the marker, one URI, and nothing
	// else. A trailing "and EX-02" would otherwise be dropped in silence,
	// which is how the old free-form citations lost half their examples.
	citation = regexp.MustCompile(commentLead + regexp.QuoteMeta(Marker) + `\s+(\S+)\s*$`)
	attempt  = regexp.MustCompile(commentLead + regexp.QuoteMeta(Marker) + `\b`)
)

// skipDirs are trees that hold other people's code. Hidden directories are
// skipped too, which is what keeps a nested worktree in .claude out of a scan
// of the repository around it.
var skipDirs = map[string]bool{"node_modules": true, "vendor": true, "dist": true, "target": true}

// maxFileSize keeps a scan away from data files a repository happens to carry.
const maxFileSize = 4 << 20

// Options carry what a checkout cannot answer for itself, all optional: an
// origin read from git fills the rest.
type Options struct {
	Repo        string
	Rev         string
	Base        string
	URLTemplate string
	Now         func() time.Time
}

// Scan reads a checkout and reports what its tests claim to automate. The
// warnings name lines that meant to be citations and are not, because a
// malformed claim that scans as nothing is the one failure nobody would see.
func Scan(root string, opts Options) (Report, []string, error) {
	now := opts.Now
	if now == nil {
		now = time.Now
	}
	report := Report{
		Repo:        opts.Repo,
		Rev:         opts.Rev,
		GeneratedAt: now().UTC().Truncate(time.Second),
		Citations:   []Citation{},
	}
	var warnings []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() {
			if path != root && (strings.HasPrefix(name, ".") || skipDirs[name]) {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}
		info, err := d.Info()
		if err != nil || info.Size() > maxFileSize {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		found, warned, err := scanFile(path, filepath.ToSlash(rel), opts, report.Rev)
		if err != nil {
			return err
		}
		report.Citations = append(report.Citations, found...)
		warnings = append(warnings, warned...)
		return nil
	})
	if err != nil {
		return Report{}, nil, err
	}
	return report, warnings, nil
}

func scanFile(path, rel string, opts Options, rev string) ([]Citation, []string, error) {
	f, err := os.Open(path) // #nosec G304 -- the path comes from walking the root the caller named
	if err != nil {
		return nil, nil, err
	}
	defer f.Close() //nolint:errcheck // read-only

	head := make([]byte, 8000)
	n, _ := f.Read(head)
	if bytes.IndexByte(head[:n], 0) >= 0 {
		return nil, nil, nil
	}
	if _, err := f.Seek(0, 0); err != nil {
		return nil, nil, err
	}

	var found []Citation
	var warnings []string
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1<<20)
	for line := 1; scanner.Scan(); line++ {
		text := scanner.Text()
		if !attempt.MatchString(text) {
			continue
		}
		m := citation.FindStringSubmatch(text)
		if m == nil {
			warnings = append(warnings, fmt.Sprintf("%s:%d: %s must be followed by one livt URI and nothing else", rel, line, Marker))
			continue
		}
		if strings.ContainsAny(m[1], "{}") {
			// A template is documentation showing the shape, not a claim about
			// this file. Skipping it silently is what lets the marker be
			// written about in a README without inventing a citation.
			continue
		}
		if _, ok := uri.Parse(m[1]); !ok {
			warnings = append(warnings, fmt.Sprintf("%s:%d: %q is not a livt URI", rel, line, m[1]))
			continue
		}
		found = append(found, Citation{URI: m[1], File: rel, Line: line, URL: lineURL(opts, rev, rel, line)})
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", rel, err)
	}
	return found, warnings, nil
}

func lineURL(opts Options, rev, rel string, line int) string {
	if opts.URLTemplate == "" || opts.Base == "" || rev == "" {
		return ""
	}
	return strings.NewReplacer(
		"{base}", opts.Base,
		"{rev}", rev,
		"{path}", rel,
		"{line}", strconv.Itoa(line),
	).Replace(opts.URLTemplate)
}
