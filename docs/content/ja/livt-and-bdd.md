---
title: livt と BDD - livt
---

# livt と BDD

## AI 時代の BDD

BDD では、合意した仕様を1つの一貫した形で持ち続けます。

ただ、実例マッピングは会話の場で終わりません。残った疑問点は開発のあいだに解消され、ルールも変わっていきます。livt は実例マッピングそのものを形式にして記録し、feature file に当たる内容を含んだまま、デザイナーやプロダクトマネージャー（PdM）を含む全職能で、完了まで追えるようにします。

AI エージェントの進化で、状況が変わりました。合意したルールの自動化は、エージェントに任せられます。テストピラミッドに沿って自動化すると、実例マッピングのルールや具体例が、別々のレイヤや言語のテストで自動化されます。livt は、そのための型とツールを CLI と AI プラグインとして渡し、レイヤや言語の違いから来る課題を、エージェントと一緒に解きます。

## プラクティスと違うところ

### 実例マッピングを残す

Matt Wynne の実例マッピングでは、カードは 25 分の会話のための道具で、ストーリーの準備ができているかも示します。会話のあと、ルールと具体例は Gherkin に移り、ボードは役目を終えます。

それでも livt はマップを残します。ただ、残したマップは一人で埋めていく仕様書になりがちです。そこで livt は、チームが合意したあとにだけ書き起こします。[書き起こし用の skill](https://github.com/boykush/livt/tree/main/plugins/livt-discovery) がまずボードをそのまま写し、構造を整える変更は別の差分にしてレビューします。大事なのはあくまで協働で、livt はその結果を残すだけです。

### 定式化についての立場

livt は、feature file ではなく実例マッピングを正とし、その形式のまま全職能の持ち物にします。livt での定式化は、そのマッピングを書き換えることです。Gherkin は、要るならそこから生成するものと考えていますが、今の livt はまだ生成しません。

## BDD を学ぶなら

livt が描いているものの多くは、BDD コミュニティが練り上げて公開してきたものです。

- Dan North — [Introducing BDD](https://dannorth.net/introducing-bdd/)
- Matt Wynne — [Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/)
- Seb Rose — [Keep your scenarios BRIEF](https://cucumber.io/blog/bdd/keep-your-scenarios-brief/)
- Gáspár Nagy、Seb Rose — [The BDD Books](https://bddbooks.com/)：Discovery、Formulation（Leanpub に日本語版あり）
