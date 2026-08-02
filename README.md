<h1 align="center">livt</h1>

<p align="center">
  <b>Collaborate on board. Make it living in text.</b>
</p>

<p align="center">
  <a href="https://boykush.github.io/livt/">Documentation</a> |
  <a href="https://boykush.github.io/livt/getting-started.html">Getting Started</a> |
  <a href="https://boykush.github.io/livt/guides/stories.html">Guides</a> |
  <a href="https://boykush.github.io/livt/reference/commands.html">Commands</a>
</p>

## What is livt?

livt is a CLI tool that captures collaborative discovery outcomes as living text. It bridges the gap between synchronous discovery sessions (like [User Story Mapping](https://boykush.github.io/livt/guides/story-maps.html) and [Example Mapping](https://boykush.github.io/livt/guides/example-mappings.html)) and development artifacts.

Discovery outcomes are written as plain text files (YAML, Markdown) and visualized as boards:

![Story Map board](docs/src/images/story-map.png)

## Features

- **Stories as Markdown** -- Write stories with YAML frontmatter, keep them alongside your code
- **Story Maps** -- Visualize activities, steps, and stories with release slices on a board
- **Example Mappings** -- Render rules, examples, and questions as color-coded sticky notes
- **Personas** -- Name who each story is for, and browse every actor beside the stories that name it
- **Ubiquitous Language** -- Keep shared terms as Markdown files and browse them as a glossary table
- **Static HTML output** -- `livt build` generates a standalone site, no runtime required
- **Local dev server** -- `livt serve` builds and serves with one command

## Installation

See the [Installation guide](https://boykush.github.io/livt/installation.html) for details on how release artifacts are verified.

### Prebuilt binaries

Download the archive for your platform from [GitHub Releases](https://github.com/boykush/livt/releases).

Every artifact is listed in `checksums.txt` (SHA-256) and carries a Sigstore-signed [build provenance attestation](https://github.com/boykush/livt/attestations), so you can verify it was built by this repository's release workflow before running it:

```bash
gh attestation verify livt_<version>_<os>_<arch>.tar.gz --repo boykush/livt
```

### mise

```bash
mise use "github:boykush/livt@<version>"
```

The [`github` backend](https://mise.jdx.dev/dev-tools/backends/github.html) downloads the release binary and, by default, verifies its GitHub artifact attestation (Sigstore build provenance). With [`lockfile = true`](https://mise.jdx.dev/dev-tools/mise-lock.html), the checksum and provenance are also pinned in `mise.lock`.

To build from source instead, use the [`go` backend](https://mise.jdx.dev/dev-tools/backends/go.html) (requires Go; verified by the [Go checksum database](https://sum.golang.org/)):

```bash
mise use "go:github.com/boykush/livt@<version>"
```

### go install

```bash
go install github.com/boykush/livt@<version>
```

## Quick Start

```bash
# Create the directory structure
mkdir -p stories personas discoveries/usm discoveries/example-mappings ubiquitous

# Create your first story
cat <<'EOF' > stories/my-first-story.md
---
name: My first story
---

As a user
I want to do something
So that I get value
EOF

# Build and serve
livt serve
```

Open http://localhost:3000 in your browser.

See the [Getting Started guide](https://boykush.github.io/livt/getting-started.html) for more details.

## File Structure

```
stories/
  {story-key}.md                     # Story files
discoveries/
  usm/
    {map-name}.yaml                  # Story map files
  example-mappings/
    {story-key}.yaml                 # Example mapping files
```

See [File Structure reference](https://boykush.github.io/livt/reference/file-structure.html) for output details.

## Contributing

Issues and pull requests are welcome — see [CONTRIBUTING.md](CONTRIBUTING.md). For vulnerability reports, see [SECURITY.md](SECURITY.md).

## License

[MIT](LICENSE)
