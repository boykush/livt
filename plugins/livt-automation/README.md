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

2. **Install this plugin** in each implementation repository:

   ```
   /plugin install livt-automation@boykush/livt
   ```

The agent can then call `list_stories` / `list_story_maps` / `list_terms` and read resources such as `livt://story-map/{map_name}`, `livt://story/{story_key}`, `livt://mapping/{story_key}`, and `livt://ubiquitous/{term_key}`. See the [`livt mcp` command reference](https://github.com/boykush/livt/blob/main/docs/src/reference/commands.md) for the full tool and resource list. A rule whose `status` is `proposed` is a candidate awaiting agreement, not spec to automate yet.

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

Coming from `livt-mcp` 0.x: the configuration is unchanged; only the name moved from the mechanism to the station it serves.
