package app

import (
	"context"
	"encoding/json"

	"google.golang.org/protobuf/types/known/structpb"

	types "github.com/decentrio/sorobook-api/types/v1"
	"github.com/decentrio/xdr-converter/converter"
)

func (k Keeper) Transaction(ctx context.Context, request *types.TransactionRequest) (*types.TransactionResponse, error) {
	var transaction types.Transaction

	err := k.dbHandler.Table(TRANSACTION_TABLE).Where("hash = ?", request.Hash).First(&transaction).Error
	if err != nil {
		return &types.TransactionResponse{
			Found:       false,
			Transaction: &types.TransactionInfo{},
		}, err
	}

	txInfo, err := convertToTxInfo(&transaction)
	if err != nil {
		return &types.TransactionResponse{
			Found:       false,
			Transaction: &types.TransactionInfo{},
		}, err
	}

	return &types.TransactionResponse{
		Found:       false,
		Transaction: txInfo,
	}, nil
}

func (k Keeper) TransactionsAtLedgerSeq(ctx context.Context, request *types.TransactionsAtLedgerSeqRequest) (*types.TransactionsAtLedgerSeqResponse, error) {
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = PAGE_SIZE
	}

	offset := (page - 1) * pageSize
	var txs []*types.Transaction
	err := k.dbHandler.Table(TRANSACTION_TABLE).Where("ledger = ?", request.Ledger).Limit(pageSize).Offset(offset).Find(&txs).Error
	if err != nil {
		return &types.TransactionsAtLedgerSeqResponse{}, err
	}

	var infos []*types.TransactionInfo

	for _, item := range txs {
		txInfo, err := convertToTxInfo(item)
		if err != nil {
			return &types.TransactionsAtLedgerSeqResponse{}, err
		}
		infos = append(infos, txInfo)
	}

	return &types.TransactionsAtLedgerSeqResponse{
		Txs:  infos,
	}, nil
}

func (k Keeper) TransactionsAtLedgerHash(ctx context.Context, request *types.TransactionsAtLedgerHashRequest) (*types.TransactionsAtLedgerHashResponse, error) {
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = PAGE_SIZE
	}

	offset := (page - 1) * pageSize
	var txs []*types.Transaction
	err := k.dbHandler.Table(TRANSACTION_TABLE).Joins("JOIN ledgers ON transactions.ledger = ledgers.seq").
		Where("ledgers.hash = ?", request.LedgerHash).
		Limit(pageSize).
		Offset(offset).
		Find(&txs).Error
	if err != nil {
		return &types.TransactionsAtLedgerHashResponse{}, err
	}

	var infos []*types.TransactionInfo

	for _, item := range txs {
		txInfo, err := convertToTxInfo(item)
		if err != nil {
			return &types.TransactionsAtLedgerHashResponse{}, err
		}
		infos = append(infos, txInfo)
	}

	return &types.TransactionsAtLedgerHashResponse{
		Txs:  infos,
	}, nil
}

func (k Keeper) TransactionsByAddress(ctx context.Context, request *types.TransactionsByAddressRequest) (*types.TransactionsByAddressResponse, error) {
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = PAGE_SIZE
	}

	offset := (page - 1) * pageSize
	var txs []*types.Transaction
	err := k.dbHandler.Table(TRANSACTION_TABLE).Where("source_address = ?", request.Address).Order("ledger DESC").Limit(pageSize).Offset(offset).Find(&txs).Error
	if err != nil {
		return &types.TransactionsByAddressResponse{
			Txs:  []*types.TransactionInfo{},
		}, err
	}

	var infos []*types.TransactionInfo

	for _, item := range txs {
		txInfo, err := convertToTxInfo(item)
		if err != nil {
			return &types.TransactionsByAddressResponse{}, err
		}
		infos = append(infos, txInfo)
	}

	return &types.TransactionsByAddressResponse{
		Txs:  infos,
	}, nil
}

func (k Keeper) ContractTransactions(ctx context.Context, request *types.ContractTransactionsRequest) (*types.ContractTransactionsResponse, error) {
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = PAGE_SIZE
	}

	offset := (page - 1) * pageSize
	var txs []*types.Transaction
	err := k.dbHandler.Table(TRANSACTION_TABLE).
		Joins("JOIN contracts ON transactions.hash = contracts.tx_hash").
		Where("contract_id = ?", request.Contract).
		Order("ledger DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&txs).Error
	if err != nil {
		return &types.ContractTransactionsResponse{}, err
	}

	var infos []*types.TransactionInfo

	for _, item := range txs {
		txInfo, err := convertToTxInfo(item)
		if err != nil {
			return &types.ContractTransactionsResponse{}, err
		}
		infos = append(infos, txInfo)
	}

	return &types.ContractTransactionsResponse{
		Txs:  infos,
	}, nil
}

func (k Keeper) UserContractTransactions(ctx context.Context, request *types.UserContractTransactionsRequest) (*types.UserContractTransactionsResponse, error) {
	var txs []*types.Transaction
	err := k.dbHandler.Table(TRANSACTION_TABLE).
		Joins("JOIN contracts ON transactions.hash = contracts.tx_hash").
		Where("contract_id = ?", request.Contract).
		Where("source_address = ?", request.Address).
		Order("ledger DESC").
		Find(&txs).Error
	if err != nil {
		return &types.UserContractTransactionsResponse{
			Txs: []*types.TransactionInfo{},
		}, err
	}

	var infos []*types.TransactionInfo

	for _, item := range txs {
		txInfo, err := convertToTxInfo(item)
		if err != nil {
			return &types.UserContractTransactionsResponse{
				Txs: []*types.TransactionInfo{},
			}, err
		}
		infos = append(infos, txInfo)
	}

	return &types.UserContractTransactionsResponse{
		Txs: infos,
	}, nil
}

func convertToTxInfo(tx *types.Transaction) (*types.TransactionInfo, error) {
	envelopeJson := &structpb.Struct{}
	envelopeData, err := converter.MarshalJSONEnvelopeXdr(tx.EnvelopeXdr)
	if err != nil {
		return &types.TransactionInfo{}, err
	}
	if err := json.Unmarshal(envelopeData, envelopeJson); err != nil {
		return &types.TransactionInfo{}, err
	}

	resultMetaJson := &structpb.Struct{}
	resultMetaData, err := converter.MarshalJSONResultMetaXdr(tx.ResultMetaXdr)
	if err != nil {
		return &types.TransactionInfo{}, err
	}
	if err := json.Unmarshal(resultMetaData, resultMetaJson); err != nil {
		return &types.TransactionInfo{}, err
	}

	resultJson := &structpb.Struct{}
	resultData, err := converter.MarshalJSONResultXdr(tx.ResultXdr)
	if err != nil {
		return &types.TransactionInfo{}, err
	}
	if err := json.Unmarshal(resultData, resultJson); err != nil {
		return &types.TransactionInfo{}, err
	}

	return &types.TransactionInfo{
		Hash:             tx.Hash,
		Status:           tx.Status,
		Ledger:           tx.Ledger,
		ApplicationOrder: tx.ApplicationOrder,
		Envelope:         envelopeJson,
		Result:           resultJson,
		ResultMeta:       resultMetaJson,
		SourceAddress:    tx.SourceAddress,
	}, nil
}
