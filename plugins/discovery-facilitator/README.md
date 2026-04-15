# Discovery Facilitator Plugin

Facilitation skills and expert agents for BDD Example Mapping and User Story Mapping workshops.

## Overview

This plugin provides AI-powered facilitation for two key discovery practices:

- **User Story Mapping** (Jeff Patton) — Discover activities, tasks, and stories through narrative-driven dialogue
- **Example Mapping** (BDD Discovery) — Explore concrete examples, extract rules, and capture questions

Each practice includes a facilitator skill for running sessions and an expert agent for consulting and reviewing artifacts.

## Skills

### `/usm-facilitator`

Facilitate a User Story Mapping session. Guides the team through discovering activities, breaking them into steps, and slicing releases. Outputs to `discoveries/usm/{map-name}.yaml`.

### `/example-mapping-facilitator`

Facilitate an Example Mapping session on a story. Walks through concrete examples, extracts business rules, and captures unresolved questions. Outputs to `discoveries/example-mappings/{story-key}.yaml`.

## Agents

### `usm-expert`

Expert consultant grounded in Jeff Patton's User Story Mapping methodology. Use for reviewing USM artifacts, challenging scope, and answering questions about narrative flow, backbone structure, and release slicing.

### `bdd-expert`

Expert in Behaviour-Driven Development processes (Discovery, Formulation, Automation), Example Mapping, and Gherkin syntax. Use for reviewing stories and example mappings, answering BDD practice questions, and consulting on artifact consistency.

## Install

```
/plugin install discovery-facilitator@boykush/livt
```
