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
```

- `livt.yaml` configures the site build — see [Configuration](./configuration.md). Without it, every setting takes its default
- Opportunity keys, like story keys, are derived from filenames (without extension) and must be kebab-case
- An [opportunity canvas](../guides/opportunities.md#the-opportunity-canvas) filename must match an opportunity key to link them
- A story map filename that matches an opportunity key marks the map as the journey mapped for that opportunity. A map whose key matches no opportunity stands in as its own, named by the map — which is how every livt repository behaved before opportunities were files
- Story keys are derived from filenames (without extension)
- Story keys must be kebab-case: lowercase letters, numbers, and hyphens
- The `stories/` directory is the committed story registry, and `stories/{story-key}.md` provides story key uniqueness
- Example mapping filenames must match story keys to link them
- Term keys are derived from filenames, and a term's [context](../guides/ubiquitous-language.md#contexts) from the directory holding it. The path is what makes a term unique, so the same key can sit at the root and under a context as two separate terms
- A context is optional and one directory deep; terms nested deeper are not addressable and are left out of the glossary
- A term is anchored as `ubiquitous.html#{term-key}`, or `ubiquitous.html#{ctx}/{term-key}` when it is scoped

## Output

`livt build` generates the following structure:

```
dist/
  index.html                              # Example mappings overview (home)
  opportunities.html                      # Opportunities overview
  story-maps.html                         # Story maps overview
  stories.html                            # Story list
  ubiquitous.html                         # Ubiquitous language table
  tasks.html                              # Open questions and un-automated rules
  opportunity/
    {opportunity-key}.html                # Opportunity detail pages
  opportunity-canvas/
    {opportunity-key}.html                # Opportunity canvas sheets
  story/
    {story-key}.html                      # Story detail pages
  mapping/
    {story-key}.html                      # Example mapping boards
  story-map/
    {map-name}.html                       # Story map boards
```

Every page shares a left sidebar that links the five resource types (Example
Mappings, Opportunities, Story Maps, Stories, Ubiquitous Language) and, below
them, Tasks. The overview pages render each example mapping, opportunity
canvas, and story map as a preview card.

`tasks.html` gathers what the livt repository leaves unfinished, so neither kind has to
be hunted for board by board:

- **Open Questions** — every `questions` entry across the example mappings.
  These close by a conversation, so they feed the next discovery session.
- **Un-automated Rules** — every rule with no `automated: true` recorded. These
  close by a test, so they read as the list of behaviour still to build.

[Retired](../guides/example-mappings.md#retiring-an-item) items are on neither
list, and off the boards as well: nothing can close them, so they would sit here
forever.

Each item names the story it came from and links to its own sticky on that
story's mapping board. Both lists are filtered together by opportunity, and the
selection is mirrored in the `?opportunity=` query parameter so a filtered view
is shareable.

Items carry no status beyond being listed. A rule records its automation issue
URLs but not their open/closed state, and the build never queries the issue
tracker, so an item leaves this page by being resolved — a question becoming a
rule, a rule becoming automated — rather than by moving through states.
