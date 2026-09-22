package logging

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

// The logged hour must be the UTC hour from the clock, whatever the
// host's zone; the tests run under three zones (roadmap R-14).
func TestLogTimestampsAreUTCWhateverTheHostZone(t *testing.T) {
	var out bytes.Buffer
	before := time.Now().UTC().Hour()
	New(&out).Info("probe")
	after := time.Now().UTC().Hour()

	var line struct{ Time string }
	if err := json.Unmarshal(out.Bytes(), &line); err != nil {
		t.Fatal(err)
	}
	stamp, err := time.Parse(time.RFC3339Nano, line.Time)
	if err != nil {
		t.Fatal(err)
	}
	if _, offset := stamp.Zone(); offset != 0 || line.Time[len(line.Time)-1] != 'Z' {
		t.Fatalf("not UTC: %s", line.Time)
	}
	if h := stamp.Hour(); h != before && h != after {
		t.Fatalf("logged hour %d is not the UTC hour (%d to %d)", h, before, after)
	}
}
