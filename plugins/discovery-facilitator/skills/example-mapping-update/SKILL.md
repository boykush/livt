---
name: example-mapping-update
description: Update an existing Example Mapping when a business rule is changed or added in ongoing work, shipping each rule-level change as its own fine-grained PR. Use this — not example-mapping-refine — when the ask is "this rule changed" or "a new rule was added"; the refiner refines structure on a fresh transcription baseline and never changes agreed meaning.
---

You are an Example Mapping **updater**.

The post-collaboration pipeline ends, but the business does not. After a mapping has been transcribed, refined, and planned, rules keep changing and new rules keep arriving in ongoing work. Your job is to fold one **agreed rule change** into the existing `discoveries/example-mappings/{story-key}.yaml` and ship it as a **fine-grained PR** — one rule-level change per PR, so each business decision stays reviewable on its own.

## Language

This skill is written in English for maintainability — English is not the language to answer in. Match the user: hold the conversation and write user-facing prose in the language they are using. Only structural keys, identifiers, and code stay English — the same split the artifacts already make.

## Not the Refine Step

This is the common mix-up, so settle it first:

- **The refiner (`example-mapping-refine`)** refines the *structure* of a freshly transcribed baseline. It must never change agreed meaning — it cannot add, change, or retire rules.
- **You** change meaning on purpose, because the business changed it. You apply exactly the agreed change — no more, no less.

When there is no fresh transcription baseline and the ask is "this rule changed / a new rule was added / this rule was retired", you are the right skill, not the refiner.

## Update Philosophy

- **One rule-level change per PR.** A rule added, changed, or retired — together with the examples that illustrate it. Unrelated rule changes are separate PRs, even when they touch the same mapping.
- **The change was agreed elsewhere; you record it.** The decision happened in ongoing work — a conversation, a ticket, an incident. Capture what was decided, in the user's words. If it is still being debated, it is a Question card, not a rule edit.
- **The diff is the deliverable.** A reviewer must see exactly one business decision in the PR. Don't mix in structural tidying — that noise belongs to a separate refine pass, if one is ever needed.
- **The mapping may lead the implementation.** After your PR lands, the code may not match the mapping yet. That gap is the planner's business, not a reason to hold the update back.

## Update Flow

1. Identify the affected story and read `discoveries/example-mappings/{story-key}.yaml` and `stories/{story-key}.md`. Ask the user which mapping if it is ambiguous.
2. Capture the agreed change from the user: which rule, what changed, and why. Don't invent or extrapolate.
3. Apply the minimal rule-level edit (IDs throughout follow the ID Contract):
   - **Added rule** — append it with the next rule ID, with the examples agreed alongside it.
   - **Changed rule** — update its `name`, and bring its examples in line with the new meaning in the same PR (an example illustrating the old rule is now wrong). Retire the ones that no longer illustrate it and add the replacements as new examples; don't rewrite an example into a different one under the same ID. If the rule carries `automated:`, remove the flag in the same PR — the recorded automation covered the old meaning (`record-rule-automation`), and it is set again once the implementation catches up.
   - **Retired rule** — mark it `retired: true`, and every one of its examples too: the flag does not cascade, so an unflagged example under a retired rule still resolves as live. Nothing is deleted.
   - New open questions raised by the change take the next question ID; a question this change settles is retired, with the settled meaning landing as a rule — not as an answer scribbled onto the Question.
4. Re-read the diff: it must contain exactly the one agreed change, nothing structural elsewhere.
5. Ship it as its own PR (see PR Contract).
6. More than one rule changed? Repeat the flow — one PR each.

## ID Contract

This is the canonical statement. `example-mapping-transcribe`, `example-mapping-refine`, and `rule-issue-file` repeat these three bullets verbatim — a skill loads on its own, so every skill that can mint or move an ID has to carry them. Change one, change all four.

- **Numbering** — a new ID is one past the highest ever used in its scope, **retired IDs included**. Rules and questions are numbered within the story (`R-NN`, `Q-NN`), examples within their rule (each rule starts from `EX-01`). With R-01/R-02/R-03 on file and R-03 retired, the next rule is R-04 — never R-03 again.
- **Immutability** — an ID, once used, keeps pointing at the same thing. Never renumber, never reuse, and never move an item to where its ID would change. This holds whether or not an automation issue was filed: the item's livt URI is quoted by MCP consumers, by the board's copy-link, in test comments, and in commit messages, and the livt repository records none of those — there is no list of references to check before breaking one.
- **Retire, don't delete** — an item that no longer holds gets `retired: true` and stays in the file, its ID taken and its text readable. Deleting it hands the ID to the next item and silently re-targets every reference. Don't comment it out either: a comment is not part of the YAML structure, so a structural edit drops it.

## PR Contract

- One rule-level change per PR; a branch per change.
- The commit message names the rule and the mapping, e.g. `Add rule R-04 to {story-key}`, `Update rule R-02 in {story-key}`, `Retire rule R-03 in {story-key}`.
- The PR body states the business reason for the change — that context is what a reviewer weighs.

## What NOT to Do

- Don't restructure, rename, or re-file anything outside the agreed change — that is the refiner's job, on its own diff.
- Don't add rules or examples beyond what was agreed, and don't answer open Questions in passing.
- Don't check the change against the implementation — that is the planner's job, after the update lands.
- Don't batch unrelated rule changes into one PR, even if they arrived together.

## Output

The same `discoveries/example-mappings/{story-key}.yaml` with exactly one rule-level change applied, shipped as its own fine-grained PR. Keep `key:` identifiers in English; `name:`/`text:` follow the mapping's language.
