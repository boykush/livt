---
name: record-example-mapping
description: Record an Example Mapping session into discoveries/example-mappings/{story-key}.yaml — the board's rules, examples, and questions committed verbatim as the baseline every later diff is read against. Use after a session on a story's board, or when that board was reworked; it never changes what the room agreed, and never tidies what it wrote — sharpening the mapping is formulate-example-mapping's commit, laid over yours. A rule proposed, changed, or added outside a session routes to change-rule; this skill never reads the implementation.
---

You **record** an example mapping — the record station of the discovery ring, at the story level.

The Example Mapping session already happened: people worked it out together on a board (Miro, sticky notes, a photo), and the content was agreed in the room. Your job is not to facilitate, improve, or re-discover it. It is *talk, then record*: write the board out so the conversation leaves a faithful residue. Tidying it for a reader is the next station's work, and it lands as a visible diff over what you commit, never as a silent rewrite. In Ron Jeffries' three Cs, the story's card was written before the session and the session was the conversation; you write the **confirmation** — the rules and examples the team will hold the implementation to.

## Language

This skill is written in English for maintainability — English is not the language to answer in. Match the user: hold the conversation and write user-facing prose in the language they are using. Only structural keys, identifiers, and code stay English — the same split the artifacts already make.

## Where the Record Stops

You ship **one commit**: the board verbatim, with no edits at all. It is the session's residue, and the reference every later diff is read against.

Making the mapping read well — rule clarity, example naming, grouping, question phrasing — is the next station, and `formulate-example-mapping` ships it as its own commit consulting `bdd-expert`. Never do that work here. A pre-polished baseline hides the edit, and the review loses the one thing that second commit exists to show. When you are done, say the baseline is in and that formulating it is the next step.

## Where You Sit

- Before you: `write-story-card` wrote `stories/{story-key}.md`, and the team held the session on that story's board.
- You: record the board against the card's key.
- After you: `formulate-example-mapping` lays the structural edit over your baseline, and the agreed mapping is handed outward from there by `file-story-issue` and `file-rule-issues` — the start of the delivery ring. A rule that changes later, or is proposed without a session, is `change-rule`'s.

## Recording the Baseline

1. Read the story from `stories/{story-key}.md` to get the scope and key.
2. Take in the board — a screenshot/photo, exported text, or pasted sticky-note contents. Ask the user for it if it was not provided.
3. Walk the board card by card and capture every card, by color:
   - Yellow (Story) → the `story` key
   - Blue (Rule) → a `rules[]` entry
   - Green (Example) → an `examples[]` entry under its rule
   - Red (Question) → a `questions[]` entry
   - Pink (Ubiquitous language) → an `ubiquitous[]` term reference (see Ubiquitous Language below)
4. Write to `discoveries/example-mappings/{story-key}.yaml`.
5. Read the YAML back against the board and confirm nothing was dropped or altered.
6. Commit it as the baseline (see Commit Contract).

### Baseline principles

- Capture every card. Dropping a card silently is the worst failure mode — when unsure whether something is a card, include it.
- Keep the board's wording. Record the team's phrasing, not your paraphrase.
- Keep examples attached to the rule they sat under on the board. Don't re-file them.
- Don't merge near-duplicate cards — if the board has two, the YAML has two.
- Preserve disagreement. If something was left as a Question (red card), it stays a Question — never answer it.
- Record agreement as the board shows it. A rule the board itself marks as not yet agreed takes `status: proposed`; every other rule was agreed in the room and takes no status. Never mark one proposed on your own reading of it.
- If the board is genuinely ambiguous (illegible card, an example with no rule), ask the user rather than guessing.
- Stay scoped to this story's board. New stories noted on the board are captured separately, not folded in here.
- Faithfulness beats correctness here. A gap or an awkward phrasing is recorded as-is; `formulate-example-mapping` is where it may be addressed, in the open.

## YAML Format

```yaml
story: {story-key}

rules:
  - id: R-01
    name: {rule description, as written on the board}
    examples:
      - id: EX-01
        name: {concrete example, as written on the board}

questions:
  - id: Q-01
    text: {unresolved question, as written on the board}

ubiquitous:
  - {term-key}
  - {ctx}/{term-key}
```

- IDs follow the ID Contract below. A board recorded for the first time simply numbers from 01 in each scope.
- Keep `key:` identifiers in English; `name:`/`text:` follow the board's language.
- `ubiquitous` lists the board's pink stickies as term references (see Ubiquitous Language below).

## ID Contract

Canonical statement in `change-rule`; these bullets are verbatim from it. You mint every ID in the mapping, so they start with you — and a tidier numbering is never a reason to renumber.

- **Numbering** — a new ID is one past the highest ever used in its scope, **retired IDs included**. Rules and questions are numbered within the story (`R-NN`, `Q-NN`), examples within their rule (each rule starts from `EX-01`). With R-01/R-02/R-03 on file and R-03 retired, the next rule is R-04 — never R-03 again.
- **Immutability** — an ID, once used, keeps pointing at the same thing. Never renumber, never reuse, and never move an item to where its ID would change. This holds whether or not an automation issue was filed: the item's livt URI is quoted by MCP consumers, by the board's copy-link, in test comments, and in commit messages, and the livt repository records none of those — there is no list of references to check before breaking one.
- **Retire, don't delete** — a rule that no longer holds gets `status: retired` and an example or question gets `retired: true`; either way it stays in the file, its ID taken and its text readable. Deleting it hands the ID to the next item and silently re-targets every reference. Don't comment it out either: a comment is not part of the YAML structure, so a structural edit drops it.
- **Say where the spec went** — when something took the retired item's place, add `superseded_by:` beside `retired:`, listing the successors as **livt URIs** so a reference landing on the retired item reads on instead of stopping. A list, because an item can split into two; livt URIs, because a successor can live in another mapping and a bare `R-05` names nothing. Leave the field off when nothing replaced it. Only the pointer goes in the YAML — *why* it was retired belongs to the commit, where it is written once and cannot drift.

## Ubiquitous Language

Pink stickies are the words the room agreed on, and the record is their only path off the board. A term that never reaches `ubiquitous/` can never be looked up. `record-story-map` carries this section verbatim — a skill loads on its own, so both record skills have to hold it. Change one, change both.

- Capture every pink sticky as an `ubiquitous:` entry. The list holds **term references**, not display names: kebab-case English, `{ctx}/{term-key}` when the board scopes the word to a context. Mint the key from the term, and ask the user when the wording gives no obvious one.
- When the board carries the term's agreed **definition**, write `ubiquitous/{term-key}.md` with the board's word as `name:` and the definition as the body.
- When it carries only the word, capture the reference and tell the user which terms are still undefined. An unknown key renders as a plain pink card, so a reference without a file is a safe and visible state — inventing a definition is not.
- A term already in `ubiquitous/` keeps its definition. The board naming a word does not license rewriting what the glossary already says about it.

## Re-recording an Existing Mapping

When the team reworks a board that has already been recorded, `discoveries/example-mappings/{story-key}.yaml` already holds IDs. Record onto that file — never write a fresh one over it, which would re-mint every ID:

- A card already on file keeps its ID, even where the board reworded it. Follow the board for the text; leave the ID alone.
- A card that is new on the board takes the next ID.
- A card that has left the board is closed where it sits — `status: retired` on a rule, `retired: true` on an example or a question — gone from the spec, still in the file. A rule that has gone takes each of its examples with it: neither spelling cascades, so an unflagged example still resolves as live.
- An example the board moved under a different rule cannot keep its ID, since example IDs are rule-scoped. Retire it under the old rule and add it under the new one with a fresh ID, with `superseded_by:` on the retired one naming the new URI: the board's grouping is honoured, the old URI still resolves to the same text, and it says where the example went.

## Commit Contract

Canonical statement in `change-rule`; these bullets are verbatim from it.

- **The commit is livt's unit.** One decision per commit — a rule-level change, a record's baseline, the structural edit on top of it, a story's card. It is the smallest thing a reviewer can weigh on its own, and the only split livt asks for.
- **Never fold two decisions into one commit.** That is the one thing no later grouping can undo: a reviewer reading a combined diff cannot tell which change carried which reason, and neither can the history.
- **The message carries the why.** The subject names what changed; the body states the reason a reviewer weighs. It is written once, in the commit, where it cannot drift from the diff it explains.
- **How commits are grouped into pull requests is yours.** One per commit, one per session, one per chat thread — that is your team's branch and review convention, and no skill here has an opinion on it. Sending several up together is not batching, as long as each decision arrived as its own commit.

Yours is one: `Record {story-key} example mapping (baseline)`, committed with no edits. The structural edit that follows is a second unit, and a second skill's.

## Completeness Checks (record, not session health)

Before the baseline commit, verify the record is complete — not whether the map is "good" (that is `formulate-example-mapping`'s question, and the session's before it):

- Every card on the board appears in the YAML; no card was dropped.
- Every green card sits under the same rule it sat under on the board.
- Every red card is preserved as a Question — none were silently answered.
- Every pink card appears in `ubiquitous:`, and every term whose definition the board carried has its `ubiquitous/{term-key}.md`.
- Wording matches the board; nothing was paraphrased away.
- Re-recording only: every ID that was already on file still names the same card, and nothing was deleted to make room.

## Output

`discoveries/example-mappings/{story-key}.yaml` as one commit: the verbatim baseline every later diff is read against. Plus the `ubiquitous/{term-key}.md` files the board's definitions produced, a note of which referenced terms are still undefined, and the hand-off to `formulate-example-mapping`.
