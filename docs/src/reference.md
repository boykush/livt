# Reference

These docs keep no reference. A command reference, a guide to each file format, a table of what the MCP server returns — each would be a second copy of something the tool already says, and a copy goes stale unless someone remembers to update it. That is the failure livt exists to prevent, so livt does not ship one.

Each detail is written once, where it cannot drift from what it describes, and where the reader who needs it — nowadays, usually a coding agent — already looks:

| To know | Ask |
|---|---|
| What a command takes | `livt <command> --help` |
| How to write your livt repository: its files, their formats, the rules for changing them | the skills that write it, under [plugins](https://github.com/boykush/livt/tree/main/plugins) |
| How an agent reads your livt repository, and how it cites what it reads | the `livt mcp` server you run over it, which describes its own tools and resources |
| What a livt URI in your livt repository points at | `livt resolve <uri>` |
| What livt itself does, and why, rule by rule | livt's own livt repository: its [example mappings](https://github.com/boykush/livt/tree/main/discoveries/example-mappings), rendered on the [live demo](https://boykush.github.io/livt/demo/) |

The last row is livt used on itself. livt keeps its own rules in a livt repository, so "what does livt do when…" is answered the way livt answers it for any team: by the rule that was agreed, with the reason in the commit that agreed it.
