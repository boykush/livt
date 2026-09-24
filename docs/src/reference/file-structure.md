# File Structure

## Input

```
project-root/
  livt.yaml                               # Site config (optional)
  opportunities/
    {opportunity-key}.md                  # Opportunity files
  stories/
    {story-key}.md                        # Story files
  discoveries/
    opportunity-canvases/
      {opportunity-key}.yaml              # Opportunity canvas files
    usm/
      {opportunity-key}.yaml              # Story map files
    example-mappings/
      {story-key}.yaml                    # Example mapping files
  ubiquitous/
    {term-key}.md                         # Terms holding across contexts
    {ctx}/
      {term-key}.md                       # Terms scoped to one context
  automations/
    {owner}/{repo}.json                   # Collected automation reports
```

- `livt.yaml` configures the site build — see [Configuration](./configuration.md). Without it, every setting takes its default
- Opportunity keys, like story keys, are derived from filenames (without extension) and must be kebab-case
- An [opportunity canvas](../guides/opportunities.md#the-opportunity-canvas) filename must match an opportunity key to link them
- A story map filename that matches an opportunity key marks the map as the journey mapped for that opportunity. The filename is the whole join, so an opportunity has one map, and the key is what addresses it. A map whose key matches no opportunity stands in as its own, named by the map — which is how every livt repository behaved before opportunities were files
- Story keys are derived from filenames (without extension)
- Story keys must be kebab-case: lowercase letters, numbers, and hyphens
- The `stories/` directory is the story registry, and `stories/{story-key}.md` provides story key uniqueness
- An example mapping whose filename matches a story key maps that story. A mapping needs no story file: without one it is called by its [`name`](../guides/example-mappings.md#a-mapping-without-a-story), or else by its key
- Term keys are derived from filenames, and a term's [context](../guides/ubiquitous-language.md#contexts) from the directory holding it. The path is what makes a term unique, so the same key can sit at the root and under a context as two separate terms
- A context is optional and one directory deep; terms nested deeper are not addressable and are left out of the glossary
- A term is anchored as `ubiquitous.html#{term-key}`, or `ubiquitous.html#{ctx}/{term-key}` when it is scoped
- `automations/` holds the reports [`livt automations`](./commands.md#livt-automations) collects from implementation repositories' tests, one file per repository. They are generated, not written by hand: livt derives from them which rules and examples are automated — see [Automating a rule](../guides/example-mappings.md#automating-a-rule)

## Output

`livt build` generates the following structure:

```
dist/
  index.html                              # Example mappings overview (home)
  opportunities.html                      # Opportunities overview
  story-maps.html                         # Story maps overview
  stories.html                            # Story list
  ubiquitous.html                         # Ubiquitous language table
  tasks.html                              # Open questions, proposed and un-automated rules
  opportunity/
    {opportunity-key}.html                # Opportunity detail pages
  opportunity-canvas/
    {opportunity-key}.html                # Opportunity canvas sheets
  opportunity-progress/
    {opportunity-key}.html                # Opportunity dashboards
  story/
    {story-key}.html                      # Story detail pages
  mapping/
    {story-key}.html                      # Example mapping boards
  story-map/
    {opportunity-key}.html                # Story map boards
  diff.html                               # The diff, only with --diff
```

Every page shares a left sidebar that links the five resource types (Example
Mappings, Opportunities, Story Maps, Stories, Ubiquitous Language) and, below
them, Tasks. The overview pages render each example mapping, opportunity
canvas, and story map as a preview card.

`tasks.html` gathers what the livt repository leaves unfinished, so none of it
has to be hunted for board by board:

- **Open Questions** — every `questions` entry across the example mappings.
  These close by a conversation, so they feed the next discovery session.
- **Proposed Rules** — every rule marked
  [`status: proposed`](../guides/example-mappings.md#proposing-a-rule). These
  close by agreement — accepted, or rejected when turned down — so they are the
  decisions still open.
- **Un-automated Rules** — every accepted rule that no collected
  [report](../guides/example-mappings.md#automating-a-rule) cites and that
  carries no deprecated `automated:`, including one whose examples are cited but
  not the rule itself. These close by a test, so
  they read as the list of behaviour still to build. A proposed rule is never
  here, even once a test covers it: a test cannot close what is not agreed yet.

[Closed](../guides/example-mappings.md#retiring-an-item) rules and retired
examples and questions are on none of the lists, and off the boards as well:
nothing can close them again, so they would sit here forever.

Each item names the mapping it came from, as the board's yellow sticky reads,
and links to its own sticky on that board. All three lists are filtered together by opportunity, and
the selection is mirrored in the `?opportunity=` query parameter so a filtered
view is shareable.

Items carry no issue state. A rule records its automation issue URLs but not
their open/closed state, and the build never queries the issue tracker, so an
item leaves this page only through what the livt repository records — a
question becoming a rule, a proposal being accepted, a rule becoming automated —
never through an issue closing.
