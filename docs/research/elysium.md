# What Elysium taught

Elysium is the owner's earlier prototype in the same category. It was
reviewed, measured and tested in September 2026. Its design does not
carry over (`RULES.md`, Purpose); its lessons do. Each lesson below was
measured or read in its code, and each is answered by a roadmap item.

## What it got right

- **328 integration tests that need no model**, from which
  `conformance/BEHAVIOURS.md` derives Urshanabi's baseline properties,
  each citing the test that proves it.
- **A drift report** that refuses a source whose column no longer
  matches its declared type, rather than silently coercing it (R-93).
- **Secrets as references** in configuration, with a missing reference
  refusing to load by name (R-92).
- **Per-source capability declarations**, so conditions a source can
  evaluate are pushed down and the rest run in process (R-69).
- **An identifier-keyed change log** and per-generation snapshot pinning
  — real inputs an index can consume (R-70).

## What it measured

- **Serving from a lake format was slow.** Point lookups took 10.96
  milliseconds from its mirror against 1.03 from the live source,
  because the mirror served them from files. Urshanabi's index is built
  for lookups (R-70).
- **Unauthenticated disk exhaustion.** Login fields had no length bounds
  and failed attempts were never reaped: ten requests with 900-kilobyte
  usernames grew its credential store from 100 kilobytes to 16.3
  megabytes (R-44, R-45).
- **Validation errors echoed request bodies**, including a submitted
  password, to an unauthenticated caller (R-67).
- **A truncation flag revealed hidden rows.** A user who could see no
  records of a link-secured type was told the search had hit its scan
  ceiling, disclosing that many hidden records existed (DENY-09, R-68).
- **A lint gate's failure was swallowed.** One gate set a variable the
  script never read, so the script reported success while that gate
  failed (R-127).
- **Seventeen unit tests failed on a fresh clone**, because they
  depended on local state they did not declare (R-02).
- **Ingestion held whole tables in memory**, about 1.2 kilobytes a row
  (R-53).

## What it read

- **No CI**: nothing ran its checks on push (R-01).
- **Session tokens stored in plaintext**, and expired sessions never
  deleted (R-46).
- **A CSRF comparison that was not constant-time** (R-48).
- **Identifiers whose type depended on the read path**: integers from
  one, strings from the other (QUERY-08, R-20).
- **A memory guard built and tested but never wired in** (R-33).
- **A model adapter with no deadline and no token accounting** (R-27).
- **No graceful shutdown and no source validation at startup** (R-60,
  R-61).
- **Documentation that had gone stale** in several places, describing
  behaviour the code had since changed (R-15).
- **Plans built around an object store** whose open-source edition was
  archived in 2026 (`RULES.md` H4a).
- **About 1,740 lines of its hardest code** — the pending-write store,
  its persistence, write log, resume and expiry — which a durable
  workflow engine replaces in Urshanabi (R-38).
