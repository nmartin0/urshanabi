# Model gateway

The only component that talks to model providers. Holds no long-lived
keys, reaches only listed provider endpoints, and loads no third-party
plugins.

- **Language role:** services language
- **Holds:** short-lived model-provider credentials
- **Separate because:** it alone holds the model providers' keys
  (R-80)
- **Speaks to:** model providers, and the agent
- **Roadmap:** R-31, R-87

**Status:** planned; nothing is built yet. When work starts, this
project gains `script/bootstrap`, `script/test` and
`script/test-integration` (`RULES.md` E5).
