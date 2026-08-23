# AGENTS.md

livt turns a livt repository — story maps, example mappings, ubiquitous language — into a living document (static site, CLI, MCP server). The repo dogfoods itself: `discoveries/`, `stories/`, and `ubiquitous/` hold a real livt repository describing livt, rendered by livt.

Setup, checks, and the conventional-commit rule are in [CONTRIBUTING.md](CONTRIBUTING.md); this file adds only what is specific to working here. Run `mise run check` before pushing.

## Language

Two languages, split by what the text *is* — not by who wrote it.

**English** — everything that is the tool or the process:

- PR titles and bodies, commit subjects and bodies, issues
- Go code and its comments, `docs/src/` (the user-facing mdBook), `plugins/*/skills/**`
- Every identifier: `key:`, `id:`, `release:`, file names, branch names

**Japanese** — the prose of the livt repository itself:

- `stories/*.md` — `name:` and the user-story body
- `discoveries/**/*.yaml` — `name:` and `text:`
- `ubiquitous/*.md` — `name:` and the definition
- `internal/i18n/ja.go` — the Japanese chrome catalog. Go code holding Japanese by design; its keys, like every identifier, stay English

The livt repository's language is not a rule livt imposes. The skills say prose follows the language the user is speaking — see any `## Language` section under `plugins/discovery-facilitator/skills/`. *This* livt repository happens to be Japanese, and stays that way for consistency.

Two consequences worth stating:

- Quoting the livt repository inside English text keeps the Japanese verbatim. A commit body naming a story map writes 協働ディスカバリー, not a translation of it.
- What livt renders *of itself* — nav labels, headings, empty states — is English here, set by `lang: en` in [livt.yaml](livt.yaml). That is this repository's choice, not a property of livt: the chrome has one catalog per language in `internal/i18n/`, and a message added to one is owed to all. The Japanese on a livt site belongs to the livt repository being rendered, not to the tool.

[CONTRIBUTING.md](CONTRIBUTING.md) welcomes issues and pull requests "in English or Japanese". That is an invitation to human contributors, and it stands. The rules above govern agents working in this repo.

## Commits and PRs

- livt repository changes (`discoveries/`, `stories/`, `ubiquitous/`) ship as `docs:`. The branch prefix mirrors the type: `docs/…`, `feat/…`, `fix/…`.
- The body explains **why** the change is right — the diff already says what changed. Wrap at ~80 columns.
- One rule-level change per PR, a branch each: `docs: add rule R-11 to file-automation-issues-to-impl-repos`. The ID and PR contracts live in [example-mapping-update/SKILL.md](plugins/discovery-facilitator/skills/example-mapping-update/SKILL.md) — a filed rule ID is immutable.
- A commit automating a rule cites it by livt URI: `Automates livt://mapping/{story-key}/rule/{rule-id}.` A bare `R-02` exists in every mapping file and so identifies nothing.
