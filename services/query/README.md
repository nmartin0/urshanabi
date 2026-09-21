# Query service

Answers every question from people and agents and enforces policy on
every answer: security values filter first, no signal derives from
hidden rows, and large results stream in columnar form.

- **Language role:** systems language
- **Holds:** index and curated-layer read access
- **Speaks to:** the gateway's and the agent's requests
- **Roadmap:** R-21, R-68, R-69, R-70, R-71, R-114, R-116

**Status:** walking skeleton. It serves the build-identity contract,
`urshanabi.build.v1.BuildService`, generating its code from
`contracts/` at build time, and logs each caller's request id.

- `script/bootstrap` installs the toolchain pinned in
  `rust-toolchain.toml` where its installer exists.
- `script/test` runs the gates: formatting, lints at the strictest level
  with warnings as errors, unit tests, and the dependency policy in
  `deny.toml` (`RULES.md` E3).
- `script/test-integration` starts the real server and calls it over
  the network.
