package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/decentrio/sorobook-api/app"
	types "github.com/decentrio/sorobook-api/types/v1"
)

func runGRPCServer() error {
	keeper := app.NewKeeper()
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	types.RegisterQueryServer(s, keeper)
	return s.Serve(lis)
}

func runHTTPServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := types.RegisterQueryHandlerFromEndpoint(ctx, mux, ":9090", opts)

	if err != nil {
		return err
	}

	// Serve Swagger UI
	http.Handle("/", mux)

	log.Println("HTTP server listening on :8080")
	return http.ListenAndServe(":8080", nil)
}

func main() {
	go func() {
		if err := runGRPCServer(); err != nil {
			log.Fatalf("failed to run gRPC server: %v", err)
		}
	}()
	if err := runHTTPServer(); err != nil {
		log.Fatalf("failed to run HTTP server: %v", err)
	}
}
