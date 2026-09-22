# Compliance control matrix

Every control Urshanabi implements or plans, mapped to the baselines its buyers use (roadmap R-125): the federal control catalogue (SP 800-53 Rev. 5), the secure software development framework (SP 800-218, cited by task, such as PO.3.2), the zero-trust guidance (SP 800-207 and SP 800-207A), and the technical safeguards of the US health-data law's Security Rule (45 CFR Part 164, Subpart C, cited by section, such as 164.312(b)), designed to its proposed revision.

`script/check-compliance` enforces the rules of this file: every implemented control cites evidence that exists — an executable check, and the self-test proving that check fails when it should — and every planned control names the roadmap items that will deliver it. An implemented control without evidence fails the build.

SP 800-53 references for implemented controls follow the mapping published in SP 800-218 itself wherever one exists; the others were checked against the catalogue's text. References for planned controls are confirmed again when each becomes implemented.

## Controls

| ID | Control | Baselines | Roadmap | Status | Evidence |
|---|---|---|---|---|---|
| UC-01 | Every change runs every security check in CI, and the same checks run locally from one command | SP800-53: SA-11, SA-15, SA-15(1); SP800-218: PO.4.1, PO.4.2, PW.7.2, PW.8.2 | R-01, R-127 | implemented | `script/cibuild`; self-test "a failing check fails script/test" |
| UC-02 | Every development tool is pinned by version and SHA-256, checked against an independent published record, and refused on mismatch | SP800-53: SA-15, SR-3, SR-4, SR-4(3); SP800-218: PO.3.2, PW.4.4, PW.6.1 | R-10, R-86 | implemented | `script/bootstrap`; self-test "a checksum mismatch is refused" |
| UC-03 | The CI definition runs with least privilege: its one action pinned to a commit, no permissions by default, no secrets and no privileged triggers, scanned at the strictest level | SP800-53: AC-6, CM-7, SA-15; SP800-218: PO.3.2, PO.5.1 | R-86 | implemented | `script/check-workflows`; self-test "an unpinned action fails the workflow check" |
| UC-04 | No credential is ever committed, anywhere in the history | SP800-53: IA-5(7); SP800-218: PS.1.1 | R-86, R-92 | implemented | `script/check-secrets`; self-test "a committed credential fails the secret check" |
| UC-05 | A breaking change to any service contract fails the build | SP800-53: SA-10, SA-15(1); SP800-218: PO.4.2 | R-06 | implemented | `contracts/script/test`; self-test "a committed field removal fails" |
| UC-06 | No component can reach into another's code | SP800-53: SA-8; SP800-218: PO.1.2 | R-127 | implemented | `script/check-boundaries`; self-test "a service reaching into another fails" |
| UC-07 | Every gate is itself tested, and every fault found becomes a permanent self-test | SP800-53: SA-11, SA-15; SP800-218: RV.3.2, RV.3.4 | R-127 | implemented | `script/check-scripts` |
| UC-08 | Security decisions and requirements are tracked, and no document can cite one that does not exist | SP800-53: SA-8, SA-10, SA-17; SP800-218: PW.1.2 | R-15 | implemented | `script/check-docs`; self-test "a missing roadmap item fails" |
| UC-09 | This matrix: every implemented control proves itself with evidence | SP800-53: SA-15(1); SP800-218: PO.4.1 | R-125 | implemented | `script/check-compliance`; self-test "an implemented control without evidence fails" |
| UC-10 | Errors reveal only what is safe; arguments that may carry data are never logged | SP800-53: SI-11; SP800-218: PW.5.1 | R-06, R-67 | implemented | `services/gateway/script/test` |
| UC-11 | Every timestamp is UTC, and no code reads the host's timezone | SP800-53: AU-8; 45CFR: 164.312(b) | R-14, R-128 | planned | — |
| UC-12 | Every person and every service has a unique identity, and reaches only what it has been granted | SP800-53: AC-3, AC-6, IA-2; 45CFR: 164.312(a)(1), 164.312(a)(2)(i) | R-44, R-68, R-81 | planned | — |
| UC-13 | Multi-factor authentication for every account that can reach regulated data | SP800-53: IA-2(1), IA-2(2); 45CFR: 164.312(d) | R-49 | planned | — |
| UC-14 | Every connection encrypted, with the hybrid post-quantum key exchange | SP800-53: SC-8, SC-8(1); 45CFR: 164.312(e)(1), 164.312(e)(2)(ii) | R-91, R-112 | planned | — |
| UC-15 | Every store encrypted at rest | SP800-53: SC-28, SC-28(1); 45CFR: 164.312(a)(2)(iv) | R-91 | planned | — |
| UC-16 | Sessions end after inactivity and expire | SP800-53: AC-12; 45CFR: 164.312(a)(2)(iii) | R-45, R-46 | planned | — |
| UC-17 | A tamper-evident audit record of every access and change | SP800-53: AU-2, AU-9; 45CFR: 164.312(b), 164.312(c)(1) | R-23, R-88 | planned | — |
| UC-18 | A departing person's access ends everywhere within one hour | SP800-53: AC-2, PS-4 | R-84 | planned | — |
| UC-19 | Services are segmented by identity, not network location | SP800-53: AC-4, SC-7; SP: 800-207A | R-81, R-82, R-126 | planned | — |
| UC-20 | Restoration from backup within 72 hours, tested on a schedule | SP800-53: CP-9, CP-10 | R-64 | planned | — |
| UC-21 | Vulnerability scans at least every six months, penetration tests at least yearly | SP800-53: RA-5, CA-8; SP800-218: RV.1.2 | R-65 | planned | — |
| UC-22 | A machine-generated asset inventory and network map, reviewed yearly | SP800-53: CM-8 | R-123 | planned | — |
| UC-23 | Every production build's fingerprint is published so customers can verify what runs | SP800-53: SA-8, SR-4; SP800-218: PS.2.1, PS.3.1, PS.3.2 | R-121 | planned | — |
| UC-24 | The single writer of every store is named in a checked map | SP800-53: SA-8, SA-17; SP800-218: PW.1.2 | R-79 | implemented | `script/check-writers`; self-test "a store with two owners fails" |
| UC-25 | Every store refuses writes from anyone but its owner | SP800-53: AC-3, AC-6 | R-79, R-104 | planned | — |
