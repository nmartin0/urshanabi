// Package caller says where a request came from (roadmap R-59).
//
// A client-address header is a claim by whoever sent it, so it counts
// only from a proxy we have said we trust. From anyone else it is
// ignored entirely and the connection's own address is used, because an
// address taken on trust can be set to anything: a caller could forge
// another tenant's address in our logs, or hide its own.
package caller

import (
	"net"
	"net/http"
	"net/netip"
	"strings"
)

// Forwarded is the header a proxy uses to pass on the address it saw.
const Forwarded = "X-Forwarded-For"

// Trusted reads the proxies we trust from a configuration value: a list
// of addresses or ranges, separated by commas. Anything unreadable is
// left out, so a mistake in configuration trusts less, never more.
func Trusted(configured string) []netip.Prefix {
	var trusted []netip.Prefix
	for _, entry := range strings.Split(configured, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if prefix, err := netip.ParsePrefix(entry); err == nil {
			trusted = append(trusted, prefix)
			continue
		}
		if address, err := netip.ParseAddr(entry); err == nil {
			trusted = append(trusted, netip.PrefixFrom(address, address.BitLen()))
		}
	}
	return trusted
}

// Of returns the address to record for r: the connection's own address,
// unless it comes from a proxy we trust, in which case the address that
// proxy passed on.
func Of(r *http.Request, trusted []netip.Prefix) string {
	peer := address(r.RemoteAddr)
	if !peer.IsValid() || !among(peer, trusted) {
		return text(peer, r.RemoteAddr)
	}
	// Read the chain from the right, stepping over proxies we trust:
	// the first address that is not one of ours is the caller.
	passed := strings.Split(r.Header.Get(Forwarded), ",")
	for i := len(passed) - 1; i >= 0; i-- {
		candidate, err := netip.ParseAddr(strings.TrimSpace(passed[i]))
		if err != nil {
			continue
		}
		if among(candidate, trusted) {
			continue
		}
		return candidate.String()
	}
	return text(peer, r.RemoteAddr)
}

// address reads the address out of a "host:port" pair.
func address(remote string) netip.Addr {
	host, _, err := net.SplitHostPort(remote)
	if err != nil {
		host = remote
	}
	parsed, err := netip.ParseAddr(host)
	if err != nil {
		return netip.Addr{}
	}
	return parsed.Unmap()
}

// text is the address as it should be recorded.
func text(peer netip.Addr, remote string) string {
	if peer.IsValid() {
		return peer.String()
	}
	return remote
}

// among says whether address falls inside any of the ranges.
func among(address netip.Addr, ranges []netip.Prefix) bool {
	for _, allowed := range ranges {
		if allowed.Contains(address) {
			return true
		}
	}
	return false
}
