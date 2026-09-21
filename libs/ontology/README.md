# Ontology library

The one implementation of what the ontology means: its model, its
compiler and its validation. The ingestion, indexer, query and write
services and the control plane's compiler all use it, so no two
components can disagree about the ontology (R-51).

- **Language role:** systems language
- **Holds:** nothing
- **Speaks to:** nothing on the network; it is a library
- **Roadmap:** R-51, R-70

**Status:** planned; nothing is built yet. When work starts, this
project gains `script/bootstrap`, `script/test` and
`script/test-integration` (`RULES.md` E5).
