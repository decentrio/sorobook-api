package app

import (
	grpc "google.golang.org/grpc"
)

type AppModule interface {
	RegisterServices(server *grpc.Server)
}
