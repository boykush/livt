---
name: render-site
description: Render livt's docs site — the HTML in docs/site — from its text in docs/content, after the text changes or when the board excerpt on the top page should catch up with the record. Keeps the design the site already has and changes only what the text or the record changed. Use whenever docs/content/ is edited, or when asked to update, re-render, or fix the docs site.
---

You **render the docs site**. Its text lives in `docs/content/{en,ja}/*.md`, one file per page and per language; the HTML in `docs/site/` is made from it and committed, and Pages publishes that directory as it is. The text is what a change edits. You carry it into the HTML, so what the site says is always what the text says.

## Language

The site's text follows `AGENTS.md`: English under `docs/content/en/`, Japanese under `docs/content/ja/`, each written in its own language rather than translated from the other. If only one language's text changed, say so and ask whether the other is owed the same change, written by the maintainer or drafted by you for them to review. Never fill it with a translation on your own. Talk to the maintainer in the language they are using.

## What Maps Where

- `docs/content/en/{page}.md` becomes `docs/site/{page}.html`; `docs/content/ja/{page}.md` becomes `docs/site/ja/{page}.html`.
- `site.md` is the text every page of that language shows around its own: the header, the language switch, the footer. It has no page of its own.
- Front matter `title:` is the page's `<title>`; `description:`, where given, is its description meta.
- Every heading, paragraph and list item is shown on the page, word for word. Where it sits and how it is split across elements is the design's call: a canvas box, a sticky, a label beside its answer.
- A comment (`<!-- … -->`) is a note to you, and is not shown. `<!-- board: … -->` names the board excerpt the top page draws.

## Keep the Design

Change the HTML only where the text or the record changed. Don't restyle a page, reorder its sections, or rewrite `docs/site/style.css` unless the maintainer asks for a design change. A page that is new, or a structure that is new (a section the text adds), follows the patterns the other pages already use, with the same header, footer, stylesheet and language switch.

Pages link each other relatively, since the site lives under `/livt/`, and each links the same page in the other language.

## The Board Excerpt

The top page draws part of livt's own board, the one its `<!-- board: … -->` comment names, the way the demo draws it. It quotes the record, not the text:

- Read the stickies from `discoveries/example-mappings/{story-key}.yaml` and the story's name from `stories/{story-key}.md`. The Japanese page copies each one verbatim. The English page translates them, and its caption says so.
- A sticky's ✓ and its test link come from the automation reports in `automations/`: the citation's file and line, and the URL the report gives. When the demo's rendering is in doubt, build it (`go run . build -o <dir>`) and match the mapping's page.
- A question longer than a sticky ends with `（…）` or `(…)` rather than being reworded.

If a sticky the comment names has been retired or rewritten, draw what the record says now, and tell the maintainer, who may want a different excerpt.

## Redirects

`docs/site` also holds redirect pages for addresses the mdBook site served, which published release notes link: `installation.html`, `getting-started.html`, `introduction.html`, and the pages under `guides/` and `reference/`. Leave them alone. A page you rename or remove needs one of its own, in the same form.

## Check

`go test ./docs/` fails when a page does not show a passage of its text, does not link what the text links, or carries another title, and when a page has no text behind it. Run it, then look at the pages: serve `docs/site` (`python3 -m http.server --directory docs/site`, or the `docs-site` preview) and open each changed page in both languages, at desktop and phone width.

## Commit

The text and the HTML made from it go in one `docs:` commit, so neither is ever ahead of the other. The body says why the text changed; the diff under `docs/site/` is the rendering of it.
