# livt MCP Plugin

Connect an implementation repo's coding agent to a locally running livt MCP server, so it can fetch the discovery spec — stories, example mappings, and rules — straight from the master instead of reading livt's source or a stale copy.

## Overview

This is the **consumer-side** counterpart to `discovery-facilitator`: that plugin authors the discovery master; this one lets your repo read it. It bundles only an MCP server configuration — no skills or agents — pointing your agent at a livt server you run locally.

```
 livt mcp --http localhost:5488      ← one local server, holds the master
        │  /mcp
        ├── repo A  (livt-mcp plugin)
        ├── repo B  (livt-mcp plugin)
        └── repo C  (livt-mcp plugin)
```

One server backs every repo on your machine: no per-repo checkout of the master, no `--root` to configure. Each tool and resource payload carries a `spec_version` (the master's git revision) so your agent can detect drift.

## Setup

1. **Run the server** from your livt master checkout (one long-running process):

   ```bash
   livt mcp --http localhost:5488
   ```

   Keep `git pull` current on that checkout — the spec and its `spec_version` are read per request, so updates are served live.

2. **Install this plugin** in each implementation repo:

   ```
   /plugin install livt-mcp@boykush/livt
   ```

The agent can then call `list_stories` and read `livt://mapping/{story_key}` resources. See the [`livt mcp` command reference](https://github.com/boykush/livt/blob/main/docs/src/reference/commands.md) for the full tool and resource list.

## Configuration

The server URL defaults to `http://localhost:5488/mcp`. To use a different port or host, set `LIVT_MCP_URL` before launching your agent:

```bash
export LIVT_MCP_URL=http://localhost:8398/mcp
```

This setup assumes local use with no authentication; the server binds to localhost and is not meant to be exposed to a network.

## Using stdio instead (no server)

For a single repo without running a server, connect over stdio: the client spawns `livt mcp` as a subprocess. This does **not** use this plugin — add the server directly to the repo's `.mcp.json`, pointing `LIVT_ROOT` at your master checkout:

```json
{
  "mcpServers": {
    "livt": {
      "type": "stdio",
      "command": "livt",
      "args": ["mcp"],
      "env": { "LIVT_ROOT": "/path/to/master" }
    }
  }
}
```

stdio needs the master's path per repo (`LIVT_ROOT`, or `--root`), so — unlike the shared HTTP URL — it can't ship as a turnkey plugin config. Hence it's documented here rather than bundled as a plugin.
