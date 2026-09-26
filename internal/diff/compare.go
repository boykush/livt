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

// Line is one line of an entry's diff: the field it renders, and — when the
// line is one half of a rewording — that field's value broken into what stayed
// and what moved. Most of what a livt repository holds is prose, and two long
// sentences side by side leave the reader to find the difference by eye.
type Line struct {
	Op    Op
	Field Field
	Parts []Part
}

// Part is a run of a line's value, marked as changed or not.
type Part struct {
	Changed bool
	Text    string
}

// Status is what happened to a URI's record: a URI the head does not hold, one
// the base did not, and one both hold differently. It is not what happened to
// the decision — for that, see Became.
type Status string

const (
	StatusAdded    Status = "added"
	StatusRemoved  Status = "removed"
	StatusModified Status = "modified"
)

// Became is what the change did to the item as a decision, which is the reading the
// site shows. It is not Status: retiring a rule keeps its entry and edits one
// field, deleting takes the entry away, and to a reader both mean the rule
// stopped being asked for. Deletion is no category of its own — the ID contract
// forbids it, so it is an accident rather than a kind of change, and a shelf of
// its own beside the two that are meant to happen would say otherwise.
type Became string

const (
	BecameAdded     Became = "added"
	BecameChanged   Became = "changed"
	BecameWithdrawn Became = "withdrawn"
)

// became reads the two revisions' standing, not their records. An item that was
// agreed and is not was withdrawn, however that was written down.
func became(wasSpec, isSpec bool) Became {
	switch {
	case wasSpec && !isSpec:
		return BecameWithdrawn
	case !wasSpec && isSpec:
		return BecameAdded
	}
	return BecameChanged
}

// Change is one URI's diff. Page is where it lands on the site being built,
// empty when the site does not hold it — a removed URI has nowhere to go, and a
// link to a page that was never built is worse than no link.
type Change struct {
	URI    string
	Kind   uri.Kind
	Parent string
	// ParentTitle is what the parent is called, carried on the child because a
	// mapping heading its changed rules may have no change of its own to be
	// named by.
	ParentTitle string
	Title       string
	Status      Status
	// Became is what the change did to the item as a decision. Every surface that
	// names a change to a reader names this one, so that one sticky wears one
	// word: a retired rule is a modified record, and what happened to it is
	// that it was withdrawn.
	Became Became
	Lines  []Line
	Page   string
}

// Compare pairs the two snapshots by URI. A URI both hold with the same fields
// is not a change and is left out: the page is what moved, not the repository.
func Compare(base, head *Snapshot) []Change {
	var placed []placedChange
	for i, e := range head.Entries {
		before, held := base.get(e.URI)
		switch {
		case !held:
			placed = append(placed, placedChange{anchor: i, change: change(e, StatusAdded, became(false, e.Live), added(e.Fields))})
		case !sameFields(before.Fields, e.Fields):
			b := became(before.Live, e.Live)
			placed = append(placed, placedChange{anchor: i, change: change(e, StatusModified, b, lines(b, before, e))})
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
		placed = append(placed, placedChange{anchor: anchor, after: after, change: change(e, StatusRemoved, BecameWithdrawn, withdrawn(e, Entry{}))})
	}
	return placed
}

func change(e Entry, status Status, became Became, lines []Line) Change {
	return Change{URI: e.URI, Kind: e.Kind, Parent: e.Parent, Title: e.Title, Status: status, Became: became, Lines: lines}
}

// lines draws the change. A withdrawal is drawn as the removal it is, whichever
// way the file recorded it: retiring adds the line that retires, and read as an
// addition it would say the opposite of what happened.
func lines(b Became, before, after Entry) []Line {
	if b == BecameWithdrawn {
		return withdrawn(before, after)
	}
	return diffFields(before.Fields, after.Fields)
}

// withdrawn is the statement that stopped holding, and what replaced it. The
// line recording the withdrawal is left out — the change says that once already,
// and saying it again as an addition is what made a retirement look like the
// opposite of a deletion when both are the same thing to a reader.
func withdrawn(before, after Entry) []Line {
	var out []Line
	for _, f := range before.Fields {
		if f.Label == "" {
			out = append(out, Line{Op: OpDel, Field: f})
		}
	}
	for _, f := range after.Fields {
		if f.Label == LabelSupersededBy {
			out = append(out, Line{Op: OpContext, Field: f})
		}
	}
	return out
}

func added(fields []Field) []Line   { return ops(fields, OpAdd) }
func removed(fields []Field) []Line { return ops(fields, OpDel) }

func ops(fields []Field, op Op) []Line {
	out := make([]Line, 0, len(fields))
	for _, f := range fields {
		out = append(out, Line{Op: op, Field: f})
	}
	return out
}

func sameFields(a, b []Field) bool {
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

// diffFields is the line diff git diff draws — the longest common subsequence
// stays as context and everything else is a removal beside its replacement —
// with each reworded pair then broken down further. An entry holds its own
// fields and nothing else, so the quadratic table is over a handful of lines.
func diffFields(a, b []Field) []Line {
	lcs := lcsTable(len(a), len(b), func(i, j int) bool { return a[i] == b[j] })
	var lines []Line
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] == b[j]:
			lines = append(lines, Line{Op: OpContext, Field: a[i]})
			i, j = i+1, j+1
		case lcs[i+1][j] >= lcs[i][j+1]:
			lines = append(lines, Line{Op: OpDel, Field: a[i]})
			i++
		default:
			lines = append(lines, Line{Op: OpAdd, Field: b[j]})
			j++
		}
	}
	lines = append(lines, removed(a[i:])...)
	lines = append(lines, added(b[j:])...)
	return markReworded(lines)
}

// markReworded pairs each removal with the addition that replaced it and breaks
// both down to what actually moved. Only a pair naming the same field is a
// rewording; two different fields are two changes, and marking them up against
// each other would invent a relation the record does not have.
func markReworded(lines []Line) []Line {
	for i := 0; i+1 < len(lines); i++ {
		del, add := lines[i], lines[i+1]
		if del.Op != OpDel || add.Op != OpAdd || del.Field.Label != add.Field.Label {
			continue
		}
		delParts, addParts, ok := parts(del.Field.Value, add.Field.Value)
		if !ok {
			continue
		}
		lines[i].Parts, lines[i+1].Parts = delParts, addParts
		i++
	}
	return lines
}

// similarEnough is the share of a line that has to survive a rewording for the
// breakdown to help. Below it the two are different sentences rather than one
// sentence edited, and marking up the few characters they happen to share
// scatters highlights through both.
const similarEnough = 0.4

// parts breaks a reworded pair down to the runs that stayed and the runs that
// moved. It works in runes: most of what a livt repository holds is prose, and
// for a language that does not space its words the character is the only unit
// the diff can honestly claim to see.
func parts(before, after string) (delParts, addParts []Part, ok bool) {
	a, b := []rune(before), []rune(after)
	lcs := lcsTable(len(a), len(b), func(i, j int) bool { return a[i] == b[j] })
	common := lcs[0][0]
	if common == 0 || float64(common) < similarEnough*float64(max(len(a), len(b))) {
		return nil, nil, false
	}
	var del, add []Part
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] == b[j]:
			del, add = appendPart(del, false, a[i]), appendPart(add, false, b[j])
			i, j = i+1, j+1
		case lcs[i+1][j] >= lcs[i][j+1]:
			del = appendPart(del, true, a[i])
			i++
		default:
			add = appendPart(add, true, b[j])
			j++
		}
	}
	for ; i < len(a); i++ {
		del = appendPart(del, true, a[i])
	}
	for ; j < len(b); j++ {
		add = appendPart(add, true, b[j])
	}
	return del, add, true
}

// appendPart grows the run being built rather than starting a new one, so a
// changed phrase is highlighted once instead of character by character.
func appendPart(parts []Part, changed bool, r rune) []Part {
	if n := len(parts); n > 0 && parts[n-1].Changed == changed {
		parts[n-1].Text += string(r)
		return parts
	}
	return append(parts, Part{Changed: changed, Text: string(r)})
}

// lcsTable is the longest-common-subsequence table both diffs walk, over lines
// or over runes. [i][j] is the length of the longest subsequence shared by the
// two tails starting there.
func lcsTable(n, m int, equal func(i, j int) bool) [][]int {
	table := make([][]int, n+1)
	for i := range table {
		table[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if equal(i, j) {
				table[i][j] = table[i+1][j+1] + 1
				continue
			}
			table[i][j] = max(table[i+1][j], table[i][j+1])
		}
	}
	return table
}
