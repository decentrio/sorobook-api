package app

import (
	grpc "google.golang.org/grpc"
)

type App struct {
	Modules   []AppModule
	Server    *grpc.Server
}

func NewApp(server *grpc.Server, modules []AppModule) *App {
	return &App{
		Modules:   modules,
		Server:    server,
	}
}

func (app *App) RegisterServices() {
	for _, module := range app.Modules {
		module.RegisterServices(app.Server)
	}
}
