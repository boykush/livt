# livt-delivery

Skills for the **delivery ring** of a livt repository: carrying an agreed example mapping to the implementation and bringing the result back.

## The ring

```
livt repository ─▶ file ─▶ build ─▶ collect ─▶ review ─▶ livt repository
                             │         │
                     implementation repository
                        (livt-automation)
```

The boundary with the [discovery ring](../livt-discovery/README.md) is the **record**: no discovery skill reads code, and nothing here decides what the mapping says. Delivery starts when an agreed record is handed outward, and the implementation is first read on the far side of that hand-off, in the implementation repository. Each **station** keeps its name whether or not a skill sits at it:

- **File** — hand the agreed rules, and the story's context, to where the team tracks its work, as issues carrying livt URI backpointers. The created URLs are written back; the livt repository, not the tracker, is the record of what is filed where.
- **Build** — the implementation repository's own station: test-driven automation of the rules, citing the spec by livt URI. That end of the ring is [livt-automation](../livt-automation/README.md).
- **Collect** — what the tests say they automate, read out of the implementation repository and carried back as a report. No skill sits here: `livt automations` runs in that repository's CI and opens a pull request when the citations change. A machine can do it because nothing is being judged — the citation is a claim its author made while writing the test.
- **Review** — the human station. The coverage is read here, and what the team learns goes back round as rule changes. Whether a test really covers the rule it names is judged in the review that lands the test, over in the implementation repository, not from a distance here.

## What these skills teach, and what they leave to you

A skill that ships one team's tracker, template, and labels as everyone's is a skill every other team has to fork — and a fork drifts away from livt's own rules along with the parts it meant to change. So these skills teach **the record contract and nothing else**:

- **livt's** — what a filed issue must carry (the livt URIs of the rule, its examples, and the story, plus `spec_version`), where the record lives (`issues:` in the livt repository, never in a tracker), how dedupe is decided (against that record, never a tracker search or its link graph), and that filing an issue says nothing about whether the rule is automated. `file-rule-issues` holds the canonical statement; the other one repeats verbatim the bullets it needs, the way `change-rule` holds the ID contract.
- **Yours** — the destination (the implementation repositories a story declares in `repos:`, or a separate tracker you name), the tool that files and reads (a tracker CLI, its MCP server, or your own hands), the template, the labels, the fields, and whether the tracker links a parent to a child at all. A URL pasted into `issues:` by hand counts exactly as one a skill filed, so wrapping a station in a skill of your own costs you nothing: yours owns the shape of the issue, livt's owns what the record says about it. Yours too is how the commits these skills write are grouped into pull requests — livt's unit of change is the commit, and the branch and review convention on top of it is your team's.

The livt repository records **links, not ticket state** — that is what lets the line sit here. Whether the work behind an issue landed is not read off the tracker at all: it is answered by the tests that cite the rule.

### Reference implementation: GitHub

The skills name no tracker, so this is the worked example rather than a bundled default — the shape a GitHub team's own wrapper can start from, with `gh` already authenticated.

Filing, and reading an issue's state and what closed it:

```bash
gh issue create --repo {owner}/{repo} --title "{title}" --body-file {body-file}
gh issue view {issue-url} --json number,state,stateReason,closedAt,closedByPullRequestsReferences
gh pr view {pr-url} --json title,url,mergedAt,files
```

GitHub takes a sub-issue from another repository, so a story issue in one repo can parent rule issues filed in several. The link is the `addSubIssue` GraphQL mutation, over node IDs:

```bash
# node ID of an issue (run for the parent and each child, each in its own {owner}/{repo})
gh api graphql \
  -f query='query($owner: String!, $name: String!, $number: Int!) {
    repository(owner: $owner, name: $name) { issue(number: $number) { id } }
  }' -f owner={owner} -f name={repo} -F number={issue-number}

gh api graphql \
  -f query='mutation($parentId: ID!, $childId: ID!) {
    addSubIssue(input: { issueId: $parentId, subIssueId: $childId }) {
      issue { number }
    }
  }' -f parentId={parent-node-id} -f childId={child-node-id}
```

Parent links are write-only sugar for the tracker's UI. Dedupe and parenthood are always answered by the livt repository, so nothing here is ever read back to decide anything.

## Skills

- **`/file-story-issue`** — file a story-level issue carrying the story body and backpointers, deduped against the story frontmatter's `issues:`, with the created URL written back. Adopts the mapping's rule issues as children where the tracker supports it.
- **`/file-rule-issues`** — file automation issues for a story's rules, one per rule × destination, deduped by the mapping's own record, each carrying the rule, its examples, and backpointers; the created URL is written back to the rule's `issues:`. Only rules with `status: accepted` are filed. Holds the canonical statement of the record contract.

Every skill here is a plain [Agent Skill](https://agentskills.io) — no subagents, no hooks, no slash-command-only behaviour — so it works in any conformant runtime, and the runtime-specific glue stays on your side of the line.

## Install

```
/plugin install livt-delivery@boykush/livt
```

## Coming from 4.x

`/inspect-automation` is gone, with the `automated:` flag it maintained. A rule is automated when a test cites it, livt collects the citations from the implementation repository and derives the status, and there is no longer a record here for a skill to propose changes to.

The judgment did not disappear; it moved to where it can be made. Whether a test really covers the rule it names is a question for the review that lands the test, in the repository that holds it — not for someone reading an issue's state from a distance, which is all the old station could see.

One thing it did is not replaced. It caught a rule whose text had moved on since its automation was recorded; a citation names a rule and not its wording, so that drift is now invisible. It is carried as an open question rather than quietly dropped.

## Coming from 3.x

`/plan-story` is gone. Holding an agreed mapping against the implementation and the design is no longer a skill here — that work now largely happens inside the session that produces the mapping, and the station's position had gone ambiguous anyway: it called itself the first stop of the delivery ring while the delivery story map never carried it and the discovery map did.

What it guarded survives it. No discovery skill reads code; they simply no longer point at one that would. And the ring boundary is the record itself — delivery starts when an agreed record is handed outward, and the implementation is first read on the far side of that hand-off. A question the code settles still reaches the mapping, as a proposal through [`/change-rule`](../livt-discovery/README.md).

## Coming from 1.x

The skills and their inputs are unchanged; what left them is the tracker. `gh issue create`, `gh issue view`, `gh pr view`, and the `addSubIssue` GraphQL calls moved out of the three delivery skills and into the reference implementation above, and the issue-body markdown shrank to the backpointer block — the rest of the body is your template now. A team that was already filing to GitHub loses nothing: the recipe is a scroll away rather than a fork away. A team that had forked a skill to change the tracker should re-read the boundary above; most of what those forks changed is no longer in the skill.

Coming from `discovery-facilitator` 0.x: `/example-mapping-plan` → removed, `/story-issue-file` → `/file-story-issue`, `/rule-issue-file` → `/file-rule-issues`, `/rule-automation-sync` → removed. The other skills moved to [livt-discovery](../livt-discovery/README.md).
