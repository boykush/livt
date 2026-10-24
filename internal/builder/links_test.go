package builder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// linkRepo lays out a livt repository under one root, the way cmd/build does,
// holding a term, a scoped term, a policy, a story, and a mapping to link to.
func linkRepo(t *testing.T) Builder {
	t.Helper()
	root := t.TempDir()
	b := Builder{
		Root:             root,
		OpportunitiesDir: filepath.Join(root, "opportunities"),
		CanvasesDir:      filepath.Join(root, "discoveries", "opportunity-canvases"),
		MappingsDir:      filepath.Join(root, "discoveries", "example-mappings"),
		StoriesDir:       filepath.Join(root, "stories"),
		USMDir:           filepath.Join(root, "discoveries", "usm"),
		UbiquitousDir:    filepath.Join(root, "design", "ubiquitous"),
		PoliciesDir:      filepath.Join(root, "design", "policies"),
		OutDir:           t.TempDir(),
	}
	writeFile(t, filepath.Join(b.UbiquitousDir, "opportunity.md"), "---\nname: オポチュニティ\n---\n\n機会。\n")
	writeFile(t, filepath.Join(b.UbiquitousDir, "billing", "invoice.md"), "---\nname: 請求書\n---\n\n請求。\n")
	writeFile(t, filepath.Join(b.PoliciesDir, "lists-filter.md"), "---\nname: 一覧は絞り込める\n---\n\n[](livt://ubiquitous/opportunity)で絞り込める。\n")
	writeFile(t, filepath.Join(b.PoliciesDir, "design-system", "tiles.md"), "---\nname: タイルで並ぶ\n---\n\n本文。\n")
	writeFile(t, filepath.Join(b.StoriesDir, "checkout.md"), "---\nname: 購入する\n---\n")
	writeFile(t, filepath.Join(b.MappingsDir, "checkout.yaml"), "rules:\n  - id: R-01\n    name: 在庫がなければ買えない\n    examples:\n      - id: EX-01\n        name: 在庫0なら買えない\n")
	return b
}

func render(b Builder, src, sourceDir, prefix string) string {
	return string(b.markdown(src, sourceDir, prefix))
}

// livt:automates livt://mapping/link-artifacts-from-body-text/rule/R-01/example/EX-01
// A livt URI written as a link's destination becomes a link to its row.
func TestMarkdownLinksALivtURIToItsPage(t *testing.T) {
	b := linkRepo(t)
	got := render(b, "[オポチュニティ](livt://ubiquitous/opportunity)", b.UbiquitousDir, "")
	if want := `<a href="ubiquitous.html#opportunity">オポチュニティ</a>`; !strings.Contains(got, want) {
		t.Errorf("got %s, want it to contain %s", got, want)
	}
}

// livt:automates livt://mapping/link-artifacts-from-body-text/rule/R-01/example/EX-02
// A rule's URI lands on its sticky, written relative to the page the body is on.
func TestMarkdownLinksARuleURIToItsSticky(t *testing.T) {
	b := linkRepo(t)
	got := render(b, "[ルール](livt://mapping/checkout/rule/R-01)", b.StoriesDir, "../")
	if want := `<a href="../mapping/checkout.html#rule-R-01">ルール</a>`; !strings.Contains(got, want) {
		t.Errorf("got %s, want it to contain %s", got, want)
	}
}

// livt:automates livt://mapping/link-artifacts-from-body-text/rule/R-02/example/EX-01
// livt:automates livt://mapping/link-artifacts-from-body-text/rule/R-02/example/EX-02
// An empty link text is filled with the name of what the link points at.
func TestMarkdownNamesAnEmptyLinkAfterItsTarget(t *testing.T) {
	b := linkRepo(t)
	cases := map[string]string{
		"[](livt://ubiquitous/opportunity)":                   `<a href="ubiquitous.html#opportunity">オポチュニティ</a>`,
		"[](livt://ubiquitous/billing/invoice)":               `<a href="ubiquitous.html#billing/invoice">請求書</a>`,
		"[](livt://mapping/checkout/rule/R-01)":               `<a href="mapping/checkout.html#rule-R-01">在庫がなければ買えない</a>`,
		"[](livt://policy/design-system/tiles)":               `<a href="policies.html#design-system/tiles">タイルで並ぶ</a>`,
		"[](livt://story/checkout)":                           `<a href="story/checkout.html">購入する</a>`,
		"[](livt://mapping/checkout/rule/R-01/example/EX-01)": `<a href="mapping/checkout.html#rule-R-01-example-EX-01">在庫0なら買えない</a>`,
	}
	for src, want := range cases {
		if got := render(b, src, b.UbiquitousDir, ""); !strings.Contains(got, want) {
			t.Errorf("%s rendered %s, want it to contain %s", src, got, want)
		}
	}
}

// livt:automates livt://mapping/link-artifacts-from-body-text/rule/R-03/example/EX-01
// A path from the livt repository's root lands where the term's URI does.
func TestMarkdownLinksARootPathLikeItsURI(t *testing.T) {
	b := linkRepo(t)
	got := render(b, "[オポチュニティ](/design/ubiquitous/opportunity.md)", b.PoliciesDir, "")
	if want := `<a href="ubiquitous.html#opportunity">オポチュニティ</a>`; !strings.Contains(got, want) {
		t.Errorf("got %s, want it to contain %s", got, want)
	}
}

// livt:automates livt://mapping/link-artifacts-from-body-text/rule/R-03/example/EX-02
// A path relative to the file the body is in lands there too.
func TestMarkdownLinksARelativePathLikeItsURI(t *testing.T) {
	b := linkRepo(t)
	got := render(b, "[オポチュニティ](../ubiquitous/opportunity.md)", b.PoliciesDir, "")
	if want := `<a href="ubiquitous.html#opportunity">オポチュニティ</a>`; !strings.Contains(got, want) {
		t.Errorf("got %s, want it to contain %s", got, want)
	}
	got = render(b, "[請求書](../../ubiquitous/billing/invoice.md)", filepath.Join(b.PoliciesDir, "design-system"), "")
	if want := `<a href="ubiquitous.html#billing/invoice">請求書</a>`; !strings.Contains(got, want) {
		t.Errorf("got %s, want it to contain %s", got, want)
	}
}

// livt:automates livt://mapping/link-artifacts-from-body-text/rule/R-04/example/EX-01
// livt:automates livt://mapping/link-artifacts-from-body-text/rule/R-04/example/EX-02
// A reference the livt repository cannot follow keeps its text and loses the link.
func TestMarkdownUnwrapsAReferenceItCannotFollow(t *testing.T) {
	b := linkRepo(t)
	cases := map[string]string{
		"まだない[用語](livt://ubiquitous/not-yet)だ":     "<p>まだない用語だ</p>",
		"[](livt://ubiquitous/not-yet)":            "<p>livt://ubiquitous/not-yet</p>",
		"[ルール](livt://mapping/checkout/rule/R-09)": "<p>ルール</p>",
		"[**太字**](/design/ubiquitous/not-yet.md)":  "<p><strong>太字</strong></p>",
		"[README](/README.md)":                     "<p>README</p>",
	}
	for src, want := range cases {
		if got := render(b, src, b.UbiquitousDir, ""); !strings.Contains(got, want) || strings.Contains(got, "<a ") {
			t.Errorf("%s rendered %s, want %s with no link", src, got, want)
		}
	}
}

// livt:automates livt://mapping/link-artifacts-from-body-text/rule/R-04/example/EX-03
// A URL outside the livt repository, or an anchor on the page, is left as written.
func TestMarkdownLeavesOtherLinksAsWritten(t *testing.T) {
	b := linkRepo(t)
	cases := map[string]string{
		"[livt](https://example.com/a.md)": `<a href="https://example.com/a.md">livt</a>`,
		"[mail](mailto:a@example.com)":     `<a href="mailto:a@example.com">mail</a>`,
		"[上へ](#top)":                       `<a href="#top">上へ</a>`,
	}
	for src, want := range cases {
		if got := render(b, src, b.UbiquitousDir, ""); !strings.Contains(got, want) {
			t.Errorf("%s rendered %s, want it to contain %s", src, got, want)
		}
	}
}

// livt:automates livt://mapping/link-artifacts-from-body-text/rule/R-05/example/EX-01
// A term's definition is Markdown, and the links in it resolve.
func TestBuildGlossaryRendersDefinitionsAsMarkdown(t *testing.T) {
	b := linkRepo(t)
	writeFile(t, filepath.Join(b.UbiquitousDir, "story.md"), "---\nname: ストーリー\n---\n\n**価値**の単位。[](livt://ubiquitous/opportunity)に紐づく。\n")
	if err := b.buildGlossary(); err != nil {
		t.Fatal(err)
	}
	out, err := os.ReadFile(filepath.Join(b.OutDir, "ubiquitous.html"))
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)
	for _, want := range []string{"<strong>価値</strong>", `<a href="ubiquitous.html#opportunity">オポチュニティ</a>`} {
		if !strings.Contains(html, want) {
			t.Errorf("glossary lacks %s", want)
		}
	}
}

// livt:automates livt://mapping/link-artifacts-from-body-text/rule/R-06/example/EX-01
// livt:automates livt://mapping/link-artifacts-from-body-text/rule/R-06/example/EX-02
// Every policy in force is a row anchored by its reference, a scoped one shown
// with its context, and its body resolves links like any other.
func TestBuildPoliciesRendersAnchoredRows(t *testing.T) {
	b := linkRepo(t)
	writeFile(t, filepath.Join(b.PoliciesDir, "withdrawn.md"), "---\nname: 退役した決まり\nstatus: retired\n---\n")
	if err := b.buildPolicies(); err != nil {
		t.Fatal(err)
	}
	out, err := os.ReadFile(filepath.Join(b.OutDir, "policies.html"))
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)
	for _, want := range []string{
		`id="lists-filter"`,
		`id="design-system/tiles"`,
		">design-system</span>",
		`<a href="ubiquitous.html#opportunity">オポチュニティ</a>`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("policies page lacks %s", want)
		}
	}
	if strings.Contains(html, "退役した決まり") {
		t.Error("a retired policy is still on the table")
	}
}
