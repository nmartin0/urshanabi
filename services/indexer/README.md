# Indexer

Keeps the ontology index current from the curated layer's change log,
stamping every write with its source version so the index never goes
backwards. The only writer of the index.

- **Language role:** systems language
- **Holds:** index write access
- **Separate because:** it alone holds the index's write credentials
  (R-80)
- **Speaks to:** the change log and the serving engine
- **Roadmap:** R-70, R-105, R-111

**Status:** planned; nothing is built yet. When work starts, this
project gains `script/bootstrap`, `script/test` and
`script/test-integration` (`RULES.md` E5).
