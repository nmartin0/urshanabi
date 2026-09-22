# Contracts

Every service, event and error contract in Urshanabi, and the only place
any of them exists (`RULES.md` E2). Contracts are written in a
schema-first binary interface language; services call each other over
the standard RPC protocol for it, and the gateway also serves browsers
over a browser-compatible variant (roadmap R-06). Clients are generated
from here at build time and never committed.

## Guardrails

Decided by the owner, so that no single company can strand Urshanabi
(R-06). `script/test` enforces each one:

- **Core syntax only.** Every file uses the language's version 3 syntax,
  which every generator in our four languages reads.
- **No outside imports** beyond the language's own well-known types. The
  standard error status is carried by the RPC protocol itself, so our
  contracts define only what goes inside it.
- **No hosted services.** No dependency on a hosted schema registry and
  no remote code-generation plugins: the toolchain runs locally, pinned
  like every other tool.
- **Explicit compatibility rules.** The breaking-change rules are the
  strictest category the toolchain offers, covering wire, JSON and
  generated source, so another tool could enforce the same rules.

## Conventions every service follows

- **Errors** carry an `urshanabi.common.v1.ErrorDetail`. Only its
  `safe_args` may be logged. Services never send sentences meant for
  people: the edge renders each `reason` in the reader's language
  (R-129).
- **Time** is UTC (R-128): every instant uses the standard UTC
  timestamp type. A future civil time is carried as its local
  date-time and timezone identifier; a calendar date as a date.
- **Request id** (R-08): every request carries an `x-request-id` header;
  a service generates one if absent and passes it on every call it
  makes.
- **Idempotency** (R-106): every state-changing request carries an
  `idempotency-key` header, as the standards-track draft describes.

## Checks

`script/test` formats, lints and checks the guardrails, then compares
the contracts with the previous commit (or, in CI, with the commit the
change is based on) and fails on any breaking change.

Every application error carries an `urshanabi.common.v1.ErrorDetail`.
Its reason and safe arguments may be logged; its unsafe arguments, and
the error's free-text message, never are, since either could echo
customer data (roadmap R-06). An error without a detail is taken to come
from the transport itself.
