---
title: livt and BDD - livt
---

# livt and BDD

## Where livt departs from the practice

### Formulation: the Example Mapping is the source

In mainstream BDD, developers and testers write the rules and examples up as a feature file (Gherkin), and that file becomes the source of the specification. In livt, a YAML file in Example Mapping's own format is the source, and [formulation](https://boykush.github.io/livt/demo/ubiquitous.html#formulation) is rewriting that file. A feature file, where a team wants one, would be generated from it, which livt does not do yet.

Because the format stays Example Mapping's, product managers and designers can take part in answering its questions and changing its rules afterwards. All anyone needs to know is Example Mapping, so there is less to learn.

<!-- figure: The practice above, livt below. Practice: Example Mapping → (written up by developers and testers) → feature file (Gherkin) (the source; developers, testers); the feature file stays with development and testing. livt: Example Mapping → (recorded, then rewritten) → YAML in Example Mapping's format (the source; the three amigos), with a dashed line on to a feature file (generated if wanted); product managers and designers join in answering its questions and changing its rules. -->

### Automation: coding agents take it on

Coding agents can now take on automating what a team agreed. Following the test pyramid, one Example Mapping's rules and examples end up automated at different layers and in different languages. livt gives the agent a shape to follow and the tools to follow it, as a CLI and plugins.

<!-- figure: One Example Mapping's rules and examples, cited with `livt:automates` by a backend test (Go), a frontend test (TypeScript) and an end-to-end test, the citations gathered into the automation report. The tests are made up for the example. -->

## Learn it from the people who wrote it down

Most of what livt renders was worked out and published by the BDD community.

- Dan North — [Introducing BDD](https://dannorth.net/introducing-bdd/)
- Matt Wynne — [Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/)
- Seb Rose — [Keep your scenarios BRIEF](https://cucumber.io/blog/bdd/keep-your-scenarios-brief/)
- Gáspár Nagy, Seb Rose — [The BDD Books](https://bddbooks.com/): Discovery, Formulation (Japanese editions on Leanpub)
