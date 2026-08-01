# File Structure

## Input

```
project-root/
  stories/
    {story-key}.md                        # Story files
  discoveries/
    usm/
      {map-name}.yaml                     # Story map files
    example-mappings/
      {story-key}.yaml                    # Example mapping files
  ubiquitous/
    {term-key}.md                         # Ubiquitous language term files
```

- Story keys are derived from filenames (without extension)
- Story keys must be kebab-case: lowercase letters, numbers, and hyphens
- The `stories/` directory is the committed story registry, and `stories/{story-key}.md` provides story key uniqueness
- Example mapping filenames must match story keys to link them
- Term keys are derived from `ubiquitous/{term-key}.md` filenames and used as `ubiquitous.html#{term-key}` link anchors

## Output

`livt build` generates the following structure:

```
dist/
  index.html                              # Overview: open questions and un-automated rules (home)
  example-mappings.html                   # Example mappings overview
  story-maps.html                         # Story maps overview
  stories.html                            # Story list
  ubiquitous.html                         # Ubiquitous language table
  story/
    {story-key}.html                      # Story detail pages
  mapping/
    {story-key}.html                      # Example mapping boards
  story-map/
    {map-name}.html                       # Story map boards
```

Every page shares a left sidebar that links the home page (Overview) and the
four resource types (Example Mappings, Story Maps, Stories, Ubiquitous
Language). The example mapping and story map overviews render each board as a
preview card.

The home page gathers what the master leaves unfinished, so neither has to be
hunted for board by board:

- **Open Questions** — every `questions` entry across the example mappings.
  These close by a conversation, so they feed the next discovery session.
- **Un-automated Rules** — every rule with no `automated: true` recorded. These
  close by a test, so they read as the list of behaviour still to build.

Each item names the story it came from and links to its own sticky on that
story's mapping board. Both lists are filtered together by opportunity, and the
selection is mirrored in the `?opportunity=` query parameter so a filtered view
is shareable.
