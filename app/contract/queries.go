package contract

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"google.golang.org/protobuf/types/known/structpb"

	app "github.com/decentrio/sorobook-api/app"
	types "github.com/decentrio/sorobook-api/types/contract"
	"github.com/decentrio/xdr-converter/converter"
	"github.com/stellar/go/xdr"
)

func (k Keeper) ContractEntry(ctx context.Context, request *types.ContractEntryRequest) (*types.ContractEntryResponse, error) {
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

func (k Keeper) ContractKeys(ctx context.Context, request *types.ContractKeysRequest) (*types.ContractKeysResponse, error) {
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
		Found:    true,
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

func (k Keeper) ContractInvoke(ctx context.Context, request *types.ContractInvokeRequest) (*types.ContractInvokeResponse, error) {
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

func (k Keeper) ContractInvokes(ctx context.Context, request *types.ContractInvokesRequest) (*types.ContractInvokesResponse, error) {
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

func (k Keeper) ContractInvokesAtLedger(ctx context.Context, request *types.ContractInvokesAtLedgerRequest) (*types.ContractInvokesAtLedgerResponse, error) {
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

func (k Keeper) ContractInvokesByUser(ctx context.Context, request *types.ContractInvokesByUserRequest) (*types.ContractInvokesByUserResponse, error) {
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

func convertToInvokeInfo(data *types.ContractInvoke) (*types.ContractInvokeInfo, error) {
	argsJson := &structpb.Struct{}

	argsData, err := converter.MarshalJSONInvokeContractArgsXdr(data.Args)
	if err != nil {
		return &types.ContractInvokeInfo{}, err
	}
	if err := json.Unmarshal(argsData, argsJson); err != nil {
		return &types.ContractInvokeInfo{}, err
	}

	return &types.ContractInvokeInfo{
		Hash:         data.Hash,
		ContractId:   data.ContractId,
		FunctionName: data.FunctionName,
		FunctionType: data.FunctionType,
		Args:         argsJson,
	}, nil
}

func (k Keeper) ContractKeyXdr(ctx context.Context, request *types.ContractKeyXdrRequest) (*types.ContractKeyXdrResponse, error) {
	key_xdr := ""
	if request.KeyName != "" && request.KeyType != "" {
		xdrType, data, err := convertToData(request.KeyType, request.KeyName)
		if err != nil {
			return &types.ContractKeyXdrResponse{}, err
		}

		xdrKey, err := xdr.NewScVal(xdrType, data)
		if err != nil {
			return &types.ContractKeyXdrResponse{}, err
		}

		bytes, err := xdrKey.MarshalBinary()
		if err != nil {
			return &types.ContractKeyXdrResponse{}, err
		}
		key_xdr = hex.EncodeToString(bytes)
	}

	return &types.ContractKeyXdrResponse{
		KeyXdr: key_xdr,
	}, nil
}

func convertToData(key_type string, key_name string) (xdr.ScValType, interface{}, error) {
	values := strings.Split(key_name, "#")

	switch key_type {
	case "bool":
		if len(values) == 1 {
			return convertToDataBool(values[0])
		}
	case "Uint32":
		if len(values) == 1 {
			return convertToDataUint32(values[0])
		}
	case "Int32":
		if len(values) == 1 {
			return convertToDataInt32(values[0])
		}
	case "Uint64":
		if len(values) == 1 {
			return convertToDataUint64(values[0])
		}
	case "Int64":
		if len(values) == 1 {
			return convertToDataInt64(values[0])
		}
	case "TimePoint":
		if len(values) == 1 {
			return convertToDataTimePoint(values[0])
		}
	case "Duration":
		if len(values) == 1 {
			return convertToDataDuration(values[0])
		}
	case "UInt128Parts":
		if len(values) == 2 {
			return convertToDataUInt128Parts(values[0], values[1])
		}
	case "Int128Parts":
		if len(values) == 2 {
			return convertToDataInt128Parts(values[0], values[1])
		}
	case "UInt256Parts":
		if len(values) == 4 {
			return convertToDataUInt256Parts(values[0], values[1], values[2], values[3])
		}
	case "Int256Parts":
		if len(values) == 4 {
			return convertToDataInt256Parts(values[0], values[1], values[2], values[3])
		}
	case "ScBytes":
		if len(values) == 1 {
			return convertToDataScBytes(values[0])
		}
	case "ScString":
		if len(values) == 1 {
			return convertToDataScString(values[0])
		}
	case "ScSymbol":
		if len(values) == 1 {
			return convertToDataScSymbol(values[0])
		}
	case "ScNonceKey":
		if len(values) == 1 {
			return convertToDataScNonceKey(values[0])
		}
	default:
		return 0, nil, errors.New("convert false")
	}

	return 0, nil, errors.New("convert false")
}

func convertToDataBool(value string) (xdr.ScValType, interface{}, error) {
	data, err := strconv.ParseBool(value)
	if err != nil {
		return 0, nil, err
	}

	return xdr.ScValTypeScvBool, data, err
}

func convertToDataUint32(value string) (xdr.ScValType, interface{}, error) {
	data, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, nil, err
	}

	return xdr.ScValTypeScvU32, xdr.Uint32(data), err
}

func convertToDataInt32(value string) (xdr.ScValType, interface{}, error) {
	data, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, nil, err
	}

	return xdr.ScValTypeScvI32, xdr.Int32(data), err
}

func convertToDataUint64(value string) (xdr.ScValType, interface{}, error) {
	data, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	return xdr.ScValTypeScvU64, xdr.Uint64(data), err
}

func convertToDataInt64(value string) (xdr.ScValType, interface{}, error) {
	data, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	return xdr.ScValTypeScvI64, xdr.Int64(data), err
}

func convertToDataTimePoint(value string) (xdr.ScValType, interface{}, error) {
	data, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	return xdr.ScValTypeScvTimepoint, xdr.TimePoint(xdr.Uint64(data)), err
}

func convertToDataDuration(value string) (xdr.ScValType, interface{}, error) {
	data, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	return xdr.ScValTypeScvDuration, xdr.Duration(xdr.Uint64(data)), err
}

func convertToDataUInt128Parts(value1 string, value2 string) (xdr.ScValType, interface{}, error) {
	data1, err := strconv.ParseUint(value1, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	data2, err := strconv.ParseUint(value2, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	var data xdr.UInt128Parts = xdr.UInt128Parts{
		Hi: xdr.Uint64(data1),
		Lo: xdr.Uint64(data2),
	}

	return xdr.ScValTypeScvU128, data, err
}

func convertToDataInt128Parts(value1 string, value2 string) (xdr.ScValType, interface{}, error) {
	data1, err := strconv.ParseInt(value1, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	data2, err := strconv.ParseUint(value2, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	var data xdr.Int128Parts = xdr.Int128Parts{
		Hi: xdr.Int64(data1),
		Lo: xdr.Uint64(data2),
	}

	return xdr.ScValTypeScvI128, data, err
}

func convertToDataUInt256Parts(value1 string, value2 string, value3 string, value4 string) (xdr.ScValType, interface{}, error) {
	data1, err := strconv.ParseUint(value1, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	data2, err := strconv.ParseUint(value2, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	data3, err := strconv.ParseUint(value3, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	data4, err := strconv.ParseUint(value4, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	var data xdr.UInt256Parts = xdr.UInt256Parts{
		HiHi: xdr.Uint64(data1),
		HiLo: xdr.Uint64(data2),
		LoHi: xdr.Uint64(data3),
		LoLo: xdr.Uint64(data4),
	}

	return xdr.ScValTypeScvU256, data, err
}

func convertToDataInt256Parts(value1 string, value2 string, value3 string, value4 string) (xdr.ScValType, interface{}, error) {
	data1, err := strconv.ParseInt(value1, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	data2, err := strconv.ParseUint(value2, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	data3, err := strconv.ParseUint(value3, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	data4, err := strconv.ParseUint(value4, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	var data xdr.Int256Parts = xdr.Int256Parts{
		HiHi: xdr.Int64(data1),
		HiLo: xdr.Uint64(data2),
		LoHi: xdr.Uint64(data3),
		LoLo: xdr.Uint64(data4),
	}

	return xdr.ScValTypeScvI256, data, err
}

func convertToDataScBytes(value string) (xdr.ScValType, interface{}, error) {
	return xdr.ScValTypeScvBytes, xdr.ScBytes([]byte(value)), nil
}

func convertToDataScString(value string) (xdr.ScValType, interface{}, error) {
	return xdr.ScValTypeScvString, xdr.ScString(value), nil
}

func convertToDataScSymbol(value string) (xdr.ScValType, interface{}, error) {
	return xdr.ScValTypeScvSymbol, xdr.ScSymbol(value), nil
}

func convertToDataScNonceKey(value string) (xdr.ScValType, interface{}, error) {
	data1, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, nil, err
	}

	var data xdr.ScNonceKey = xdr.ScNonceKey{
		Nonce: xdr.Int64(data1),
	}

	return xdr.ScValTypeScvLedgerKeyNonce, data, err
}
