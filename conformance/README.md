# Conformance suite

The behaviour specification (`BEHAVIOURS.md`) and the black-box suite
that checks it, reaching each system through a thin driver.

- **Language role:** POSIX shell for now, like the journeys; the long-term
  language is still to be decided
- **Holds:** nothing
- **Speaks to:** each system under test, through its driver
- **Roadmap:** R-07

**Status:** the suite runs, test-first. `script/test-integration`
starts the stack from its container images and runs every executable
in `properties/`, each reaching Urshanabi only through
`drivers/urshanabi`, which runs two replicas of every service and
alternates requests between them (roadmap R-18); the run fails unless
both answered. A property exits 0 when it holds, 1 when it does
not, and 2 when Urshanabi cannot yet be held to it, naming the roadmap
item it waits for. Each property is seen to fail before it passes.

The object-storage compatibility suite (roadmap R-139) runs in the same
pass: every store in `stores/drivers/` is started by its driver, and
every check in `stores/checks/` runs against it through the signing
helper `stores/request`. Each run generates its own secret; none is
committed.

The table catalog's checks run in the same pass too: every catalog in
`catalogs/drivers/` starts over its own development store, and every
check in `catalogs/checks/` runs against it through `catalogs/request`.
`catalogs/runtime` finds the runtime the catalog needs, version 21 or
newer, wherever the machine keeps it.

