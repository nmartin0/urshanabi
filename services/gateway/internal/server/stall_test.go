package server

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"urshanabi/services/gateway/internal/dependency"
)

// kindOf reads the kind out of a failure's answer.
func kindOf(t *testing.T, body string) string {
	t.Helper()
	var answered map[string]string
	if err := json.Unmarshal([]byte(body), &answered); err != nil {
		t.Fatalf("the answer is not an object: %s", body)
	}
	return answered["kind"]
}

// With the query service stalled, a request that needs it fails at its
// deadline, a request that does not keeps answering at once, and once
// the dependency has failed enough times it is left to recover and
// requests fail immediately instead of waiting (roadmap R-83).
func TestAStalledDependencyCostsOnlyItsOwnRequests(t *testing.T) {
	deadline := 150 * time.Millisecond
	guarded = dependency.New(8, 3, time.Minute, deadline)
	t.Cleanup(func() { guarded = nil })

	stalled := make(chan struct{})
	t.Cleanup(func() { close(stalled) })
	web := start(t, &fakeQuery{stall: stalled})

	// A request that needs the stalled service fails at its deadline.
	started := time.Now()
	resp, body := get(t, web.URL+"/v1/builds", "test-stall-1")
	waited := time.Since(started)
	if resp.StatusCode != http.StatusBadGateway || kindOf(t, string(body)) != "unreachable" {
		t.Fatalf("answered %d %s", resp.StatusCode, body)
	}
	if waited < deadline || waited > 2*time.Second {
		t.Errorf("the request waited %v, not about its deadline of %v", waited, deadline)
	}

	// A request that does not need it is unaffected.
	started = time.Now()
	resp, _ = get(t, web.URL+"/nowhere", "test-stall-2")
	if quick := time.Since(started); quick > deadline {
		t.Errorf("a request needing nothing waited %v, longer than the deadline", quick)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("a request needing nothing answered %d", resp.StatusCode)
	}

	// After enough failures the dependency is left to recover, and
	// further requests fail at once rather than waiting again.
	for i := 0; i < 2; i++ {
		get(t, web.URL+"/v1/builds", "test-stall-more")
	}
	started = time.Now()
	resp, body = get(t, web.URL+"/v1/builds", "test-stall-3")
	resting := time.Since(started)
	if resp.StatusCode != http.StatusBadGateway || kindOf(t, string(body)) != "unreachable" {
		t.Fatalf("while resting, answered %d %s", resp.StatusCode, body)
	}
	if resting > deadline/2 {
		t.Errorf("while resting, the request still waited %v", resting)
	}
}
