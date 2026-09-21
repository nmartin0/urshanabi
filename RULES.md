# RULES

How work is done on Urshanabi, by every person and every agent.

---

# Purpose: what Urshanabi is, and what it is not

**Urshanabi takes an organisation's isolated data silos, cleans the
data through pipelines, and sets an ontology on the final, distilled
data.** The ontology is presented to people, and to agents, alike.

**It belongs to an emerging category.** One established platform is
the first instance and the leader, and the category is about to fill
with systems of this kind. Urshanabi is one of them, built to its own
design.

**It is not a copy.** Not of the leading platform, and not of Elysium,
the owner's earlier prototype. Neither one's architecture or
limitations constrain Urshanabi. They are inspiration and a loose map:
evidence of what the category needs, and of what went wrong along the
way.

**Every stage follows industry-standard practice.** Data sourcing,
pipelines, storage, the ontology, serving, security, service
architecture, delivery: for each, the question is what current
practice and research establish as best, not what the leader or
Elysium happened to do.

Part one is hard rules. They have no exceptions unless the owner
authorizes one explicitly, and that authorization is recorded where
the exception lives.

Part two is the working method, learned while building Elysium,
where every rule was learned by getting something wrong. The method
carries over; Elysium's design does not. Where the mistake
is instructive it is named, because a rule without its reason gets
worked around the first time it is inconvenient.

Part three is the engineering standards: how code is chosen,
arranged, checked and tested. Like part one, they bind every change;
unlike part one, each is enforced by a gate in CI rather than by
review alone.

Keep this file precise rather than long. Instructions that restate
what a reader can learn from the code make work worse, not better.

---

# Part one: hard rules

## H1. Never name a company or a product.

**Urshanabi is our own product, and everything written for it stays
free of other organisations' names.** No company, vendor, competitor,
product, open-source project, programming language, research system
or published tool is named.

**Scope: everything written in or for this repository.** Prose,
documentation, code comments, identifiers we choose, test names,
filenames, configuration comments, commit messages, user-facing text
and error messages.

**Describe the thing by what it does.** Not the name of a durable
workflow engine, but "the durable workflow engine". Not a language's
name, but its role: the systems language, the services language.
The test: can a reader understand what the thing does without the
name? If not, the description is not finished.

**Where a name is unavoidable, it appears only there.** Dependency
manifests, lockfiles and import statements must name what they
depend on; that is the one place a name may appear, and it is not
repeated in any comment beside it. A development tool's own
configuration file counts as a manifest, since the tool fixes its
name and location; the file carries settings, never commentary about
the tool.

**The name check enforces this rule (confirmed by the owner).** Its
forbidden terms live in `script/names.list`, its configuration, as
committed banned-term lists do in prose linters. A code line that must
invoke a tool by name — the shell's equivalent of an import — may end
with the marker `# name-ok`; the check refuses the marker in Markdown
and on comment lines, so it can never excuse prose. Reviewed English
phrases that happen to contain a listed word, such as the title of H8,
are listed in `script/names.allow`.

**A banned-term list catches only the names it knows.** So whenever a
name is recorded in the named register kept outside the repository, it
is added to `script/names.list` in the same change. When the list was
first automated it lacked about 150 names researched after it was
written; none had reached the repository, but only because the
documents had been written carefully.

**Provisional, pending the owner's confirmation:** generic technical
vocabulary and published protocol designations — HTTP, TLS, SQL,
JSON, YAML, OIDC, SAML, SCIM, RFC numbers — are treated as vocabulary
rather than product names. So are the designations of standards and
regulations that buyers require us to meet — control-catalogue and
publication numbers such as SP 800-53, ISO/IEC numbers, and legal
citations such as 45 CFR Part 164 — because compliance cannot be
shown without citing them. The bodies that issue them are described,
not named.

**Carrying anything over from Elysium means scrubbing it on the way.**
Elysium's documents and comments name companies and products freely.
Nothing crosses into this repository until those names are removed
and the sentence still reads correctly without them.

## H2. Never copy code from outside our own repositories.

**Not verbatim, not lightly edited, not "adapted".** This covers
documentation examples, answers on forums, tutorials, other projects'
source, specification appendices, and anything recalled from memory
that reproduces someone else's code.

**The method is: read to understand, close it, write our own.** The
implementation comes from our understanding of the problem, not from
the text of someone else's solution.

**Our own repositories are the exception.** Elysium and Urshanabi
belong to the same owner, and code moves between them freely —
subject to H1's scrubbing.

**If our implementation would be near-identical to a specific
outside source** because the problem allows only one shape, stop and
ask the owner before writing it.

**The same applies to prose.** Outside text is paraphrased, never
quoted into this repository.

## H3. Research precedent; record it without names or quotations.

Checking what an established system does before inventing a pattern
remains the rule — a problem with a name has usually been solved
badly several times, and the failure modes are written down.

**What changes is how it is recorded.** A commit message or design
note describes the precedent in our own words: "an established
ontology platform never touches the source row on delete; it records
an edit". No name, no quotation, no link to the source in the
repository.

**Reverse a committed decision when research contradicts it**, and
say so in the record. A design document that hides its wrong turns
teaches nobody.

**And check the precedent applies.** A standard way to load remote
interface code into a host page provides no isolation at all, which
disqualifies it for untrusted code. The precedent was right and the
context was different.

## H4. Dependencies are installed and used, never modified.

**Every dependency is installed unmodified, from its package index,
and used through its public, documented interface.** Extending one
through its own documented extension point is normal use.

**Using a dependency is fine regardless of its licence.** Importing,
linking and calling create no modification obligations. **Anything we
actually modify must be permissively licensed** — and under this
rule we modify nothing, which keeps the question moot.

A licence check that fails the build on any copyleft dependency
enforces the wrong rule. Check what we modify, which is nothing.

## H4a. Dependencies are chosen for licence, governance and exit.

**"Open source" is not a sufficient test.** Three precedents, each
recent: a widely used object store stopped publishing binaries,
entered maintenance mode, and had its repository archived, leaving a
high-severity fix unshipped in any official image. A service mesh
stopped producing open-source stable releases; stable builds now come
only from its commercial steward. A streaming log adopted for being a
single binary turned out to be source-available, barring use as a
commercial service.

Before a dependency is adopted, record:

- **Licence.** Anything shipped inside a cell must carry a licence
  that meets the recognised open-source definition, because
  disconnected installs distribute it.
  Source-available licences are refused. Strong copyleft is accepted
  only unmodified, in a separate process, and flagged for buyers.
- **Governance.** Who controls the roadmap, how many organisations
  maintain it, and whether releases are published openly. A single
  vendor needs a written exit plan.
- **Exit.** The dependency sits behind an interface Urshanabi owns,
  so replacing it is a migration, not a redesign.
- **Disconnection.** It runs with no network access, and any
  telemetry it sends by default is switched off and verified off.
- **Supply chain.** It is installed pinned by hash, with published
  provenance verified where it exists.

**Tools count as dependencies.** A vulnerability scanner running in
CI was compromised and used to steal the publishing credentials of a
widely used library, which then shipped credential-stealing releases.
Every tool in CI is pinned by digest and kept away from publishing
credentials.

## H5. Breaking changes need explicit authorization.

A change to a persisted format, a configuration schema, a security
invariant, or anything a deployment depends on is made only after
the owner authorizes it in conversation. "It seemed better" is not
authorization.

## H6. Security invariants.

These hold everywhere. Each has a test that fails if it stops being
true.

- **Authorization is decided in one place.** Two places deciding one
  question is how they drift.
- **Attenuation only.** No mechanism may let a request do something
  the requesting user could not already do.
- **Authority is never stored.** Store what was asked for; decide
  again at the point of use, against the current configuration.
- **Data does not cross compartments.** What an action reads must not
  be written somewhere readable by people who could not read it.
- **"Does not exist" and "exists but denied" are indistinguishable**
  to the caller — same shape, same status, same message. The real
  reason goes to the audit trail only.
- **Fail safe.** A missing grant, role or setting denies. A missing
  environment variable stops the load; it never becomes an empty
  string.
- **Names, never contents.** Status and error output may name a
  source or a role; never its path, connection details or grants.
- **No credential in any error message.** An error is the likeliest
  place for a secret to reach a log someone pastes.
- **Compute per recipient; never filter after assembly.** Filtering
  after assembly is where systems leak, because the unfiltered thing
  existed.
- **A display filter is not a security decision.** The interface may
  hide what the server already refused; it is never the only check.

---

## H7. Decide from industry-standard practice; never improvise.

**When a choice is uncertain, research it before making it.** Look up
what established practice and current research say for that stage —
data sourcing, pipelines, storage, service architecture, security,
delivery — and record the evidence with the decision, described
without names (H1, H3).

A plausible design made up on the spot is not a decision. This
project has already corrected one: pipeline scheduling had been given
to the durable engine that runs approvals, when industry practice
treats data orchestration and durable execution as different
categories that are routinely confused.

## H8. Close calls go to the owner.

**When the research leaves the options roughly balanced, stop and ask
the owner.** Present the options, the evidence for each, the costs,
and a recommendation if there is one. Do not choose, and do not build
on a choice that has not been made.

A close call recorded in the roadmap is marked [NEEDS OWNER] and
stays open until the owner decides.

---

# Part two: the working method

## 1. Verify directly. Never assume.

**Read the code before describing it. Run the command before
recommending it. Measure the cost before claiming it.**

The recurring failure is asserting a general case from a check of a
specific one. **A diagnostic that has never been executed is a guess
with a prompt in front of it.**

**A search is not proof.** "Nothing runs in parallel" was wrong
because the search did not cover the directory where it did. Check a
claim where it would actually live, and check the negative case.

**Read before assuming something is missing.** Roadmap entries have
repeatedly described defects the code had already ruled out.

**Your own notes decay.** Re-verify before building from a note
written a month ago.

## 2. Audit your own work before delivering it.

**After the work is done and before it is handed over, read what you
changed as though somebody else wrote it.**

- Review the whole diff, not the parts you remember editing.
- List untracked files as well — a diff does not show new files.
- Re-read every comment against the code beside it. A comment
  describing an earlier draft is worse than none.
- Every claim in the commit message is measured, or removed.
- Every new test fails when the change is reverted. If you have not
  run that control, you do not know.
- Look for anything nothing calls.

**This is a separate step, not a feeling.**

## 3. Real tests, with real negative controls.

**A test that passes when the feature is removed tests nothing.**

**The control is how you find out whether you wrote a test.** State
the property in one sentence, then name the code change that would
make it false. Break the code, watch the right test fail for the
right reason, restore from a backup copy, watch it pass.

- **A test written beside the code tends to assert its shape rather
  than its property.** Only breaking it separates the two.
- **Check the opposite direction.** A guard that only ever fires is
  decoration; pair every "must not appear" with a "legitimate case
  still passes".
- **A mocked callback tests the component and skips the wire.** For
  every handler a test mocks, one test drives the real one end to
  end.
- **A fixture that does not match the interface proves nothing.**
  Read the response model; mock the interface, not the model.
- **A concurrency test must force the interleaving,** not hope for
  it.
- **Asserting that a word appears in source is not a test.** Source
  scans are for drift checks that count or compare sets.
- **Run coverage on the files you touched,** not the total. An
  average hid a file that had never executed.
- **When a cost cannot be caught behaviourally, pin the mechanism**
  at source level and say in the test why no behavioural control
  could work.
- **When a guarantee cannot be observed, delete the test and write
  down why** where the test would have been.

## 4. Measure, and name the measurement.

**"This is faster" is not a finding. "858 ms and 66 MB became 49 ms
and 3.2 MB" is.** Measure the thing that will bite, separate the
steps when a result surprises you, and state what a measurement does
not cover.

**Measure the cost of a fix, not only the size of the problem.** A
decision not to fix something is a claim about cost. "This touches
every read path" with no number in it is an estimate dressed as a
finding; counting the call sites takes one search.

**Numbers in a commit message are measurements from that session.**

## 5. Say what is still open.

**Nothing ships without an honest account of what it does not do** —
in the commit, in the code, and in the plan.

Every non-trivial file ends with a notes section for the next person
or agent: what is resolved and why, and what is deferred and why. A
deferred entry is a found, considered decision, not a hidden gap.
Resolved entries are kept as the project's memory.

**A known limitation names what would close it,** or a test that
asserts it is still true.

## 6. No speculative code.

**Do not build for a caller that does not exist.** Extract a shared
abstraction at the second real caller, not the first.

**Remove duplication only when two copies could silently drift** on a
real correctness or security property. Similar-looking text is not
duplicated logic.

**A workaround that works is still the wrong answer.** Look for code
that compensates for a structure rather than changing it, and audit
for it periodically.

## 7. Read the signature before you call it.

Invented method names and guessed argument shapes are the most
common avoidable error. Check the method exists, the argument order,
the return type, and that the fixture has the data the assertion
needs. **When an exact-match edit fails twice, read the file** — it
is not what you think it is.

## 8. Do not truncate the output you are diagnosing from.

Keeping only the last line discards the name of what failed. When
something fails unexpectedly, run it again without the pipe before
concluding anything. And read the log you already have.

## 9. Layers are enforced, not aspirational.

If an architecture says one layer must not reach another, a tool
fails when it does. A needed exception is named individually; the
check is never weakened.

## 10. Comments explain why.

A comment that restates the code is noise. **Write down the wrong
turn** — what was tried first and why it failed — and correct the
comment when the code changes.

## 11. One change per commit; the message carries the reasoning.

Subject at most 72 characters, body wrapped at 72. Say why, what was
measured, what was tried and rejected, what is still open, and which
controls fired. Someone reading the history in a year should
reconstruct the decision without finding you.

**Squash within a change, never across.** A bug in work already
handed over gets its own commit; those fix commits are often the most
useful documentation in the repository.

## 12. Run the feature before committing it.

**Tests passing is not the feature working.** Exercise it through the
real entry point, against the real configuration and the real data,
with a throwaway probe. The commit says what was run, not only what
passed.

**Anything visual is verified in a real browser.** The simulated
document unit tests run against computes no layout, stacking or
cascade; a fully green suite has shipped a visibly broken screen.

## 13. Know what the environment hides.

A container is not the target machine. Runtime version, installed
tools, timezone, free disk, and whether a running process has
reloaded new code have each caused a real failure. Pin a non-UTC
timezone in tests — one behind UTC, one observing daylight saving.

## 14. Hand over everything, in the same block.

If a change needs a rebuild, a restart, a migration or a new
dependency, that step goes in the same block as the command that
applies the change. **A step written as prose beside a command block
is a step that gets skipped.** A stale server does not always return
an error; a new field arriving empty renders as a correct-looking
empty state.

## 15. Delivery without push access.

Agents working here have no write access to the remote. Work is
delivered as patches.

- **One patch per commit, numbered, applied in order.** A combined
  patch is all-or-nothing and buries each commit's reasoning.
- **Dry-run every patch against a fresh clone of the remote head**
  before presenting it, and check the commit count in it.
- **Clear stale patches from the output folder** before generating,
  so a stale one cannot be applied.
- **Confirm the previous patch landed** before starting the next.
- **Check what a hard reset would discard before running one.** It
  is silent.
- **A failed patch application leaves state behind** that blocks
  every later one; clear it before the next attempt, and treat an
  application as done only when the head commit has changed.
- **Generated state is never tracked.** A patch cannot delete a file
  that changes on its own.
- **Install from the lockfile exactly** except when deliberately
  changing a dependency, in which case the lockfile change belongs in
  that commit.

## 16. Working with the owner.

- **Give a synopsis before a run of work**: what is about to be
  built and two or three specifics about how.
- **For anything larger than one commit, agree the shape first.**
  Changing a design before building is much cheaper than after.
- **Say what to look at and where to click** when handing over
  anything visible.
- **Write scripts for anything the owner must verify by hand,** which
  print the steps and say what a failure means.
- **Say plainly when something cannot be tested,** rather than
  shipping a test that passes vacuously.
- **Ask for the fact rather than guessing at it** — which browser,
  the full output, the server log.
- **When a fix does not work twice, change the approach,** not the
  fix.
- **Flag when context is running low** before starting something that
  cannot be finished well.

## 17. Prefer the boring order.

When a list holds one interesting item and several dull blocking
ones, the dull ones come first. Order work by what depends on what,
and write that ordering down separately from the grouping by subject.

## 18. One list of what is open.

There is exactly one backlog. Several planning documents each keeping
their own list drift apart; in Elysium an entry marked "blocking
everything below it" had been fixed weeks earlier and no list said
so. Design documents hold reasoning; the backlog holds the work.

---

# Part three: engineering standards

How code is written, arranged, checked and tested. Each standard was
researched against current practice (H7). Named tool choices live
only in dependency manifests and tool configuration (H1).

## E1. Standard library first, then the canonical library, then further.

**Use the language's standard library and official toolchain
whenever they do the job.** Only when they do not, take the
canonical third-party library: the one the language's community
treats as the default, widely adopted and actively maintained, and
preferably governed by a foundation or community rather than one
company. It must still pass H4a.

**Only when neither exists, or the canonical choice fails H4a,
compare non-canonical libraries with writing our own.** Record the
comparison with the decision. When the options are roughly balanced,
the owner decides (H8).

**Never write our own** cryptography, TLS, authentication protocols,
random numbers for security, or parsers of untrusted formats where a
vetted library exists. Home-made security code fails in ways its
tests do not show.

Why: every dependency is attack surface, maintenance and exit risk.
A widely used library's credential-stealing releases in 2026 reached
its users through the supply chain (H4a); the standard library and
the canonical choice carry the least of that risk.

*Test:* every entry in the dependency register records which rung it
sits on and why a lower rung would not do.

## E2. One repository; every component its own project.

**Urshanabi lives in one repository.** Practice favours one
repository when a small team owns most services and changes cross
service boundaries, since one change and one review cover the whole
of it; one repository still lets every service build and deploy on
its own.

**Every component is its own project**: its own directory, package
manifest, tests, scripts (E5), container build and a short README
saying what it is and which contracts it speaks. It can be built,
tested and shipped as if the rest of the repository did not exist.

**Dependencies are shared per language, not per component**, as each
language's maintainers advise. The services language keeps one module
for the whole repository: its maintainers say one repository with one
module is almost always simpler, and several modules need awkward
redirections to share code. The systems language keeps one workspace,
whose members share a lockfile and a build directory so shared
libraries are compiled once. A language's projects are split only
where their requirements genuinely conflict: the agent and the
pipelines are separate projects, each with its own lockfile, because
the orchestrator's installation constraints and the agent's model
libraries pull in different directions. Diamond conflicts arise only
inside one program, so sharing versions between separate services
costs nothing in isolation.

**Isolation is enforced at the code boundary, never left to
convention.**

- Services-language code lives under each service's `internal/`
  directory, which that language's compiler forbids any other
  service to import.
- A systems-language service may depend on `libs/` only, never on
  another service; a check fails the build otherwise, because that
  language has no compiler-enforced equivalent.
- The interface depends on nothing but the client generated from
  `contracts/`, and develops against a mock generated from the same
  contracts.
- `script/check-boundaries` fails if any manifest reaches from one
  component into another.

**A contract exists in exactly one place.** No copy of a contract
file lives anywhere outside `contracts/`, and generated code is
produced at build time rather than committed (R-11). If generated
code is ever committed, CI regenerates it and fails on any
difference.

A heavyweight multi-language build system is adopted only if
measurement shows the native tools cannot keep builds fast.

**Layout:**

```
AGENTS.md  RULES.md  ROADMAP.md  README.md  LICENSE
<systems-language workspace manifest and lockfile>
<services-language module manifest and checksums>
script/            top-level scripts (E5)
contracts/         every service, event and error contract
libs/ontology/     the one ontology library (R-51)
services/<name>/   one project per deployable service
pipelines/         pipeline definitions and transformations
ui/                the interface and its design system
conformance/       the behaviour specification and its suite
e2e/               the few end-to-end journeys
load/              load generation against the objectives (R-66)
deploy/            cell, bundle and release definitions
docs/              design documents and the research behind them
```

**`AGENTS.md`** at the root tells coding agents how to work here and
points them to this file, as other large repositories now do.

## E3. Every language passes the same six gates.

1. **Format** — the language's official formatter; unformatted code
   fails.
2. **Types** — the compiler or type checker at its strictest setting.
   Where typing is gradual, the canonical checker in strict mode is
   the authoritative gate, whatever faster checker an editor uses,
   and unannotated code fails.
3. **Lint** — the canonical linter, warnings treated as errors, its
   policy declared once per workspace in the manifest rather than
   scattered through files or command lines.
4. **Dead code** — unused code, exports and dependencies fail the
   build.
5. **Audit** — known vulnerabilities and licences (H4a) are checked
   on every build.
6. **Test** — the language's test runner, with race or sanitizer
   checks wherever the language has them.

**Per language:**

- **Systems language.** Unsafe code is forbidden except in named,
  reviewed modules. No panicking shortcut in non-test code: a panic
  in a service is a crash an attacker can trigger. The minimum
  compiler version is pinned. Every parser of external input has a
  fuzz target. Mutation testing runs periodically on the core
  library to prove the tests catch real faults.
- **Services language.** Every test run uses the race detector. No
  error is ever left unchecked. The standard library comes first for
  HTTP, logging, TLS, JSON, testing and fuzzing.
- **Agent and pipeline language.** Typed throughout. Installs are
  locked and hash-pinned.
- **Frontend language.** The strictest compiler options; no escape
  hatch to untyped values; unhandled promises are errors.
- **SQL** is linted and formatted like any other code.
- **Contracts** are linted, formatted and checked for breaking
  changes (roadmap R-06).

## E4. Test in layers, weighted toward the boundaries.

The established strategy for microservices has five layers; current
practice weights them toward integration, because most microservice
faults occur where services meet.

1. **Unit** — pure logic, fast; property-based wherever inputs are
   wide.
2. **Integration** — each adapter against the real dependency in an
   ephemeral container, never against a mock of infrastructure we
   run ourselves.
3. **Component** — one whole service, its collaborators replaced by
   test doubles that speak the real contract.
4. **Contract** — consumer-driven expectations verified in the
   provider's build (R-78), with schema breaking-change checks (R-06).
5. **End-to-end** — few, covering only critical journeys, in an
   ephemeral environment per change. Never the first line of defence.

**Where each layer lives:**

- **Unit** — in each project, beside its code; nothing external.
- **Integration** — in each project, marked as integration; real
  dependencies in temporary containers.
- **Component** — in each project; the whole service, collaborators
  replaced by test doubles.
- **Contract** — the consumer records expectations and the provider
  verifies them in its own build; the expectation files live in
  `contracts/`.
- **Conformance** — `conformance/`; the whole system, black-box.
- **End-to-end** — `e2e/`; a temporary full environment.
- **Load** — `load/`; against the scale objectives.

**Tests are discovered, never listed.** Each language's own tool
finds every test in a component, and the top-level script finds every
component by walking the tree. Of four well-known microservice
repositories inspected, three had tests that never ran or failures
that were swallowed: a fixed list of directories that omitted
existing tests, unit tests no build or pipeline ever invoked, and a
loop whose exit status was only its last iteration's. The one that
discovered its tests did not. A component whose test files are not
all collected fails.

**Slow layers carry a marker.** Integration tests are marked so the
fast layers run on their own; the marker, not a list, decides what
runs where.

Alongside the layers: the black-box conformance suite (R-07); fuzzing
of every external input; mutation testing; fault injection for
deadlines, breakers and retries (R-83); load against the scale
objectives (R-66); security tests for every denial property; and the
agent evaluations (R-28, R-72). Every test declares its prerequisites
and skips by name when one is missing (R-02).

## E5. Same scripts everywhere; the top-level script only delegates.

**Every project has the same script names**, so anyone can build and
test any component without learning it first — a convention one large
software company adopted across all its projects:

- `script/bootstrap` — installs the project's dependencies.
- `script/test` — runs its gates (E3) and its fast test layers.
- `script/test-integration` — runs the tests that need real
  dependencies in containers.

**The top-level `script/` directory holds the same names, plus
`cibuild` for CI.** Each top-level script does only what a person
could do by hand: `script/each-component` walks every component and
runs that component's script of the same name inside its directory,
and `script/test` first runs every repository check — each executable
`script/check-*`, discovered, never listed. CI runs exactly
`script/cibuild`, which a person can run too; the CI definition does
nothing else.

**How components are found.** A top-level directory with its own
README is a component; one without is a group, whose subdirectories
are components (`libs/`, `services/`); `docs/` and `script/` are
neither. A component holding only Markdown files is planned, and is
reported rather than run. A top-level directory that is neither a
component nor a group fails.

**Tools are pinned.** `script/tools.lock` records every tool's
version, download address and SHA-256 for each supported platform,
each cross-checked against an independent published record when
pinned. `script/bootstrap` installs them into the ignored `.tools/`
directory and refuses any download whose checksum does not match.

**The top-level scripts are plain POSIX `sh`**, linted by the
canonical shell linter in POSIX mode, and they:

- change to the repository root first, so they behave the same from
  any directory;
- find components by walking the tree, never from a fixed list;
- treat a component holding anything but Markdown without the named
  script as a failure, not a skip, so code can never arrive untested;
- run every component even after a failure, then report every
  failure together — so they deliberately do not stop at the first
  error;
- fail if they found no components at all, or `script/test` found no
  repository checks, so a run that tested nothing can never pass;
- never place a command whose exit status matters on the left of a
  pipe, because the option that propagates such failures entered the
  POSIX standard only in its 2024 edition and older shells lack it.

**The scripts are tested themselves.** `script/check-scripts` builds
throwaway repositories and breaks every gate on purpose — a failing
component, code without a script, an empty tree, an unknown
directory, a failing check, a planted name, a missing citation, a
component reaching into another, a flawed script, a checksum
mismatch — and expects failure every time, with matching cases that
must pass so a gate cannot simply fail everything.

Why: Elysium's lint script reported success while a gate failed, and
three of the four repositories inspected had the same class of fault
(E4). A convenience script that can pass while testing nothing is
worse than none.
