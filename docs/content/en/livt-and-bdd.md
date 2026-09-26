---
title: livt and BDD - livt
---

# livt and BDD

## Where livt departs from the practice

### Discovery: it keeps the Example Mapping

In Matt Wynne's Example Mapping, the cards serve a 25-minute conversation, and their job ends with it.

livt keeps the Example Mapping itself. Its questions get answered and its rules change all through development, so the three amigos who met for the conversation carry on with the same map until it is done.

<!-- figure: The practice above, livt below. Practice: Example Mapping (the three amigos, 25 minutes); its job ends with the conversation. livt: Example Mapping (the three amigos) → (questions answered, rules changed) → done; the same three amigos follow the same map to the end. -->

### Formulation: the Example Mapping is the source

In mainstream BDD, developers and testers write the rules and examples up as Gherkin, and that feature file becomes the source of the specification. In livt the Example Mapping is the source, and [formulation](https://boykush.github.io/livt/demo/ubiquitous.html#formulation) is rewriting it. Gherkin, where a team wants it, would be generated from it, which livt does not do yet.

<!-- figure: The practice above, livt below. Practice: Example Mapping → (written up by developers and testers) → feature file (the source; developers, testers); the feature file stays with development and testing. livt: Example Mapping (the source; the three amigos) → (rewritten) → Example Mapping, with a dashed line on to Gherkin (generated if wanted). -->

### Automation: coding agents take it on

Coding agents can now take on automating what a team agreed. Following the test pyramid, one Example Mapping's rules and examples end up automated at different layers and in different languages. livt gives the agent a shape to follow and the tools to follow it, as a CLI and plugins.

<!-- figure: One Example Mapping's rules and examples, cited with `livt:automates` by a backend test (Go), a frontend test (TypeScript) and an end-to-end test, the citations gathered into the automation report. The tests are made up for the example. -->

## Learn it from the people who wrote it down

Most of what livt renders was worked out and published by the BDD community.

- Dan North — [Introducing BDD](https://dannorth.net/introducing-bdd/)
- Matt Wynne — [Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/)
- Seb Rose — [Keep your scenarios BRIEF](https://cucumber.io/blog/bdd/keep-your-scenarios-brief/)
- Gáspár Nagy, Seb Rose — [The BDD Books](https://bddbooks.com/): Discovery, Formulation (Japanese editions on Leanpub)
