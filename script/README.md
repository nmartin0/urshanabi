# Top-level scripts

The scripts here only delegate: each walks every component and runs
that component's script of the same name inside its directory, then
reports every failure together (`RULES.md` E5).

- `bootstrap` — installs every component's dependencies.
- `test` — every component's gates and fast test layers.
- `test-integration` — every component's integration tests.
- `cibuild` — what CI runs.
- `check-boundaries` — fails if any component reaches into another.

**Status:** planned; the scripts and the tests of the scripts
themselves arrive with roadmap R-127.
