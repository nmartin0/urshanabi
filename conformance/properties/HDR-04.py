"""HDR-04: no machine-readable description of the interface is served.

See conformance/BEHAVIOURS.md.
"""

import support

DESCRIPTIONS = [
    "/openapi.json",  # name-ok
    "/openapi.yaml",  # name-ok
    "/swagger.json",  # name-ok
    "/v1/openapi.json",  # name-ok
    "/docs",
    "/api-docs",
]

for path in DESCRIPTIONS:
    answered = support.request("GET", path)
    if answered.status != 404:
        support.broken(f"{path} answered {answered.status}, not 404")
support.holds()
