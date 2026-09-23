// Package server is the gateway's web API. In the walking skeleton it
// answers one question: which builds are running (roadmap R-121).
package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"google.golang.org/grpc/metadata"                    // name-ok
	wire "google.golang.org/protobuf/encoding/protojson" // name-ok
	"google.golang.org/protobuf/types/known/timestamppb" // name-ok

	"urshanabi/services/gateway/internal/failure"
	buildv1 "urshanabi/services/gateway/internal/gen/urshanabi/build/v1"
	"urshanabi/services/gateway/internal/identity"
	"urshanabi/services/gateway/internal/requestid"
)

// callTimeout bounds each call to another service.
const callTimeout = 5 * time.Second

// Server answers the web API.
type Server struct {
	self  *buildv1.BuildInfo
	query buildv1.BuildServiceClient
	log   *slog.Logger
}

// New returns a server reporting self and asking query for its build.
func New(self identity.Build, query buildv1.BuildServiceClient, log *slog.Logger) *Server {
	return &Server{
		self: &buildv1.BuildInfo{
			Component:         self.Component,
			Version:           self.Version,
			SourceRevision:    self.Revision,
			SourceCommittedAt: timestamppb.New(self.CommittedAt),
			ArtifactSha256:    self.ArtifactSHA256,
		},
		query: query,
		log:   log,
	}
}

// Handler routes the web API.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/builds", s.builds)
	// The router's own refusals answer in the same vocabulary as every
	// other failure (roadmap R-22), not in the framework's words.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		kind := failure.Empty
		if r.URL.Path == "/v1/builds" {
			kind = failure.Refused // the route exists; this method does not
		}
		status := http.StatusNotFound
		if kind == failure.Refused {
			status = http.StatusMethodNotAllowed
		}
		writeFailure(w, status, kind, requestid.FromContext(r.Context()))
	})
	return requestid.Middleware(securityHeaders(mux))
}

// builds answers with the gateway's build first, then the query service's.
func (s *Server) builds(w http.ResponseWriter, r *http.Request) {
	id := requestid.FromContext(r.Context())
	ctx, cancel := context.WithTimeout(metadata.AppendToOutgoingContext(r.Context(), requestid.MetadataKey, id), callTimeout)
	defer cancel()

	answer, err := s.query.GetBuildInfo(ctx, &buildv1.GetBuildInfoRequest{})
	if err != nil {
		// The caller learns only that a service is unavailable, never why:
		// error text can reveal internals (SP 800-53 SI-11).
		s.log.Error("GetBuildInfo failed", append([]any{"request_id", id}, failure.Attrs(err)...)...)
		writeFailure(w, http.StatusBadGateway, failure.Kind(err), id)
		return
	}
	builds := append([]*buildv1.BuildInfo{s.self}, answer.GetBuilds()...)
	body, err := wire.Marshal(&buildv1.GetBuildInfoResponse{Builds: builds})
	if err != nil {
		s.log.Error("encoding the answer failed", "request_id", id, "error", err)
		writeFailure(w, http.StatusInternalServerError, failure.Misconfigured, id)
		return
	}
	s.log.Info("GetBuildInfo", "request_id", id, "builds", len(builds))
	writeJSON(w, http.StatusOK, body)
}

func writeJSON(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

// writeFailure answers with the failure's kind and the request's id, and
// nothing else: a caller learns what it can do, never what went wrong
// inside (roadmap R-22, SP 800-53 SI-11).
func writeFailure(w http.ResponseWriter, status int, kind string, id string) {
	writeJSON(w, status, []byte(`{"kind":"`+kind+`","request_id":"`+id+`"}`))
}
