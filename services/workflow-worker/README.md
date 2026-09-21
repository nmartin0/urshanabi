# Workflow worker

Runs approvals, expiries, retries and automations on the durable
workflow engine, with maker-checker rules for consequential actions.

- **Language role:** services language
- **Holds:** nothing beyond its workload identity
- **Speaks to:** the durable workflow engine and the write service
- **Roadmap:** R-38, R-42, R-57

**Status:** planned; nothing is built yet. When work starts, this
project gains `script/bootstrap`, `script/test` and
`script/test-integration` (`RULES.md` E5).
