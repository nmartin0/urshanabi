# Roadmap

The one list of work for Urshanabi, per `RULES.md` §18. Design
reasoning lives in design documents; what gets built, and in what
order, lives here.

**Every item exists to prevent a flaw Elysium had.** Each names the
Elysium evidence, what Urshanabi builds instead, and the test that
proves the prevention — a test that must fail if the prevention is
removed, per `RULES.md` §3. Flaws Elysium later fixed are included:
the goal is never to have them at all.

Ordering is by dependency, per `RULES.md` §17. Phases follow the
architecture proposal. Nothing in a later phase starts until what it
depends on is done.

## Evidence tags

- **[measured]** — measured directly against Elysium at commit
  `f4ea94e`, during its external review or since.
- **[code]** — confirmed by reading Elysium's code at that commit.
- **[docs]** — recorded in Elysium's own documents: its history, or
  measurements it took. Not re-checked here, because history cannot
  be and those measurements were not repeated.
- **fixed late** — added to any tag when Elysium had the flaw and
  later fixed it.
- **[precedent]** — established practice elsewhere, described by what
  it does. Never named, per `RULES.md` H1.

**Excluded or narrowed after checking.** Three flaws in Elysium's own
documents were stale: log rotation is listed as missing but ships as
a configuration file; migrations are listed as absent but exist, ad
hoc and unversioned (R-05 covers the real gap); and its README says
links cannot cross data sources, which they now can.

---

# Phase 0 — Foundations

Before any service holds real data. Everything here is cheap now and
expensive once data or users exist.

### R-01 Continuous integration from the first commit

- **Elysium:** no CI of any kind. Its lint script called itself the
  CI equivalent, and nothing ran it on push. [measured]
- **Urshanabi:** every push runs every lint, type, dead-code,
  boundary, lockfile and naming check, and every test suite, on a
  clean checkout.
- **Done when:** a deliberately broken commit fails the pipeline, for
  each check in turn.

### R-02 Tests declare their prerequisites

- **Elysium:** on a fresh clone, 15 of 1,849 unit tests failed with
  `assert 0 > 0`. They depended on seeded and synced data nobody had
  declared. [measured]
- **Urshanabi:** a test needing data declares a fixture that builds
  it, or skips with the exact command that would. The fast tier needs
  nothing outside the repository.
- **Done when:** a CI job runs the fast tier on a checkout with no
  generated data and every test passes or skips with a message.

### R-03 Browser tests in CI, not beside it

- **Elysium:** the simulated document that unit tests run against
  computes no layout, so 20 passing tests coexisted with a visibly
  broken screen. Real-browser tests existed and were run by hand.
  [docs]
- **Urshanabi:** real-browser tests run in CI against the built
  interface for anything layout, hit-testing or cascade dependent.
- **Done when:** a change that covers one element with another fails
  CI.

### R-04 One end-to-end test per wire

- **Elysium:** three shipped commits crashed on a user's first real
  run, each behind a green suite; three more bugs hid behind mocked
  callbacks. [docs]
- **Urshanabi:** every handler a test mocks has one test driving the
  real one through its real entry point.
- **Done when:** replacing a real handler with a stub fails at least
  one test.

### R-05 Versioned migrations before any store persists

- **Elysium:** stores create their tables if absent. One store has an
  ad hoc, unversioned migration function; none records a schema
  version. [code]
- **Urshanabi:** every store carries a schema version and ordered
  migrations from its first table.
- **Done when:** a store created at version 1 opens at version 2 with
  its data intact, and a downgrade is refused.

### R-06 Service contracts with breaking-change detection

- **Elysium:** one process; contracts were function signatures.
- **Precedent:** one vendor's open-source contract toolchain was
  built to scale its own microservices, but its compatibility checker
  was never released. The canonical choices are a schema-first binary
  interface language for internal calls, the industry-standard HTTP
  API description (governed by a neutral foundation) for the public
  edge, and the standard event-API description for streams. That
  vendor's error model separates safe arguments from unsafe ones, so
  the schema itself declares what may reach a log. [precedent]
- **Urshanabi:** internal calls defined in a schema-first binary
  interface language, with breaking-change checks in CI; the public
  HTTP description generated from it; the changelog and audit streams
  described with the standard event-API description. Every error
  argument is declared safe or unsafe, and only safe arguments are
  ever logged.
- **Done when:** removing a field from a contract fails CI, and
  logging an unsafe error argument fails a test.

### R-07 The conformance suite runs against Elysium

- **Urshanabi:** `conformance/BEHAVIOURS.md` implemented as a
  black-box HTTP suite, passing against Elysium first.
- **Done when:** the suite passes against Elysium in CI, and flipping
  one uniform-denial response in a test build fails it.

### R-08 A request id on every route from the first route

- **Elysium:** request correlation existed and one route of 39
  created an id, so most reads wrote untracked audit lines.
  [docs, fixed late]
- **Precedent:** diagnosing incidents by logging into individual
  servers is the recognised sign that observability was left too
  late. [precedent]
- **Urshanabi:** every request carries an id from the gateway through
  every service, into every audit record and log line, with traces,
  metrics and logs emitted to a vendor-neutral telemetry standard.
- **Done when:** a test enumerates every route and fails if any
  response lacks the id.

### R-09 Integrity checks that cannot be stripped

- **Elysium:** data-integrity checks were language assertions, which
  an optimisation flag removes; it had to refuse to start with that
  flag set. [code]
- **Urshanabi:** integrity checks are ordinary code paths that return
  or raise explicitly, in every language role.
- **Done when:** a release build with every optimisation enabled
  still rejects a corrupt record.

### R-10 Reproducible installs, exactly

- **Elysium:** dependencies were bounded rather than locked until
  late, and a comment still says locking was deferred after lockfiles
  shipped. [measured]
- **Urshanabi:** every language role installs from a hashed lockfile;
  CI fails when a lockfile drifts from its manifest.
- **Done when:** editing a manifest without regenerating its lockfile
  fails CI.

### R-11 Generated state is never tracked

- **Elysium:** generated databases and lake files were committed
  twice, leaving the tree permanently dirty and blocking patches for
  several rounds. [docs, fixed late]
- **Urshanabi:** ignore rules cover every generated path before the
  first generator exists.
- **Done when:** running every generator leaves version control
  reporting a clean tree.

### R-12 A development environment that survives a reboot

- **Elysium:** development data lived in a directory cleared on
  reboot, and was lost twice. [docs]
- **Urshanabi:** one command builds the whole local stack, with data
  under the user's own directory.
- **Done when:** a reboot followed by one command restores a working
  stack.

### R-13 Fixtures with volume

- **Elysium:** fixtures held 4 customers and 7 transactions, so paging
  and cut-off behaviour could not be seen. [measured]
- **Urshanabi:** a small fixture for speed and a volume fixture for
  paging, truncation and memory ceilings, both in CI.
- **Done when:** the volume fixture exercises every paging and
  truncation path in `BEHAVIOURS.md`.

### R-14 Tests pin a non-UTC timezone

- **Elysium:** its method notes that a UTC container hides every
  timezone bug. [docs]
- **Urshanabi:** CI runs date-sensitive tests in one zone behind UTC
  and one observing daylight saving.
- **Done when:** a deliberate local-time conversion fails in both.

### R-15 Documentation that cannot silently drift

- **Elysium:** 648 KB of top-level markdown — more than twice its
  source — across six overlapping roadmaps. Six documented claims
  were found stale, including a "real hole" already closed and two
  limitations already lifted. 54% of the backend is prose. [measured]
- **Urshanabi:** one roadmap (this file). Design documents hold
  reasoning, not work. A known limitation names what would close it,
  or cites a test asserting it still holds.
- **Done when:** fixing a documented limitation without updating its
  entry fails the test that cites it.

### R-16 The naming rule is enforced

- **Elysium:** its documents and comments name a competitor and many
  products throughout, and instruct quoting them.
- **Urshanabi:** CI scans everything written for the repository for
  company, product and language names.
- **Done when:** adding a product name to a comment fails CI.

### R-17 Repository visibility and licence

- **Elysium:** its licence describes unpublished proprietary source;
  its repository is publicly readable. [measured]
- **Urshanabi:** private from the first commit, or licensed
  deliberately. An owner decision.

### R-66 Stated scale objectives, each with a load test

- **Elysium:** no scale target anywhere. Before this item, the
  largest number in this roadmap was a 10-million-row sync. [measured]
- **Precedent:** an established ontology platform publishes its
  limits per object type: tens of billions of objects per type, a
  10-million-object ceiling on a link traversal's result, 10,000
  objects per edit, and metered indexing throughput. [precedent]
- **Urshanabi:** written objectives for objects per type, query
  latency at the 95th percentile, indexing lag, rows synced per hour
  and cell throughput. Every later scale item is measured against
  them.
- **Done when:** a change that breaks an objective fails a scheduled
  load run.

### R-78 Consumer-driven contract tests

- **Elysium:** one process, so there were no consumers to verify
  against.
- **Precedent:** contract testing, where producers and consumers
  verify compatibility independently, is the established way to
  catch breaking changes without full integration suites.
  [precedent]
- **Urshanabi:** each consumer publishes what it relies on; each
  producer's CI verifies against every consumer's expectations.
- **Done when:** a producer change that breaks one consumer's
  recorded expectation fails the producer's build.

### R-79 One writer per store

- **Elysium:** one process owned every store, so ownership was never
  written down. [code]
- **Precedent:** shared databases are the most cited cause of
  distributed monoliths, though sharing inside one repository with
  changes shipped together has worked; the rule underneath both is
  clear ownership. [precedent]
- **Urshanabi:** a written map naming the single writer of every
  store, enforced by credentials: only the owner holds write access.
- **Done when:** a service attempting to write a store it does not own
  is refused by the store, not by convention.

### R-80 Decompose on evidence

- **Elysium:** a single process throughout. [code]
- **Precedent:** over-splitting before a domain justifies it adds
  coordination cost with no benefit; the standard advice is to start
  modular and split when pressure proves it necessary, and to split
  by capability rather than by technical layer. [precedent]
- **Urshanabi:** every service contract exists from Phase 0, but
  Phase 1 ships as few processes as the security boundaries allow.
  The write-credential boundary justifies its own process; other
  splits wait for a measurement.
- **Done when:** each separately deployed service cites the
  measurement or security boundary that justified it.

### R-85 Every dependency passes the selection record

- **Elysium:** dependencies were chosen well but case by case; its
  own plans still name an object store whose open-source edition has
  since been archived. [code]
- **Precedent:** see `RULES.md` H4a — an archived object store, a
  mesh whose stable releases went vendor-only, and a streaming log
  that proved source-available. [precedent]
- **Urshanabi:** a dependency register records licence, governance,
  exit interface, disconnection and supply chain for every dependency.
  CI checks licences against an allowlist.
- **Done when:** adding a dependency with a source-available licence,
  or with no register entry, fails CI.

### R-86 The build pipeline is hardened against supply-chain attack

- **Elysium:** no CI, so no pipeline to attack yet. [measured]
- **Precedent:** in 2026 a widely used model-gateway library shipped
  credential-stealing releases after its publishing credentials were
  taken through a compromised vulnerability scanner in its CI; the
  malicious code ran at interpreter start-up without being imported,
  and the source repository was never changed. [precedent]
- **Urshanabi:** every CI action and tool is pinned by digest;
  workflows triggered by untrusted input get no secrets; publishing
  uses short-lived federated credentials, never stored tokens;
  release artifacts are rebuilt and compared with source; installs
  are audited for start-up hooks.
- **Done when:** a pipeline change that uses an unpinned tool, or
  exposes a secret to an untrusted trigger, fails a policy check.

### R-94 Nothing leaves a cell unless enabled

- **Elysium:** no telemetry of its own. [code]
- **Precedent:** an incremental streaming database considered for
  Urshanabi sends anonymous usage statistics by default. [precedent]
- **Urshanabi:** no component sends anything outside a cell unless an
  administrator enables it; dependency telemetry is switched off in
  every shipped configuration.
- **Done when:** a cell running under an egress capture sends
  nothing outward over a full test cycle.

### R-95 The interface component library is our own decision

- **Elysium:** its interface is built on a component library
  published by the company whose product is the precedent. [code]
- **Urshanabi:** keep it for Phase 1 behind our own component
  wrappers, and decide deliberately whether to replace it — a
  competitor's library carries its visual identity and its roadmap.
- **Owner:** [NEEDS OWNER] Keep with an exit plan, or replace.
- **Done when:** no screen imports the library directly; every use
  passes through our own wrappers.

---

# Phase 1 — The read path

### R-18 No shared in-process state

- **Elysium:** per-object locks, the pending-write working copy, the
  loaded configuration, lockouts, rate limits, ten single-writer
  embedded databases, a local audit file and a catalog in an embedded
  database — all tied it to one process, so a deploy was downtime.
  [code]
- **Urshanabi:** every service is stateless; state lives in the
  stores built for it.
- **Done when:** the conformance suite passes against two replicas of
  every service with requests alternating between them.

### R-19 The lake format does not serve point lookups

- **Elysium:** a mirror search took 10.96 ms against 1.03 ms live;
  86% was the table-format library reloading metadata and decoding
  manifests, twice per search. [measured]
- **Urshanabi:** table metadata cached per snapshot, one scan per
  search, and lookups served by the query engine.
- **Done when:** a benchmark gate in CI fails if a mirror search
  regresses past a threshold set from the first measured baseline.

### R-20 Ids are strings on every read path

- **Elysium:** on the live path, detail returned link values as
  integers while search returned the same ids as strings. [measured]
- **Urshanabi:** one id representation everywhere, including links.
- **Done when:** QUERY-08 passes against both read paths.

### R-21 One uniform-denial convention

- **Elysium:** denial is uniform within each endpoint but shaped
  differently across them — 200 with nulls, 200 with empty fields,
  404, an empty list — and unauthenticated reads return 401 while
  writes return 403. [measured]
- **Urshanabi:** one convention, decided once and documented in
  `BEHAVIOURS.md` before the first route exists.
- **Done when:** a test fetches a nonexistent and a forbidden
  resource from every endpoint and asserts byte-identical responses.

### R-22 Failure kinds in our own vocabulary

- **Elysium:** source status reports the runtime's exception class
  name, which reveals the implementation language and library.
  [code]
- **Urshanabi:** failures map to a closed set of our own kinds:
  unreachable, refused, empty, misconfigured.
- **Done when:** no response body anywhere contains a runtime type
  name.

### R-23 Audit that is durable and affordable

- **Elysium:** one record per field access; at 200,000 objects, 2.5 s
  of a 3.4 s aggregate was audit I/O. [docs, fixed late — bulk reads
  now write one record]
- **Urshanabi:** one record per bulk read, every denial its own
  record, durable before the response leaves, write-once retention.
- **Done when:** a read of 10,000 objects writes one grant record, and
  a service killed after responding has written it.

### R-24 Administrators can trace any request

- **Elysium:** only a caller's own trace is readable; debugging
  another user's request is not served. [docs]
- **Urshanabi:** caller-scoped traces plus an audited, admin-only
  trace view.
- **Done when:** an administrator reads a user's trace, the read is
  itself audited, and an ordinary user still gets an empty list.

### R-25 Say how authoritative a count is

- **Elysium:** the interface shows "N of M matches" unconditionally,
  which overstates a live deployment and understates a mirror.
  [docs]
- **Urshanabi:** responses say whether a count is pinned to a
  snapshot.
- **Done when:** the interface's wording differs between the two, in
  a browser test.

### R-26 The agent's context is budgeted

- **Elysium:** the loop overflowed its own 4,096-token window at hop
  four of a default eight (4,664 tokens, 114%); a step could return
  every matching id. [docs]
- **Urshanabi:** a step returns a count and a sample; context is
  budgeted per hop and never compacted.
- **Done when:** an eight-hop run over the volume fixture stays under
  the window, and a test fails if any step returns more than its cap.

### R-27 Model calls have deadlines and accounting

- **Elysium:** the model interface accepts no timeout and returns a
  bare string, discarding token counts, so questions had no deadline
  and could not be costed. [code]
- **Urshanabi:** every model call has a deadline and returns a result
  carrying token counts.
- **Done when:** a stalled model call ends at its deadline with a
  clear error, and every answer records its token cost.

### R-28 An evaluation harness with a regression gate

- **Elysium:** model-dependent tests were pass or fail on one run,
  and deselected by default. [docs]
- **Urshanabi:** each case runs many times; rates are compared
  against a stored baseline with a stated tolerance, nightly.
- **Done when:** a prompt change that doubles the invalid-step rate
  fails the nightly run.

### R-29 Retrieved data is treated as able to steer

- **Elysium:** ontology data enters the same prompt as instructions,
  so an injected value can steer which objects are read and what the
  answer claims — within the user's own permissions. [docs]
- **Urshanabi:** prompt content derived from data is marked, recorded
  and evaluated; a planner that never sees raw data is assessed.
- **Done when:** an injection fixture is recorded as injected in the
  trace, and the harness reports its success rate.

### R-30 Prompt caches never cross users

- **Elysium:** a shared 103-character prompt prefix let one user
  probe another's by timing; fixed to 2 characters, with users in the
  same role still sharing a whole system prompt. [docs, fixed late]
- **Urshanabi:** per-user content at the head of every prompt, and no
  prefix-cache reuse across users on shared inference.
- **Done when:** the shared-prefix test fails past its budget, and a
  second user's request never hits the first user's cache.

### R-31 Production inference, not a desktop runner

- **Elysium:** a desktop model runner with 3–4B models; its own notes
  say the query screen was too slow to use. [docs]
- **Urshanabi:** hosted models or a batching inference server, chosen
  per tenant.
- **Done when:** a load test at the target concurrency meets the
  answer-latency objective.

### R-32 Output modes are chosen by measurement

- **Elysium:** its research found hard schema-constrained output
  raised validity but lowered accuracy. [docs]
- **Urshanabi:** any change of output mode is decided by the harness.
- **Done when:** the harness reports validity and accuracy together.

### R-33 Nothing built that nothing uses

- **Elysium:** a memory guard was built and tested but never wired
  into the agent. [code]
- **Urshanabi:** a component ships with its caller, or not at all.
- **Done when:** the dead-code check fails on an uncalled public
  component.

### R-34 Starter questions reveal nothing

- **Elysium:** example questions name specific objects, which
  discloses their existence to users who cannot see them; deferred.
  [docs]
- **Urshanabi:** starter questions carry no ids.
- **Done when:** validation rejects a starter containing an id.

### R-35 Known limits refuse clearly

- **Elysium:** one search or field lookup may touch only one
  storage; a filter spanning two databases is unsolved. [code]
- **Urshanabi:** until lifted, the limit refuses with a clear error
  rather than returning a partial answer.
- **Done when:** a two-storage filter fails validation with a named
  reason.

### R-36 Complete security headers, one owner each

- **Elysium:** no permissions-policy header, and the
  transport-security header is left to an undocumented proxy.
  [measured]
- **Urshanabi:** the full header set, each header set in exactly one
  documented place.
- **Done when:** a test asserts every header on every route, and a
  duplicate setter fails it.

### R-37 Multi-source reads run in parallel

- **Elysium:** a field backed by another source waits for each read in
  turn, so latency adds rather than overlaps. [code]
- **Urshanabi:** independent source reads within a request overlap.
- **Done when:** three sources with a 100 ms delay each answer in
  about 100 ms, not 300 ms.

### R-67 Validation errors never reflect request bodies

- **Elysium:** an unauthenticated login missing its username returns
  the submitted password verbatim in the error body, because the
  default validation error includes the offending input. [measured]
- **Urshanabi:** one error shape for every route; no submitted value
  is ever echoed.
- **Done when:** a request carrying a secret-shaped value in any
  field fails validation, and the value appears nowhere in the
  response or the logs.

### R-68 No signal derived from hidden rows

- **Elysium:** a user who can see no transactions is told the scan
  hit its ceiling, on both read paths, because the ceiling bounds a
  scan taken before the security check for link-secured types.
  [measured]
- **Precedent:** a published study found row-level security leaks the
  size of hidden row sets through timing when the caller's condition
  runs first. A mature database enforces the security condition
  first, except for operators certified unable to leak. A search
  product documents that post-filtered counts overstate what is
  visible. [precedent]
- **Urshanabi:** the security condition runs first, with a reviewed
  list of operators allowed ahead of it; counts, totals and
  truncation come only from rows that passed it.
- **Done when:** DENY-09 passes with a fixture above the ceiling, and
  response time does not measurably separate a query over many hidden
  rows from one over none.

### R-69 Execution tiers with published limits

- **Elysium:** two tiers, observed working. Conditions a source
  declares it can evaluate are pushed to it; the rest run in process;
  a 10,000-row scan ceiling bounds both. No tier beyond one process.
  The per-source capability declaration is worth carrying over.
  [measured]
- **Precedent:** simple filters and aggregations push down to
  storage, sets up to 100,000 run in memory, and larger work moves to
  distributed compute, with published ceilings at each step.
  [precedent]
- **Urshanabi:** the tier is chosen by estimated size, every limit is
  published, and work past the last limit is refused by name.
- **Done when:** the same query answered in each tier returns
  identical results, and a query past the last limit is refused with
  a named reason.

### R-70 An indexed object layer — PROPOSED, AWAITING THE OWNER

- **Elysium:** meaning is resolved at read time and nothing is
  indexed; its transform stage deliberately materialises no
  per-type tables. Its identifier-keyed changelog and per-generation
  snapshot pinning are real, working inputs an index would consume.
  [code]
- **Precedent:** an established platform separates indexing from
  querying so each scales horizontally. Enterprise search names the
  choice: early binding indexes permissions with content; late
  binding checks each result at query time; systems combine them,
  falling back to late binding when early cannot express a rule. An
  authorization system at very large scale answers stale-permission
  exposure with a freshness token: check against data at least as
  fresh as a given moment. Incrementally maintained views run in
  production, with modes where a read waits for changes to arrive.
  [precedent]
- **Urshanabi, proposed:** an indexer fed incrementally from the
  changelog, and a separate query service. Early binding by default:
  the index holds security values, so filtering, counting and paging
  see visible rows only. A security-freshness watermark per object
  type advances with every security-relevant change. When the index
  is behind it, results are re-verified against the freshest source
  and counts are reported as unavailable, never overstated. Live
  reads remain available per type.
- **What it costs:** exposure bounded by the watermark rather than
  zero; custody of another copy of customer data; edits merged into
  the index; indexing to operate.
- **Owner:** [NEEDS OWNER] It reverses Elysium's read-time
  resolution, so `RULES.md` H5 requires explicit authorization first.
- **Done when:** after a reclassification, no member of the former
  audience sees the object once the watermark has advanced, measured;
  no count includes an invisible row.

### R-71 Aggregate inference is controlled

- **Elysium:** aggregates are exact at any group size. [code]
- **Precedent:** suppressing aggregates over very small groups is
  standard practice. A national statistics agency moved to
  differential privacy after reconstructing data for 17% of its
  population from tables protected by record swapping. A data
  platform offers differentially private aggregates, blocking
  row-level reads. The costs are documented: a spent privacy budget
  answers nothing more, and small groups lose accuracy. [precedent]
- **Urshanabi:** aggregates below a configurable group size are
  suppressed by default; differential privacy is opt-in for
  designated sensitive types, with a per-user budget.
- **Done when:** every aggregate route suppresses a group below the
  threshold, and a designated type refuses row-level reads and stops
  answering when a budget is spent.

### R-72 Agent defenses are evaluated adaptively

- **Elysium:** model tests pass or fail on one run; injection is not
  measured. [docs]
- **Precedent:** research has converged on enforcing agent security
  outside the model, with a planner that never sees untrusted data.
  A 2026 study warns these defenses are validated only against fixed
  attacks — the method under which adaptive attacks later broke
  twelve in-model defenses at over 90% success. [precedent]
- **Urshanabi:** the out-of-band pattern, and adaptive,
  defense-aware attacks in the nightly evaluation.
- **Done when:** the nightly run includes adaptive attacks and fails
  on a rise in their success rate.

### R-81 Services authenticate each other

- **Elysium:** one process, so there were no internal calls.
- **Precedent:** the federal guidance for microservices requires
  mutual authentication between services, a token service and key
  management. [precedent]
- **Urshanabi:** every internal call is mutually authenticated, with
  short-lived, automatically rotated workload identities.
- **Done when:** a call without a valid workload identity is refused,
  and an expired identity is refused after rotation.

### R-82 Each service names who may call it

- **Elysium:** one process. [code]
- **Precedent:** the same guidance applies attribute-based access
  control between services, not only to users. [precedent]
- **Urshanabi:** each service declares its permitted callers; all
  other callers are denied by default.
- **Done when:** a call from an undeclared service is refused, even
  with a valid identity and a valid user token.

### R-83 Deadlines, circuit breakers and bulkheads

- **Elysium:** no deadline on model calls and none on requests; one
  process meant one failure domain. [code]
- **Precedent:** synchronous call chains are a recognised cause of
  cascading failure; the established remedies are timeouts, circuit
  breaking, throttling and isolation of resources. [precedent]
- **Urshanabi:** a deadline set at the gateway travels with every
  call and shrinks at each hop; calls to a failing dependency trip a
  breaker; each dependency has its own bounded resources.
- **Done when:** with one dependency stalled, requests that do not
  need it keep meeting their latency objective, and requests that do
  fail at their deadline with a clear error.

### R-87 Model-provider credentials are isolated

- **Elysium:** its model adapter held its connection details in the
  application process. [code]
- **Precedent:** the model-gateway compromise in R-86 was valuable
  precisely because a gateway concentrates every provider's keys.
  [precedent]
- **Urshanabi:** the model gateway holds no long-lived keys; it
  receives short-lived, scoped credentials from the secret store,
  runs in its own process with an egress allowlist of provider
  endpoints only, and loads no third-party plugins.
- **Done when:** reading the gateway's environment and memory at rest
  yields no reusable provider key, and a call to an unlisted endpoint
  is refused.

### R-88 The audit trail is tamper-evident

- **Elysium:** audit records are appended to files and are complete,
  but nothing proves after the fact that none was altered or removed.
  [code]
- **Precedent:** transparency logs place each entry in an append-only
  hash tree; proving an entry is included costs a few dozen hashes at
  any size, signed checkpoints prove the tree only grew, and an
  external witness removes the need to trust whoever stores them. A
  commercial gateway sells this as a licensed feature. [precedent]
- **Urshanabi:** every audit record is added to such a tree; signed
  checkpoints are published to a witness the customer controls.
- **Done when:** altering or deleting any stored audit record makes
  verification against the last checkpoint fail.

### R-91 Encryption in transit and at rest, keyed per tenant

- **Elysium:** loopback-only transport; files at rest unencrypted.
  [code]
- **Urshanabi:** every connection encrypted; every store encrypted at
  rest with keys held in a key service; per-tenant keys once tenancy
  is decided (R-52), so one tenant's data can be destroyed by
  destroying its key.
- **Done when:** a store's files read without the key service yield
  no plaintext, and rotating a key re-encrypts without downtime.

### R-92 Secrets are references, never values

- **Elysium:** configuration may carry an environment-variable
  reference instead of a password, and a missing variable refuses to
  load, naming the field. [code]
- **Urshanabi:** the same rule, resolved from the secret store, with
  secrets never written to logs, errors, bundles or the published
  manifest.
- **Done when:** a bundle or manifest containing a secret value fails
  validation, and a missing reference refuses to load by name.

---

# Phase 2 — Writes

### R-38 Pending writes are durable by construction

- **Elysium:** the approval queue lived in process memory and was lost
  on restart; persistence was added later. [docs, fixed late]
- **Urshanabi:** a proposal is a durable workflow from its first
  version.
- **Done when:** killing every service mid-approval loses nothing, and
  a crash mid-apply resumes correctly.

### R-39 Writes serialize per object across replicas

- **Elysium:** per-object locks were in-process. [code]
- **Urshanabi:** compare-and-swap against expected values, plus one
  in-flight write per object across the fleet.
- **Done when:** two replicas racing on one object produce one success
  and one refusal, with the interleaving forced.

### R-40 Retention runs on a clock

- **Elysium:** expired artifacts are deleted only when something
  reads the store — deliberately, since nothing was harmed. [code]
- **Urshanabi:** expiry and retention are enforced on schedule,
  because a commercial product has deletion obligations.
- **Done when:** an expired record is gone within its window with no
  read in between.

### R-41 Queues cannot be flooded

- **Elysium:** its design notes warned that an automation firing
  across many objects floods the queue with writes that expire —
  15 minutes by default — before anyone reads them. Bulk actions have
  since gained a 1,000-object cap; the flooding risk for automations
  remains. [docs; the cap and default confirmed in code]
- **Urshanabi:** fan-out caps, an execute-once-for-all option, and a
  stated answer when a condition matches more than the cap.
- **Done when:** a condition matching past the cap produces one clear
  refusal, not a flood.

### R-42 Four-eyes against automation is decided

- **Elysium:** undecided. [docs]
- **Urshanabi:** an owner decision, recorded before automations can
  propose writes.

### R-43 Nothing holds a stale slice of configuration

- **Elysium:** after a reload, a surviving component still held the
  old roles, so creating a user with a new role failed until restart;
  pending writes were audited against the startup configuration; the
  thread pool's size was fixed at construction. [docs; the stale
  roles copy fixed late, the pool size confirmed in code]
- **Urshanabi:** every component reads configuration through the
  generation pinned to the current request.
- **Done when:** a reload adding a role lets a user be created with it
  immediately, with no restart.

---

# Phase 3 — Identity and control plane

### R-44 Login input is bounded before anything persists

- **Elysium:** usernames and passwords had no length limit. Ten
  unauthenticated requests grew the credentials database from 100 KB
  to 16.3 MB, and an unbounded password makes every hash
  arbitrarily expensive. [measured]
- **Urshanabi:** both fields bounded at the edge.
- **Done when:** an oversized login is rejected and the test asserts
  no record was written.

### R-45 Expiring state expires by construction

- **Elysium:** failed-login records and expired sessions were never
  deleted. [measured]
- **Urshanabi:** both live in a store whose keys expire.
- **Done when:** an expired entry is absent without any cleanup job
  having run.

### R-46 Session tokens hashed at rest

- **Elysium:** stored in plain text, so reading the credentials store
  equals hijacking every live session. [measured]
- **Urshanabi:** only a hash of each token is stored.
- **Done when:** a test reading the store finds no value that works
  as a session cookie.

### R-47 Idle timeout as well as absolute expiry

- **Elysium:** a 24-hour absolute cap, no idle timeout; acknowledged.
  [code]
- **Urshanabi:** both.
- **Done when:** a session idle past its limit is refused before its
  absolute expiry.

### R-48 Constant-time token comparison

- **Elysium:** the anti-forgery check compares with ordinary
  equality. Practically unexploitable, and free to fix. [code]
- **Urshanabi:** constant-time comparison for every secret.

### R-49 Single sign-on first; passwords for break-glass

- **Elysium:** local passwords were the only login. [code]
- **Urshanabi:** federated sign-on, directory provisioning and
  multi-factor authentication through an identity provider.
- **Done when:** a provisioned user signs in, and deprovisioning ends
  their sessions.

### R-50 Policy that can be analysed

- **Elysium:** grants are strings, so "does this change widen any
  role?" had no answer, and "why can this user not see X?" had at
  least four possible causes nobody could tell apart.
  [docs]
- **Urshanabi:** an analyzable policy language, a reload gate that
  blocks unintended widening, and an admin-only explainer.
- **Done when:** a reload granting a role a new field is blocked until
  confirmed, and the explainer names the missing grant.

### R-51 Configuration is authored as code

- **Urshanabi:** ontology and policy authored in version control,
  compiled and signed by the control plane, never edited around it.
  Published beside the data as an allow-listed manifest.
- **Done when:** the control plane refuses an unsigned bundle, and
  the published manifest contains nothing outside its allow-list.

### R-52 Tenancy is a product decision, made first

- **Elysium:** single-tenant by construction — recorded as a business
  decision, and it fixes pricing at one deployment per customer.
  [docs]
- **Urshanabi:** an owner decision between hosted, customer-cloud and
  disconnected cells, recorded before this phase starts.

### R-73 Metering, quotas and visible cost

- **Elysium:** request rate, errors and duration per route, and a
  20-question window per user; nothing per tenant. [code]
- **Precedent:** an established platform meters compute and indexing
  throughput per object type; its users report limited cost
  visibility. [precedent]
- **Urshanabi:** per-tenant metering and quotas, indexing throughput
  limits, a cost estimate shown before a query runs, and usage
  visible to the customer.
- **Done when:** in a load test a noisy tenant is throttled without
  affecting another, and estimates fall within a stated tolerance.

### R-74 The ontology is portable

- **Elysium:** its ontology is a bespoke format. [code]
- **Precedent:** lock-in is the most consistent criticism of the
  established platform. An open, vendor-neutral semantic-model
  specification was published in 2026 under a permissive licence and
  has entered an open-source foundation's incubator. [precedent]
- **Urshanabi:** its own format is documented openly, and exports to
  the open specification.
- **Done when:** an export-then-import round trip preserves types,
  links and metric definitions.

### R-84 A revoked user stops everywhere, quickly

- **Elysium:** disabling a user ended their sessions; there were no
  internal tokens to outlive them. [code]
- **Precedent:** a review of 62 studies lists insufficient token
  invalidation among recurring microservice vulnerabilities; the
  stale-permission exposure it causes is a named problem in
  large-scale authorization. [precedent]
- **Urshanabi:** internal user tokens live for a bounded, short time;
  disabling a user or revoking a grant publishes a revocation that
  every enforcing service honours within a stated bound; high-risk
  actions check revocation directly.
- **Done when:** after a user is disabled, no service accepts their
  token beyond the stated bound, measured.

---

# Phase 4 — Sync and automation

### R-53 Sync within a fixed memory ceiling

- **Elysium:** a sync held the whole table four times over, peaking
  at about 1.2 KB per row; 10 million rows would fail. [docs]
- **Urshanabi:** sync streams in bounded batches, and refuses a table
  it cannot hold rather than being killed.
- **Done when:** a 10-million-row table syncs under a fixed memory
  limit in CI.

### R-54 A pointer is committed only after what it names is durable

- **Elysium:** a full disk left the catalog pointing at metadata that
  was never written; re-syncing could not repair it. [docs, fixed
  late]
- **Urshanabi:** every pointer swap follows a durability point.
- **Done when:** a fault injected between write and swap leaves the
  previous snapshot readable.

### R-55 Bounded raw storage from the first sync

- **Elysium:** copy-on-write grew a 20,000-row table with one change
  from 177 KB to 839 KB over five syncs. [docs, fixed late]
- **Urshanabi:** the raw layer keeps two snapshots from day one;
  history lives in a changelog that grows only with real changes.
- **Done when:** fifty unchanged syncs leave storage flat.

### R-56 History is custodial, so it is durable

- **Elysium:** its roadmap recognised that a changelog holds history
  no source can return, making the mirror a system of record.
  [docs]
- **Urshanabi:** the changelog ships only with durable object storage
  and verified backups.
- **Done when:** deleting a cell and restoring it recovers the
  changelog.

### R-57 One scheduler, running once

- **Elysium:** no scheduler; its own roadmap found the scheduler and
  multi-worker questions were the same decision. [code]
- **Urshanabi:** the durable workflow engine schedules, with
  single-execution semantics.
- **Done when:** three replicas produce exactly one run of a
  scheduled job.

### R-58 Failures are pushed, not discovered

- **Elysium:** a sync refused overnight was discovered only when
  someone opened a panel; notifications came later. [docs, fixed late]
- **Urshanabi:** a refused sync, a stale table or a failed integrity
  check notifies its administrators when it happens.
- **Done when:** a refused sync produces one notification per eligible
  administrator, with repeats suppressed.

### R-75 Distributed compute, batch and streaming

- **Elysium:** single-process sync of whole tables. [docs]
- **Precedent:** autoscaling distributed batch and streaming engines
  run beside single-node engines, which handle terabyte-scale inputs
  of the right shape. Streaming object types there give up user edits
  and multi-source objects. [precedent]
- **Urshanabi:** single-node by default, distributed beyond one
  node, and streaming indexing designed to merge edits from the
  start.
- **Done when:** a transform over data larger than one node's memory
  completes, and a streaming type accepts an edit.

### R-76 Classification follows lineage

- **Elysium:** provenance is recorded per table, and the published
  manifest carries the configuration's security declarations, but
  nothing propagates to derived data. [measured]
- **Precedent:** markings propagate automatically to every derived
  dataset; stopping propagation is explicit in code and visible in
  the lineage graph; a change can be simulated before it applies.
  [precedent]
- **Urshanabi:** the same, with every stop requiring review.
- **Done when:** a derived dataset inherits its inputs' compartments,
  and an unreviewed stop fails validation.

### R-89 Data is published only after it passes its checks

- **Elysium:** a sync writes straight to the tables readers see; a
  drift refusal keeps the previous snapshot, but there is no staging.
  [code]
- **Precedent:** the established platform versions data on branches,
  and its pipelines halt on failed health checks before downstream
  data is affected. The table format already in use supports branches
  natively. [precedent]
- **Urshanabi:** each sync writes to a branch, runs its declared
  checks there, and publishes only on success.
- **Done when:** a sync that fails a check leaves readers on the
  previous snapshot, and the failure is reported by name.

### R-90 Identity resolution in a gold layer

- **Elysium:** its fusion design places resolution upstream of the
  ontology, makes declared joins primary, keeps inference off by
  default with every merge approved, requires unmerging, preserves
  field provenance, and leaves the resolved link's own classification
  open. None of it is built. [docs]
- **Urshanabi:** that design, built in a gold layer fed by the
  changelog, with merges and unmerges passing through the approval
  workflow.
- **Done when:** a declared join produces one subject whose fields
  keep their source classification, and an approved unmerge restores
  both.

### R-93 A source that changes shape is refused, not absorbed

- **Elysium:** a column whose values no longer match the declared type
  produces a drift report naming the column, the type and an example,
  and the sync is refused rather than silently coerced. [code]
- **Urshanabi:** the same, per source, with severity tiers and a
  notification to administrators (R-58).
- **Done when:** a source column changing type refuses the sync with a
  report naming the column.

---

# Phase 5 — Commercial hardening

### R-59 Network posture decided in one place

- **Elysium:** bound to loopback for convenience, with TLS
  termination, trusted proxy headers and binding left undecided.
  [measured]
- **Urshanabi:** a documented posture: where TLS ends, which proxies
  may set client-address headers, and binding as a configuration
  value.
- **Done when:** a client-address header from an untrusted source is
  ignored.

### R-60 Graceful shutdown drains our own work

- **Elysium:** no shutdown handling; a restart mid-write could split
  an audit pair or a lake commit. [code]
- **Urshanabi:** every service stops accepting work, finishes what it
  holds, then exits.
- **Done when:** a shutdown during a write leaves no half-written
  record.

### R-61 Startup validates what configuration cannot

- **Elysium:** configuration was validated at load; source
  reachability, credentials and mirror shape were not checked until a
  first query failed. [code]
- **Urshanabi:** each cell checks every source, credential and table
  shape at startup and reports by name.
- **Done when:** an unreachable source is reported at startup, not on
  first use.

### R-62 Read credentials are checked, not assumed

- **Elysium:** documented that sources should use read-only
  credentials, and could not check it. [docs]
- **Urshanabi:** at startup, each read connection tests whether its
  credential can write, and refuses or warns by policy.
- **Done when:** a read connection holding write privilege is flagged.

### R-63 Upgrades, not only installs

- **Elysium:** its install script is a fresh install, not an upgrade
  path. [docs]
- **Precedent:** agents inside each environment pull declarative
  plans; a hub releases a plan only when its constraints hold, and
  prioritises recalling a bad release. [precedent]
- **Urshanabi:** every release upgrades the previous one in place,
  migrations included. Cells pull releases the way they pull
  configuration, staged by constraint, with automatic recall.
- **Done when:** CI upgrades the last release's cell to the new one
  with data intact.

### R-64 Restore is exercised, not assumed

- **Elysium:** backup existed before restore did. [docs, fixed late]
- **Urshanabi:** backups are restored and checked on a schedule.
- **Done when:** a scheduled job restores the latest backup and runs
  the conformance suite against it.

### R-65 An independent security audit and penetration test

- **Urshanabi:** before the first external customer.

### R-77 Ephemeral infrastructure, per-workload egress

- **Elysium:** neither; its only egress control is the engine
  refusing to attach other databases. [code]
- **Precedent:** nodes live at most 48 hours and containers at most
  72, so every service is built for failover and a compromise cannot
  persist; network egress is allowed per workload by container-level
  firewall rules. [precedent]
- **Urshanabi:** an enforced maximum lifetime, and egress denied by
  default with a declared allowlist per workload.
- **Done when:** a container past its maximum age is replaced, and a
  workload's connection to an undeclared host is refused.
