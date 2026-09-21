# Audit service

The single writer of the tamper-evident audit log: every record enters
an append-only hash tree whose signed checkpoints are published to a
witness the customer controls.

- **Language role:** services language
- **Holds:** audit-log write access
- **Speaks to:** audit events from every service
- **Roadmap:** R-23, R-88

**Status:** planned; nothing is built yet. When work starts, this
project gains `script/bootstrap`, `script/test` and
`script/test-integration` (`RULES.md` E5).
