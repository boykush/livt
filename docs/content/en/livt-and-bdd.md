---
title: livt and BDD - livt
---

# livt and BDD

livt keeps what a team agreed in a discovery conversation, so the decision can still be read and cited when someone implements it months later.

## BDD, now that coding agents write the code

In BDD, a team keeps what it agreed as one consistent specification.

An Example Mapping does not end with the conversation, though. Its open questions get answered during development, and its rules change. livt records the Example Mapping itself, in its own format: it holds what a feature file would, and lets every role follow it through to completion, designers and product managers as much as developers.

Coding agents have changed the situation. Automating an agreed rule is now work an agent can take on. Following the test pyramid, an Example Mapping's rules and examples end up automated at different layers and in different languages. livt gives the agent a shape to follow and the tools to follow it, as a CLI and plugins, and works through what those differences bring with it.

## What it builds on

- **Behaviour-Driven Development** — Dan North, [Introducing BDD](https://dannorth.net/introducing-bdd/) (2006). The three phases livt is organised around, Discovery, Formulation and Automation, are those Gáspár Nagy and Seb Rose set out in [The BDD Books](https://bddbooks.com/).
- **Example Mapping** — Matt Wynne, [Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/) (2015). The four cards, their colours, the timebox, and reading the board as a signal about the story. livt's board is that format, drawn from a file.
- **User Story Mapping** — Jeff Patton, [story mapping](https://jpattonassociates.com/story-mapping/).
- **Opportunity Canvas** — Jeff Patton, [Opportunity Canvas](https://jpattonassociates.com/opportunity-canvas/).
- **Ubiquitous Language** — Eric Evans, Domain-Driven Design. It runs through BDD's Formulation as well. livt keeps it as a glossary next to the boards that use its terms.

## What livt adds

- **A file, not a photo.** Each board becomes YAML, versioned in a repository.
  `discoveries/example-mappings/collect-automations.yaml`
- **IDs that stay put.** An ID never changes once filed, and a URI points at it from anywhere.
  `livt://mapping/collect-automations/rule/R-01`
- **Questions kept beside the rules.** A question stays in view until resolved, whoever holds it, then points to the rule that settled it.
  `Q-01 → R-08`
- **Automation, reported from the tests.** Put a rule's URI in a test comment, and livt gathers the tests rule by rule.
  `// livt:automates livt://mapping/collect-automations/rule/R-01`

## Where livt departs from the practice

### It keeps the Example Mapping

In Matt Wynne's Example Mapping, the cards serve a 25-minute conversation and show whether the story is ready. Afterwards the rules and examples move into Gherkin, and the board has done its job.

livt keeps the map anyway. A kept map can turn into a spec that someone fills in alone, so livt writes a mapping only after the team agrees: [the recording skills](https://github.com/boykush/livt/tree/main/plugins/livt-discovery) first copy the board as it was, then restructure it in a separate diff for review. The conversation still comes first. livt keeps its outcome and does not replace it.

### Its view of Formulation

livt takes the Example Mapping, not a feature file, as the source, and keeps its format so that every role owns it. In livt, formulation means rewriting that mapping. Gherkin, where a team wants it, would be generated from it, though livt does not generate it today.

## What livt does not do

- **Run tests.** Tests stay in your repositories and run with the tools you already use.
- **Teach BDD.** The people listed here teach it.
- **Redefine terms.** Rule, example, question, discovery and formulation keep the meaning the practice gives them.

## Learn it from the people who wrote it down

Most of what livt renders was worked out and published by the BDD community.

- Dan North — [Introducing BDD](https://dannorth.net/introducing-bdd/)
- Matt Wynne — [Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/)
- Seb Rose — [Keep your scenarios BRIEF](https://cucumber.io/blog/bdd/keep-your-scenarios-brief/)
- Gáspár Nagy, Seb Rose — [The BDD Books](https://bddbooks.com/): Discovery, Formulation (Japanese editions on Leanpub)
