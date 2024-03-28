package app

import (
	"context"

	types "github.com/decentrio/sorobook-api/types/v1"
)

func (q Keeper) Ledger(ctx context.Context, request *types.LedgerRequest) (*types.LedgerResponse, error) {
	var ledger types.Ledger
	
	err := q.dbHandler.Table("ledgers").Where("sequence = ?", request.Sequence).First(&ledger).Error
	if err != nil {
		return &types.LedgerResponse{
			Found: false,
			Ledger: &types.Ledger{},
		}, err
	}

	return &types.LedgerResponse{
		Ledger: &ledger,
		Found: true,
	}, nil
}
