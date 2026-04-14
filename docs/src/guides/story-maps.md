# Story Maps

Story maps are YAML files stored in `discoveries/usm/`. They define the structure of a User Story Map with activities, steps, stories, and release slices.

## Format

```yaml
name: Map Name

activities:
  - key: activity-key
    name: Activity Name
    steps:
      - key: step-key
        name: Step Name
        stories:
          - key: story-key
            name: Story Card Name

releases:
  - name: Release Name
    stories:
      - story-key
```

## Releases

- Each release defines a horizontal divider on the board
- Stories listed in a release appear above its divider
- A release without `name` defaults to "Release N" based on position
- Stories not in any release appear below all dividers
- A story cannot belong to multiple releases

## Example

`discoveries/usm/collaborative-discovery.yaml`:

```yaml
name: Collaborative Discovery

activities:
  - key: story-mapping
    name: Story Mapping
    steps:
      - key: discover-stories
        name: Discover stories
        stories:
          - key: confirm-story-context
            name: Confirm story context
          - key: confirm-story-map
            name: Confirm story map
      - key: slice-releases
        name: Slice into releases
        stories:
          - key: split-release-scope
            name: Split release scope

  - key: discovery
    name: Discovery
    steps:
      - key: discover-rules
        name: Discover rules
        stories:
          - key: confirm-discovery-outcomes
            name: Confirm discovery outcomes

releases:
  - name: Walking Skeleton
    stories:
      - confirm-story-context
      - confirm-story-map
      - confirm-discovery-outcomes
  - stories:
      - split-release-scope
```

![Story map board](../images/story-map.png)
