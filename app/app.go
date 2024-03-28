package app

import (
	"gorm.io/gorm"

	"github.com/decentrio/sorobook-api/database"
	types "github.com/decentrio/sorobook-api/types/v1"
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
