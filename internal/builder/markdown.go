package builder

import (
	"bytes"
	"html/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// md leaves raw HTML escaped: a body is the author's prose, and the page must
// not run whatever markup a repository file happens to carry. Hard wraps keep
// a single newline a line break, as the plain-text rendering did before.
var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithRendererOptions(html.WithHardWraps()),
)

// markdown renders a story's or an opportunity's body for a page.
func markdown(src string) template.HTML {
	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		return template.HTML(template.HTMLEscapeString(src))
	}
	return template.HTML(buf.String()) //nolint:gosec // goldmark escapes raw HTML without WithUnsafe
}
