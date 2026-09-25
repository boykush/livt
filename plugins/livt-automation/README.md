# livt-automation

Connect an implementation repository's coding agent to a locally running livt MCP server, so it can read the spec — story maps, stories, example mappings, rules, and the ubiquitous language — straight from the livt repository and cite it by livt URI while automating rules as tests.

## Where it sits

This is the implementation repository's end of the **delivery ring** — BDD's Automation phase, where an agreed rule becomes a failing test and then passing code. [livt-discovery](../livt-discovery/README.md) and [livt-delivery](../livt-delivery/README.md) author the livt repository; this plugin lets your repository read it. It bundles only an MCP server configuration — no skills — pointing your agent at a livt server you run locally.

```
 livt mcp --http localhost:5488      ← one local server, holds the livt repository
        │  /mcp
        ├── repo A  (livt-automation plugin)
        ├── repo B  (livt-automation plugin)
        └── repo C  (livt-automation plugin)
```

One server backs every repository on your machine: no per-repo checkout of the livt repository, no `--root` to configure. Each tool and resource payload carries a `spec_version` (the livt repository's git revision) so your agent can detect drift, and every rule, example, and question carries its own `uri` — the form to quote in a test comment, a commit message, or an issue, since a bare `R-02` exists in every mapping and names nothing on its own.

## Setup

1. **Run the server** from your livt repository checkout (one long-running process):

   ```bash
   livt mcp --http localhost:5488
   ```

   Keep `git pull` current on that checkout — the spec and its `spec_version` are read per request, so updates are served live.

2. **Install this plugin** in each implementation repository, adding livt's marketplace first if you haven't already:

   ```
   /plugin marketplace add boykush/livt
   /plugin install livt-automation@livt-claude-code-plugins
   ```

The agent can then read the livt repository through the server's tools and resources, which the server lists and describes itself, so the agent needs nothing beyond the connection. A rule whose `status` is `proposed` is a candidate awaiting agreement, not spec to automate yet.

## Configuration

The server URL defaults to `http://localhost:5488/mcp`. To use a different port or host, set `LIVT_MCP_URL` before launching your agent:

```bash
export LIVT_MCP_URL=http://localhost:8398/mcp
```

This setup assumes local use with no authentication; the server binds to localhost and is not meant to be exposed to a network.

## Using stdio instead (no server)

For a single repository without running a server, connect over stdio: the client spawns `livt mcp` as a subprocess. This does **not** use this plugin — add the server directly to the repository's `.mcp.json`, pointing `LIVT_ROOT` at your livt repository checkout:

```json
{
  "mcpServers": {
    "livt": {
      "type": "stdio",
      "command": "livt",
      "args": ["mcp"],
      "env": { "LIVT_ROOT": "/path/to/livt-repository" }
    }
  }
}
```

stdio needs the livt repository's path per repository (`LIVT_ROOT`, or `--root`), so — unlike the shared HTTP URL — it can't ship as a turnkey plugin config. Hence it's documented here rather than bundled as a plugin.

## Sending automations back

Reading the spec is the **build** station of the [delivery ring](../livt-delivery/README.md); this is the next one, **collect**. Telling the livt repository which rules your tests now cover is what keeps its board reporting the present instead of a guess, and no skill sits at the station — a machine can do it because nothing is being judged. The citation is a claim its author made while writing the test. A test makes it with a marker and a livt URI on a comment line above itself:

```go
// livt:automates livt://mapping/{story_key}/rule/{rule_id}/example/{example_id}
func TestAnExpiredCardIsRejected(t *testing.T) {
```

The comment syntax is your language's — livt looks for the marker and the URI and never reads the structure around them. `livt automations <path>` walks a checkout and collects every such line into a report, and the livt repository's build derives each rule's status from the reports committed under `automations/`. A livt URI written without the marker is a reference, not a claim, and is not collected.

The report goes back as a **pull request**, not a push: a pull request is where a citation can be checked. It lands at `automations/{owner}/{repo}.json` — one file and one branch per implementation repository, force-updated rather than accumulating, because a report is a snapshot and the open pull request should carry the latest one rather than a queue of them. None of that is yours to arrange. **Your side never scans and never writes a report**: it says which revision it merged, and the livt repository reads that revision itself. So the report is produced by one livt — the livt repository's — rather than by whichever version each implementation repository happened to install, and the path, the branch, the commit and the pull request stay where they are decided.

### The workflow

livt ships the step as an action, so nothing below spells out what a marker is or where a report goes. Set `LIVT_REPOSITORY` to your livt repository, and put a token that can send it a repository dispatch in `LIVT_REPOSITORY_TOKEN` — a GitHub App installation token, or a fine-grained PAT, with Contents write on that repository. `GITHUB_TOKEN` cannot stand in for it even when your livt repository *is* this repository: a dispatch it signs starts no workflow run, so nothing would wake to read the revision.

The actions carry no versioning of their own: they live in the livt repository, so livt's release tags are theirs. Pin by commit SHA with the release in a comment, as you would any other action, and pin `livt-version` to that same release when you want the binary to move with the action rather than tracking the newest one.

```yaml
name: livt-automations

on:
  push:
    branches: [main]

permissions: {}

jobs:
  notify:
    runs-on: ubuntu-latest
    permissions:
      contents: read # checkout, and fetch the base the gate reads
    steps:
      - uses: actions/checkout@v7
        with:
          persist-credentials: false

      - uses: boykush/livt/actions/notify@<commit-sha> # <livt release>
        with:
          livt-repository: your-org/your-livt-repository
          token: ${{ secrets.LIVT_REPOSITORY_TOKEN }}
          livt-version: <livt release> # omit to install the newest
```

That is the whole of it. The action installs livt, asks it whether this push could have changed what your tests claim, and — only then — sends the revision. It writes nothing, reads nothing of the livt repository, and needs no checkout of it.

### The livt repository's side

The other half is an action too, and one workflow there serves every implementation repository:

```yaml
name: automations

on:
  repository_dispatch:
    types: [automations-rev]

permissions: {}

jobs:
  collect:
    runs-on: ubuntu-latest
    permissions:
      contents: read # checkout, and read the committed report to compare against
    steps:
      - uses: actions/checkout@v7
        with:
          persist-credentials: false

      - uses: boykush/livt/actions/collect@<commit-sha> # <livt release>
        with:
          repository: ${{ github.event.client_payload.repository }}
          rev: ${{ github.event.client_payload.rev }}
          token: ${{ secrets.REPO_WRITER_TOKEN }}
          # read-token: only when the implementation repository is private
```

It fetches that one revision of that one repository — no history, no checkout of anything else — scans it, compares the claims against the report already committed, and opens the pull request when they differ. Nothing in it names a participating repository: the dispatch says which one, and the report is refused if it does not name the same one back. A repository joins by adding the notify step to its CI and being handed a token, and the livt repository is not edited to let it in.

`automations-rev` is livt's, the way the marker is. It is the one string both sides have to agree on, and the `types:` line is the one place a livt repository writes it — an action can ship steps, never a trigger. Get it wrong and the dispatch is accepted and nothing runs, which is quiet; it is a setup-time mistake, made once, while you are watching for the first report.

### Why the walk is gated, and why removal counts

Collecting walks every file. That is sub-second on a small repository and not free on a monorepo, so `livt automations changed <base>` answers whether walking could tell you anything new, and nothing is said when it cannot — no dispatch, and so no fetch and no walk on the other side either. The judgment is livt's rather than your CI's: livt is what decides what a marker is, so a step spelling `livt:automates` itself would keep a second copy of that definition, free to drift the moment the real one moves. This is the one thing that has to run on your side, because the diff it reads is yours.

A line that merely moved is not a reason to walk, and the reasoning is worth holding onto because it is easy to get backwards. A report is a **dated snapshot, not a live view**: its `rev`, `file` and `line` all come from one revision, and the forge URL is pinned to that revision. A later commit that moves a cited line does not break the link — the report simply describes an earlier revision, which it says out loud. A rename is the same, and `changed` answers false for both.

Removal *is* a reason. A deleted test leaves the board claiming a rule is automated when nothing covers it any more, which is the board lying about the present rather than being merely out of date. One question over added and removed lines catches both, since a deleted file shows its marker lines as removals.

A range livt cannot read — a new branch, a force-push, a revision since gone — answers **true**, not false, and says why on stderr. The two mistakes are not the same size: guessing false drops a claim, while guessing true costs one walk. That is why the fetch of the base revision above is allowed to fail quietly.

When the gate does open, the livt repository fetches that one revision and rescans the whole of it rather than patching the old report: the answer decides whether to read at all, never what the report says.

The comparison that follows is on the claims, not on the bytes, and it runs on the livt side against the report already committed there — which is why it is not in your workflow. A report never matches itself byte for byte — `rev`, `generated_at` and the urls built from `rev` all move on their own — so comparing files would open a pull request on every run. Comparing the set of cited livt URIs asks the question the pull request is actually about, and it is the same question the gate asked, one grain finer: could the claims have changed, and did they.

## What this plugin teaches, and what it leaves to you

**What a claim is, and nothing about how you make one true.** No skill ships here to say it, but the line has moved all the same: it used to sit at the read surface, with a livt URI in a test comment as the whole contract between the two repositories. The collect station put the claim itself on livt's side — two repositories that spell a claim differently produce a report about something other than what the tests wrote — and nothing else crossed over with it.

- **livt's** — the read surface it always shipped (the tools, the resources, the `spec_version` on every payload, the livt URI addressing each rule, example, and question), and now the claim that travels back the other way. The marker `livt:automates` followed by one livt URI and nothing else; a URI written without the marker as a reference rather than a claim, since production code cites rules for context too; a URI carrying a placeholder as documentation, so the form can be written about without inventing a citation; a marked line that does not read as one URI as a warning rather than a silent drop. Then the report it collects into — each citation a livt URI with its file and line, the whole of it naming the repository, the revision, and when it was read, and carrying no verdict about any of that — and what a diff must hold to be worth a walk at all, which is livt's for the same reason the marker is.
- **Yours** — everything about building against it. The test framework, the comment syntax your language spells the line in, where in the file the citation sits, whether a block wrapping a rule's cases carries the rule's own citation above them, what counts as done. The scan reads lines and URIs and never the structure around them, which is what spares livt any knowledge of your framework and leaves all of that to you as house style, for your own `AGENTS.md` or a skill you write. So is taking part at all: your repository scans itself, so joining is one step added to your CI and one credential held on your side, with nothing about it declared in the livt repository.

livt still ships no test framework and no definition of done. What it does ship now is the step itself, as two actions — so the marker, the walk, the report, where it lands and what the pull request says are all livt's, and none of them is a recipe you keep a copy of. What stays yours is everything above the citation: the framework, the comment syntax, where in the file the line sits, what counts as covered.

## Coming from 5.x

The collect workflow is replaced, not adjusted. A 5.x one installed livt, scanned this repository, checked the livt repository out, committed the report and opened the pull request. All of that is gone from your side: you tell the livt repository which revision you merged, and it reads that revision itself.

Recopy [the workflow](#the-workflow). What you keep is the trigger and a token; what you drop is the second checkout, the branch, the commit message, the pull request body, and the report file itself. The token narrows with it — Contents write on the livt repository is enough now, where Pull requests write was needed to open one.

An old workflow keeps working. Nothing on the livt side refuses a pull request opened the way 5.x opened it, so nothing breaks the day you upgrade livt. What you carry until you recopy is a version of livt on your side producing a report read by a different version on theirs, and a copy of the livt repository's file layout that no longer has to be there.

The livt repository gains a workflow of its own, listening for `automations-rev`. [The livt repository's side](#the-livt-repositorys-side) has it. Until it exists, a notify step dispatches into silence.

## Coming from 4.x

Nothing an implementation repository holds breaks. Only a rule's or an example's citation is ever attached to a board, and an agent reads each payload as it comes back, so the reshaping below costs a session nothing.

One case is worth checking. `list_stories` and `list_example_mappings` match `opportunity` against the key `list_opportunities` hands out, where they matched a story map's display name. A name written into a filter that outlives a session — your own `AGENTS.md`, a skill, a saved prompt — now returns nothing rather than failing, which is the one way this goes quiet instead of loud.

The rest is shape, and reaches only a URI or a payload written down somewhere. `livt://story-map/` takes an opportunity key where it took a percent-encoded display name, an opportunity carries `story_map` as one object where it carried a `story_maps` list, and each opportunities ref gains `key`. All of it needs livt 0.15.0 or later, as the two `livt automations` commands above do.

## Coming from 3.x

Nothing here changed. The three plugins carry one version between them, so this one moves when either of the others breaks — [livt-discovery](../livt-discovery/README.md) 4.0 split the formulation station out of the record, and [livt-delivery](../livt-delivery/README.md) 4.0 dropped the plan station.

## Coming from 1.x

The configuration is unchanged. [livt-delivery](../livt-delivery/README.md) 2.0 is where the change landed, and it is worth reading if this repository's agent also files or inspects issues.

Coming from `livt-mcp` 0.x: the configuration is unchanged there too; only the name moved from the mechanism to the station it serves.
