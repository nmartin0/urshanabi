"""HDR-02: every response specific to its caller forbids caching (conformance/BEHAVIOURS.md)."""

import support

for method, path in support.routes():
    answered = support.request(method, path)
    if answered.header("cache-control") != "no-store":
        support.broken(f"{method} {path} does not forbid caching")
support.holds()
