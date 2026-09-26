---
title: livt とプラクティス - livt
---

# livt とプラクティス

livt は、発見でチームが協働して決めたことを残して、何か月か後に実装する人が読んだり引用したりできるようにするためのツールです。

## 土台にしたプラクティス

- **振る舞い駆動開発（BDD）** — Dan North「[Introducing BDD](https://dannorth.net/introducing-bdd/)」（2006）。livt が軸にしている発見・定式化・自動化の 3 フェーズは、Gáspár Nagy と Seb Rose が [The BDD Books](https://bddbooks.com/) で整理したものです。
- **実例マッピング** — Matt Wynne「[Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/)」（2015）。4 種類のカードと色、タイムボックス、ボードの様子からストーリーの状態を読むこと。livt のボードは、この形式をファイルから描いたものです。
- **ユーザーストーリーマッピング** — Jeff Patton「[story mapping](https://jpattonassociates.com/story-mapping/)」。
- **オポチュニティキャンバス** — Jeff Patton「[Opportunity Canvas](https://jpattonassociates.com/opportunity-canvas/)」。
- **ユビキタス言語** — Eric Evans『ドメイン駆動設計』。BDD の定式化にも通じる考え方です。livt では、用語を使っているボードのすぐそばに用語集として置きます。

## livt が足したもの

- **写真ではなくファイル**　ボードを YAML にして、リポジトリでバージョン管理する
  `discoveries/example-mappings/collect-automations.yaml`
- **変わらない ID**　一度付いた ID は変わらず、URI でどこからでも指せる
  `livt://mapping/collect-automations/rule/R-01`
- **疑問点もルールと並べて**　誰が持っていても解消まで追え、解消したら決着したルールを指す
  `Q-01 → R-08`
- **テストから集める自動化レポート**　テストのコメントに URI を書くと、ルールごとに集まる
  `// livt:automates livt://mapping/collect-automations/rule/R-01`

## プラクティスと違うところ

### 実例マッピングを残す

Matt Wynne の実例マッピングでは、カードは 25 分の会話のための道具で、ストーリーの準備ができているかも示します。会話のあと、ルールと具体例は Gherkin に移り、ボードは役目を終えます。

それでも livt がマップを残すのは、実例マッピングが会話の場で終わらないからです。ルールは変わり、疑問点は開発のあいだに解消されていきます。その移り変わりを、疑問点も含めて全職能で追える形で残します。ただ、残したマップは一人で埋めていく仕様書になりがちです。そこで livt は、チームが合意したあとにだけ書き起こします。[書き起こし用の skill](https://github.com/boykush/livt/tree/main/plugins/livt-discovery) がまずボードをそのまま写し、構造を整える変更は別の差分にしてレビューします。大事なのはあくまで協働で、livt はその結果を残すだけです。

### 定式化についての立場

livt は、feature file ではなく実例マッピングを正とし、その形式のまま全職能の持ち物にします。livt での定式化は、そのマッピングを書き換えることです。Gherkin は、要るならそこから生成するものと考えていますが、今の livt はまだ生成しません。

### 自動化はエージェントに任せる

AI エージェントの進化で、状況が変わりました。合意したルールの自動化は、エージェントに任せられます。テストピラミッドに沿って自動化すると、実例マッピングのルールや具体例が、別々のレイヤや言語のテストで自動化されます。livt は、そのための型とツールを CLI と AI プラグインとして渡し、レイヤや言語の違いから来る課題を、エージェントと一緒に解きます。

## livt がしないこと

- **テストの実行**　テストはそれぞれのリポジトリに置いたまま、いま使っているツールで動かします。
- **BDD を教えること**　ここに挙げた人たちの資料で学べます。
- **用語の再定義**　ルール、具体例、疑問点、発見、定式化は、プラクティスでの意味のまま使います。

## BDD を学ぶなら

livt が描いているものの多くは、BDD コミュニティが練り上げて公開してきたものです。

- Dan North — [Introducing BDD](https://dannorth.net/introducing-bdd/)
- Matt Wynne — [Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/)
- Seb Rose — [Keep your scenarios BRIEF](https://cucumber.io/blog/bdd/keep-your-scenarios-brief/)
- Gáspár Nagy、Seb Rose — [The BDD Books](https://bddbooks.com/)：Discovery、Formulation（Leanpub に日本語版あり）
