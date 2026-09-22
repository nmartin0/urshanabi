package failure

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"  // name-ok
	"google.golang.org/grpc/status" // name-ok

	commonv1 "urshanabi/services/gateway/internal/gen/urshanabi/common/v1"
	"urshanabi/services/gateway/internal/logging"
)

const secret = "SELECT name FROM patients"

func logged(err error) string {
	var out bytes.Buffer
	logging.New(&out).Error("call failed", Attrs(err)...)
	return out.String()
}

func TestAnUnsafeArgumentIsNeverLogged(t *testing.T) {
	st, err := status.New(codes.InvalidArgument, "bad query: "+secret).WithDetails(&commonv1.ErrorDetail{
		Reason:     "QUERY_INVALID",
		SafeArgs:   []*commonv1.ErrorArgument{{Name: "object_type", Value: "Patient"}},
		UnsafeArgs: []*commonv1.ErrorArgument{{Name: "query_text", Value: secret}},
	})
	if err != nil {
		t.Fatal(err)
	}
	line := logged(st.Err())
	if strings.Contains(line, "patients") {
		t.Fatalf("an unsafe argument or the free-text message was logged: %s", line)
	}
	for _, want := range []string{`"code":"InvalidArgument"`, `"reason":"QUERY_INVALID"`, `"arg.object_type":"Patient"`} {
		if !strings.Contains(line, want) {
			t.Errorf("missing %s in %s", want, line)
		}
	}
}

func TestATransportErrorKeepsItsMessage(t *testing.T) {
	line := logged(status.Error(codes.Unavailable, "connection refused"))
	if !strings.Contains(line, `"code":"Unavailable"`) || !strings.Contains(line, "connection refused") {
		t.Fatalf("the transport's message was lost: %s", line)
	}
}

func TestAnErrorFromElsewhereIsLoggedAsItIs(t *testing.T) {
	if line := logged(errors.New("local failure")); !strings.Contains(line, "local failure") {
		t.Fatalf("got %s", line)
	}
}
