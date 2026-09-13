---
name: write-story-card
description: Write a story's card — create stories/{story-key}.md for a story candidate on a User Story Map and stamp the key back onto that candidate — so the story can go into its own conversation, an Example Mapping. Use when the team picks a candidate for detailed discovery, or conceives a story that sits on no map; it registers the card and changes no meaning. Rules and examples route to record-example-mapping.
---

You **write a story's card** — the card station of the discovery ring.

A User Story Map holds story candidates hanging under the backbone. Most stay lightweight — captured by `name:` alone. When the team picks one to take into its own **conversation** (an Example Mapping), it gets a **card**: `stories/{story-key}.md`, with a stable key the map records. This is the Card of Ron Jeffries' three Cs — a reminder of a conversation still to come, not a requirement. The conversation follows on the story's board, and `record-example-mapping` writes the confirmation — rules and examples — against the card's key.

## Language

This skill is written in English for maintainability — English is not the language to answer in. Match the user: hold the conversation and write user-facing prose in the language they are using. Only structural keys, identifiers, and code stay English — the same split the artifacts already make (the key stays English; the card's body follows the map's language).

## Where You Sit

- Before you: the map was recorded (`record-story-map`). Candidates live in `discoveries/usm/{map-name}.yaml`.
- You: write one candidate's card and stamp its `key:` onto the map.
- After you: the team prepares and holds the story's Example Mapping; `record-example-mapping` reads `stories/{story-key}.md` for the key and scope.

You do **not** change agreed meaning. Reframing a story or restructuring the map belongs to `record-story-map`'s edit commit; changing a rule later is `change-rule`'s. You take a candidate as it stands and give it a card.

## Card Flow

1. Identify the source map and the candidate. The candidate is matched by its `name:` in `discoveries/usm/{map-name}.yaml`. If the map or candidate is ambiguous, ask the user.
2. Decide the story **key** — kebab-case, English (lowercase letters, digits, hyphens; e.g. `preview-story-map-in-browser`). Ask the user if they have a preferred key; otherwise propose one from the candidate's intent.
3. Check the preconditions before writing anything:
   - the candidate appears **exactly once** in the map (refuse if it is not unique),
   - the candidate has **no `key:` yet** (refuse if it already has a card),
   - `stories/{key}.md` **does not already exist** — that filesystem path is what guarantees key uniqueness.
4. Create `stories/{key}.md`:
   - frontmatter `name:` = the candidate's name, verbatim from the map.
   - optionally a story body in **As a / I want / So that** form (persona / goal / benefit) — include all three together or none. Write the body in the **map's language**; the key stays English.
5. Write the same `key:` back onto the matching candidate in `discoveries/usm/{map-name}.yaml`. Touch only that one candidate — preserve the order, structure, indentation, and wording of everything else in the file.
6. Read both files back: `stories/{key}.md` exists with the right frontmatter, and the candidate now carries `key: {key}`.
7. Commit the card (see Commit Contract).

## The Card Is the Team's

You write the card's `name:` and, when asked, its As a / I want / So that body. Whatever the team gathers to prepare the conversation — links, context, a board to walk into the session with, a declaration of where the story's work will be tracked — is theirs to add to the card afterwards. Leave the file plain and easy to extend, and don't invent preparation the team did not ask for.

## Standalone Stories

When a story is conceived directly — with no candidate on a map — create `stories/{key}.md` from a name given by the user, following the same key and body rules, and skip the map write-back. This is the secondary path; a card written from a map is the norm.

## Key Contract

- Keys are kebab-case and **English** — they are used as the filename `stories/{key}.md` and as cross-references from story maps and example mappings.
- Uniqueness is enforced by the filesystem: two stories cannot share `stories/{key}.md`. Never reuse a key.
- Once the card is written, the key is **owned** by `stories/{key}.md`.

## Commit Contract

Writing a card does not change agreed meaning, so it ships as a single, plain git commit — not a PR that asks for review of a diff. The commit pairs the new `stories/{key}.md` with the one-line `key:` addition to the map. The message names the key, e.g. `Write story card {key}`.

## What NOT to Do

- Don't reword, restructure, re-file, or reformat anything in the map beyond adding the one `key:`.
- Don't re-key a candidate that already has one, and don't overwrite an existing `stories/{key}.md`.
- Don't reframe the story or rewrite its name — take the candidate as agreed. Reframing is the record's edit commit, on the map.
- Don't add Example Mapping content (rules, examples, questions) — that is `record-example-mapping`'s job, after the conversation.
