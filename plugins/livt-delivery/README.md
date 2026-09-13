# livt-delivery

Skills for the **delivery ring** of a livt repository: carrying an agreed example mapping to the implementation and bringing the result back.

## The ring

```
livt repository ─▶ plan ─▶ file ─▶ build ─▶ inspect ─▶ review ─▶ livt repository
                                     │
                          implementation repository
                          (livt-automation)
```

The boundary with the [discovery ring](../livt-discovery/README.md) is the **first read of the implementation**: no discovery skill reads code, and `plan-story` is where that starts. Each skill is named for its **station**:

- **Plan** — sprint planning (Scrum). The agreed mapping is the backlog item; hold each rule against the current implementation and design and work out how to realize it. Questions the code settles come back to discovery as proposals, through `change-rule`.
- **File** — hand the agreed rules, and the story's context, to where the team tracks its work, as issues carrying livt URI backpointers. The created URLs are written back; the livt repository, not the tracker, is the record of what is filed where.
- **Build** — the implementation repository's own station: test-driven automation of the rules, citing the spec by livt URI. That end of the ring is [livt-automation](../livt-automation/README.md).
- **Inspect** — Scrum's inspection. Hold each rule's `automated:` against what actually happened — the issues the rule records, the changes that closed them, the mapping's own history — and propose every correction for review with the evidence attached. A closed issue is a trigger, never proof.
- **Review** — the human station. The proposals are judged here, the coverage is read here, and what the team learns goes back round as rule changes.

## What these skills assume about your tracker

The livt repository records **links, not ticket state**: a rule's `issues:` and a story's `issues:` hold URLs, their state lives at the URL, and a URL pasted in by hand counts exactly as one a skill filed. So the skills assume no particular tracker:

- **Theirs** — the backpointers (livt URIs and `spec_version`) in every issue body, the write-back of the created URL, and dedupe against that record.
- **Yours** — the destination (the implementation repositories a story declares in `repos:`, or a separate tracker you name), the tool that files (`gh` for GitHub, a tracker's MCP server or CLI, or your own hands), the template, labels, and fields. Parent/child links between a story issue and its rule issues are written where the tracker supports them and skipped where it does not.
- **Inspection** reads issue state with the same tools, and reports a URL it cannot read as *unverifiable* rather than guessing.

## Skills

- **`/plan-story`** — plan a story's implementation from its agreed example mapping, rule by rule, against the implementation and design the user (or the story's frontmatter) points at: `feasible`, `needs-decision`, `contradicts`, `resolved-question`, or `unverifiable`, each citing where it looked. A written plan for a human decision — never an edit of the mapping, never an audit of what is built.
- **`/file-story-issue`** — file a story-level issue carrying the story body and backpointers, deduped against the story frontmatter's `issues:`, with the created URL written back. Adopts same-destination rule issues as children where the tracker supports it.
- **`/file-rule-issues`** — file automation issues for a story's rules, one per rule × destination, deduped by the mapping's own record, each carrying the rule, its examples, and backpointers; the created URL is written back to the rule's `issues:`. Proposed and retired rules are skipped.
- **`/inspect-automation`** — propose setting `automated:` where every recorded issue closed with tests behind it, and unsetting it where a rule changed after the flag was set — one rule's record per PR, evidence attached, unverifiable reported as such.

## Install

```
/plugin install livt-delivery@boykush/livt
```

Coming from `discovery-facilitator` 0.x: `/example-mapping-plan` → `/plan-story`, `/story-issue-file` → `/file-story-issue`, `/rule-issue-file` → `/file-rule-issues`, `/rule-automation-sync` → `/inspect-automation`. The other skills moved to [livt-discovery](../livt-discovery/README.md).
