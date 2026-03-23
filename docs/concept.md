# livt - Living Text

> Collaborate on board. Make it living in text.

## Outcome

Collaborative outcomes evolve alongside the product as living text.

## Opportunities

- Discovery-phase artifacts (Example Mapping, USM, etc.) live in synchronous collaboration tools (Miro, Figjam) and become stale after the session — they lack version control, AI integration, and consistency checks
- As Rules and Examples evolve through development, changes are not tracked back to Discovery artifacts, breaking the link between Discovery and Formulation

## Solutions

- Provide a CLI tool that captures collaborative outcomes as text files (YAML, Markdown)
- Track consistency across artifacts via ID-based references with automated checks
- Treat Formulation artifacts (Gherkin scenarios) as generated output, not the source of truth — the master lives in Discovery artifacts

## Experiments

- Dogfood livt's own development process using livt's format
