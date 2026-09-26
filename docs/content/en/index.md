---
title: livt - Living Text
description: BDD for teams whose code is written by coding agents. livt keeps what your team agreed on its boards, each in its own format, owned by every role on the team.
---

# Collaborate on board. Make it living in text.

Teams decide what to build on a whiteboard, often an online one such as Miro or FigJam. livt keeps what they agreed in the repository, each board in its own format, so that every role on the team owns it. Coding agents read it and carry it through the flow BDD has always aimed at: discovery, formulation and automation.

[Open the live demo](https://boykush.github.io/livt/demo/)

[Get started](https://github.com/boykush/livt#getting-started)

<!-- board: livt://mapping/collect-automations — story; rule R-01 with examples EX-01 and EX-03; rule R-02 with examples EX-01 and EX-02; question Q-01; the test line that cites R-01. The stickies are quoted from the record, translated on this page. -->

An excerpt of livt's own example mapping, translated from Japanese. The test below marks R-01 as the rule it automates. livt collects such marks into a report, and each sticky's ✓ comes from it.

[See the whole board](https://boykush.github.io/livt/demo/mapping/collect-automations.html)

## Why livt

<!-- The boxes below are livt's own opportunity canvas, for livt as a whole. Their headings are the labels the demo gives a canvas's boxes. -->

### Problems

- There is no place to carry an Example Mapping with open questions through to completion.
- Rule changes belong to every role, but feature files tend to stay with developers and testers.
- Automation spreads across layers and languages, and it is hard to tell which rule is automated where.

### Customers & Users

- Every role on a team that practises BDD: product managers, designers, developers, testers
- Engineers who leave implementation to coding agents

### Solutions Today

- Formulate rules and examples as Gherkin, and automate them as an executable specification
- Copy the agreement into tickets or documents, and stop syncing from there
- Keep the board as it was, and look back at it later

### Solution Idea

- Record the Example Mapping itself, and give every role one place to follow it until its questions are resolved
- Keep rule changes as diffs of the Example Mapping, readable by every role
- Leave automation to coding agents, and gather the rules their tests cite to show which rule is automated where

### What Will Users Do To Get Value?

- Hand an agreed rule to a coding agent, and have it automated
- Follow questions being answered and rules changing in one place, whatever their role
- See which rules are not automated yet, and pick the next one to close

## Get started

Install livt and add its plugins to your coding agent. The README has the steps.

[Read Getting started](https://github.com/boykush/livt#getting-started)

[Releases](https://github.com/boykush/livt/releases)
