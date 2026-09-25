# Example Mappings

Example mappings are YAML files stored in `discoveries/example-mappings/`. They capture the rules, examples, and questions discovered during an Example Mapping session for a story. The format is [Matt Wynne's](https://cucumber.io/blog/bdd/example-mapping-introduction/); what livt adds to it, and where it departs, is [livt and BDD](../livt-and-bdd.md).

## Format

```yaml
name: Mapping name

rules:
  - id: R-01
    name: Rule description
    examples:
      - id: EX-01
        name: Example description
    issues:
      - https://github.com/owner/repo/issues/1
  - id: R-02
    name: Rule the spec no longer asks for
    status: retired
    superseded_by:
      - livt://mapping/story-key/rule/R-03
  - id: R-03
    name: Rule put forward, not yet agreed
    status: proposed

questions:
  - id: Q-01
    text: Question text

ubiquitous:
  - term-key
```

- The filename is the story key: `{story-key}.yaml` maps the story in `stories/{story-key}.md`, and the board links to that story's page when there is one
- `name` is optional: the mapping's own name. The board, its overview tile, the Tasks page, and the diff call the mapping by it, falling back to the story's name and then to the key. livt reads it apart from the story's `name` and never compares the two; keeping them the same where both are written is left to whoever writes them, which the [livt-discovery skills](https://github.com/boykush/livt/tree/main/plugins/livt-discovery) do. See [A mapping without a story](#a-mapping-without-a-story)
- IDs must be unique within their rule or question list
- `ubiquitous` is optional: each entry is a [ubiquitous language](./ubiquitous-language.md) term key, rendered as a pink sticky linking to `ubiquitous.html#{term-key}`. A key with no matching term file renders as a plain pink card.
- `issues` is optional: the rule's automation Issue URLs on implementation repos (Issue URLs only). The livt repository records the links; their state lives at the URL target. A rule without `issues` is unlinked.
- Whether a rule is automated is not a field. The tests that automate it say so, and livt reads it from them — see [Automating a rule](#automating-a-rule). A deprecated `automated:` is still read where a mapping carries one.
- `status` is optional and applies to a rule: `proposed` while the rule is put forward but not yet agreed, `accepted` once it is, `rejected` when the proposal was turned down, `retired` when spec it once was stopped holding. Absent means accepted, and any other value fails the build. See [Proposing a rule](#proposing-a-rule).
- `retired` is optional and applies to an example or a question: it records that the item is no longer part of the spec. Absent means live. It is not a rule field — a rule closes through `status`; see [Retiring an item](#retiring-an-item).
- `superseded_by` is optional and goes with a closed rule or a retired example or question: the [livt URIs](../reference/uri.md) of whatever took its place. Absent means nothing did.

## A mapping without a story

Not every mapping starts from a story card. A bug fix whose expected behaviour is already clear can go straight to the rule it broke and an example that reproduces it, and from there into a failing test. Such a mapping is a file with no `stories/{story-key}.md` beside it:

```yaml
name: Login rejects a password with a full-width space

rules:
  - id: R-01
    name: A password is compared exactly as it was typed
    examples:
      - id: EX-01
        name: A password starting with a full-width space logs in
```

Give it a `name`: without one, the board is called by its key. Otherwise it is a mapping like any other — it renders, lists on the Tasks page, diffs, and resolves by livt URI. What it lacks is a story page to link to and, since a mapping reaches its opportunities through its story, any opportunity to be filtered under. An agent finds it through the MCP server's `list_example_mappings`, which lists every mapping, story or not.

## Proposing a rule

A rule can go on the board before it is agreed. `status: proposed` marks it as put forward and awaiting agreement, and from there it moves on the way an [architecture decision record](https://adr.github.io/) does:

| ADR status | In an example mapping |
|------------|-----------------------|
| Proposed | `status: proposed` |
| Accepted | `status: accepted`, or no `status` at all |
| Rejected | `status: rejected` |
| Deprecated | `status: retired` |
| Superseded | `status: retired` with [`superseded_by`](#saying-where-the-spec-went) |

- **A proposal is a candidate answer.** A Question card records what the team does not know yet; a proposed rule puts one answer on the table, examples and all, for the team to agree to or turn down.
- **Accepting it is a one-line diff.** `proposed` becomes `accepted`, and that line is what the review approves.
- **Rejecting it is a one-line diff too.** `proposed` becomes `rejected`. Its ID stays taken like any closed rule's, and the record reads as a proposal that never became spec rather than as a rule dropped later — the distinction a bare retirement flag could not draw.
- **Absent means accepted.** Every rule written before the field existed was agreed when it went on the board, so an existing mapping reads as it always has.

A proposed rule stays on the board, drawn pale with a dashed edge and stamped *proposed*, and its examples are drawn pale with it. The [Tasks page](../reference/file-structure.md) lists it under Proposed Rules — it closes by agreement, not by a test — and never under Un-automated Rules, even once a test covers it. The MCP server and `livt resolve` return `status` on every rule, so an agent choosing what to automate takes the `accepted` ones and passes over the rest without knowing the default.

## Automating a rule

A mapping does not record whether a rule is automated. The tests do: a test in an implementation repository that automates a rule says so on a comment line above itself, with the `livt:automates` marker and the rule's [livt URI](../reference/uri.md).

```go
// livt:automates livt://mapping/{story-key}/rule/{rule-id}
func TestAnExpiredCardIsRejected(t *testing.T) {
```

The comment syntax is whatever the test's language uses; livt looks only for the marker and the URI. [`livt automations`](../reference/commands.md#livt-automations) collects every marked line in the implementation repository into a report, which goes back to the livt repository as a pull request and lands under `automations/`. livt reads the reports there, and a rule is automated when one of them cites it.

- **The test is the record.** Which tests cover a rule is the implementation's state, not a decision, so the mapping does not keep it. A citation is written beside the test that makes it true and read again from that test each time the report is collected, so the mapping has nothing to keep in step by hand.
- **A rule and its examples are cited separately.** An example is cited by its own URI, `livt://mapping/{story-key}/rule/{rule-id}/example/{example-id}`, and that automates the example alone: a rule whose examples are all cited is still not automated until a test cites the rule itself. The test's author already said which one they meant, and livt does not infer the other.
- **One marker, one URI, one line.** A test that automates several points writes a line for each. A line with anything after its URI is not collected: `livt automations` names it and exits non-zero rather than dropping it in silence.
- **Without the marker, a URI is a reference.** Code quotes the spec for context too, so a livt URI in a comment counts as automation only when the marker comes before it.
- **A citation names the rule, not its wording.** Reword a rule and the tests citing it go on counting, so it keeps reading as automated until they catch up. Nothing unsets that for you: say so in the commit that changes the rule, where the review will see it.
- **Filing is not automating.** [`issues`](#format) link the work to automate a rule; filing one, or closing it, changes nothing here.

On the board, a rule or an example that a test cites is stamped ✓ *automated* on its own sticky, and lists its citations: the repository, file and line of each, linked to that line on the forge when livt could build the link. The [Tasks page](../reference/file-structure.md) lists every accepted rule that no report cites under Un-automated Rules. Neither says whether the tests pass — livt runs no tests, and shows what a test claims without deciding whether the claim is true.

Setting this up in an implementation repository — the CI step that collects the report and opens its pull request — is covered by the [livt-automation plugin](https://github.com/boykush/livt/blob/main/plugins/livt-automation/README.md#sending-automations-back).

`automated: true` on a rule once recorded all of this as a judgment written into the mapping. It is deprecated rather than gone: a rule still carrying the line reads as automated, so a repository keeps the board it had while its rules gain their citations. Nothing maintains the flag any more — no skill writes it and no command moves it — so it stays where a citation would have been withdrawn. Delete the line as each rule is cited from the test that automates it. livt will stop reading it once the status derived from citations can carry a board on its own; no release is promised for that, and the fallback stays until it can.

## Retiring an item

A rule that no longer holds takes a closed `status` — `rejected` or `retired`; an example or a question that no longer holds is marked `retired: true`. Neither is deleted, and neither is commented out:

- **Deleting frees the ID.** With `R-01`/`R-02`/`R-03` on file, deleting `R-03` makes `R-02` the highest, so the next rule takes `R-03` back. A `livt://mapping/{story-key}/rule/R-03` reference already quoted in an Issue or a test comment then resolves to a *different* rule instead of failing — the quietest way for a reference to break. Retired items keep their IDs taken: new IDs are numbered from the max including them.
- **Commenting out loses the record.** A comment is not part of the YAML structure, so any tool that rewrites the file drops it. Both spellings are fields and survive.

A closed item stays readable in the file and still resolves by its livt URI, saying so the way it does in the YAML — a rule through its `status`, an example or a question through `retired: true`. It leaves the board and the [Tasks page](../reference/file-structure.md): a retired question is not an open question, and a retired rule is not waiting for a test.

`retired: true` on a **rule** was once the spelling for both closed statuses, folded onto `status` as the mapping was read. It is not read any more: a rule still carrying that line is a rule with no `status`, which means accepted — back on the board, and back on the Tasks page. Write `status: retired` instead, or `status: rejected` where the rule was a proposal that was turned down.

### Saying where the spec went

Retirement alone tells a reference it has stopped, not where to go. `superseded_by` adds that half — the livt URIs of whatever took the item's place:

```yaml
rules:
  - id: R-02
    name: Rule the spec no longer asks for
    status: retired
    superseded_by:
      - livt://mapping/story-key/rule/R-05
      - livt://mapping/another-story/rule/R-01

questions:
  - id: Q-01
    text: Question that turned into a rule
    retired: true
    superseded_by:
      - livt://mapping/story-key/rule/R-05
```

- **It is a list**, so a rule that split into two names both.
- **It holds livt URIs, not bare ids.** A successor can live in another mapping — the rule moved to the story that actually owns it — and `R-05` on its own names nothing, since ids restart in every mapping.
- **A settled question points at the rule that settled it.** The answer lands as a rule; the Question card never carries one.
- **Nothing replaced it?** Leave `superseded_by` off. Plenty of retirements are just the business no longer asking.

Only the pointer is structured. *Why* the item was retired belongs to the commit that retired it, where it is written once and cannot drift — a second copy in the YAML would. Tooling reads the pointer back as URIs and stops there: the successor is one read away for whoever needs it, and inlining its text would spend a consumer's context on a hop most of them never take.

## Visual Layout

The board renders cards in the [Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/) format:

- **Yellow** card: Story (top), reading the mapping's `name` when it has one and the story's otherwise; it links to the story page when there is one
- **Blue** cards: Rules (row below story)
- **Pale blue, dashed** cards: Proposed rules, stamped *proposed*; their examples are drawn pale too
- **Green** cards: Examples (stacked under their rule)
- **Red** cards: Questions (separate column)
- **Pink** cards: Ubiquitous language terms (referenced via `ubiquitous`, below the board)

### Reading it back as a list

The board keeps the layout the session left on the wall. To read a mapping back later, switch it to **List** with the toggle above it:

- Each rule is one row: its ID, name, automated mark and citations, and issue links, with the number of examples folded under it
- Examples fold under their rule. Click a row to open it
- Questions follow the rules, one row each

The choice is remembered in the browser, so the next mapping you open is shown the same way; the previews on the overview page stay boards. A link to a sticky lands on it in either view — a link to a folded example opens its rule first.

## Example

`discoveries/example-mappings/confirm-discovery-outcomes.yaml`:

```yaml
rules:
  - id: R-01
    name: An example mapping can be rendered as a sticky view with only a story reference
    examples:
      - id: EX-01
        name: A YAML with only a story reference displays a single yellow Story card

  - id: R-02
    name: Cards are laid out following the Example Mapping format
    examples:
      - id: EX-01
        name: Rules are displayed as blue cards in a row below the Story card
      - id: EX-02
        name: Examples are displayed as green cards stacked under their Rule
      - id: EX-03
        name: Questions are displayed as red cards in a separate column

questions: []
```

![Example mapping board](../images/example-mapping.png)
