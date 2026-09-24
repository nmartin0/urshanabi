# Top-level scripts

Everything CI does, runnable locally with one command (`RULES.md` E5):

```sh
script/cibuild    # exactly what CI runs: bootstrap, test, test-integration
```

- `bootstrap` — installs the pinned tools from `tools.lock` into
  `.tools/`, verifying every checksum, then every component's own
  bootstrap.
- `test` — every repository check (`check-*`, discovered), then every
  component's own `script/test`. It clears every variable that could
  name an external service or a credential first, so no test may
  quietly depend on one (R-02).
- `test-integration` — every component's own `script/test-integration`.
- `each-component NAME` — runs every component's `script/NAME`; the
  one place components are discovered.
- `systems-gates PACKAGE` — the systems-language gates, and
  `systems-toolchain` — its toolchain, each shared by every component
  in that language so none can drift.
- `services-gates DIR` and `services-toolchain` — the same for the
  services language.
- `server` — builds every service's image and runs the whole local
  stack with one command; `server stop` removes it (R-12).
- `report-tool-updates` — reports any pinned tool with a newer release;
  run weekly by its own workflow, never by `test`, so a new upstream
  release cannot fail anyone's push.

Repository checks, each runnable on its own:

- `check-inventory` — the asset inventory matches the repository;
  `inventory` regenerates it (R-123).
- `check-licences` — each recorded licence is the one the dependency
  itself declares (H4a, R-85).
- `check-names` — no name outside manifests and tool configuration
  (H1); configured by `names.list`, `names.allow` and `names.exempt`.
- `check-docs` — every roadmap item, behaviour property and rule cited
  anywhere exists, and every roadmap item is complete (R-15).
- `check-unused` — a library component ships with its caller, or
  records why it is ahead of it (R-33).
- `check-shell` — every POSIX script passes the shell linter.
- `check-workflows` — the CI definition is valid and free of known
  security flaws (R-86).
- `check-secrets` — no credential anywhere in the history.
- `check-boundaries` — no component reaches into another (E2).
- `check-compliance` — every implemented control in
  `docs/compliance.md` cites a passing check and self-test (R-125).
- `check-writers` — every store in `docs/writers.md` has exactly one
  writer, and says how the store enforces it (R-79).
- `check-scale` — every objective in `docs/scale.md` has a number, a
  unit and a roadmap item; a tested one cites its load test (R-66).
- `check-ontology-parsers` — only `libs/ontology` parses ontology
  definitions (R-51).
- `check-generated` — every generator runs and leaves version control
  seeing no change (R-11).
- `check-dependencies` — every dependency we chose has a register
  entry, and its licence is one the allowlist accepts (H4a, R-85).
- `check-constant-time` — every comparison of a secret uses a
  constant-time function (R-48).
- `check-decomposition` — every separately deployed service records
  why it is separate: what it alone holds, or the measurement that
  proved the split necessary (R-80).
- `check-rungs` — every direct dependency records its rung on the
  standard-library ladder, beside itself in its manifest (E1, R-127).
- `check-scripts` — the scripts' own tests: every gate is broken on
  purpose and must fail (R-127).

Supported platforms: the four pinned in `tools.lock`, covering the two
common Unix-like operating systems on x86-64 and ARM. Requirements: a
POSIX shell, `awk`, `tar`, a SHA-256 utility, and the command-line
download tool that `bootstrap` invokes.
