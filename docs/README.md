# Docs

The architecture of Urshanabi and the research behind it. `ROADMAP.md`
holds the work; these files hold the reasoning. Where the two disagree,
the roadmap is right and the document here is stale.

- `architecture.md` — what Urshanabi is made of and how the parts fit.
- `ontology-format.md` — how the ontology is written down: the open
  specification of Urshanabi's own format.
- `compliance.md` — every control mapped to the baselines buyers use,
  each implemented one proved by a test.
- `writers.md` — the single writer of every store, and how the store
  enforces it.
- `scale.md` — what a single cell must sustain, each objective
  awaiting or citing its load test.
- `research/technology-review.md` — every technology role checked
  against the newest alternatives, and each language checked against the
  libraries it uses.
- `research/repository-structure.md` — how the repository is arranged,
  and the evidence for it.
- `research/testing.md` — how Urshanabi is tested, and why.
- `research/security.md` — untrusted machines, post-quantum
  cryptography, the service mesh, inference channels and compliance.
- `research/federation.md` — peering between installations and searching
  systems where they live.
- `research/elysium.md` — what the owner's earlier prototype taught,
  measured.

## What is deliberately not here

Everything in this repository follows `RULES.md` H1, so these files
describe technologies by role and never by name. Two records must name
things and therefore live outside the repository: the named dependency
register, which maps every role to the product chosen for it, with
licence, governance and verification status; and the work order for
Elysium's own maintainers, which belongs to Elysium. Names appear in
this repository only in dependency manifests and tool configuration,
once a choice is built.
