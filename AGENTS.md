# AGENTS.md

livt turns a livt repository — opportunities, story maps, example mappings, ubiquitous language — into a living document (static site, CLI, MCP server). The repo dogfoods itself: `opportunities/`, `discoveries/`, `stories/`, and `ubiquitous/` hold a real livt repository describing livt, rendered by livt.

Setup, checks, and the conventional-commit rule are in [CONTRIBUTING.md](CONTRIBUTING.md); this file adds only what is specific to working here. Run `mise run check` before pushing.

## Language

Two languages, split by what the text *is* — not by who wrote it.

**English** — everything that is the tool or the process:

- PR titles and bodies, commit subjects and bodies, issues
- Go code and its comments, `docs/src/` (the user-facing mdBook), `plugins/*/skills/**`
- Every identifier: `key:`, `id:`, `release:`, file names, branch names

**Japanese** — the prose of the livt repository itself:

- `opportunities/*.md` — `name:` and the opportunity's statement
- `stories/*.md` — `name:` and the user-story body
- `discoveries/**/*.yaml` — `name:` and `text:`
- `ubiquitous/*.md` — `name:` and the definition
- `internal/i18n/ja.go` — the Japanese chrome catalog. Go code holding Japanese by design; its keys, like every identifier, stay English
- `docs/po/ja.po` — the Japanese translation of `docs/src/`. Its `msgid`s are the English verbatim and only its `msgstr`s are Japanese; code, identifiers and names repeat their `msgid` unchanged

The livt repository's language is not a rule livt imposes. The skills say prose follows the language the user is speaking — see any `## Language` section under `plugins/livt-discovery/skills/`. *This* livt repository happens to be Japanese, and stays that way for consistency.

Consequences worth stating:

- Quoting the livt repository inside English text keeps the Japanese verbatim. A commit body naming a story map writes 協働ディスカバリー, not a translation of it.
- What livt renders *of itself* — nav labels, headings, empty states — follows `lang:` in [livt.yaml](livt.yaml), which this repository sets to `ja`, matching the artifacts the site renders. That is a per-repository setting and not a property of livt, whose own default is `en`. So a string the site shows is never a literal in a template or in `internal/domain/`: it goes in `internal/i18n/en.go` and is owed a translation in `ja.go`, and the templates reach it with `{{t "…"}}`. The one exception is the wordmark and its tagline — livt's own name and slogan, English wherever livt is named.
- `ja.go`, `docs/po/ja.po` and `ubiquitous/` name the same things and must agree. A term renamed in one is renamed in the others, or a board ends up labelled 具体例 above a glossary that calls it 実例 — the drift ubiquitous language exists to prevent, on the site whose job is to show it.
- A change to `docs/src/` is owed its translation in `docs/po/ja.po`; [CONTRIBUTING.md](CONTRIBUTING.md#docs) has the steps.

[CONTRIBUTING.md](CONTRIBUTING.md) welcomes issues and pull requests "in English or Japanese". That is an invitation to human contributors, and it stands. The rules above govern agents working in this repo.

## Documentation

- **`docs/` never restates what livt already says elsewhere.** Each kind of detail has one home, beside what it describes — [Reference](docs/src/reference.md) says where — and a change in behaviour updates that home in the same diff, never a page in `docs/`.
- **`docs/` says why livt exists and when it fits by pointing at the record.** livt's reasons are its opportunities and story maps, and the introduction links them on the live demo instead of summarising them: a summary there fell behind once, still listing two opportunities after a third was taken on.
- **Explaining what livt does starts from the example mapping that decided it** — the files, or `go run . resolve livt://mapping/{story-key}` — cited by livt URI. Code that disagrees with it is a gap to report, not a document to update.

## Commits and PRs

- livt repository changes (`opportunities/`, `discoveries/`, `stories/`, `ubiquitous/`) ship as `docs:`. The branch prefix mirrors the type: `docs/…`, `feat/…`, `fix/…`.
- The body explains **why** the change is right — the diff already says what changed. Wrap at ~80 columns.
- One rule-level change per commit: `docs: add rule R-11 to file-automation-issues-to-impl-repos`. The ID and commit contracts live in [change-rule/SKILL.md](plugins/livt-discovery/skills/change-rule/SKILL.md) — a filed rule ID is immutable.
- How those commits are grouped into a PR is this repository's call, not livt's: a session's worth of work ships as one PR on one branch, however many commits it took. Never fold two decisions into one commit to make that grouping tidier.
- A commit automating a rule cites it by livt URI: `Automates livt://mapping/{story-key}/rule/{rule-id}.` A bare `R-02` exists in every mapping file and so identifies nothing.
- A plugin's version moves in a commit of its own, merged last — after the changes it covers and immediately before the release tag, so the three plugins carry one number per release rather than one per pull request (`plugins/*/.claude-plugin/plugin.json` and `.claude-plugin/marketplace.json`, which they share). The marketplace serves them by path, so merging is releasing, and every hour between a breaking change and its bump is an hour of skills labelled with a number they have outgrown: keep that window inside the release, and never cut a release without closing it.
