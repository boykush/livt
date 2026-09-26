---
name: file-rule-issues
description: File automation issues for the business rules of an agreed example mapping to where the story's work is tracked — one issue per rule × destination, deduped against the mapping's own record, with the created URL written back to the rule. Holds the canonical statement of livt's record contract. The backpointers, the write-back, and the dedupe are livt's; the tracker, the template, and the tool that files are the team's. Use when agreed rules are ready to hand to implementation repositories for test-driven automation; give a story-key to file all unfiled rules, or add a rule-id for one. Story-level issues route to file-story-issue.
---

You **file rule issues** — the file station of the delivery ring, at the rule level.

After a story's example mapping is agreed, its business rules wait to be automated — as native tests in the implementation repositories. Your job is to send each rule there as an **automation issue**: a pointer carrying the rule, its examples, and backpointers to the livt repository. Filing reads no code on the far side; the issue is the entire handoff.

## Language

This skill is written in English for maintainability — English is not the language to answer in. Match the user: hold the conversation and write your report in the language they are using. Issue bodies already follow the livt repository's language; `key:` identifiers, code, and tool commands stay English.

## The Record Contract

This is the canonical statement. `file-story-issue` repeats the bullets it needs, verbatim — a skill loads on its own, so every skill that reads or writes the record has to carry them. Change one, change all.

- **The record lives in the livt repository.** A rule's `issues:` and a story's frontmatter `issues:` hold the automation issues filed for them. Neither is stored in a tracker, and neither is derived from one.
- **`issues:` holds issue URLs and nothing else** — no PR links, no test links. A URL pasted in by hand counts exactly as one a skill filed.
- **Dedupe reads the record, never the tracker.** An item is unfiled in a destination when its `issues:` holds no URL there. Never a tracker search, never a tracker's parent/child graph: links are written into a tracker, never read back out of it.
- **The write-back is the filing.** An issue created and not recorded did not happen, so the created URL lands in `issues:` before you are done — a **working-tree edit only**, no commit and no PR, riding the team's normal review flow.
- **Filing is not automating.** An issue opened here, or closed over there, says nothing about whether a rule is automated: that is answered by the tests citing it, which livt collects from the implementation repository.

## Yours and the Team's

Three things are yours — the contract above, applied:

- **Backpointers** — the livt URIs of the rule and its examples, and the `spec_version`, go in the body, so the issue can be followed back to the exact points in the livt repository they were cut from.
- **Write-back** — the created URL lands in the rule's `issues:`.
- **Dedupe** — against that record, and nothing else.

Everything else is the team's: the destination, the tool that files, the template, the labels, the fields, and whether the tracker can link a parent to a child at all. livt records links, not ticket state, so it does not care which tracker a URL points at. The story's frontmatter `repos:` declares its implementation repositories (`owner/repo`), and their tracker is the default destination; when the team tracks work elsewhere — a separate repository, a project board, another tracker — the user names it and you file there instead. Where the tracker has its own issue template, fill it and keep the backpointers.

This skill ships no tracker knowledge on purpose: one team's answer shipped as everyone's is what makes a skill need forking. Use the tool the user side provides, and where a team has wrapped this station in a skill of its own, that skill owns the shape of the issue — you own what the record says about it.

## Inputs

- **story-key** (required), and optionally a **rule-id**. With a rule-id you file that one rule; with only a story-key you file every rule not yet filed to a destination.
- **Destination** — the tracker of the story's `repos:` by default, or one the user names. Ask if `repos:` is missing and no destination was given. The livt repository itself is a valid destination; nothing in the flow assumes the destination is a different repository.
- **Living document URL** (optional) — the base URL where this project's `livt build` output is published (often in the repo's README or Pages workflow). It buys the reader a browsable link, nothing more; the livt URIs are the backpointer either way, so never block filing on finding it.

## Filing Flow

1. Read `discoveries/example-mappings/{story-key}.yaml` and `stories/{story-key}.md`. Resolve the destination(s) from the story's `repos:`, or from what the user named. A story with no card has no `repos:` to read, so ask for the destination.
2. Select the rules to file (rule-id → that one; story-key only → all), keeping only rules with `status: accepted`: a `proposed` rule is not asked for yet, and a `rejected` or `retired` one is not asked for any more, so neither has anything to automate. Then dedupe each **rule × destination** pair against the record. A link to one destination never blocks filing to another.
3. Record the revision of the livt repository: `git rev-parse --short HEAD`.
4. Compose each issue (see Issue Content) and file it with the tool the team uses — a tracker's CLI or MCP server. Never check out the target repository. With no tool that reaches the destination, hand the composed body to the user and take the created URL back; the record treats it exactly as one you filed.
5. Where the story records a story issue and the tracker supports parent/child links, attach the new issue under it (see Parent Linking).
6. Write the created URL back to the rule's `issues:` — append to the list, creating it if absent. Touch nothing else in the file.
7. Report filed and skipped pairs per rule × destination, and remind the user the write-back is uncommitted.

## Issue Content

The body carries the rule and its live examples — quoted verbatim, in the mapping's language, a retired example left out because it no longer illustrates the rule — and then the backpointers, which are exact:

```markdown
## livt repository

- rule: `livt://mapping/{story-key}/rule/{rule-id}`
- examples: `livt://mapping/{story-key}/rule/{rule-id}/example/{example-id}`, …
- story: `livt://story/{story-key}`
- spec_version: `{short rev}`
- living document: {living-doc-url}/mapping/{story-key}.html#rule-{rule-id}
```

Every reference is a **livt URI** — the citation form the implementation repository carries onward into test comments. A bare `{rule-id}` names nothing on its own: rule and example ids restart in every mapping, so each example is quoted by its own URI and a test can cite the one it covers. `spec_version` pins which revision of the livt repository the issue was cut from. The living document anchor is a convenience for humans in a browser; it depends on where the site is deployed, and an issue in someone else's tracker is not yours to edit later, so it never replaces the URI. No site published? Drop that line. No story card? Drop the `story:` line too: there is no story resource for it to name.

Above that block the shape is the team's — their issue template, or a plain heading naming the rule and a list of its examples when nothing is prescribed.

## Parent Linking

Parenthood comes from the livt repository's structure — story ⊃ rule — so the parent is the story issue recorded in the story's frontmatter `issues:`, whichever destination it lives in: the one in this destination when there is one, otherwise the one recorded elsewhere; with several elsewhere and none here, ask which. No story issue, no such link, or a link the tracker refuses → the rule issue stands alone. That is a supported state, not an error.

Whether a tracker takes a child from another repository is its answer, not yours to assume ahead of it — attempt the link and report what came back. GitHub takes one, through the `addSubIssue` GraphQL mutation. The link is write-only sugar for the tracker's UI: never read a parent/child graph back to decide anything.

## IDs Are Forever

Your issue quotes the rule-id and every example-id in it, so you depend on this contract — but filing is not what creates it. The IDs were already immutable; an unfiled rule is not a free ID. You never mint or retire one yourself: the write-back touches `issues:` and nothing else.

The half you rely on, verbatim from the canonical statement in `change-rule`:

- **Immutability** — an ID, once used, keeps pointing at the same thing. Never renumber, never reuse, and never move an item to where its ID would change. This holds whether or not an automation issue was filed: the item's livt URI is quoted by MCP consumers, by the board's copy-link, in test comments, and in commit messages, and the livt repository records none of those — there is no list of references to check before breaking one.

## What NOT to Do

- Don't check out or read the implementation repositories — the issue is a pointer, not a synchronized copy.
- Don't consult the tracker to decide what is already filed, and don't re-file a rule × destination pair the record already holds — but don't let an existing link stop you filing the same rule to a *different* declared destination.
- Don't put PR or test links in `issues:` — what goes there is an issue's URL, and a test is found by the citation it carries.
- Don't cite the rule or its examples by bare id, and don't let the living-document URL stand in for the livt URI.
- Don't commit or open a PR for the write-back.
- Don't file story-level issues — that is `file-story-issue`'s job. Missing story issue? Suggest running it; don't improvise one.
- Don't assume a tracker, and don't rule a parent/child link out because the two issues sit in different repositories — both are the tracker's answer, not yours.

## Output

One automation issue per unfiled rule × destination, each linked under the story issue where one is recorded and the tracker supports it, and the mapping YAML in the working tree with every created URL appended to its rule's `issues:`. A closing report of what was filed, what was skipped and why, and the uncommitted write-back.
