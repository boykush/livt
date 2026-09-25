# Contributing to livt

Thanks for your interest in contributing! Issues and pull requests are welcome, in English or Japanese.

## Development setup

Tools and tasks are managed with [mise](https://mise.jdx.dev/):

```bash
git clone https://github.com/boykush/livt.git
cd livt
mise install
```

`mise install` installs the pinned toolchain (Go, golangci-lint, ...) and, via the `postinstall` hook, the [hk](https://hk.jdx.dev/) pre-commit hooks — so the same checks run locally and in CI.

## Checks and tests

Every CI check is a mise task (see `.mise.toml`):

```bash
mise run check   # fmt + vet + lint + tidy + build + test
mise run test    # go test -race ./...
```

## Docs

The documentation site is plain HTML in `docs/site/`, in English and in Japanese, and Pages publishes it as it is. It says why livt exists and when it fits. Every other detail has a home of its own, beside what it describes; [Reference](docs/content/en/reference.md) says where.

What the site says is its text: one Markdown file per page in `docs/content/en/` and `docs/content/ja/`, each language written in its own language rather than translated from the other. A change edits the text, then renders the HTML from it with the [render-site](.apm/skills/render-site/SKILL.md) skill, and the two go in one commit. `mise run check` fails when a page no longer shows its text.

To look at the site:

```bash
python3 -m http.server --directory docs/site
```

## Commits and pull requests

- Use conventional-commit style titles (`feat:`, `fix:`, `docs:`, ...) — the release changelog is grouped by them.
- Target `main`; CI must pass.

## Releasing (maintainers)

Releases are built by [GoReleaser](https://goreleaser.com/) in the [release workflow](.github/workflows/release.yml), which runs the full checks, publishes platform binaries and `checksums.txt` to GitHub Releases, and signs a build provenance attestation for every artifact (see [SECURITY.md](SECURITY.md)).

```bash
git tag v0.x.y
git push origin v0.x.y
```

The published notes are GoReleaser's changelog, grouped by commit type. Once the release is up, the [`enrich-release-notes`](.apm/skills/enrich-release-notes/SKILL.md) skill writes the part above it for people who use livt — who is affected, what to do, how to upgrade — and publishes it after you approve the draft.

Dry-run the release build locally without publishing:

```bash
mise run snapshot
```
