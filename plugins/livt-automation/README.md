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

The agent can then call `list_stories` / `list_example_mappings` / `list_story_maps` / `list_terms` and read resources such as `livt://story-map/{map_name}`, `livt://story/{story_key}`, `livt://mapping/{story_key}`, and `livt://ubiquitous/{term_key}`. See the [`livt mcp` command reference](https://github.com/boykush/livt/blob/main/docs/src/reference/commands.md) for the full tool and resource list. A rule whose `status` is `proposed` is a candidate awaiting agreement, not spec to automate yet.

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
// livt:automates livt://mapping/checkout/rule/R-02/example/EX-01
func TestAnExpiredCardIsRejected(t *testing.T) {
```

The comment syntax is your language's — livt looks for the marker and the URI and never reads the structure around them. `livt automations <path>` walks a checkout and collects every such line into a report, and the livt repository's build derives each rule's status from the reports committed under `automations/`. A livt URI written without the marker is a reference, not a claim, and is not collected.

The report goes back as a **pull request**, not a push: a pull request is where a citation can be checked. It lands at `automations/{owner}/{repo}.json` — one file and one branch per implementation repository, force-updated rather than accumulating, because a report is a snapshot and the open pull request should carry the latest one rather than a queue of them. The credential lives on your side, one token with write access to the livt repository; the livt repository holds none, which is what keeps joining the ring a change to your repository alone.

### The workflow

Needs livt 0.15.0 or later, where `livt automations` and its `changed` subcommand arrived. Set `LIVT_REPOSITORY` to your livt repository, and put a token with write access to it in `LIVT_REPOSITORY_TOKEN` — a GitHub App installation token, or a fine-grained PAT with Contents and Pull requests write. If your livt repository *is* this repository, drop the second checkout and use the default `GITHUB_TOKEN` throughout.

```yaml
name: livt-automations

on:
  push:
    branches: [main]

permissions: {}

# One report branch per implementation repository, so two merges must not
# race to rewrite it.
concurrency:
  group: ${{ github.workflow }}
  cancel-in-progress: false

env:
  LIVT_REPOSITORY: your-org/your-livt-repository

jobs:
  collect:
    runs-on: ubuntu-latest
    permissions:
      contents: read # checkout
    steps:
      - uses: actions/checkout@v7
        with:
          persist-credentials: false

      - name: Install livt
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          set -euo pipefail
          gh release download --repo boykush/livt \
            --pattern 'livt_*_linux_amd64.tar.gz' --output - | tar -xz -C /usr/local/bin livt

      # livt answers whether this push could have changed what the report
      # claims. A base it cannot read answers true, so the fetch may fail.
      - name: Decide whether to collect
        id: gate
        env:
          BEFORE: ${{ github.event.before }}
          REMOTE: https://x-access-token:${{ secrets.GITHUB_TOKEN }}@github.com/${{ github.repository }}.git
        run: |
          set -euo pipefail
          git fetch --no-tags --depth=1 "$REMOTE" "$BEFORE" 2>/dev/null || true
          echo "collect=$(livt automations changed "$BEFORE")" >> "$GITHUB_OUTPUT"

      # Collected before the livt repository is checked out: an untracked
      # directory would leave the tree dirty, and a scan that cannot name its
      # revision drops every line URL with it.
      - name: Collect
        if: steps.gate.outputs.collect == 'true'
        run: livt automations . --out "$RUNNER_TEMP/report.json"

      - name: Checkout the livt repository
        if: steps.gate.outputs.collect == 'true'
        uses: actions/checkout@v7
        with:
          repository: ${{ env.LIVT_REPOSITORY }}
          token: ${{ secrets.LIVT_REPOSITORY_TOKEN }}
          path: livt-repository

      - name: Open the report's pull request
        if: steps.gate.outputs.collect == 'true'
        working-directory: livt-repository
        env:
          GH_TOKEN: ${{ secrets.LIVT_REPOSITORY_TOKEN }}
        run: |
          set -euo pipefail
          report="automations/$GITHUB_REPOSITORY.json"
          branch="automations/$GITHUB_REPOSITORY"
          mkdir -p "$(dirname "$report")"
          cp "$RUNNER_TEMP/report.json" "$report"
          # rev, generated_at and the urls built from rev move on every run, so
          # the claims are compared and the snapshot's own dating is not.
          claims() { jq -S '[.citations[].livt_uri] | unique' "$1"; }
          if git show "HEAD:$report" > "$RUNNER_TEMP/committed.json" 2>/dev/null \
             && [ "$(claims "$RUNNER_TEMP/committed.json")" = "$(claims "$report")" ]; then
            exit 0
          fi
          git config user.name 'github-actions[bot]'
          git config user.email '41898282+github-actions[bot]@users.noreply.github.com'
          git switch --create "$branch"
          git add "$report"
          git commit --message "chore: collect automations from $GITHUB_REPOSITORY"
          git push --force origin "HEAD:refs/heads/$branch"
          if [ -z "$(gh pr list --repo "$LIVT_REPOSITORY" --head "$branch" \
                       --state open --json number --jq '.[].number')" ]; then
            gh pr create --repo "$LIVT_REPOSITORY" --head "$branch" \
              --title "chore: collect automations from $GITHUB_REPOSITORY" \
              --body "Collected from $GITHUB_REPOSITORY at $GITHUB_SHA."
          fi
```

### Why the walk is gated, and why removal counts

Collecting walks every file. That is sub-second on a small repository and not free on a monorepo, so `livt automations changed <base>` answers whether walking could tell you anything new, and the walk is skipped when it cannot. The judgment is livt's rather than your CI's: livt is what decides what a marker is, so a step spelling `livt:automates` itself would keep a second copy of that definition, free to drift the moment the real one moves.

A line that merely moved is not a reason to walk, and the reasoning is worth holding onto because it is easy to get backwards. A report is a **dated snapshot, not a live view**: its `rev`, `file` and `line` all come from one revision, and the forge URL is pinned to that revision. A later commit that moves a cited line does not break the link — the report simply describes an earlier revision, which it says out loud. A rename is the same, and `changed` answers false for both.

Removal *is* a reason. A deleted test leaves the board claiming a rule is automated when nothing covers it any more, which is the board lying about the present rather than being merely out of date. One question over added and removed lines catches both, since a deleted file shows its marker lines as removals.

A range livt cannot read — a new branch, a force-push, a revision since gone — answers **true**, not false, and says why on stderr. The two mistakes are not the same size: guessing false drops a claim, while guessing true costs one walk. That is why the fetch of the base revision above is allowed to fail quietly.

When the gate does open, the whole repository is rescanned rather than the old report patched: the answer decides whether to walk, never what the report says.

The comparison that follows is on the claims, not on the bytes. A report never matches itself byte for byte — `rev`, `generated_at` and the urls built from `rev` all move on their own — so comparing files would open a pull request on every run. Comparing the set of cited livt URIs asks the question the pull request is actually about, and it is the same question the gate asked, one grain finer: could the claims have changed, and did they.

## What this plugin teaches, and what it leaves to you

Nothing and everything, respectively — and that is the clearest statement of the line the other two plugins draw in prose. livt ships the **read surface**: the tools, the resources, the `spec_version` on every payload, and the livt URI that addresses each rule, example, and question. How your agent builds against what it reads — the test framework, where a citation goes in a test comment, what counts as done — is your repository's, and belongs in its own `AGENTS.md` or a skill you write. A livt URI in a test comment is the whole contract between the two.

## Coming from 3.x

Nothing here changed. The three plugins carry one version between them, so this one moves when either of the others breaks — [livt-discovery](../livt-discovery/README.md) 4.0 split the formulation station out of the record, and [livt-delivery](../livt-delivery/README.md) 4.0 dropped the plan station.

## Coming from 1.x

The configuration is unchanged. [livt-delivery](../livt-delivery/README.md) 2.0 is where the change landed, and it is worth reading if this repository's agent also files or inspects issues.

Coming from `livt-mcp` 0.x: the configuration is unchanged there too; only the name moved from the mechanism to the station it serves.
