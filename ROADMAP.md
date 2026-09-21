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
- **Urshanabi:** every inter-service interface is a versioned
  contract, checked for breaking changes in CI.
- **Done when:** removing a field from a contract fails CI.

### R-07 The conformance suite runs against Elysium

- **Urshanabi:** `conformance/BEHAVIOURS.md` implemented as a
  black-box HTTP suite, passing against Elysium first.
- **Done when:** the suite passes against Elysium in CI, and flipping
  one uniform-denial response in a test build fails it.

### R-08 A request id on every route from the first route

- **Elysium:** request correlation existed and one route of 39
  created an id, so most reads wrote untracked audit lines.
  [docs, fixed late]
- **Urshanabi:** every request carries an id from the gateway through
  every service, into every audit record and log line.
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
- **Urshanabi:** every release upgrades the previous one in place,
  migrations included.
- **Done when:** CI upgrades the last release's cell to the new one
  with data intact.

### R-64 Restore is exercised, not assumed

- **Elysium:** backup existed before restore did. [docs, fixed late]
- **Urshanabi:** backups are restored and checked on a schedule.
- **Done when:** a scheduled job restores the latest backup and runs
  the conformance suite against it.

### R-65 An independent security audit and penetration test

- **Urshanabi:** before the first external customer.
