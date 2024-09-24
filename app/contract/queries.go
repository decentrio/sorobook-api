package contract

import (
	"context"
	"encoding/hex"
	"encoding/json"

	"google.golang.org/protobuf/types/known/structpb"

	app "github.com/decentrio/sorobook-api/app"
	types "github.com/decentrio/sorobook-api/types/contract"
	"github.com/decentrio/xdr-converter/converter"
)

func (k Keeper) ContractEntry(_ context.Context, request *types.ContractEntryRequest) (*types.ContractEntryResponse, error) {
	var entry types.ContractEntryInfo

	query := k.dbHandler.Table(app.CONTRACT_TABLE).
		Where("contract_id = ?", request.ContractId).
		Where("key_xdr = ?", request.KeyXdr)

	if request.Ledger != 0 {
		query = query.Where("ledger >= ?", request.Ledger).Where("ledger >= ?", request.Ledger)
	} else {
		query = query.Where("is_newest = true")
	}

	err := query.First(&entry).Error
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

func (k Keeper) ContractData(_ context.Context, request *types.ContractDataRequest) (*types.ContractDataResponse, error) {
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

	query := k.dbHandler.Table(app.CONTRACT_TABLE).
		Where("contract_id = ?", request.ContractId).
		Where("is_newest = ?", true)

	if request.Ledger != 0 {
		query = query.Where("ledger >= ?", request.Ledger).Where("ledger >= ?", request.Ledger)
	} else {
		query = query.Where("is_newest = true")
	}

	if request.KeyXdr != "" {
		keyByte, err := hex.DecodeString(request.KeyXdr)
		if err != nil {
			return &types.ContractDataResponse{}, err
		}
		query = query.Where("key_xdr = ?", keyByte)
	}

	err := query.Limit(pageSize).
		Offset(offset).Find(&entries).Error
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

func (k Keeper) ContractKeys(_ context.Context, request *types.ContractKeysRequest) (*types.ContractKeysResponse, error) {
	var entries []*types.ContractEntry
	var keys []*structpb.Struct

	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = app.PAGE_SIZE
	}

	offset := (page - 1) * pageSize

	query := k.dbHandler.Table(app.CONTRACT_TABLE).
		Where("contract_id = ?", request.ContractId).
		Where("is_newest = ?", true)

	if request.Ledger != 0 {
		query = query.Where("ledger >= ?", request.Ledger).Where("ledger >= ?", request.Ledger)
	} else {
		query = query.Where("is_newest = true")
	}

	err := query.Limit(pageSize).
		Offset(offset).Find(&entries).Error
	if err != nil {
		return &types.ContractKeysResponse{
			Keys: []*structpb.Struct{},
		}, err
	}

	for _, entry := range entries {
		keyJSON := &structpb.Struct{}
		keyData, err := converter.MarshalJSONContractKeyXdr(entry.KeyXdr)
		if err != nil {
			return &types.ContractKeysResponse{
				Keys: []*structpb.Struct{},
			}, err
		}
		if err := json.Unmarshal(keyData, keyJSON); err != nil {
			return &types.ContractKeysResponse{
				Keys: []*structpb.Struct{},
			}, err
		}

		keys = append(keys, keyJSON)
	}

	return &types.ContractKeysResponse{
		Keys: keys,
	}, nil
}

func (k Keeper) UserInteractionContracts(_ context.Context, request *types.UserInteractionContractsRequest) (*types.UserInteractionContractsResponse, error) {
	var contracts []string

	err := k.dbHandler.Table(app.CONTRACT_TABLE).
		Select("contracts.contract_id").
		Distinct().
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

func (k Keeper) ContractCode(_ context.Context, request *types.ContractCodeRequest) (*types.ContractCodeResponse, error) {
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
		Found:    true,
	}, nil
}

func (k Keeper) ContractCodes(_ context.Context, request *types.ContractCodesRequest) (*types.ContractCodesResponse, error) {
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
		Order("created_ledger DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&data).Error
	if err != nil {
		return &types.ContractCodesResponse{}, err
	}

	return &types.ContractCodesResponse{
		Data: data,
	}, nil
}

func (k Keeper) ContractsAtLedger(_ context.Context, request *types.ContractsAtLedgerRequest) (*types.ContractsAtLedgerResponse, error) {
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

func (k Keeper) ContractInvoke(_ context.Context, request *types.ContractInvokeRequest) (*types.ContractInvokeResponse, error) {
	var data types.ContractInvoke

	err := k.dbHandler.Table(app.INVOKE_TXS).
		Where("hash = ?", request.Hash).
		First(&data).Error
	if err != nil {
		return &types.ContractInvokeResponse{
			Found: false,
		}, err
	}

	info, err := convertToInvokeInfo(&data)
	if err != nil {
		return &types.ContractInvokeResponse{
			Found: false,
		}, err
	}

	return &types.ContractInvokeResponse{
		Info:  info,
		Found: true,
	}, nil
}

func (k Keeper) ContractInvokes(_ context.Context, request *types.ContractInvokesRequest) (*types.ContractInvokesResponse, error) {
	var data []*types.ContractInvoke
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = app.PAGE_SIZE
	}

	offset := (page - 1) * pageSize

	query := k.dbHandler.Table(app.INVOKE_TXS).
		Joins("JOIN transactions ON transactions.hash = invoke_transactions.hash").
		Where("contract_id = ?", request.ContractId).
		Order("transactions.ledger DESC")

	if request.FunctionName != "" {
		query = query.Where("function_name = ?", request.FunctionName)
	}

	err := query.Limit(pageSize).
		Offset(offset).
		Find(&data).Error
	if err != nil {
		return &types.ContractInvokesResponse{}, err
	}

	var infos []*types.ContractInvokeInfo

	for _, item := range data {
		invokeInfo, err := convertToInvokeInfo(item)
		if err != nil {
			return &types.ContractInvokesResponse{}, err
		}
		infos = append(infos, invokeInfo)
	}

	return &types.ContractInvokesResponse{
		Data: infos,
	}, nil
}

func (k Keeper) ContractInvokesAtLedger(_ context.Context, request *types.ContractInvokesAtLedgerRequest) (*types.ContractInvokesAtLedgerResponse, error) {
	var data []*types.ContractInvoke
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = app.PAGE_SIZE
	}

	offset := (page - 1) * pageSize

	err := k.dbHandler.Table(app.INVOKE_TXS).
		Joins("JOIN transactions ON transactions.hash = invoke_transactions.hash").
		Where("contract_id = ?", request.ContractId).
		Where("transactions.ledger = ?", request.Ledger).
		Limit(pageSize).
		Offset(offset).
		Find(&data).Error
	if err != nil {
		return &types.ContractInvokesAtLedgerResponse{}, err
	}

	var infos []*types.ContractInvokeInfo

	for _, item := range data {
		invokeInfo, err := convertToInvokeInfo(item)
		if err != nil {
			return &types.ContractInvokesAtLedgerResponse{}, err
		}
		infos = append(infos, invokeInfo)
	}

	return &types.ContractInvokesAtLedgerResponse{
		Data: infos,
	}, nil
}

func (k Keeper) ContractInvokesByUser(_ context.Context, request *types.ContractInvokesByUserRequest) (*types.ContractInvokesByUserResponse, error) {
	var data []*types.ContractInvoke
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = app.PAGE_SIZE
	}

	offset := (page - 1) * pageSize

	query := k.dbHandler.Table(app.INVOKE_TXS).
		Joins("JOIN transactions ON transactions.hash = invoke_transactions.hash").
		Where("transactions.source_address = ?", request.Address)

	if request.ContractId != "" {
		query = query.Where("contract_id = ?", request.ContractId)
	}

	if request.Ledger != 0 {
		query = query.Where("transactions.ledger = ?", request.Ledger)
	}

	err := query.Limit(pageSize).
		Offset(offset).
		Find(&data).Error
	if err != nil {
		return &types.ContractInvokesByUserResponse{}, err
	}

	var infos []*types.ContractInvokeInfo

	for _, item := range data {
		invokeInfo, err := convertToInvokeInfo(item)
		if err != nil {
			return &types.ContractInvokesByUserResponse{}, err
		}
		infos = append(infos, invokeInfo)
	}

	return &types.ContractInvokesByUserResponse{
		Data: infos,
	}, nil
}

func convertToEntryInfo(entry *types.ContractEntry) (*types.ContractEntryInfo, error) {
	keyJSON := &structpb.Struct{}
	keyData, err := converter.MarshalJSONContractKeyXdr(entry.KeyXdr)
	if err != nil {
		return &types.ContractEntryInfo{}, err
	}
	if err := json.Unmarshal(keyData, keyJSON); err != nil {
		return &types.ContractEntryInfo{}, err
	}

	valueJSON := &structpb.Struct{}
	valueData, err := converter.MarshalJSONContractValueXdr(entry.ValueXdr)
	if err != nil {
		return &types.ContractEntryInfo{}, err
	}
	if err := json.Unmarshal(valueData, valueJSON); err != nil {
		return &types.ContractEntryInfo{}, err
	}

	return &types.ContractEntryInfo{
		Key:   keyJSON,
		Value: valueJSON,
	}, nil
}

func convertToInvokeInfo(data *types.ContractInvoke) (*types.ContractInvokeInfo, error) {
	argsJSON := &structpb.Struct{}

	argsData, err := converter.MarshalJSONInvokeContractArgsXdr(data.Args)
	if err != nil {
		return &types.ContractInvokeInfo{}, err
	}
	if err := json.Unmarshal(argsData, argsJSON); err != nil {
		return &types.ContractInvokeInfo{}, err
	}

	return &types.ContractInvokeInfo{
		Hash:         data.Hash,
		ContractId:   data.ContractId,
		FunctionName: data.FunctionName,
		FunctionType: data.FunctionType,
		Args:         argsJSON,
	}, nil
}
