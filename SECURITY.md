# Security Policy

## Supported versions

Only the latest release receives security fixes.

## Reporting a vulnerability

Report vulnerabilities privately via [GitHub security advisories](https://github.com/boykush/livt/security/advisories/new) ("Report a vulnerability" on the Security tab) — please do not open a public issue. Reports can be written in English or Japanese.

This is a personal open source project maintained on a best-effort basis, but you can expect an initial response within a week.

## Release integrity

Release artifacts are built and published exclusively by the [release workflow](.github/workflows/release.yml) on GitHub Actions:

- Binaries are built from the module as served by the Go module proxy and verified against the [Go checksum database](https://sum.golang.org/) — the same code `go install` would fetch.
- Every artifact is listed in `checksums.txt` (SHA-256).
- Every artifact carries a Sigstore-signed [build provenance attestation](https://github.com/boykush/livt/attestations) tying it to the exact source commit and workflow run. Verify before running:

  ```bash
  gh attestation verify livt_<version>_<os>_<arch>.tar.gz --repo boykush/livt
  ```

Installs via `go install` or the mise `go:` backend compile from source and are likewise verified by the Go checksum database.
