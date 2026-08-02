# Example Mappings

Example mappings are YAML files stored in `discoveries/example-mappings/`. They capture the rules, examples, and questions discovered during an Example Mapping session for a story.

## Format

```yaml
story: story-key

rules:
  - id: R-01
    name: Rule description
    examples:
      - id: EX-01
        name: Example description
    issues:
      - https://github.com/owner/repo/issues/1
    automated: true
  - id: R-02
    name: Rule the spec no longer asks for
    retired: true

questions:
  - id: Q-01
    text: Question text

ubiquitous:
  - term-key
```

- `story` is optional (links to the corresponding story detail page)
- IDs must be unique within their rule or question list
- `ubiquitous` is optional: each entry is a [ubiquitous language](./ubiquitous-language.md) term key, rendered as a pink sticky linking to `ubiquitous.html#{term-key}`. A key with no matching term file renders as a plain pink card.
- `issues` is optional: the rule's automation Issue URLs on implementation repos (Issue URLs only). The livt repository records the links; their state lives at the URL target. A rule without `issues` is unlinked.
- `automated` is optional: records the judgment that the rule is actually automated by tests, which is independent of Issues being filed or closed. Absent means not automated. Set it when the rule's automation lands; unset it when the rule changes.
- `retired` is optional and applies to a rule, an example, or a question: it records that the item is no longer part of the spec. Absent means live.

## Retiring an item

An item that no longer holds is marked `retired: true` — never deleted, and never commented out:

- **Deleting frees the ID.** With `R-01`/`R-02`/`R-03` on file, deleting `R-03` makes `R-02` the highest, so the next rule takes `R-03` back. A `livt://mapping/{story-key}/rule/R-03` reference already quoted in an Issue or a test comment then resolves to a *different* rule instead of failing — the quietest way for a reference to break. Retired items keep their IDs taken: new IDs are numbered from the max including them.
- **Commenting out loses the record.** A comment is not part of the YAML structure, so any tool that rewrites the file drops it. `retired: true` is a field and survives.

A retired item stays readable in the file and still resolves by its livt URI, carrying `retired: true` so the reader can tell. It leaves the board and the [Tasks page](../reference/file-structure.md): a retired question is not an open question, and a retired rule is not waiting for a test.

## Visual Layout

The board renders cards in the [Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/) format:

- **Yellow** card: Story (top)
- **Blue** cards: Rules (row below story)
- **Green** cards: Examples (stacked under their rule)
- **Red** cards: Questions (separate column)
- **Pink** cards: Ubiquitous language terms (referenced via `ubiquitous`, below the board)

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
