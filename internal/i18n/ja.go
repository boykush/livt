package i18n

// ja mirrors en key for key, in the words the Japanese practice uses: 具体例
// and 疑問点 for the stickies, and アクティビティ・ユーザータスク・ユーザー
// ストーリー for the map's three layers, as the Japanese edition of Patton's
// ユーザーストーリーマッピング names them — which also calls the board itself
// ストーリーマップ. These are this repository's ubiquitous language too; the
// two are one vocabulary and must not drift apart.
var ja = Catalog{
	"lang.code": "ja",

	"nav.example-mappings": "実例マッピング",
	"nav.story-maps":       "ストーリーマップ",
	"nav.stories":          "ユーザーストーリー",
	"nav.ubiquitous":       "ユビキタス言語",
	"nav.tasks":            "タスク",

	"label.story":           "ユーザーストーリー",
	"label.rule":            "ルール",
	"label.example":         "具体例",
	"label.question":        "疑問点",
	"label.questions":       "疑問点",
	"label.opportunity":     "オポチュニティ",
	"label.activity":        "アクティビティ",
	"label.step":            "ユーザータスク",
	"label.context":         "文脈",
	"label.example-mapping": "実例マッピング",

	"filter.all": "すべて",

	"badge.copy-link": "リンクをコピー",
	"badge.copied":    "✓ コピーしました",

	"empty.example-mappings": "実例マッピングはまだありません。",
	"empty.story-maps":       "ストーリーマップはまだありません。",
	"empty.stories":          "ユーザーストーリーが見つかりません。",
	"empty.terms":            "用語が見つかりません。",

	"tasks.questions":                "未解決の疑問点",
	"tasks.questions-hint":           "会話で閉じる",
	"tasks.questions-empty":          "未解決の疑問点はありません。",
	"tasks.questions-empty-filtered": "このオポチュニティに未解決の疑問点はありません。",
	"tasks.rules":                    "未自動化のルール",
	"tasks.rules-hint":               "テストで閉じる",
	"tasks.rules-empty":              "すべてのルールが自動化されています。",
	"tasks.rules-empty-filtered":     "このオポチュニティのルールはすべて自動化されています。",

	"glossary.term":       "用語",
	"glossary.key":        "キー",
	"glossary.definition": "定義",

	"story.metadata":    "メタデータ",
	"story.related":     "関連",
	"story.description": "説明",

	"mapping.automated-legend": "自動化済みルール",
	"mapping.automated-badge":  "✓ 自動化済み",
	"mapping.automated-title":  "テストで自動化済み",
	"nav.opportunities":        "オポチュニティ",

	"label.story-map":          "ストーリーマップ",
	"label.opportunity-canvas": "オポチュニティキャンバス",

	"empty.opportunities":      "オポチュニティはまだありません。",
	"empty.opportunity-canvas": "オポチュニティキャンバスはまだありません。",

	"opportunity.opportunity": "オポチュニティ",

	"canvas.zone.facts":    "確かめられる事実",
	"canvas.zone.solution": "解決策",
	"canvas.zone.value":    "価値についての仮説",
	"canvas.empty-box":     "未記入",

	"canvas.solution-ideas":             "解決策のアイデア",
	"canvas.solution-ideas.prompt":      "具体的なプロダクト・機能・改善のアイデアを挙げる。",
	"canvas.problems":                   "課題",
	"canvas.problems.prompt":            "この解決策が扱う課題を、見込みユーザーや顧客は今どう抱えているか。",
	"canvas.users-and-customers":        "顧客とユーザー",
	"canvas.users-and-customers.prompt": "この解決策が扱う課題を抱えているのは、どんなユーザー・顧客か。",
	"canvas.solutions-today":            "現在の解決手段",
	"canvas.solutions-today.prompt":     "ユーザーは今、その課題にどう対処しているか。",
	"canvas.business-challenges":        "事業上の課題",
	"canvas.business-challenges.prompt": "ユーザーと顧客のその課題は、自分たちの事業にどう響いているか。",
	"canvas.user-value":                 "ユーザーは価値を得るために何をするか",
	"canvas.user-value.prompt":          "この解決策が届いたとき、対象のユーザーは何をするか。",
	"canvas.user-metrics":               "指標",
	"canvas.user-metrics.prompt":        "価値を得るまでの一連の行動のうち、実際にそうしたと示せるものとして何を測れるか。",
	"canvas.adoption-strategy":          "採用の道筋",
	"canvas.adoption-strategy.prompt":   "顧客とユーザーは、この解決策をどう知り、使い方を覚え、使い続けるようになるか。",
	"canvas.business-impact":            "事業へのインパクト",
	"canvas.business-impact.prompt":     "この解決策がうまくいったとき、どの事業指標が動くか。",
	"canvas.budget":                     "予算",
	"canvas.budget.prompt":              "この課題を解いてその成果を得るために、どれだけの費用やチームの時間を充てるか。",
}
