---
name: file-story-issue
description: File a story-level issue from the livt repository to where the story's work is tracked — carrying the story body and backpointers, deduped against the story frontmatter's own issues record, with the created URL written back there. The tracker and its format are the team's; the backpointers and the write-back are the contract. Use when a story's context should be handed to an implementation repository, with or without rule issues; per-rule automation issues route to file-rule-issues.
---

You **file a story issue** — the file station of the delivery ring, at the story level.

A story gives its rules their context — the persona, the goal, the benefit. Your job is to hand that context to the team's tracker as one **story-level issue**, which then acts as the parent bundling the story's rule-level automation issues filed in the same place. The story issue also stands entirely on its own: a story with no rule issues is a supported state, not a missing half.

## Language

This skill is written in English for maintainability — English is not the language to answer in. Match the user: hold the conversation and write your report in the language they are using. Issue bodies already follow the livt repository's language; `key:` identifiers, code, and tool commands stay English.

## The livt Repository Is the Record

The story file — not the tracker — holds the truth about what is filed where:

- `stories/{story-key}.md`'s frontmatter may carry `issues:` — a list of **issue URLs** filed for this story.
- No URL in a tracker means the story is **unfiled** there. That record is the dedupe test — never a tracker search, never a tracker's parent/child graph (you may write links into it, you never read them back).
- Writing the URL back after filing **is** the record; a URL pasted in by hand counts exactly the same.
- The singular `issue:` some stories carry is **unrelated free metadata** (any link the team parked there). Never read it for dedupe, never edit it, never migrate it into `issues:`.

## What Is Yours and What Is the Team's

livt records links, not ticket state, so it does not care which tracker a URL points at. Three things are yours; everything else follows the team:

- **Backpointers** — the livt URI of the story and the `spec_version` go in the body, so the issue can be followed back to the point in the spec it was cut from.
- **Write-back** — the created URL lands in the story's `issues:`.
- **Dedupe** — against that record, and nothing else.

The destination, the tool that files, the template, labels, and fields are the team's. The story's frontmatter `repos:` declares its implementation repositories (`owner/repo`), and their tracker is the default destination; when the team tracks work elsewhere — a separate repository, a project board, another tracker — the user names it and you file there instead. Where the tracker has its own issue template, fill it and keep the backpointers; the body below is the shape when nothing else is prescribed.

## Inputs

- **story-key** (required).
- **Destination** — the tracker of the story's `repos:` by default; with several declared, or a tracker the user names, confirm which one(s). Ask if `repos:` is missing and no destination was given. The livt repository itself is a valid destination; nothing in the flow assumes the destination is a different repository.

## Filing Flow

1. Read `stories/{story-key}.md` — frontmatter and body — and `discoveries/example-mappings/{story-key}.yaml` if it exists (you'll need its rules' `issues:` for adoption).
2. Dedupe per **story × destination**: skip a destination when the frontmatter `issues:` already holds a URL there. A link to one destination never blocks filing to another.
3. Record the spec rev of the livt repository: `git rev-parse --short HEAD`.
4. Compose the issue (see Issue Content) and file it with the tool at hand — no checkout of the target, ever. For a GitHub repository with `gh` authenticated:

   ```
   gh issue create --repo {owner}/{repo} --title "{story name}" --body-file {body-file}
   ```

   For another tracker, use the tool the user side provides (an MCP server, a CLI). With no tool, hand the composed body to the user to file, and take the created URL back — the record treats that URL exactly as one you filed.
5. **Adopt existing rule issues** where the tracker supports a parent/child link (see Parent Linking): every issue URL in the mapping's rules' `issues:` that lives in the same destination becomes a child of the new story issue. Parenthood does not depend on filing order — rule issues filed earlier are adopted now; rule issues filed later attach themselves (`file-rule-issues`'s job). A tracker without such links skips this step; the record needs none of it.
6. Write the created URL back to the story's frontmatter `issues:` — append to the list, creating it if absent. Touch nothing else in the file. **Working-tree edit only**: no commit, no PR — the write-back rides the normal review flow.
7. Report what was filed, which rule issues were adopted, what was skipped, and remind the user the write-back is uncommitted.

## Issue Content

The body carries the story verbatim — in the story's language — plus backpointers to the livt repository:

```markdown
{story body — the As a / I want / So that text, as written}

## livt repository

- story: `livt://story/{story-key}`
- spec_version: `{short rev}`
```

The **livt URI** is the citation form the implementation repository carries onward; `spec_version` pins which revision of the livt repository the issue was cut from. If you know the living document URL (where `livt build` output is published), add its story page — `{living-doc-url}/story/{story-key}.html` — below as a browsable convenience for humans. It depends on where the site is deployed, so it never replaces the URI; don't block filing on finding it.

## Parent Linking

Parenthood comes from the livt repository's structure — story ⊃ rule — so the children are the rule issues recorded in the mapping for the same destination. The link is write-only sugar for the tracker's UI; on GitHub it is the sub-issues GraphQL API:

```
# node ID of an issue (run for the new story issue and each rule issue)
gh api graphql \
  -f query='query($owner: String!, $name: String!, $number: Int!) {
    repository(owner: $owner, name: $name) { issue(number: $number) { id } }
  }' -f owner={owner} -f name={repo} -F number={issue-number}

# attach a rule issue as a sub-issue of the story issue
gh api graphql \
  -f query='mutation($parentId: ID!, $childId: ID!) {
    addSubIssue(input: { issueId: $parentId, subIssueId: $childId }) {
      issue { number }
    }
  }' -f parentId={story-issue-node-id} -f childId={rule-issue-node-id}
```

Never read a tracker's parent/child graph back to decide anything — dedupe and parenthood are always answered by the livt repository.

## What NOT to Do

- Don't file rule-level automation issues — that is `file-rule-issues`'s job, and a story issue with zero children is fine.
- Don't check out or read the implementation repositories — the issue is a pointer, not a synchronized copy.
- Don't commit or open a PR for the write-back; leave the working tree for the user's normal review flow.
- Don't touch the singular `issue:` frontmatter field, and don't use it for dedupe — only `issues:` records filings.
- Don't consult the tracker (search or its link graph) to decide what is already filed — the story's record is the only dedupe source.
- Don't re-file a story × destination pair that is already linked, and don't let an existing link stop you filing the same story to a *different* declared destination.
- Don't assume a tracker. A destination your tools cannot reach is filed by hand and recorded the same way.

## Output

One story-level issue in the chosen destination, carrying the story body and backpointers, with same-destination rule issues adopted as children where the tracker supports it, and `stories/{story-key}.md` in the working tree with the created URL appended to frontmatter `issues:`. A closing report of what was filed, what was adopted, what was skipped and why, and the uncommitted write-back.
