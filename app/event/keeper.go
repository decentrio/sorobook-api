package event

import (
	"gorm.io/gorm"

	types "github.com/decentrio/sorobook-api/types/event"
)

type Keeper struct {
	dbHandler *gorm.DB
	types.UnimplementedEventQueryServer
}

func NewKeeper(db *gorm.DB) *Keeper {
	return &Keeper{
		dbHandler: db,
	}
}

var _ types.EventQueryServer = Keeper{}
