# Ingestion

Pulls data from each source into the raw layer in bounded batches,
refuses a source whose shape has changed, and only appends: every
other table write goes through the reference implementation (R-110).

- **Language role:** systems language
- **Holds:** the source credentials, verified read-only
- **Separate because:** it alone holds the source systems'
  credentials, which nothing else may use (R-80)
- **Speaks to:** the source systems, and the raw layer
- **Roadmap:** R-53, R-55, R-62, R-93, R-110

**Status:** planned; nothing is built yet. When work starts, this
project gains `script/bootstrap`, `script/test` and
`script/test-integration` (`RULES.md` E5).
