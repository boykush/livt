---
name: change-rule
description: Change one business rule in an existing Example Mapping — propose, accept, reject, change, or retire it, with its examples — as its own fine-grained commit. This is discovery's asynchronous lane — a rule enters the mapping as a proposal carrying status `proposed`, without a session, and the review stands in for the conversation. Use for "this rule changed", "a rule was added", "propose this rule"; a session's outcome routes to record-example-mapping, and this skill never reads the implementation.
---

You **change a rule** — the asynchronous lane of the discovery ring.

The conversation in the room is not the only way a rule reaches the record. Rules change in ongoing work, agents find gaps while planning or building, and people think of rules between sessions. Your job is to fold exactly one **rule-level change** into the existing `discoveries/example-mappings/{story-key}.yaml` and ship it as a **fine-grained commit** — one business decision per commit, so each stays reviewable on its own. When the change is not agreed yet, it enters as a **proposal**: `status: proposed`, examples and all, and the review is where the team talks it over. Agreement is then a one-line diff.

## Language

This skill is written in English for maintainability — English is not the language to answer in. Match the user: hold the conversation and write user-facing prose in the language they are using. Only structural keys, identifiers, and code stay English — the same split the artifacts already make.

## Not the Record Step

This is the common mix-up, so settle it first:

- **The record (`record-example-mapping`)** writes a session's outcome: the board, verbatim. **The formulation (`formulate-example-mapping`)** sharpens how that record reads, over the baseline it committed. Neither may change agreed meaning — neither can add, change, or retire rules.
- **You** change meaning on purpose, because the business changed it or is being asked to. You apply exactly that change — no more, no less.

When there is no fresh session and the ask is "this rule changed / a new rule was added / this rule was retired / propose this rule", you are the right skill, not the record.

## Two Lanes Into One Record

- **The synchronous lane** — people around a board at the same time, then `record-example-mapping` and `formulate-example-mapping`. This is Example Mapping as designed: a structured conversation.
- **The asynchronous lane (you)** — a proposal lands on the board pale and dashed, the Tasks page lists it under Proposed Rules ("closed by agreement"), and the review accepts it, reworks it, or turns it down. In Ron Jeffries' three Cs the proposal is the card, the review thread — wherever your team holds it — is the conversation, and the one-line acceptance is the confirmation.
- **Where proposals come from** — `inspect-automation` finding a rule the implementation has drifted from; an agent in an implementation repository reading the spec through `livt-automation` and noticing a gap; a person between sessions. Record all of these as proposals unless the user says the business already agreed.
- **When the lane stalls** — a proposal nobody agrees to asynchronously stays on the Tasks page, beside the open questions, as the agenda of the next session. The two lanes do not compete: the asynchronous one runs ahead, the synchronous one catches what it leaves.

## Change Philosophy

- **One rule-level change per commit.** A rule proposed, added, changed, or retired — together with the examples that illustrate it. A second rule-level change is a second commit, even when it touches the same mapping and arrived in the same conversation.
- **You record decisions; you don't make them.** The decision happened in ongoing work — a conversation, a ticket, an incident — or is being put forward for one. Capture what was decided, in the user's words. Still being debated? With no candidate answer on the table it is a Question card; a concrete rule put forward for agreement is a rule with `status: proposed`.
- **The diff is the deliverable.** A reviewer must see exactly one business decision in the commit. Don't mix in structural tidying — that noise belongs to `formulate-example-mapping`, if it is ever needed.
- **The mapping may lead the implementation.** After your change lands, the code may not match the mapping yet. That gap is the delivery ring's business, not a reason to hold the change back.

## Change Flow

1. Identify the affected story and read `discoveries/example-mappings/{story-key}.yaml` and `stories/{story-key}.md`. Ask the user which mapping if it is ambiguous.
2. Capture the change from the user: which rule, what changed, and why — and whether it is agreed or only put forward. Don't invent or extrapolate.
3. Apply the minimal rule-level edit (IDs throughout follow the ID Contract):
   - **Added rule** — append it with the next rule ID, with the examples agreed alongside it.
   - **Changed rule** — update its `name`, and bring its examples in line with the new meaning in the same commit (an example illustrating the old rule is now wrong). Retire the ones that no longer illustrate it and add the replacements as new examples, pointing each retired one at its replacement with `superseded_by:`; don't rewrite an example into a different one under the same ID. If the rule carries `automated:`, remove the flag in the same commit — the recorded automation covered the old meaning, and it is set again once the implementation catches up.
   - **Retired rule** — set `status: retired`, and mark every one of its examples `retired: true`: the status does not cascade, so an unflagged example under a retired rule still resolves as live. Nothing is deleted. Where another rule took over, add `superseded_by:`; where the business simply stopped asking, leave it off.
   - **Proposed rule** — append it like an added rule, with `status: proposed` beside its `name`. It is a candidate, not spec: it goes on the board pale and waits for agreement.
   - **Accepted proposal** — change its `status: proposed` to `status: accepted`. If the agreement reworded it, the new wording lands in the same commit, as for a changed rule.
   - **Rejected proposal** — change its `status: proposed` to `status: rejected`, and mark its examples `retired: true` as for any closed rule. `rejected` rather than `retired` is what makes the record read as a proposal that never became spec, not as a rule dropped later.
   - New open questions raised by the change take the next question ID; a question this change settles is retired, with the settled meaning landing as a rule — not as an answer scribbled onto the Question. Point the retired question at that rule with `superseded_by:`: the rule is where its answer lives. A proposal settles nothing yet: a question it would answer stays open until the proposal is accepted, and is retired in that commit.
4. Re-read the diff: it must contain exactly the one change, nothing structural elsewhere.
5. Commit it on its own (see Commit Contract).
6. More than one rule changed? Repeat the flow — one commit each.

## ID Contract

This is the canonical statement. `record-example-mapping` and `formulate-example-mapping` repeat these bullets verbatim, and `file-rule-issues` the immutability half — a skill loads on its own, so every skill that can mint, move, or quote an ID has to carry them. Change one, change all.

- **Numbering** — a new ID is one past the highest ever used in its scope, **retired IDs included**. Rules and questions are numbered within the story (`R-NN`, `Q-NN`), examples within their rule (each rule starts from `EX-01`). With R-01/R-02/R-03 on file and R-03 retired, the next rule is R-04 — never R-03 again.
- **Immutability** — an ID, once used, keeps pointing at the same thing. Never renumber, never reuse, and never move an item to where its ID would change. This holds whether or not an automation issue was filed: the item's livt URI is quoted by MCP consumers, by the board's copy-link, in test comments, and in commit messages, and the livt repository records none of those — there is no list of references to check before breaking one.
- **Retire, don't delete** — a rule that no longer holds gets a closed `status` (`rejected` or `retired`) and an example or question gets `retired: true`; either way it stays in the file, its ID taken and its text readable. Deleting it hands the ID to the next item and silently re-targets every reference. Don't comment it out either: a comment is not part of the YAML structure, so a structural edit drops it.
- **Say where the spec went** — when something took the closed item's place, add `superseded_by:` beside the status or the flag, listing the successors as **livt URIs** so a reference landing on the retired item reads on instead of stopping. A list, because an item can split into two; livt URIs, because a successor can live in another mapping and a bare `R-05` names nothing. Leave the field off when nothing replaced it. Only the pointer goes in the YAML — *why* it was retired belongs to the commit, where it is written once and cannot drift.

## Commit Contract

This is the canonical statement. `record-story-map`, `record-example-mapping`, `formulate-example-mapping`, `write-story-card`, and `inspect-automation` repeat these bullets verbatim — a skill loads on its own, so every skill that writes a commit has to carry them. Change one, change all.

- **The commit is livt's unit.** One decision per commit — a rule-level change, a record's baseline, the structural edit on top of it, a story's card. It is the smallest thing a reviewer can weigh on its own, and the only split livt asks for.
- **Never fold two decisions into one commit.** That is the one thing no later grouping can undo: a reviewer reading a combined diff cannot tell which change carried which reason, and neither can the history.
- **The message carries the why.** The subject names what changed; the body states the reason a reviewer weighs. It is written once, in the commit, where it cannot drift from the diff it explains.
- **How commits are grouped into pull requests is yours.** One per commit, one per session, one per chat thread — that is your team's branch and review convention, and no skill here has an opinion on it. Sending several up together is not batching, as long as each decision arrived as its own commit.

For you, that unit is one rule-level change:

- The commit message names the rule and the mapping, e.g. `Add rule R-04 to {story-key}`, `Update rule R-02 in {story-key}`, `Retire rule R-03 in {story-key}`, `Propose rule R-05 in {story-key}`, `Accept rule R-05 in {story-key}`.
- The body states the business reason for the change — that context is what a reviewer weighs, and for a proposal it is the opening of the conversation.

## What NOT to Do

- Don't restructure, rename, or re-file anything outside the agreed change — that is `formulate-example-mapping`'s, on its own diff.
- Don't add rules or examples beyond what was agreed or proposed, and don't answer open Questions in passing.
- Don't accept a proposal on your own reading of a conversation. `status: accepted` records that the business agreed, so it takes the user saying so.
- Don't check the change against the implementation — discovery never reads code; the gap closes on the delivery side, after the change lands.
- Don't fold two rule-level changes into one commit, even if they arrived together.

## Output

The same `discoveries/example-mappings/{story-key}.yaml` with exactly one rule-level change applied, shipped as its own fine-grained commit. Keep `key:` identifiers in English; `name:`/`text:` follow the mapping's language.
