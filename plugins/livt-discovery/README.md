# livt-discovery

Skills for the **discovery ring** of a livt repository: the conversations that decide what to build — Opportunity Canvas, User Story Mapping, Example Mapping — and their record as text.

## The ring

```
opportunity ─▶ story map ─▶ card ─▶ conversation ─▶ record ─▶ formulation ─▶ livt repository
                                        ▲                                           │
                                        └── what remains open ◀─────────────────────┘
```

Live facilitation happens with people on a board (Miro, sticky notes). These skills start where the board ends, and each is named for the **station** it serves rather than for how it edits a file:

- **Record** — *talk, then record* (Patton). A session's outcome is written out as text and committed, and the record never changes what the room agreed. On a story map that is two commits, the board verbatim and a structural edit consulting `usm-expert`; on an example mapping the edit has a station of its own.
- **Formulate** — BDD's Formulation, at the story level. The recorded mapping is made into something every reader reads the same behaviour out of — rules that assert, examples named for what they show — as a structural edit over the committed baseline, whose diff is what the review reads. Mainstream BDD formulates by writing Gherkin; livt derives Gherkin from the mapping instead, so the act is the rewrite of the mapping itself. It runs again whenever the wording needs it, with no new board.
- **Card** — the Card of Ron Jeffries' three Cs. A story candidate picked from the map gets `stories/{story-key}.md` and a stable key, so its own conversation — an Example Mapping — can be held and recorded against it. A story taken straight to its examples, such as a bug fix, can go without one: its mapping carries the name instead.
- **Change** — discovery's **asynchronous lane**. A rule proposed, changed, or retired outside a session lands in the mapping as its own fine-grained commit; a proposal carries `status: proposed`, the review stands in for the conversation, and agreement is a one-line diff.

No skill here reads the implementation. The first read of the code is where the [delivery ring](../livt-delivery/README.md) begins.

## What these skills teach, and what they leave to you

Two authorities make a discovery skill true, and neither of them is your team:

- **livt's** — the YAML the record writes, the ID contract (numbering past retired IDs, immutability, retire-don't-delete, `superseded_by`), the key contract, the name contract (a story and its mapping, when both are named, carry one name — livt never compares them, so the skills keep them together), and the line between what a record may change and what only an agreed decision may. `change-rule` holds the canonical statement of the ID contract; `record-example-mapping` and `formulate-example-mapping` repeat it verbatim, as does `record-story-map` for the ubiquitous-language section it shares with the record. `write-story-card` holds the name contract, and `record-example-mapping` repeats it. This half moves when livt moves.
- **The practice's** — Patton on opportunities and story maps, Jeffries' three Cs, North on BDD and Wynne on Example Mapping. The expert skills below are that knowledge and nothing else, which is why they are the part of this plugin that is worth reading even without livt: they do not move when livt moves. Each names its sources, so the summary can be checked against — and traded for — the original.

**Yours** is the board and the session: which tool the room uses, how facilitation is run, what else a card carries once it exists, and your own branch and review conventions — livt's unit of change is the commit, and how commits are grouped into pull requests is a decision no skill here makes for you. The record skills take a board as a photo, an export, or pasted text, and no skill here assumes a board tool — a session held on paper records the same way.

Every skill here is a plain [Agent Skill](https://agentskills.io) — no subagents, no hooks, no slash-command-only behaviour — so the whole discovery ring works in any conformant runtime, and the runtime-specific glue stays on your side of the line. The record skills consult the expert skills as peer skills for exactly that reason.

## Skills

### Record

- **`/record-story-map`** — record a User Story Mapping session into `discoveries/usm/{opportunity-key}.yaml`, the one map of the opportunity it serves: backbone (activities, user tasks), story cards, release slices, and the board's ubiquitous terms, as a verbatim baseline plus a structural edit consulting `usm-expert`.
- **`/record-example-mapping`** — record an Example Mapping session into `discoveries/example-mappings/{story-key}.yaml`: rules, examples, questions, and terms, committed verbatim as the baseline. Re-recording a reworked board keeps every ID.

### Formulate

- **`/formulate-example-mapping`** — formulate a recorded example mapping: rule wording, example naming, grouping, and question phrasing, consulting `bdd-expert`, as one structural edit over the committed baseline. It runs on a mapping recorded minutes or months ago, and changes expression and structure only — never what the room agreed.

### Card

- **`/write-story-card`** — write a story's card: create `stories/{story-key}.md` for a candidate on a story map and stamp the key back onto the map. The bridge from the map to the story's own conversation; it registers the card and changes no meaning.

### Change

- **`/change-rule`** — change one business rule in an existing example mapping — propose, accept, reject, change, or retire it with its examples — as its own fine-grained commit. The canonical statements of the ID contract and the commit contract live here.

## Expert skills

The expert skills are the **knowledge backend** the record and formulation skills consult for their structural edit; you can also invoke them standalone (`/opportunity-expert`, `/usm-expert`, `/bdd-expert`) for ad-hoc consulting.

- **`opportunity-expert`** — Jeff Patton's Opportunity Canvas: framing an opportunity, keeping verifiable facts apart from assumptions about value, and supporting the decision of whether to take it on at all, including the decision not to.
- **`usm-expert`** — Jeff Patton's User Story Mapping: narrative flow, backbone structure, release slicing, and story scope.
- **`bdd-expert`** — Behaviour-Driven Development as its community teaches it: the three phases (Discovery, Formulation, Automation), Matt Wynne's Example Mapping, and Gherkin syntax. Its `Sources` section is where to send anyone who wants more than the summary.

## Install

```
/plugin marketplace add boykush/livt
/plugin install livt-discovery@livt-claude-code-plugins
```

## Coming from 4.x

`/record-story-map` and `/write-story-card` file a story map under its opportunity's key, not under the map's own name:

| Was | Now |
|-----|-----|
| `discoveries/usm/{map-name}.yaml` | `discoveries/usm/{opportunity-key}.yaml` |

livt joins a map to its opportunity by filename, so a map filed under its own name was joined to nothing — and in livt 0.15.0 the map's URI and its page moved to the key as well. Rename a map already on file to its opportunity's key once; its name stays its title. Where no opportunity file exists yet, the key is agreed with you and names the opportunity.

A story and its mapping are kept to one name. A mapping can name its own board and livt never compares that name with its story's, so the skills are the only place the two can be held together: a card names its story, a story with no card is named by its mapping, and where both are named they are one string. The skills that read a story also stop requiring one — a board whose story has no card records fine, and its yellow sticky becomes the mapping's `name:`.

New story keys are minted at four words or more. A key already on file keeps its words, since renaming it would break every URI citing it. It is short keys that collide over years: a two-word key fits every story on its topic, the first card takes it, and it goes on claiming the topic after the others arrive.

`automated:` left the record contract with `/inspect-automation`, which is gone from [livt-delivery](../livt-delivery/README.md). No skill writes the flag any more. livt still reads one a mapping carries, and will until the derived status can carry a board without it, so delete the line as each rule gains its citation. When a rule is reworded nothing unsets its automation for you — the tests go on citing it — and `/change-rule` asks you to say so in the commit body, where the review reads it.

## Coming from 3.x

`/record-example-mapping` now ships the baseline and stops. The structural edit it used to commit on top is `/formulate-example-mapping`:

| Was | Now |
|-----|-----|
| `/record-example-mapping` (baseline + edit) | `/record-example-mapping` (baseline) then `/formulate-example-mapping` (edit) |

Nothing about the two commits changed — the same baseline, the same edit over it, in the same order. What changed is that the edit can now be asked for on its own, against a mapping recorded long ago, which the record skill could not do without a board in hand.

`/record-story-map` is unchanged and still ships both commits. A story map's structural edit is narrative work, not Formulation: the word is defined over rules and examples, and stretching it to cover a backbone would buy a symmetry livt does not mean.

## Coming from 1.x

Nothing here changed shape — the skills, their names, and their inputs are as they were. What is new is that the boundary above is written down, and it is the [delivery plugin](../livt-delivery/README.md) that moved: its three skills gave their tracker recipes back to you. If you had forked a discovery skill, the likely reason was a board or review convention, and that is still yours to hold outside the skill.

Coming from `discovery-facilitator` 0.x? The skills were renamed for their stations and split across two plugins:

| Was | Now |
|-----|-----|
| `/usm-transcribe` + `/usm-refine` | `/record-story-map` (two commits) |
| `/example-mapping-transcribe` + `/example-mapping-refine` | `/record-example-mapping` (two commits) |
| `/story-commit` | `/write-story-card` |
| `/example-mapping-update` | `/change-rule` |
| `/example-mapping-plan` | removed |
| `/story-issue-file` | `/file-story-issue` — in livt-delivery |
| `/rule-issue-file` | `/file-rule-issues` — in livt-delivery |
| `/rule-automation-sync` | removed — automation is collected from the tests |
