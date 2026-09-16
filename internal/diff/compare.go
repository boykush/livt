package diff

import (
	"sort"

	"github.com/boykush/livt/internal/uri"
)

// Op is what one line did between the two revisions, spelled the way git diff
// spells it so the page can be read without a legend.
type Op string

const (
	OpContext Op = " "
	OpAdd     Op = "+"
	OpDel     Op = "-"
)

// Line is one line of an entry's diff.
type Line struct {
	Op   Op
	Text string
}

// Status is what happened to a URI. There are only three: a URI the head does
// not hold, one the base did not, and one both hold differently.
type Status string

const (
	StatusAdded    Status = "added"
	StatusRemoved  Status = "removed"
	StatusModified Status = "modified"
)

// Change is one URI's diff. Page is where it lands on the site being built,
// empty when the site does not hold it — a removed URI has nowhere to go, and a
// link to a page that was never built is worse than no link.
type Change struct {
	URI    string
	Kind   uri.Kind
	Parent string
	Title  string
	Status Status
	Lines  []Line
	Page   string
}

// Compare pairs the two snapshots by URI. A URI both hold with the same lines
// is not a change and is left out: the page is what moved, not the repository.
func Compare(base, head *Snapshot) []Change {
	var placed []placedChange
	for i, e := range head.Entries {
		before, held := base.get(e.URI)
		switch {
		case !held:
			placed = append(placed, placedChange{anchor: i, change: change(e, StatusAdded, added(e.Lines))})
		case !sameLines(before.Lines, e.Lines):
			placed = append(placed, placedChange{anchor: i, change: change(e, StatusModified, diffLines(before.Lines, e.Lines))})
		}
	}
	placed = append(placed, removals(base, head)...)
	sort.SliceStable(placed, func(i, j int) bool {
		if placed[i].anchor != placed[j].anchor {
			return placed[i].anchor < placed[j].anchor
		}
		return placed[i].after < placed[j].after
	})
	changes := make([]Change, 0, len(placed))
	for _, p := range placed {
		changes = append(changes, p.change)
	}
	return changes
}

// placedChange carries where a change belongs in the head's reading order, so
// the diff lists what changed in the order the livt repository reads today.
type placedChange struct {
	change Change
	anchor int
	after  int
}

// removals anchor to the last entry that survived before them, so a deleted
// rule is listed where it used to sit rather than in a heap at the end.
func removals(base, head *Snapshot) []placedChange {
	var placed []placedChange
	anchor, after := -1, 0
	for _, e := range base.Entries {
		if i, held := head.index(e.URI); held {
			anchor, after = i, 0
			continue
		}
		after++
		placed = append(placed, placedChange{anchor: anchor, after: after, change: change(e, StatusRemoved, removed(e.Lines))})
	}
	return placed
}

func change(e Entry, status Status, lines []Line) Change {
	return Change{URI: e.URI, Kind: e.Kind, Parent: e.Parent, Title: e.Title, Status: status, Lines: lines}
}

func added(lines []string) []Line   { return ops(lines, OpAdd) }
func removed(lines []string) []Line { return ops(lines, OpDel) }

func ops(lines []string, op Op) []Line {
	out := make([]Line, 0, len(lines))
	for _, l := range lines {
		out = append(out, Line{Op: op, Text: l})
	}
	return out
}

func sameLines(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// diffLines is the line diff git diff draws: the longest common subsequence
// stays as context and everything else is a removal beside its replacement.
// An entry holds its own fields and nothing else, so the quadratic table is
// over a handful of lines.
func diffLines(a, b []string) []Line {
	lcs := make([][]int, len(a)+1)
	for i := range lcs {
		lcs[i] = make([]int, len(b)+1)
	}
	for i := len(a) - 1; i >= 0; i-- {
		for j := len(b) - 1; j >= 0; j-- {
			if a[i] == b[j] {
				lcs[i][j] = lcs[i+1][j+1] + 1
				continue
			}
			lcs[i][j] = max(lcs[i+1][j], lcs[i][j+1])
		}
	}
	var lines []Line
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] == b[j]:
			lines = append(lines, Line{Op: OpContext, Text: a[i]})
			i, j = i+1, j+1
		case lcs[i+1][j] >= lcs[i][j+1]:
			lines = append(lines, Line{Op: OpDel, Text: a[i]})
			i++
		default:
			lines = append(lines, Line{Op: OpAdd, Text: b[j]})
			j++
		}
	}
	lines = append(lines, removed(a[i:])...)
	return append(lines, added(b[j:])...)
}
