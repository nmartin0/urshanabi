# Security

The research behind Urshanabi's security design, as of September 2026.

## Can a service be secure on computers we do not trust?

**Partly.** Against a malicious operator, a cloud administrator or
compromised host software, yes: hardware-isolated confidential computing
keeps a workload's memory encrypted, and remote attestation lets a
workload prove it is genuine, unmodified code on genuine hardware before
any key is released to it.

**Against someone with physical access to the hardware, no one can
promise it.** Confidential computing's threat models have always
excluded physical attacks. In 2025, researchers built a memory-bus
interposer from under 1,000 dollars of off-the-shelf parts that
extracted secret keys, including attestation keys, from fully updated
servers of both major technologies, and used them to forge attestations:
to claim a workload was protected while reading its data and returning
false results. In September 2026 an active interposer costing under 200
dollars broke the integrity of up-to-date systems.

**The strongest precedent** for running sensitive work on infrastructure
its operator does not fully trust rests on five requirements: stateless
computation, a small verifiable trusted base, no privileged runtime
access, non-targetability — compromising one machine must not let an
attacker steer a chosen user's requests to it — and verifiable
transparency, publishing every production build. It does not rely on
confidential computing alone against side channels.

Urshanabi adopts all five (R-119 to R-121) and adds a physical-trust
ceiling for every cell (R-122).

## Post-quantum cryptography

Encrypted traffic can be recorded now and decrypted later by a quantum
computer. The hybrid key exchange that defends against this was
standardised as RFC 10024 in August 2026 and is the default in major
browsers and in one services-language standard library. A national
security algorithm timeline recommends hybrid deployment now and
post-quantum-only operation by 2033, and a 2026 measurement found
government and defence adoption effectively absent. The hybrid algorithm
is not yet in validated cryptographic modules, so cells that must use
validated cryptography run classical key exchange until it is (R-112).

## The service mesh

Zero-trust guidance for cloud-native applications (SP 800-207A) moves
security from network location to identity: every service authenticated
and authorised by who it is, enforced through gateways, proxies and
workload identities, typically by a service mesh. A sidecar-less mesh
runs one proxy per node, which is efficient and supports the hybrid
exchange, but a failure of that shared proxy affects every workload on
the node, and a proxy outside a confidential workload's boundary would
see its plaintext. So standard cells use the sidecar-less mesh, while
confidential cells terminate encryption inside each protected workload
and issue certificates only to attested workloads, as an existing
confidential mesh does (R-126).

## Inference channels

Row-level security prevents direct access but not disclosure through
aggregates, counts or timing. A published study measured row-level
security leaking the size of hidden row sets through timing when the
caller's condition ran first; a mature database orders the security
condition first for that reason. Elysium's search truncation flag
revealed rows a user could not see (see `elysium.md`). Hence DENY-09 and
R-68: the security condition runs first, and nothing a caller sees
derives from rows they cannot. Small aggregates are suppressed by
default and differential privacy is available for designated types, at
the documented cost of spent budgets and reduced accuracy for small
groups (R-71).

## Compliance

- **Health data.** A proposed revision of the US health-data law's
  Security Rule, not final as of mid-2026, would make encryption at rest
  and in transit and multi-factor authentication mandatory, with
  vulnerability scans every six months, yearly penetration tests,
  restoration within 72 hours, and a yearly asset inventory and network
  map. Urshanabi designs to the proposed revision (R-123 to R-125).
- **Federal baselines.** The federal control catalogue (SP 800-53 Rev.
  5), zero-trust guidance (SP 800-207, SP 800-207A), microservice
  security guidance (SP 800-204) and the secure software development
  framework (SP 800-218) are mapped control by control, each to a test
  (R-125).

## Supply chain

In 2026 a widely used model-gateway library shipped credential-stealing
releases after its publishing credentials were stolen through a
compromised vulnerability scanner in its CI. Tools in CI are attack
surface, every dependency is chosen for licence, governance and exit
(`RULES.md` H4a), and the standard library comes first (E1).
