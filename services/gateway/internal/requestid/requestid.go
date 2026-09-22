// Package requestid gives every request an id that follows it through
// every service and log line (roadmap R-08).
package requestid

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// Header is the HTTP header carrying the request id.
const Header = "X-Request-Id"

// MetadataKey carries the request id between services.
const MetadataKey = "x-request-id"

// Valid reports whether id is a plain identifier: 1 to 128 letters,
// digits, dots, hyphens or underscores. Anything else could forge log
// lines, since ids come from callers.
func Valid(id string) bool {
	if len(id) == 0 || len(id) > 128 {
		return false
	}
	for _, c := range []byte(id) {
		ok := c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' ||
			c == '-' || c == '_' || c == '.'
		if !ok {
			return false
		}
	}
	return true
}

// New returns a fresh random id.
func New() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b) // never fails: crypto/rand panics rather than return an error
	return hex.EncodeToString(b)
}

// FromHTTP returns the caller's id if it is valid, and a fresh one otherwise.
func FromHTTP(r *http.Request) string {
	if id := r.Header.Get(Header); Valid(id) {
		return id
	}
	return New()
}

type contextKey struct{}

// Middleware gives every request an id before any route sees it -- the
// caller's if valid, a fresh one otherwise -- and sets it on every
// response, including the router's own refusals.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := FromHTTP(r)
		w.Header().Set(Header, id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKey{}, id)))
	})
}

// FromContext returns the id Middleware gave the request.
func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(contextKey{}).(string)
	return id
}
