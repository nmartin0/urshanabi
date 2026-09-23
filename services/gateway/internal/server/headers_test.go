package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// theHeaders is every security header the gateway owns, with the value
// it must carry. securityHeaders is the only place any of them is set
// (roadmap R-36).
var theHeaders = map[string]string{
	"X-Content-Type-Options":    "nosniff",
	"X-Frame-Options":           "DENY",
	"Referrer-Policy":           "no-referrer",
	"Content-Security-Policy":   "default-src 'none'; frame-ancestors 'none'; base-uri 'none'",
	"Cache-Control":             "no-store",
	"Permissions-Policy":        "geolocation=(), camera=(), microphone=(), payment=(), usb=(), interest-cohort=()",
	"Strict-Transport-Security": "max-age=63072000; includeSubDomains; preload",
}

func TestEveryHeaderOnEveryRoute(t *testing.T) {
	web := start(t, &fakeQuery{})
	for _, c := range []struct{ method, path string }{
		{"GET", "/v1/builds"},
		{"POST", "/v1/builds"},
		{"GET", "/nowhere"},
		{"DELETE", "/"},
	} {
		req, _ := http.NewRequest(c.method, web.URL+c.path, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		_ = resp.Body.Close()
		for name, value := range theHeaders {
			if got := resp.Header.Get(name); got != value {
				t.Errorf("%s %s: %s is %q, want %q", c.method, c.path, name, got, value)
			}
			if len(resp.Header.Values(name)) != 1 {
				t.Errorf("%s %s: %s is set %d times, want once", c.method, c.path, name, len(resp.Header.Values(name)))
			}
		}
	}
}

// A second place setting one of these headers is a second owner, and
// R-36 allows exactly one. The answer alone cannot show this -- two
// setters agreeing look like one -- so the source says it.
func TestOnlyOnePlaceSetsThem(t *testing.T) {
	owner := "headers.go"             // name-ok
	root := filepath.Join("..", "..") // the gateway's own source
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") { // name-ok
			return err
		}
		if strings.HasSuffix(path, owner) || strings.HasSuffix(path, "_test.go") || strings.Contains(path, "internal/gen") { // name-ok
			return nil
		}
		content, err := os.ReadFile(path) // name-ok
		if err != nil {
			return err
		}
		for name := range theHeaders {
			if strings.Contains(string(content), `Set("`+name) {
				t.Errorf("%s also sets %s; %s is its one owner (R-36)", path, name, owner)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("the gateway's source is unreadable: %v", err)
	}
}
