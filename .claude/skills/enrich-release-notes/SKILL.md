---
name: enrich-release-notes
description: Enrich a published livt release's notes for the people who use livt — who is affected and what to do, the highlights, and the upgrade steps — above the changelog GoReleaser generated, which stays byte-for-byte. Use at release time, once the release workflow has published a v* tag, or to rework an earlier release's notes. Publishes only after the maintainer approves the draft.
---

You **enrich a release's notes**. GoReleaser publishes every `v*` tag with a changelog grouped by commit type: complete, and written for whoever reads commits. You write the part above it for whoever uses livt — what changed for them, and what they have to do about it.

## Language

The notes are English, like everything that is the tool or the process (see `AGENTS.md`). Quote this repository's livt repository verbatim — a term from `ubiquitous/` stays in its own words. Talk to the maintainer in the language they are using.

## When This Runs

After the release workflow has published the tag: at release time, or later, to rework an older release. Take the tag from the maintainer; without one, use the latest release (`gh release view --json tagName --jq .tagName`).

Don't create the release ahead of the tag push to write into it. `.goreleaser.yaml` sets `release.mode: keep-existing`, so GoReleaser would keep that body and skip its own changelog.

## What livt Ships

Each reader uses one or more of these, and the notes answer for each:

- **The `livt` binary** — the CLI, the site it builds, and `livt mcp`, built from the tag. It changes only when `main.go`, `cmd/`, `internal/`, `go.mod`, or `go.sum` do.
- **The reader's own livt repository** — the files the binary reads. A renamed field, a removed flag, or a new file layout breaks it even when every command still runs.
- **The Claude Code plugins** — `livt-discovery`, `livt-delivery`, `livt-automation`, served from `main` by the marketplace in `.claude-plugin/marketplace.json`, by path rather than by tag. Their shared version is in `plugins/*/.claude-plugin/plugin.json`.
- **The docs site** — `docs/src/`, deployed from `main`: why livt exists and when it fits. It keeps no reference, so a release rarely moves it.

This repository's own livt repository (`discoveries/`, `stories/`, `ubiquitous/`, `opportunities/`) is livt describing itself: mention it only when a term the plugins speak moved with it. CI, agent config, and dependency bumps stay in the changelog alone.

## Gather

Let `$prev` be the tag before `$tag`: `git describe --tags --abbrev=0 "$tag^"`.

1. Save the published body: `gh release view "$tag" --json body --jq .body`. Everything from `## Changelog` down, the attestation footer included, is GoReleaser's and goes back unchanged. Keep the copy — publishing it again is how the edit is undone.
2. Read every commit in the range with its body: `git log --no-merges --format='%h %s%n%n%b' "$prev..$tag"`. The body says why. A `BREAKING CHANGE:` footer is the maintainer's own statement of what broke and what replaces it.
3. Check each deliverable against the diff, not against the subjects:
   - binary — `git diff --stat "$prev..$tag" -- main.go cmd internal go.mod go.sum ':(exclude)*_test.go' ':(exclude)**/testdata/**'`. Empty means no code changed, and the notes say so.
   - the reader's livt repository — `BREAKING CHANGE:` footers that name a field, a file, or a flag, and the diff under `plugins/*/skills/`, where the formats it holds are written down.
   - plugins — `git diff --stat "$prev..$tag" -- plugins .claude-plugin`, and the version in each `plugin.json` at both tags. If `plugins/` changed and the version did not move, stop and tell the maintainer: by `AGENTS.md` it moves immediately before the tag, and the upgrade steps cannot promise a number the plugins do not carry.
   - docs — `git diff --stat --diff-filter=A "$prev..$tag" -- docs/src`, for new pages worth a link.
4. Read the migration notes the maintainers already wrote: each plugin README's `Coming from N.x` section at `$tag`. Link them rather than re-deriving them.

Describe the release at `$tag`. What `main` gained since belongs to the next release's notes.

## Write

Above `## Changelog`, in this order:

1. **One or two sentences** on what the release is about.
2. **A table** — `| If you use | What changed | What to do |` — with a row for the binary and one per plugin. A deliverable that did not change gets a row saying `Nothing`, so no reader has to infer it from an absence.
3. **`## Highlights`** — a `###` per change a reader will notice, breaking changes first. The heading names the change in the reader's terms — `/plan-story` is removed — not in the commit's. A renamed, split, or removed command gets a before/after table. Close each with a `More:` link to the README section that covers it, or to the page on the live demo that shows it.
4. **`### Also in this release`** — a short bullet per smaller change a reader can see: a new option, a fix to behaviour they may have hit, a new docs page.
5. **`## Upgrading`** — `### The livt binary` and `### Plugins`, each with its steps or `Nothing to do`.
   - Binary: link the README's [Getting started](https://github.com/boykush/livt/blob/$tag/README.md#getting-started) rather than repeating it, and give every change to the reader's livt repository as a before/after table.
   - Plugins: the commands below, then a restart of Claude Code; `claude plugin list` shows the version to expect. Name each plugin as the CLI does, `<plugin>@<marketplace>`, where `<marketplace>` is the `name` in `.claude-plugin/marketplace.json`.

     ```bash
     claude plugin marketplace update <marketplace>
     claude plugin update <plugin>@<marketplace>
     ```

   - Then what the reader changes in their own workflow: a command to stop calling, a new one to run after another.

`gh release view v0.14.0` is a worked example of the shape.

### Style

- **No hard wraps.** GitHub renders a single newline in release notes as a line break, so a paragraph wrapped at 80 columns shows as ragged lines. One paragraph, one line.
- **Pin links to the tag** — `https://github.com/boykush/livt/blob/$tag/…` — so they keep saying what the release said after `main` moves.
- **Spell things as they are typed**: commands, fields, and files in code spans, exactly.
- **Lead with what to do.** The why is a sentence; the commit bodies and the READMEs hold the rest, and the notes link to them.
- **Only what you can trace.** Every claim comes from a commit, a diff, or a README at `$tag` — never from a subject line alone.

## Check

Before showing the draft:

- Every link resolves: `curl -s -o /dev/null -w '%{http_code}' -L <url>` prints `200`.
- Every command matches the tool that runs it; read its `--help` when unsure.
- Everything from `## Changelog` down is byte-for-byte the saved body.

## Publish

Editing a release is public. Show the draft as a file, and publish only on the maintainer's go-ahead:

```bash
gh release edit "$tag" --notes-file <draft>
```

Then read it back. `gh release view "$tag" --json body --jq .body` matches the draft, and the rendered page — `gh api "repos/boykush/livt/releases/tags/$tag" -H 'Accept: application/vnd.github.html+json' --jq .body_html` — has its tables and no `<br>` above the changelog.
