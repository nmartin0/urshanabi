# Top-level scripts

Everything CI does, runnable locally with one command (`RULES.md` E5):

```sh
script/cibuild    # exactly what CI runs: bootstrap, test, test-integration
```

- `bootstrap` — installs the pinned tools from `tools.lock` into
  `.tools/`, verifying every checksum, then every component's own
  bootstrap.
- `test` — every repository check (`check-*`, discovered), then every
  component's own `script/test`.
- `test-integration` — every component's own `script/test-integration`.
- `each-component NAME` — runs every component's `script/NAME`; the
  one place components are discovered.
- `systems-gates PACKAGE` — the systems-language gates, and
  `systems-toolchain` — its toolchain, each shared by every component
  in that language so none can drift.
- `report-tool-updates` — reports any pinned tool with a newer release;
  run weekly by its own workflow, never by `test`, so a new upstream
  release cannot fail anyone's push.

Repository checks, each runnable on its own:

- `check-names` — no name outside manifests and tool configuration
  (H1); configured by `names.list`, `names.allow` and `names.exempt`.
- `check-docs` — every roadmap item, behaviour property and rule cited
  anywhere exists, and every roadmap item is complete (R-15).
- `check-shell` — every POSIX script passes the shell linter.
- `check-workflows` — the CI definition is valid and free of known
  security flaws (R-86).
- `check-secrets` — no credential anywhere in the history.
- `check-boundaries` — no component reaches into another (E2).
- `check-scripts` — the scripts' own tests: every gate is broken on
  purpose and must fail (R-127).

Supported platforms: the four pinned in `tools.lock`, covering the two
common Unix-like operating systems on x86-64 and ARM. Requirements: a
POSIX shell, `awk`, `tar`, a SHA-256 utility, and the command-line
download tool that `bootstrap` invokes.
