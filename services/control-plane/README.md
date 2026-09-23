# Control plane

Manages configuration and tenants, compiles the ontology by calling
the ontology library, and signs configuration bundles. Keeps signing
keys away from the front door.

- **Language role:** services language
- **Holds:** the signing keys
- **Separate because:** it alone holds the signing keys (R-80)
- **Speaks to:** cells, through signed bundles
- **Roadmap:** R-51, R-52, R-73

**Status:** planned; nothing is built yet. When work starts, this
project gains `script/bootstrap`, `script/test` and
`script/test-integration` (`RULES.md` E5).
