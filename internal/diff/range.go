package diff

import (
	"fmt"
	"os"
	"strings"

	"github.com/boykush/livt/internal/uri"
)

// revSep is git's own two-dot spelling, and means here what it means to
// `git diff`: the tree at one revision against the tree at another.
const revSep = ".."

// Range is the two revisions a build renders a diff between. An empty Head is
// the working tree — the common case, since a change is usually reviewed from
// the branch that makes it, and it is also what the rest of the site renders.
type Range struct {
	Base string
	Head string
}

// ParseRange reads the --diff argument: "<base>..<head>", or "<base>" alone
// against the working tree.
func ParseRange(s string) (Range, error) {
	base, head, found := strings.Cut(s, revSep)
	if !found {
		if s == "" {
			return Range{}, fmt.Errorf("--diff takes a revision: <base>%s<head>, or <base> alone against the working tree", revSep)
		}
		return Range{Base: s}, nil
	}
	// A three-dot range is git's merge-base spelling, which says something this
	// does not do. Cut would read it as a head beginning with a dot — a ref name
	// git itself forbids — so it is refused by name rather than silently.
	if strings.HasPrefix(head, ".") {
		return Range{}, fmt.Errorf("--diff %q: only the two-dot range %s is understood", s, revSep)
	}
	if base == "" {
		return Range{}, fmt.Errorf("--diff %q: no base revision before %s", s, revSep)
	}
	return Range{Base: base, Head: head}, nil
}

// Result is one build's diff: what changed, tallied by how, and the revisions
// it was taken between.
type Result struct {
	// Base and Head are short hashes as git resolved them, so the page says
	// which revisions were actually read rather than what was typed. Head is
	// empty when the head is the working tree, which has no hash to print.
	Base     string
	Head     string
	Changes  []Change
	Added    int
	Removed  int
	Modified int
}

// Compute reads both revisions and pairs them. The working tree is read first
// and kept whatever the range is: it is what the site being built holds, so it
// is also what decides whether a change has a page to link to.
func (r Range) Compute(root string, dirs Dirs) (*Result, error) {
	if err := requireRepo(root); err != nil {
		return nil, err
	}
	// Read under root like every other side, so the tree git is asked about and
	// the tree the site is built from are the same one.
	site, err := Scan(dirs.under(root))
	if err != nil {
		return nil, err
	}
	base, baseRev, err := scanRev(root, r.Base, dirs)
	if err != nil {
		return nil, err
	}
	head, headRev := site, ""
	if r.Head != "" {
		if head, headRev, err = scanRev(root, r.Head, dirs); err != nil {
			return nil, err
		}
	}
	result := &Result{Base: baseRev, Head: headRev, Changes: Compare(base, head)}
	for i, c := range result.Changes {
		result.Changes[i].Page = page(site, c.URI)
		if parent, held := head.get(c.Parent); held {
			result.Changes[i].ParentTitle = parent.Title
		}
		switch c.Status {
		case StatusAdded:
			result.Added++
		case StatusRemoved:
			result.Removed++
		case StatusModified:
			result.Modified++
		}
	}
	return result, nil
}

// scanRev exports one revision to a temporary directory and snapshots it there.
// The export is removed on the way out: it is read once, and a build should not
// leave a copy of an old revision beside the site it wrote.
func scanRev(root, rev string, dirs Dirs) (*Snapshot, string, error) {
	short, err := resolveRev(root, rev)
	if err != nil {
		return nil, "", err
	}
	dir, err := os.MkdirTemp("", "livt-diff-")
	if err != nil {
		return nil, "", err
	}
	defer os.RemoveAll(dir)
	if err := exportRev(root, rev, dir, dirs.paths()); err != nil {
		return nil, "", err
	}
	snapshot, err := Scan(dirs.under(dir))
	if err != nil {
		return nil, "", fmt.Errorf("read revision %s: %w", short, err)
	}
	return snapshot, short, nil
}

// page is where a change lands on the site being built, empty when the site
// does not hold that URI. That is what leaves a removed URI unlinked without
// the page having to reason about it.
func page(site *Snapshot, u string) string {
	if _, held := site.get(u); !held {
		return ""
	}
	parsed, ok := uri.Parse(u)
	if !ok {
		return ""
	}
	return parsed.Page()
}
