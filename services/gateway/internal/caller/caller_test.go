package caller

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// asking returns a request from remote, claiming claimed.
func asking(remote, claimed string) *http.Request {
	r := httptest.NewRequest("GET", "/v1/builds", nil)
	r.RemoteAddr = remote
	if claimed != "" {
		r.Header.Set(Forwarded, claimed)
	}
	return r
}

func TestAClaimFromAnUntrustedCallerIsIgnored(t *testing.T) {
	trusted := Trusted("10.0.0.0/8")
	got := Of(asking("203.0.113.7:4444", "198.51.100.23"), trusted)
	if got != "203.0.113.7" {
		t.Fatalf("recorded %q; a caller we do not trust set its own address", got)
	}
}

func TestAClaimFromATrustedProxyIsUsed(t *testing.T) {
	trusted := Trusted("10.0.0.0/8")
	got := Of(asking("10.1.2.3:4444", "198.51.100.23"), trusted)
	if got != "198.51.100.23" {
		t.Fatalf("recorded %q, not the address the proxy passed on", got)
	}
}

func TestOurOwnProxiesAreSteppedOver(t *testing.T) {
	trusted := Trusted("10.0.0.0/8, 192.0.2.1")
	got := Of(asking("10.1.2.3:4444", "198.51.100.23, 192.0.2.1, 10.9.9.9"), trusted)
	if got != "198.51.100.23" {
		t.Fatalf("recorded %q, not the caller beyond our own proxies", got)
	}
}

func TestWithNoProxyTrustedNothingIsTakenOnTrust(t *testing.T) {
	for _, configured := range []string{"", "not an address", "10.0.0.0/999"} {
		got := Of(asking("203.0.113.7:4444", "198.51.100.23"), Trusted(configured))
		if got != "203.0.113.7" {
			t.Errorf("with %q configured, recorded %q", configured, got)
		}
	}
}

func TestAnUnreadableClaimLeavesTheConnectionsAddress(t *testing.T) {
	trusted := Trusted("10.0.0.0/8")
	for _, claimed := range []string{"nonsense", "", " , ,"} {
		if got := Of(asking("10.1.2.3:4444", claimed), trusted); got != "10.1.2.3" {
			t.Errorf("with %q passed on, recorded %q", claimed, got)
		}
	}
}
