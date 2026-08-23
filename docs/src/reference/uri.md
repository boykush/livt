# livt URI

A **livt URI** names one point in the livt repository — an opportunity, its
canvas, a rule, an example, a question, a story, a story map, or a ubiquitous
language term. It is the form to
reach for whenever a reference has to survive outside the livt repository: a test comment,
an issue body, a commit message.

```
livt://mapping/{story-key}/rule/{rule-id}
```

The MCP server addresses its resources by these URIs — see
[Commands](./commands.md#resources) for what each one returns — and the
generated site anchors its stickies to them. Neither depends on where the site
is deployed.

## Why not a bare ID

`R-02` on its own is not a reference. IDs are numbered within the file that
holds them, so every example mapping in this repository defines an `R-02`, and
every rule that has examples defines an `EX-01`. Even `R-13 EX-01` narrows it to
nothing — the reader still cannot tell which mapping was meant, and neither can
a search. The story key is the part that makes it addressable:

```
livt://mapping/trace-test-to-rule/rule/R-02/example/EX-01
```

## Relative to a repository

The widening stops there: a livt URI never says **which** livt repository
it belongs to. Its interpretation is uniform — one point in the livt
repository at hand — the same shape `http://localhost/` has, where every
reader agrees on the meaning while the referent depends on where it is
resolved. Which repository answers is supplied by the consumer's own
declaration (`--root` / `LIVT_ROOT`), so scope lives in the workspace, not
in the name. That is what keeps a citation short, and what lets the same
reference work in every checkout of the same repository.

The first segment is drawn from a closed set — `opportunity`,
`opportunity-canvas`, `mapping`, `story`, `story-map`, `ubiquitous` — so a future form that names a repository in
that position stays open without breaking any URI written today.

## URI to page

`livt build` renders the livt repository as a static site, and this is where each URI
lands in it. Paths are relative to the output root.

| livt URI | Page |
|---|---|
| `livt://mapping/{story-key}` | `mapping/{story-key}.html` |
| `livt://mapping/{story-key}/rule/{rule-id}` | `mapping/{story-key}.html#rule-{rule-id}` |
| `livt://mapping/{story-key}/rule/{rule-id}/example/{example-id}` | `mapping/{story-key}.html#rule-{rule-id}-example-{example-id}` |
| `livt://mapping/{story-key}/question/{question-id}` | `mapping/{story-key}.html#question-{question-id}` |
| `livt://story/{story-key}` | `story/{story-key}.html` |
| `livt://opportunity/{opportunity-key}` | `opportunity/{opportunity-key}.html` |
| `livt://opportunity-canvas/{opportunity-key}` | `opportunity-canvas/{opportunity-key}.html` |
| `livt://story-map/{map-name}` | `story-map/{map-name}.html` |
| `livt://ubiquitous/{term-key}` | `ubiquitous.html#{term-key}` |
| `livt://ubiquitous/{ctx}/{term-key}` | `ubiquitous.html#{ctx}/{term-key}` |

An example's anchor repeats its rule for the same reason its URI does: `EX-01`
recurs under every rule of a board, so `#example-EX-01` would land on whichever
one happened to be first. The sticky's own badge still shows the local `EX-01` —
that is the ID the livt repository numbers — while the link behind it carries the rule.

A term takes a [context](../guides/ubiquitous-language.md#contexts) the same
way, and for the same reason: `invoice` can mean one thing in billing and
another in shipping, so the context is part of the address rather than a label
on it. The context is optional — a term whose meaning holds everywhere keeps the
one-segment form — and the two shapes never resolve to each other.

An opportunity is addressed by key rather than by display name, unlike a story
map. The key is what its filename and its canvas already join on, so it is the
identifier the livt repository carries. The canvas sits beside the opportunity
rather than under it, the way a mapping sits beside its story — the two are
joined by key, and either can exist without the other.

A [retired](../guides/example-mappings.md#retiring-an-item) item has no sticky,
so its URI lands on the board with nothing to scroll to. The URI still resolves:
ask the tooling, which answers with the item and `retired: true`. Deriving a
page is how a reference is *shown*, not how it is *resolved*.

## Store the URI, render the URL

Nothing in the livt repository records where the site is deployed. The deployment URL is
prefixed onto the paths above only while a page is being rendered, which is why
the same livt repository can be served from a local `livt serve`, from GitHub Pages, and
from an internal host at once:

```
livt://mapping/trace-test-to-rule/rule/R-02
  -> mapping/trace-test-to-rule.html#rule-R-02
  -> https://boykush.github.io/livt/demo/mapping/trace-test-to-rule.html#rule-R-02
```

A deployed URL is therefore a fine thing to paste into a chat and a poor thing to
commit: it goes stale the moment the site moves, and it buries the identity it
was meant to carry. Keep the URI; let the URL be derived.

## Where each form comes from

- **A URL, to send to a person** — take it from the board. Every sticky shows its
  own ID in the bottom-right corner; clicking it copies the deployed URL of that
  sticky, so what you paste already points at the exact card.
- **A livt URI, to leave in an artifact** — take it from the tooling rather than
  the board. Every rule, example, and question the MCP server returns carries
  its own `uri`, so an agent writing a test can quote the point it is proving
  without composing the URI by hand.

Keeping the two apart is what stops the site from becoming a second source of
identity. The board is for reading and sharing; the URI is what the livt repository
actually knows about itself.

## Reading one back

A URI left in a test comment is only worth as much as the ability to follow it,
and that has to work for a reader who is not running an MCP client — CI, an
editor, or someone with just a checkout. [`livt resolve`](./commands.md#livt-resolve)
turns any of the shapes above back into the point it names:

```bash
$ livt resolve livt://mapping/trace-test-to-rule/rule/R-02
{ "spec_version": "...", "rule": { "id": "R-02", ... } }

$ livt resolve livt://mapping/trace-test-to-rule/rule/R-02 \
    --format url --base-url https://boykush.github.io/livt
https://boykush.github.io/livt/mapping/trace-test-to-rule.html#rule-R-02
```

The URL form is derived from the table above rather than restated, so a link the
CLI hands out and the anchor the build writes cannot drift apart.
