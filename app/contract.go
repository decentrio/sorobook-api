package app

import (
	"context"

	types "github.com/decentrio/sorobook-api/types/v1"
)

func (k Keeper) ContractEntry(ctx context.Context, request *types.ContractEntryRequest) (*types.ContractEntryResponse, error) {
	var entry types.ContractEntryInfo

	err := k.dbHandler.Table("contracts").
		Where("contract_id = ?", request.ContractId).
		Where("key_xdr = ?", request.KeyXdr).
		First(&entry).Error

	if err != nil {
		return &types.ContractEntryResponse{
			Found: false,
			Entry: &types.ContractEntryInfo{},
		}, err
	}

	return &types.ContractEntryResponse{
		Entry: &entry,
		Found: true,
	}, nil
}

func (k Keeper) ContractData(ctx context.Context, request *types.ContractDataRequest) (*types.ContractDataResponse, error) {
	return nil, nil
}
