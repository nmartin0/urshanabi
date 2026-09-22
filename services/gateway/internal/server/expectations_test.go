package server

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// The gateway records what it relies on from the query service in
// contracts/expectations, which that service's own build verifies
// (roadmap R-78). This test keeps the record honest from this side: what
// the gateway reads out of an answer must be recorded there, so the two
// can never drift apart.
func TestWhatTheGatewayReliesOnIsRecorded(t *testing.T) {
	// Every field the gateway reads from a build the query service returns.
	relied := []string{
		"builds.0.component",
		"builds.0.version",
		"builds.0.source_revision",
		"builds.0.source_committed_at.seconds",
	}
	path := filepath.Join("..", "..", "..", "..", "contracts", "expectations", "gateway-on-query.json")
	content, err := os.ReadFile(path) // name-ok
	if err != nil {
		t.Fatalf("the recorded expectations are unreadable: %v", err)
	}
	var recorded struct {
		Consumer     string `json:"consumer"`
		Producer     string `json:"producer"`
		Interactions []struct {
			Expect []struct {
				Path string `json:"path"`
			} `json:"expect"`
		} `json:"interactions"`
	}
	if err := json.Unmarshal(content, &recorded); err != nil {
		t.Fatalf("the recorded expectations are malformed: %v", err)
	}
	if recorded.Consumer != "gateway" || recorded.Producer != "query" {
		t.Fatalf("the file records %q on %q", recorded.Consumer, recorded.Producer)
	}
	paths := map[string]bool{}
	for _, interaction := range recorded.Interactions {
		for _, expectation := range interaction.Expect {
			paths[expectation.Path] = true
		}
	}
	for _, field := range relied {
		if !paths[field] {
			t.Errorf("the gateway reads %s but does not record relying on it", field)
		}
	}
}
