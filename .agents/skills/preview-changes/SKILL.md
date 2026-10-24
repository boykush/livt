---
name: preview-changes
description: Serve the livt site with the diff of the current branch and open it in a browser, so the maintainer sees what the change does to the board — the diff page, and each changed page with its mark. Use right after editing the livt repository (opportunities/, discoveries/, stories/, design/, livt.yaml) or anything that changes how the site renders (internal/builder, templates, i18n), and whenever asked to preview, show, or diff a change.
---

You **preview a change** the way a reviewer of it would read it: as livt renders it, not as a YAML diff. `livt serve --diff <base>` builds the site with a diff page at `diff.html`, marks every changed page, and rebuilds and reloads as the files change. What the diff shows is decided in `livt://mapping/review-diff-between-revisions`; this skill only opens it.

## Start

Start the `livt-diff` preview (`preview_start` with that name). Its command lives in `.claude/launch.json`: the base is the merge-base with `main`, and the head is the working tree, so the diff is everything this branch does, committed or not — main's newer commits never show up reversed. It serves on its own port and output directory, beside the plain `livt-serve` the session already started.

Without the desktop browser pane (Codex, a terminal session), run that same command in the background and open its URL with the Playwright MCP server, or give the maintainer the URL.

Reuse a preview that is already running: it rebuilds on every edit to the livt repository. A change to Go code is not one of its inputs, so restart the preview after editing `internal/` or `cmd/`.

## Walk the Change

1. Open `/diff.html`. Its head counts what was added, changed and withdrawn; check they match what the change meant to do.
2. Follow each entry to the page it lands on, and look at the mark there. An entry with no link is one whose page is gone; the list that held it says so by its count.
3. Take a screenshot of the diff page and of the page that matters most to the change, so the maintainer sees it without opening anything.

The site's language follows `lang:` in `livt.yaml`, so this repository's preview is in Japanese. Talk to the maintainer in the language they are using.

## Report

Say what the diff shows, by livt URI and what happened to each, and point out anything the change did not intend: a URI that moved without being edited, a rule withdrawn rather than changed, a page with no mark. An empty diff after an edit is a finding too — the edit is not something livt reads.

Where the rendering disagrees with the example mapping that decided it, that is a gap to report, not something to work around here. The output directory is outside the repository; nothing from it is committed.
