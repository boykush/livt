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
  index.html                              # Example mappings overview (home)
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

Every page shares a left sidebar that links the four resource types
(Example Mappings, Story Maps, Stories, Ubiquitous Language). The overview
pages render each example mapping and story map as a preview card.
