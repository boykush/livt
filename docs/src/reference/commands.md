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
that exposes the discovery master (story maps, stories, example mappings, and
the ubiquitous language). An implementation repo's coding agent can then fetch
the spec for a story or rule without reading livt's source.

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
| `list_stories` | `opportunity` (optional) — an opportunity name, matched exactly against a story map display name; keeps only the stories on that map, and an unknown name yields an empty list | Every story with its key and name. Each entry links to its story resource (`uri`); stories that have an example mapping also include `example_mapping_uri`, and stories on a story map carry `opportunities` — one map name plus story map resource URI per map they sit on. |
| `list_story_maps` | — | Every story map with its name and its story map resource URI (`uri`). |

### Resources

The spec itself is exposed as resources, addressable by URI
(story map → story → mapping → rule, with ubiquitous terms linked from
mappings and story maps):

| URI | Returns |
|-----|---------|
| `livt://story-map/{map_name}` | A story map: activities, steps, story cards, and releases. Committed story cards link to their story resource. `{map_name}` is the map's display name (percent-encoded) — the same identifier the build output uses for `story-map/{name}.html`. |
| `livt://story/{story_key}` | The story's name, body, and frontmatter meta (e.g. `issue`), plus `example_mapping_uri` when a mapping exists and `opportunities` — the story maps the story sits on, as map name plus story map resource URI. |
| `livt://mapping/{story_key}` | The story's example mapping (rules, examples, questions, ubiquitous terms). Each rule carries its own `uri`, and `ubiquitous_terms` resolves each referenced term to its resource URI. |
| `livt://mapping/{story_key}/rule/{rule_id}` | A single rule and its examples, plus its recorded automation: `issues` (automation Issue URLs) and `automated` (whether the rule is automated by tests). Rules inside `livt://mapping/{story_key}` carry the same fields. |
| `livt://ubiquitous/{term_key}` | A ubiquitous language term's name and definition. |

Read them with `resources/read`; all appear in `resources/templates/list`. The
server advertises templates only — there is no concrete resource list and no
change notification (subscribe); every read is served fresh from disk.

Every tool and resource payload also includes a `spec_version` field -- the
short git revision of the master -- so consumers can tell which version of the
spec they are reading and detect drift.

## `livt version`

Print the version of livt.

```bash
livt version
```
