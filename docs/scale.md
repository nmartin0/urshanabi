# Scale objectives

What a single Urshanabi cell must sustain (roadmap R-66). The owner
decided on 2026-09-21 that the first release is proven against the
largest enterprises and governments, so every objective below is set for
that size; smaller cells need less of everything.

The figures are targets proposed from precedent, not measurements: the
leading platform's published limits, established response-time
thresholds, and the workflow engine maintainers' sizing guidance. The
first load tests may revise them, and any revision is recorded here with
its evidence.

`script/check-scale` enforces this file: every objective has a number
and a unit, names the roadmap items that will prove it, and an objective
marked tested must cite an executable load test under `load/`.

## Objectives, per cell

| ID | Objective | Target | Roadmap | Status | Evidence |
|---|---|---|---|---|---|
| SO-01 | Objects held in the ontology | at least 50,000,000,000 objects | R-70 | untested | — |
| SO-02 | People working at once | at least 20,000 people | R-69, R-70 | untested | — |
| SO-03 | Rows ingested | at least 2,000,000,000 rows per hour | R-53 | untested | — |
| SO-04 | Durable workflows started | at least 1,000 workflows per second at peak | R-57 | untested | — |
| SO-05 | Opening one record | at most 100 ms at the 95th percentile | R-70 | untested | — |
| SO-06 | A filtered search | at most 500 ms at the 95th percentile | R-69, R-70 | untested | — |
| SO-07 | An aggregate over a billion objects | at most 2 s at the 95th percentile | R-69 | untested | — |
| SO-08 | A link traversal's results | at least 10,000,000 objects, larger results streamed rather than refused | R-69, R-114 | untested | — |
| SO-09 | Real-time freshness, source to screen | at most 5 s at the 95th percentile | R-107 | untested | — |
| SO-10 | Near-real-time freshness, source to screen | at most 5 min at the 95th percentile | R-107 | untested | — |
| SO-11 | An edit becoming searchable | at most 1 s after it is acknowledged | R-70 | untested | — |

## Where the figures come from

- **Response times.** About a tenth of a second feels instant, a second
  keeps a person's flow of thought, and ten seconds is the limit of
  attention. Opening a record must feel instant; searches and aggregates
  must not break flow.
- **The leading platform's current limits.** A link traversal returns at
  most 10 million objects, one call loads at most 100,000 into memory,
  and an action edits at most 10,000 objects across 50 types. There is
  no fixed limit on objects per type: the index's disk space is the
  ceiling. SO-08 matches the traversal limit and goes further by
  streaming larger results instead of refusing them.
- **The durable workflow engine.** Its history shard count is fixed when
  a cluster is provisioned and bounds how many hosts the cluster can
  ever grow to. Its maintainers advised at least 4,000 shards for about
  500 workflows a second; for SO-04's 1,000, every cell's cluster is
  provisioned with 8,192 shards, confirmed by a load test before the
  first cluster is created (R-57).
