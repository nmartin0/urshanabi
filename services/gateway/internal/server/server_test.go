package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	rpc "google.golang.org/grpc"                  // name-ok
	"google.golang.org/grpc/codes"                // name-ok
	"google.golang.org/grpc/credentials/insecure" // name-ok
	"google.golang.org/grpc/metadata"             // name-ok
	"google.golang.org/grpc/status"               // name-ok

	buildv1 "urshanabi/services/gateway/internal/gen/urshanabi/build/v1"
	commonv1 "urshanabi/services/gateway/internal/gen/urshanabi/common/v1"
	"urshanabi/services/gateway/internal/identity"
	"urshanabi/services/gateway/internal/logging"
	"urshanabi/services/gateway/internal/requestid"
)

// fakeQuery stands in for the query service over a real network
// connection, recording the request id it receives.
type fakeQuery struct {
	buildv1.UnimplementedBuildServiceServer
	mu     sync.Mutex
	gotID  string
	fail   bool
	detail error
}

func (f *fakeQuery) id() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.gotID
}

func (f *fakeQuery) GetBuildInfo(ctx context.Context, _ *buildv1.GetBuildInfoRequest) (*buildv1.GetBuildInfoResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if ids := md.Get(requestid.MetadataKey); len(ids) == 1 {
		f.mu.Lock()
		f.gotID = ids[0]
		f.mu.Unlock()
	}
	if f.detail != nil {
		return nil, f.detail
	}
	if f.fail {
		return nil, errors.New("disk /var/lib/secret-path is full")
	}
	return &buildv1.GetBuildInfoResponse{Builds: []*buildv1.BuildInfo{{Component: "query"}}}, nil
}

// start runs the fake query service and a gateway calling it.
func start(t *testing.T, fake *fakeQuery) *httptest.Server {
	return startLogging(t, fake, io.Discard)
}

// startLogging is start, with the gateway's log written to w.
func startLogging(t *testing.T, fake *fakeQuery, w io.Writer) *httptest.Server {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	upstream := rpc.NewServer()
	buildv1.RegisterBuildServiceServer(upstream, fake)
	go func() { _ = upstream.Serve(lis) }() // name-ok
	t.Cleanup(upstream.Stop)

	conn, err := rpc.NewClient(lis.Addr().String(), rpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	self := identity.Build{Component: "gateway", Version: "development", Revision: strings.Repeat("a", 40), CommittedAt: time.Unix(1, 0).UTC()}
	web := httptest.NewServer(New(self, buildv1.NewBuildServiceClient(conn), logging.New(w)).Handler())
	t.Cleanup(web.Close)
	return web
}

func get(t *testing.T, url, id string) (*http.Response, []byte) {
	t.Helper()
	req, _ := http.NewRequest("GET", url, nil)
	if id != "" {
		req.Header.Set(requestid.Header, id)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp, body
}

func TestTheGatewayReportsItselfThenTheQueryService(t *testing.T) {
	fake := &fakeQuery{}
	web := start(t, fake)
	resp, body := get(t, web.URL+"/v1/builds", "test-1")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d: %s", resp.StatusCode, body)
	}
	var got struct{ Builds []struct{ Component string } }
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Builds) != 2 || got.Builds[0].Component != "gateway" || got.Builds[1].Component != "query" {
		t.Fatalf("builds: %s", body)
	}
}

func TestTheRequestIdReachesTheQueryService(t *testing.T) {
	fake := &fakeQuery{}
	web := start(t, fake)
	resp, _ := get(t, web.URL+"/v1/builds", "trace-me")
	if fake.id() != "trace-me" || resp.Header.Get(requestid.Header) != "trace-me" {
		t.Fatalf("query saw %q, response header %q", fake.id(), resp.Header.Get(requestid.Header))
	}
	resp, _ = get(t, web.URL+"/v1/builds", "")
	if fake.id() == "" || fake.id() != resp.Header.Get(requestid.Header) {
		t.Fatalf("an assigned id was not passed on: query %q, header %q", fake.id(), resp.Header.Get(requestid.Header))
	}
}

func TestAFailureDownstreamRevealsNothingAboutIt(t *testing.T) {
	web := start(t, &fakeQuery{fail: true})
	resp, body := get(t, web.URL+"/v1/builds", "test-2")
	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("status %d", resp.StatusCode)
	}
	if strings.Contains(string(body), "secret-path") || !strings.Contains(string(body), "test-2") {
		t.Fatalf("the error body leaks detail or lacks the request id: %s", body)
	}
}

func TestTheRoutersOwnRefusalsCarryAnId(t *testing.T) {
	web := start(t, &fakeQuery{})
	for _, c := range []struct{ method, path string }{{"GET", "/nope"}, {"POST", "/v1/builds"}} {
		req, _ := http.NewRequest(c.method, web.URL+c.path, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		_ = resp.Body.Close()
		if !requestid.Valid(resp.Header.Get(requestid.Header)) {
			t.Errorf("%s %s answered %d with no request id", c.method, c.path, resp.StatusCode)
		}
	}
}

func TestTheGatewayNeverLogsADownstreamUnsafeArgument(t *testing.T) {
	st, err := status.New(codes.Internal, "failed on SECRET-VALUE").WithDetails(&commonv1.ErrorDetail{
		Reason:     "INDEX_UNAVAILABLE",
		UnsafeArgs: []*commonv1.ErrorArgument{{Name: "query_text", Value: "SECRET-VALUE"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	var log syncBuffer
	web := startLogging(t, &fakeQuery{detail: st.Err()}, &log)
	resp, _ := get(t, web.URL+"/v1/builds", "test-3")
	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("status %d", resp.StatusCode)
	}
	if line := log.String(); strings.Contains(line, "SECRET-VALUE") || !strings.Contains(line, "INDEX_UNAVAILABLE") {
		t.Fatalf("the gateway's log leaks the unsafe argument or lacks the reason: %s", line)
	}
}

// syncBuffer is a log destination safe to read while servers write to it.
type syncBuffer struct {
	mu  sync.Mutex
	out strings.Builder
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.out.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.out.String()
}
