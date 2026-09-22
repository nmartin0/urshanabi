"""What every property written in this language needs: the system under
test, reached only through its driver (roadmap R-07).

The driver decides where Urshanabi answers and alternates between its
replicas (R-18), so properties never address a replica themselves.
"""

import os
import subprocess
import urllib.error
import urllib.request


def driver() -> str:
    """Where the driver of the system under test is."""
    return os.environ["DRIVER"]


class Answer:
    """One response: its status and its headers, the names lower-cased."""

    def __init__(self, status: int, headers: dict[str, str]) -> None:
        self.status = status
        self.headers = headers

    def header(self, name: str) -> str:
        """The header's value, or the empty string when it is absent."""
        return self.headers.get(name.lower(), "")


def routes() -> list[tuple[str, str]]:
    """One route that answers and two the router refuses."""
    listed = subprocess.run(
        [driver(), "routes"], capture_output=True, text=True, check=True
    )
    return [(r.split(":", 1)[0], r.split(":", 1)[1]) for r in listed.stdout.split()]


def request(method: str, path: str) -> Answer:
    """Send one unauthenticated request to the next replica."""
    where = subprocess.run(
        [driver(), "base"], capture_output=True, text=True, check=True
    )
    asked = urllib.request.Request(where.stdout.strip() + path, method=method)
    try:
        with urllib.request.urlopen(asked, timeout=30) as answered:
            status = answered.status if answered.status is not None else 0
            return Answer(status, {k.lower(): v for k, v in answered.headers.items()})
    except urllib.error.HTTPError as refused:
        status = refused.status if refused.status is not None else 0
        return Answer(status, {k.lower(): v for k, v in refused.headers.items()})


def holds() -> None:
    """End a property that holds."""
    raise SystemExit(0)


def broken(why: str) -> None:
    """End a property that does not hold, saying why."""
    print(why)
    raise SystemExit(1)


def pending(why: str) -> None:
    """End a property Urshanabi cannot yet be held to, saying what it waits for."""
    print(why)
    raise SystemExit(2)
