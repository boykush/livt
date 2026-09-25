<h1 align="center">livt</h1>

<p align="center">
  <b>Collaborate on board. Make it living in text.</b>
</p>

<p align="center">
  <a href="https://boykush.github.io/livt/">Why livt</a> |
  <a href="#getting-started">Getting Started</a> |
  <a href="https://boykush.github.io/livt/demo/">Live Demo</a>
</p>

## What is livt?

livt keeps what a discovery session agreed as text in a repository, shows it as the board it came from, and serves it to the coding agents that build it — so the agreement is still readable, and still quotable, when someone implements it months later.

![Story Map board](docs/src/images/story-map.png)

[Why livt](https://boykush.github.io/livt/) says which problems it is built against and when it is the right tool.

## Getting started

Install livt. [mise](https://mise.jdx.dev/) checks the release's build provenance before installing it:

```bash
mise use "github:boykush/livt@<version>"
```

A binary from [GitHub Releases](https://github.com/boykush/livt/releases), or `go install github.com/boykush/livt@<version>`, works as well; [SECURITY.md](SECURITY.md#release-integrity) says how each release is built and how to verify it.

Then give your coding agent the skills that write a livt repository. Each plugin under [plugins/](plugins) says in its README what it is for and how to install it; [livt-discovery](plugins/livt-discovery/README.md) is where a team starts.

After a session, give your agent the board — a photo, an export, or the stickies pasted as text — and ask it to record it. In the repository it wrote to, `livt serve` shows the result at http://localhost:3000.

There is no command or format reference to read: each detail lives where it cannot drift, and [Reference](https://boykush.github.io/livt/reference.html) says where.

## Contributing

Issues and pull requests are welcome — see [CONTRIBUTING.md](CONTRIBUTING.md). For vulnerability reports, see [SECURITY.md](SECURITY.md).

## Acknowledgements

livt renders practices it did not invent. [livt and BDD](https://boykush.github.io/livt/livt-and-bdd.html) credits the people who wrote them down, says what livt added and where it departs on purpose, and links where their maintainers are funded.

## License

[MIT](LICENSE)
