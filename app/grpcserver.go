package app

import (
	"context"

	types "github.com/decentrio/sorobook-api/types/v1"
)

var _ types.QueryServer = Keeper{}

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

func (q Keeper) Transaction(ctx context.Context, request *types.TransactionRequest) (*types.TransactionResponse, error) {
	return nil, nil
}

func (q Keeper) Event(ctx context.Context, request *types.EventRequest) (*types.EventResponse, error) {
	return nil, nil
}

func (q Keeper) ContractEntry(ctx context.Context, request *types.ContractEntryRequest) (*types.ContractEntryResponse, error) {
	return nil, nil
}
