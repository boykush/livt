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
name: ストーリーマップ
---

アクティビティ・ステップ・ストーリーをリリーススライスと共に俯瞰するボード。
```

## Visual Layout

`livt build` renders every term as a row on a single page at `ubiquitous.html`,
with **Term**, **Key**, and **Definition** columns. Each row carries an
`id={term-key}` anchor, so a term is linkable as `ubiquitous.html#{term-key}`.

## Referencing terms from boards

Story maps and example mappings can declare the terms they use with a top-level
`ubiquitous` list of term keys:

```yaml
ubiquitous:
  - story-map
  - story
```

Referenced terms render as pink stickies below the board, each linking to its
glossary row (`ubiquitous.html#{term-key}`). A key with no matching
`ubiquitous/{term-key}.md` file renders as a plain pink card, so references
degrade gracefully. See [Story Maps](./story-maps.md) and
[Example Mappings](./example-mappings.md).
