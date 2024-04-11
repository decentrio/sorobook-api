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

	err := k.dbHandler.Table(CONTRACT_TABLE).
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
	var entries []*types.ContractEntry
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	offset := (page - 1) * PAGE_SIZE * 2

	err := k.dbHandler.Table(CONTRACT_TABLE).
		Where("contract_id = ?", request.ContractId).
		Where("is_newest = ?", true).
		Limit(PAGE_SIZE * 2).
		Offset(offset).
		Find(&entries).Error

	if err != nil {
		return &types.ContractDataResponse{}, err
	}

	var infos []*types.ContractEntryInfo

	for _, entry := range entries {
		info, err := convertToEntryInfo(entry)
		if err != nil {
			return &types.ContractDataResponse{}, err
		}

		infos = append(infos, info)
	}

	return &types.ContractDataResponse{
		Data: infos,
		Page: int32(page),
	}, nil
}

func (k Keeper) ContractKeys(ctx context.Context, request *types.ContractKeysRequest) (*types.ContractKeysResponse, error) {
	var entries []*types.ContractEntry
	var keys []*structpb.Struct

	err := k.dbHandler.Table(CONTRACT_TABLE).
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

func (k Keeper) UserInteractionContracts(ctx context.Context, request *types.UserInteractionContractsRequest) (*types.UserInteractionContractsResponse, error) {
	var contracts []string

	err := k.dbHandler.Table(CONTRACT_TABLE).
		Select("contracts.contract_id").
		Joins("JOIN transactions ON transactions.hash = contracts.tx_hash").
		Where("source_address = ?", request.Address).
		Find(&contracts).Error

	if err != nil {
		return &types.UserInteractionContractsResponse{}, err
	}

	return &types.UserInteractionContractsResponse{
		Contracts: contracts,
	}, nil
}

func convertToEntryInfo(entry *types.ContractEntry) (*types.ContractEntryInfo, error) {
	keyJson := &structpb.Struct{}
	keyData, err := converter.MarshalJSONContractKeyXdr(entry.KeyXdr)
	if err != nil {
		return &types.ContractEntryInfo{}, err
	}
	if err := json.Unmarshal(keyData, keyJson); err != nil {
		return &types.ContractEntryInfo{}, err
	}

	valueJson := &structpb.Struct{}
	valueData, err := converter.MarshalJSONContractValueXdr(entry.ValueXdr)
	if err != nil {
		return &types.ContractEntryInfo{}, err
	}
	if err := json.Unmarshal(valueData, valueJson); err != nil {
		return &types.ContractEntryInfo{}, err
	}

	return &types.ContractEntryInfo{
		Key:   keyJson,
		Value: valueJson,
	}, nil
}
