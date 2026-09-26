---
title: livt と BDD - livt
---

# livt と BDD

## AI 時代の BDD

AI エージェントの進化で、合意したルールの自動化はエージェントに任せられるようになりました。テストピラミッドに沿うと、1つの実例マッピングのルールや具体例が、別々のレイヤや言語のテストで自動化されます。livt はそのための型とツールを、CLI と AI プラグインとして渡します。

<!-- figure: 1つの実例マッピングのルールと具体例を、バックエンド（Go）・フロントエンド（TypeScript）・E2E のテストが `livt:automates` で引用し、引用が自動化レポートに集まる。テストはこの例のための架空のもの。 -->

## プラクティスと違うところ

### 実例マッピングを残す

<!-- figure: 上にプラクティス、下に livt。プラクティス：付箋のボード →（25 分の会話）→ Gherkin、ボードは役目を終える。livt：付箋のボード →（合意）→ 実例マッピング →（疑問点の解消・ルールの変更）→ 完了、全職能が同じマップを追う。 -->

Matt Wynne の実例マッピングでは、カードは 25 分の会話のための道具です。会話のあと、ルールと具体例は Gherkin に移ります。

livt は、実例マッピングそのものを残します。疑問点の解消やルールの変更は開発のあいだも続き、デザイナーを含む全職能が担うからです。

### 定式化についての立場

<!-- figure: 上にプラクティス、下に livt。プラクティス：実例マッピング →（書き写す）→ feature file（正）。livt：実例マッピング（正）→（書き換える）→ 実例マッピング、そこから点線で Gherkin（要るなら生成）。 -->

主流の BDD では feature file が仕様の正ですが、livt では実例マッピングが正です。livt の[定式化](https://boykush.github.io/livt/demo/ubiquitous.html#formulation)は、そのマッピングを書き換えることです。Gherkin は要るならそこから生成するもので、今の livt はまだ生成しません。

## BDD を学ぶなら

livt が描いているものの多くは、BDD コミュニティが練り上げて公開してきたものです。

- Dan North — [Introducing BDD](https://dannorth.net/introducing-bdd/)
- Matt Wynne — [Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/)
- Seb Rose — [Keep your scenarios BRIEF](https://cucumber.io/blog/bdd/keep-your-scenarios-brief/)
- Gáspár Nagy、Seb Rose — [The BDD Books](https://bddbooks.com/)：Discovery、Formulation（Leanpub に日本語版あり）
