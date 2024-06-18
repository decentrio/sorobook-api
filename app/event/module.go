package event

import (
	types "github.com/decentrio/sorobook-api/types/event"
	"google.golang.org/grpc"
)

type AppModule struct {
	keeper Keeper
}

func NewAppModule(
	keeper Keeper,
) AppModule {
	return AppModule{
		keeper: keeper,
	}
}

func (am AppModule) RegisterServices(server *grpc.Server) {
	types.RegisterEventQueryServer(server, am.keeper)
}
