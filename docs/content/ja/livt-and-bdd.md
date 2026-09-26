---
title: livt と BDD - livt
---

# livt と BDD

## プラクティスと違うところ

### 発見：実例マッピングを残す

Matt Wynne の実例マッピングでは、カードは 25 分の会話のための道具です。会話のあと、ルールと具体例は開発者とテスターが Gherkin に書き起こし、その手元に置かれます。

livt は、実例マッピングそのものを残します。疑問点の解消やルールの変更は開発のあいだも続き、デザイナーを含む全職能が担うからです。

<!-- figure: 上にプラクティス、下に livt。プラクティス：実例マッピング（PdM・開発者・テスターで 25 分）→（開発者とテスターが書き起こす）→ Gherkin（開発者・テスター）、実例マッピングは会話で役目を終え、Gherkin は開発とテストの手元に置かれる。livt：実例マッピング（全職能）→（疑問点の解消・ルールの変更）→ 完了、全職能が同じマップを完了まで追う。 -->

### 定式化：実例マッピングを正にする

主流の BDD では feature file が仕様の正ですが、livt では実例マッピングが正です。livt の[定式化](https://boykush.github.io/livt/demo/ubiquitous.html#formulation)は、そのマッピングを書き換えることです。Gherkin は要るならそこから生成するもので、今の livt はまだ生成しません。

<!-- figure: 上にプラクティス、下に livt。プラクティス：実例マッピング →（開発者とテスターが書き写す）→ feature file（正、開発者・テスター）。livt：実例マッピング（正、全職能）→（書き換える）→ 実例マッピング、そこから点線で Gherkin（要るなら生成）。 -->

### 自動化：AI エージェントに任せる

AI エージェントの進化で、合意したルールの自動化はエージェントに任せられるようになりました。テストピラミッドに沿うと、1つの実例マッピングのルールや具体例が、別々のレイヤや言語のテストで自動化されます。livt はそのための型とツールを、CLI と AI プラグインとして渡します。

<!-- figure: 1つの実例マッピングのルールと具体例を、バックエンド（Go）・フロントエンド（TypeScript）・E2E のテストが `livt:automates` で引用し、引用が自動化レポートに集まる。テストはこの例のための架空のもの。 -->

## BDD を学ぶなら

livt が描いているものの多くは、BDD コミュニティが練り上げて公開してきたものです。

- Dan North — [Introducing BDD](https://dannorth.net/introducing-bdd/)
- Matt Wynne — [Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/)
- Seb Rose — [Keep your scenarios BRIEF](https://cucumber.io/blog/bdd/keep-your-scenarios-brief/)
- Gáspár Nagy、Seb Rose — [The BDD Books](https://bddbooks.com/)：Discovery、Formulation（Leanpub に日本語版あり）
