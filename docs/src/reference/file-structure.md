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
```

- Story keys are derived from filenames (without extension)
- Example mapping filenames must match story keys to link them

## Output

`livt build` generates the following structure:

```
dist/
  index.html                              # Story list
  story/
    {story-key}.html                      # Story detail pages
  mapping/
    {story-key}.html                      # Example mapping boards
  story-map/
    {map-name}.html                       # Story map boards
```
