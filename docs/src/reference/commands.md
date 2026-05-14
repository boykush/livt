# Commands

## `livt serve`

Build artifacts and start a local server.

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

## `livt story init`

Initialize a story Markdown file.

```bash
livt story init <key> --name "Story name"
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--name` | | | Story display name |
| `--stories-dir` | | `stories` | Stories directory |

## `livt version`

Print the version of livt.

```bash
livt version
```
