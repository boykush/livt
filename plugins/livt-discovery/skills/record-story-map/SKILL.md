---
name: record-story-map
description: Record a User Story Mapping session into discoveries/usm/{map-name}.yaml as two commits on one PR — the board as it stands, committed verbatim as the baseline, then a structural edit consulting the usm-expert skill, whose diff is what the review reads. Use after a story-mapping session, or when a mapped board was reworked; it never changes what the room agreed. A story's card routes to write-story-card, a rule-level change to change-rule.
---

You **record** a story map — the record station of the discovery ring.

The mapping session already happened: people built the map together on a board (Miro, sticky notes, a photo), and the narrative was worked out in the room. Your job is not to facilitate, challenge scope, or re-discover the journey. It is what Patton calls *talk, then record*: write the board out so the conversation leaves a faithful residue, then tidy its structure where that helps a reader — as a visible diff, never as a silent rewrite.

## Language

This skill is written in English for maintainability — English is not the language to answer in. Match the user: hold the conversation and write user-facing prose in the language they are using. Only structural keys, identifiers, and code stay English — the same split the artifacts already make.

## Two Commits, One PR

The record ships as one PR carrying two commits, in this order:

1. **Baseline** — the board verbatim, committed with no edits at all. This is the session's residue, and the reference every later diff is read against.
2. **Edit** — structure only, consulting the `usm-expert` skill: narrative flow, granularity, story framing, ordering. The diff of this commit is the deliverable. A reviewer reads it to see exactly what the record changed, and to confirm that no meaning moved.

Never fold the edit into the baseline. A pre-polished baseline hides the edit, and the review loses the one thing the second commit exists to show. When the board needs nothing, the PR carries the baseline alone — an empty edit is a fine outcome, a disguised one is not.

## Where You Sit

The discovery ring runs opportunity → story map → card → conversation → record, and back to the conversation through what remains open. You are the record of the story-map level.

- Before you: the session on the board.
- After you: `write-story-card` takes a candidate on the map into its own conversation (an Example Mapping) by writing its card and stamping its `key:` back onto the map.
- Later: a rule that changes in ongoing work is `change-rule`'s, on the story's own mapping. You never touch the example mappings.

## Recording the Baseline

1. Take in the board — a screenshot/photo, exported text, or pasted sticky-note contents. Ask the user for it if it was not provided.
2. Walk the backbone left to right and capture, in order:
   - **Activities** (top row) → `activities[]`
   - **User tasks** (second row, under each activity) → `steps[]`
   - **Story cards** (hanging below a task) → `stories[]` under that step, top-to-bottom as on the board
   - **Release slices** (the horizontal dividers cutting across the whole board) → `releases[]`, top-to-bottom, and stamp every card sitting above a divider with that slice's `release:`
3. Capture the pink stickies — the terms the room agreed on — as `ubiquitous:` term references (see Ubiquitous Language below).
4. Write to `discoveries/usm/{map-name}.yaml`.
5. Read the YAML back against the board and confirm nothing was dropped, reordered, or altered.
6. Commit it as the baseline (see Commit Contract).

### Baseline principles

- Capture every card and keep its position. Order is meaning in a story map — preserve left-to-right backbone order and top-to-bottom story priority.
- Keep the board's wording. Record the team's phrasing, not your paraphrase.
- Keep stories hanging under the same task they sat under on the board. Don't re-file them.
- Don't merge near-duplicate cards or "tidy" the granularity. If the board has two, the YAML has two.
- Assign a `key:` only when the board clearly intends a stable identifier; otherwise capture the story with `name:` alone.
- If the board is genuinely ambiguous, ask the user rather than guessing.
- Faithfulness beats correctness here. An awkward seam or a thin column is recorded as-is; the edit commit is where it may be addressed, in the open.

## Editing the Structure

1. Consult the `usm-expert` skill for a structural critique — does the backbone read as a coherent narrative, is granularity consistent, do stories carry clear user value?
2. Apply the edit to the YAML:
   - **Narrative flow** — sharpen activity and task names so the backbone reads left-to-right as a coherent user journey.
   - **Granularity** — lift a too-detailed task up or break a too-big one down so the backbone stays at a consistent level.
   - **Story framing** — phrase stories as user-valuable units; keep the same intent.
   - **Ordering** — fix backbone or priority order only where the board's order was clearly an artifact of the room, not a decision.
3. Re-read the diff against the baseline: every change is structural, none is a scope change.
4. Commit it as the edit, and open the PR.

### What the edit never touches

- Don't add activities, tasks, or stories that weren't mapped. Discovering more is the session's job.
- Don't re-slice releases or re-prioritize against the team's decisions — raise the concern in the PR body instead.
- Don't mint a `key:` the baseline did not intend; writing a card is `write-story-card`'s job.
- Keep `key:` identifiers in **English** (they back `stories/{key}.md` filenames and `step` cross-references); `name:` follows the board's language.

## YAML Format

```yaml
name: {map name}

activities:
  - key: {activity-key}
    name: {activity name}
    steps:
      - key: {step-key}
        name: {user task name}
        stories:
          - name: {story card, as written on the board}
          - key: {story-key}
            name: {story card with a stable key}
            release: {release-id}

releases:
  - id: {release-id}
    name: {slice name, as written on the board}

ubiquitous:
  - {term-key}
  - {ctx}/{term-key}
```

- `key:` and `id:` identifiers stay in **English** (story keys are used as filenames `stories/{key}.md`); `name:` follows the board's language.
- A story card may be captured with only `name:` when it is still lightweight and has no stable key yet.
- `releases:` lists the slices top-to-bottom, matching the board. A card carries the `release:` of the divider it sits above; a card below every divider carries none. A card without a `key:` can still carry a `release:`.
- `ubiquitous` lists the board's pink stickies as term references (see Ubiquitous Language below).

## Ubiquitous Language

Pink stickies are the words the room agreed on, and the record is their only path off the board. A term that never reaches `ubiquitous/` can never be looked up. `record-example-mapping` carries this section verbatim — a skill loads on its own, so both record skills have to hold it. Change one, change both.

- Capture every pink sticky as an `ubiquitous:` entry. The list holds **term references**, not display names: kebab-case English, `{ctx}/{term-key}` when the board scopes the word to a context. Mint the key from the term, and ask the user when the wording gives no obvious one.
- When the board carries the term's agreed **definition**, write `ubiquitous/{term-key}.md` with the board's word as `name:` and the definition as the body.
- When it carries only the word, capture the reference and tell the user which terms are still undefined. An unknown key renders as a plain pink card, so a reference without a file is a safe and visible state — inventing a definition is not.
- A term already in `ubiquitous/` keeps its definition. The board naming a word does not license rewriting what the glossary already says about it.

## Commit Contract

- Baseline: `Record {map-name} story map (baseline)`, committed with no edits.
- Edit: `Edit {map-name} story map structure`, on top of the baseline.
- One PR holding both. The PR body says what the edit changed and why, so the review can confirm that no meaning moved.

## Completeness Checks (record, not map health)

Before the baseline commit, verify the record is complete — not whether the map is "good" (that is the edit's question, and the session's before it):

- Every backbone card and story card on the board appears in the YAML; none was dropped.
- Backbone order (left to right) and story order (top to bottom) match the board.
- Every story hangs under the same task it sat under on the board.
- Every release slice on the board appears in `releases:`, in the board's top-to-bottom order, and every card above a divider carries its `release:`.
- Every pink card appears in `ubiquitous:`, and every term whose definition the board carried has its `ubiquitous/{term-key}.md`.
- Wording matches the board; nothing was paraphrased away.

## Output

`discoveries/usm/{map-name}.yaml` on one PR with two commits: the verbatim baseline, and a structural edit whose diff the review reads. Plus the `ubiquitous/{term-key}.md` files the board's definitions produced, and a note of which referenced terms are still undefined.
