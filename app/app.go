package app

import (
	"github.com/decentrio/sorobook-api/database"
	types "github.com/decentrio/sorobook-api/types/v1"

)

type Keeper struct {
	dbHandler database.DBHandler
	types.UnimplementedQueryServer
}


func NewKeeper() *Keeper {
	dbHandler := database.NewDBHandler()

	return &Keeper{
		dbHandler: *dbHandler,
	}
}
