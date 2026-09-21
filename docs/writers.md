# One writer per store

Every store Urshanabi keeps, and the single identity allowed to write it
(roadmap R-79). Ownership is enforced by the store refusing everyone
else, never by convention.

**One writer means one owning identity.** An owner may write through
more than one mechanism, but always under its own identity and
credentials. Ingestion, for example, appends raw data through the
systems-language table library, and runs compaction and retention on the
same tables through the table format's reference implementation (R-110)
— both as ingestion.

**Two conditions make per-table enforcement work.** The table catalog
grants privileges on individual tables and vends short-lived storage
credentials scoped to one table's storage prefix; its privileges are
enforced correctly only when each table's directory holds that table's
files alone and the directory hierarchy matches the namespace hierarchy.
Every table is laid out that way. And because one of the catalog's
releases documented privilege changes taking up to an hour to apply,
that delay is measured against the one-hour access-removal bound (R-84)
before R-84 is closed.

Services never write each other's stores: a component that needs
another's data calls its contract. Audit records, for example, are
submitted to the audit service, which alone writes the audit log.

`script/check-writers` keeps this file well-formed: every store has
exactly one owner, the owner is an existing component or a named
platform role, and every row says how the store enforces it.

## Stores

| Store | Owner | Enforced by | Roadmap |
|---|---|---|---|
| Raw-layer tables, their change logs and their maintenance | `services/ingestion` | Per-table write grants in the table catalog, with vended storage credentials scoped to each table's prefix | R-55, R-104, R-110 |
| Cleaned-layer tables | `pipelines` | Per-table write grants and scoped credentials; each table's single owning pipeline, run once at a time by the orchestrator | R-98, R-104 |
| Curated-layer tables, including merged edits published back | `pipelines` | Per-table write grants and scoped credentials, as above | R-102, R-104 |
| Ontology index | `services/indexer` | The serving engine's write credential is held by the indexer alone | R-70, R-105 |
| Edits store | `services/write` | A database role with write privileges granted only to the write service | R-102, R-103 |
| Sessions, login attempts and idempotency records | `services/gateway` | A key-value store user with write access granted only to the gateway | R-45, R-46, R-106 |
| Configuration, tenants and signed bundles | `services/control-plane` | A database role with write privileges granted only to the control plane | R-51, R-73 |
| Audit log | `services/audit` | Storage write credentials held by the audit service alone; other services submit records through its contract | R-23, R-88 |
| The durable workflow engine's persistence | platform: durable workflow engine | The engine alone holds its database credentials; the workflow worker uses its API | R-57 |
| The data orchestrator's metadata | platform: data orchestrator | The orchestrator alone holds its database credentials; pipelines use its API | R-96 |
| The table catalog's own records | platform: table catalog | The catalog alone holds its database credentials; every component commits through its API | R-104 |

Event-log topics join this table as each is designed: every topic has
exactly one producing component, enforced by per-topic write permissions
in the event log.
