# Ubiquitous Language

Ubiquitous language terms are Markdown files with YAML frontmatter, stored in the
`ubiquitous/` directory. Each file is a single term, and livt renders them as a
table you can browse like a database.

## Format

```markdown
---
name: Term display name
---

Definition in Markdown.
```

The `name` field in frontmatter is the display term. The **term key** is derived
from the filename (without `.md`) and must be kebab-case, using lowercase
letters, numbers, and hyphens. The body is the term's definition.

## Example

`ubiquitous/story-map.md`:

```markdown
---
name: Story Map
---

A board to overview activities, steps, and stories alongside release slices.
```

## Contexts

A term can be scoped to a **context** by putting it in a directory:

```
ubiquitous/
  story-map.md                # holds across contexts
  billing/
    invoice.md                # what "invoice" means in billing
  shipping/
    invoice.md                # what "invoice" means in shipping
```

A context is optional. Cut one when a word means different things to different
parts of the domain — the DDD bounded context — and leave a term at the root
when its meaning holds everywhere. Most masters start with every term at the
root and grow contexts only where the language actually diverges.

The directory is what makes a term unique, so `invoice` above is three separate
terms, not one with three labels. Each has its own definition, its own row, and
its own address. A term reference is therefore `{ctx}/{term-key}` for a scoped
term and `{term-key}` for a context-free one, and the two never resolve to each
other: a bare `invoice` names the root term even when `billing/invoice` exists.

Contexts are one directory deep. `ubiquitous/billing/eu/invoice.md` is not a
term livt can address, and is left out of the glossary.

## Visual Layout

`livt build` renders every term as a row on a single page at `ubiquitous.html`,
with **Term**, **Key**, and **Definition** columns. A scoped term shows its
context beneath its key. Each row carries an `id={term-ref}` anchor, so a term is
linkable as `ubiquitous.html#{term-key}` or `ubiquitous.html#{ctx}/{term-key}`.

When a master cuts at least one context, the page carries a filter bar above the
table. Selecting a context narrows the table to that context's terms, and the
selection is mirrored in the `?context=` query parameter so a filtered view is
shareable. Context-free terms are not part of any single context, so they drop
out while a filter is on rather than being repeated under every chip.

## Referencing terms from boards

Story maps and example mappings can declare the terms they use with a top-level
`ubiquitous` list of term references:

```yaml
ubiquitous:
  - story-map
  - billing/invoice
```

Referenced terms render as pink stickies below the board, each linking to its
glossary row, with the context shown on the sticky when it has one. A reference
with no matching file renders as a plain pink card, so references degrade
gracefully. See [Story Maps](./story-maps.md) and
[Example Mappings](./example-mappings.md).
