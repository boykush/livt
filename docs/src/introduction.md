# livt - Living Text

> Collaborate on board. Make it living in text.

## Outcome

Collaborative outcomes evolve alongside the product as living text.

## Opportunities

- **Stale Discovery**: Discovery-phase artifacts are not persisted after synchronous collaboration sessions
- **Discovery-Development Gap**: Persisted discovery artifacts are not leveraged in the development process

## Solutions

- Provide a CLI tool that captures collaborative outcomes as text files (YAML, Markdown)
- Track consistency across artifacts via ID-based references with automated checks
- Treat Formulation artifacts (Gherkin scenarios) as generated output, not the source of truth — that lives in the Discovery artifacts

## Live Demo

livt dogfoods itself: its own discovery artifacts — stories, story maps, example mappings, and ubiquitous language — are published with `livt build` as a [live demo](https://boykush.github.io/livt/demo/). It shows exactly what the guides below describe.
