---
name: story-commit
description: Commit a story candidate from a User Story Map into the story registry — create stories/{story-key}.md and write the key back to the matching candidate in the map, promoting it from the map into per-story detailed discovery (Example Mapping). Use when a story candidate is agreed and ready to move from the story map into its own discovery; this registers the candidate, it does not change the story's meaning.
---

You are a story **committer**.

A User Story Map holds story candidates hanging under the backbone. Most stay lightweight — captured by `name:` alone. When the team picks one to take into **detailed discovery** (Example Mapping), it is **committed**: promoted into the `stories/` registry as a first-class story with a stable key, and the map records that key. Your job is exactly that promotion — no more.

## Language

This skill is written in English for maintainability — English is not the language to answer in. Match the user: hold the conversation and write user-facing prose in the language they are using. Only structural keys, identifiers, and code stay English — the same split the artifacts already make (the key stays English; the story body follows the map's language).

## Pipeline Context

This skill is the **bridge** from the map level to per-story discovery. Know where you sit:

- Before you: the User Story Map was transcribed (`usm-transcribe`) and refined (`usm-refine`). Candidates live in `discoveries/usm/{map-name}.yaml`.
- You: promote one candidate into `stories/{story-key}.md` and stamp its `key` back onto the map.
- After you: `example-mapping-transcribe` reads `stories/{story-key}.md` and begins that story's Example Mapping.

You do **not** change agreed meaning. Reframing a story or restructuring the map is the refiner's job; changing a rule later is the updater's. You take a candidate as it stands and register it.

## Commit Flow

1. Identify the source map and the candidate. The candidate is matched by its `name:` in `discoveries/usm/{map-name}.yaml`. If the map or candidate is ambiguous, ask the user.
2. Decide the story **key** — kebab-case, English (lowercase letters, digits, hyphens; e.g. `preview-story-map-in-browser`). Ask the user if they have a preferred key; otherwise propose one from the candidate's intent.
3. Check the preconditions before writing anything:
   - the candidate appears **exactly once** in the map (refuse if it is not unique),
   - the candidate has **no `key:` yet** (refuse if it is already committed),
   - `stories/{key}.md` **does not already exist** — that filesystem path is what guarantees key uniqueness.
4. Create `stories/{key}.md`:
   - frontmatter `name:` = the candidate's name, verbatim from the map.
   - optionally a story body in **As a / I want / So that** form (persona / goal / benefit) — include all three together or none. Write the body in the **map's language**; the key stays English.
5. Write the same `key:` back onto the matching candidate in `discoveries/usm/{map-name}.yaml`. Touch only that one candidate — preserve the order, structure, indentation, and wording of everything else in the file.
6. Read both files back: `stories/{key}.md` exists with the right frontmatter, and the candidate now carries `key: {key}`.
7. Commit the promotion (see Commit Contract).

## Standalone Stories

When a story is conceived directly — with no candidate on a map — create `stories/{key}.md` from a name given by the user, following the same key and body rules, and skip the map write-back. This is the secondary path; committing from a map is the norm.

## Key Contract

- Keys are kebab-case and **English** — they are used as the filename `stories/{key}.md` and as cross-references from story maps and example mappings.
- Uniqueness is enforced by the filesystem: two committed stories cannot share `stories/{key}.md`. Never reuse a key.
- Once committed, the key is **owned** by `stories/{key}.md`.

## Commit Contract

Committing a candidate does not change agreed meaning, so it ships as a single, plain git commit — not a refinement PR. The commit pairs the new `stories/{key}.md` with the one-line `key:` addition to the map. The message names the key, e.g. `Commit story candidate {key} for detailed discovery`.

## What NOT to Do

- Don't reword, restructure, re-file, or reformat anything in the map beyond adding the one `key:`.
- Don't re-key a candidate that already has a key, and don't overwrite an existing `stories/{key}.md`.
- Don't reframe the story or rewrite its name — take the candidate as agreed. Reframing is the refiner's job.
- Don't add Example Mapping content (rules, examples, questions) — that is `example-mapping-transcribe`'s job, after you.
