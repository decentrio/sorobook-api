package contract

import (
	"gorm.io/gorm"

	"github.com/decentrio/sorobook-api/database"
	types "github.com/decentrio/sorobook-api/types/contract"
)

const (
	PAGE_SIZE         = 20
	LEDGER_TABLE      = "ledgers"
	TRANSACTION_TABLE = "transactions"
	CONTRACT_TABLE    = "contracts"
	EVENT_TABLE       = "wasm_contract_events"
	TRANSFER_TABLE    = "asset_contract_transfer_events"
	MINT_TABLE        = "asset_contract_mint_events"
	BURN_TABLE        = "asset_contract_burn_events"
	CLAWBACK_TABLE    = "asset_contract_clawback_events"
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
