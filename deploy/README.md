# Deployment

The release's deployment package, which runs a cell on any conformant
container orchestrator, and the live configuration of every hosted
cell: which release it runs, its sizing and its secret references
(roadmap R-136, R-137). Secrets are referenced, never stored (R-92).
While the repository is public (R-17), cells are named by opaque
identifiers, and no customer's name, endpoint, sizing or other
identifying detail may appear here.

- **Language role:** declarative configuration
- **Holds:** nothing; secrets are references (R-92)
- **Speaks to:** the orchestration platform
- **Roadmap:** R-52, R-63, R-77, R-126

**Status:** planned; nothing is built yet. When work starts, this
project gains `script/bootstrap`, `script/test` and
`script/test-integration` (`RULES.md` E5).
