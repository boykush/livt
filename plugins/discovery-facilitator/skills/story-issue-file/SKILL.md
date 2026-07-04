---
name: story-issue-file
description: File a story-level issue from the story master to a declared implementation repository — carrying the story body and backpointers, deduped against the story frontmatter's own issues record, with the created URL written back there. Existing rule issues in the same repository are adopted as sub-issues. Use when a story's context should be handed to an implementation repo, with or without rule issues; per-rule automation issues route to rule-issue-file.
---

You are a story issue **filer**.

A story gives its rules their context — the persona, the goal, the benefit. Your job is to hand that context to an implementation repository as one **story-level issue**, which then acts as the parent bundling the story's rule-level automation issues in that repository. The story issue also stands entirely on its own: a story with no rule issues is a supported state, not a missing half.

## Master Is the Record

The story file — not GitHub — holds the truth about what is filed where:

- `stories/{story-key}.md`'s frontmatter may carry `issues:` — a list of **issue URLs** filed for this story.
- No URL in a repository means the story is **unfiled** there. That record is the dedupe test — never GitHub search, never GitHub's sub-issue graph (you write sub-issue links, you don't read them).
- Writing the URL back after filing **is** the record; a URL pasted in by hand counts exactly the same.
- The singular `issue:` some stories carry is **unrelated free metadata** (any link the team parked there). Never read it for dedupe, never edit it, never migrate it into `issues:`.

## Inputs

- **story-key** (required).
- **Target repository** — from the story's frontmatter `repos:` (a list of `owner/repo`); with several declared, confirm which one(s) with the user. Ask if `repos:` is missing. The master's own repository is a valid target; nothing in the flow assumes the target is a different repository.

## Filing Flow

1. Read `stories/{story-key}.md` — frontmatter and body — and `discoveries/example-mappings/{story-key}.yaml` if it exists (you'll need its rules' `issues:` for adoption).
2. Dedupe per **story × repository**: skip a repository when the frontmatter `issues:` already holds a URL there. A link to one repository never blocks filing to another.
3. Record the spec rev of the master: `git rev-parse --short HEAD`.
4. Compose the issue (see Issue Content) and file it with existing `gh` auth — no checkout of the target, ever:

   ```
   gh issue create --repo {owner}/{repo} --title "{story name}" --body-file {body-file}
   ```

5. **Adopt existing rule issues**: every issue URL in the mapping's rules' `issues:` that lives in the same repository becomes a sub-issue of the new story issue (see Sub-issue Linking). Parenthood does not depend on filing order — rule issues filed earlier are adopted now; rule issues filed later attach themselves (`rule-issue-file`'s job).
6. Write the created URL back to the story's frontmatter `issues:` — append to the list, creating it if absent. Touch nothing else in the file. **Working-tree edit only**: no commit, no PR — the write-back rides the normal review flow.
7. Report what was filed, which rule issues were adopted, what was skipped, and remind the user the write-back is uncommitted.

## Issue Content

The body carries the story verbatim — in the story's language — plus backpointers to the master:

```markdown
{story body — the As a / I want / So that text, as written}

## Master

- story: `{story-key}`
- spec_version: `{short rev}`
```

If you know the living document URL (where `livt build` output is published), add its story page — `{living-doc-url}/story/{story-key}.html` — to the Master section; don't block filing on it.

## Sub-issue Linking

Parenthood comes from the master's structure — story ⊃ rule — so the children are the rule issues recorded in the mapping for the same repository. Link via GitHub's sub-issues GraphQL API:

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

The link is write-only sugar for GitHub's UI. Never read the sub-issue graph back to decide anything — dedupe and parenthood are always answered by the master.

## What NOT to Do

- Don't file rule-level automation issues — that is `rule-issue-file`'s job, and a story issue with zero sub-issues is fine.
- Don't check out or read the implementation repositories — the issue is a pointer, not a synchronized copy.
- Don't commit or open a PR for the write-back; leave the working tree for the user's normal review flow.
- Don't touch the singular `issue:` frontmatter field, and don't use it for dedupe — only `issues:` records filings.
- Don't consult GitHub (search or sub-issue graph) to decide what is already filed — the story's record is the only dedupe source.
- Don't re-file a story × repository pair that is already linked, and don't let an existing link stop you filing the same story to a *different* declared repository.

## Output

One story-level issue in the chosen repository, carrying the story body and backpointers, with same-repository rule issues adopted as sub-issues, and `stories/{story-key}.md` in the working tree with the created URL appended to frontmatter `issues:`. A closing report of what was filed, what was adopted, what was skipped and why, and the uncommitted write-back.
