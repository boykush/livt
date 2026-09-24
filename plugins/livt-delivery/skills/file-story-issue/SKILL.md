---
name: file-story-issue
description: File a story-level issue from the livt repository to where the story's work is tracked — carrying the story body and backpointers, deduped against the story frontmatter's own record, with the created URL written back there. The backpointers, the write-back, and the dedupe are livt's; the tracker, the template, and the tool that files are the team's. Use when a story's context should be handed to an implementation repository, with or without rule issues; per-rule automation issues route to file-rule-issues.
---

You **file a story issue** — the file station of the delivery ring, at the story level.

A story gives its rules their context — the persona, the goal, the benefit. Your job is to hand that context to the team's tracker as one **story-level issue**, which then acts as the parent bundling the story's rule-level automation issues, wherever they were filed. The story issue also stands entirely on its own: a story with no rule issues is a supported state, not a missing half.

## Language

This skill is written in English for maintainability — English is not the language to answer in. Match the user: hold the conversation and write your report in the language they are using. Issue bodies already follow the livt repository's language; `key:` identifiers, code, and tool commands stay English.

## The Record Contract

Verbatim from the canonical statement in `file-rule-issues` — the bullets this station depends on. A skill loads on its own, so it carries them rather than pointing at them. Change one, change all.

- **The record lives in the livt repository.** A rule's `issues:` and a story's frontmatter `issues:` hold the automation issues filed for them. Neither is stored in a tracker, and neither is derived from one.
- **`issues:` holds issue URLs and nothing else** — no PR links, no test links. A URL pasted in by hand counts exactly as one a skill filed.
- **Dedupe reads the record, never the tracker.** An item is unfiled in a destination when its `issues:` holds no URL there. Never a tracker search, never a tracker's parent/child graph: links are written into a tracker, never read back out of it.
- **The write-back is the filing.** An issue created and not recorded did not happen, so the created URL lands in `issues:` before you are done — a **working-tree edit only**, no commit and no PR, riding the team's normal review flow.

One thing at the story level is yours alone: the singular `issue:` some stories carry is **unrelated free metadata** — any link the team parked there. Never read it for dedupe, never edit it, never migrate it into `issues:`.

## Yours and the Team's

Three things are yours — the contract above, applied:

- **Backpointers** — the livt URI of the story and the `spec_version` go in the body, so the issue can be followed back to the point in the spec it was cut from.
- **Write-back** — the created URL lands in the story's `issues:`.
- **Dedupe** — against that record, and nothing else.

Everything else is the team's: the destination, the tool that files, the template, the labels, the fields, and whether the tracker can link a parent to a child at all. livt records links, not ticket state, so it does not care which tracker a URL points at. The story's frontmatter `repos:` declares its implementation repositories (`owner/repo`), and their tracker is the default destination; when the team tracks work elsewhere — a separate repository, a project board, another tracker — the user names it and you file there instead. Where the tracker has its own issue template, fill it and keep the backpointers.

This skill ships no tracker knowledge on purpose: one team's answer shipped as everyone's is what makes a skill need forking. Use the tool the user side provides, and where a team has wrapped this station in a skill of its own, that skill owns the shape of the issue — you own what the record says about it.

## Inputs

- **story-key** (required). A mapping whose story has no card has no story to file; its rules go out on their own through `file-rule-issues`.
- **Destination** — the tracker of the story's `repos:` by default; with several declared, or a tracker the user names, confirm which one(s). Ask if `repos:` is missing and no destination was given. The livt repository itself is a valid destination; nothing in the flow assumes the destination is a different repository.

## Filing Flow

1. Read `stories/{story-key}.md` — frontmatter and body — and `discoveries/example-mappings/{story-key}.yaml` if it exists (you'll need its rules' `issues:` for adoption).
2. Dedupe per **story × destination** against the record. A link to one destination never blocks filing to another.
3. Record the spec rev of the livt repository: `git rev-parse --short HEAD`.
4. Compose the issue (see Issue Content) and file it with the tool the team uses — a tracker's CLI or MCP server. Never check out the target repository. With no tool that reaches the destination, hand the composed body to the user and take the created URL back; the record treats it exactly as one you filed.
5. **Adopt existing rule issues** where the tracker supports a parent/child link (see Parent Linking): every issue URL in the mapping's rules' `issues:` becomes a child of the new story issue, save those whose own destination already records a story issue — that nearer story issue is their parent. Parenthood does not depend on filing order — rule issues filed earlier are adopted now; rule issues filed later attach themselves (`file-rule-issues`'s job). A tracker without such links, or one that refuses this one, skips the step; the record needs none of it.
6. Write the created URL back to the story's frontmatter `issues:` — append to the list, creating it if absent. Touch nothing else in the file.
7. Report what was filed, which rule issues were adopted, what was skipped, and remind the user the write-back is uncommitted.

## Issue Content

The body carries the story verbatim — the As a / I want / So that text as written, in the story's language — and then the backpointers, which are exact:

```markdown
## livt repository

- story: `livt://story/{story-key}`
- spec_version: `{short rev}`
```

The **livt URI** is the citation form the implementation repository carries onward; `spec_version` pins which revision of the livt repository the issue was cut from. If you know the living document URL (where `livt build` output is published), add its story page — `{living-doc-url}/story/{story-key}.html` — below as a browsable convenience for humans. It depends on where the site is deployed, so it never replaces the URI; don't block filing on finding it.

Around that block the shape is the team's — their issue template when they have one.

## Parent Linking

Parenthood comes from the livt repository's structure — story ⊃ rule — so the children are the rule issues recorded in the mapping, whichever destination each was filed to.

Whether a tracker takes a child from another repository is its answer, not yours to assume ahead of it — attempt the link and report what came back. GitHub takes one, through the `addSubIssue` GraphQL mutation. The link is write-only sugar for the tracker's UI: never read a parent/child graph back to decide anything.

## What NOT to Do

- Don't file rule-level automation issues — that is `file-rule-issues`'s job, and a story issue with zero children is fine.
- Don't check out or read the implementation repositories — the issue is a pointer, not a synchronized copy.
- Don't consult the tracker to decide what is already filed, and don't re-file a story × destination pair the record already holds — but don't let an existing link stop you filing the same story to a *different* declared destination.
- Don't touch the singular `issue:` frontmatter field, and don't use it for dedupe.
- Don't commit or open a PR for the write-back.
- Don't assume a tracker, and don't rule a parent/child link out because the two issues sit in different repositories — both are the tracker's answer, not yours.

## Output

One story-level issue in the chosen destination, carrying the story body and backpointers, with the mapping's rule issues adopted as children where the tracker supports it, and `stories/{story-key}.md` in the working tree with the created URL appended to frontmatter `issues:`. A closing report of what was filed, what was adopted, what was skipped and why, and the uncommitted write-back.
