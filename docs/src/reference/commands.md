# Commands

## `livt serve`

Build artifacts and start a local server.

While the server is running, livt watches every directory the build reads
(`opportunities`, `discoveries/opportunity-canvases`,
`discoveries/example-mappings`, `stories`, `discoveries/usm`, and `ubiquitous`)
along with `livt.yaml`. When a file changes, livt rebuilds and reloads the page
in the browser automatically, so you can preview refinements while editing.

```bash
livt serve [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--port` | `-p` | `3000` | Port to listen on |
| `--out` | `-o` | `dist` | Output directory |
| `--diff` | | (off) | Also render the diff between two revisions — see [Reviewing a change](#reviewing-a-change) |

## `livt build`

Build static HTML from artifacts without starting a server.

```bash
livt build [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--out` | `-o` | `dist` | Output directory |
| `--diff` | | (off) | Also render the diff between two revisions — see [Reviewing a change](#reviewing-a-change) |

Both commands read [`livt.yaml`](./configuration.md) from the directory they run
in, which is where the site's language is set.

### Reviewing a change

A pull request against a livt repository is a YAML diff: indentation, list
ordering, quoting. `--diff` renders what changed as the spec instead, one
[livt URI](./uri.md) at a time.

```bash
livt serve --diff main
```

The argument is git's own two-dot range — `<base>..<head>` — or, as above, a
single revision against your working tree, which is the shape a review usually
takes: you are on the branch that makes the change. Revisions are read with
`git archive`, so the command must run inside the repository and a revision git
cannot resolve fails the build rather than producing a site.

The site is otherwise the site you always get. What is added is `diff.html`,
which groups every changed URI under the resource type it belongs to, and a mark
on each changed item where it lives — a rule's sticky on its board, a term's row
in the glossary — leading into the diff at that URI. The diff takes no entry in
the nav: it is not a kind of thing the livt repository holds, so the way in is
the thing that changed.

Every mark says one of three words. An item was **added**, it was **changed**,
or it was **withdrawn** — it stopped being spec. That is a reading of the spec
rather than of the file: a rule
[retired](../guides/example-mappings.md#retiring-an-item) keeps its entry and
edits one field, a deleted one loses its entry, and to a reader both mean the
rule is no longer asked for.

A withdrawal is drawn as the removal it is, whichever way the file recorded it:
retiring a rule *adds* the line that retires it, and drawn off the record it
would read as the opposite of what happened. The line that records the
withdrawal is left out — the mark says it once — and the successor stays, since
where the spec went is what the reader came for.

A withdrawn sticky goes back on its board, drawn the way a proposal is — pale,
dashed, struck through — since a proposal is not spec *yet* and this is not spec
*any more*. Only what left in this range comes back: a rule retired long ago
stays off the wall. A resource gone entirely has no page left to put anything
on, so its list says how many of its kind went.

## `livt mcp`

Run an MCP ([Model Context Protocol](https://modelcontextprotocol.io)) server
that exposes the livt repository (story maps, stories, example mappings, and
the ubiquitous language). An implementation repo's coding agent can then fetch
the spec for a story or rule without reading livt's source.

The livt repository usually lives in a separate checkout from the consumer, so point at
it with `--root` or the `LIVT_ROOT` environment variable. The flag takes
precedence; both default to the current directory.

```bash
livt mcp [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--root` | `$LIVT_ROOT`, then `.` | Path to the root of the livt repository |
| `--http` | (off; stdio) | Serve over Streamable HTTP at this address (e.g. `localhost:5488`) instead of stdio; the MCP endpoint is `<addr>/mcp` |

### Transports

By default the server runs over **stdio**, spawned per consumer — the client
launches `livt mcp` as a subprocess. stdio is the recommended form: the
binding to the livt repository is read at spawn time from the consumer's own
environment, so each workspace resolves against what it declares, and a
session never outlives the declaration it was started with.

Reach for **`--http`** when reachability is the problem — a client that
cannot spawn processes, or several repos sharing one machine-wide server
without each holding a checkout of the livt repository. It serves Streamable
HTTP from one long-running process:

```bash
livt mcp --http localhost:5488
```

Each consumer points its MCP client at `http://localhost:5488/mcp`. The server
is stateless and read-only, so one process backs many clients; keep `git pull`
current on its checkout and the served spec (and `spec_version`) updates live.
The binding trades the other way from stdio: the server fixes its repository
once, at start, for every client — not per workspace.
For distributing this client configuration to implementation repos, see the
[`livt-automation` plugin](https://github.com/boykush/livt/tree/main/plugins/livt-automation).
This mode assumes local use with no authentication — the server is meant to bind
to localhost, not a public network.

### Tools

| Tool | Arguments | Returns |
|------|-----------|---------|
| `list_opportunities` | — | Every [opportunity](../guides/opportunities.md) with its key, name, and statement. Each entry links to its opportunity resource (`uri`), to its canvas (`canvas_uri`) when one has been filled in, and to the story maps mapped for it. A missing canvas or story map is the record that the opportunity has not been taken that far. |
| `list_stories` | `opportunity` (optional) — an opportunity name, matched exactly against the name a story's opportunity chip carries; keeps only the stories on that map, and an unknown name yields an empty list | Every story with its key and name. Each entry links to its story resource (`uri`); stories that have an example mapping also include `example_mapping_uri`, and stories on a story map carry `opportunities` — one map name plus story map resource URI per map they sit on. |
| `list_example_mappings` | `opportunity` (optional) — as for `list_stories`; a mapping whose key sits on no map, as one with no story usually does, matches none | Every example mapping, including one whose story has no card, with its story key and mapping resource URI (`uri`). A mapping carries `name` when it has one of its own and `story_uri` when a story file exists — read that story for the name of a mapping with none of its own — and `opportunities`, the story maps its key sits on. |
| `list_story_maps` | — | Every story map with its name and its story map resource URI (`uri`). |
| `list_terms` | — | Every [ubiquitous language](../guides/ubiquitous-language.md) term with its key, display name, and term resource URI (`uri`) — including terms no board references. A term scoped to a context also carries `ctx`, which is part of what identifies it. |

### Resources

The spec itself is exposed as resources, addressable by URI
(opportunity → story map → story → mapping → rule → example, with the canvas,
questions, and ubiquitous terms linked alongside):

| URI | Returns |
|-----|---------|
| `livt://opportunity/{opportunity_key}` | An opportunity: its name, its statement (whose problem, and what the business gets from solving it), and its frontmatter meta. Carries `canvas_uri` when a canvas has been filled in, and `story_maps` — the maps whose key matches, as map name plus story map resource URI. |
| `livt://opportunity-canvas/{opportunity_key}` | The [Opportunity Canvas](../guides/opportunities.md#the-opportunity-canvas) filled in for an opportunity, as its ten `boxes` — each with its `key`, printed `number`, heading, the `prompt` it asks, and its `items`. Every box is returned, unanswered ones with an empty `items`: a blank box records a question the opportunity has not answered. |
| `livt://story-map/{map_name}` | A story map: activities, steps, story cards, and releases. Story cards that have a story file link to their story resource. `{map_name}` is the map's display name (percent-encoded) — the same identifier the build output uses for `story-map/{name}.html`. |
| `livt://story/{story_key}` | The story's name, body, and frontmatter meta (e.g. `issue`), plus `example_mapping_uri` when a mapping exists and `opportunities` — the story maps the story sits on, as map name plus story map resource URI. |
| `livt://mapping/{story_key}` | The story's example mapping (rules, examples, questions, ubiquitous terms), with its own `name` when it has one — the only name a mapping with no story file carries. Each rule, example, and question carries its own `uri`, and `ubiquitous_terms` resolves each referenced term to its resource URI. [Retired](../guides/example-mappings.md#retiring-an-item) entries are listed too, saying how they closed — the mapping is the structural record their ids are numbered from. |
| `livt://mapping/{story_key}/rule/{rule_id}` | A single rule and its examples, its `status` — `proposed` while it awaits agreement, `accepted` once agreed, `rejected` or `retired` once closed, and present on every rule — plus its recorded automation: `issues` (automation Issue URLs) and `automated` (whether the rule is automated by tests). Rules inside `livt://mapping/{story_key}` carry the same fields. |
| `livt://mapping/{story_key}/rule/{rule_id}/example/{example_id}` | A single example of a rule. Example ids are numbered within their rule, so the address carries `{rule_id}` — `EX-01` alone does not identify an example. |
| `livt://mapping/{story_key}/question/{question_id}` | A single question. Questions hang off the mapping rather than off a rule, so the address stops at `{story_key}`. |
| `livt://ubiquitous/{term_key}` | A ubiquitous language term's name and definition. This shape addresses a term whose meaning holds across contexts. |
| `livt://ubiquitous/{ctx}/{term_key}` | A term scoped to one [context](../guides/ubiquitous-language.md#contexts), carrying `ctx` alongside its key. A context is optional and part of the address, so the same `{term_key}` can name one term at the root and another inside a context; the two never resolve to each other. |

A closed rule and a retired example or question keep resolving by their URI and
say so as they do in the YAML — a rule through its `status`, an example or a
question through `retired: true` — so a reference to one reads as closed rather
than failing (or, worse, landing on whatever reused its id). A live example or
question omits the flag. When something took the item's place, `superseded_by`
lists it as livt URIs, so the reference reads on instead of stopping; a
retirement that replaced nothing omits that field too.

Read them with `resources/read`; all appear in `resources/templates/list`. The
server advertises templates only — there is no concrete resource list and no
change notification (subscribe); every read is served fresh from disk.

Every tool and resource payload also includes a `spec_version` field -- the
short git revision of the livt repository -- so consumers can tell which version of the
spec they are reading and detect drift.

### Citing the livt repository

The handshake carries `instructions`, so every session tells the consuming agent
how to reference the livt repository in what it produces: copy the `uri` from the result
verbatim, and never write a bare `id`. Ids are unique only within one mapping
file, so `R-02` exists in every mapping and identifies nothing on its own.

```go
// livt://mapping/place-order-with-saved-card/rule/R-13/example/EX-01
```

That applies wherever a reference leaves the livt repository -- a test comment, an issue
body, a commit message, a PR description. A published living-document URL is a
convenience link for humans, not the citation form: it depends on where the site
is deployed, and the livt URI does not.

A URI cited this way is read back with [`livt resolve`](#livt-resolve), which
needs no MCP client -- so the citation stays followable for CI, an editor, or a
reader who only has a checkout.

## `livt resolve`

Resolve a [livt URI](./uri.md) against the livt repository, without running an
MCP client. A rule, example, or question cited in a test comment, an issue body,
or a commit message can then be read back by whoever needs it — CI, an editor,
or an agent that is not connected to livt over MCP.

```bash
livt resolve <uri> [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--root` | `$LIVT_ROOT`, then `.` | Path to the root of the livt repository |
| `--format` | `json` | Output form: `json` or `url` |
| `--base-url` | — | Root of the deployed site; required by `--format url` |

Every URI shape the MCP server exposes as a resource resolves here: mappings,
rules, examples, questions, stories, story maps, and ubiquitous language terms.

### Output forms

**`--format json`** (the default) prints the same payload an MCP
`resources/read` of that URI serves, `spec_version` included — the two surfaces
resolve through one code path, so a consumer sees one shape whichever way it
asked:

```bash
livt resolve livt://mapping/trace-test-to-rule/rule/R-04
```

**`--format url`** prints the item's page on the deployed site, using the same
[URI-to-page derivation](./uri.md#uri-to-page) the site build anchors its
stickies to:

```bash
$ livt resolve livt://mapping/trace-test-to-rule/rule/R-04 \
    --format url --base-url https://boykush.github.io/livt
https://boykush.github.io/livt/mapping/trace-test-to-rule.html#rule-R-04
```

The livt repository is still read in this form, even though the page path derives from
the URI alone: a link to a page that was never built is worse than an error.

### When a URI does not resolve

Both cases exit non-zero and explain themselves on stderr, because they need
different fixes:

| | |
|---|---|
| **Malformed URI** | The string is not a livt URI at all. The error lists every shape one can take. Fix the URI. |
| **Nothing to resolve** | The URI is well formed but the livt repository holds no such item, e.g. `rule "R-99" not found in story "trace-test-to-rule"`. Fix the reference, or add the item. |

## `livt version`

Print the version of livt.

```bash
livt version
```
