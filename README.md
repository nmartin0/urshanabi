# Urshanabi

Urshanabi takes an organisation's isolated data silos, cleans the data
through pipelines, and sets an ontology on the final, curated data,
presented to people and to agents alike.

## Where to start

- `RULES.md` — how work is done here, including the hard rules and the
  engineering standards. Read it first.
- `ROADMAP.md` — the one list of work, in the order the product
  depends on it.
- `AGENTS.md` — instructions for coding agents.
- `docs/` — the architecture and the research behind it.
- `conformance/BEHAVIOURS.md` — the properties every release must
  keep.

## Layout

Every component is its own project, with its own README saying what it
is for (`RULES.md` E2):

```
script/            top-level scripts that delegate to each project
contracts/         every service, event and error contract
libs/ontology/     the one ontology library
services/<name>/   one project per deployable service
pipelines/         pipeline definitions and transformations
ui/                the interface and its design system
conformance/       the behaviour specification and its suite
e2e/               the few end-to-end journeys
load/              load generation against the scale objectives
deploy/            the release package and hosted cells' configuration
docs/              architecture and research
```

## Checking your work

```sh
script/cibuild
```

runs exactly what CI runs — the same checks, the same tools at the
same pinned versions (`script/README.md`).

## Status

Phase 0 is under way. The repository checks and CI run on every push,
and the walking skeleton runs end to end: a request crosses from the
gateway to the query service and back. Every roadmap item's status says
what is done and what remains.

## Licence

Proprietary; all rights reserved. See `LICENSE`.
