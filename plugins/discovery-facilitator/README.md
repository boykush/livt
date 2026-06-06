# Discovery Facilitator Plugin

Skills and expert agents for turning collaborative discovery sessions — BDD Example Mapping and User Story Mapping — into reviewed, reconciled artifacts.

## Overview

Live facilitation happens with people on a board (Miro, sticky notes). This plugin handles the **post-collaboration** workflow: writing the board out and refining it through a three-step pipeline.

```
Transcribe ──▶ Review (清書) ──▶ Reconcile
  (skill)        (expert agent)     (skill)
```

1. **Transcribe** — write the agreed board out to YAML faithfully, committed with no review as the diff baseline.
2. **Review / 清書** — an expert agent brushes up the structure and ships it as a reviewable PR diff on top of the baseline. The review's value is exactly that diff.
3. **Reconcile** — diff the result against the current implementation and design to surface gaps and omissions.

## Skills

### `/usm-transcriber`

Transcribe a completed User Story Mapping board into `discoveries/usm/{map-name}.yaml`. Faithfully captures the backbone (activities, steps) and story cards in board order, and commits the result as the review-free baseline.

### `/example-mapping-transcriber`

Transcribe a completed Example Mapping board into `discoveries/example-mappings/{story-key}.yaml`. Faithfully captures rules, examples, and questions, and commits the result as the review-free baseline.

### `/example-mapping-reconciler`

Reconcile a reviewed example mapping against the current implementation and design. Diffs rules/examples/questions against what the code and design actually do, and surfaces contradictions, missing implementations, and omissions in the collaboration outcome.

## Agents

### `usm-expert`

Expert consultant grounded in Jeff Patton's User Story Mapping methodology. Serves as the **reviewer (清書)** for story-map artifacts — challenging scope and answering questions about narrative flow, backbone structure, and release slicing.

### `bdd-expert`

Expert in Behaviour-Driven Development processes (Discovery, Formulation, Automation), Example Mapping, and Gherkin syntax. Serves as the **reviewer (清書)** for stories and example mappings, answering BDD practice questions and consulting on artifact consistency.

## Install

```
/plugin install discovery-facilitator@boykush/livt
```
