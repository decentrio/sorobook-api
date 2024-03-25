package app

import (
	"github.com/decentrio/sorobook-api/database"
)

type Keeper struct {
	dbHandler database.DBHandler
}


func NewKeeper() *Keeper {
	dbHandler := database.NewDBHandler()

	return &Keeper{
		dbHandler: *dbHandler,
	}
}
