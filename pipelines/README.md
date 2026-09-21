# Pipelines

Pipeline definitions in Urshanabi's own format, compiled to the data
orchestrator, and the transformations that turn raw data into cleaned
and curated data. A separate project because it follows the
orchestrator's installation constraints.

- **Language role:** agent and pipeline language, and SQL
- **Holds:** nothing beyond its workload identity
- **Speaks to:** the orchestrator and the table format
- **Roadmap:** R-89, R-90, R-96, R-97, R-98, R-107

**Status:** planned; nothing is built yet. When work starts, this
project gains `script/bootstrap`, `script/test` and
`script/test-integration` (`RULES.md` E5).
