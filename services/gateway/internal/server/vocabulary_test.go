package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"  // name-ok
	"google.golang.org/grpc/status" // name-ok

	"urshanabi/services/gateway/internal/failure"
	commonv1 "urshanabi/services/gateway/internal/gen/urshanabi/common/v1"
)

// Every failure a caller can provoke answers in our own vocabulary and
// says nothing else (roadmap R-22).
func TestEveryFailureAnswersInOurOwnVocabulary(t *testing.T) {
	leaky, err := status.New(codes.Internal, "sql: no rows in *pgconn.PgError").WithDetails(&commonv1.ErrorDetail{
		Reason:     "INDEX_UNAVAILABLE",
		UnsafeArgs: []*commonv1.ErrorArgument{{Name: "query_text", Value: "SELECT secret"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	kinds := map[string]bool{failure.Unreachable: true, failure.Refused: true, failure.Empty: true, failure.Misconfigured: true}
	for _, c := range []struct {
		what   string
		fake   *fakeQuery
		method string
		path   string
	}{
		{"a downstream failure", &fakeQuery{detail: leaky.Err()}, "GET", "/v1/builds"},
		{"a downstream service that is not there", &fakeQuery{fail: true}, "GET", "/v1/builds"},
		{"an unknown route", &fakeQuery{}, "GET", "/nowhere"},
		{"a method the route refuses", &fakeQuery{}, "POST", "/v1/builds"},
	} {
		web := start(t, c.fake)
		req, _ := http.NewRequest(c.method, web.URL+c.path, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("%s: %v", c.what, err)
		}
		body := read(t, resp)
		var answered map[string]string
		if err := json.Unmarshal([]byte(body), &answered); err != nil {
			t.Errorf("%s: the answer is not an object of strings: %s", c.what, body)
			continue
		}
		if len(answered) != 2 || !kinds[answered["kind"]] || answered["request_id"] == "" {
			t.Errorf("%s: answered %s, which is not a kind and a request id", c.what, body)
		}
		// Nothing in the answer may name a type, a library or a language.
		for _, leak := range []string{"rpc", "grpc", "sql", "pgconn", "Error", "exception", "SELECT", "*", "internal/"} { // name-ok
			if strings.Contains(body, leak) {
				t.Errorf("%s: the answer contains %q: %s", c.what, leak, body)
			}
		}
	}
}

// A caller cannot put another address in our logs (roadmap R-59).
func TestAForgedCallerAddressNeverReachesTheLog(t *testing.T) {
	var log syncBuffer
	web := startLogging(t, &fakeQuery{}, &log)
	req, _ := http.NewRequest("GET", web.URL+"/v1/builds", nil)
	req.Header.Set("X-Forwarded-For", "198.51.100.23") // name-ok
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	said := log.String()
	if strings.Contains(said, "198.51.100.23") {
		t.Errorf("the log records the address the caller claimed: %s", said)
	}
	if !strings.Contains(said, `"caller":"127.0.0.1"`) {
		t.Errorf("the log does not record the address the connection came from: %s", said)
	}
}
