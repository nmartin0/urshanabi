# RULES

How work is done on Urshanabi, by every person and every agent.

Part one is hard rules. They have no exceptions unless the owner
authorizes one explicitly, and that authorization is recorded where
the exception lives.

Part two is the working method, carried over from Elysium, where
every rule was learned by getting something wrong. Where the mistake
is instructive it is named, because a rule without its reason gets
worked around the first time it is inconvenient.

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
repeated in any comment beside it.

**Provisional, pending the owner's confirmation:** generic technical
vocabulary and published protocol designations — HTTP, TLS, SQL,
JSON, YAML, OIDC, SAML, SCIM, RFC numbers — are treated as vocabulary
rather than product names.

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
