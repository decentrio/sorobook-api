package app

import (
	grpc "google.golang.org/grpc"
)

type AppModule interface { //nolint
	RegisterServices(server *grpc.Server)
}
