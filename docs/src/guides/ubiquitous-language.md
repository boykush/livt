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
`id={term-key}` anchor, so a term is linkable as `ubiquitous.html#{term-key}` —
for example, to reference it from a story map or example mapping board.
