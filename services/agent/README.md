# Agent

Urshanabi's own agent. It handles untrusted model output, so it holds
no data secrets and acts only as the user, within its envelope. A
separate project because its model libraries conflict with the
pipelines' requirements.

- **Language role:** agent and pipeline language
- **Holds:** nothing
- **Speaks to:** the query service, the model gateway and other agents
- **Roadmap:** R-26 to R-30, R-72, R-115

**Status:** planned; nothing is built yet. When work starts, this
project gains `script/bootstrap`, `script/test` and
`script/test-integration` (`RULES.md` E5).
