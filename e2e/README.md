# End-to-end journeys

The few end-to-end journeys, run in a temporary full environment.
Never the first line of defence (`RULES.md` E4).

- **Language role:** to be decided
- **Holds:** nothing
- **Speaks to:** the whole system
- **Roadmap:** R-127

**Status:** the walking skeleton's journey, run twice: against the
services' binaries, then against their container images. A request
enters the gateway over HTTP, crosses to the query service, and returns
both builds, each reporting the commit it was built from, with the
request id in both services' logs. Each image must hold its binary
alone, run unprivileged, and rebuild to the same image. The journey is
written in POSIX shell; the language for later journeys is still to be
decided.

- `script/test-integration` builds both services and their images and
  runs the journey against each.
