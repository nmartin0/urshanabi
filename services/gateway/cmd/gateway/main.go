// Command gateway is Urshanabi's front door: the web API that people,
// agents and applications call.
package main

import (
	"context"
	"net"
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
		Handler:           server.New(self, buildv1.NewBuildServiceClient(conn), log).Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}
	address := env("URSHANABI_LISTEN", "127.0.0.1:8080")
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Error("cannot listen", "address", address, "error", err)
		os.Exit(1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	log.Info("listening", "address", lis.Addr().String())
	// Serve returns only after in-flight requests finish, so the connection
	// they use is closed after them, not under them.
	if err := server.Serve(ctx, web, lis, 10*time.Second); err != nil {
		log.Error("the gateway stopped", "error", err)
		os.Exit(1)
	}
	log.Info("stopped")
}

func env(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}
