package app

import (
	"gorm.io/gorm"

	"github.com/decentrio/sorobook-api/database"
	types "github.com/decentrio/sorobook-api/types/v1"
)

const (
	PAGE_SIZE = 10
	LEDGER_TABLE = "ledgers"
	TRANSACTION_TABLE = "transactions"
	CONTRACT_TABLE = "contracts"
	EVENT_TABLE = "wasm_contract_events"
)

type Keeper struct {
	dbHandler *gorm.DB
	types.UnimplementedQueryServer
}

func NewKeeper() *Keeper {
	dbHandler := database.NewDBHandler()

	return &Keeper{
		dbHandler: dbHandler,
	}
}

var _ types.QueryServer = Keeper{}