package builder

import (
	"bytes"
	"html/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

// md leaves raw HTML escaped: a body is the author's prose, and the page must
// not run whatever markup a repository file happens to carry. Hard wraps keep
// a single newline a line break, as the plain-text rendering did before.
var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithRendererOptions(html.WithHardWraps()),
)

// markdown renders a body for a page. sourceDir is the directory of the file
// the body was read from, and prefix the way back to the output root from the
// page it lands on: a link to another point of the livt repository resolves
// against the one and is written relative to the other.
func (b *Builder) markdown(src, sourceDir, prefix string) template.HTML {
	source := []byte(src)
	doc := md.Parser().Parse(text.NewReader(source))

	var links []*ast.Link
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if l, ok := n.(*ast.Link); ok && entering {
			links = append(links, l)
		}
		return ast.WalkContinue, nil
	})
	for _, l := range links {
		b.resolveLink(l, sourceDir, prefix)
	}

	var buf bytes.Buffer
	if err := md.Renderer().Render(&buf, source, doc); err != nil {
		return template.HTML(template.HTMLEscapeString(src))
	}
	return template.HTML(buf.String()) //nolint:gosec // goldmark escapes raw HTML without WithUnsafe
}

// resolveLink points a link at the page its livt URI or file path lands on,
// naming it after its target when the author left the text empty. A link the
// livt repository cannot follow is unwrapped to its text, so a reader is never
// sent to a page that was not built.
func (b *Builder) resolveLink(l *ast.Link, sourceDir, prefix string) {
	dest := string(l.Destination)
	p, ok, ours := b.bodyLinkURI(dest, sourceDir)
	if !ours {
		return
	}
	name, held := "", false
	if ok {
		name, held = b.linkName(p)
	}
	if held && p.Page() != "" {
		l.Destination = []byte(prefix + p.Page())
		if l.ChildCount() == 0 {
			l.AppendChild(l, ast.NewString([]byte(name)))
		}
		return
	}

	parent := l.Parent()
	if l.ChildCount() == 0 {
		parent.InsertBefore(parent, l, ast.NewString([]byte(dest)))
	}
	for c := l.FirstChild(); c != nil; {
		next := c.NextSibling()
		parent.InsertBefore(parent, l, c)
		c = next
	}
	parent.RemoveChild(parent, l)
}
