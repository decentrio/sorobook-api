package contract

import (
	"gorm.io/gorm"

	types "github.com/decentrio/sorobook-api/types/contract"
)

type Keeper struct {
	dbHandler *gorm.DB
	types.UnimplementedContractQueryServer
}

func NewKeeper(db *gorm.DB) *Keeper {
	return &Keeper{
		dbHandler: db,
	}
}

var _ types.ContractQueryServer = Keeper{}
