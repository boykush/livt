# Commands

## `livt serve`

Build artifacts and start a local server.

While the server is running, livt watches the input directories
(`discoveries/example-mappings`, `stories`, `discoveries/usm`, and `ubiquitous`).
When a file changes, livt rebuilds and reloads the page in the browser
automatically, so you can preview refinements while editing.

```bash
livt serve [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--port` | `-p` | `3000` | Port to listen on |
| `--out` | `-o` | `dist` | Output directory |

## `livt build`

Build static HTML from artifacts without starting a server.

```bash
livt build [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--out` | `-o` | `dist` | Output directory |

## `livt story commit`

Commit a story candidate to detailed discovery by creating its story Markdown file.
The story key must be kebab-case. When committing from a story map, livt creates the story file first and then writes the key back to the matching story candidate in the story map.

```bash
livt story commit --usm discoveries/usm/map.yaml --candidate "Story name" --key story-key
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--key` | | | Story key, in kebab-case |
| `--usm` | | | Story map YAML file |
| `--candidate` | | | Story candidate name in the story map |
| `--name` | | | Story display name for standalone story creation |
| `--as` | | | Story persona |
| `--want` | | | Story goal |
| `--so-that` | | | Story benefit |
| `--stories-dir` | | `stories` | Stories directory |

## `livt version`

Print the version of livt.

```bash
livt version
```
