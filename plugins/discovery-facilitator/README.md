# Discovery Facilitator Plugin

Skills and expert agents for turning collaborative discovery sessions — BDD Example Mapping and User Story Mapping — into reviewed, reconciled artifacts.

## Overview

Live facilitation happens with people on a board (Miro, sticky notes). This plugin handles the **post-collaboration** workflow: writing the board out and refining it through a three-step pipeline, where every step is a skill.

```
Transcribe ──▶ Review ──▶ Reconcile
  (skill)     (skill +     (skill)
            expert agent)
```

1. **Transcribe** — write the agreed board out to YAML faithfully, committed with no review as the diff baseline.
2. **Review** — brush up the structure and ship it as a reviewable PR diff on top of the baseline, backed by an expert agent. The review's value is exactly that diff.
3. **Reconcile** — diff the result against the current implementation to surface gaps and omissions.

Between the map and per-story discovery sits a bridge: when a story candidate on a reviewed User Story Map is ready for detailed discovery, the **Commit** skill promotes it into the `stories/` registry — creating its story file and stamping the key back onto the map — so its Example Mapping can begin.

The pipeline runs once per session, but the mapping keeps living afterwards: when a business rule is changed or added in ongoing work, the **Update** skill folds that change in as a fine-grained PR — one rule-level change per PR — instead of re-running the review.

## Skills

### Transcribe

- **`/usm-transcriber`** — transcribe a completed User Story Mapping board into `discoveries/usm/{map-name}.yaml`. Faithfully captures the backbone (activities, steps) and story cards in board order, committed as the review-free baseline.
- **`/example-mapping-transcriber`** — transcribe a completed Example Mapping board into `discoveries/example-mappings/{story-key}.yaml`. Faithfully captures rules, examples, and questions, committed as the review-free baseline.

### Review

- **`/usm-reviewer`** — brush up a transcribed story-map baseline (backbone narrative, granularity, story framing) and ship it as a reviewable PR diff, consulting `usm-expert`.
- **`/example-mapping-reviewer`** — brush up a transcribed example-mapping baseline (rule/example clarity, naming, grouping) and ship it as a reviewable PR diff, consulting `bdd-expert`.

### Commit

- **`/story-committer`** — commit a story candidate from a User Story Map into the story registry: create `stories/{story-key}.md` and write the key back to the matching candidate in the map. The bridge from the map to per-story detailed discovery — run it when a candidate is ready to move from the story map into its own Example Mapping.

### Reconcile

- **`/example-mapping-reconciler`** — reconcile a reviewed example mapping against the current implementation. Diffs rules/examples/questions against what the code actually does and surfaces contradictions, missing implementations, and omissions in the collaboration outcome. The implementation's location is specified by the user.

### Update

- **`/example-mapping-updater`** — fold a rule change or addition from ongoing work into an existing example mapping and ship it as a fine-grained PR, one rule-level change per PR. The right skill when "a rule changed" or "a rule was added" — the reviewer never changes agreed meaning.

## Agents

The expert agents are the **knowledge backend** the review skills consult; they can also be used standalone for ad-hoc consulting.

### `usm-expert`

Expert consultant grounded in Jeff Patton's User Story Mapping methodology — narrative flow, backbone structure, release slicing, and story scope.

### `bdd-expert`

Expert in Behaviour-Driven Development processes (Discovery, Formulation, Automation), Example Mapping, and Gherkin syntax.

## Install

```
/plugin install discovery-facilitator@boykush/livt
```
