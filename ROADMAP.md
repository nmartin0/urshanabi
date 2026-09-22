# Roadmap

The one list of work for Urshanabi (RULES.md §18): what is built, in
what order, and why.

**What Urshanabi is.** A system that takes an organisation's isolated
data silos, cleans the data through pipelines, and sets an ontology on
the final, curated data, presented to people and to agents alike. It
belongs to an emerging category whose first instance leads it. It
copies neither that leader nor Elysium, the owner's earlier prototype;
every stage follows industry-standard practice, with both as a loose
map (RULES.md, Purpose).

**How the roadmap is organised.** Stages in the order the product
depends on them: foundations, then the data, then the ontology and its
serving, then agents, actions, identity and tenancy, hardening, and
federation.
Nothing in a stage starts until what it depends on is done.

## How to read an item

- **Class** — measured against the leading platform's public
  documentation:
  - *foundation*: practice any commercial system needs, whatever its
    category.
  - *parity*: a capability the leader has, which Urshanabi matches.
  - *improvement*: the leader has it; Urshanabi does it better.
  - *new*: neither the leader nor Elysium has it, as far as public
    documentation shows.
  - *open decision*: undecided, waiting for the owner.
- **Status** — *planned*; *in progress*, saying what is done and what
  remains; or *done*, citing its evidence. Paths cited must exist.
- **Outcome** — what Urshanabi does. Always first.
- **Practice and precedent** — the industry practice and precedent
  behind the outcome, described without names (RULES.md H1).
- **Learned from Elysium** — evidence from the prototype, tagged by
  how it is known.
- **Owner** — a decision the owner must make (RULES.md H8).
- **Decision** — how an owner decision was settled, and when.
- **Done when** — a test that fails if the outcome is removed.

## Evidence tags

- **[measured]** — measured against Elysium at commit `f4ea94e` or
  later.
- **[code]** — confirmed by reading Elysium's code.
- **[docs]** — recorded in Elysium's documents, not re-checked.
- **fixed late** — Elysium had the flaw and later fixed it.
- **[precedent]** — established practice elsewhere, described by what
  it does.

Item numbers are stable because commits and documents refer to them;
they are not an order. A retired item keeps its number and is struck
through.

---

# Phase 0 — Foundations

Before any stage holds real data. Engineering and security practice
any commercial system needs, whatever its category.

### R-01 Continuous integration from the first commit

- **Class:** foundation
- **Status:** done — every push runs `script/cibuild`, and every check
  is broken on purpose by a self-test (`script/cibuild`,
  `script/check-scripts`).
- **Outcome:** every push runs every lint, type, dead-code, boundary,
  lockfile and naming check, and every test suite, on a clean checkout.
- **Learned from Elysium:** no CI of any kind. Its lint script called
  itself the CI equivalent, and nothing ran it on push. [measured]
- **Done when:** a deliberately broken commit fails the pipeline, for
  each check in turn.

### R-02 Tests declare their prerequisites

- **Class:** foundation
- **Status:** planned
- **Outcome:** a test needing data declares a fixture that builds it,
  or skips with the exact command that would. The fast tier needs
  nothing outside the repository.
- **Learned from Elysium:** on a fresh clone, 15 of 1,849 unit tests
  failed with `assert 0 > 0`. They depended on seeded and synced data
  nobody had declared. [measured]
- **Done when:** a CI job runs the fast tier on a checkout with no
  generated data and every test passes or skips with a message.

### R-03 Browser tests in CI, not beside it

- **Class:** foundation
- **Status:** planned
- **Outcome:** real-browser tests run in CI against the built
  interface for anything layout, hit-testing or cascade dependent.
- **Learned from Elysium:** the simulated document that unit tests run
  against computes no layout, so 20 passing tests coexisted with a
  visibly broken screen. Real-browser tests existed and were run by
  hand. [docs]
- **Done when:** a change that covers one element with another fails
  CI.

### R-04 One end-to-end test per wire

- **Class:** foundation
- **Status:** in progress — the gateway's tests mock the query
  service, and the walking skeleton's journey drives the real one
  through its real entry point, as binaries and as container images
  (`e2e/script/test-integration`); remaining: the same for every handler
  mocked from now on.
- **Outcome:** every handler a test mocks has one test driving the
  real one through its real entry point.
- **Learned from Elysium:** three shipped commits crashed on a user's
  first real run, each behind a green suite; three more bugs hid behind
  mocked callbacks. [docs]
- **Done when:** replacing a real handler with a stub fails at least
  one test.

### R-05 Versioned migrations before any store persists

- **Class:** foundation
- **Status:** planned
- **Outcome:** every store carries a schema version and ordered
  migrations from its first table.
- **Learned from Elysium:** stores create their tables if absent. One
  store has an ad hoc, unversioned migration function; none records a
  schema version. [code]
- **Done when:** a store created at version 1 opens at version 2 with
  its data intact, and a downgrade is refused.

### R-06 Service contracts with breaking-change detection

- **Class:** foundation
- **Status:** in progress — breaking changes fail CI, and a test fails
  if the gateway logs an unsafe error argument or a producer's free-text
  message; each was seen to fail on a planted leak
  (`contracts/script/test`,
  `services/gateway/internal/failure/failure_test.go`); remaining: a
  decimal type that crosses every contract exactly, never as a
  floating-point number, which the ontology's decimal properties
  require.
- **Outcome:** internal calls defined in a schema-first binary
  interface language, with breaking-change checks in CI; the public HTTP
  description generated from it; the changelog and audit streams
  described with the standard event-API description. Every error
  argument is declared safe or unsafe, and only safe arguments are ever
  logged.
- **Practice and precedent:** one vendor's open-source contract
  toolchain was built to scale its own microservices, but its
  compatibility checker was never released. The canonical choices are a
  schema-first binary interface language for internal calls, the
  industry-standard HTTP API description (governed by a neutral
  foundation) for the public edge, and the standard event-API
  description for streams. That vendor's error model separates safe
  arguments from unsafe ones, so the schema itself declares what may
  reach a log. [precedent]
- **Learned from Elysium:** one process; contracts were function
  signatures.
- **Decision:** the schema-first binary interface language at its
  version 3 syntax, its standard RPC protocol between services, and that
  protocol's browser-compatible variant at the gateway (owner,
  2026-09-21), chosen over a fully foundation-governed pair of HTTP and
  event description standards. The language is steered by one company,
  and the RPC protocol, though hosted by a foundation, largely by the
  same one; the toolchain CLI and the frontend runtime come from a
  second company. Accepted because the format is open, permissively
  licensed and independently implemented, including in the
  implementations we use; the alternative's specifications are
  foundation-governed, but its tools rest largely on volunteers, it
  lacks typed binary contracts and streaming, and bulk data (R-114)
  needs the RPC protocol regardless. Portability guardrails, each
  enforced by `contracts/script/test`: version 3 syntax only; no imports
  beyond the language's own well-known types; no hosted registry and no
  remote generation plugins; the strictest breaking-change category.
  Each single-company piece has a register entry and an exit. The public
  HTTP description is chosen with the gateway, and the event
  descriptions with the first stream.
- **Done when:** removing a field from a contract fails CI, and
  logging an unsafe error argument fails a test.

### R-07 A conformance suite, built test-first with Urshanabi

- **Class:** foundation
- **Status:** planned
- **Outcome:** `conformance/BEHAVIOURS.md` implemented as black-box
  property tests against Urshanabi, through a thin driver that
  translates each property into Urshanabi's own interface. The suite is
  written test-first: each property's test lands with the feature that
  meets it, starting with the walking skeleton. Elysium is the evidence
  that each property is achievable and testable, not a target, so no
  driver is written for it.
- **Decision:** the suite is built against Urshanabi only (owner,
  2026-09-21). A driver for Elysium would be throwaway work against a
  prototype interface Urshanabi deliberately does not keep, would fail
  in its known flaws and need each recorded as a divergence, and would
  risk shaping the tests around Elysium's shapes. Each test is proved
  instead by a planted fault in a test build of Urshanabi.
- **Done when:** every property passes against Urshanabi, and for each
  property a test build with a planted fault — such as one
  uniform-denial response flipped — fails the suite.

### R-08 A request id on every route from the first route

- **Class:** foundation
- **Status:** in progress — the gateway gives every request an id
  before any route sees it, including the router's own refusals, keeps a
  caller's valid id, refuses ids that could forge log lines, and passes
  it to the query service; both log it
  (`services/gateway/internal/requestid/requestid.go`,
  `e2e/script/test-integration`); remaining: a test enumerating every
  route as routes are added, the id in audit records, and emission to
  the vendor-neutral telemetry standard.
- **Outcome:** every request carries an id from the gateway through
  every service, into every audit record and log line, with traces,
  metrics and logs emitted to a vendor-neutral telemetry standard.
- **Practice and precedent:** diagnosing incidents by logging into
  individual servers is the recognised sign that observability was left
  too late. [precedent]
- **Learned from Elysium:** request correlation existed and one route
  of 39 created an id, so most reads wrote untracked audit lines. [docs,
  fixed late]
- **Done when:** a test enumerates every route and fails if any
  response lacks the id.

### R-09 Integrity checks that cannot be stripped

- **Class:** foundation
- **Status:** planned
- **Outcome:** integrity checks are ordinary code paths that return or
  raise explicitly, in every language role.
- **Learned from Elysium:** data-integrity checks were language
  assertions, which an optimisation flag removes; it had to refuse to
  start with that flag set. [code]
- **Done when:** a release build with every optimisation enabled still
  rejects a corrupt record.

### R-10 Reproducible installs, exactly

- **Class:** foundation
- **Status:** in progress — tools are pinned and checksum-verified;
  the systems-language workspace builds only from its lockfile, and the
  services-language gates refuse a module that differs from its recorded
  hash or an untidy manifest (`script/tools.lock`, `Cargo.lock`,
  `script/services-gates`); remaining: the other two languages'
  lockfiles.
- **Outcome:** every language role installs from a hashed lockfile; CI
  fails when a lockfile drifts from its manifest.
- **Learned from Elysium:** dependencies were bounded rather than
  locked until late, and a comment still says locking was deferred after
  lockfiles shipped. [measured]
- **Done when:** editing a manifest without regenerating its lockfile
  fails CI.

### R-11 Generated state is never tracked

- **Class:** foundation
- **Status:** done — every component's generator runs on every push,
  and the check fails unless version control then sees exactly what it
  saw before, every untracked file listed on its own; the systems
  language's build directory must be ignored too
  (`script/check-generated`).
- **Outcome:** ignore rules cover every generated path before the
  first generator exists.
- **Learned from Elysium:** generated databases and lake files were
  committed twice, leaving the tree permanently dirty and blocking
  patches for several rounds. [docs, fixed late]
- **Done when:** running every generator leaves version control
  reporting a clean tree.

### R-12 A development environment that survives a reboot

- **Class:** foundation
- **Status:** in progress — one command builds both services' images
  and runs the stack, replacing any running copy; the journey stops the
  stack, confirms it is down, restores it with the one command and
  requires the full answer (`script/server`,
  `e2e/script/test-integration`); remaining: each store's data under the
  user's own directory, as stores arrive.
- **Outcome:** one command builds the whole local stack, with data
  under the user's own directory.
- **Learned from Elysium:** development data lived in a directory
  cleared on reboot, and was lost twice. [docs]
- **Done when:** a reboot followed by one command restores a working
  stack.

### R-135 Every service ships as a verifiable container image

- **Class:** foundation
- **Status:** in progress — both services build into images holding
  their static binary alone, on no base, run unprivileged, stamped with
  the commit's time; the journey proves each holds one file, runs as a
  non-root user, and rebuilds without the cache to the same image, and
  each check was seen to fail on a planted fault
  (`services/query/Containerfile`, `services/gateway/Containerfile`,
  `e2e/script/test-integration`); remaining: a bill of materials and
  signature for every image, publication to the transparency log, and a
  cluster that refuses an image lacking them.
- **Outcome:** each service's container build lives in its own
  directory (`RULES.md` E2) and produces an image in the open container
  image format: built in stages, on a minimal base pinned by digest,
  running as a non-root user with no shell (R-120). Builds are
  reproducible, stamped with the source commit's time rather than the
  build time (R-121); every image carries a bill of materials and a
  signature, published to the transparency log. Clusters run images with
  the standard foundation-governed container runtime, so no
  image-building tool is installed on servers.
- **Practice and precedent:** every microservice repository inspected
  keeps each service's container build beside its source. Images in the
  open format run on any conformant runtime, and reproducible, signed
  images with bills of materials are standard supply-chain practice.
  [precedent]
- **Learned from Elysium:** shipped no container images. [code]
- **Done when:** two builds of the same commit produce identical
  images, and an image without a valid signature and bill of materials
  is refused by the cluster.

### R-136 Cells run on any conformant orchestrator, tested on one reference

- **Class:** foundation
- **Status:** planned
- **Outcome:** Urshanabi assumes only a certified, conformant
  container orchestrator, so it runs on any conformant distribution.
  Each cell is one cluster, dedicated to one customer while hosted
  (R-52); the metadata-only control plane runs in a cluster of its own.
  Traffic enters through the orchestrator's standard gateway interface.
  CI and developer machines run a lightweight distribution, and the
  product is released only after passing on the reference distribution.
- **Practice and precedent:** a security-hardened distribution of the
  orchestrator ships defaults that pass the CIS benchmark, enables FIPS
  140-2 validated cryptography, supports SELinux and air-gapped
  installation, and stays closely aligned with upstream; its lightweight
  sibling, a foundation sandbox project since 2020, suits edge sites,
  development and CI. The most widely used ingress controller retired in
  March 2026, and distributions moved to the standard gateway interface.
  [precedent]
- **Learned from Elysium:** ran as a single process, with no
  orchestrator. [code]
- **Decision:** the reference and production distribution is the
  security-hardened, government-oriented one; the lightweight sibling
  runs on developer machines and in CI (owner, 2026-09-21). It is one
  company's, contained because the product assumes only conformance, so
  changing distribution is an operations change, not a product change.
- **Done when:** the end-to-end suite passes on the reference
  distribution in its hardened profile and on the lightweight one, and a
  release that fails either is not published.

### R-13 Fixtures with volume

- **Class:** foundation
- **Status:** planned
- **Outcome:** a small fixture for speed and a volume fixture for
  paging, truncation and memory ceilings, both in CI.
- **Learned from Elysium:** fixtures held 4 customers and 7
  transactions, so paging and cut-off behaviour could not be seen.
  [measured]
- **Done when:** the volume fixture exercises every paging and
  truncation path in `BEHAVIOURS.md`.

### R-14 Tests prove no code depends on the host's timezone

- **Class:** foundation
- **Status:** in progress — the systems and services languages' tests
  run under UTC, São Paulo and New York, and self-tests prove a
  deliberate local-time read fails in both non-UTC zones only, in each
  (`script/systems-gates`, `script/services-gates`,
  `script/check-scripts`); remaining: the other two languages.
- **Outcome:** date-sensitive tests run three times: under UTC, under
  a zone behind UTC, and under a zone observing daylight saving. The
  zone is set by the test scripts themselves, so a local run and CI run
  behave the same (`RULES.md` E5). Correct UTC-only code passes in every
  zone; code that reads the host's local time fails.
- **Practice and precedent:** code that accidentally uses the
  machine's local time looks correct on a machine set to UTC, because
  the two coincide there. One team's daily reports were shifted by five
  hours for a year because CI and staging ran in UTC while production
  did not. Established projects run their tests under non-UTC zones,
  including a daylight-saving one, and set the zone in the test runner
  rather than only in CI so local runs are covered too. [precedent]
- **Learned from Elysium:** its method notes that a UTC container
  hides every timezone bug. [docs]
- **Done when:** a deliberate local-time conversion fails under both
  non-UTC zones, and passes under UTC alone.

### R-128 UTC is the one timezone

- **Class:** foundation
- **Status:** in progress — the rule is recorded, contracts use the
  UTC timestamp type, and both services' logs are UTC, proved under
  three zones (`services/query/src/logging.rs`,
  `services/gateway/internal/logging/logging.go`); remaining: the
  three-zone tests in the other languages, and civil-time scheduling.
- **Outcome:** every instant is stored, processed, compared and
  transmitted in UTC: in contracts as the standard UTC timestamp type,
  in text as an RFC 3339 timestamp ending in Z. Only the edge converts
  to local time, for the person reading it, using their chosen timezone.
  No code asks the host for its timezone. Two exceptions are stored as
  what they are: a future civil time, such as an automation every day at
  9am in one city or a deadline set in a jurisdiction's local time, is
  kept as its local date-time and timezone identifier, with its UTC
  instant derived and recomputed whenever the timezone database changes;
  and a calendar date is a date, not a timestamp.
- **Practice and precedent:** machine-generated timestamps are
  instants, and storing them in UTC is sound; but converting future
  local times to UTC goes wrong when governments change timezone rules,
  so the local time and zone are kept as the truth and the instant
  derived from them. An insurer that stored coverage end dates as
  timestamps found there was no single instant at which coverage ends
  everywhere. [precedent]
- **Learned from Elysium:** timestamps were UTC in its own records but
  untested against other zones. [code]
- **Done when:** every timestamp crossing a contract is UTC, and an
  automation scheduled for 9am local time still runs at 9am local time
  after a simulated timezone-rule change.

### R-15 Documentation that cannot silently drift

- **Class:** foundation
- **Status:** in progress — every cited item, property and rule is
  checked (`script/check-docs`); remaining: documented limitations
  linked to the tests that cite them.
- **Outcome:** one roadmap (this file). Design documents hold
  reasoning, not work. A known limitation names what would close it, or
  cites a test asserting it still holds.
- **Learned from Elysium:** 648 KB of top-level markdown — more than
  twice its source — across six overlapping roadmaps. Six documented
  claims were found stale, including a "real hole" already closed and
  two limitations already lifted. 54% of the backend is prose.
  [measured]
- **Done when:** fixing a documented limitation without updating its
  entry fails the test that cites it.

### R-16 The naming rule is enforced

- **Class:** foundation
- **Status:** done — the name check runs on every push
  (`script/check-names`).
- **Outcome:** CI scans everything written for the repository for
  company, product and language names.
- **Learned from Elysium:** its documents and comments name a
  competitor and many products throughout, and instruct quoting them.
- **Done when:** adding a product name to a comment fails CI.

### R-17 Repository visibility and licence

- **Class:** foundation
- **Status:** done — an explicit proprietary notice states the terms
  (`LICENSE`).
- **Outcome:** private from the first commit, or licensed
  deliberately. An owner decision.
- **Learned from Elysium:** its licence describes unpublished
  proprietary source; its repository is publicly readable. [measured]
- **Audit finding:** Urshanabi's own repository answers
  unauthenticated requests today, so this roadmap and the security
  design in it are publicly readable, and there is no licence file.
  [measured]
- **Decision:** public for the time being, with an explicit
  proprietary notice: all rights reserved to Nicholas Martin, and no
  licence granted except by separate written agreement (owner,
  2026-09-21). The notice replaces the implicit default of an unlicensed
  public repository.
- **Done when:** an unauthenticated request for the repository is
  refused, or a licence file states the terms deliberately chosen.

### R-18 No shared in-process state

- **Class:** foundation
- **Status:** planned
- **Outcome:** every service is stateless; state lives in the stores
  built for it.
- **Learned from Elysium:** per-object locks, the pending-write
  working copy, the loaded configuration, lockouts, rate limits, ten
  single-writer embedded databases, a local audit file and a catalog in
  an embedded database — all tied it to one process, so a deploy was
  downtime. [code]
- **Done when:** the conformance suite passes against two replicas of
  every service with requests alternating between them.

### R-33 Nothing built that nothing uses

- **Class:** foundation
- **Status:** planned
- **Outcome:** a component ships with its caller, or not at all.
- **Learned from Elysium:** a memory guard was built and tested but
  never wired into the agent. [code]
- **Done when:** the dead-code check fails on an uncalled public
  component.

### R-66 Stated scale objectives, each with a load test

- **Class:** foundation
- **Status:** in progress — the objectives are written and checked
  (`docs/scale.md`, `script/check-scale`); remaining: a load test for
  each objective, run on a schedule.
- **Outcome:** written objectives for objects per type, query latency
  at the 95th percentile, pipeline freshness, rows ingested per hour and
  cell throughput, derived from the workloads of the customers Urshanabi
  targets. The leading platform's published limits are a reference
  point, not a target.
- **Practice and precedent:** an established ontology platform's
  current documentation caps a link traversal's result at 10 million
  objects, loads at most 100,000 into memory per call, and lets one
  action edit at most 10,000 objects across 50 types; it sets no fixed
  limit on objects per type, the index's disk space being the ceiling.
  Established response-time thresholds put instant at about a tenth of a
  second, uninterrupted flow at a second, and the limit of attention at
  ten seconds. [precedent]
- **Learned from Elysium:** no scale target anywhere. Before this
  item, the largest number in this roadmap was a 10-million-row sync.
  [measured]
- **Decision:** the first release is proven against the largest
  enterprises and governments (owner, 2026-09-21). The objectives, per
  cell, are in `docs/scale.md`, kept honest by `script/check-scale`: at
  least 50 billion objects, 20,000 people at once, 2 billion rows
  ingested an hour and 1,000 workflows a second at peak; a record opened
  within 100 ms, a filtered search within 500 ms and a billion-object
  aggregate within 2 s, each at the 95th percentile. The figures are
  proposed from precedent and revised only with load-test evidence.
- **Done when:** the objectives are written in Phase 0, each later
  item that meets one adds its load test, and a change that breaks a
  tested objective fails a scheduled load run.

### R-78 Consumer-driven contract tests

- **Class:** foundation
- **Status:** planned
- **Outcome:** each consumer publishes what it relies on; each
  producer's CI verifies against every consumer's expectations.
- **Practice and precedent:** contract testing, where producers and
  consumers verify compatibility independently, is the established way
  to catch breaking changes without full integration suites. [precedent]
- **Learned from Elysium:** one process, so there were no consumers to
  verify against.
- **Done when:** a producer change that breaks one consumer's recorded
  expectation fails the producer's build.

### R-79 One writer per store

- **Class:** foundation
- **Status:** in progress — the map is written and checked
  (`docs/writers.md`, `script/check-writers`); remaining: enforcement by
  each store once it exists.
- **Outcome:** a written map, `docs/writers.md`, naming the single
  writer of every store, enforced by credentials: only the owner holds
  write access, and the store refuses everyone else. One writer is one
  owning identity: an owner may write through more than one mechanism,
  but always as itself — ingestion appends raw data and runs raw-table
  maintenance through the reference implementation (R-110), both as
  ingestion. Tables are enforced by per-table grants in the table
  catalog, with storage credentials scoped to each table's prefix, which
  requires each table's files in a directory of their own matching the
  namespace hierarchy.
- **Practice and precedent:** shared databases are the most cited
  cause of distributed monoliths, though sharing inside one repository
  with changes shipped together has worked; the rule underneath both is
  clear ownership. [precedent]
- **Learned from Elysium:** one process owned every store, so
  ownership was never written down. [code]
- **Done when:** a service attempting to write a store it does not own
  is refused by the store, not by convention.

### R-80 Decompose on evidence

- **Class:** foundation
- **Status:** planned
- **Outcome:** every service contract exists from Phase 0, but the
  first release ships as few processes as the security boundaries allow,
  and each process exists because of the secrets it alone holds:
  ingestion (source credentials), indexer (index write), query service
  (index read), write service (data write and writeback), model gateway
  (model-provider keys), control plane (signing keys), audit service
  (the audit log), gateway (sessions), and the agent, which holds no
  data secrets because it handles untrusted model output. Approvals and
  automations share one workflow worker because they share one engine
  and the same secrets. Other splits wait for a measurement.
- **Practice and precedent:** over-splitting before a domain justifies
  it adds coordination cost with no benefit; the standard advice is to
  start modular and split when pressure proves it necessary, and to
  split by capability rather than by technical layer. [precedent]
- **Learned from Elysium:** a single process throughout. [code]
- **Done when:** each separately deployed service cites the
  measurement or security boundary that justified it.

### R-85 Every dependency passes the selection record

- **Class:** foundation
- **Status:** planned
- **Outcome:** a dependency register records licence, governance, exit
  interface, disconnection and supply chain for every dependency. CI
  checks licences against an allowlist.
- **Practice and precedent:** see `RULES.md` H4a — an archived object
  store, a mesh whose stable releases went vendor-only, and a streaming
  log that proved source-available. [precedent]
- **Learned from Elysium:** dependencies were chosen well but case by
  case; its own plans still name an object store whose open-source
  edition has since been archived. [code]
- **Done when:** adding a dependency with a source-available licence,
  or with no register entry, fails CI.

### R-86 The build pipeline is hardened against supply-chain attack

- **Class:** foundation
- **Status:** done — the CI definition is pinned and scanned at the
  strictest level, and tools are pinned and checksum-verified
  (`script/check-workflows`, `script/bootstrap`).
- **Outcome:** every CI action and tool is pinned by digest; workflows
  triggered by untrusted input get no secrets; publishing uses
  short-lived federated credentials, never stored tokens; release
  artifacts are rebuilt and compared with source; installs are audited
  for start-up hooks.
- **Practice and precedent:** in 2026 a widely used model-gateway
  library shipped credential-stealing releases after its publishing
  credentials were taken through a compromised vulnerability scanner in
  its CI; the malicious code ran at interpreter start-up without being
  imported, and the source repository was never changed. [precedent]
- **Learned from Elysium:** no CI, so no pipeline to attack yet.
  [measured]
- **Done when:** a pipeline change that uses an unpinned tool, or
  exposes a secret to an untrusted trigger, fails a policy check.

### R-94 Nothing leaves a cell unless enabled

- **Class:** foundation
- **Status:** planned
- **Outcome:** no component sends anything outside a cell unless an
  administrator enables it; dependency telemetry is switched off in
  every shipped configuration.
- **Practice and precedent:** an incremental streaming database
  considered for Urshanabi sends anonymous usage statistics by default.
  [precedent]
- **Learned from Elysium:** no telemetry of its own. [code]
- **Done when:** a cell running under an egress capture sends nothing
  outward over a full test cycle.

### R-95 The interface is our own design system

- **Class:** improvement
- **Status:** planned
- **Outcome:** Urshanabi's own design system: layered design tokens
  (primitive, semantic, component) and its own components, built on
  unstyled, accessible behaviour primitives installed as ordinary
  dependencies, never copied in (RULES.md H2). The aim is better, not
  only different: the current accessibility standard at AA level, and
  stated density and speed targets for the data-heavy screens.
- **Learned from Elysium:** its interface is built on a component
  library published by the company whose product is the precedent.
  [code]
- **Practice and precedent:** a product with its own design language
  builds it as layered tokens on accessible primitives rather than
  overriding a third-party look; teams that adopted a styled library
  early report spending much of their time fighting it. The common
  shortcut of copying component source into a project is excluded by
  RULES.md H2. [precedent]
- **Decision:** our own creation, but better (owner, 2026-09-21).
- **Done when:** no screen imports the previous library; an
  accessibility audit passes at AA; the density and speed targets are
  met on the data-heavy screens.
### R-129 Every word on screen is translatable

- **Class:** parity
- **Status:** planned
- **Outcome:** every string the interface shows lives in a message
  catalogue, written in the standard message syntax with plural and
  gender rules, never in code. English, Spanish and Portuguese ship
  first; adding a language is adding a catalogue. Dates, numbers,
  plurals and sort order come from the standard Unicode locale data
  built into every browser, applied to UTC instants in each person's
  timezone (R-128). Errors reach the interface as a stable reason and
  safe arguments, and the interface renders the sentence in the reader's
  language. Layouts use direction-neutral styling, so right-to-left
  languages mirror correctly. A person's language comes from their own
  setting, then their browser.
- **Practice and precedent:** the Unicode locale data is used by all
  major browsers and phones. The newer version of the standard message
  syntax became a stable specification in 2025, but its implementations
  remain previews and production use is rare, so messages use the
  established syntax, which the new version's data model can represent.
  The leading platform lets application builders translate fixed text in
  its application builder, but only listed kinds of text. A
  pseudo-language that lengthens and marks every string is a standard
  test for hard-coded text and truncation. [precedent]
- **Learned from Elysium:** every string was English, written inline.
  [code]
- **Done when:** switching to Spanish or Portuguese leaves no
  untranslated string, verified in CI by the pseudo-language, and every
  error a service returns renders in the reader's language.

### R-121 Every production build is publicly verifiable

- **Class:** foundation
- **Status:** in progress — the gateway answers with its own build
  identity and the query service's; both report the commit they were
  built from, say so when built from uncommitted changes, and refuse a
  release built from them; the walking skeleton's journey proves the two
  agree (`e2e/script/test-integration`); remaining: published
  fingerprints, and attestation checked against them.
- **Outcome:** the fingerprint of every production build is published
  to the tamper-evident log (R-88), and every component's attestation is
  checked against it, so a customer can verify exactly which code is
  running on their data.
- **Practice and precedent:** the strongest precedent for running on
  machines the operator does not fully trust publishes every production
  build for inspection, keeps a transparency log, and makes verifiable
  transparency one of five core requirements. [precedent]
- **Learned from Elysium:** published nothing about what ran. [code]
- **Done when:** a customer can match a running component's
  attestation to a published build fingerprint, and an unpublished build
  fails attestation.

### R-125 A compliance control matrix from the first commit

- **Class:** foundation
- **Status:** done — every implemented control cites a passing check
  and self-test (`docs/compliance.md`, `script/check-compliance`).
- **Outcome:** a maintained matrix maps every control Urshanabi
  implements to the baselines its buyers use: the federal control
  catalogue (SP 800-53 Rev. 5), the zero-trust architecture guidance (SP
  800-207 and SP 800-207A), the microservice security series (SP
  800-204), the secure software development framework (SP 800-218), and
  the US health-data law's Security Rule (45 CFR Part 164, Subpart C),
  designed to its proposed revision rather than the current text. Each
  control names the roadmap item and the test that proves it.
- **Practice and precedent:** the proposed revision of the health-data
  Security Rule, not final as of mid-2026, makes encryption at rest and
  in transit and multi-factor authentication mandatory, with scans every
  six months, yearly penetration tests, 72-hour restoration and a yearly
  asset inventory and network map; its core mandates are expected to
  survive. The zero-trust guidance for cloud-native applications moves
  security from network location to identity, enforced through gateways,
  proxies and workload identities. [precedent]
- **Learned from Elysium:** had no control mapping. [code]
- **Done when:** every control in the matrix cites a test that passes
  in CI, and a control without a test fails the build.

### R-127 The engineering standards are enforced from the first commit

- **Class:** foundation
- **Status:** in progress — repository gates, boundaries and the
  scripts' own tests run in CI, and the systems and services languages'
  gates run on their first components, each broken on purpose by a
  self-test (`script/check-scripts`, `script/systems-gates`,
  `script/services-gates`); remaining: the other two languages' gates,
  vulnerability and licence gates for the services language, and the
  dependency-rung check.
- **Outcome:** `RULES.md` part three is enforced in CI from the first
  commit: every dependency records its rung on the standard-library
  ladder (E1); every component is its own project, with dependencies
  shared per language and isolation enforced at code boundaries (E2);
  every language passes the six gates (E3); tests are discovered,
  never listed, and written in the layers of E4; and every project and
  the repository root share the same delegating scripts (E5).
- **Practice and precedent:** repositories are kept whole when code is
  shared and changes coordinated; heavyweight multi-language build
  systems pay off only at very large scale. The established
  microservice testing strategy has unit, integration, component,
  contract and end-to-end layers, and current practice weights
  integration most because faults cluster at service boundaries.
  [precedent]
- **Learned from Elysium:** its lint script had five gates and nothing
  ran it on push. One gate's failure was silently discarded: the
  lock-file check set a variable the script never read, so the script
  exited successfully while that gate failed. A test that deliberately
  breaks each gate would have caught it. [measured]
- **Done when:** a deliberate violation of each gate, in each language
  present, fails CI; a dependency without a recorded rung fails the
  build; a service importing another service's code fails
  `script/check-boundaries`; and the top-level script fails against a
  planted failing component, a component with no script, and an empty
  tree.

---

# Phase 1 — Data foundation: sourcing, pipelines, storage

Silos become raw, then cleaned, then curated data. Everything later
reads curated data, so this stage comes first. Order of work: secrets
and encryption, then orchestration and layer contracts, then
ingestion, then transformation and publication, then history, curation
and lineage.

### R-92 Secrets are references, never values

- **Class:** foundation
- **Status:** planned
- **Outcome:** every credential in configuration is a reference to the
  secret store, never a value; a missing reference refuses to load and
  names the field; secrets never reach logs, errors, bundles or the
  published manifest.
- **Learned from Elysium:** configuration may carry an
  environment-variable reference instead of a password, and a missing
  variable refuses to load, naming the field. [code]
- **Done when:** a bundle or manifest containing a secret value fails
  validation, and a missing reference refuses to load by name.

### R-91 Encryption in transit and at rest, keyed per tenant

- **Class:** foundation
- **Status:** planned
- **Outcome:** every connection encrypted; every store encrypted at
  rest with keys held in a key service; per-tenant keys once tenancy is
  decided (R-52), so one tenant's data can be destroyed by destroying
  its key.
- **Learned from Elysium:** loopback-only transport; files at rest
  unencrypted. [code]
- **Done when:** a store's files read without the key service yield no
  plaintext, and rotating a key re-encrypts without downtime.

### R-96 Pipelines run under an asset-aware data orchestrator

- **Class:** parity
- **Status:** planned
- **Outcome:** pipelines are declared in Urshanabi's own format, as
  data assets with their dependencies, freshness, checks and lineage,
  and compiled to a foundation-governed data orchestrator that keeps
  them reconciled; events trigger runs and the orchestrator keeps the
  record. Approvals and actions stay on durable execution (R-57); one
  engine does not serve both.
- **Practice and precedent:** data-pipeline orchestration and durable
  execution are distinct categories that teams routinely confuse. Data
  orchestrators model assets, lineage, freshness and dataset-aware
  scheduling; durable execution is for long-running business processes
  and human-in-the-loop flows. The common pattern is that events trigger
  and an orchestrator tracks, using a shared lineage event format.
  [precedent]
- **Learned from Elysium:** no scheduler; its roadmap found scheduling
  and running several workers to be one decision. [docs]
- **Decision:** the foundation-governed, task-based standard that has
  gained asset-aware scheduling (owner, 2026-09-21). Chosen for
  governance under RULES.md H4a: no company can take it private or
  change its terms. The alternative fitted the asset model more
  naturally but is owned by one company, acquired in 2026. Because
  customers never touch the orchestrator and our own format compiles to
  it, the design-fit gap is ours to absorb, and switching later stays a
  migration, not a rebuild.
- **Done when:** a failed upstream load marks every downstream asset
  stale, visibly, without that tracking having been built by hand.

### R-98 Each layer has an owner and a quality contract

- **Class:** parity
- **Status:** planned
- **Outcome:** the raw, cleaned and curated layers each publish what
  they guarantee — schema, freshness, completeness — and name who owns
  it. The ontology binds only to curated data whose contract it accepts.
- **Practice and precedent:** practitioners warn that teams get the
  layer names right and the ownership wrong: the layers are a contract
  about who is responsible for quality at each stage, and ingestion
  should never write the cleaned layer directly. [precedent]
- **Learned from Elysium:** raw, change-log and current layers with
  drift refusal, but their guarantees were implicit. [code]
- **Done when:** an ontology type bound to a table without an accepted
  contract fails validation, and a contract breach blocks publication.

### R-62 Read credentials are checked, not assumed

- **Class:** improvement
- **Status:** planned
- **Outcome:** at startup, each read connection tests whether its
  credential can write, and refuses or warns by policy.
- **Learned from Elysium:** documented that sources should use
  read-only credentials, and could not check it. [docs]
- **Done when:** a read connection holding write privilege is flagged.

### R-53 Ingestion within a fixed memory ceiling

- **Class:** parity
- **Status:** planned
- **Outcome:** ingestion streams in bounded batches, and refuses a
  table it cannot hold rather than being killed.
- **Learned from Elysium:** a sync held the whole table four times
  over, peaking at about 1.2 KB per row; 10 million rows would fail.
  [docs]
- **Done when:** a 10-million-row table syncs under a fixed memory
  limit in CI, and the scale objective R-66 sets for sync is met by a
  scheduled load run.

### R-55 Bounded raw storage from the first sync

- **Class:** parity
- **Status:** planned
- **Outcome:** the raw layer keeps two snapshots from day one; history
  lives in a changelog that grows only with real changes.
- **Learned from Elysium:** copy-on-write grew a 20,000-row table with
  one change from 177 KB to 839 KB over five syncs. [docs, fixed late]
- **Done when:** fifty unchanged syncs leave storage flat.

### R-93 A source that changes shape is refused, not absorbed

- **Class:** parity
- **Status:** planned
- **Outcome:** a column whose values no longer match its declared type
  produces a drift report naming the column, the type and an example,
  and ingestion of that source is refused rather than silently coerced;
  severity is tiered, and administrators are notified (R-58).
- **Learned from Elysium:** a column whose values no longer match the
  declared type produces a drift report naming the column, the type and
  an example, and the sync is refused rather than silently coerced.
  [code]
- **Done when:** a source column changing type refuses the sync with a
  report naming the column.

### R-104 One writer per table, three layers deep

- **Class:** foundation
- **Status:** planned
- **Outcome:** each table has exactly one owner (R-79): ingestion for
  every raw table, and one owning pipeline for each cleaned and curated
  table. The orchestrator runs at most one run of a pipeline at a time,
  and the table format's own conflict check remains the guarantee, so
  correctness never depends on the first two. Commits are batched rather
  than frequent.
- **Practice and precedent:** the table format commits by atomically
  swapping a pointer and, on collision, re-validating and retrying;
  practitioners report retries becoming the bottleneck with many writers
  committing every few seconds. The orchestrator can cap runs per
  pipeline and access to shared resources, but that cap has been
  silently not honoured in past versions and is enforced only by its
  scheduler. Other-language implementations of the table format have
  shipped false conflicts and missing retries. [precedent]
- **Learned from Elysium:** a single process wrote every table, so
  collisions could not occur. [code]
- **Done when:** two runs forced to overlap on one table produce one
  success and one clean retry or refusal, never a lost change; the
  systems-language table library's retry and validation are verified
  before any pipeline depends on them.

### R-110 Non-append table writes use the reference implementation

- **Class:** foundation
- **Status:** planned
- **Outcome:** raw ingestion appends through the systems-language
  table library, and all reads use it; every overwrite, delete,
  row-level update, merge and compaction — including retention deletions
  (R-40) — runs through the table format's reference implementation in
  the distributed engines. The systems-language library takes over each
  operation only once it ships and passes our conformance tests.
- **Practice and precedent:** the systems-language table library
  supports only fast appends; overwrite, partition replacement, file
  deletion, row-level changes and compaction are listed as missing,
  blocking even a foundation project's native write path. Another
  language's library fixed a bug in August 2026 where retrying a
  conflicting commit could resurrect deleted rows, because it replayed
  changes instead of re-applying them as the reference implementation
  does. [precedent]
- **Learned from Elysium:** wrote its mirror with a library in its own
  language, which lacked commit retries until recently. [docs]
- **Done when:** a check fails if any systems-language component
  issues a non-append table operation; and a forced conflict during a
  compaction never resurrects a deleted row.

### R-111 Tables use the format's third version

- **Class:** foundation
- **Status:** planned
- **Outcome:** every table is created at the table format's third
  version: row lineage gives each row a permanent id and the sequence
  number of the change that last touched it, which becomes the index
  version (R-105) and enables change capture straight from tables;
  deletion vectors make updates and deletes compact; built-in table
  encryption keys support R-91.
- **Practice and precedent:** the third version was ratified in 2025
  and shipped across major engines through 2026; a proposed fourth
  version would consolidate each commit into a single file, easing
  commit contention. [precedent]
- **Learned from Elysium:** used the second version. [code]
- **Done when:** a created table reports version 3, and an index write
  carries the row's last-updated sequence number.

### R-107 Freshness is declared per object type

- **Class:** improvement
- **Status:** planned
- **Outcome:** every object type declares its freshness: real-time
  (changes streamed from the source and processed continuously, target
  in seconds), near-real-time (small frequent batches, target in
  minutes) or scheduled (periodic runs). Each has a stated target and an
  alert when it is missed. A source that cannot stream its changes is
  checked on an interval, and its freshness can never be declared
  tighter than that interval.
- **Practice and precedent:** the leading platform streams changes
  into its index in seconds to minutes; practitioners state freshness as
  a service-level objective per dataset. End-to-end freshness is set by
  the pipeline, not the index, which is only the last step. [precedent]
- **Learned from Elysium:** read sources live or from a mirror
  refreshed on demand, with no declared targets. [code]
- **Done when:** a real-time type shows a source change on screen
  within its target under load, and a missed target raises an alert
  naming the type.

### R-97 Transformations are versioned, tested code

- **Class:** parity
- **Status:** planned
- **Outcome:** every transformation from raw to curated is code in
  version control, with tests and data checks, reviewed like any change
  and previewable on a branch before it publishes.
- **Practice and precedent:** transformation as code with tests and
  incremental models is the prevailing practice; the cleaned layer
  applies schema enforcement, deduplication and standardisation, and
  lineage from raw to cleaned is kept for compliance. [precedent]
- **Learned from Elysium:** its transform stage was code with
  provenance recorded, but individual transformations carried no data
  tests. [code]
- **Done when:** a transformation change that breaks a declared check
  fails review before it can publish.

### R-89 Data is published only after it passes its checks

- **Class:** parity
- **Status:** planned
- **Outcome:** each sync writes to a branch, runs its declared checks
  there, and publishes only on success.
- **Practice and precedent:** the established platform versions data
  on branches, and its pipelines halt on failed health checks before
  downstream data is affected. The table format already in use supports
  branches natively. [precedent]
- **Learned from Elysium:** a sync writes straight to the tables
  readers see; a drift refusal keeps the previous snapshot, but there is
  no staging. [code]
- **Done when:** a sync that fails a check leaves readers on the
  previous snapshot, and the failure is reported by name.

### R-54 A pointer is committed only after what it names is durable

- **Class:** parity
- **Status:** planned
- **Outcome:** every pointer swap follows a durability point.
- **Learned from Elysium:** a full disk left the catalog pointing at
  metadata that was never written; re-syncing could not repair it.
  [docs, fixed late]
- **Done when:** a fault injected between write and swap leaves the
  previous snapshot readable.

### R-56 History is custodial, so it is durable

- **Class:** parity
- **Status:** planned
- **Outcome:** the changelog ships only with durable object storage
  and verified backups.
- **Learned from Elysium:** its roadmap recognised that a changelog
  holds history no source can return, making the mirror a system of
  record. [docs]
- **Done when:** deleting a cell and restoring it recovers the
  changelog.

### R-90 Identity resolution in the curated layer

- **Class:** parity
- **Status:** planned
- **Outcome:** a curation stage, fed by the change log, resolves
  records from different sources that describe the same subject.
  Declared joins are primary; inference is off by default and every
  proposed merge is approved through the approval workflow; unmerging
  restores the originals; each field keeps its provenance and
  classification.
- **Learned from Elysium:** its fusion design places resolution
  upstream of the ontology, makes declared joins primary, keeps
  inference off by default with every merge approved, requires
  unmerging, preserves field provenance, and leaves the resolved link's
  own classification open. None of it is built. [docs]
- **Decision:** settled by precedent (2026-09-21): the high-water
  mark. A derived asset is as sensitive as its most sensitive input,
  formally the least upper bound of their labels. Each merged field
  keeps its own label; only the link itself takes the combined label of
  the fields that established it, because the link reveals that those
  records correspond. This avoids the rule's known failure, every record
  drifting to the top level.
- **Done when:** a declared join produces one subject whose fields
  keep their source classification, and an approved unmerge restores
  both.

### R-76 Classification follows lineage

- **Class:** parity
- **Status:** planned
- **Outcome:** classifications follow the data from each source
  through the raw, cleaned and curated layers into the ontology, and
  every stop in propagation requires review.
- **Practice and precedent:** markings propagate automatically to
  every derived dataset; stopping propagation is explicit in code and
  visible in the lineage graph; a change can be simulated before it
  applies. [precedent]
- **Learned from Elysium:** provenance is recorded per table, and the
  published manifest carries the configuration's security declarations,
  but nothing propagates to derived data. [measured]
- **Done when:** a derived dataset inherits its inputs' compartments,
  and an unreviewed stop fails validation.

### R-75 Distributed compute, batch and streaming

- **Class:** parity
- **Status:** planned
- **Outcome:** single-node by default. Beyond one node, the
  foundation-governed distributed batch engine runs with a native
  accelerator built on the same query engine as the read plane, and jobs
  are written against that engine's standard remote-client protocol, so
  a faster compatible engine can replace it without rewriting jobs.
  Streaming uses the foundation-governed stream processor at version 2,
  whose state lives on object storage and recovers in seconds. Streaming
  indexing merges edits from the start.
- **Practice and precedent:** autoscaling distributed batch and
  streaming engines run beside single-node engines, which handle
  terabyte-scale inputs of the right shape. Streaming object types there
  give up user edits and multi-source objects. [precedent]
- **Learned from Elysium:** single-process sync of whole tables.
  [docs]
- **Done when:** a transform over data larger than one node's memory
  completes, and a streaming type accepts an edit.

### R-58 Failures are pushed, not discovered

- **Class:** parity
- **Status:** planned
- **Outcome:** a refused sync, a stale table or a failed integrity
  check notifies its administrators when it happens.
- **Learned from Elysium:** a sync refused overnight was discovered
  only when someone opened a panel; notifications came later. [docs,
  fixed late]
- **Done when:** a refused sync produces one notification per eligible
  administrator, with repeats suppressed.

### R-40 Retention runs on a clock

- **Class:** parity
- **Status:** planned
- **Outcome:** expiry and retention are enforced on schedule, because
  a commercial product has deletion obligations.
- **Learned from Elysium:** expired artifacts are deleted only when
  something reads the store — deliberately, since nothing was harmed.
  [code]
- **Done when:** an expired record is gone within its window with no
  read in between.

---

# Phase 2 — Ontology and serving

The ontology is set on curated data and served to people, with the
security model enforced from the first request. Order of work:
configuration and service identity, then login and sessions, then
error and denial shapes, then serving and its security properties,
then audit and operations.

### R-51 Configuration is authored as code

- **Class:** parity
- **Status:** in progress — the format is specified; the shared
  library loads it and enforces its rules, reporting every problem at
  once; and a check fails if any other component parses ontology
  definitions (`docs/ontology-format.md`,
  `libs/ontology/src/validate.rs`, `script/check-ontology-parsers`);
  remaining: checking declared classifications against lineage and
  sources against the table catalogue, a check refusing breaking changes
  between versions, property visibility levels, bundle signing, and the
  allow-listed manifest.
- **Outcome:** ontology and policy authored in version control,
  compiled and signed by the control plane, never edited around it, and
  published beside the data as an allow-listed manifest. The meaning of
  the ontology is implemented exactly once, in one shared library in the
  systems language, used by the compiler, the indexer, the query service
  and the write service; the control plane calls the compiler rather
  than re-implementing it, so no two components can disagree about what
  the ontology means.
- **Decision:** the ontology is authored in TOML, one small file per
  type, link, action and metric, with one translation catalogue per
  language (owner, 2026-09-21). Only the shared library parses it, and
  in the systems language TOML has exactly one canonical, maintained
  library, while YAML's libraries are fragmented across individual
  maintainers after the long-standing one was archived in 2024
  (`RULES.md` E1). The format is specified in `docs/ontology-format.md`.
- **Done when:** the control plane refuses an unsigned bundle; the
  published manifest contains nothing outside its allow-list; and a
  check fails if any component other than the shared library parses
  ontology definitions.
### R-130 The ontology speaks every supported language

- **Class:** improvement
- **Status:** in progress — a missing translation, or a missing
  catalogue for a listed language, fails validation
  (`libs/ontology/tests/rules.rs`); remaining: every screen and agent
  showing the reader's language.
- **Outcome:** every object type, property, link type and action
  carries display names and descriptions per language in the ontology
  itself, authored with it (R-51) and checked for completeness. Every
  screen and every agent inherits them, so an ontology translated once
  is translated everywhere.
- **Practice and precedent:** the leading platform translates text
  inside individual applications rather than in the ontology they share;
  no evidence was found that its ontology's own names are translatable.
  [precedent]
- **Learned from Elysium:** had no translation of any kind. [code]
- **Done when:** a type missing a supported language's display name
  fails ontology validation, and a Spanish user sees Spanish type names
  on every screen and in agent answers.

### R-81 Services authenticate each other

- **Class:** foundation
- **Status:** planned
- **Outcome:** every internal call is mutually authenticated, with
  short-lived, automatically rotated workload identities.
- **Practice and precedent:** the federal guidance for microservices
  requires mutual authentication between services, a token service and
  key management. [precedent]
- **Learned from Elysium:** one process, so there were no internal
  calls.
- **Done when:** a call without a valid workload identity is refused,
  and an expired identity is refused after rotation.

### R-82 Each service names who may call it

- **Class:** foundation
- **Status:** planned
- **Outcome:** each service declares its permitted callers; all other
  callers are denied by default.
- **Practice and precedent:** the same guidance applies
  attribute-based access control between services, not only to users.
  [precedent]
- **Learned from Elysium:** one process. [code]
- **Done when:** a call from an undeclared service is refused, even
  with a valid identity and a valid user token.

### R-84 A revoked user stops everywhere, quickly

- **Class:** foundation
- **Status:** planned
- **Outcome:** internal user tokens live for a bounded, short time;
  disabling a user or revoking a grant publishes a revocation that every
  enforcing service honours within a stated bound; high-risk actions
  check revocation directly. A departing person's access everywhere ends
  within one hour, the strictest clock in the proposed revision of the
  US health-data law's Security Rule.
- **Practice and precedent:** a review of 62 studies lists
  insufficient token invalidation among recurring microservice
  vulnerabilities; the stale-permission exposure it causes is a named
  problem in large-scale authorization. [precedent]
- **Learned from Elysium:** disabling a user ended their sessions;
  there were no internal tokens to outlive them. [code]
- **Done when:** after a user is disabled, no service accepts their
  token beyond the stated bound, measured.

### R-44 Login input is bounded before anything persists

- **Class:** foundation
- **Status:** planned
- **Outcome:** both fields bounded at the edge.
- **Learned from Elysium:** usernames and passwords had no length
  limit. Ten unauthenticated requests grew the credentials database from
  100 KB to 16.3 MB, and an unbounded password makes every hash
  arbitrarily expensive. [measured]
- **Done when:** an oversized login is rejected and the test asserts
  no record was written.

### R-45 Expiring state expires by construction

- **Class:** foundation
- **Status:** planned
- **Outcome:** both live in a store whose keys expire.
- **Learned from Elysium:** failed-login records and expired sessions
  were never deleted. [measured]
- **Done when:** an expired entry is absent without any cleanup job
  having run.

### R-46 Session tokens hashed at rest

- **Class:** foundation
- **Status:** planned
- **Outcome:** only a hash of each token is stored.
- **Learned from Elysium:** stored in plain text, so reading the
  credentials store equals hijacking every live session. [measured]
- **Done when:** a test reading the store finds no value that works as
  a session cookie.

### R-47 Idle timeout as well as absolute expiry

- **Class:** foundation
- **Status:** planned
- **Outcome:** both.
- **Learned from Elysium:** a 24-hour absolute cap, no idle timeout;
  acknowledged. [code]
- **Done when:** a session idle past its limit is refused before its
  absolute expiry.

### R-48 Constant-time token comparison

- **Class:** foundation
- **Status:** planned
- **Outcome:** constant-time comparison for every secret.
- **Learned from Elysium:** the anti-forgery check compares with
  ordinary equality. Practically unexploitable, and free to fix. [code]
- **Done when:** a source-level check finds every comparison of a
  secret uses the constant-time function, with the test explaining why a
  timing test cannot be made reliable.

### R-67 Validation errors never reflect request bodies

- **Class:** foundation
- **Status:** planned
- **Outcome:** one error shape for every route; no submitted value is
  ever echoed.
- **Learned from Elysium:** an unauthenticated login missing its
  username returns the submitted password verbatim in the error body,
  because the default validation error includes the offending input.
  [measured]
- **Done when:** a request carrying a secret-shaped value in any field
  fails validation, and the value appears nowhere in the response or the
  logs.

### R-21 One uniform-denial convention

- **Class:** foundation
- **Status:** planned
- **Outcome:** one convention, decided once and documented in
  `BEHAVIOURS.md` before the first route exists.
- **Learned from Elysium:** denial is uniform within each endpoint but
  shaped differently across them — 200 with nulls, 200 with empty
  fields, 404, an empty list — and unauthenticated reads return 401
  while writes return 403. [measured]
- **Done when:** a test fetches a nonexistent and a forbidden resource
  from every endpoint and asserts byte-identical responses.

### R-22 Failure kinds in our own vocabulary

- **Class:** foundation
- **Status:** planned
- **Outcome:** failures map to a closed set of our own kinds:
  unreachable, refused, empty, misconfigured.
- **Learned from Elysium:** source status reports the runtime's
  exception class name, which reveals the implementation language and
  library. [code]
- **Done when:** no response body anywhere contains a runtime type
  name.

### R-70 The ontology is served from an index over curated data

- **Class:** improvement
- **Status:** planned
- **Outcome:** an indexer, fed incrementally from the curated layer's
  change log, maintains what the ontology serves; a separate query
  service answers from it, so each scales on its own. Early binding by
  default: the index holds security values, so filtering, counting and
  paging see visible rows only. A per-type security-freshness watermark
  advances with every security-relevant change; when the index is behind
  it, results are re-verified against the freshest curated data and
  counts are reported as unavailable, never overstated. Edits from
  approved actions are merged into what is served.
- **Practice and precedent:** an established platform separates
  indexing from querying so each scales horizontally. Enterprise search
  names the choice: early binding indexes permissions with content; late
  binding checks each result at query time; systems combine them,
  falling back to late binding when early cannot express a rule. An
  authorization system at very large scale answers stale-permission
  exposure with a freshness token: check against data at least as fresh
  as a given moment. Incrementally maintained views run in production,
  with modes where a read waits for changes to arrive. [precedent]
- **Learned from Elysium:** resolved meaning at read time and indexed
  nothing. Its change log and per-generation snapshot pinning are real
  inputs. Its mirror answered in 10.96 ms against 1.03 ms live because
  it served point lookups from a lake format, which the index must not
  do. [code, measured]
- **What it costs:** exposure bounded by the watermark rather than
  zero; custody of another copy of customer data; edits merged into the
  index; indexing to operate.
- **Decision:** a serving engine whose data structures are updated in
  place and searchable immediately, with structured filtering, text,
  vector retrieval and machine-learned ranking in one query (owner,
  2026-09-21). It has served billions of documents for over a decade,
  including about 800,000 queries a second across one large deployment.
  Chosen over a foundation-governed engine whose immutable segments
  rewrite every edit and publish changes about once a second. Governance
  risk accepted and contained: one company steers it, with its founding
  customer holding a stake and board seat; its licence is permissive, so
  released versions can always be used and forked; releases are pinned;
  any licence or governance change triggers a review; and because the
  index is disposable and sits behind our own interface (R-102), moving
  to another engine is a reindex, not a migration.
- **Done when:** after a reclassification, no member of the former
  audience sees the object once the watermark has advanced, measured; no
  count includes an invisible row; an edit is searchable within one
  second of acknowledgement; and a documented exit drill rebuilds the
  index on the alternative engine from curated data.
### ~~R-19 The lake format does not serve point lookups~~

- **Retired.** Merged into R-70, which now covers how the ontology is
  served.

### R-69 Execution tiers with published limits

- **Class:** parity
- **Status:** planned
- **Outcome:** the tier is chosen by estimated size, every limit is
  published, and work past the last limit is refused by name. The first
  release ships two tiers; the distributed tier depends on R-75.
- **Practice and precedent:** simple filters and aggregations push
  down to storage, sets up to 100,000 run in memory, and larger work
  moves to distributed compute, with published ceilings at each step.
  [precedent]
- **Learned from Elysium:** two tiers, observed working. Conditions a
  source declares it can evaluate are pushed to it; the rest run in
  process; a 10,000-row scan ceiling bounds both. No tier beyond one
  process. The per-source capability declaration is worth carrying over.
  [measured]
- **Done when:** the same query answered in each tier returns
  identical results, and a query past the last limit is refused with a
  named reason.

### R-68 No signal derived from hidden rows

- **Class:** new
- **Status:** planned
- **Outcome:** the security condition runs first, with a reviewed list
  of operators allowed ahead of it; counts, totals and truncation come
  only from rows that passed it.
- **Practice and precedent:** a published study found row-level
  security leaks the size of hidden row sets through timing when the
  caller's condition runs first. A mature database enforces the security
  condition first, except for operators certified unable to leak. A
  search product documents that post-filtered counts overstate what is
  visible. [precedent]
- **Learned from Elysium:** a user who can see no transactions is told
  the scan hit its ceiling, on both read paths, because the ceiling
  bounds a scan taken before the security check for link-secured types.
  [measured]
- **Done when:** DENY-09 passes with a fixture above the ceiling, and
  response time does not measurably separate a query over many hidden
  rows from one over none.

### R-71 Aggregate inference is controlled

- **Class:** new
- **Status:** planned
- **Outcome:** aggregates below a configurable group size are
  suppressed by default; differential privacy is opt-in for designated
  sensitive types, with a per-user budget.
- **Practice and precedent:** suppressing aggregates over very small
  groups is standard practice. A national statistics agency moved to
  differential privacy after reconstructing data for 17% of its
  population from tables protected by record swapping. A data platform
  offers differentially private aggregates, blocking row-level reads.
  The costs are documented: a spent privacy budget answers nothing more,
  and small groups lose accuracy. [precedent]
- **Learned from Elysium:** aggregates are exact at any group size.
  [code]
- **Done when:** every aggregate route suppresses a group below the
  threshold, and a designated type refuses row-level reads and stops
  answering when a budget is spent.

### R-20 One id representation everywhere

- **Class:** foundation
- **Status:** planned
- **Outcome:** every object id, including every link value, has one
  representation in every response.
- **Learned from Elysium:** on the live path, detail returned link
  values as integers while search returned the same ids as strings.
  [measured]
- **Done when:** QUERY-08 passes in full.

### R-105 The index never goes backwards

- **Class:** foundation
- **Status:** planned
- **Outcome:** every write to the search index carries the source
  change's sequence number as its version, and the index rejects any
  write older than what it holds. Rejected writes are counted, and a
  rising count raises an alarm.
- **Practice and precedent:** the search engine accepts externally
  supplied versions and rejects lower ones, for distributed processes
  that cannot guarantee order. One migration of 850 million documents
  found its default reindex would have silently overwritten 13.5 million
  newer documents with stale copies while reporting zero failures.
  [precedent]
- **Learned from Elysium:** nothing indexed, so the problem never
  arose. [code]
- **Done when:** delivering two updates to one object in reverse order
  leaves the newer one in the index, and the stale one is counted.


### R-106 Every state-changing request can be safely retried

- **Class:** foundation
- **Status:** planned
- **Outcome:** every state-changing endpoint accepts an idempotency
  key. The key, a fingerprint of the request and the response are stored
  durably, scoped to the caller, for a day. A retry returns the stored
  response; the same key with a different request, or while the first is
  still running, is refused.
- **Practice and precedent:** payment platforms established the
  pattern before a standards draft described it; keys typically expire
  after a day, and a request fingerprint catches a key reused with
  different content. The draft defines distinct errors for a missing
  required key, a key reused with a different request, and a retry that
  arrives while the first is in flight. [precedent]
- **Learned from Elysium:** no retry protection; a repeated request
  was a new request. [code]
- **Done when:** a request retried after a dropped connection is
  applied once, and reusing a key with a different body is refused.

### R-108 Staleness is always visible

- **Class:** improvement
- **Status:** planned
- **Outcome:** every record and result shows how current it is, and
  anything past its declared freshness target (R-107) is marked as stale
  rather than presented as current.
- **Practice and precedent:** stale data shown as current is a
  recognised failure of analytical systems; freshness indicators are
  standard in data products. [precedent]
- **Learned from Elysium:** reported whether reads were live or from
  its mirror, and when it last synced, on a separate screen. [code]
- **Done when:** a record whose type misses its target is visibly
  marked stale on every screen that shows it.

### R-109 Screens update live

- **Class:** improvement
- **Status:** planned
- **Outcome:** when a record changes, every open screen showing it
  updates without a reload; a person's own edits appear everywhere
  immediately (R-102), and a save confirms only once the change is
  visible.
- **Practice and precedent:** the leading platform's documentation
  says only one of its applications supports live data refresh, even for
  streaming object types, and others must be reloaded; serving engines
  that make writes visible immediately support live screens directly.
  [precedent]
- **Learned from Elysium:** screens refreshed only when reloaded.
  [code]
- **Done when:** an edit made in one session appears in another open
  session within the freshness target, without a reload.

### R-114 Large results travel as columnar streams

- **Class:** foundation
- **Status:** planned
- **Outcome:** bulk results move between services and to analytical
  clients as columnar streams over a standard columnar SQL wire
  protocol, including parallel streams for distributed reads;
  request-and-response contracts (R-06) stay for everything else.
- **Practice and precedent:** row-oriented database protocols
  transpose columnar results into rows and back; one practitioner
  describes a query computed in 300 milliseconds taking forty seconds to
  hand over. The columnar protocol also serves older row-based clients.
  [precedent]
- **Learned from Elysium:** returned every result as row-shaped JSON.
  [code]
- **Done when:** a ten-million-row result transfers at least an order
  of magnitude faster than the row-shaped equivalent, measured.

### R-116 Query plans use the standard cross-engine representation

- **Class:** foundation
- **Status:** planned
- **Outcome:** the query service can express its plans in the standard
  cross-engine query plan representation, so an accelerated engine —
  including GPU-native engines now emerging — can execute them without
  changing any caller.
- **Practice and precedent:** a 2026 research engine executing
  standard plans on GPUs reported about 7 times the speed of a CPU
  engine at the same hardware cost, and up to 12.5 times in a
  distributed setting; it plugs in through that standard representation.
  [precedent]
- **Learned from Elysium:** had no plan representation outside its own
  process. [code]
- **Done when:** a plan exported in the standard representation and
  re-imported produces identical results.

### R-131 Search works across languages

- **Class:** foundation
- **Status:** planned
- **Outcome:** every text field carries a language tag, set from
  source metadata or detected, and is processed with that language's
  rules for keyword search. Semantic search uses multilingual
  embeddings, so a query in one language finds relevant records written
  in another. Results say which language each record is in, and a
  person's own language is preferred without hiding others.
- **Practice and precedent:** the chosen serving engine supports
  several languages in one schema but not cross-lingual retrieval by
  itself; its documentation recommends multilingual embeddings for that.
  Research found models trained for cross-lingual retrieval consistently
  beat keyword matching and translating documents first, but quality
  varies widely between models and language pairs. [precedent]
- **Learned from Elysium:** indexed text without language awareness.
  [code]
- **Done when:** relevance is measured per supported language and
  across each pair, against stated targets, and a query in Portuguese
  finds a matching English record.

### R-112 Key exchange is post-quantum hybrid

- **Class:** foundation
- **Status:** planned
- **Outcome:** every TLS connection negotiates the hybrid post-quantum
  key exchange standardised as RFC 10024, internally and at the edge; a
  cryptographic inventory lists every key and algorithm in use;
  post-quantum signatures follow as their standards and libraries
  mature. A cell that must use validated cryptography runs in validated
  mode with classical key exchange, because the hybrid algorithm is not
  yet in validated modules; every other cell uses the hybrid exchange,
  and validated cells switch when validated modules include it.
- **Practice and precedent:** the hybrid exchange is the default in
  major browsers and in one services-language standard library since its
  2025 release; a national security algorithm timeline recommends hybrid
  deployment now and post-quantum-only operation by 2033; a 2026
  measurement found government and defence adoption effectively absent,
  while encrypted data can be harvested now and decrypted later.
  [precedent]
- **Learned from Elysium:** used classical TLS through a reverse
  proxy. [code]
- **Done when:** every internal and external connection negotiates the
  hybrid exchange, and a component that cannot fails the inventory
  check.

### R-119 No single machine can reach everyone's data

- **Class:** foundation
- **Status:** planned
- **Outcome:** no service instance holds standing credentials for all
  data. Credentials are short-lived and scoped to one tenant and one
  request, and requests are distributed so that nobody can steer a
  chosen user's work onto a machine they control.
- **Practice and precedent:** the strongest precedent requires that
  compromising one node must not let an attacker target a particular
  user, even with physical attacks in the supply chain, contrasting it
  with the common design in which every application server holds
  credentials for the whole database. [precedent]
- **Learned from Elysium:** one process held credentials for every
  source. [code]
- **Done when:** a compromised query-service instance, simulated, can
  read only the requests it happens to be serving, and cannot choose
  which user's requests it receives.

### R-120 Stateless work and no privileged runtime access

- **Class:** foundation
- **Status:** in progress — both production images contain their
  binary alone, with no shell or debugging tool, and the journey fails
  if one gains a file (`e2e/script/test-integration`); remaining: the
  memory scan proving no plaintext survives a request, and the reviewed
  channels operators observe through.
- **Outcome:** services keep no customer plaintext once a request
  completes, and production workloads have no shell, debugger or
  administrative path; operators observe them only through reviewed,
  structured channels.
- **Practice and precedent:** the same precedent requires stateless
  computation and that no staff member can extend their privilege, even
  during an outage. [precedent]
- **Learned from Elysium:** anyone with the host could inspect the
  process. [code]
- **Done when:** a production image contains no shell or debugging
  tool, and a memory scan after a request finds none of its plaintext.

### R-126 The mesh fits the cell

- **Class:** foundation
- **Status:** planned
- **Outcome:** standard cells use a sidecar-less service mesh for
  mutual TLS with the hybrid post-quantum exchange and identity-based
  policy. Confidential cells terminate encryption inside each protected
  workload, never in a shared per-node proxy, and issue mesh
  certificates only to workloads that pass attestation. Our own services
  can also run post-quantum TLS in-process.
- **Practice and precedent:** the zero-trust guidance for cloud-native
  applications places enforcement in a service mesh with workload
  identities; the sidecar-less mode runs one proxy per node and supports
  the hybrid exchange when configured, though not yet in
  validated-cryptography mode, and a failure of that shared proxy
  affects every workload on the node. A confidential mesh issues
  certificates only to attested workloads and terminates TLS in a proxy
  inside each protected pod; a per-node proxy outside the protected
  boundary would see plaintext. [precedent]
- **Learned from Elysium:** one process, no mesh. [code]
- **Done when:** in a confidential cell, no process outside a
  protected workload ever holds plaintext traffic, verified by test; in
  a standard cell, every connection negotiates the hybrid exchange.

### R-25 Say how authoritative a count is

- **Class:** improvement
- **Status:** planned
- **Outcome:** responses say whether a count is pinned to a snapshot.
- **Learned from Elysium:** the interface shows "N of M matches"
  unconditionally, which overstates a live deployment and understates a
  mirror. [docs]
- **Done when:** the interface's wording differs between the two, in a
  browser test.

### R-35 Cross-source questions are answered by the curated layer

- **Class:** improvement
- **Status:** planned
- **Outcome:** sources are joined in the pipelines, so every ontology
  type is served from curated data and no serving query spans stores;
  the limit disappears by design. Anything past a published serving
  limit is refused by name.
- **Learned from Elysium:** one search could draw on only one storage,
  because the ontology was resolved live across separate sources. [code]
- **Done when:** a type combining fields from two sources is served
  from one curated table, and no serving query touches more than one
  store.

### R-23 Audit that is durable and affordable

- **Class:** parity
- **Status:** planned
- **Outcome:** one record per bulk read, every denial its own record,
  durable before the response leaves, write-once retention.
- **Learned from Elysium:** one record per field access; at 200,000
  objects, 2.5 s of a 3.4 s aggregate was audit I/O. [docs, fixed late —
  bulk reads now write one record]
- **Done when:** a read of 10,000 objects writes one grant record, and
  a service killed after responding has written it.

### R-88 The audit trail is tamper-evident

- **Class:** new
- **Status:** planned
- **Outcome:** every audit record is added to such a tree; signed
  checkpoints are published to a witness the customer controls.
- **Practice and precedent:** transparency logs place each entry in an
  append-only hash tree; proving an entry is included costs a few dozen
  hashes at any size, signed checkpoints prove the tree only grew, and
  an external witness removes the need to trust whoever stores them. A
  commercial gateway sells this as a licensed feature. [precedent]
- **Learned from Elysium:** audit records are appended to files and
  are complete, but nothing proves after the fact that none was altered
  or removed. [code]
- **Done when:** altering or deleting any stored audit record makes
  verification against the last checkpoint fail.

### R-24 Administrators can trace any request

- **Class:** parity
- **Status:** planned
- **Outcome:** caller-scoped traces plus an audited, admin-only trace
  view.
- **Learned from Elysium:** only a caller's own trace is readable;
  debugging another user's request is not served. [docs]
- **Done when:** an administrator reads a user's trace, the read is
  itself audited, and an ordinary user still gets an empty list.

### R-36 Complete security headers, one owner each

- **Class:** foundation
- **Status:** planned
- **Outcome:** the full header set, each header set in exactly one
  documented place.
- **Learned from Elysium:** no permissions-policy header, and the
  transport-security header is left to an undocumented proxy. [measured]
- **Done when:** a test asserts every header on every route, and a
  duplicate setter fails it.

### R-83 Deadlines, circuit breakers and bulkheads

- **Class:** foundation
- **Status:** planned
- **Outcome:** a deadline set at the gateway travels with every call
  and shrinks at each hop; calls to a failing dependency trip a breaker;
  each dependency has its own bounded resources.
- **Practice and precedent:** synchronous call chains are a recognised
  cause of cascading failure; the established remedies are timeouts,
  circuit breaking, throttling and isolation of resources. [precedent]
- **Learned from Elysium:** no deadline on model calls and none on
  requests; one process meant one failure domain. [code]
- **Done when:** with one dependency stalled, requests that do not
  need it keep meeting their latency objective, and requests that do
  fail at their deadline with a clear error.

### R-74 The ontology is portable

- **Class:** new
- **Status:** in progress — the format is specified openly, and every
  example in the specification is tested against the library
  (`docs/ontology-format.md`, `libs/ontology/tests/specification.rs`);
  remaining: metrics as expressive as the open specification's — ratios,
  distinct counts, filters and several types — so the round trip
  preserves them, then the export and its round-trip test.
- **Outcome:** its own format is documented openly, and exports to the
  open specification.
- **Practice and precedent:** lock-in is the most consistent criticism
  of the established platform. The open, vendor-neutral semantic-model
  specification published in January 2026 under a permissive licence was
  accepted into a software foundation's incubator in July 2026 and
  renamed. Its core specification covers datasets, metrics, dimensions
  and relationships, with an extension mechanism for what core cannot
  express; its ontology specification, still a development version, is
  fact-oriented, with entity and value types related through verbalised
  relationships; and it carries no security model, warning that an
  exported model is not a security boundary. [precedent]
- **Learned from Elysium:** its ontology is a bespoke format. [code]
- **Decision:** Urshanabi keeps its own format, specified openly in
  `docs/ontology-format.md`, and exports to that specification: object
  types as entity types identified by their primary keys, properties as
  relationships to value types, links as relationships between entity
  types, and metrics as its metrics. Actions, edit rules, freshness and
  translations travel in its extension mechanism; security labels are
  never exported.
- **Done when:** an export-then-import round trip preserves types,
  links and metric definitions.

### R-101 Live access only to already-curated external data

- **Class:** parity
- **Status:** planned
- **Outcome:** the ontology is served from curated data. An external
  table may be served live, without copying, only if it is already
  curated in a governed store and meets the curated-layer contract
  (R-98). Raw silos always pass through the pipelines, because serving
  them live would skip the cleaning the product exists to do.
- **Practice and precedent:** the leading platform offers virtual
  tables beside its indexed objects. Access without copying is advised
  where an organisation already holds curated datasets in a governed
  warehouse; one virtualization vendor's own advice is to materialise
  any heavily queried dataset; and virtualization without a semantic
  layer is described as risky for an organisation. [precedent]
- **Learned from Elysium:** resolved everything live by design, then
  made its mirror the default read path. [code]
- **Decision:** settled by precedent (2026-09-21).
- **Done when:** binding a live external table without an accepted
  curated-layer contract fails validation, and a raw source cannot be
  bound live at all.
### ~~R-37 Multi-source reads run in parallel~~

- **Retired.** Applies only if R-101 chooses a live mode; otherwise
  sources are joined in the pipelines.

---

# Phase 3 — Agents

The same ontology, presented to agents: through the open agent-tool
protocol, and through Urshanabi's own agent.

### R-99 Agents reach the ontology through the open agent-tool protocol

- **Class:** parity
- **Status:** planned
- **Outcome:** the ontology is offered to any agent through the open,
  foundation-governed protocol agents use to reach tools and data, under
  the same policy as every other caller. Urshanabi's own agent is one
  client among many.
- **Practice and precedent:** that protocol moved to neutral
  foundation governance in December 2025 and is supported by the major
  agent clients. [precedent]
- **Learned from Elysium:** its own agent was the only way in. [code]
- **Done when:** an external agent client, authenticated as a user,
  sees exactly what that user sees: the uniform-denial properties pass
  through the protocol.

### R-132 Agents work in the person's language

- **Class:** foundation
- **Status:** planned
- **Outcome:** the agent reads questions and answers in the person's
  language, using the ontology's translated names (R-130); every agent
  evaluation runs in each supported language, and a language whose rates
  fall below target is not offered.
- **Practice and precedent:** model quality varies by language, so
  evaluation per language, not an English result assumed to carry over,
  is the standard for multilingual agents. [precedent]
- **Learned from Elysium:** its agent answered in English only. [code]
- **Done when:** the evaluation suite reports rates per language, and
  a Spanish question receives a Spanish answer naming Spanish type
  names.

### R-115 Agents work with other organisations' agents

- **Class:** parity
- **Status:** planned
- **Outcome:** Urshanabi's agent can be reached by, and can delegate
  to, other agents through the foundation-governed agent-to-agent
  protocol, always as the user and within the agent envelope (the
  security model in `docs/architecture.md`); tools and data stay on the
  agent-tool protocol (R-99).
- **Practice and precedent:** the agent-to-agent protocol reached
  production use across the major clouds with over 150 supporting
  organisations in its first year, and is described as complementary to
  the agent-tool protocol: one connects agents to peers, the other to
  tools. [precedent]
- **Learned from Elysium:** its agent could talk to nothing outside
  itself. [code]
- **Done when:** an external agent delegating a task receives only
  what the delegating user may see, and a delegated action goes through
  approvals like any other.

### R-26 The agent's context is budgeted

- **Class:** improvement
- **Status:** planned
- **Outcome:** a step returns a count and a sample; context is
  budgeted per hop and never compacted.
- **Learned from Elysium:** the loop overflowed its own 4,096-token
  window at hop four of a default eight (4,664 tokens, 114%); a step
  could return every matching id. [docs]
- **Done when:** an eight-hop run over the volume fixture stays under
  the window, and a test fails if any step returns more than its cap.

### R-27 Model calls have deadlines and accounting

- **Class:** foundation
- **Status:** planned
- **Outcome:** every model call has a deadline and returns a result
  carrying token counts.
- **Learned from Elysium:** the model interface accepts no timeout and
  returns a bare string, discarding token counts, so questions had no
  deadline and could not be costed. [code]
- **Done when:** a stalled model call ends at its deadline with a
  clear error, and every answer records its token cost.

### R-28 An evaluation harness with a regression gate

- **Class:** foundation
- **Status:** planned
- **Outcome:** each case runs many times; rates are compared against a
  stored baseline with a stated tolerance, nightly.
- **Learned from Elysium:** model-dependent tests were pass or fail on
  one run, and deselected by default. [docs]
- **Done when:** a prompt change that doubles the invalid-step rate
  fails the nightly run.

### R-29 Retrieved data is treated as able to steer

- **Class:** new
- **Status:** planned
- **Outcome:** prompt content derived from data is marked, recorded
  and evaluated; a planner that never sees raw data is assessed.
- **Learned from Elysium:** ontology data enters the same prompt as
  instructions, so an injected value can steer which objects are read
  and what the answer claims — within the user's own permissions. [docs]
- **Done when:** an injection fixture is recorded as injected in the
  trace, and the harness reports its success rate.

### R-30 Prompt caches never cross users

- **Class:** new
- **Status:** planned
- **Outcome:** per-user content at the head of every prompt, and no
  prefix-cache reuse across users on shared inference.
- **Learned from Elysium:** a shared 103-character prompt prefix let
  one user probe another's by timing; fixed to 2 characters, with users
  in the same role still sharing a whole system prompt. [docs, fixed
  late]
- **Done when:** the shared-prefix test fails past its budget, and a
  second user's request never hits the first user's cache.

### R-31 Production inference, not a desktop runner

- **Class:** parity
- **Status:** planned
- **Outcome:** hosted models or a batching inference server, chosen
  per tenant. A distributed inference layer that routes requests by
  cached prompt prefixes may be added for throughput only with caches
  salted per user, so no two users ever share a cache entry (R-30).
- **Learned from Elysium:** a desktop model runner with 3–4B models;
  its own notes say the query screen was too slow to use. [docs]
- **Done when:** a load test at the target concurrency meets the
  answer-latency objective.

### R-32 Output modes are chosen by measurement

- **Class:** foundation
- **Status:** planned
- **Outcome:** any change of output mode is decided by the harness.
- **Learned from Elysium:** its research found hard schema-constrained
  output raised validity but lowered accuracy. [docs]
- **Done when:** the harness reports validity and accuracy together.

### R-34 Starter questions reveal nothing

- **Class:** foundation
- **Status:** planned
- **Outcome:** starter questions carry no ids.
- **Learned from Elysium:** example questions name specific objects,
  which discloses their existence to users who cannot see them;
  deferred. [docs]
- **Done when:** validation rejects a starter containing an id.

### R-72 Agent defenses are evaluated adaptively

- **Class:** new
- **Status:** planned
- **Outcome:** the out-of-band pattern, and adaptive, defense-aware
  attacks in the nightly evaluation.
- **Practice and precedent:** research has converged on enforcing
  agent security outside the model, with a planner that never sees
  untrusted data. A 2026 study warns these defenses are validated only
  against fixed attacks — the method under which adaptive attacks later
  broke twelve in-model defenses at over 90% success. [precedent]
- **Learned from Elysium:** model tests pass or fail on one run;
  injection is not measured. [docs]
- **Done when:** the nightly run includes adaptive attacks and fails
  on a rise in their success rate.

### R-87 Model-provider credentials are isolated

- **Class:** improvement
- **Status:** planned
- **Outcome:** the model gateway holds no long-lived keys; it receives
  short-lived, scoped credentials from the secret store, runs in its own
  process with an egress allowlist of provider endpoints only, and loads
  no third-party plugins.
- **Practice and precedent:** the model-gateway compromise in R-86 was
  valuable precisely because a gateway concentrates every provider's
  keys. [precedent]
- **Learned from Elysium:** its model adapter held its connection
  details in the application process. [code]
- **Done when:** reading the gateway's environment and memory at rest
  yields no reusable provider key, and a call to an unlisted endpoint is
  refused.

---

# Phase 4 — Actions and automation

Changes proposed, approved and applied; writeback to sources only when
enabled.

### R-38 Pending writes are durable by construction

- **Class:** parity
- **Status:** planned
- **Outcome:** a proposal is a durable workflow from its first
  version.
- **Learned from Elysium:** the approval queue lived in process memory
  and was lost on restart; persistence was added later. [docs, fixed
  late]
- **Done when:** killing every service mid-approval loses nothing, and
  a crash mid-apply resumes correctly.

### R-39 Writes serialize per object across replicas

- **Class:** parity
- **Status:** planned
- **Outcome:** compare-and-swap against expected values, plus one
  in-flight write per object across the fleet.
- **Learned from Elysium:** per-object locks were in-process. [code]
- **Done when:** two replicas racing on one object produce one success
  and one refusal, with the interleaving forced.

### R-57 Approvals and actions run on durable execution, once

- **Class:** parity
- **Status:** planned
- **Outcome:** a foundation-governed durable workflow engine runs
  approvals, expiries, action retries and automations with
  single-execution semantics. Its partition count cannot be changed
  after creation, so it is sized from the scale objectives (R-66) before
  first use; small cells keep its history in the relational database,
  and large ones move to a store that grows by adding machines.
  Pipelines are scheduled by the data orchestrator (R-96).
- **Practice and precedent:** durable execution suits long-running
  business processes such as approvals, with exactly-once runs and
  persistent state, and is a different category from data orchestration
  (R-96). The chosen engine's design is proven at over 12 billion
  workflow runs and 270 billion steps a month at one large deployment.
  Published tests of the same design on a relational database found the
  database, not the engine, to be the limit. [precedent]
- **Learned from Elysium:** no scheduler; its own roadmap found the
  scheduler and multi-worker questions were the same decision. [code]
- **Decision:** the foundation-governed engine with the largest proven
  scale (owner, 2026-09-21). Chosen over a lighter foundation-governed
  engine that reuses our database but is the least proven at scale, and
  over a single-company engine of the same design, whose scale it
  matches. The costs accepted: its own cluster in every cell, and a
  foundation membership still at the entry tier, so it sits behind our
  own workflow interface (RULES.md H4a). Every cell's cluster is
  provisioned with 8,192 history shards, derived from the scale
  objectives (R-66): the engine's maintainers advised at least 4,000 for
  about 500 workflows a second, and the objective is 1,000. The count is
  confirmed by a load test before the first cluster is created, because
  it can never be changed.
- **Done when:** three replicas produce exactly one run of a scheduled
  expiry; the partition count is recorded with the calculation that
  produced it; and a load run at the R-66 objective completes without
  the database saturating.
### R-100 Writeback to sources is off by default

- **Class:** parity
- **Status:** planned
- **Outcome:** an approved action changes Urshanabi's own state at
  once. Pushing it back to a source is separate: enabled per source by
  an administrator, delivered through an outbox, and re-checked against
  the source's current value at push time, refusing a stale change.
- **Practice and precedent:** the leading platform keeps edits in its
  own object layer and merges them when serving; a transactional outbox
  is the standard way to deliver changes reliably to another system.
  [precedent]
- **Learned from Elysium:** writeback off by default and never
  pre-configured, with an outbox and a push-time re-check planned.
  [docs]
- **Done when:** with writeback disabled nothing reaches the source,
  and with it enabled a changed source value refuses the push by name.

### R-102 Edits and source data meet under declared rules

- **Class:** parity
- **Status:** in progress — an editable property without a declared
  rule fails validation, and a property that exists only through edits
  must keep them (`libs/ontology/tests/rules.rs`); remaining: the edits
  store, reconciliation on refresh, and the winning source shown on each
  field.
- **Outcome:** every editable field declares how a user's edit and
  fresh source data are reconciled — the edit persists, the most recent
  value wins, or a named source has priority — with no silent default.
  Fields that exist only through edits always keep them. Edits live in
  their own durable store, the index remains disposable, and every
  served field records which source won. The merged result is published
  back to the curated layer so downstream pipelines see it.
- **Practice and precedent:** master data management sets such rules
  per attribute, not per system, because no source is best for every
  field, and its standard strategies include source priority, most
  recent, most complete and a persistent manual override; each field
  records which source won. The leading platform offers edits-win and
  most-recent strategies per object type, keeps edits durably outside
  its index, and publishes merged results for downstream pipelines.
  [precedent]
- **Learned from Elysium:** edits and source data never met, because
  it served sources live. [code]
- **Done when:** an editable field without a declared rule fails
  validation; a pipeline refresh after an edit produces the value the
  declared rule dictates; and the winning source is visible on the
  field.


### R-138 Actions reach the leading platform's breadth

- **Class:** parity
- **Status:** planned
- **Outcome:** an action can create, modify, create-or-modify and
  delete objects, create and delete many-to-many links, run a function,
  call an external system, notify people and start a pipeline, each as a
  rule in the ontology's format (`docs/ontology-format.md`). Submission
  criteria decide who may submit an action and when, and a parameter may
  take the submitting person or the submission time as its value. Today
  an action only modifies the object it names (R-51).
- **Practice and precedent:** the leading platform's action rules
  create, modify, create-or-modify and delete objects, create and delete
  many-to-many links, run functions, call external systems before or
  after the edits, notify users, trigger builds and apply scenario
  edits; one-to-many links change through their foreign-key property.
  Actions carry submission criteria, and string and timestamp properties
  can take the current user or submission time as fixed values.
  [precedent]
- **Done when:** each rule kind has a test that applies it and one
  that refuses it when its submission criteria fail.

### R-103 Every change states the version it was based on

- **Class:** foundation
- **Status:** planned
- **Outcome:** every object carries a version tag; every edit must
  send the tag it was based on. A change sent without one is refused as
  a missing precondition, and a stale one is refused with the current
  version, so the person can review and retry. The check and the write
  happen atomically, and an action reads all its objects at one
  consistent version.
- **Practice and precedent:** the web's standard for conditional
  requests refuses a change whose version no longer matches with
  "precondition failed", preventing lost updates, and can require a
  precondition on every change. The leading platform loads all objects
  at the same versions throughout an action. [precedent]
- **Learned from Elysium:** it compared expected values before
  applying a write, but within one process. [code]
- **Done when:** two people editing one object from the same version
  produce one success and one refusal carrying the current version; an
  edit without a version is refused.

### R-41 Queues cannot be flooded

- **Class:** parity
- **Status:** planned
- **Outcome:** fan-out caps, an execute-once-for-all option, and a
  stated answer when a condition matches more than the cap.
- **Learned from Elysium:** its design notes warned that an automation
  firing across many objects floods the queue with writes that expire —
  15 minutes by default — before anyone reads them. Bulk actions have
  since gained a 1,000-object cap; the flooding risk for automations
  remains. [docs; the cap and default confirmed in code]
- **Done when:** a condition matching past the cap produces one clear
  refusal, not a flood.

### R-42 Automated changes follow risk-based maker-checker

- **Class:** parity
- **Status:** planned
- **Outcome:** every action type declares its risk. An automation acts
  as its owner, so a consequential action it proposes needs approval
  from a human other than that owner; action types declared low-risk may
  apply automatically. No one ever approves their own proposal.
- **Learned from Elysium:** undecided. [docs]
- **Practice and precedent:** maker-checker is required for sensitive
  operations by banking regulators and by financial-reporting and
  payment-card standards, with no self-approval; routine low-risk
  actions approve quickly while high-risk ones route to the required
  checker. [precedent]
- **Decision:** settled by precedent (2026-09-21).
- **Done when:** an automation-proposed consequential write waits for
  a human other than the automation's owner; a declared low-risk type
  applies automatically; self-approval is refused.
### R-43 Nothing holds a stale slice of configuration

- **Class:** parity
- **Status:** planned
- **Outcome:** every component reads configuration through the
  generation pinned to the current request.
- **Learned from Elysium:** after a reload, a surviving component
  still held the old roles, so creating a user with a new role failed
  until restart; pending writes were audited against the startup
  configuration; the thread pool's size was fixed at construction.
  [docs; the stale roles copy fixed late, the pool size confirmed in
  code]
- **Done when:** a reload adding a role lets a user be created with it
  immediately, with no restart.

---

# Phase 5 — Identity, tenancy and control plane

Federated identity, tenants, policy analysis and metering: what turns
one deployment into a product.

### R-49 Single sign-on first; passwords for break-glass

- **Class:** parity
- **Status:** planned
- **Outcome:** federated sign-on, directory provisioning and
  multi-factor authentication through an identity provider. Multi-factor
  authentication is mandatory for every account that can reach regulated
  data, and people reached through federation are identity-proofed to
  identity assurance level 2 and authenticator assurance level 2 of SP
  800-63.
- **Learned from Elysium:** local passwords were the only login.
  [code]
- **Done when:** a provisioned user signs in, and deprovisioning ends
  their sessions.

### R-50 Policy that can be analysed

- **Class:** improvement
- **Status:** planned
- **Outcome:** an analyzable policy language, a reload gate that
  blocks unintended widening, and an admin-only explainer.
- **Learned from Elysium:** grants are strings, so "does this change
  widen any role?" had no answer, and "why can this user not see X?" had
  at least four possible causes nobody could tell apart. [docs]
- **Done when:** a reload granting a role a new field is blocked until
  confirmed, and the explainer names the missing grant.

### R-52 Tenancy is a product decision, made first

- **Class:** parity
- **Status:** done — the decision is recorded here and in the
  deployment shape of `docs/architecture.md`.
- **Outcome:** an owner decision between hosted, customer-cloud and
  disconnected cells, recorded before this phase starts.
- **Learned from Elysium:** single-tenant by construction — recorded
  as a business decision, and it fixes pricing at one deployment per
  customer. [docs]
- **Decision:** hosted cells come first (owner, 2026-09-21).
  Customer-cloud and disconnected cells are out of scope until hosted is
  established; the single cell artifact stays designed for all three, so
  neither is ruled out. Because R-113 requires confidential computing
  before a cell is shared between customers, hosted begins with
  dedicated cells, one per customer, unless that isolation is ready
  first. Work that requires in-country infrastructure run by an
  accredited operator, such as a government's classified information,
  waits for customer-cloud cells.
- **Done when:** the decision is recorded here, with the models that
  are out of scope named.

### R-137 All deployment material lives in this repository

- **Class:** parity
- **Status:** planned
- **Outcome:** each service's container build, the release's
  deployment package for any conformant cluster, and the live
  configuration of every hosted cell — which release it runs, its sizing
  and its secret references — live in this repository, under `deploy/`.
  The concerns that usually separate them are handled here instead:
  configuration-only changes do not trigger rebuilds; the
  configuration's history is read by path; changes to live configuration
  require a second person's approval; and automated updates to it are
  confined to its own path, so they cannot trigger build loops. Secrets
  are never stored, only referenced (R-92). While the repository is
  public (R-17), cells are named by opaque identifiers, and no
  customer's name, endpoint, sizing or other identifying detail appears
  in it.
- **Practice and precedent:** the leading continuous-deployment tool's
  documentation highly recommends a separate repository for deployment
  configuration: to change configuration without rebuilding, to keep a
  clean audit trail, to deploy services from several repositories as one
  unit, to separate production access from code access, and to avoid
  automated commits triggering build loops. [precedent]
- **Learned from Elysium:** had no deployment configuration. [code]
- **Owner:** [NEEDS OWNER] whether customer-identifying configuration
  will be kept by making the repository private before the first hosted
  customer, or by keeping the mapping from opaque cell identifiers to
  customers outside it.
- **Decision:** all deployment material stays in this repository
  (owner, 2026-09-21), rather than splitting live hosted configuration
  into a separate repository as the precedent recommends.
- **Done when:** a configuration-only change deploys without a
  rebuild, a change to live configuration without a second approval is
  refused, and a check fails if any file under `deploy/` holds a
  customer-identifying detail while the repository is public.

### R-73 Metering, quotas and visible cost

- **Class:** improvement
- **Status:** planned
- **Outcome:** per-tenant metering and quotas, indexing throughput
  limits, a cost estimate shown before a query runs, and usage visible
  to the customer.
- **Practice and precedent:** an established platform meters compute
  and indexing throughput per object type; its users report limited cost
  visibility. [precedent]
- **Learned from Elysium:** request rate, errors and duration per
  route, and a 20-question window per user; nothing per tenant. [code]
- **Done when:** in a load test a noisy tenant is throttled without
  affecting another, and estimates fall within a stated tolerance.

### R-113 Tenant data is protected while in use

- **Class:** foundation
- **Status:** planned
- **Outcome:** pooled hosted cells run tenant workloads inside
  hardware-isolated confidential environments, and every workload must
  prove by remote attestation that it is genuine, unmodified code on
  genuine hardware before the key service releases any key; other cells
  may enable the same. Confidential computing is never the only
  protection: physical interposer attacks in 2025 and 2026 extracted
  attestation keys and broke integrity on fully updated hardware, so
  R-119 to R-122 limit what any single compromised machine can reach.
- **Practice and precedent:** hardware-isolated confidential computing
  is a production default across mainstream servers and major clouds in
  2026, with overhead reported below 5 percent for one of the two main
  technologies; the container-orchestrator integration became a
  foundation incubating project in July 2026; attestation gates key
  release so only verified code receives secrets. [precedent]
- **Learned from Elysium:** had no notion of tenants or of protecting
  data in memory. [code]
- **Done when:** a workload whose attestation fails receives no keys,
  and an operator with host access cannot read tenant data from memory,
  verified by test.

### R-122 Physical trust caps what a cell may hold

- **Class:** foundation
- **Status:** planned
- **Outcome:** every cell declares the physical trust of its hardware
  and site, and that sets the highest classification it may hold: a cell
  whose machines could be physically accessed by an adversary never
  holds data above its ceiling, whatever its software protections.
- **Practice and precedent:** confidential-computing threat models
  have always excluded physical attacks; in 2025 an interposer built for
  under 1,000 dollars extracted attestation keys from fully updated
  servers and forged attestations, and in September 2026 an active
  interposer costing under 200 dollars broke integrity on up-to-date
  hardware. [precedent]
- **Learned from Elysium:** no classification ceilings of any kind.
  [code]
- **Done when:** labelling a cell's site as physically untrusted
  refuses any data above its ceiling, at ingestion and at peering.

---

# Phase 6 — Commercial hardening

What an independent reviewer and a first customer will check.

### R-59 Network posture decided in one place

- **Class:** foundation
- **Status:** planned
- **Outcome:** a documented posture: where TLS ends, which proxies may
  set client-address headers, and binding as a configuration value.
- **Learned from Elysium:** bound to loopback for convenience, with
  TLS termination, trusted proxy headers and binding left undecided.
  [measured]
- **Done when:** a client-address header from an untrusted source is
  ignored.

### R-60 Graceful shutdown drains our own work

- **Class:** foundation
- **Status:** planned
- **Outcome:** every service stops accepting work, finishes what it
  holds, then exits.
- **Learned from Elysium:** no shutdown handling; a restart mid-write
  could split an audit pair or a lake commit. [code]
- **Done when:** a shutdown during a write leaves no half-written
  record.

### R-61 Startup validates what configuration cannot

- **Class:** foundation
- **Status:** planned
- **Outcome:** each cell checks every source, credential and table
  shape at startup and reports by name.
- **Learned from Elysium:** configuration was validated at load;
  source reachability, credentials and mirror shape were not checked
  until a first query failed. [code]
- **Done when:** an unreachable source is reported at startup, not on
  first use.

### R-63 Upgrades, not only installs

- **Class:** foundation
- **Status:** planned
- **Outcome:** every release upgrades the previous one in place,
  migrations included. Cells pull releases the way they pull
  configuration, staged by constraint, with automatic recall.
- **Practice and precedent:** agents inside each environment pull
  declarative plans; a hub releases a plan only when its constraints
  hold, and prioritises recalling a bad release. [precedent]
- **Learned from Elysium:** its install script is a fresh install, not
  an upgrade path. [docs]
- **Done when:** CI upgrades the last release's cell to the new one
  with data intact.

### R-64 Restore is exercised, not assumed

- **Class:** foundation
- **Status:** planned
- **Outcome:** backups are restored and checked on a schedule. The
  stated restoration target is 72 hours at most.
- **Learned from Elysium:** backup existed before restore did. [docs,
  fixed late]
- **Done when:** a scheduled job restores the latest backup and runs
  the conformance suite against it.

### R-65 An independent security audit and penetration test

- **Class:** foundation
- **Status:** planned
- **Outcome:** before the first external customer, and after any
  change to authentication, authorization or tenancy. Vulnerability
  scans run at least every six months and penetration tests at least
  every twelve.
- **Learned from Elysium:** reviewed once, externally, which found an
  unauthenticated flaw its own suite had missed. [measured]
- **Done when:** every finding is fixed or accepted in writing, and
  each fix has a test that fails if it is reverted.

### R-77 Ephemeral infrastructure, per-workload egress

- **Class:** parity
- **Status:** planned
- **Outcome:** an enforced maximum lifetime, and egress denied by
  default with a declared allowlist per workload.
- **Practice and precedent:** nodes live at most 48 hours and
  containers at most 72, so every service is built for failover and a
  compromise cannot persist; network egress is allowed per workload by
  container-level firewall rules. [precedent]
- **Learned from Elysium:** neither; its only egress control is the
  engine refusing to attach other databases. [code]
- **Done when:** a container past its maximum age is replaced, and a
  workload's connection to an undeclared host is refused.
### R-123 Asset inventory and network map, yearly

- **Class:** foundation
- **Status:** planned
- **Outcome:** a machine-generated inventory of every component,
  dependency, key and data flow, with a network map, is reviewed at
  least every twelve months and on every significant change.
- **Practice and precedent:** the proposed revision of the health-data
  Security Rule requires both, reviewed at least every twelve months.
  [precedent]
- **Learned from Elysium:** none. [code]
- **Done when:** the inventory is regenerated in CI and a component
  missing from it fails the build.

### R-124 Ready to be a health-data business associate

- **Class:** foundation
- **Status:** planned
- **Outcome:** Urshanabi can sign business associate agreements: its
  controls, breach notification, subcontractor terms and data return or
  destruction are documented and meet the Security Rule, so health
  customers can put regulated data into it.
- **Practice and precedent:** a service that holds a covered entity's
  health data is its business associate and must meet the Security Rule;
  the proposed revision strengthens business-associate oversight.
  [precedent]
- **Learned from Elysium:** never handled regulated health data.
  [docs]
- **Done when:** a template agreement exists, and every obligation in
  it maps to a control in the matrix (R-125).

---

# Phase 7 — Federation

Reaching beyond one installation: sharing with other organisations
under policy, and searching systems where they live. Designed in from
the start so earlier work never blocks it; built after the product is
established.

### R-117 Installations peer under enforceable policy

- **Class:** improvement
- **Status:** planned
- **Outcome:** separate installations — ours or any other conforming
  participant's — share selected objects over the open dataspace
  protocol being standardised through ISO/IEC. Each connection declares
  which types may flow, in which direction, and a classification
  ceiling; usage conditions travel with the data and are enforced by the
  recipient; exchange continues over low-bandwidth or disconnected links
  by queueing; schemas stay in step through signed bundles; conflicting
  values reconcile under R-102; classifications follow the data (R-76).
  Designed now: identifiers are translatable between installations,
  labels are portable, and nothing in the data model assumes a single
  installation.
- **Practice and precedent:** the leading platform peers objects
  between its own installations over proprietary connections that set
  allowed types, direction and a classification ceiling, and keep
  working when disconnected. Data spaces, defined in ISO/IEC 20151,
  share data between organisations under agreed policies, protocols and
  semantic models; their open protocol and trust protocol have been
  submitted for international standardisation, and a foundation-governed
  reference connector negotiates contracts, enforces usage policies and
  audits exchanges in production ecosystems. [precedent]
- **Learned from Elysium:** a single installation with no sharing.
  [code]
- **Decision:** peering is in scope, designed now and built later
  (owner, 2026-09-21).
- **Done when:** two installations exchange an object type under a
  ceiling, an object above the ceiling never crosses, and a conforming
  third-party participant can exchange with Urshanabi.

### R-118 Federated search, promoted through the pipelines

- **Class:** parity
- **Status:** planned
- **Outcome:** opt-in per source, Urshanabi can search external
  systems in place. Results are marked external and unverified, cannot
  be acted on, and are filtered both by the source's own permissions —
  searching as the user where the source allows — and by a
  classification ceiling per source; no count or result reveals anything
  the user cannot see (DENY-09). Every federated request states its
  purpose of use, and federated users are identity-proofed (R-49).
  Bringing a record in runs it through the pipelines, so everything in
  the ontology stays curated. Each source has time limits, rate limits
  and full audit.
- **Practice and precedent:** the US national health-data exchange
  framework is a network of networks: a discovery request is broadcast,
  matching records are retrieved where they live, every request states
  its purpose, users are identity-proofed to SP 800-63 IAL2 and AAL2,
  and it reached about 500 million records by February 2026. The leading
  platform lets users search external systems and promote records into
  its ontology, each datum tethered to its source. [precedent]
- **Learned from Elysium:** read only its own sources and mirror.
  [code]
- **Decision:** federated search is in scope, opt-in per source,
  following the health-exchange model (owner, 2026-09-21).
- **Done when:** a federated search returns external results marked
  unverified, a request without a stated purpose is refused, an
  unpromoted result cannot be acted on, and a promoted record appears
  only after its pipeline checks pass.
