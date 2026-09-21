# Technology review

Every technology role in Urshanabi was checked in September 2026 against
the newest and emerging alternatives, for massive scale, fast
processing, strong security and heavy distributed computing. Each got
one of four verdicts: **keep**, **upgrade** (same choice, newer version
or better configuration), **replace**, or **add** (a capability that was
missing). Where the only evidence for an alternative came from its own
vendor, that is said.

## Findings that changed the plan

1. **The systems-language table library can only append.** Overwrites,
   deletes, row-level updates and merges, and compaction were all still
   missing, and another language's library had just fixed a retry bug
   that could resurrect deleted rows. Every non-append write therefore
   goes through the table format's reference implementation in the
   distributed engines (R-110).
2. **The table format moves to its third version.** Row lineage gives
   every row a permanent identity and the sequence number of its last
   change, which becomes the index version; deletion vectors make
   updates cheaper; built-in table encryption keys support per-tenant
   keys (R-111).
3. **Post-quantum key exchange was added.** The hybrid exchange was
   standardised in August 2026 and is the default in major browsers and
   in one services-language standard library, while a 2026 measurement
   found government and defence adoption effectively absent (R-112).
4. **Confidential computing was added** for pooled hosted cells, with
   attestation-gated keys (R-113).
5. **Large results travel as columnar streams**, because row-by-row
   protocols can take far longer to hand over a result than to compute
   it (R-114).
6. **Agent-to-agent interoperability was added** alongside the
   agent-tool protocol (R-115).
7. **Query plans use the standard cross-engine representation**, the
   hook for GPU-native engines now emerging in research (R-116).

## Verdicts by role

| Role | Verdict | Why |
|---|---|---|
| Systems language | Keep | Memory safety without garbage collection; the default for new data engines |
| Services language | Keep | Post-quantum TLS by default; the durable workflow engine's primary kits |
| Agent and pipeline language | Keep | The model and data ecosystems are strongest there |
| Frontend language | Keep | Its compiler was rewritten natively in 2026, roughly ten times faster |
| Table format | Upgrade to version 3 | The format contest is effectively over; the others now publish into it |
| Table catalog | Keep | Became a top-level foundation project in February 2026 |
| Query engine | Keep | Emits standard plans for future accelerated engines |
| Distributed batch | Upgrade | A native accelerator inside the engine; a replacement written natively in the systems language claims 4 to 6 times the speed on a public benchmark but is young and single-company, so jobs use the standard remote protocol that lets it replace the engine later |
| Streaming | Upgrade to version 2 | State on object storage, recovery in seconds; rival streaming databases' wins are vendor benchmarks against older versions |
| Event log | Keep | Object-storage "diskless" topics were accepted as the direction in March 2026 but not yet shipped; forks that ship them are excluded |
| Serving engine | Replaced after review | See below |
| Policy engine | Keep | Formally verified and analysable |
| Model serving | Add a distributed layer | Only with prompt caches salted per user (R-30, R-31) |

## The serving engine

The foundation-governed search engine first chosen rewrites every
changed document and makes it searchable only after a periodic refresh,
about a second by default. A reading by id is always current, and a
write can wait until it is visible. The alternative, chosen by the
owner, updates data in place and makes it searchable immediately,
combines filters, text, vectors and learned ranking in one query, and
has served about 800,000 queries a second across one large deployment.
Its throughput advantage comes from its own benchmark; its governance is
a single company with its founding customer as a stakeholder. Because
the index is disposable and rebuilt from curated data, leaving it is a
reindex, not a migration, and an exit drill proves that (R-70).

## Language and library couplings

Each language was checked against the libraries it relies on.

- **Best-fit pairings.** The systems language is the home of the query
  engine, the columnar format, the formally verified policy engine and
  the differential-privacy library. The services language is the primary
  kit for the durable workflow engine, the language of the
  transparency-log library, and has a top-tier official kit for the
  agent-tool protocol, where the systems language's is second-tier. The
  agent and pipeline language is native to the data orchestrator.
- **Deliberate mismatches.** The systems language against the table
  format, whose reference implementation lives elsewhere, is contained
  by R-110. The systems-language indexer against the serving engine,
  whose high-throughput feeder is written for another runtime: the
  indexer must compute security labels with the same ontology library as
  the query service, and the serving engine's write interface is plain
  HTTP/2 built for high concurrency, so retries and throttling are
  implemented from its documented behaviour. The event log's only
  official client is for another runtime; the community clients in use
  are established practice.
- **One pairing corrected.** Streaming transformations are written in
  SQL, the stream processor's first-class interface, not in the agent
  and pipeline language, whose interface there lags.
- **One platform issue.** Several foundation-governed components run on
  a runtime that gained hybrid post-quantum key exchange only in its
  release planned for September 2026. The service mesh therefore carries
  post-quantum encryption between services, whatever runtime each uses
  (R-126).

## Watching, not adopting

A distributed engine written natively in the systems language compatible
with the batch engine's protocol; GPU-native SQL engines; diskless
event-log topics; single-file commits in the table format's fourth
version; a streaming analytics component for very large data views; and
a columnar format designed for vectors and multimodal data.
