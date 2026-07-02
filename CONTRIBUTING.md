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

The documentation site (`docs/`) is built with mdBook:

```bash
mdbook serve docs
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

Dry-run the release build locally without publishing:

```bash
mise run snapshot
```
