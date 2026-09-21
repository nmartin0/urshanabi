# Ontology library

The one implementation of what the ontology means: its model, its
compiler and its validation. The ingestion, indexer, query and write
services and the control plane's compiler all use it, so no two
components can disagree about the ontology (R-51).

- **Language role:** systems language
- **Holds:** nothing
- **Speaks to:** nothing on the network; it is a library
- **Roadmap:** R-51, R-70

**Status:** loads an ontology written in the format
`docs/ontology-format.md` specifies, and refuses one that breaks its
rules, reporting every problem at once. Checking classifications
against lineage, and sources against the table catalogue, waits for
the pipelines that provide them.

- `script/test` runs the systems-language gates (`script/systems-gates`).
- `script/test-integration` loads a complete fixture ontology in three
  languages, breaks each rule in turn, and parses every example in the
  specification.
