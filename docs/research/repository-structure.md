# Repository structure

How the repository is arranged (`RULES.md` E2 and E5), and the evidence
for it. Researched and audited in September 2026, including four
well-known microservice repositories cloned and read in full.

## One repository

Practice favours one repository when a small team owns most services and
changes cross service boundaries, since one change and one review cover
the whole of it; separate repositories suit genuinely independent teams
that need hard isolation and their own release schedule. One repository
decides where code lives, not how it ships: services still build and
deploy independently. The cost is that one repository trades enforced
isolation for cheap coordination, so isolation has to be checked by
tests rather than left to convention.

Heavyweight multi-language build systems model every file across
languages for hermetic builds, but their setup cost pays off only at
very large scale. Each language's own toolchain, with a thin script
layer on top, is the right size for Urshanabi.

## Dependencies: shared per language

Two strategies are recognised: every project keeps its own dependencies,
or a single version of each dependency serves the whole repository. The
single-version rule prevents diamond conflicts, where two libraries in
one program need incompatible versions of a third, but in its best-known
form it depends on whole-tree upgrades, dedicated tooling and people
whose job is keeping dependencies current.

The first proposal gave every service its own manifest and lockfile. An
audit against each language's own guidance corrected it:

- **Services language.** Its maintainers say one repository with one
  module is almost always simpler, and several modules need awkward
  redirections to share code. One module; each service is still its own
  program, container and deployment.
- **Systems language.** A workspace shares one lockfile and one build
  directory precisely so shared libraries compile once; per-service
  lockfiles would have compiled the ontology library four times.
- **Agent and pipeline language.** Its package manager's documentation
  says shared workspaces do not suit members with conflicting
  requirements, and that one member may import another's dependencies
  undetected. The agent and the pipelines are separate projects.

Isolation is therefore enforced at the code boundary: the services
language's compiler forbids importing another service's internal
directory; a check forbids systems-language services depending on each
other; the interface knows only its generated client.

## Four repositories compared

| | Cloud-native reference demo | Observability showcase | Banking demo | Production log system |
|---|---|---|---|---|
| Layout | One directory per service | One directory per service | Grouped by domain, then service | Programs in one directory, code in another |
| Dependencies | Every service separate | Every service separate | One shared build for one language; separate projects for another | One module for the whole system |
| Interfaces | Central directory, but services copy it and commit generated code | Central directory; CI regenerates code and fails on any difference | No schema directory | Client library as its own versioned module |
| Unit tests | Beside the code; few | Beside the code; 15 files | Per service | Beside the code; 1,091 files |
| Multi-service tests | Deploys everything to a cluster | Top-level suite against the running system | Browser tests inside the CI configuration directory | Top-level integration suite, separated by a build marker |
| Load | Its own component | Its own component | Its own component | — |

The two demos give every service its own manifest because nearly every
service is in a different language, to show each one off. The two
layouts closest to production share one build per language, as Urshanabi
does.

## The lesson that mattered most

Three of the four had tests that never ran or failures that were
swallowed:

- the reference demo tested its services-language code from a fixed list
  of directories, so two existing test packages never ran in CI;
- the showcase had fifteen unit-test files, but only two services' tests
  ran anywhere;
- the banking demo's local test target looped over a fixed list, and a
  shell loop reports only its last iteration's status, so a failure in
  the first service was reported as success — reproduced here.

The production system discovered every package instead of listing them,
and made piped commands report failures. Hence the rules: tests are
discovered, never listed (E4); the top-level scripts discover
components, treat a missing script as a failure, report every failure,
fail when they found nothing, and are tested themselves (E5).

## The script convention

One large software company standardised the same script names across all
its projects, so anyone can set up and test any project the moment they
clone it; each script does one unit of work so they compose, and CI runs
the same test script a person runs. A top-level delegating makefile was
rejected because recursive builds are notoriously hard to order; each
language's own tool owns its dependency graph and the top-level script
stays simple. The scripts are plain POSIX `sh`, and never put a
status-bearing command on the left of a pipe, because the option that
propagates such failures entered the POSIX standard only in its 2024
edition.
