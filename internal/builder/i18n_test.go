package builder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/boykush/livt/internal/i18n"
)

// filledBuilder writes a small livt repository holding one of everything, so a
// build over it reaches every page and the branches that only appear once a
// resource has content — an automated rule, a question, a scoped term.
func filledBuilder(t *testing.T) Builder {
	t.Helper()
	b := emptyDirsBuilder(t)
	writeFile(t, filepath.Join(b.StoriesDir, "checkout.md"),
		"---\n"+
			"name: カートの中身を確認して注文する\n"+
			"issue: https://github.com/boykush/livt/issues/1\n"+
			"---\n\n"+
			"買い物客として\n注文を確定したい\n")
	writeFile(t, filepath.Join(b.MappingsDir, "checkout.yaml"),
		"rules:\n"+
			"  - id: R-01\n"+
			"    name: 在庫のない商品は注文できない\n"+
			"    examples:\n"+
			"      - id: EX-01\n"+
			"        name: 在庫0の商品はカートに入らない\n"+
			"    automated: true\n"+
			"    issues:\n"+
			"      - https://github.com/boykush/livt/issues/2\n"+
			"  - id: R-02\n"+
			"    name: 送料は購入金額で決まる\n"+
			"questions:\n"+
			"  - id: Q-01\n"+
			"    text: 予約商品はどう扱うか\n"+
			"ubiquitous:\n"+
			"  - cart\n")
	writeFile(t, filepath.Join(b.USMDir, "shopping.yaml"),
		"name: 買い物ジャーニー\n"+
			"ubiquitous:\n"+
			"  - cart\n"+
			"releases:\n"+
			"  - id: v1\n"+
			"    name: 最初のリリース\n"+
			"activities:\n"+
			"  - key: browse\n"+
			"    name: 商品を探す\n"+
			"    steps:\n"+
			"      - key: search\n"+
			"        name: 検索する\n"+
			"        stories:\n"+
			"          - key: checkout\n"+
			"            name: カートの中身を確認して注文する\n"+
			"            release: v1\n")
	writeFile(t, filepath.Join(b.OpportunitiesDir, "shopping.md"), "---\nname: 買い物体験\nissue: https://github.com/boykush/livt/issues/3\n---\n\n買い物客が迷わず注文まで進めない。\n")
	writeFile(t, filepath.Join(b.CanvasesDir, "shopping.yaml"), demoCanvasYAML)
	writeFile(t, filepath.Join(b.UbiquitousDir, "cart.md"), "---\nname: カート\n---\n\n購入予定の商品を貯めておく入れ物。\n")
	writeFile(t, filepath.Join(b.UbiquitousDir, "order", "cart.md"), "---\nname: 注文カート\n---\n\n注文の対象として確定したカート。\n")
	return b
}

// pages is every page a filled build produces.
var pages = []string{
	"index.html",
	"opportunities.html",
	"story-maps.html",
	"stories.html",
	"ubiquitous.html",
	"tasks.html",
	filepath.Join("opportunity", "shopping.html"),
	filepath.Join("opportunity-canvas", "shopping.html"),
	filepath.Join("story", "checkout.html"),
	filepath.Join("mapping", "checkout.html"),
	filepath.Join("story-map", "買い物ジャーニー.html"),
}

// This is what lets Catalog.T treat an unknown key as an error: every page is
// rendered in every language here, so a message id a template invented — or one
// a catalog never gained — fails before a user's build is what finds it. The
// empty repository is built too, since the empty states sit on the far side of
// branches a filled one never takes.
func TestBuildRendersEveryPageInEveryLanguage(t *testing.T) {
	for _, lang := range i18n.Langs {
		t.Run(string(lang), func(t *testing.T) {
			b := filledBuilder(t)
			b.Lang = lang
			if err := b.Build(); err != nil {
				t.Fatal(err)
			}
			code := i18n.Of(lang)["lang.code"]
			for _, page := range pages {
				html := readPage(t, b, page)
				if want := `<html lang="` + code + `">`; !strings.Contains(html, want) {
					t.Errorf("%s does not declare %s", page, want)
				}
			}

			empty := emptyDirsBuilder(t)
			empty.Lang = lang
			if err := empty.Build(); err != nil {
				t.Fatalf("empty repository: %v", err)
			}
		})
	}
}

// The chrome is translated; the livt repository's own prose is not. A team
// writing its stories in Japanese and reading the site in English is the same
// arrangement as the reverse, and neither language setting may rewrite what the
// repository says.
func TestBuildTranslatesTheChromeAndLeavesTheProseAlone(t *testing.T) {
	b := filledBuilder(t)
	b.Lang = i18n.Ja
	if err := b.Build(); err != nil {
		t.Fatal(err)
	}

	for page, wants := range map[string][]string{
		"index.html":                              {"実例マッピング"},
		"story-maps.html":                         {"ストーリーマップ"},
		"stories.html":                            {"オポチュニティ", "すべて"},
		"ubiquitous.html":                         {"用語", "キー", "定義", "文脈"},
		"tasks.html":                              {"未解決の疑問", "会話で閉じる", "未自動化のルール", "テストで閉じる"},
		filepath.Join("story", "checkout.html"):   {"メタデータ", "関連", "説明"},
		filepath.Join("mapping", "checkout.html"): {"ルール", "具体例", "疑問点", "✓ 自動化済み", "ユビキタス言語"},
		filepath.Join("story-map", "買い物ジャーニー.html"): {"オポチュニティ", "アクティビティ", "ユーザータスク", "ストーリー"},
		"opportunities.html":                                 {"オポチュニティ"},
		filepath.Join("opportunity", "shopping.html"):        {"メタデータ", "オポチュニティ"},
		filepath.Join("opportunity-canvas", "shopping.html"): {"オポチュニティキャンバス", "確かめられる事実", "価値についての仮説", "解決策のアイデア", "未記入"},
	} {
		html := readPage(t, b, page)
		for _, want := range wants {
			if !strings.Contains(html, want) {
				t.Errorf("%s is missing the Japanese label %q", page, want)
			}
		}
	}

	// The story's own words, rendered as written whichever language the chrome is in.
	story := readPage(t, b, filepath.Join("story", "checkout.html"))
	if !strings.Contains(story, "カートの中身を確認して注文する") {
		t.Error("the story lost its name")
	}
}

// A Builder that was never handed a language still renders — the zero value is
// the default, not a broken site.
func TestBuildFallsBackToTheDefaultLanguage(t *testing.T) {
	b := filledBuilder(t)
	if err := b.Build(); err != nil {
		t.Fatal(err)
	}
	html := readPage(t, b, "index.html")
	if want := `<html lang="` + i18n.Of(i18n.Default)["lang.code"] + `">`; !strings.Contains(html, want) {
		t.Fatalf("index.html does not declare %s", want)
	}
	if !strings.Contains(html, "Example Mappings") {
		t.Fatal("index.html is not in the default language")
	}
}

func readPage(t *testing.T, b Builder, page string) string {
	t.Helper()
	out, err := os.ReadFile(filepath.Join(b.OutDir, page))
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}
