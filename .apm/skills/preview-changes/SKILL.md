---
name: preview-changes
description: Serve the livt site with the diff of the current branch and open it in a browser, so the user sees what the change does to the board — the diff page, and each changed page with its mark. Use right after editing the livt repository (opportunities/, discoveries/, stories/, design/, livt.yaml), and whenever asked to preview, show, or diff a change.
---

You **preview a change** the way a reviewer of it would read it: as livt renders it, not as a YAML diff. `livt serve --diff <base>` builds the site with a diff page at `diff.html`, marks every changed page, and rebuilds and reloads as the files change. What the diff shows is decided in `livt://mapping/review-diff-between-revisions`; this skill only opens it.

## Language

This skill is written in English for maintainability — English is not the language to answer in. Match the user: talk and report in the language they are using. The site itself follows `lang:` in `livt.yaml`, whatever language you speak.

## Start

Reuse a preview that is already running: it rebuilds on every edit to the livt repository.

If the project defines a preview for this — in the Claude desktop app, a `.claude/launch.json` configuration that runs `livt serve --diff` — start that one. A repository that builds livt from its own source binds the command there, so the preview renders the livt being changed.

Otherwise run it yourself, in the background:

```bash
livt serve --diff "$(git merge-base HEAD origin/HEAD)" -p 3002 -o "${TMPDIR:-/tmp}/livt-diff"
```

- The base is where this branch left the default branch, so the diff is everything the branch does, committed or not, and the default branch's newer commits never show up reversed. Without `origin/HEAD`, use the default branch by name.
- The port and output directory keep clear of a plain `livt serve` on 3000 and its `dist/`, and nothing lands in the repository.

Open the URL in the browser you have: the desktop app's browser pane, or the Playwright MCP server. With neither, give the user the URL.

`livt serve` watches the livt repository, not livt itself: when the livt you run changes — upgraded, or rebuilt from source — restart the preview.

## Walk the Change

1. Open `/diff.html`. Its head counts what was added, changed and withdrawn; check they match what the change meant to do.
2. Follow each entry to the page it lands on, and look at the mark there. An entry with no link is one whose page is gone; the list that held it says so by its count.
3. Take a screenshot of the diff page and of the page that matters most to the change, so the user sees it without opening anything.

## Report

Say what the diff shows, by livt URI and what happened to each, and point out anything the change did not intend: a URI that moved without being edited, a rule withdrawn rather than changed, a page with no mark. An empty diff after an edit is a finding too — the edit is not something livt reads.

Where the rendering disagrees with the example mapping that decided it, that is a gap to report, not something to work around here.
