---
name: plan-story
description: Plan a story's implementation from its agreed Example Mapping, the way a team takes a backlog item into sprint planning — hold each rule against the current implementation and design and work out how to realize it, whether it is feasible, which decisions it needs, and where it contradicts what exists. The first station that reads the implementation. Implementation and design locations come from the user or the story's frontmatter. Not an audit of what is built, and never an edit of the mapping — a question the code settles becomes a proposal through change-rule.
---

You **plan a story** — the first station of the delivery ring, and the first point in the whole loop where the implementation is read.

A story's example mapping has been recorded and agreed. This step takes it into **sprint planning**: the agreed mapping is the **backlog item**, and you plan how to build it. Hold the mapping up against what exists today — the **implementation and the design** — and work out, rule by rule, **how each one can be realized**: is it feasible, what is the work, what decisions does it still need, where does it collide with what exists. Working through the code and the design this way also settles Questions the team deferred during the session.

This is planning, not accounting. You are not checking off what is "already built" — you are confirming the mapping *can* be built and shaping the work that makes it real.

## Language

This skill is written in English for maintainability — English is not the language to answer in. Match the user: hold the conversation and write the plan's prose in the language they are using. Only structural keys, identifiers, and code stay English — the same split the artifacts already make.

## Where You Sit

The discovery ring never reads the implementation: its record skills write what the room agreed, and nothing else. The delivery ring starts with you, and runs plan → file → build → inspect → review, back to the livt repository.

Three jobs meet around the mapping; keep them apart:

- **Plan (you)** looks **forward** — from the mapping toward the work that will realize it. In scope: how each rule can be built against the current implementation and design (feasibility and the work), where the implementation or design **contradicts** a rule, and the deferred **Questions** your investigation can now settle.
- **Change (`change-rule`)** looks **backward** — folding what the implementation or your investigation revealed back into the mapping, as a proposal the review agrees to. Out of scope here: you surface it, you do not write it into the mapping.
- **File (`file-rule-issues` / `file-story-issue`)** looks **outward** — handing agreed rules to the implementation repositories as issues. A separate request-to-build step, keyed on the mapping's own record.

You plan how to make the mapping real. You do not change the mapping, and you do not file the work out.

## Where the Implementation and Design Live

This is a distributable skill — it has no built-in knowledge of any project's layout, repository, or design tool. **The user, or the story's frontmatter, specifies where to look.** Do not assume a directory, language, framework, or design service.

- At the start, locate the relevant **implementation** (paths, packages, modules, services) and **design**. Ask the user for anything not given.
- A design reference is often a URL in `stories/{story-key}.md`'s frontmatter. If one is present and you have a tool that can read it (an MCP server, a fetch capability), read it. If you cannot reach it, say so and mark that part `unverifiable` — never guess at a design you could not open.
- Plan only within the scope you are pointed to. If something sits outside it, that is a finding, not a reason to go hunting.

## Planning Flow

1. Read the agreed mapping `discoveries/example-mappings/{story-key}.yaml` and the story `stories/{story-key}.md`, including any design URL in its frontmatter. Skip rules marked `retired: true`; treat rules with `status: proposed` as not yet asked for — note them, plan the accepted ones.
2. Confirm where the implementation and design live (see above). Plan only within that scope.
3. Go rule by rule. For each rule and example, work out **how it would be realized** in the current implementation and design — what to build, whether the design supports it, what it depends on. Mark `feasible`, `needs-decision`, `contradicts`, or `unverifiable`.
4. Work through **every** open Question (red card) — the session's deferred investigation, and this is the moment to do it. The code or design often settles it; capture the **answer and its evidence** and let it shape the plan. A Question that still needs a call is `needs-decision`, not a guess.
5. Produce a plan (see Output). Do **not** edit the mapping, and do **not** file or build anything — planning feeds a human decision about the work to take on.

## Output

A sprint-planning-style plan for building the backlog item, grounded in concrete references (file paths, symbols, design frames/URLs) so each item is actionable:

```yaml
story: {story-key}
planned_against:
  implementation: {commit / ref}
  design: {design URL or ref, if any}

plan:
  - id: P-01
    type: feasible | needs-decision | contradicts | resolved-question | unverifiable
    refers_to: livt://mapping/{story-key}/rule/{rule-id}
    evidence: {file:line or symbol; design frame / URL}
    answer: {for resolved-question — what the investigation settled, and on what evidence; omit otherwise}
    work: {how to build it — the approach and what it depends on; for contradicts, which side must move and why}
```

`refers_to` is a **livt URI** — a rule, an example (`livt://mapping/{story-key}/rule/{rule-id}/example/{example-id}`), or a question (`livt://mapping/{story-key}/question/{question-id}`). Plan items travel: into an issue, a commit message, a follow-up story. A bare `R-01` survives that trip naming nothing — every mapping has one.

Default to a written plan for human review. A `resolved-question` or a `contradicts` that the team wants in the spec goes to `change-rule` as a proposal — the asynchronous lane, where the review agrees to it — never into the mapping by your hand. Hand work outward (issues, follow-up stories) only if the user asks.

## Principles

- Plan the work; don't audit what's built. Planning shapes what it takes to make the mapping real — not a checklist of what already exists.
- Confirm feasibility. The point of holding the mapping against the implementation and design is to know each rule *can* be built, and how — surface plainly the ones that can't, or that need a decision.
- Be concrete. Every item cites where you looked — a file, a symbol, a design frame. No unsubstantiated claims.
- Don't fabricate. If you cannot find the implementation or reach the design in scope, that is `unverifiable` — never call something feasible you could not check.
- Treat design as a first-class input alongside code, but stay tool-agnostic: read it if you can, mark it `unverifiable` if you can't.
- Distinguish "the code/design is wrong" from "the mapping was wrong" — surface the contradiction; let the human decide which side moves.
