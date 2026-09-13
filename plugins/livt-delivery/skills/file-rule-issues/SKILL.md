---
name: file-rule-issues
description: File automation issues for business rules from an example mapping to where the story's work is tracked — one issue per rule × destination, deduped against the mapping's own issues record, with the created URL written back to the rule. The tracker and its format are the team's; the backpointers and the write-back are the contract. Use when agreed rules are ready to hand to implementation repositories for test-driven automation; give a story-key to file all unfiled rules, or add a rule-id for one. Story-level issues route to file-story-issue.
---

You **file rule issues** — the file station of the delivery ring, at the rule level.

After a story's example mapping is agreed, its business rules wait to be automated — as native tests in the implementation repositories. Your job is to send each rule there as an **automation issue**: a pointer carrying the rule, its examples, and backpointers to the livt repository. livt never reads the implementation repositories' code; the issue is the entire handoff.

## Language

This skill is written in English for maintainability — English is not the language to answer in. Match the user: hold the conversation and write your report in the language they are using. Issue bodies already follow the livt repository's language; `key:` identifiers, code, and tool commands stay English.

## The livt Repository Is the Record

The mapping YAML — not the tracker — holds the truth about what is filed where:

- Each rule in `discoveries/example-mappings/{story-key}.yaml` may carry `issues:` — a list of **issue URLs** and nothing else (no PR or test links).
- A rule with no issue URL in a destination is **unfiled** there. That record is the dedupe test — never a tracker search, never a tracker's parent/child graph (you may write links into it, you never read them back).
- Writing the URL back after filing **is** the record; an issue URL pasted onto a rule by hand counts exactly the same as one you filed.
- `automated:` is not yours to touch. It records that the rule's automation actually exists — filing (or even closing) an issue does not make that true.

## What Is Yours and What Is the Team's

livt records links, not ticket state, so it does not care which tracker a URL points at. Three things are yours; everything else follows the team:

- **Backpointers** — the livt URIs of the rule and its examples, and the `spec_version`, go in the body, so the issue can be followed back to the exact points in the spec it was cut from.
- **Write-back** — the created URL lands in the rule's `issues:`.
- **Dedupe** — against that record, and nothing else.

The destination, the tool that files, the template, labels, and fields are the team's. The story's frontmatter `repos:` declares its implementation repositories (`owner/repo`), and their tracker is the default destination; when the team tracks work elsewhere — a separate repository, a project board, another tracker — the user names it and you file there instead. Where the tracker has its own issue template, fill it and keep the backpointers; the body below is the shape when nothing else is prescribed.

## Inputs

- **story-key** (required), and optionally a **rule-id**. With a rule-id you file that one rule; with only a story-key you file every rule not yet filed to a destination.
- **Destination** — the tracker of the story's `repos:` by default, or one the user names. Ask if `repos:` is missing and no destination was given. The livt repository itself is a valid destination; nothing in the flow assumes the destination is a different repository.
- **Living document URL** (optional) — the base URL where this project's `livt build` output is published (often in the repo's README or Pages workflow). It buys the reader a browsable link, nothing more; the livt URIs are the backpointer either way, so never block filing on finding it.

## Filing Flow

1. Read `discoveries/example-mappings/{story-key}.yaml` and `stories/{story-key}.md`. Resolve the destination(s) from the story's `repos:`, or from what the user named.
2. Select the rules to file (rule-id → that one; story-key only → all), skipping any rule marked `retired: true` — the spec no longer asks for it, so there is nothing to automate — and any rule with `status: proposed`, which the spec does not ask for yet. Then dedupe each **rule × destination** pair: skip it when the rule's `issues:` already holds a URL in that destination. A link to one destination never blocks filing to another.
3. Record the spec rev of the livt repository: `git rev-parse --short HEAD`.
4. Compose each issue (see Issue Content) and file it with the tool at hand — no checkout of the target, ever. For a GitHub repository with `gh` authenticated:

   ```
   gh issue create --repo {owner}/{repo} --title "Automate {rule-id}: {rule name}" --body-file {body-file}
   ```

   For another tracker, use the tool the user side provides (an MCP server, a CLI). With no tool, hand the composed body to the user to file, and take the created URL back — the record treats that URL exactly as one you filed.
5. If the story's frontmatter `issues:` records a story issue **in the same destination**, and the tracker supports a parent/child link, attach the new issue as its child (see Parent Linking). No story issue, or no such link → the rule issue stands alone; that is a supported state, not an error.
6. Write the created URL back to the rule's `issues:` in the mapping YAML — append to the list, creating it if absent. Touch nothing else in the file. **Working-tree edit only**: no commit, no PR — the write-back rides the normal review flow.
7. Report filed and skipped pairs per rule × destination, and remind the user the write-back is uncommitted.

## Issue Content

The body carries the rule, its examples, and backpointers to the livt repository — quote rule and example names verbatim, in the mapping's language:

```markdown
## Rule

**{rule-id}** — {rule name}

### Examples

- `livt://mapping/{story-key}/rule/{rule-id}/example/{example-id}` — {example name}
- …

## livt repository

- rule: `livt://mapping/{story-key}/rule/{rule-id}`
- story: `livt://story/{story-key}`
- spec_version: `{short rev}`
- living document: {living-doc-url}/mapping/{story-key}.html#rule-{rule-id}
```

List only the rule's live examples; a retired one no longer illustrates it. Every reference is a **livt URI** — the citation form the implementation repository carries onward into test comments. A bare `{rule-id}` names nothing on its own: rule and example ids restart in every mapping. `spec_version` pins which revision of the livt repository the issue was cut from. The living document anchor is a convenience link for humans in a browser; it depends on where the site is deployed, and an issue in someone else's tracker is not yours to edit later, so it never replaces the URI. No site published? Drop that line — the URIs stand alone.

## Parent Linking

Parenthood comes from the livt repository's structure — story ⊃ rule — so the parent is the story issue recorded in the story's frontmatter `issues:` for the same destination. The link is write-only sugar for the tracker's UI; on GitHub it is the sub-issues GraphQL API:

```
# node ID of an issue (run for the parent and the new issue)
gh api graphql \
  -f query='query($owner: String!, $name: String!, $number: Int!) {
    repository(owner: $owner, name: $name) { issue(number: $number) { id } }
  }' -f owner={owner} -f name={repo} -F number={issue-number}

# attach the new issue as a sub-issue of the story issue
gh api graphql \
  -f query='mutation($parentId: ID!, $childId: ID!) {
    addSubIssue(input: { issueId: $parentId, subIssueId: $childId }) {
      issue { number }
    }
  }' -f parentId={parent-node-id} -f childId={child-node-id}
```

Never read a tracker's parent/child graph back to decide anything — dedupe and parenthood are always answered by the livt repository.

## IDs Are Forever

Your issue quotes the rule-id and every example-id in it, so you depend on this contract — but filing is not what creates it. The IDs were already immutable; an unfiled rule is not a free ID. You never mint or retire one yourself: the write-back touches `issues:` and nothing else.

The half you rely on, verbatim from the canonical statement in `change-rule`:

- **Immutability** — an ID, once used, keeps pointing at the same thing. Never renumber, never reuse, and never move an item to where its ID would change. This holds whether or not an automation issue was filed: the item's livt URI is quoted by MCP consumers, by the board's copy-link, in test comments, and in commit messages, and the livt repository records none of those — there is no list of references to check before breaking one.

## What NOT to Do

- Don't check out or read the implementation repositories — the issue is a pointer, not a synchronized copy.
- Don't commit or open a PR for the write-back; leave the working tree for the user's normal review flow.
- Don't file story-level issues — that is `file-story-issue`'s job. Missing story issue? Suggest running it; don't improvise one.
- Don't put PR or test links in `issues:`, and don't set or unset `automated:`.
- Don't cite the rule or its examples by bare id, and don't let the living-document URL stand in for the livt URI.
- Don't consult the tracker (search or its link graph) to decide what is already filed — the mapping's record is the only dedupe source.
- Don't re-file a rule × destination pair that is already linked, and don't let an existing link stop you filing the same rule to a *different* declared destination.
- Don't assume a tracker. A destination your tools cannot reach is filed by hand and recorded the same way.

## Output

One automation issue per unfiled rule × destination, each linked under the story issue where one exists there and the tracker supports it, and the mapping YAML in the working tree with every created URL appended to its rule's `issues:`. A closing report of what was filed, what was skipped and why, and the uncommitted write-back.
