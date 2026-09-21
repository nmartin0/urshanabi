# Federation

What "federated" means in practice, the precedent for each meaning, and
how Urshanabi applies it (roadmap Phase 7). Researched in September
2026.

## Three meanings

Today Urshanabi reads only its own curated store. "Federation" covers
three different things, easily confused:

1. **Federated search** — querying other systems' stores where they
   live, without copying them first.
2. **Peering** — separate installations sharing selected objects with
   each other.
3. **Federated identity** — single sign-on across organisations.
   Urshanabi already has this (R-49); it only shares the word.

## Federated search

The leading platform lets users search external systems in place and
promote the records they need into its ontology, with every datum tied
to its source and access controlled down to individual attributes.

The largest precedent is the US national health-data exchange framework:
a network of networks that organisations join once to exchange with all
the others. A discovery request is broadcast across the networks and
matching records are retrieved from wherever they live. Every request
must state its purpose, users are identity-proofed to SP 800-63 identity
and authenticator assurance level 2, and the framework carried about 500
million records by February 2026.

**Urshanabi (R-118)**, decided by the owner: opt-in per source; results
marked external and unverified and not actionable; filtered by the
source's own permissions and by a classification ceiling per source; no
count or result revealing anything the user cannot see; a stated purpose
on every request; identity-proofed federated users; time limits, rate
limits and audit per source. Bringing a record in runs it through the
pipelines, so everything in the ontology stays curated. The value is
data minimisation: search everywhere, copy only what is needed,
including sources that may be searched but not copied.

## Peering

The leading platform peers objects between its own installations over
connections that each declare which types may flow, in which direction,
and a classification ceiling; peering keeps working over low-bandwidth
or disconnected links by queueing, and both sides must keep the object
type's properties in step.

The open precedent is the data space, defined in ISO/IEC 20151 as a
governance framework for trustworthy data sharing through agreed
policies, semantic models, protocols and processes. Its open protocol
and a companion trust protocol have been submitted for international
standardisation, and a foundation-governed reference connector
negotiates contracts, enforces usage policies and audits exchanges in
production ecosystems; one national data-exchange layer is adopting the
protocol.

**Urshanabi (R-117)**, decided by the owner: designed now, built later,
over the open protocol, so Urshanabi can peer with any conforming
participant rather than only itself — an improvement on the leader's
closed approach. Each connection declares types, direction and a
classification ceiling; usage conditions travel with the data;
conflicting values reconcile under R-102; classifications follow the
data (R-76). Designed in from the start: identifiers translatable
between installations, portable labels, and nothing in the data model
that assumes a single installation.
