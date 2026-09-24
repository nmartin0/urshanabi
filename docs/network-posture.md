# Network posture

Decided in one place, because Elysium left these undecided and bound to
loopback for convenience (roadmap R-59).

## Where transport security ends

Inside a cell, a service speaks only to the mesh, and the mesh carries
every hop under mutual transport security (R-126). Traffic from outside
the cell ends its transport security at the cell's ingress, which is the
only component that holds a public certificate. A service therefore
never terminates a public connection itself, and never needs a
certificate of its own beyond the mesh's.

## Which proxies may say where a request came from

None, unless configuration names them. A client-address header is a
claim by whoever sent it: from a caller we do not trust it is ignored
entirely, and the connection's own address is recorded instead.
`URSHANABI_TRUSTED_PROXIES` names the addresses and ranges whose claims
count, and reading the chain steps over our own proxies to the first
address beyond them. A value that cannot be read trusts nothing, so a
mistake in configuration can only trust less.

## What a service binds to

An address in configuration, never a compiled-in default that assumes a
network. `URSHANABI_LISTEN` decides it, and it is loopback only where
nothing else should reach the service. In a cell, a service binds to the
address the mesh gives it.

## What this rules out

- A caller setting another caller's address in our logs or traces.
- A service accidentally reachable from outside its cell because its
  binding was decided in code rather than configuration.
- Transport security ending somewhere nobody documented, which is what
  made Elysium's posture impossible to state.
