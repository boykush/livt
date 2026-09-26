---
title: livt - Living Text
description: livt keeps what your team agreed on its boards, each in its own format, owned by every role on the team. Coding agents read it and carry it through to automation.
---

# Collaborate on board. Make it living in text.

Teams decide what to build on a whiteboard, often an online one such as Miro or FigJam. livt keeps what they agreed in the repository, each board in its own format, so that every role on the team owns it. Coding agents read it and follow it as one thread, from the opportunity to the tests.

[Open the live demo](https://boykush.github.io/livt/demo/)

[Get started](https://github.com/boykush/livt#getting-started)

<!-- board: livt://mapping/collect-automations — story; rule R-01 with examples EX-01 and EX-03; rule R-02 with examples EX-01 and EX-02; question Q-01; the test line that cites R-01. The stickies are quoted from the record, translated on this page. -->

An excerpt of livt's own example mapping, translated from Japanese. The test below marks R-01 as the rule it automates. livt collects such marks into a report, and each sticky's ✓ comes from it.

[See the whole board](https://boykush.github.io/livt/demo/mapping/collect-automations.html)

## What livt keeps

<!-- The glossary owns what each one is, so this says only its part in the flow and links the name to the glossary's entry. -->

- [Opportunity Canvas](https://boykush.github.io/livt/demo/ubiquitous.html#opportunity-canvas) — why to build
- [Story Map](https://boykush.github.io/livt/demo/ubiquitous.html#story-map) — what to build
- [Example Mapping](https://boykush.github.io/livt/demo/ubiquitous.html#example-mapping) — how it should behave
- [Ubiquitous Language](https://boykush.github.io/livt/demo/ubiquitous.html#ubiquitous-language) — which words to use

The people behind each practice are credited on [livt and BDD](livt-and-bdd.html).

## Why livt

<!-- The boxes below are livt's own opportunity canvas, for livt as a whole. Their headings are the labels the demo gives a canvas's boxes. -->

### Problems

- Decisions made in discovery keep changing during development, yet stay on the board and go stale.
- Opportunities, stories and rules sit on separate boards, with no way to follow one to the next.
- What was decided goes unread from the implementation side, and how much of it is implemented and automated is hard to tell.

### Customers & Users

- Every role on a team that runs discovery through development together: product managers, designers, developers, testers
- Engineers who leave implementation to coding agents

### Solutions Today

- Keep the board as it was, share its URL, and look back at it later
- Copy the decisions into tickets or documents, and stop syncing from there
- Have people mark in tickets or flags whether something is implemented or automated

### Solution Idea

- Record each practice's board in the repository in its own format, and layer changes on it as diffs
- Link opportunities to their story maps and Example Mappings with URIs
- Leave automation to coding agents, gather the rules their tests cite, and show where each opportunity stands

### What Will Users Do To Get Value?

- Have a coding agent record what the board decided, and follow changes in the same place
- Follow an opportunity down to its Example Mappings on one site
- See where each opportunity stands, and pick the next place to work on

## Get started

Install livt and add its plugins to your coding agent. The README has the steps.

[Read Getting started](https://github.com/boykush/livt#getting-started)

[Releases](https://github.com/boykush/livt/releases)
