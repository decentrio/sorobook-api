package app

import (
	"context"

	types "github.com/decentrio/sorobook-api/types/v1"
)

func (k Keeper) Ledger(ctx context.Context, request *types.LedgerRequest) (*types.LedgerResponse, error) {
	var ledger types.Ledger

	err := k.dbHandler.Table("ledgers").Where("sequence = ?", request.Sequence).First(&ledger).Error
	if err != nil {
		return &types.LedgerResponse{
			Found:  false,
			Ledger: &types.Ledger{},
		}, err
	}

	return &types.LedgerResponse{
		Ledger: &ledger,
		Found:  true,
	}, nil
}

func (k Keeper) LedgerHash(ctx context.Context, request *types.LedgerHashRequest) (*types.LedgerHashResponse, error) {
	var ledger types.Ledger

	err := k.dbHandler.Table("ledgers").Where("hash = ?", request.Hash).First(&ledger).Error
	if err != nil {
		return &types.LedgerHashResponse{
			Found:  false,
			Ledger: &types.Ledger{},
		}, err
	}

	return &types.LedgerHashResponse{
		Ledger: &ledger,
		Found:  true,
	}, nil
}
