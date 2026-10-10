---
name: propose-rule
description: Carry one business rule from proposal to agreement in an existing Example Mapping — put it forward with status `proposed` and the people whose agreement it needs, then accept or reject it — each step as its own fine-grained commit. Use when a rule was found without the people who have to agree to it, as when one person works a story through with an AI — "propose this rule", "who has to agree to this", "accept the proposal". A rule already agreed routes to change-rule, a session's outcome to record-example-mapping, and this skill never reads the implementation.
---

You **propose a rule** and see it through to agreement.

Not every rule is found with everyone who has to agree to it in the room. One person works a story through with an AI, an agent notices a gap while building, someone thinks of a rule between sessions. The rule is concrete — it has examples — but agreement between people has not happened. Your job is to put that rule on the record as what it is, a **proposal**: `status: proposed`, and the people whose agreement it needs, so that the proposal says who it is waiting for. Later, when they have agreed or declined, you record that too. Each step is one **fine-grained commit**.

A team that maps its examples together, in one room, agrees before it records and has no use for you. `record-example-mapping` writes what they agreed, and `change-rule` changes it afterwards.

## Language

This skill is written in English for maintainability — English is not the language to answer in. Match the user: hold the conversation and write user-facing prose in the language they are using. Only structural keys, identifiers, and code stay English — the same split the artifacts already make.

## Not change-rule

- **`change-rule`** changes a rule the team has already agreed: adds one, rewords one, retires one. The agreement happened before it was called.
- **You** handle a rule whose agreement is still to come, from the moment it is put forward to the moment it is accepted or turned down.

A proposal that is accepted becomes an agreed rule, and from then on it is `change-rule`'s. With no candidate answer on the table there is nothing to propose: that is a Question card.

## Who Has to Agree

- **Ask; do not guess.** Who has to agree to a rule is the user's to say. Ask when they have not said, and write what they answer into `decision_makers`.
- **Name people, as the forge mentions them** — `"@alice"`. A check ties each entry to an identity mechanically, so write the handle and nothing else. A team cannot be matched to a person yet; name the people.
- **The author may be one of them.** Someone who decides a rule and also writes it is named like anyone else, and counts by having written the change.
- **Leave out everyone else.** People who should hear about the rule but do not decide it are not decision makers. An absence from the list says so.
- **It is optional.** A team that does not hold its rules to named people writes no `decision_makers`, and a proposal without it is still a proposal.

## Flow

1. Identify the affected story and read `discoveries/example-mappings/{story-key}.yaml` and `stories/{story-key}.md` — a story with no card has only its mapping, scoped by the mapping's own `name:`. Ask the user which mapping if it is ambiguous.
2. Capture what the user wants: a rule put forward, a proposal accepted, or a proposal turned down — and why. Don't invent or extrapolate.
3. Check whether what you touch has merged: read the default branch's copy of the mapping (`git show {default-branch}:discoveries/example-mappings/{story-key}.yaml`). A rule it does not hold was added on this branch and is still in review — a **draft** (see ID Contract), which is rewritten in place or deleted rather than closed.
4. Apply the one step (IDs throughout follow the ID Contract):
   - **Propose** — append the rule with the next rule ID, its examples, `status: proposed`, and `decision_makers` when the user names them. It goes on the board pale and waits.
   - **Name who has to agree** — add or change `decision_makers` on a proposal already on file. Nothing else about the rule moves.
   - **Accept** — change `status: proposed` to `status: accepted`, and nothing else. The diff is that one line, which is what the people named are asked to approve. If the agreement reworded the rule, the new wording lands in the same commit; a rule that came to mean something else is a different rule — reject this one, propose its successor under a new ID, and point the rejected one at it with `superseded_by:`.
   - **Reject** — change `status: proposed` to `status: rejected`, and mark its examples `retired: true` as for any closed rule. `rejected` rather than `retired` is what makes the record read as a proposal that was never agreed, not as a rule dropped later.
   - A question the proposal would answer stays open until the proposal is accepted, and is retired in that commit, pointing at the rule with `superseded_by:`.
5. Re-read the diff: it must contain exactly the one step.
6. Commit it on its own (see Commit Contract).

## Agreement Is Not Yours to Declare

- **Accepting writes a request, not a fact.** A commit that says `accepted` is what the people named are asked to approve. Where the repository runs `livt agreements check`, that commit does not merge until every one of them stands behind it; you neither run that judgment nor stand in for it.
- **Do not take a conversation as approval.** A message saying someone agreed is not their approval. Where a check is in place it will ask for the approval itself; where none is, accept only when the user tells you the people named have agreed.
- **Never write `agreements/`.** That record is written by the check from who actually wrote and approved a change. A name added by hand is refused by the same check.
- **Do not accept on a pull request that will rewrite for you.** Where the repository has the check write, a proposal everyone approved is turned accepted on its own branch. Leave it proposed and let the approvals arrive.

## ID Contract

Canonical statement in `change-rule`; these bullets are verbatim from it. A skill loads on its own, so every skill that can mint, move, or quote an ID has to carry them.

- **Numbering** — a new ID is one past the highest ever used in its scope, **retired IDs included**. Rules and questions are numbered within the story (`R-NN`, `Q-NN`), examples within their rule (each rule starts from `EX-01`). With R-01/R-02/R-03 on file and R-03 retired, the next rule is R-04 — never R-03 again.
- **Immutability** — an ID, once used, keeps pointing at the same thing. Never renumber, never reuse, and never move an item to where its ID would change. This holds whether or not an automation issue was filed: the item's livt URI is quoted by MCP consumers, by the board's copy-link, in test comments, and in commit messages, and the livt repository records none of those — there is no list of references to check before breaking one.
- **Retire, don't delete** — a rule that no longer holds gets a closed `status` (`rejected` or `retired`) and an example or question gets `retired: true`; either way it stays in the file, its ID taken and its text readable. Deleting it hands the ID to the next item and silently re-targets every reference. Don't comment it out either: a comment is not part of the YAML structure, so a structural edit drops it.
- **Say what replaced it** — when something took the closed item's place, add `superseded_by:` beside the status or the flag, listing the successors as **livt URIs** so a reference landing on the retired item reads on instead of stopping. A list, because an item can split into two; livt URIs, because a successor can live in another mapping and a bare `R-05` names nothing. Leave the field off when nothing replaced it. Only the pointer goes in the YAML — *why* it was retired belongs to the commit, where it is written once and cannot drift.
- **A draft is rewritten, not retired** — the bullets above bind an ID once its item reaches the default branch: the branch the published site is built from and the implementation repositories read, and so the first place anything outside the review can point at it. An item this branch added and the default branch does not hold yet is a draft. It is edited like any other text in review — reworded past its meaning, moved, or deleted, a deleted draft giving its number back — and never retired: `retired` says a decision was agreed and then withdrawn, and on an item the default branch never held it writes the review's back-and-forth into the record as history. The default branch's copy of the file is the test: an ID it holds is bound wherever the item is changed next, and one you cannot check is treated as bound too.

## Commit Contract

Canonical statement in `change-rule`; these bullets are verbatim from it. A skill loads on its own, so every skill that writes a commit has to carry them.

- **The commit is livt's unit.** One decision per commit — a rule-level change, a record's baseline, the structural edit on top of it, a story's card. It is the smallest thing a reviewer can weigh on its own, and the only split livt asks for.
- **Never fold two decisions into one commit.** That is the one thing no later grouping can undo: a reviewer reading a combined diff cannot tell which change carried which reason, and neither can the history.
- **The message carries the why.** The subject names what changed; the body states the reason a reviewer weighs. It is written once, in the commit, where it cannot drift from the diff it explains.
- **How commits are grouped into pull requests is yours.** One per commit, one per session, one per chat thread — that is your team's branch and review convention, and no skill here has an opinion on it. Sending several up together is not batching, as long as each decision arrived as its own commit.

For you, that unit is one step in a proposal's life:

- The commit message names the rule and the mapping, e.g. `Propose rule R-05 in {story-key}`, `Accept rule R-05 in {story-key}`, `Reject rule R-05 in {story-key}`.
- The body states the business reason. For a proposal it is the opening of the conversation; for an acceptance it is why the people named agreed.

## What NOT to Do

- Don't propose a rule the team has already agreed — that is `change-rule`'s, and it takes no status.
- Don't guess who has to agree, and don't name a team.
- Don't accept a proposal on your own reading of a conversation, and don't write the agreement record.
- Don't fold two steps into one commit: proposing two rules is two commits, and so is accepting two.
- Don't check the proposal against the implementation — discovery never reads code.

## Output

The same `discoveries/example-mappings/{story-key}.yaml` with exactly one step applied to one rule, shipped as its own fine-grained commit. Keep `key:` identifiers in English; `name:`/`text:` follow the mapping's language.
