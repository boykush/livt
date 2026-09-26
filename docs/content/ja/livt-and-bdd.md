---
title: livt と BDD - livt
---

# livt と BDD

livt はテストを実行しません。Cucumber の代わりでもありません。

発見でチームが協働して決めたことを残して、何か月か後に実装する人が読んだり引用したりできるようにするためのツールです。

## 土台にしたプラクティス

- **振る舞い駆動開発（BDD）** — Dan North「[Introducing BDD](https://dannorth.net/introducing-bdd/)」（2006）。livt が軸にしている発見・定式化・自動化の 3 フェーズは、コミュニティ自身の整理です（[Cucumber の BDD ドキュメント](https://cucumber.io/docs/bdd/)）。
- **実例マッピング** — Matt Wynne「[Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/)」（2015）。4 種類のカードと色、タイムボックス、ボードの様子からストーリーの状態を読むこと。livt のボードは、この形式をファイルから描いたものです。
- **ユーザーストーリーマッピング** — Jeff Patton「[story mapping](https://jpattonassociates.com/story-mapping/)」。
- **オポチュニティキャンバス** — Jeff Patton「[Opportunity Canvas](https://jpattonassociates.com/opportunity-canvas/)」。
- **ユビキタス言語** — Eric Evans『ドメイン駆動設計』。BDD の定式化にも通じる考え方です。livt では、用語を使っているボードのすぐそばに用語集として置きます。

## livt が足したもの

- **写真ではなくファイル**　マッピングはリポジトリの YAML ファイルで、対応するコードと一緒にバージョン管理されます。
- **変わらない ID**　ルール・具体例・疑問点には、一度付いたら変わらない ID と、そこを指す URI が付きます。Issue やテストで引用したルールが、あとから別の意味になることはありません。
- **テストから集める自動化レポート**　ルールには、自動化のために起票した Issue が記録されます。テストの側はコメントにルールの URI を書いて、どのルールを自動化したかを示します。livt はそれを集めて自動化レポートにし、ルールごとの自動化の状況をそこから導きます。テストが通るかどうかは見ません。

## プラクティスと違うところ

### 実例マッピングを残す

Matt Wynne の実例マッピングでは、カードは 25 分の会話の副産物です。ストーリーの準備ができているかを示した時点で役目を終えるので、ボードを二度と見返さないチームも間違ってはいません。

それでも livt がマップを残すのは、半年後に合意済みのルールを確かめたい開発者が、写真からは探せないからです。ただ、残したマップは一人で埋めていく仕様書になりがちです。そこで livt は、チームが合意したあとにだけ書き起こします。[書き起こし用の skill](https://github.com/boykush/livt/tree/main/plugins/livt-discovery) がまずボードをそのまま写し、構造を整える変更は別の差分にしてレビューします。大事なのはあくまで協働で、livt はその結果を残すだけです。

### 定式化についての立場

主流の BDD では、スリーアミーゴスが一緒に書く Gherkin のフィーチャーファイルが共有の成果物です。livt は逆の立場で、発見の成果物を正とし、Gherkin はそこから生成するものと考えています。

ただ、今の livt は Gherkin を生成しませんし、この立場に依存する機能もまだありません。livt の中でいちばん間違っている可能性が高い部分なので、違うと思ったら [Issue](https://github.com/boykush/livt/issues) で教えてください。

## livt がしないこと

- **テストの実行**　Cucumber、SpecFlow、Behave の代わりにはなりません。Gherkin を出力することがあっても、それらのツールがそのまま実行できる標準の構文にします。
- **協働の進行**　タイマーも投票もボードもありません。協働は、チームがいつも使っているオンラインホワイトボードなどで進めてください。livt の出番はそのあとです。
- **BDD を教えること**　ここに挙げた人たちの資料で学べます。
- **用語の再定義**　ルール、具体例、疑問点、発見、定式化は、プラクティスでの意味のまま使います。

## BDD を学ぶなら

livt が描いているものの多くは、BDD コミュニティが練り上げて公開してきたものです。メンテナーは 2023 年から寄付で活動しています。

- Dan North — [Introducing BDD](https://dannorth.net/introducing-bdd/)
- Matt Wynne — [Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/)
- Seb Rose — [Keep your scenarios BRIEF](https://cucumber.io/blog/bdd/keep-your-scenarios-brief/)
- Gáspár Nagy、Seb Rose — [The BDD Books](https://bddbooks.com/)：Discovery、Formulation（Leanpub に日本語版あり）
- Cucumber — [Cucumber School](https://school.cucumber.io/)、無料と有料のコース
- Cucumber — [Open Collective](https://opencollective.com/cucumber)、メンテナーへの寄付の窓口
