package server

import "net/http"

// securityHeaders sets, on every response the gateway sends, the headers
// that stop a browser from sniffing a type, framing the response,
// sending a referrer or keeping a copy (conformance HDR-01, HDR-02). The
// gateway serves data, never pages, so its content policy allows
// nothing, and nothing it answers is safe to cache: every response may
// be specific to its caller.
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; base-uri 'none'")
		h.Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
