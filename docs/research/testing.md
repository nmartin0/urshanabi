# Testing

How Urshanabi is tested (`RULES.md` E3 and E4), and why.

## The layers

The established strategy for testing microservices, published in 2014
and still the reference, has five layers:

1. **Unit** — pure logic, fast; property-based where inputs are wide.
2. **Integration** — each adapter against the real dependency, in a
   temporary container, never a mock of infrastructure we run ourselves.
3. **Component** — one whole service, its collaborators replaced by test
   doubles that speak the real contract.
4. **Contract** — each consumer records what it relies on, and the
   provider verifies those expectations in its own build.
5. **End-to-end** — a few critical journeys, in a temporary full
   environment.

A widely cited refinement for microservices, the "honeycomb", weights
integration most, because individual services often hold little logic
and most faults occur where services meet. Urshanabi follows it: the
most effort goes to integration and contract tests, and end-to-end tests
are never the first line of defence.

## Around the layers

- **Conformance** — the behaviour properties in
  `conformance/BEHAVIOURS.md`, checked black-box against Urshanabi
  through one driver, written test-first; Elysium is the evidence that
  each property is testable, not a target (R-07).
- **Fuzzing** of every parser of external input.
- **Mutation testing** on the core library, to prove the tests catch
  real faults.
- **Fault injection** for deadlines, circuit breakers and retries
  (R-83).
- **Load** against the scale objectives (R-66).
- **Security** — every denial property, including that no count, flag or
  timing reveals hidden rows (DENY-09).
- **Agent evaluations** — rates compared, adaptive attacks run nightly
  (R-28, R-72).

## Where tests live

Unit, integration and component tests live inside each component, in its
language's usual place, with integration tests carrying a marker so the
fast layers run on their own. Contract expectation files live in
`contracts/`, so no separate broker is needed. Conformance, end-to-end
and load tests are top-level projects of their own (`conformance/`,
`e2e/`, `load/`).

## Discovered, never listed

Of four well-known microservice repositories inspected
(`repository-structure.md`), three had tests that existed but never ran,
or failures that were silently swallowed, every time because something
was listed by hand. Elysium had the same fault: its lint script reported
success while its lock-file gate failed, because that gate set a
variable the script never read. So every language's own tool discovers
the tests in a component, the top-level script discovers the components,
and a component whose test files are not all collected fails. Every gate
is proved by breaking it deliberately and expecting failure (R-127).

## Gates in every language

Format, strictest type checking, lint with warnings as errors, dead-code
detection, vulnerability and licence audit, and tests with race or
sanitizer checks where the language has them (E3). For gradually typed
languages the canonical checker in strict mode is the authoritative
gate, whatever faster checker an editor uses. Tool choices follow the
standard-library-first ladder (E1).
