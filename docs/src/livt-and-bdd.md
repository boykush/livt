# livt and BDD

livt renders practices it did not invent. This page says which parts are the community's, which parts livt added, where livt departs on purpose, and what livt deliberately leaves alone.

If you arrived from the BDD community, the short version: **livt is not an alternative to Cucumber.** It runs no tests. It keeps what a Discovery session decided, so the decision is still readable — and still quotable — when someone implements it months later.

## What livt takes

- **Behaviour-Driven Development** — Dan North, [Introducing BDD](https://dannorth.net/introducing-bdd/) (2006). The three phases livt organizes around, Discovery / Formulation / Automation, are the community's own framing: [Cucumber's BDD docs](https://cucumber.io/docs/bdd/).
- **Example Mapping** — Matt Wynne, [Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/) (2015). The four cards, their colours, the timebox, and reading the board back as a signal about the story. livt's board is that format, rendered from a file.
- **Ubiquitous language** — Eric Evans' term from Domain-Driven Design, and a thread through BDD's Formulation. livt keeps it as a [glossary](./guides/ubiquitous-language.md) beside the artifacts that use it.
- **Opportunity Canvas and User Story Mapping** — Jeff Patton. Not BDD, but under the same rule; the [Opportunities](./guides/opportunities.md) and [Story Maps](./guides/story-maps.md) guides link the originals.

## What livt adds

None of this is in the practice. It is livt's, and it moves when livt moves:

- **A file instead of a photograph.** The mapping is YAML in the repository, versioned alongside the code it specifies.
- **Identity.** Every rule, example, and question has an ID that is stable once filed, numbered past retired ones, and addressable as a [livt URI](./reference/uri.md) — so a rule quoted from an Issue or a test comment cannot silently come to mean something else.
- **A lifecycle.** A rule carries `status: proposed`, `accepted`, `rejected`, or `retired`, borrowed from [architecture decision records](https://adr.github.io/), and `superseded_by` says where the spec went. See [Example Mappings](./guides/example-mappings.md).
- **An asynchronous lane.** A rule can be proposed and agreed in a pull request review, between sessions, when the room does not need to meet for it.
- **Links outward.** A rule records the Issues filed to automate it, and the tests that automate it cite it back by livt URI. livt collects those citations from the implementation repository — what each test claims to cover, not whether it passes — and the site rolls them up.

## Where livt departs

### An Example Mapping is supposed to be disposable

This is the real departure, and it is worth naming rather than smuggling.

Wynne's Example Mapping is a twenty-five minute conversation whose product is shared understanding. The cards are a byproduct and a gauge: too many blue and the story is too big, too many red and it is not ready. Once they have said that, they have done their job. A team that photographs the board and never opens the photo again has not failed at the practice.

livt keeps the map, because of the gap it was built for: the developer who needs the agreed rule six months later cannot query a photograph. But a durable artifact invites maintenance, and maintenance invites treating the map as a specification to be completed rather than a conversation to be had. That anti-pattern is real, and no tool can rule it out. livt's structural answer is that the record is written after the room agrees and never changes what was agreed — the [record skills](https://github.com/boykush/livt/tree/main/plugins/livt-discovery) write the board verbatim first, and only then a structural edit whose diff is what the review reads.

So, plainly: **if you have to choose, choose the conversation.** A team that runs good sessions and keeps no record is doing BDD. A team with immaculate YAML and no session is not.

### livt takes a position on Formulation that the community does not share

Mainstream BDD treats the feature file as the shared artifact: the Gherkin is what the Three Amigos produce together, and it is the living documentation. livt's [stated intent](./introduction.md) runs the other way — Discovery artifacts are the source of truth, and Gherkin would be generated output.

Be clear about what that is worth today: **livt does not generate Gherkin at all.** Nothing in the tool depends on the position yet. It is a view under test, it is the part of livt most likely to be wrong, and it belongs in an argument rather than in a release note. If you think it is wrong, [say so](https://github.com/boykush/livt/issues) — that is a more useful contribution than agreement.

## What livt does not do

- **Run tests.** livt is not an alternative to Cucumber, SpecFlow, Behave, or anything else that executes specifications. If it ever emits Gherkin, it will emit the standard syntax for those tools to run, never a dialect of its own.
- **Facilitate the session.** No timer, no voting, no board. Live facilitation happens with people, on whatever board the room already uses; livt starts where the board ends.
- **Teach BDD.** These guides describe file formats and what the site renders. The practice is taught by the people below, and taught better.
- **Redefine the community's words.** Where a term is the practice's — rule, example, question, discovery, formulation — it keeps the practice's meaning. Where livt needed a term of its own, the name says so: livt repository, livt URI, implementation repository.

## The people who wrote this down

Most of what livt renders was worked out and published by a community that gives away the introduction and sells the depth, and whose maintainers have run on donations since 2023.

- Dan North — [Introducing BDD](https://dannorth.net/introducing-bdd/)
- Matt Wynne — [Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/)
- Seb Rose — [Keep your scenarios BRIEF](https://cucumber.io/blog/bdd/keep-your-scenarios-brief/)
- Gáspár Nagy and Seb Rose — [The BDD Books](https://bddbooks.com/): *Discovery* and *Formulation*, with Japanese editions on Leanpub
- [Cucumber School](https://school.cucumber.io/) — the courses, free and paid
- [Cucumber on Open Collective](https://opencollective.com/cucumber) — where the maintainers are funded

If livt saved your team an argument it had already had, some of that is theirs.
