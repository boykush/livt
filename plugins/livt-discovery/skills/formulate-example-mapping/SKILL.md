---
name: formulate-example-mapping
description: Formulate a recorded Example Mapping — rule wording, example naming, grouping, and question phrasing — as one structural edit over a committed baseline, consulting the bdd-expert skill, whose diff is what the review reads. Use after record-example-mapping, or to work a mapping recorded long ago; it changes expression and structure only, never what the room agreed. A rule whose meaning changed routes to change-rule, and this skill never reads the implementation.
---

You **formulate** an example mapping — the formulation station of the discovery ring, at the story level.

The mapping is already on file: a session was recorded verbatim minutes ago, or the file has been sitting there for months. Your job is to make it a document every reader on the team reads the same behaviour out of — rules that assert, examples whose names say what they show. You change how the mapping is written. You never change what it says.

This is BDD's **Formulation**, the phase mainstream practice performs by writing Gherkin. livt derives Gherkin from the mapping rather than the other way round, so the act is the rewrite of the mapping itself — and the mapping stays the source of truth.

## Language

This skill is written in English for maintainability — English is not the language to answer in. Match the user: hold the conversation and write user-facing prose in the language they are using. Only structural keys, identifiers, and code stay English — the same split the artifacts already make.

## Where You Sit

- Before you: `record-example-mapping` wrote the board out verbatim and committed it. That commit is your baseline.
- You: lay a structural edit over it as its own commit. That diff is the deliverable — a reviewer reads it to see exactly what you changed, and to confirm that no meaning moved.
- After you: the agreed mapping is handed outward by `file-story-issue` and `file-rule-issues` — the start of the delivery ring. A rule whose *meaning* changes later is `change-rule`'s.

You are not once-only. A mapping formulated when it was recorded is formulated again whenever its wording turns out to read two ways — the same station, a fresh diff, no board required. `record-example-mapping` runs again only when a board does.

## The Baseline Is Committed

Check this before you touch anything: the mapping's last change is in a commit, and the working tree is clean of it. Everything below is defined against what is committed — that is what a baseline *is*.

If the record is still sitting uncommitted, stop and say so. Committing it is `record-example-mapping`'s last step, and folding your edit into it destroys the one thing this station exists to produce: a pre-polished baseline hides the edit, and the review can no longer see what moved.

Then read, in this order: `stories/{story-key}.md` for the story's scope and key — or, for a story with no card, the mapping's own `name:` — and the mapping as it stands — as it is, not as you would have recorded it.

## The Edit

1. Consult the `bdd-expert` skill for a structural critique — anti-patterns, rule/example balance, naming, whether the story should be split.
2. Apply the edit to `discoveries/example-mappings/{story-key}.yaml`:
   - **Rule clarity** — sharpen vague rule names into crisp business rules; keep the team's intent.
   - **Example naming** — make examples concrete and memorable ("the one where…"); keep the same scenario.
   - **Grouping** — example IDs are rule-scoped, so moving an example to the rule it actually illustrates always changes its ID. Move it only while nothing can be pointing at it: a just-recorded baseline where no rule carries `issues:` and no test cites it yet. Otherwise — and whenever you are unsure — retire it where it sits and add it under the right rule with a fresh ID, with `superseded_by:` on the retired one naming the new URI, so its old URI keeps resolving to the same text and says where the example went.
   - **Splitting** — if the map shows too many rules (story too large), recommend a split in the commit message; don't silently shard.
   - **Question phrasing** — make a Question precise without answering it.
3. Re-read the diff against the baseline: every change is structural, none is a meaning change. A line you cannot defend as expression is a line that belongs to `change-rule`.
4. Commit it (see Commit Contract).

## What the Edit Never Touches

- Don't add rules or examples that weren't discovered.
- Don't fold in rule changes or additions decided after the session — `change-rule` ships those as their own fine-grained commits.
- Don't resolve or delete open Questions.
- Don't check the mapping against the implementation or design — no discovery skill reads code.
- Don't write Gherkin — example mapping stays low-tech, and the Gherkin livt derives from the mapping is Automation's output, not yours.
- Don't change the mapping's key or its `name:`. The name is the story's, not a rule's, so sharpening it is no part of formulation — and where the story also has a card, the two have to stay the same string.
- Don't retire rules or questions. The regrouped example is the only retirement in your remit; retiring a rule takes an agreed business decision, which is `change-rule`'s.
- Don't change a rule's `status:`. Accepting or turning down a proposal is a business decision too, and `change-rule`'s.
- Don't touch `ubiquitous:`. The pink stickies are the room's and the record is their path off the board; sharpening a rule's wording is not a licence to name a term nobody agreed. Tell the user what the wording seems to want instead.

## ID Contract

Canonical statement in `change-rule`; these bullets are verbatim from it. Regrouping is the one place you mint an ID — and a tidier numbering is never a reason to renumber, a temptation this station puts right in front of you.

- **Numbering** — a new ID is one past the highest ever used in its scope, **retired IDs included**. Rules and questions are numbered within the story (`R-NN`, `Q-NN`), examples within their rule (each rule starts from `EX-01`). With R-01/R-02/R-03 on file and R-03 retired, the next rule is R-04 — never R-03 again.
- **Immutability** — an ID, once used, keeps pointing at the same thing. Never renumber, never reuse, and never move an item to where its ID would change. This holds whether or not an automation issue was filed: the item's livt URI is quoted by MCP consumers, by the board's copy-link, in test comments, and in commit messages, and the livt repository records none of those — there is no list of references to check before breaking one.
- **Retire, don't delete** — a rule that no longer holds gets a closed `status` (`rejected` or `retired`) and an example or question gets `retired: true`; either way it stays in the file, its ID taken and its text readable. Deleting it hands the ID to the next item and silently re-targets every reference. Don't comment it out either: a comment is not part of the YAML structure, so a structural edit drops it.
- **Say where the spec went** — when something took the closed item's place, add `superseded_by:` beside the status or the flag, listing the successors as **livt URIs** so a reference landing on the retired item reads on instead of stopping. A list, because an item can split into two; livt URIs, because a successor can live in another mapping and a bare `R-05` names nothing. Leave the field off when nothing replaced it. Only the pointer goes in the YAML — *why* it was retired belongs to the commit, where it is written once and cannot drift.

## Commit Contract

Canonical statement in `change-rule`; these bullets are verbatim from it.

- **The commit is livt's unit.** One decision per commit — a rule-level change, a record's baseline, the structural edit on top of it, a story's card. It is the smallest thing a reviewer can weigh on its own, and the only split livt asks for.
- **Never fold two decisions into one commit.** That is the one thing no later grouping can undo: a reviewer reading a combined diff cannot tell which change carried which reason, and neither can the history.
- **The message carries the why.** The subject names what changed; the body states the reason a reviewer weighs. It is written once, in the commit, where it cannot drift from the diff it explains.
- **How commits are grouped into pull requests is yours.** One per commit, one per session, one per chat thread — that is your team's branch and review convention, and no skill here has an opinion on it. Sending several up together is not batching, as long as each decision arrived as its own commit.

Yours is one commit: `Formulate {story-key} example mapping`, on top of the baseline. The body says what the edit changed and why, so the review can confirm that no meaning moved.

## Output

`discoveries/example-mappings/{story-key}.yaml`, as one structural edit over a committed baseline, whose diff is what the review reads. When the mapping needs nothing, commit nothing and say so — an empty formulation is a fine outcome, a disguised one is not.
