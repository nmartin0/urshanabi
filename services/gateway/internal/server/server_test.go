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
	"testing"
	"time"

	rpc "google.golang.org/grpc"                  // name-ok
	"google.golang.org/grpc/credentials/insecure" // name-ok
	"google.golang.org/grpc/metadata"             // name-ok

	buildv1 "urshanabi/services/gateway/internal/gen/urshanabi/build/v1"
	"urshanabi/services/gateway/internal/identity"
	"urshanabi/services/gateway/internal/logging"
	"urshanabi/services/gateway/internal/requestid"
)

// fakeQuery stands in for the query service over a real network
// connection, recording the request id it receives.
type fakeQuery struct {
	buildv1.UnimplementedBuildServiceServer
	gotID string
	fail  bool
}

func (f *fakeQuery) GetBuildInfo(ctx context.Context, _ *buildv1.GetBuildInfoRequest) (*buildv1.GetBuildInfoResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if ids := md.Get(requestid.MetadataKey); len(ids) == 1 {
		f.gotID = ids[0]
	}
	if f.fail {
		return nil, errors.New("disk /var/lib/secret-path is full")
	}
	return &buildv1.GetBuildInfoResponse{Builds: []*buildv1.BuildInfo{{Component: "query"}}}, nil
}

// start runs the fake query service and a gateway calling it.
func start(t *testing.T, fake *fakeQuery) *httptest.Server {
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
	web := httptest.NewServer(New(self, buildv1.NewBuildServiceClient(conn), logging.New(io.Discard)).Handler())
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
	if fake.gotID != "trace-me" || resp.Header.Get(requestid.Header) != "trace-me" {
		t.Fatalf("query saw %q, response header %q", fake.gotID, resp.Header.Get(requestid.Header))
	}
	resp, _ = get(t, web.URL+"/v1/builds", "")
	if fake.gotID == "" || fake.gotID != resp.Header.Get(requestid.Header) {
		t.Fatalf("an assigned id was not passed on: query %q, header %q", fake.gotID, resp.Header.Get(requestid.Header))
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
