package contract

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"

	app "github.com/decentrio/sorobook-api/app"
	types "github.com/decentrio/sorobook-api/types/contract"
	"github.com/decentrio/xdr-converter/converter"
)

func (k Keeper) ContractEntry(ctx context.Context, request *types.ContractEntryRequest) (*types.ContractEntryResponse, error) {
	var entry types.ContractEntryInfo

	err := k.dbHandler.Table(app.CONTRACT_TABLE).
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
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = app.PAGE_SIZE
	}

	offset := (page - 1) * pageSize

	err := k.dbHandler.Table(app.CONTRACT_TABLE).
		Where("contract_id = ?", request.ContractId).
		Where("is_newest = ?", true).
		Limit(pageSize).
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
	}, nil
}

func (k Keeper) ContractKeys(ctx context.Context, request *types.ContractKeysRequest) (*types.ContractKeysResponse, error) {
	var entries []*types.ContractEntry
	var keys []*structpb.Struct

	err := k.dbHandler.Table(app.CONTRACT_TABLE).
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

	err := k.dbHandler.Table(app.CONTRACT_TABLE).
		Select("contracts.contract_id").
		Distinct().
		Joins("JOIN transactions ON transactions.hash = contracts.tx_hash").
		Where("source_address = ?", request.Address).
		Find(&contracts).Error
	fmt.Print(contracts)
	if err != nil {
		return &types.UserInteractionContractsResponse{}, err
	}

	return &types.UserInteractionContractsResponse{
		Contracts: contracts,
	}, nil
}

func (k Keeper) ContractCode(ctx context.Context, request *types.ContractCodeRequest) (*types.ContractCodeResponse, error) {
	var data types.ContractCode

	err := k.dbHandler.Table(app.CONTRACT_CODES).
		Where("contract_id = ?", request.ContractId).
		First(&data).Error

	if err != nil {
		return &types.ContractCodeResponse{
			Found: false,
		}, err
	}

	return &types.ContractCodeResponse{
		Contract: &data,
		Found: true,
	}, nil
}

func (k Keeper) ContractCodes(ctx context.Context, request *types.ContractCodesRequest) (*types.ContractCodesResponse, error) {
	var data []*types.ContractCode
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = app.PAGE_SIZE
	}

	offset := (page - 1) * pageSize

	err := k.dbHandler.Table(app.CONTRACT_CODES).
		Limit(pageSize).
		Order("created_ledger DESC").
		Offset(offset).
		Find(&data).Error

	if err != nil {
		return &types.ContractCodesResponse{}, err
	}


	return &types.ContractCodesResponse{
		Data: data,
	}, nil
}

func (k Keeper) ContractsAtLedger(ctx context.Context, request *types.ContractsAtLedgerRequest) (*types.ContractsAtLedgerResponse, error) {
	var data []*types.ContractCode
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = app.PAGE_SIZE
	}

	offset := (page - 1) * pageSize

	err := k.dbHandler.Table(app.CONTRACT_CODES).
		Limit(pageSize).
		Where("created_ledger = ?", request.Ledger).
		Offset(offset).
		Find(&data).Error

	if err != nil {
		return &types.ContractsAtLedgerResponse{}, err
	}


	return &types.ContractsAtLedgerResponse{
		Data: data,
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
