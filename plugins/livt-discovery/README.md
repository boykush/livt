# livt-discovery

Skills for the **discovery ring** of a livt repository: the conversations that decide what to build — Opportunity Canvas, User Story Mapping, Example Mapping — and their record as text.

## The ring

```
opportunity ─▶ story map ─▶ card ─▶ conversation ─▶ record ─▶ livt repository
                                        ▲                          │
                                        └── what remains open ◀────┘
```

Live facilitation happens with people on a board (Miro, sticky notes). These skills start where the board ends, and each is named for the **station** it serves rather than for how it edits a file:

- **Record** — *talk, then record* (Patton). A session's outcome is written out as two commits on one PR: the board verbatim as the baseline, then a structural edit, consulting an expert skill, whose diff is what the review reads. The record never changes what the room agreed.
- **Card** — the Card of Ron Jeffries' three Cs. A story candidate picked from the map gets `stories/{story-key}.md` and a stable key, so its own conversation — an Example Mapping — can be held and recorded against it.
- **Change** — discovery's **asynchronous lane**. A rule proposed, changed, or retired outside a session lands in the mapping as its own fine-grained PR; a proposal carries `status: proposed`, the PR review stands in for the conversation, and agreement is a one-line diff. What stays unagreed waits on the Tasks page as the next session's agenda.

No skill here reads the implementation. The first read of the code is where the [delivery ring](../livt-delivery/README.md) begins.

## Skills

### Record

- **`/record-story-map`** — record a User Story Mapping session into `discoveries/usm/{map-name}.yaml`: backbone (activities, user tasks), story cards, release slices, and the board's ubiquitous terms, as a verbatim baseline plus a structural edit consulting `usm-expert`.
- **`/record-example-mapping`** — record an Example Mapping session into `discoveries/example-mappings/{story-key}.yaml`: rules, examples, questions, and terms, as a verbatim baseline plus a structural edit consulting `bdd-expert`. Re-recording a reworked board keeps every ID.

### Card

- **`/write-story-card`** — write a story's card: create `stories/{story-key}.md` for a candidate on a story map and stamp the key back onto the map. The bridge from the map to the story's own conversation; it registers the card and changes no meaning.

### Change

- **`/change-rule`** — change one business rule in an existing example mapping — propose, accept, reject, change, or retire it with its examples — as its own fine-grained PR. The canonical statement of the ID contract (numbering past retired IDs, immutability, retire-don't-delete, `superseded_by`) lives here.

## Expert skills

The expert skills are the **knowledge backend** the record skills consult for their structural edit; you can also invoke them standalone (`/opportunity-expert`, `/usm-expert`, `/bdd-expert`) for ad-hoc consulting.

They are plain [Agent Skills](https://agentskills.io), so the consultation travels with the plugin's `skills/` and works in any conformant runtime — the record skills reference them as peer skills, not as runtime-specific subagents.

- **`opportunity-expert`** — Jeff Patton's Opportunity Canvas: framing an opportunity, keeping verifiable facts apart from assumptions about value, and supporting the decision of whether to take it on at all, including the decision not to.
- **`usm-expert`** — Jeff Patton's User Story Mapping: narrative flow, backbone structure, release slicing, and story scope.
- **`bdd-expert`** — Behaviour-Driven Development (Discovery, Formulation, Automation), Example Mapping, and Gherkin syntax.

## Install

```
/plugin install livt-discovery@boykush/livt
```

Coming from `discovery-facilitator` 0.x? The skills were renamed for their stations and split across two plugins:

| Was | Now |
|-----|-----|
| `/usm-transcribe` + `/usm-refine` | `/record-story-map` (one PR, two commits) |
| `/example-mapping-transcribe` + `/example-mapping-refine` | `/record-example-mapping` (one PR, two commits) |
| `/story-commit` | `/write-story-card` |
| `/example-mapping-update` | `/change-rule` |
| `/example-mapping-plan` | `/plan-story` — in [livt-delivery](../livt-delivery/README.md) |
| `/story-issue-file` | `/file-story-issue` — in livt-delivery |
| `/rule-issue-file` | `/file-rule-issues` — in livt-delivery |
| `/rule-automation-sync` | `/inspect-automation` — in livt-delivery |
