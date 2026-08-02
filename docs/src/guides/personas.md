# Personas

Personas are Markdown files with YAML frontmatter, stored in the `personas/`
directory. Each file is one actor the stories are written for, and livt renders
them as a table you can browse alongside the stories that name them.

## Format

```markdown
---
name: Persona display name
---

Description in Markdown.
```

The `name` field in frontmatter is the display name. The **persona key** is
derived from the filename (without `.md`) and must be kebab-case, using lowercase
letters, numbers, and hyphens. The body describes who the actor is.

## Example

`personas/coding-agent.md`:

```markdown
---
name: Coding agent
---

An AI agent that reads the livt repository over MCP and automates rules
test-first in an implementation repository.
```

## Naming a persona from a story

A story declares its actor with a `persona` key in its frontmatter:

```markdown
---
name: Automate from the spec, test-first
persona: coding-agent
---

As a coding agent
I want to automate rules from the example mapping test-first
So that the spec keeps following the implementation
```

`persona` is reserved the way `name` is: it addresses a persona rather than
becoming a [metadata](./stories.md) row. A story names at most one — a user
story has one "as a" line — and naming none is fine, so a story written before
its actor was committed still parses.

The story's body keeps its "as a" line. The line is prose, which is exactly why
the key exists: prose is what drifts into a second name for one actor, and the
key is what says the two lines mean the same person.

## Why not a context

A ubiquitous language term takes an optional
[context](./ubiquitous-language.md#contexts) because a word can mean different
things to different parts of the domain. A persona takes none. An actor named
twice is the same actor, and cutting personas by directory would invite the
drift the list exists to stop, so `personas/` is flat and a persona key is one
segment.

## Visual Layout

`livt build` renders every persona as a row on a single page at `personas.html`,
with **Persona**, **Key**, **Description**, and **Stories** columns. The Stories
column is the reverse of the frontmatter link — every story that named this
persona — so the page answers "whose stories are these" rather than reading as a
second glossary. A persona no story names yet keeps its row and says so.

Each row carries an `id={persona-key}` anchor, so a persona is linkable as
`personas.html#{persona-key}`. A story's page and the Stories list both show the
actor as a chip linking there. A `persona` key with no matching file renders as a
plain chip, so the link degrades gracefully.

## Reading personas from tooling

The MCP server lists the cast with `list_personas` and serves each one at
`livt://persona/{persona-key}`; a story resource carries its own `persona`
resolved to that URI. See [Commands](../reference/commands.md#resources) and
[livt URI](../reference/uri.md).
