# Configuration

`livt.yaml` at the root of the livt repository configures the site build:

```yaml
lang: ja
```

The file is optional and every field falls back to a default, so a repository
without one builds exactly as it did before there was one.

| Field | Default | Values | Description |
|-------|---------|--------|-------------|
| `lang` | `en` | `en`, `ja` | Language of the site chrome |

Both [`livt build`](./commands.md#livt-build) and
[`livt serve`](./commands.md#livt-serve) read it from the directory the command
runs in — the same place `stories/` and `discoveries/` are resolved from.
`livt serve` watches it like any other input, so an edit rebuilds and reloads
the page you have open.

A value livt does not support is an error rather than a silent fall back to the
default, since the alternative is a site built in a language nobody asked for:

```console
$ livt build
Error: livt.yaml: unknown lang "de" (supported: en, ja)
```

## `lang`

`lang` is the language of everything livt renders *of itself*: the sidebar,
page headings and titles, board legends, filter chips, table headers, empty
states, the `lang` attribute of every page, and the headings and prompts of the
[Opportunity Canvas](../guides/opportunities.md#the-opportunity-canvas) — the
questions printed on the sheet, which are livt's words rather than yours.

It does not touch the livt repository's own prose. Story names and bodies,
rules, examples, questions, and term definitions are rendered as written,
whichever language they are written in. The two are independent on purpose — a
team writing its stories in Japanese and reading the site in English is as
ordinary an arrangement as the reverse, and the setting is about the frame, not
the contents.

The setting reaches the site only. `livt mcp` and `livt resolve` answer in
English whatever it says, canvas headings included: those payloads are a
contract read by agents, and a field that changed language with a site setting
would not be one. What they carry *from* the livt repository — names, rules,
definitions — is verbatim there too.

Two languages are supported today, `en` and `ja`. A new one is a catalog in
`internal/i18n/`; every catalog answers for every message, so a language is
either complete or not offered.
