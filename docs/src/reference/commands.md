# Commands

## `livt serve`

Build artifacts and start a local server.

While the server is running, livt watches the input directories
(`discoveries/example-mappings`, `stories`, `discoveries/usm`, and `ubiquitous`).
When a file changes, livt rebuilds and reloads the page in the browser
automatically, so you can preview refinements while editing.

```bash
livt serve [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--port` | `-p` | `3000` | Port to listen on |
| `--out` | `-o` | `dist` | Output directory |

## `livt build`

Build static HTML from artifacts without starting a server.

```bash
livt build [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--out` | `-o` | `dist` | Output directory |

## `livt mcp`

Run an MCP ([Model Context Protocol](https://modelcontextprotocol.io)) server
that exposes the discovery master (stories and example mappings). An
implementation repo's coding agent can then fetch the spec for a story or rule
without reading livt's source.

The master usually lives in a separate checkout from the consumer, so point at
it with `--root` or the `LIVT_ROOT` environment variable. The flag takes
precedence; both default to the current directory.

```bash
livt mcp [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--root` | `$LIVT_ROOT`, then `.` | Path to the livt project root holding the discovery master |
| `--http` | (off; stdio) | Serve over Streamable HTTP at this address (e.g. `localhost:5488`) instead of stdio; the MCP endpoint is `<addr>/mcp` |

### Transports

By default the server runs over **stdio**, spawned per consumer — the client
launches `livt mcp` as a subprocess. This suits a single repo with the master
checked out alongside it.

Pass **`--http`** to instead serve over Streamable HTTP from one long-running
process, so several repos on the same machine can share a single server without
each holding a checkout of the master:

```bash
livt mcp --http localhost:5488
```

Each consumer points its MCP client at `http://localhost:5488/mcp`. The server
is stateless and read-only, so one process backs many clients; keep `git pull`
current on its checkout and the served spec (and `spec_version`) updates live.
For distributing this client configuration to implementation repos, see the
[`livt-mcp` plugin](https://github.com/boykush/livt/tree/main/plugins/livt-mcp).
This mode assumes local use with no authentication — the server is meant to bind
to localhost, not a public network.

### Tools

| Tool | Arguments | Returns |
|------|-----------|---------|
| `list_stories` | — | Every story with its key and name. Stories that have an example mapping include `example_mapping_uri`, the resource URI below. |

### Resources

The spec itself is exposed as resources, addressable by URI (story → mapping → rule):

| URI | Returns |
|-----|---------|
| `livt://mapping/{story_key}` | The story's example mapping (rules, examples, questions, ubiquitous terms). Each rule carries its own `uri`. |
| `livt://mapping/{story_key}/rule/{rule_id}` | A single rule and its examples. |

Read them with `resources/read`; both appear in `resources/templates/list`.

Every tool and resource payload also includes a `spec_version` field -- the git
revision of the master -- so consumers can tell which version of the spec they
are reading and detect drift.

## `livt version`

Print the version of livt.

```bash
livt version
```
