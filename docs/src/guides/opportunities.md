# Opportunities

An **opportunity** is something the product could take on: a user problem
together with the business benefit of solving it, held as one unit of
consideration. It is what a story map is *for* — you decide whether to pursue an
opportunity before you map the journey that serves it.

Opportunities are Markdown files with YAML frontmatter, stored in the
`opportunities/` directory. Their canvases are YAML files stored in
`discoveries/opportunity-canvases/`.

## Format

```markdown
---
name: Opportunity display name
---

Whose problem this is, and what the business gets from solving it.
```

The `name` field is a short label — it is what the filter chips and the
navigation show. The **body is the opportunity itself**: the sentence saying
whose problem it is and what solving it is worth. This is the same split a
[story](./stories.md) makes between its name and its narrative, and it is there
for the same reason — a label fits on a card, a statement does not.

Name the **solution idea you are weighing**, the way Patton fills his canvas —
"starting with the solution idea I've been asked for and working backwards"
([Opportunity Canvas](https://jpattonassociates.com/opportunity-canvas/)).
"Collaborative discovery" still fails as a name: it names a topic, and a topic
cannot be declined. "Living fair copy" names an idea, and an idea is something
you can weigh — build it, or decide against it. The problem the idea answers is
not lost; it is the body's opening sentence, and the canvas's left zone holds it
as verifiable fact. The topic already has a home — it is what the
[story map](./story-maps.md) is called, and a map and its opportunity are free to
carry different names precisely because they answer different questions: the map
says which journey was drawn, the opportunity names the idea that made it worth
drawing.

The **opportunity key** is derived from the filename (without `.md`), and must
be kebab-case. Key uniqueness is enforced by the filesystem.

## Example

`opportunities/collaborative-discovery.md`:

```markdown
---
name: Living fair copy
repos:
  - boykush/livt
---

Discovery outcomes sit in the board tool and go stale. The raw record is rough, and
nothing versions it or checks it for consistency. Transcribing it, refining it into a
fair copy, and keeping it updated as rules change saves re-running the same argument
and lets the outcome stand as the specification.
```

Any frontmatter field beyond `name` is kept and shown as metadata, and a field
whose value is a URL renders as a link — the same treatment story frontmatter gets.

## The opportunity canvas

The [Opportunity Canvas](https://jpattonassociates.com/opportunity-canvas/) holds
an opportunity on a single sheet. It is a discovery session's outcome, so it
lives in `discoveries/` beside the story maps and example mappings, and it joins
its opportunity **by filename**:

```
opportunities/collaborative-discovery.md                        # the opportunity
discoveries/opportunity-canvases/collaborative-discovery.yaml   # its canvas
```

This is the same filename join an [example mapping](./example-mappings.md) makes
with its story, and it carries the same meaning: either file can exist without
the other. An opportunity with no canvas has not been thought through yet; a
canvas with no opportunity file still renders.

### Format

```yaml
canvas:
  solution-ideas:
    - A specific product, feature, or enhancement idea
  problems:
    - A problem users have today
  users-and-customers:
    - Who has that problem
  solutions-today:
    - How they address it now
  business-challenges:
    - What those problems cost the business
  user-value:
    - What users will do with the solution
  user-metrics:
    - What you could measure to show they did it
  adoption-strategy:
    - How they discover and adopt it
  business-impact:
    - Which business metrics move
  budget:
    - What you would spend to find out

ubiquitous:
  - term-key
```

Every box is a **list**, because a box on a canvas holds sticky notes rather
than a paragraph. Keeping them apart is what lets the board render one card per
idea, the way it sat in the room.

Every key is optional. A box left out renders as an empty box on the sheet
rather than disappearing from it — a blank box is the visible record of a
question the opportunity has not answered yet, and that is worth seeing.

`ubiquitous` is optional and works exactly as it does on the other boards: each
entry is a [ubiquitous language](./ubiquitous-language.md) term key, rendered as
a pink sticky below the sheet.

### The ten boxes

| YAML key | Box | Zone |
|---|---|---|
| `solution-ideas` | 1. Solution Ideas | Solution |
| `problems` | 2. Problems | Verifiable facts |
| `users-and-customers` | 3. Users and Customers | Verifiable facts |
| `solutions-today` | 4. Solutions Today | Verifiable facts |
| `business-challenges` | 5. Business Challenges | Verifiable facts |
| `user-value` | 6. What Will Users Do To Get Value? | Assumptions about value |
| `user-metrics` | 7. User Metrics | Assumptions about value |
| `adoption-strategy` | 8. Adoption Strategy | Assumptions about value |
| `business-impact` | 9. Business Impact | Assumptions about value |
| `budget` | 10. Budget | Solution |

The numbers are the order Jeff Patton recommends filling the boxes in. The
**zones** are how the sheet is laid out, which is a different thing: the facts
on the left, the solution down the middle, the assumptions about value on the
right. Reading the sheet left to right walks back from the idea to the problem
it solves, then forward to the value it would create.

The split is the point of the canvas — what you can go and check sits apart from
what you are only assuming until the thing ships.

## Linking a story map

A story map serves an opportunity when its **filename key matches**:

```
opportunities/collaborative-discovery.md          # the opportunity
discoveries/usm/collaborative-discovery.yaml      # the journey mapped for it
```

No field connects them; the filename does, as everywhere else in livt. When the
two match:

- The story map board names its opportunity, and links to it
- Every story on that map carries the **opportunity's** name on its chip, and the
  chip links to the opportunity's page
- The Stories, Example Mappings, and Tasks lists filter on that name

A story map whose key names no opportunity file keeps working exactly as it did
before opportunities were files of their own: the map stands in as its own
opportunity, named by the map. Nothing has to be migrated.

## Where opportunities sit

```
Opportunity  ──▶  Story Map  ──▶  Story  ──▶  Example Mapping
   why, and         the journey     what one     how it must
   whether at all   that serves it  person does  behave
```

An opportunity is deliberately allowed to sit alone. One with no story map is
one that has not been taken on — still a candidate, or decided against. livt
records no status for this: whether the map exists **is** the record, the same
way a rule's automation is recorded rather than inferred.
