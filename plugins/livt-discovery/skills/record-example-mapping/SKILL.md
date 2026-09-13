---
name: record-example-mapping
description: Record an Example Mapping session into discoveries/example-mappings/{story-key}.yaml as two commits on one PR — the board's rules, examples, and questions verbatim as the baseline, then a structural edit consulting the bdd-expert skill, whose diff is what the review reads. Use after a session on a story's board, or when that board was reworked; it never changes what the room agreed. A rule proposed, changed, or added outside a session routes to change-rule; checking against the implementation is plan-story's.
---

You **record** an example mapping — the record station of the discovery ring, at the story level.

The Example Mapping session already happened: people worked it out together on a board (Miro, sticky notes, a photo), and the content was agreed in the room. Your job is not to facilitate, improve, or re-discover it. It is *talk, then record*: write the board out so the conversation leaves a faithful residue, then tidy its structure where that helps a reader — as a visible diff, never as a silent rewrite. In Ron Jeffries' three Cs, the story's card was written before the session and the session was the conversation; you write the **confirmation** — the rules and examples the team will hold the implementation to.

## Language

This skill is written in English for maintainability — English is not the language to answer in. Match the user: hold the conversation and write user-facing prose in the language they are using. Only structural keys, identifiers, and code stay English — the same split the artifacts already make.

## Two Commits, One PR

The record ships as one PR carrying two commits, in this order:

1. **Baseline** — the board verbatim, committed with no edits at all. This is the session's residue, and the reference every later diff is read against.
2. **Edit** — structure only, consulting the `bdd-expert` skill: rule clarity, example naming, grouping, question phrasing. The diff of this commit is the deliverable. A reviewer reads it to see exactly what the record changed, and to confirm that no meaning moved.

Never fold the edit into the baseline. A pre-polished baseline hides the edit, and the review loses the one thing the second commit exists to show. When the board needs nothing, the PR carries the baseline alone.

## Where You Sit

- Before you: `write-story-card` wrote `stories/{story-key}.md`, and the team held the session on that story's board.
- You: record the board against the card's key.
- After you: `plan-story` holds the agreed mapping against the implementation and design — the first read of the code, and the start of the delivery ring. A rule that changes later, or is proposed without a session, is `change-rule`'s.

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
- Faithfulness beats correctness here. A gap or an awkward phrasing is recorded as-is; the edit commit is where it may be addressed, in the open.

## Editing the Structure

1. Consult the `bdd-expert` skill for a structural critique — anti-patterns, rule/example balance, naming, whether the story should be split.
2. Apply the edit to the YAML:
   - **Rule clarity** — sharpen vague rule names into crisp business rules; keep the team's intent.
   - **Example naming** — make examples concrete and memorable ("the one where…"); keep the same scenario.
   - **Grouping** — example IDs are rule-scoped, so moving an example to the rule it actually illustrates always changes its ID. Move it only while nothing can be pointing at it: a just-recorded baseline where no rule carries `issues:` or `automated:`. Otherwise — and whenever you are unsure — retire it where it sits and add it under the right rule with a fresh ID, with `superseded_by:` on the retired one naming the new URI, so its old URI keeps resolving to the same text and says where the example went.
   - **Splitting** — if the map shows too many rules (story too large), recommend a split in the PR body; don't silently shard.
   - **Question phrasing** — make a Question precise without answering it.
3. Re-read the diff against the baseline: every change is structural, none is a meaning change.
4. Commit it as the edit, and open the PR.

### What the edit never touches

- Don't add rules or examples that weren't discovered.
- Don't fold in rule changes or additions decided after the session — `change-rule` ships those as their own fine-grained PRs.
- Don't resolve or delete open Questions.
- Don't check the mapping against the implementation or design — that is `plan-story`'s job.
- Don't write Gherkin — example mapping stays low-tech.
- Don't change the `story` key.
- Don't retire rules or questions. The regrouped example is the only retirement in your remit; retiring a rule takes an agreed business decision, which is `change-rule`'s.
- Don't change a rule's `status:`. Accepting or turning down a proposal is a business decision too, and `change-rule`'s.

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

Canonical statement in `change-rule`; these bullets are verbatim from it. You mint every ID in the mapping, so they start with you — and a tidier numbering is never a reason to renumber, a temptation the edit commit puts right in front of you.

- **Numbering** — a new ID is one past the highest ever used in its scope, **retired IDs included**. Rules and questions are numbered within the story (`R-NN`, `Q-NN`), examples within their rule (each rule starts from `EX-01`). With R-01/R-02/R-03 on file and R-03 retired, the next rule is R-04 — never R-03 again.
- **Immutability** — an ID, once used, keeps pointing at the same thing. Never renumber, never reuse, and never move an item to where its ID would change. This holds whether or not an automation issue was filed: the item's livt URI is quoted by MCP consumers, by the board's copy-link, in test comments, and in commit messages, and the livt repository records none of those — there is no list of references to check before breaking one.
- **Retire, don't delete** — an item that no longer holds gets `retired: true` and stays in the file, its ID taken and its text readable. Deleting it hands the ID to the next item and silently re-targets every reference. Don't comment it out either: a comment is not part of the YAML structure, so a structural edit drops it.
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
- A card that has left the board gets `retired: true` where it sits — gone from the spec, still in the file. A rule that has gone takes each of its examples with it: the flag does not cascade, so an unflagged example still resolves as live.
- An example the board moved under a different rule cannot keep its ID, since example IDs are rule-scoped. Retire it under the old rule and add it under the new one with a fresh ID, with `superseded_by:` on the retired one naming the new URI: the board's grouping is honoured, the old URI still resolves to the same text, and it says where the example went.

## Commit Contract

- Baseline: `Record {story-key} example mapping (baseline)`, committed with no edits.
- Edit: `Edit {story-key} example mapping structure`, on top of the baseline.
- One PR holding both. The PR body says what the edit changed and why, so the review can confirm that no meaning moved.

## Completeness Checks (record, not session health)

Before the baseline commit, verify the record is complete — not whether the map is "good" (that is the edit's question, and the session's before it):

- Every card on the board appears in the YAML; no card was dropped.
- Every green card sits under the same rule it sat under on the board.
- Every red card is preserved as a Question — none were silently answered.
- Every pink card appears in `ubiquitous:`, and every term whose definition the board carried has its `ubiquitous/{term-key}.md`.
- Wording matches the board; nothing was paraphrased away.
- Re-recording only: every ID that was already on file still names the same card, and nothing was deleted to make room.

## Output

`discoveries/example-mappings/{story-key}.yaml` on one PR with two commits: the verbatim baseline, and a structural edit whose diff the review reads. Plus the `ubiquitous/{term-key}.md` files the board's definitions produced, and a note of which referenced terms are still undefined.
