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

## `livt open`

Open a story card for conversation in your default editor.

```bash
livt open <story-key>
```

## `livt version`

Print the version of livt.

```bash
livt version
```
