"""HDR-01: every response carries the security headers (conformance/BEHAVIOURS.md)."""

import support

WANTED = {
    "x-content-type-options": "nosniff",
    "x-frame-options": "DENY",
    "referrer-policy": "no-referrer",
    "content-security-policy": "default-src 'none'; frame-ancestors 'none'; base-uri 'none'",
}

for method, path in support.routes():
    answered = support.request(method, path)
    for name, value in WANTED.items():
        if answered.header(name) != value:
            support.broken(f"{method} {path} lacks {name}: {value}")
support.holds()
