# Discovery Facilitator Plugin

Pipeline and expert skills for turning collaborative discovery sessions — Opportunity Canvas, User Story Mapping, BDD Example Mapping — into refined, planned artifacts.

## Overview

Live facilitation happens with people on a board (Miro, sticky notes). This plugin handles the **post-collaboration** workflow: writing the board out and refining it through a three-step pipeline, where every step is a skill.

```
Transcribe ──▶ Refine ──▶ Plan
  (skill)     (skill +     (skill)
            expert skill)
```

1. **Transcribe** — write the agreed board out to YAML faithfully, committed with no refinement as the diff baseline.
2. **Refine** — refine the structure and ship it as a reviewable PR diff on top of the baseline, backed by an expert skill. The refinement's value is exactly that diff.
3. **Plan** — take the refined mapping into sprint planning against the current implementation and design: work out how to build each rule — whether it's feasible, what the work is, and which decisions it needs.

Between the map and per-story discovery sits a bridge: when a story candidate on a refined User Story Map is ready for detailed discovery, the **Commit** skill promotes it into the `stories/` registry — creating its story file and stamping the key back onto the map — so its Example Mapping can begin.

The pipeline runs once per session, but the mapping keeps living afterwards: when a business rule is changed or added in ongoing work, the **Update** skill folds that change in as a fine-grained PR — one rule-level change per PR — instead of re-running the refine.

When a mapping is ready to drive implementation, the **File** skills carry it outward as GitHub issues to the story's declared implementation repositories — the livt repository, not GitHub, keeps the record of what is filed where.

## Skills

### Transcribe

- **`/usm-transcribe`** — transcribe a completed User Story Mapping board into `discoveries/usm/{map-name}.yaml`. Faithfully captures the backbone (activities, steps) and story cards in board order, committed as the refinement-free baseline.
- **`/example-mapping-transcribe`** — transcribe a completed Example Mapping board into `discoveries/example-mappings/{story-key}.yaml`. Faithfully captures rules, examples, and questions, committed as the refinement-free baseline.

### Refine

- **`/usm-refine`** — refine a transcribed story-map baseline (backbone narrative, granularity, story framing) and ship it as a reviewable PR diff, consulting `usm-expert`.
- **`/example-mapping-refine`** — refine a transcribed example-mapping baseline (rule/example clarity, naming, grouping) and ship it as a reviewable PR diff, consulting `bdd-expert`.

### Commit

- **`/story-commit`** — commit a story candidate from a User Story Map into the story registry: create `stories/{story-key}.md` and write the key back to the matching candidate in the map. The bridge from the map to per-story detailed discovery — run it when a candidate is ready to move from the story map into its own Example Mapping.

### Plan

- **`/example-mapping-plan`** — plan implementation from a refined example mapping the way a team takes a backlog item into sprint planning: hold it against the current implementation and design and work out how to realize each rule — whether it's feasible, what the work is, which decisions it needs, and where it conflicts with what exists. The implementation and design locations are specified by the user or the story's frontmatter. Not an audit of what's already built, nor a sync from the code back into the mapping.

### Update

- **`/example-mapping-update`** — fold a rule change or addition from ongoing work into an existing example mapping and ship it as a fine-grained PR, one rule-level change per PR. The right skill when "a rule changed" or "a rule was added" — the refiner never changes agreed meaning.

### File

- **`/story-issue-file`** — file a story-level issue to a declared implementation repository (story frontmatter `repos:`), carrying the story body and backpointers to the livt repository. Rule issues already filed in that repository are adopted as sub-issues; the created URL is written back to the story's frontmatter `issues:` as a working-tree edit.
- **`/rule-issue-file`** — file automation issues for business rules, one per rule × declared repository, deduped by the mapping's own record. Each issue carries the rule, its examples, and backpointers (story-key, rule-id, the living document's `#rule-{ID}` anchor, spec rev); the created URL is written back to the rule's `issues:` as a working-tree edit. Filed under the story issue as a sub-issue when one exists in that repository.

### Sync

- **`/rule-automation-sync`** — the return path of the filing skills: hold each rule's `automated:` against what actually happened in the implementation repositories, reading the state of the issues the rule already records and the mapping's own git history. Proposes setting the flag where automation landed and unsetting it where a rule moved on afterwards, one rule's record per PR with the evidence attached. A closed issue is a trigger to propose, never proof — the review decides.

## Expert skills

The expert skills are the **knowledge backend** the refine skills consult; you can also invoke them standalone (`/opportunity-expert`, `/usm-expert`, `/bdd-expert`) for ad-hoc consulting.

They are plain [Agent Skills](https://agentskills.io), so the consultation travels with the plugin's `skills/` and works in any conformant runtime — the refine skills reference them as peer skills, not as runtime-specific subagents.

### `opportunity-expert`

Expertise in Jeff Patton's Opportunity Canvas — framing an opportunity, keeping verifiable facts apart from assumptions about value, and supporting the decision of whether to take it on at all, including the decision not to.

### `usm-expert`

Expertise grounded in Jeff Patton's User Story Mapping methodology — narrative flow, backbone structure, release slicing, and story scope.

### `bdd-expert`

Expertise in Behaviour-Driven Development processes (Discovery, Formulation, Automation), Example Mapping, and Gherkin syntax.

## Install

```
/plugin install discovery-facilitator@boykush/livt
```
