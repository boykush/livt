# Contributing to livt

Thanks for your interest in contributing! Issues and pull requests are welcome, in English or Japanese.

## Development setup

Tools and tasks are managed with [mise](https://mise.jdx.dev/):

```bash
git clone https://github.com/boykush/livt.git
cd livt
mise install
```

`mise install` installs the pinned toolchain (Go, golangci-lint, ...) and, via the `postinstall` hook, the [hk](https://hk.jdx.dev/) pre-commit hooks — so the same checks run locally and in CI. One of the docs tools, mdbook-i18n-helpers, is built from source by cargo, so have [Rust](https://rustup.rs/) installed first.

## Checks and tests

Every CI check is a mise task (see `.mise.toml`):

```bash
mise run check   # fmt + vet + lint + tidy + build + test + docs:check
mise run test    # go test -race ./...
```

## Docs

The documentation site (`docs/`) is built with mdBook, in English and in Japanese:

```bash
mdbook serve docs
```

```bash
MDBOOK_BOOK__LANGUAGE=ja mdbook serve docs
```

It says why livt exists and when it fits. Every other detail has a home of its own, beside what it describes; [Reference](docs/src/reference.md) says where.

The English under `docs/src/` is the source. The Japanese is `docs/po/ja.po`, a gettext catalog that translates it passage by passage through [mdbook-i18n-helpers](https://github.com/google/mdbook-i18n-helpers); a passage the catalog does not translate yet shows in English rather than in an outdated translation.

After you change `docs/src/`, run `mise run docs:po` and translate what it leaves empty or marks `fuzzy`, dropping the `#, fuzzy` line once the translation fits. `mise run check` fails until you have.

A heading that a link points at pins its anchor, as in `## Where livt stops {#where-livt-stops}`: the Japanese book derives its anchors from the translated headings.

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
