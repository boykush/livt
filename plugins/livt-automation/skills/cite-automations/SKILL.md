---
name: cite-automations
description: Place the livt:automates citations in an implementation repository's tests — which test cites a rule, which cites an example, and where the line sits — so the livt repository's board reads what the tests cover. Use while writing or changing a test that automates a rule or an example, or to sweep a repository's tests for rules that have become citable; in nested tests every level carries its citation, and in flat tests a rule is cited once all of its examples are. The marker and what livt collects are livt's; this skill holds the placement convention livt ships as a default.
---

You **cite automations** — the claim a test makes, at the automate station of the delivery ring.

A test that automates a point of the spec says so on a comment line: the marker `livt:automates` and one livt URI. livt collects those lines and nothing else, and the livt repository's board reads a rule or an example as automated when some test cites it. Your job is to put the citations where they belong, so what the board says is what the tests do.

## Language

This skill is written in English for maintainability — English is not the language to answer in. Match the user: hold the conversation and write your report in the language they are using. Test code, comments, and livt URIs stay as the repository already writes them.

## livt's and Yours

- **livt's** — what a citation is. The marker and one livt URI on a comment line, in your language's comment syntax; a URI without the marker is a reference, not a claim; one URI per line, so a test citing several points carries several lines. The server's instructions over MCP carry this; don't restate a different spelling.
- **This skill's** — where the lines go. livt's scan reads lines and URIs and never the structure around them, so placement is a convention of the side that writes the tests. This is the one livt ships as the default. Where a repository's `AGENTS.md` or its own skill says otherwise, that wins.

Read the spec over the livt MCP server before citing: the rule's resource lists its examples, each with its own `uri`, and marks retired ones. Quote the `uri` the server returns; never assemble one from a bare `R-02`, which exists in every mapping.

## Where Citations Go

The snippets below are in Go only for illustration: the convention is the same in any language and framework, with the line written in that language's comment syntax. They write livt URIs with `{placeholders}`, which livt never collects; in a real test each is the `uri` the server returned. Rules and examples are claimed separately. Citing a rule says a test checks the rule itself; citing an example says a test checks that example. Neither stands for the other on livt's side: the scan and the board never count a rule as automated because its examples are.

### Nested tests — every level carries its citation

When the framework lets a block wrap a rule's cases (`describe`/`it`, `t.Run`, nested classes, `context` blocks, …), cite the rule on the block and each example on its case:

```go
// livt:automates livt://mapping/{story_key}/rule/{rule_id}
func TestLateOrdersAreRefused(t *testing.T) {
	// livt:automates livt://mapping/{story_key}/rule/{rule_id}/example/{example_id}
	t.Run("an order after the deadline is refused", func(t *testing.T) { … })
	// livt:automates livt://mapping/{story_key}/rule/{rule_id}/example/{example_id}
	t.Run("an order at the deadline is accepted", func(t *testing.T) { … })
}
```

### Flat tests — the rule once all its examples are cited

When each test stands alone with no block around a rule's cases, there is nowhere for the rule's own citation to sit. Each test cites the examples it checks, and the rule waits:

1. Read the rule over MCP and list its **live** examples — retired ones do not count, and a rule whose `status` is `proposed` is not spec yet, so it is not cited at all.
2. Find which of those examples **this repository's** tests cite, by their marker lines. Citations reported from other implementation repositories are theirs; a rule line here claims that the tests here cover it.
3. When every live example is cited here, add the rule's citation to **each** test that cites one of them — one more line beside the example's:

```go
// livt:automates livt://mapping/{story_key}/rule/{rule_id}
// livt:automates livt://mapping/{story_key}/rule/{rule_id}/example/{example_id}
func TestAnOrderAfterTheDeadlineIsRefused(t *testing.T) { … }

// livt:automates livt://mapping/{story_key}/rule/{rule_id}
// livt:automates livt://mapping/{story_key}/rule/{rule_id}/example/{example_id}
func TestAnOrderAtTheDeadlineIsAccepted(t *testing.T) { … }
```

Every test that makes up the rule then leads back to it, and deleting one of them does not take the rule's citation with it.

Until then, leave the rule uncited — a rule with one example uncovered is not automated, and the board should say so.

### A rule with no examples

A test that checks it cites it directly, nested or flat.

## The Inference Stays Here

Counting a rule as automated because its examples all are is a judgment, and it is yours: you make it and write it down as a citation, where a reviewer sees it in the diff. livt never makes it for you. Don't ask for the scan to infer it, and don't count a rule automated on the board's behalf.

When a rule you once cited this way has gained an example no test here cites, leave the rule's citation where it is: taking it off would make you chase every later change to the mapping, and you only see this repository. The gap stays visible anyway, because the new example's sticky on the livt repository's board carries no citation of its own. Cite the new example when a test covers it, and mention the gap in your report.

## What NOT to Do

- Don't cite a rule in a flat test before every live example of it is cited in this repository.
- Don't cite a proposed rule, or a retired rule or example.
- Don't write a bare ID or a deployed URL where the livt URI goes, and don't put two URIs on one line.
- Don't edit the livt repository or its `automations/` reports: the report is collected from your tests, never written by hand.

## Output

Test files whose marker lines follow the convention above, in the working tree for the repository's normal review. Report the rules you newly cited, and any cited rule with an example no test here covers yet.
