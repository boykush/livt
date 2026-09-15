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

Name the **opportunity**, not the subject area it sits in. "Collaborative
discovery" names a topic; "stale discovery" names something you could decide to
fix. The topic already has a home — it is what the
[story map](./story-maps.md) is called, and a map and its opportunity are free to
carry different names precisely because they answer different questions: the map
says which journey was drawn, the opportunity says why anyone drew it.

The **opportunity key** is derived from the filename (without `.md`), and must
be kebab-case. Key uniqueness is enforced by the filesystem.

## Example

`opportunities/collaborative-discovery.md`:

```markdown
---
name: Stale discovery
repos:
  - boykush/livt
---

Discovery outcomes sit in the board tool and go stale. The raw record is rough, and
nothing versions it or checks it for consistency. Recording it, and keeping it
updated as rules change, saves re-running the same argument and lets the outcome
stand as the specification.
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

## The opportunity dashboard

How far an opportunity has been taken is a **page of its own**, linked from the
opportunity under Related:

```
opportunity/collaborative-discovery.html               # what the opportunity is
opportunity-progress/collaborative-discovery.html      # how far it has got
```

This is the same split the canvas makes, for the same reason: what an
opportunity *is* reads the same on every visit, and the progress is the part
that changes. Mixing them would put a figure that moves weekly above a statement
that has not moved since the opportunity was framed. An opportunity with no
story map has no progress to read, so it gets no page and links to none.

### Four meters

The page leads with four, each drawn as a share of the whole it belongs to — a
bare `2` does not say whether a board is nearly agreed, and a filled bar does:

- **Stories through an example mapping** — of the stories its maps hang under
  the backbone, how many have had their conversation
- **Un-automated rules** — how many of its rules are waiting for a test
- **Proposed rules** — how many are waiting for an agreement
- **Open questions** — against every question its boards have asked, settled
  ones included

Each gets a meter of its own rather than becoming a segment of another's. The
rule figures do share a denominator, and one stacked bar would say that
arithmetic — but what *closes* each of them differs (a test, an agreement, a
conversation), and a segment inside someone else's bar is not a thing anyone can
go and do. They are counted the way the [Tasks page](../reference/file-structure.md)
splits them, so a proposal is never also counted as something waiting for a test.

**What is automated has no meter.** The breakdown below already says it, slice by
slice and story by story, and says *where* — a figure at the top that only
restates the section under it costs a glance and settles nothing. It was also the
one meter with nowhere to lead, since no page on the site lists automated rules.
The tiles on the Opportunities list do carry it, because a tile has no room for a
breakdown and one number is all it can offer.

### Stories, in their release slices

Below the gauges the stories are grouped by the **release slice** their map puts
them in, in the map's declared order, with the stories it left unscoped last.
Each slice carries its own coverage, because that is the number the next release
turns on rather than the opportunity's. A map that declares no release leaves
one unnamed group holding every story, and that group gets no heading — a lone
"no release" row would name a distinction the map never drew.

Each row draws its own coverage bar beside the fraction — the fraction alone does
not say that 4/5 and 19/20 are the same measure the meters take, and bars do. The
column is named once, at the head of the section, **with the bars' own colour
beside the name**: a label in the corner says what the section is, but the swatch
is what tells a reader that the blue bar three rows down is the thing being
named. A story with no example mapping says so in words rather than showing
`0/0`, which would read as a conversation that found no rules.

### What the figures link to

One link sits at page level, beside the subtitle: the opportunity's own **Example
Mappings**, narrowed to it. It hangs off no meter, because a reader who came to
read the boards is not answering a figure, and hanging it under one would turn it
into a claim about that count.

Each meter then carries **one link, in its footer**, set below a hairline in the site's
link colour and naming the page it leads to. It is the only thing on the card
that navigates: a figure that looks like it leads somewhere, paired with a link
that does not look like one, is what a dashboard cannot afford. A story row opens
its example mapping, or the story's card until there is one.



Every destination is a page that already renders what the meter counted, and it
arrives **at the list it counted** — the Tasks page keeps three, and landing on
the first leaves the reader to find the other two — narrowed to this opportunity
through the same filter those lists already carry. Nothing here is a second copy
of them.

Retired rules are counted in neither gauge, the same way the Tasks page and the
sidebar badge leave them out: a rule the board has closed is not spec anyone is
waiting on. A retired question is counted only in the *whole* an open question is
measured against — a question that has been settled is part of the record of how
much the boards have worked through.

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
