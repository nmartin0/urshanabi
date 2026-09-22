package requestid

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOnlyPlainIdentifiersAreValid(t *testing.T) {
	for id, want := range map[string]bool{
		"req-42_a.b":             true,
		"":                       false,
		"x\nforged line":         false,
		strings.Repeat("a", 129): false,
	} {
		if got := Valid(id); got != want {
			t.Errorf("Valid(%q) = %v, want %v", id, got, want)
		}
	}
}

func TestFreshIdsAreValidAndDistinct(t *testing.T) {
	a, b := New(), New()
	if !Valid(a) || a == b {
		t.Fatalf("fresh ids %q and %q", a, b)
	}
}

func TestTheMiddlewareGivesEveryResponseAnId(t *testing.T) {
	var seen string
	h := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = FromContext(r.Context())
		http.NotFound(w, r)
	}))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/anything", nil))
	if id := w.Header().Get(Header); !Valid(id) || id != seen {
		t.Fatalf("response id %q, route saw %q", id, seen)
	}
}

func TestACallersValidIdIsKeptAndAnInvalidOneReplaced(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set(Header, "from-caller")
	if got := FromHTTP(r); got != "from-caller" {
		t.Fatalf("got %q", got)
	}
	r.Header.Set(Header, "bad id")
	if got := FromHTTP(r); got == "bad id" || !Valid(got) {
		t.Fatalf("an invalid id was kept: %q", got)
	}
}
