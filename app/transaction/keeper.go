package transaction

import (
	"gorm.io/gorm"

	types "github.com/decentrio/sorobook-api/types/transaction"
)

type Keeper struct {
	dbHandler *gorm.DB
	types.UnimplementedTransactionQueryServer
}

func NewKeeper(db *gorm.DB) *Keeper {
	return &Keeper{
		dbHandler: db,
	}
}

var _ types.TransactionQueryServer = Keeper{}
