package docs

import (
	"html"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"unicode"
)

// Where each language's pages are written. content/{lang}/site.md is the text
// every page of that language shows around its own, so it is checked against
// all of them rather than against a page of its own.
var siteDirs = map[string]string{"en": "site", "ja": "site/ja"}

const chrome = "site"

var (
	commentRe  = regexp.MustCompile(`(?s)<!--.*?-->`)
	scriptRe   = regexp.MustCompile(`(?is)<script\b.*?</script>`)
	styleRe    = regexp.MustCompile(`(?is)<style\b.*?</style>`)
	tagRe      = regexp.MustCompile(`(?s)<[^>]+>`)
	linkRe     = regexp.MustCompile(`\[([^\]]*)\]\(([^)\s]+)\)`)
	titleRe    = regexp.MustCompile(`(?is)<title>(.*?)</title>`)
	descRe     = regexp.MustCompile(`(?is)<meta\s+name="description"\s+content="([^"]*)"`)
	redirectRe = regexp.MustCompile(`(?i)http-equiv="refresh"`)
)

func TestSiteShowsEveryPassageOfItsText(t *testing.T) {
	for lang, dir := range siteDirs {
		texts, err := filepath.Glob(filepath.Join("content", lang, "*.md"))
		if err != nil {
			t.Fatal(err)
		}
		var pages []string
		for _, text := range texts {
			name := strings.TrimSuffix(filepath.Base(text), ".md")
			if name == chrome {
				continue
			}
			page := filepath.Join(dir, name+".html")
			pages = append(pages, page)
			checkPage(t, text, page)
		}
		for _, page := range pages {
			checkPassages(t, filepath.Join("content", lang, chrome+".md"), page)
		}
	}
}

func TestEveryPageHasItsText(t *testing.T) {
	for lang, dir := range siteDirs {
		pages, err := filepath.Glob(filepath.Join(dir, "*.html"))
		if err != nil {
			t.Fatal(err)
		}
		for _, page := range pages {
			if redirectRe.Match(read(t, page)) {
				continue
			}
			text := filepath.Join("content", lang, strings.TrimSuffix(filepath.Base(page), ".html")+".md")
			if _, err := os.Stat(text); err != nil {
				t.Errorf("%s has no text in %s: write the page there, then render it", page, text)
			}
		}
	}
}

func checkPage(t *testing.T, text, page string) {
	t.Helper()
	if _, err := os.Stat(page); err != nil {
		t.Errorf("%s has no page at %s: render it with the render-site skill", text, page)
		return
	}
	front, _ := splitFrontMatter(string(read(t, text)))
	doc := string(read(t, page))
	if want := front["title"]; want != "" && firstGroup(titleRe, doc) != want {
		t.Errorf("%s: <title> is %q, but %s says %q", page, firstGroup(titleRe, doc), text, want)
	}
	if want := front["description"]; want != "" && firstGroup(descRe, doc) != want {
		t.Errorf("%s: the description meta is %q, but %s says %q", page, firstGroup(descRe, doc), text, want)
	}
	checkPassages(t, text, page)
}

// checkPassages fails for each heading, paragraph or list item of the text, and
// each link, that the page does not show. Whitespace and punctuation are left
// out of the comparison, so the page is free to split a passage across
// elements, as a sticky, a canvas box or a label beside its answer do.
func checkPassages(t *testing.T, text, page string) {
	t.Helper()
	_, body := splitFrontMatter(string(read(t, text)))
	doc := string(read(t, page))
	shown := normalize(visibleText(doc))
	for _, passage := range passages(body) {
		for _, link := range linkRe.FindAllStringSubmatch(passage, -1) {
			if !strings.Contains(doc, `href="`+link[2]+`"`) {
				t.Errorf("%s does not link %s, which %s links", page, link[2], text)
			}
		}
		plain := strings.NewReplacer("**", "", "`", "").Replace(linkRe.ReplaceAllString(passage, "$1"))
		if !strings.Contains(shown, normalize(plain)) {
			t.Errorf("%s does not show this passage of %s; render the page again:\n\t%s", page, text, plain)
		}
	}
}

// passages splits Markdown into the units a page shows as one run of text:
// a heading, a list item, or a paragraph. Comments are notes to whoever
// renders the page, and are not shown.
func passages(body string) []string {
	var out, para []string
	flush := func() {
		if len(para) > 0 {
			out = append(out, strings.Join(para, " "))
			para = nil
		}
	}
	for _, line := range strings.Split(commentRe.ReplaceAllString(body, ""), "\n") {
		line = strings.TrimSpace(line)
		switch {
		case line == "":
			flush()
		case strings.HasPrefix(line, "#"):
			flush()
			out = append(out, strings.TrimSpace(strings.TrimLeft(line, "#")))
		case strings.HasPrefix(line, "- "):
			flush()
			out = append(out, strings.TrimPrefix(line, "- "))
		default:
			para = append(para, line)
		}
	}
	flush()
	return out
}

func splitFrontMatter(s string) (map[string]string, string) {
	front := map[string]string{}
	if !strings.HasPrefix(s, "---\n") {
		return front, s
	}
	end := strings.Index(s[4:], "\n---\n")
	if end < 0 {
		return front, s
	}
	for _, line := range strings.Split(s[4:4+end], "\n") {
		if key, value, ok := strings.Cut(line, ":"); ok {
			front[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
	}
	return front, s[4+end+5:]
}

func visibleText(doc string) string {
	for _, re := range []*regexp.Regexp{commentRe, scriptRe, styleRe} {
		doc = re.ReplaceAllString(doc, "")
	}
	return html.UnescapeString(tagRe.ReplaceAllString(doc, ""))
}

func normalize(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			return -1
		}
		return r
	}, s)
}

func firstGroup(re *regexp.Regexp, s string) string {
	if m := re.FindStringSubmatch(s); m != nil {
		return html.UnescapeString(strings.TrimSpace(m[1]))
	}
	return ""
}

func read(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
