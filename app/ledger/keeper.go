package ledger

import (
	"gorm.io/gorm"

	types "github.com/decentrio/sorobook-api/types/ledger"
)

type Keeper struct {
	dbHandler *gorm.DB
	types.UnimplementedLedgerQueryServer
}

func NewKeeper(db *gorm.DB) *Keeper {
	return &Keeper{
		dbHandler: db,
	}
}

var _ types.LedgerQueryServer = Keeper{}
