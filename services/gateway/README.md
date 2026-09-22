# Gateway

The front door: login, sessions, the web API and the endpoint for the
open agent-tool protocol. Every state-changing request can be retried
safely.

- **Language role:** services language
- **Holds:** sessions
- **Speaks to:** browsers, API clients and external agents
- **Roadmap:** R-36, R-44 to R-48, R-99, R-106

**Status:** walking skeleton. It serves `GET /v1/builds`, answering
with its own build identity and the query service's, and gives every
request an id that follows it into the query service and both logs.

- `script/generate` generates the contract code, never committed.
- `script/test` runs the services-language gates
  (`script/services-gates`).
- `script/image` builds the container image from `Containerfile`: the
  static binary alone, run unprivileged, the same for the same commit.
- `script/test-integration` starts the real binary and requires a plain
  502, carrying the request id, when no query service answers.
