package server

import "net/http"

// securityHeaders sets, on every response the gateway sends, every
// security header the gateway owns: this function is the one place any
// of them is set (roadmap R-36). They stop a browser from sniffing a
// type, framing the response, sending a referrer or keeping a copy
// (conformance HDR-01, HDR-02), ask for no browser capability, and
// require a secure transport. The
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
		// The gateway asks for no browser capability at all: it serves
		// data, so every powerful feature stays off (roadmap R-36).
		h.Set("Permissions-Policy", "geolocation=(), camera=(), microphone=(), payment=(), usb=(), interest-cohort=()")
		// Transport security is set here, by the service that owns the
		// answer, not by a proxy nobody documents. A browser ignores it
		// over plain HTTP, so development is unaffected.
		h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		next.ServeHTTP(w, r)
	})
}
