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

## `livt version`

Print the version of livt.

```bash
livt version
```
