package app

import (
	"context"
	"encoding/json"

	"google.golang.org/protobuf/types/known/structpb"

	types "github.com/decentrio/sorobook-api/types/v1"
	"github.com/decentrio/xdr-converter/converter"
)

func (k Keeper) Event(ctx context.Context, request *types.EventRequest) (*types.EventResponse, error) {
	var event types.Event

	err := k.dbHandler.Table(EVENT_TABLE).Where("id = ?", request.Id).First(&event).Error
	if err != nil {
		return &types.EventResponse{
			Found: false,
			Event: &types.EventInfo{},
		}, err
	}

	eventInfo, err := convertToEventInfo(&event)
	if err != nil {
		return &types.EventResponse{
			Found: false,
			Event: &types.EventInfo{},
		}, err
	}

	return &types.EventResponse{
		Found: true,
		Event: eventInfo,
	}, nil
}

func (k Keeper) ContractEvents(ctx context.Context, request *types.ContractEventsRequest) (*types.ContractEventsResponse, error) {
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	offset := (page - 1) * PAGE_SIZE

	var events []*types.Event
	err := k.dbHandler.Table(EVENT_TABLE).
		Where("contract_id = ?", request.ContractId).
		Joins("JOIN transactions ON transactions.hash = wasm_contract_events.tx_hash").
		Order("transactions.ledger DESC").
		Limit(PAGE_SIZE).
		Offset(offset).
		Find(&events).Error
	if err != nil {
		return &types.ContractEventsResponse{
			Events: []*types.EventInfo{},
			Page:   int32(page),
		}, err
	}

	var infos []*types.EventInfo
	for _, item := range events {
		eventInfo, err := convertToEventInfo(item)
		if err != nil {
			return &types.ContractEventsResponse{
				Events: []*types.EventInfo{},
				Page:   int32(page),
			}, err
		}
		infos = append(infos, eventInfo)
	}

	return &types.ContractEventsResponse{
		Events: infos,
		Page:   int32(page),
	}, nil
}

func (k Keeper) EventsAtLedger(ctx context.Context, request *types.EventsAtLedgerRequest) (*types.EventsAtLedgerResponse, error) {
	var events []*types.Event
	err := k.dbHandler.Table(EVENT_TABLE).
		Joins("JOIN transactions ON transactions.hash = events.tx_hash").
		Where("contract_id = ?", request.ContractId).
		Where("transactions.ledger = ?", request.Ledger).
		Find(&events).Error
	if err != nil {
		return &types.EventsAtLedgerResponse{}, err
	}

	var infos []*types.EventInfo
	for _, item := range events {
		eventInfo, err := convertToEventInfo(item)
		if err != nil {
			return &types.EventsAtLedgerResponse{}, err
		}
		infos = append(infos, eventInfo)
	}

	return &types.EventsAtLedgerResponse{
		Events: infos,
	}, nil
}

func (k Keeper) ContractEventCount(ctx context.Context, request *types.ContractEventCountRequest) (*types.ContractEventCountResponse, error) {
	var count int64
	err := k.dbHandler.Table(EVENT_TABLE).Where("contract_id = ?", request.ContractId).Count(&count).Error
	if err != nil {
		return nil, err
	}

	return &types.ContractEventCountResponse{
		Total: count,
	}, nil
}

func convertToEventInfo(event *types.Event) (*types.EventInfo, error) {
	eventJson := &structpb.Struct{}
	eventData, err := converter.MarshalJSONContractEventBodyXdr(event.EventBodyXdr)
	if err != nil {
		return &types.EventInfo{}, err
	}
	if err := json.Unmarshal(eventData, eventJson); err != nil {
		return &types.EventInfo{}, err
	}

	return &types.EventInfo{
		Id:         event.Id,
		ContractId: event.ContractId,
		TxHash:     event.TxHash,
		Event:      eventJson,
	}, err
}
