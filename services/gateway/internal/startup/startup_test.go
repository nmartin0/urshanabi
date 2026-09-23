package startup

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"urshanabi/services/gateway/internal/logging"
)

func TestAnUnreachableDependencyIsReportedByName(t *testing.T) {
	var out bytes.Buffer
	missing := Report(context.Background(), logging.New(&out), time.Second,
		Dependency{Name: "query", Where: "127.0.0.1:50051", Ask: func(context.Context) error {
			return errors.New("connection refused")
		}},
		Dependency{Name: "catalog", Where: "127.0.0.1:8181", Ask: func(context.Context) error { return nil }},
	)
	if missing != 1 {
		t.Fatalf("%d dependencies were missing, want 1", missing)
	}
	said := out.String()
	if !strings.Contains(said, `"dependency":"query"`) || !strings.Contains(said, "127.0.0.1:50051") {
		t.Errorf("the report does not name the unreachable dependency: %s", said)
	}
	if !strings.Contains(said, `"dependency":"catalog"`) {
		t.Errorf("the report does not mention the reachable one: %s", said)
	}
	// What it cannot reach is reported, not what went wrong inside (R-22).
	if strings.Contains(said, "connection refused") {
		t.Errorf("the report repeats the transport's words: %s", said)
	}
}

func TestADependencyThatNeverAnswersIsBounded(t *testing.T) {
	var out bytes.Buffer
	started := time.Now()
	missing := Report(context.Background(), logging.New(&out), 50*time.Millisecond,
		Dependency{Name: "query", Where: "nowhere", Ask: func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		}},
	)
	if missing != 1 || time.Since(started) > time.Second {
		t.Fatalf("missing %d after %v; a startup check must be bounded", missing, time.Since(started))
	}
}
