package builder

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/boykush/livt/internal/diff"
	"github.com/boykush/livt/internal/i18n"
)

// foldMin is how many unchanged lines in a row it takes before the diff tucks
// them away: fewer than this read faster than the toggle that would hide them.
const foldMin = 6

// diffMD is md with one more inline node, the span that marks a changed run.
var diffMD = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		renderer.WithNodeRenderers(util.Prioritized(diffSpanRenderer{}, 500)),
	),
)

var kindDiffSpan = ast.NewNodeKind("DiffSpan")

// diffSpan wraps text that changed, in the classes the diff page's own lines
// already use, so a body reads with the same red and green as a rule.
type diffSpan struct {
	ast.BaseInline
	Class string
}

func (n *diffSpan) Kind() ast.NodeKind { return kindDiffSpan }

func (n *diffSpan) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{"Class": n.Class}, nil)
}

type diffSpanRenderer struct{}

func (diffSpanRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(kindDiffSpan, func(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			_, _ = w.WriteString(`<span class="` + n.(*diffSpan).Class + `">`)
		} else {
			_, _ = w.WriteString("</span>")
		}
		return ast.WalkContinue, nil
	})
}

// proseMarks says, per source line of a body, how it changed: the class its
// text takes, and the byte ranges inside it that moved.
type proseMarks struct {
	starts []int
	class  []string
	moved  [][][2]int
}

func newProseMarks(lines []diff.Line) (string, proseMarks) {
	var src strings.Builder
	m := proseMarks{class: make([]string, len(lines)), moved: make([][][2]int, len(lines))}
	for i, l := range lines {
		start := src.Len()
		m.starts = append(m.starts, start)
		src.WriteString(l.Field.Value)
		if i < len(lines)-1 {
			src.WriteByte('\n')
		}
		switch l.Op {
		case diff.OpAdd:
			m.class[i] = "diff-add"
		case diff.OpDel:
			m.class[i] = "diff-del"
		default:
			continue
		}
		if len(l.Parts) == 0 {
			continue
		}
		m.class[i] += " reworded"
		at := start
		for _, p := range l.Parts {
			if p.Changed {
				m.moved[i] = append(m.moved[i], [2]int{at, at + len(p.Text)})
			}
			at += len(p.Text)
		}
	}
	return src.String(), m
}

// lineAt is the source line an offset falls on.
func (m proseMarks) lineAt(offset int) int {
	return sort.Search(len(m.starts), func(i int) bool { return m.starts[i] > offset }) - 1
}

func (m proseMarks) changed() bool {
	for _, c := range m.class {
		if c != "" {
			return true
		}
	}
	return false
}

// diffProse renders a run of body lines as the markdown they are, with each
// changed line's text marked inline rather than its block: a story's three
// lines are one paragraph, and only the line that moved should colour.
func (b *Builder) diffProse(lines []diff.Line, sourceDir string) template.HTML {
	src, marks := newProseMarks(lines)
	source := []byte(src)
	doc := diffMD.Parser().Parse(text.NewReader(source))

	var links []*ast.Link
	var texts []ast.Node
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n := n.(type) {
		case *ast.Link:
			links = append(links, n)
		case *ast.CodeSpan:
			// Its renderer reads its children as plain text, so it is marked
			// whole from outside.
			texts = append(texts, n)
			return ast.WalkSkipChildren, nil
		case *ast.Text:
			texts = append(texts, n)
		}
		return ast.WalkContinue, nil
	})
	for _, l := range links {
		b.resolveLink(l, sourceDir, "")
	}
	for _, n := range texts {
		markText(n, marks)
	}

	var buf bytes.Buffer
	for _, run := range b.proseRuns(doc, marks) {
		var part bytes.Buffer
		for _, n := range run.blocks {
			if err := diffMD.Renderer().Render(&part, source, n); err != nil {
				return template.HTML(template.HTMLEscapeString(src))
			}
		}
		if run.folded == 0 {
			buf.Write(part.Bytes())
			continue
		}
		fmt.Fprintf(&buf, `<details class="diff-fold"><summary>%s</summary>%s</details>`,
			template.HTMLEscapeString(fmt.Sprintf(i18n.Of(b.Lang).Msg("diff.fold"), run.folded)), part.Bytes())
	}
	return template.HTML(buf.String()) //nolint:gosec // goldmark escapes raw HTML without WithUnsafe
}

// markText wraps a text node's changed runs in diffSpans, splitting it where a
// moved range starts or ends.
func markText(n ast.Node, m proseMarks) {
	if cs, ok := n.(*ast.CodeSpan); ok {
		first, _ := cs.FirstChild().(*ast.Text)
		if first == nil {
			return
		}
		if class := m.class[m.lineAt(first.Segment.Start)]; class != "" {
			wrap(cs, &diffSpan{Class: class})
		}
		return
	}
	t := n.(*ast.Text)
	line := m.lineAt(t.Segment.Start)
	if line < 0 || m.class[line] == "" {
		return
	}
	cuts := []int{t.Segment.Start, t.Segment.Stop}
	for _, r := range m.moved[line] {
		for _, c := range r {
			if c > t.Segment.Start && c < t.Segment.Stop {
				cuts = append(cuts, c)
			}
		}
	}
	sort.Ints(cuts)
	span := &diffSpan{Class: m.class[line]}
	var last *ast.Text
	for i := 0; i+1 < len(cuts); i++ {
		if cuts[i] == cuts[i+1] {
			continue
		}
		piece := ast.NewTextSegment(text.NewSegment(cuts[i], cuts[i+1]))
		piece.SetRaw(t.IsRaw())
		last = piece
		if inRanges(cuts[i], m.moved[line]) {
			moved := &diffSpan{Class: "moved"}
			moved.AppendChild(moved, piece)
			span.AppendChild(span, moved)
		} else {
			span.AppendChild(span, piece)
		}
	}
	if last == nil {
		return
	}
	last.SetSoftLineBreak(t.SoftLineBreak())
	last.SetHardLineBreak(t.HardLineBreak())
	t.Parent().ReplaceChild(t.Parent(), t, span)
}

func wrap(n ast.Node, span *diffSpan) {
	parent := n.Parent()
	parent.ReplaceChild(parent, n, span)
	span.AppendChild(span, n)
}

func inRanges(at int, ranges [][2]int) bool {
	for _, r := range ranges {
		if at >= r[0] && at < r[1] {
			return true
		}
	}
	return false
}

// proseRun is a stretch of a body's top-level blocks, folded when it is long
// and none of it changed.
type proseRun struct {
	blocks []ast.Node
	folded int
}

// proseRuns groups the top-level blocks by whether any of their source lines
// changed. A body with no change at all is an addition read whole, and stays
// open.
func (b *Builder) proseRuns(doc ast.Node, m proseMarks) []proseRun {
	var blocks []ast.Node
	var starts []int
	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		if at, ok := firstOffset(n); ok {
			blocks = append(blocks, n)
			starts = append(starts, m.lineAt(at))
		}
	}
	if !m.changed() {
		return []proseRun{{blocks: blocks}}
	}
	var runs []proseRun
	var cur []ast.Node
	curLines, curChanged := 0, false
	flush := func() {
		if len(cur) == 0 {
			return
		}
		run := proseRun{blocks: cur}
		if !curChanged && curLines >= foldMin {
			run.folded = curLines
		}
		runs = append(runs, run)
		cur, curLines = nil, 0
	}
	for i, n := range blocks {
		end := len(m.class)
		if i+1 < len(blocks) {
			end = starts[i+1]
		}
		changed := false
		for l := starts[i]; l < end; l++ {
			changed = changed || m.class[l] != ""
		}
		if changed != curChanged {
			flush()
			curChanged = changed
		}
		cur = append(cur, n)
		curLines += end - starts[i]
	}
	flush()
	return runs
}

// firstOffset is where a block's source starts, from its own lines or, for a
// container, from the first thing inside it that has any.
func firstOffset(n ast.Node) (int, bool) {
	if t, ok := n.(*ast.Text); ok {
		return t.Segment.Start, true
	}
	if n.Type() == ast.TypeBlock && n.Lines().Len() > 0 {
		return n.Lines().At(0).Start, true
	}
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if at, ok := firstOffset(c); ok {
			return at, true
		}
	}
	return 0, false
}
