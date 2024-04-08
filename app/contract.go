package app

import (
	"context"
	"encoding/json"

	"google.golang.org/protobuf/types/known/structpb"

	types "github.com/decentrio/sorobook-api/types/v1"
	"github.com/decentrio/xdr-converter/converter"
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

func (k Keeper) ContractKeys(ctx context.Context, request *types.ContractKeysRequest) (*types.ContractKeysResponse, error) {
	var entries []*types.ContractEntry
	var keys []*structpb.Struct

	err := k.dbHandler.Table("contracts").
		Where("contract_id = ?", request.ContractId).
		Limit(20).
		Find(&entries).Error
	if err != nil {
		return &types.ContractKeysResponse{
			Keys: []*structpb.Struct{},
		}, err
	}

	for _, entry := range entries {
		keyJson := &structpb.Struct{}
		keyData, err := converter.MarshalJSONContractKeyXdr(entry.KeyXdr)
		if err != nil {
			return &types.ContractKeysResponse{
				Keys: []*structpb.Struct{},
			}, err
		}
		if err := json.Unmarshal(keyData, keyJson); err != nil {
			return &types.ContractKeysResponse{
				Keys: []*structpb.Struct{},
			}, err
		}

		keys = append(keys, keyJson)
	}

	return &types.ContractKeysResponse{
		Keys: keys,
	}, nil
}
