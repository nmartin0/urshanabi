# Write service

Applies approved changes, reconciles edits with source data under each
field's declared rule, checks the version every change was based on,
and pushes changes back to sources only where writeback is enabled.

- **Language role:** systems language
- **Holds:** data-write and writeback credentials
- **Separate because:** it alone holds the credentials that write data
  and write back to sources (R-80)
- **Speaks to:** the workflow worker's approved actions
- **Roadmap:** R-39, R-100, R-102, R-103

**Status:** planned; nothing is built yet. When work starts, this
project gains `script/bootstrap`, `script/test` and
`script/test-integration` (`RULES.md` E5).
