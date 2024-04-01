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

	err := k.dbHandler.Table("transactions").Where("hash = ?", request.Hash).First(&transaction).Error
	if err != nil {
		return &types.TransactionResponse{
			Found:       false,
			Transaction: &types.TransactionInfo{},
		}, err
	}

	envelopeJson := &structpb.Struct{}
	envelopeData, err := converter.MarshalJSONEnvelopeXdr(transaction.EnvelopeXdr)
	if err != nil {
		return &types.TransactionResponse{
			Found:       false,
			Transaction: &types.TransactionInfo{},
		}, err
	}
	if err := json.Unmarshal(envelopeData, envelopeJson); err != nil {
		return &types.TransactionResponse{
			Found:       false,
			Transaction: &types.TransactionInfo{},
		}, err
	}

	resultMetaJson := &structpb.Struct{}
	resultMetaData, err := converter.MarshalJSONResultMetaXdr(transaction.ResultMetaXdr)
	if err != nil {
		return &types.TransactionResponse{
			Found:       false,
			Transaction: &types.TransactionInfo{},
		}, err
	}
	if err := json.Unmarshal(resultMetaData, resultMetaJson); err != nil {
		return &types.TransactionResponse{
			Found:       false,
			Transaction: &types.TransactionInfo{},
		}, err
	}

	resultJson := &structpb.Struct{}
	resultData, err := converter.MarshalJSONResultXdr(transaction.ResultXdr)
	if err != nil {
		return &types.TransactionResponse{
			Found:       false,
			Transaction: &types.TransactionInfo{},
		}, err
	}
	if err := json.Unmarshal(resultData, resultJson); err != nil {
		return &types.TransactionResponse{
			Found:       false,
			Transaction: &types.TransactionInfo{},
		}, err
	}

	return &types.TransactionResponse{
		Found: true,
		Transaction: &types.TransactionInfo{
			Hash:             transaction.Hash,
			Status:           transaction.Status,
			Ledger:           transaction.Ledger,
			ApplicationOrder: transaction.ApplicationOrder,
			Envelope:         envelopeJson,
			Result:           resultJson,
			ResultMeta:       resultMetaJson,
			SourceAddress:    transaction.SourceAddress,
		},
	}, nil
}
