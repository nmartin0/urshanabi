# Gateway

The front door: login, sessions, the web API and the endpoint for the
open agent-tool protocol. Every state-changing request can be retried
safely.

- **Language role:** services language
- **Holds:** sessions
- **Speaks to:** browsers, API clients and external agents
- **Roadmap:** R-36, R-44 to R-48, R-99, R-106

**Status:** planned; nothing is built yet. When work starts, this
project gains `script/bootstrap`, `script/test` and
`script/test-integration` (`RULES.md` E5).
