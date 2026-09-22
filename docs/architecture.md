# Architecture

Urshanabi takes an organisation's isolated data silos, cleans the data
through pipelines, and sets an ontology on the final, curated data,
presented to people and to agents alike. It belongs to an emerging
category whose first instance leads it, and it copies neither that
leader nor Elysium, the owner's earlier prototype: every stage follows
industry-standard practice (`RULES.md`, Purpose).

This document describes the architecture as decided. Each decision is
recorded, with its evidence, in the roadmap item cited beside it.

## 1. The product, stage by stage

```
silos ─▶ raw ─▶ cleaned ─▶ curated ─▶ index ─▶ people and agents
                              ▲                       │
                              └── approved actions ◀──┘
                                        │
                                        └─▶ writeback (off by default)
```

| Stage | What Urshanabi does | Roadmap |
|---|---|---|
| Sourcing | Identifier-keyed snapshot comparison by default, in bounded batches; change capture only where a source needs it; read-only credentials verified | R-53, R-62, R-75 |
| Raw layer | Data lands exactly as it arrived, appended, in bounded snapshots with an append-only change log | R-55, R-56 |
| Cleaned layer | A source that changes shape is refused, not absorbed; transformations are versioned, tested code | R-93, R-97 |
| Curated layer | Identity resolution with approved merges; data is published only after its checks pass on a branch | R-89, R-90 |
| Orchestration | Pipelines declared as assets in Urshanabi's own format and compiled to a foundation-governed data orchestrator | R-96 |
| Contracts | Every layer publishes what it guarantees; the ontology binds only to accepted curated data | R-98 |
| Lineage | Classifications follow data from source to ontology; stopping them needs review | R-76 |
| Ontology | Authored in version control, compiled once by one shared library, signed, exportable to an open specification | R-51, R-74 |
| Serving | A serving engine updated in place and searchable immediately; early binding with a security-freshness watermark; tiered execution | R-70, R-69 |
| People | Uniform denial, no signal from hidden rows, aggregate inference controls, live screens, visible staleness | R-21, R-68, R-71, R-108, R-109 |
| Agents | Every agent through the open agent-tool protocol under the same policy; agent-to-agent delegation within the envelope | R-99, R-115 |
| Actions | Approvals on durable execution; writeback off by default, through an outbox | R-57, R-100 |

**Freshness is declared per object type**: real-time (streamed,
seconds), near-real-time (minutes) or scheduled, each with a target and
an alert; staleness is always shown (R-107, R-108).

## 2. Deployment shape

A control plane holds only metadata. It compiles, validates, signs and
distributes configuration; everything that touches customer data runs in
a **cell**. Configuration is authored in version control, loaded as a
signed, numbered generation pinned once per request, and published
beside the data as an inert manifest. One cell artifact serves hosted,
customer-cloud and disconnected installations. Hosted cells come first,
by the owner's decision (R-52): dedicated cells, one per customer, until
confidential computing makes shared cells safe (R-113). Customer-cloud
and disconnected cells follow once hosted is established.

Cells differ in trust, and the differences are declared, not assumed:

- **Standard cells** run a sidecar-less service mesh for mutual TLS with
  the hybrid post-quantum key exchange (R-126).
- **Confidential cells**, required for pooled hosted cells, run
  workloads in hardware-isolated environments that must prove by
  attestation what they are before receiving keys; encryption terminates
  inside each protected workload (R-113, R-126).
- **Validated-cryptography cells** use validated modules with classical
  key exchange until those modules include the hybrid exchange (R-112).
- **Every cell has a physical-trust ceiling**: hardware an adversary
  could open never holds data above it (R-122).

## 3. Security model

- **One enforcement family.** Only the services that read or write data
  hold data credentials; everything else calls them as the user.
- **Policy as data** in an analysable, formally verified policy
  language, so a reload that widens any role is caught before it applies
  (R-50).
- **Security in the serving plan.** Security values are indexed with
  objects; the security condition runs before any caller condition;
  counts, totals and truncation come only from visible rows (R-68).
- **Authority propagates, never amplifies.** Every internal call carries
  the calling service's identity and the user's short-lived token;
  revocation reaches everywhere within a bound (R-81, R-84).
- **No machine can reach everyone's data**; work is stateless; there is
  no privileged runtime access; every production build is publicly
  verifiable (R-119 to R-121).
- **Agents are envelopes, not principals**: an agent's authority is the
  user's grants intersected with the agent's envelope.
- **Tamper-evident audit** with signed checkpoints published to a
  witness the customer controls (R-88).
- **Inference is controlled**: small aggregates suppressed, differential
  privacy for designated types (R-71).

## 4. Services and language roles

Every service contract exists from the start; the first release ships as
few processes as the security boundaries allow, and each process exists
because of the secrets it alone holds (R-80).

| Service | Language role | Why it is a separate process |
|---|---|---|
| Ingestion | Systems | Holds the source credentials |
| Indexer | Systems | The only writer of the index |
| Query service | Systems | Reads the index; enforces policy on every answer |
| Write service | Systems | The only holder of data-write and writeback credentials |
| Gateway | Services | Holds sessions; the front door and the agent-tool endpoint |
| Workflow worker | Services | Runs approvals, expiries, retries and automations |
| Control plane | Services | Holds the signing keys |
| Model gateway | Services | The only holder of model-provider credentials |
| Audit service | Services | The only writer of the audit log |
| Agent | Agent and pipeline | Handles untrusted model output, so holds no data secrets |
| Pipelines | Agent and pipeline, and SQL | Run under the data orchestrator |
| Interface | Frontend | Knows only the generated client |

The four systems-language services and the ontology library share one
workspace; the services-language services share one module; the agent
and the pipelines are separate projects because their requirements
conflict (`RULES.md` E2).

## 5. How writes meet

- **Edits and source data** reconcile under a rule declared per field —
  the edit persists, the most recent value wins, or a named source has
  priority — and every field shows which source won (R-102).
- **Every change states the version it was based on**; a missing or
  stale version is refused (R-103).
- **One writer per store**, named in `docs/writers.md` and enforced by
  the store itself (R-79); for tables, the table format's own conflict
  check is the final guarantee (R-104), and non-append table writes use
  the format's reference implementation (R-110).
- **The index never goes backwards**: writes carry the source's sequence
  number and older ones are rejected (R-105).
- **Every state-changing request can be retried** without being applied
  twice (R-106).

## 6. Time, languages, scale and the ontology's own form

- **UTC is the one timezone.** Every instant is stored, processed and
  transmitted in UTC; only the edge converts to a reader's local time. A
  future civil time, such as an automation every day at 9am in one city,
  keeps its local time and zone, with its instant recomputed when
  timezone rules change, and calendar dates stay dates (R-128). Tests
  run under three timezones to prove no code reads the host's (R-14).
- **Every word is translatable.** Interface text lives in message
  catalogues, English, Spanish and Portuguese first; the ontology's own
  names are translated once and inherited by every screen and agent;
  search works across languages; agents answer in the reader's language
  (R-129 to R-132).
- **Proven at the largest scale.** The first release is proven against
  the largest enterprises and governments: at least 50 billion objects,
  20,000 people at once and 1,000 workflows a second per cell, with a
  record opened within 100 ms at the 95th percentile. Each objective is
  in `docs/scale.md`, awaiting or citing its load test (R-66).
- **The ontology is authored in TOML**, one small file per type, link,
  action and metric, with one translation catalogue per language, parsed
  only by the shared ontology library. Its open specification is
  `docs/ontology-format.md`, and it exports to the open semantic-model
  specification, never exporting security labels (R-51, R-74).

## 7. Components, by role

| Role | Decision | Governance |
|---|---|---|
| Table format | The open standard, at its third version | Foundation |
| Table catalog | The foundation-governed catalog for that format | Foundation, top-level project |
| Object storage | The standard object-storage interface only: the cloud provider's store when hosted, the foundation-governed distributed store's gateway when self-hosted, a lightweight store for development and CI (R-139) | Foundation in production; the lightweight store is one maintainer's, contained by the interface (owner's decision) |
| Query engine | An embeddable columnar engine in the systems language | Foundation |
| Distributed compute | The foundation-governed batch engine with a native accelerator, jobs written against its standard remote protocol | Foundation |
| Streaming | The foundation-governed stream processor, version 2, state on object storage | Foundation |
| Event log | The foundation-governed log; diskless topics adopted when shipped upstream | Foundation |
| Data orchestrator | The foundation-governed task-based standard with asset-aware scheduling | Foundation (owner's decision) |
| Durable workflows | The engine with the largest proven scale, sized before first use | Foundation, entry tier (owner's decision) |
| Serving engine | Updated in place, immediately searchable, hybrid retrieval and ranking | One company; exit is a reindex (owner's decision) |
| Policy engine | Analysable and formally verified | Foundation |
| Interface | Our own design system on accessible primitives | Ours (owner's decision) |
| Service contracts | A schema-first binary interface language at its version 3 syntax, its standard RPC protocol between services, and that protocol's browser-compatible variant at the gateway, with portability guardrails | Steered by one company, the protocol hosted by a foundation; accepted with guardrails (owner's decision) |
| Ontology authoring | TOML, parsed only by the shared ontology library | An open format (owner's decision) |
| Agent and pipeline language | Version 3.14, built from its foundation's source and pinned exactly; its community's and foundation's own tools for packaging, formatting, linting, typing and auditing (R-10, R-127) | Foundation and community (owner's decision) |
| Container orchestration | Any conformant orchestrator; a security-hardened, government-oriented distribution as reference and production, its lightweight sibling for development and CI | One company, contained by assuming only conformance (owner's decision) |

The named register mapping each role to its product lives outside the
repository (`docs/README.md`).

## 8. Federation

Designed now, built later (Phase 7). **Peering** shares selected objects
between installations over the open dataspace protocol being
standardised through ISO/IEC, under a classification ceiling per
connection (R-117). **Federated search** queries external systems in
place, opt-in per source, with a stated purpose on every request and
identity-proofed users; results are marked unverified until promoted
through the pipelines (R-118).

## 9. Verification

- **Properties**: `conformance/BEHAVIOURS.md`, checked black-box against
  Urshanabi, written test-first (R-07).
- **Layers**: unit, integration, component, contract, end-to-end,
  weighted toward the boundaries; tests discovered, never listed
  (`RULES.md` E4).
- **Agents**: an evaluation harness comparing rates, with adaptive
  attacks run nightly (R-28, R-72).
- **Scale**: objectives from target workloads, each owned by a load test
  (R-66).
- **Compliance**: a control matrix mapping every control to its test
  (R-125).

## 10. Costs and risks

- **Operational surface**: a relational store, an expiring key-value
  store, an event log, a data orchestrator and a durable workflow
  engine, a table catalog, object storage, a serving engine and an
  identity provider. Two orchestration engines is deliberate: data
  orchestration and durable execution are different categories.
- **Custody**: curated data and its index are copies of customer data,
  so encryption, retention and deletion are obligations.
- **Four language roles**, justified by ecosystems.
- **One single-company dependency on the critical path**, the serving
  engine, contained by its licence and a disposable index.
- **Physical attacks** on confidential computing are unsolved
  industry-wide; ceilings limit what any cell may hold.
