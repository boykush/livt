---
title: livt と BDD - livt
---

# livt と BDD

## プラクティスと違うところ

### 定式化：実例マッピングを正にする

主流の BDD では、開発者とテスターがルールと具体例を Gherkin に書き起こし、その feature file が仕様の正になります。livt では実例マッピングが正で、[定式化](https://boykush.github.io/livt/demo/ubiquitous.html#formulation)はそのマッピングを書き換えることです。Gherkin は要るならそこから生成するもので、今の livt はまだ生成しません。

<!-- figure: 上にプラクティス、下に livt。プラクティス：実例マッピング →（開発者とテスターが書き起こす）→ feature file（正、開発者・テスター）、feature file は開発とテストの手元に置かれる。livt：実例マッピング（正、スリーアミーゴス）→（書き換える）→ 実例マッピング、そこから点線で Gherkin（要るなら生成）。 -->

### 自動化：AI エージェントに任せる

AI エージェントの進化で、合意したルールの自動化はエージェントに任せられるようになりました。テストピラミッドに沿うと、1つの実例マッピングのルールや具体例が、別々のレイヤや言語のテストで自動化されます。livt はそのための型とツールを、CLI と AI プラグインとして渡します。

<!-- figure: 1つの実例マッピングのルールと具体例を、バックエンド（Go）・フロントエンド（TypeScript）・E2E のテストが `livt:automates` で引用し、引用が自動化レポートに集まる。テストはこの例のための架空のもの。 -->

## BDD を学ぶなら

livt が描いているものの多くは、BDD コミュニティが練り上げて公開してきたものです。

- Dan North — [Introducing BDD](https://dannorth.net/introducing-bdd/)
- Matt Wynne — [Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/)
- Seb Rose — [Keep your scenarios BRIEF](https://cucumber.io/blog/bdd/keep-your-scenarios-brief/)
- Gáspár Nagy、Seb Rose — [The BDD Books](https://bddbooks.com/)：Discovery、Formulation（Leanpub に日本語版あり）
