---
title: livt - Living Text
description: BDD for teams whose code is written by coding agents. livt keeps what your team agreed on the board in Example Mapping's own format, owned by every role on the team.
---

# Collaborate on board. Make it living in text.

Teams decide what to build on a whiteboard, often an online one such as Miro or FigJam. livt keeps what they agreed in the repository, in Example Mapping's own format, so that every role on the team owns it. Coding agents read it and carry it through the flow BDD has always aimed at: discovery, formulation and automation.

[Open the live demo](https://boykush.github.io/livt/demo/)

[Get started](https://github.com/boykush/livt#getting-started)

<!-- board: livt://mapping/collect-automations — story; rule R-01 with examples EX-01 and EX-03; rule R-02 with examples EX-01 and EX-02; question Q-01; the test line that cites R-01. The stickies are quoted from the record, translated on this page. -->

An excerpt of livt's own example mapping, translated from Japanese. The test below marks R-01 as the rule it automates. livt collects such marks into a report, and each sticky's ✓ comes from it.

[See the whole board](https://boykush.github.io/livt/demo/mapping/collect-automations.html)

## BDD, now that coding agents write the code

The BDD community has long said that discovery matters more than formulation, and formulation more than automation. It also keeps what a team agreed as one consistent specification.

An Example Mapping does not end with the conversation, though. Its rules change and its questions get answered all through development, and a formulated feature file alone cannot follow that. livt records the Example Mapping itself, in its own format: it holds what a feature file would, and lets every role follow questions being answered and rules changing, designers and product managers as much as developers.

Coding agents have changed the situation. Automating an agreed rule is now work an agent can take on. Following the test pyramid, an Example Mapping's rules and examples end up automated at different layers and in different languages. livt gives the agent a shape to follow and the tools to follow it, as a CLI and plugins, and works through what those differences bring with it. Teams that want their specification executed as written have Cucumber for that.

<!-- The boxes below are livt's own opportunity canvas, for livt as a whole. Their headings are the labels the demo gives a canvas's boxes. -->

### Problems

- What becomes of an Example Mapping afterwards, its questions resolved and its rules changed, is more than a formulated feature file alone can follow.
- What a team agreed is hard to share across roles.
- Automation spreads across layers and repositories, and it is hard to tell which rule is automated where.

### Customers & Users

- Every role on a team that practises BDD: product managers, designers, developers, testers
- Engineers who leave implementation to coding agents

### Solutions Today

- Formulate rules and examples as Gherkin, and automate them with Cucumber as an executable specification
- Copy the agreement into tickets or documents, and stop syncing from there
- Keep the board as it was, and look back at it later

### Solution Idea

- Keep recording the Example Mapping itself, as something every role owns
- Leave automation to coding agents, and give them a shape and tools as a CLI and plugins
- Tests cite the rules they automate, and livt gathers which rule is automated where

### What Will Users Do To Get Value?

- Hand an agreed rule to a coding agent, and have it automated
- Follow questions being answered and rules changing in one place, whatever their role
- See which rules are not automated yet, and pick the next one to close

<!-- Link the opportunities rather than list them: a list here falls behind the day one is added. -->

## Each problem, recorded as an opportunity

These problems are recorded one by one as opportunities in livt's own repository. Each has its own canvas and story map, which the live demo shows.

[See the opportunities](https://boykush.github.io/livt/demo/opportunities.html)

## Get started

Install livt and add its plugins to your coding agent. The README has the steps.

[Read Getting started](https://github.com/boykush/livt#getting-started)

[Releases](https://github.com/boykush/livt/releases)
