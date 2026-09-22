// Command gateway is Urshanabi's front door: the web API that people,
// agents and applications call.
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	rpc "google.golang.org/grpc"                  // name-ok
	"google.golang.org/grpc/credentials/insecure" // name-ok

	buildv1 "urshanabi/services/gateway/internal/gen/urshanabi/build/v1"
	"urshanabi/services/gateway/internal/identity"
	"urshanabi/services/gateway/internal/logging"
	"urshanabi/services/gateway/internal/server"
)

func main() {
	log := logging.New(os.Stderr)
	self, err := identity.Current()
	if err != nil {
		log.Error("cannot read this build's identity", "error", err)
		os.Exit(1)
	}
	// Mutual TLS between services is the mesh's job (roadmap R-126), so the
	// connection itself is plain.
	conn, err := rpc.NewClient(env("URSHANABI_QUERY_ADDRESS", "127.0.0.1:50051"),
		rpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("cannot reach the query service", "error", err)
		os.Exit(1)
	}
	defer func() { _ = conn.Close() }()

	web := &http.Server{
		Addr:              env("URSHANABI_LISTEN", "127.0.0.1:8080"),
		Handler:           server.New(self, buildv1.NewBuildServiceClient(conn), log).Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() { // name-ok
		<-ctx.Done()
		log.Info("shutting down")
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = web.Shutdown(shutdown)
	}()
	log.Info("listening", "address", web.Addr)
	if err := web.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("the gateway stopped", "error", err)
		os.Exit(1)
	}
}

func env(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}
