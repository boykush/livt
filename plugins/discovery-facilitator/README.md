# Discovery Facilitator Plugin

Pipeline and expert skills for turning collaborative discovery sessions — BDD Example Mapping and User Story Mapping — into refined, planned artifacts.

## Overview

Live facilitation happens with people on a board (Miro, sticky notes). This plugin handles the **post-collaboration** workflow: writing the board out and refining it through a three-step pipeline, where every step is a skill.

```
Transcribe ──▶ Refine ──▶ Plan
  (skill)     (skill +     (skill)
            expert skill)
```

1. **Transcribe** — write the agreed board out to YAML faithfully, committed with no refinement as the diff baseline.
2. **Refine** — refine the structure and ship it as a reviewable PR diff on top of the baseline, backed by an expert skill. The refinement's value is exactly that diff.
3. **Plan** — hold the result against the current implementation and design to surface the work ahead.

Between the map and per-story discovery sits a bridge: when a story candidate on a refined User Story Map is ready for detailed discovery, the **Commit** skill promotes it into the `stories/` registry — creating its story file and stamping the key back onto the map — so its Example Mapping can begin.

The pipeline runs once per session, but the mapping keeps living afterwards: when a business rule is changed or added in ongoing work, the **Update** skill folds that change in as a fine-grained PR — one rule-level change per PR — instead of re-running the refine.

When a mapping is ready to drive implementation, the **File** skills carry it outward as GitHub issues to the story's declared implementation repositories — the master, not GitHub, keeps the record of what is filed where.

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

- **`/example-mapping-plan`** — plan implementation from a refined example mapping by holding it against the current implementation and design. Surfaces the work ahead — unimplemented rules, contradictions, and Questions the code or design already answers — scoped forward. The implementation and design locations are specified by the user or the story's frontmatter.

### Update

- **`/example-mapping-update`** — fold a rule change or addition from ongoing work into an existing example mapping and ship it as a fine-grained PR, one rule-level change per PR. The right skill when "a rule changed" or "a rule was added" — the refiner never changes agreed meaning.

### File

- **`/story-issue-file`** — file a story-level issue to a declared implementation repository (story frontmatter `repos:`), carrying the story body and backpointers to the master. Rule issues already filed in that repository are adopted as sub-issues; the created URL is written back to the story's frontmatter `issues:` as a working-tree edit.
- **`/rule-issue-file`** — file automation issues for business rules, one per rule × declared repository, deduped by the mapping's own record. Each issue carries the rule, its examples, and backpointers (story-key, rule-id, the living document's `#rule-{ID}` anchor, spec rev); the created URL is written back to the rule's `issues:` as a working-tree edit. Filed under the story issue as a sub-issue when one exists in that repository.

## Expert skills

The expert skills are the **knowledge backend** the refine skills consult; you can also invoke them standalone (`/usm-expert`, `/bdd-expert`) for ad-hoc consulting.

They are plain [Agent Skills](https://agentskills.io), so the consultation travels with the plugin's `skills/` and works in any conformant runtime — the refine skills reference them as peer skills, not as runtime-specific subagents.

### `usm-expert`

Expertise grounded in Jeff Patton's User Story Mapping methodology — narrative flow, backbone structure, release slicing, and story scope.

### `bdd-expert`

Expertise in Behaviour-Driven Development processes (Discovery, Formulation, Automation), Example Mapping, and Gherkin syntax.

## Install

```
/plugin install discovery-facilitator@boykush/livt
```
