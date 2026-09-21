# Contracts

Every service, event and error contract in Urshanabi, and the only
place any of them exists (`RULES.md` E2). Clients are generated from
here at build time; consumer expectation files for contract tests live
here too.

- **Language role:** schema-first definitions; generated code is never
  committed
- **Holds:** nothing
- **Speaks to:** every service, and the interface through its
  generated client
- **Roadmap:** R-06, R-78, R-114

**Status:** planned; nothing is built yet. When work starts, this
project gains `script/bootstrap`, `script/test` and
`script/test-integration` (`RULES.md` E5).
