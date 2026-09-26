---
title: livt - Living Text
description: BDD for teams whose code is written by coding agents. livt keeps what your team agreed on the board as text in your repository.
---

# Collaborate on board. Make it living in text.

Teams decide what to build on a whiteboard, often an online one such as Miro or FigJam. livt keeps what they agreed as text in your repository. Coding agents read that text and carry it through the flow BDD has always aimed at: discovery, formulation and automation.

[Open the live demo](https://boykush.github.io/livt/demo/)

[Get started](https://github.com/boykush/livt#getting-started)

<!-- board: livt://mapping/collect-automations — story; rule R-01 with examples EX-01 and EX-03; rule R-02 with examples EX-01 and EX-02; question Q-01; the test line that cites R-01. The stickies are quoted from the record, translated on this page. -->

An excerpt of livt's own example mapping, translated from Japanese. The test below marks R-01 as the rule it automates. livt collects such marks into a report, and each sticky's ✓ comes from it.

[See the whole board](https://boykush.github.io/livt/demo/mapping/collect-automations.html)

## BDD, now that coding agents write the code

BDD has always aimed at one flow: find the rules together, formulate them, and automate them as tests. The step between formulation and automation used to be interpretation written by hand, and a coding agent can now do it. livt gives the agent a shape to follow through that flow, as a CLI and plugins for coding agents.

<!-- The boxes below are livt's own opportunity canvas, for livt as a whole. Their headings are the labels the demo gives a canvas's boxes. -->

### Problems

- Someone writes and maintains the layer that interprets Gherkin into code, such as Cucumber's step definitions.
- A team that follows the test pyramid cannot turn every scenario into an end-to-end (large) test.
- So automation spreads across the implementation repositories, such as a backend and a frontend, and it is hard to tell which rule is automated where.

### Customers & Users

- Development teams that want to practise BDD
- Engineers who leave implementation to coding agents

### Solutions Today

- Bind Gherkin to step definitions with Cucumber and run it as tests
- Copy the agreement into tickets, and stop syncing from there

### Solution Idea

- A coding agent does the interpreting that step definitions used to do
- Every rule and example has a URI, and a test cites the URI of what it automates
- Tests can sit at any layer in any repository; livt collects their citations mechanically and shows each rule's status
- The flow from discovery through formulation to automation ships as a CLI and plugins for coding agents

### What Will Users Do To Get Value?

- Hand an agreed rule to a coding agent, and have it automated as a test at the layer that fits
- See which rules are not automated yet, and pick the next one to close

<!-- Link the opportunities rather than list them: a list here falls behind the day one is added. -->

## Each problem, recorded as an opportunity

These problems are recorded one by one as opportunities in livt's own repository. Each has its own canvas and story map, which the live demo shows.

[See the opportunities](https://boykush.github.io/livt/demo/opportunities.html)

## Get started

Install livt and add its plugins to your coding agent. The README has the steps.

[Read Getting started](https://github.com/boykush/livt#getting-started)

[Releases](https://github.com/boykush/livt/releases)
