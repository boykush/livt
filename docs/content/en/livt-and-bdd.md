---
title: livt and BDD - livt
---

# livt and BDD

livt runs no tests, and it is not an alternative to Cucumber.

It keeps what a discovery session decided, so the decision can still be read and cited when someone implements it months later.

## What it builds on

- **Behaviour-Driven Development** — Dan North, [Introducing BDD](https://dannorth.net/introducing-bdd/) (2006). The three phases livt is organised around, Discovery, Formulation and Automation, are the community's own framing, as in [Cucumber's BDD docs](https://cucumber.io/docs/bdd/).
- **Example Mapping** — Matt Wynne, [Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/) (2015). The four cards, their colours, the timebox, and reading the board as a signal about the story. livt's board is that format, drawn from a file.
- **User Story Mapping** — Jeff Patton, [story mapping](https://jpattonassociates.com/story-mapping/).
- **Opportunity Canvas** — Jeff Patton, [Opportunity Canvas](https://jpattonassociates.com/opportunity-canvas/).
- **Ubiquitous Language** — Eric Evans, Domain-Driven Design. It runs through BDD's Formulation as well. livt keeps it as a glossary next to the boards that use its terms.

## What livt adds

- **A file, not a photo.** The mapping is YAML in the repository, versioned with the code it specifies.
- **IDs that stay put.** Every rule, example and question gets an ID that never changes once filed, and a URI that points at it. A rule quoted in an issue or a test keeps meaning the same thing.
- **Automation, reported from the tests.** A rule records the issues filed to automate it. The tests that automate it mark it with its URI, and livt collects those marks into an automation report. Each rule's automation status is derived from that report: what the tests claim to cover, not whether they pass.

## Where livt departs from the practice

### It keeps the Example Mapping

In Matt Wynne's Example Mapping, the cards are a by-product of a 25-minute conversation. Once they have shown whether the story is ready, they have done their job, and a team that never looks at the board again has done nothing wrong.

livt keeps the map because a developer who needs an agreed rule six months later cannot search a photo. A kept map can turn into a spec that someone fills in alone, so livt writes a mapping only after the team agrees: [the recording skills](https://github.com/boykush/livt/tree/main/plugins/livt-discovery) first copy the board as it was, then restructure it in a separate diff for review. The session still comes first. livt keeps its outcome and does not replace it.

### Its view of Formulation

Mainstream BDD treats the Gherkin feature file as the shared artifact the Three Amigos write together. livt takes the opposite position: the discovery artifacts are the source, and Gherkin would be generated from them.

livt does not generate Gherkin today, and nothing in it depends on this position yet. It is the part of livt most likely to be wrong. If you think it is, [open an issue](https://github.com/boykush/livt/issues).

## What livt does not do

- **Run tests.** It is not an alternative to Cucumber, SpecFlow or Behave. If livt ever outputs Gherkin, it will be the standard syntax those tools run.
- **Facilitate the session.** No timer, no voting, no board. Run the session on the whiteboard your team already uses, online or in the room. livt picks up after it.
- **Teach BDD.** The people listed here teach it.
- **Redefine the community's words.** Rule, example, question, discovery and formulation keep the meaning the practice gives them.

## Learn it from the people who wrote it down

Most of what livt renders was worked out and published by the BDD community. Its maintainers have run on donations since 2023.

- Dan North — [Introducing BDD](https://dannorth.net/introducing-bdd/)
- Matt Wynne — [Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/)
- Seb Rose — [Keep your scenarios BRIEF](https://cucumber.io/blog/bdd/keep-your-scenarios-brief/)
- Gáspár Nagy, Seb Rose — [The BDD Books](https://bddbooks.com/): Discovery, Formulation (Japanese editions on Leanpub)
- Cucumber — [Cucumber School](https://school.cucumber.io/), free and paid courses
- Cucumber — [Open Collective](https://opencollective.com/cucumber), where the maintainers are funded
