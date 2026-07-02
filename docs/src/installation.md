# Installation

livt ships as a single static binary. Every release is built by the
[release workflow](https://github.com/boykush/livt/blob/main/.github/workflows/release.yml)
on GitHub Actions, and every artifact is:

- listed in `checksums.txt` (SHA-256)
- signed with a [GitHub artifact attestation](https://github.com/boykush/livt/attestations)
  (Sigstore build provenance) tying it to the exact source commit and workflow run
- built from the module as served by the Go module proxy and verified against the
  [Go checksum database](https://sum.golang.org/), so the binaries can be
  reproduced from the public source

Pick the method that matches how much of that verification you want automated.

## mise (recommended)

The [mise `github` backend](https://mise.jdx.dev/dev-tools/backends/github.html)
downloads the release binary and verifies its build provenance attestation by
default — a tampered or foreign-built asset fails to install.

```bash
mise use "github:boykush/livt@<version>"
```

Or in `mise.toml`:

```toml
[settings]
lockfile = true # pin checksum + provenance in mise.lock

[tools]
"github:boykush/livt" = "<version>"
```

With [`lockfile = true`](https://mise.jdx.dev/dev-tools/mise-lock.html), the
asset checksum and attestation provenance are recorded in `mise.lock`, so
later installs — on CI or a teammate's machine — must match bit-for-bit.

## Manual download

Download the archive for your platform and `checksums.txt` from
[GitHub Releases](https://github.com/boykush/livt/releases), then verify
before running:

```bash
# Provenance: proves the artifact was built by this repository's release
# workflow on GitHub Actions (requires the GitHub CLI)
gh attestation verify livt_<version>_<os>_<arch>.tar.gz --repo boykush/livt

# Integrity: check the SHA-256 checksum (macOS: shasum -a 256 -c)
sha256sum --check --ignore-missing checksums.txt

tar -xzf livt_<version>_<os>_<arch>.tar.gz
install livt ~/.local/bin/ # or anywhere on your PATH
```

## Build from source

Source installs are verified by the Go checksum database: everyone gets
byte-identical source for a given version, and a silently re-pushed tag is
rejected. Pin a release version rather than `latest` so installs stay
reproducible.

```bash
go install github.com/boykush/livt@<version>
```

Or with the [mise `go` backend](https://mise.jdx.dev/dev-tools/backends/go.html)
(requires Go on your PATH):

```bash
mise use "go:github.com/boykush/livt@<version>"
```
