package transaction

import (
	types "github.com/decentrio/sorobook-api/types/transaction"
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
	types.RegisterTransactionQueryServer(server, am.keeper)
}
