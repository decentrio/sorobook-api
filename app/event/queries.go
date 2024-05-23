package event

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
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = PAGE_SIZE
	}

	offset := (page - 1) * pageSize

	var events []*types.Event
	err := k.dbHandler.Table(EVENT_TABLE).
		Where("contract_id = ?", request.ContractId).
		Joins("JOIN transactions ON transactions.hash = wasm_contract_events.tx_hash").
		Order("transactions.ledger DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&events).Error
	if err != nil {
		return &types.ContractEventsResponse{}, err
	}

	var infos []*types.EventInfo
	for _, item := range events {
		eventInfo, err := convertToEventInfo(item)
		if err != nil {
			return &types.ContractEventsResponse{}, err
		}
		infos = append(infos, eventInfo)
	}

	return &types.ContractEventsResponse{
		Events: infos,
	}, nil
}

func (k Keeper) EventsAtLedger(ctx context.Context, request *types.EventsAtLedgerRequest) (*types.EventsAtLedgerResponse, error) {
	var events []*types.Event
	err := k.dbHandler.Table(EVENT_TABLE).
		Joins("JOIN transactions ON transactions.hash = wasm_contract_events.tx_hash").
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

func (k Keeper) TransferEvents(ctx context.Context, request *types.TransferEventsRequest) (*types.TransferEventsResponse, error) {
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = PAGE_SIZE
	}

	offset := (page - 1) * pageSize
	var events []*types.TranferEvent
	err := k.dbHandler.Table(TRANSFER_TABLE).
		Where("contract_id = ?", request.ContractId).
		Joins("JOIN transactions ON transactions.hash = asset_contract_transfer_events.tx_hash").
		Order("transactions.ledger DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&events).Error
	if err != nil {
		return &types.TransferEventsResponse{}, err
	}

	return &types.TransferEventsResponse{
		Events: events,
	}, nil
}

func (k Keeper) TransferEventsFrom(ctx context.Context, request *types.TransferEventsFromRequest) (*types.TransferEventsFromResponse, error) {
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = PAGE_SIZE
	}

	offset := (page - 1) * pageSize
	var events []*types.TranferEvent
	err := k.dbHandler.Table(TRANSFER_TABLE).
		Where("from_addr = ?", request.From).
		Joins("JOIN transactions ON transactions.hash = asset_contract_transfer_events.tx_hash").
		Order("transactions.ledger DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&events).Error
	if err != nil {
		return &types.TransferEventsFromResponse{}, err
	}

	return &types.TransferEventsFromResponse{
		Events: events,
	}, nil
}

func (k Keeper) TransferEventsTo(ctx context.Context, request *types.TransferEventsToRequest) (*types.TransferEventsToResponse, error) {
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = PAGE_SIZE
	}

	offset := (page - 1) * pageSize
	var events []*types.TranferEvent
	err := k.dbHandler.Table(TRANSFER_TABLE).
		Where("to_addr = ?", request.To).
		Joins("JOIN transactions ON transactions.hash = asset_contract_transfer_events.tx_hash").
		Order("transactions.ledger DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&events).Error
	if err != nil {
		return &types.TransferEventsToResponse{}, err
	}

	return &types.TransferEventsToResponse{
		Events: events,
	}, nil
}

func (k Keeper) MintEvents(ctx context.Context, request *types.MintEventsRequest) (*types.MintEventsResponse, error) {
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = PAGE_SIZE
	}

	offset := (page - 1) * pageSize
	var events []*types.MintEvent
	err := k.dbHandler.Table(MINT_TABLE).
		Where("contract_id = ?", request.ContractId).
		Joins("JOIN transactions ON transactions.hash = asset_contract_mint_events.tx_hash").
		Order("transactions.ledger DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&events).Error
	if err != nil {
		return &types.MintEventsResponse{}, err
	}

	return &types.MintEventsResponse{
		Events: events,
	}, nil
}

func (k Keeper) MintEventsAdmin(ctx context.Context, request *types.MintEventsAdminRequest) (*types.MintEventsAdminResponse, error) {
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = PAGE_SIZE
	}

	offset := (page - 1) * pageSize
	var events []*types.MintEvent
	err := k.dbHandler.Table(MINT_TABLE).
		Where("admin_addr = ?", request.Admin).
		Joins("JOIN transactions ON transactions.hash = asset_contract_mint_events.tx_hash").
		Order("transactions.ledger DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&events).Error
	if err != nil {
		return &types.MintEventsAdminResponse{}, err
	}

	return &types.MintEventsAdminResponse{
		Events: events,
	}, nil
}

func (k Keeper) MintEventsTo(ctx context.Context, request *types.MintEventsToRequest) (*types.MintEventsToResponse, error) {
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = PAGE_SIZE
	}

	offset := (page - 1) * pageSize
	var events []*types.MintEvent
	err := k.dbHandler.Table(MINT_TABLE).
		Where("to_addr = ?", request.To).
		Joins("JOIN transactions ON transactions.hash = asset_contract_mint_events.tx_hash").
		Order("transactions.ledger DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&events).Error
	if err != nil {
		return &types.MintEventsToResponse{}, err
	}

	return &types.MintEventsToResponse{
		Events: events,
	}, nil
}

func (k Keeper) BurnEvents(ctx context.Context, request *types.BurnEventsRequest) (*types.BurnEventsResponse, error) {
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = PAGE_SIZE
	}

	offset := (page - 1) * pageSize
	var events []*types.BurnEvent
	err := k.dbHandler.Table(BURN_TABLE).
		Where("contract_id = ?", request.ContractId).
		Joins("JOIN transactions ON transactions.hash = asset_contract_burn_events.tx_hash").
		Order("transactions.ledger DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&events).Error
	if err != nil {
		return &types.BurnEventsResponse{}, err
	}

	return &types.BurnEventsResponse{
		Events: events,
	}, nil
}

func (k Keeper) BurnEventsFrom(ctx context.Context, request *types.BurnEventsFromRequest) (*types.BurnEventsFromResponse, error) {
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = PAGE_SIZE
	}

	offset := (page - 1) * pageSize
	var events []*types.BurnEvent
	err := k.dbHandler.Table(BURN_TABLE).
		Where("from_addr = ?", request.From).
		Joins("JOIN transactions ON transactions.hash = asset_contract_burn_events.tx_hash").
		Order("transactions.ledger DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&events).Error
	if err != nil {
		return &types.BurnEventsFromResponse{}, err
	}

	return &types.BurnEventsFromResponse{
		Events: events,
	}, nil
}

func (k Keeper) ClawbackEvents(ctx context.Context, request *types.ClawbackEventsRequest) (*types.ClawbackEventsResponse, error) {
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = PAGE_SIZE
	}

	offset := (page - 1) * pageSize
	var events []*types.ClawbackEvent
	err := k.dbHandler.Table(CLAWBACK_TABLE).
		Where("contract_id = ?", request.ContractId).
		Joins("JOIN transactions ON transactions.hash = asset_contract_clawback_events.tx_hash").
		Order("transactions.ledger DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&events).Error
	if err != nil {
		return &types.ClawbackEventsResponse{}, err
	}

	return &types.ClawbackEventsResponse{
		Events: events,
	}, nil
}

func (k Keeper) ClawbackEventsAdmin(ctx context.Context, request *types.ClawbackEventsAdminRequest) (*types.ClawbackEventsAdminResponse, error) {
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = PAGE_SIZE
	}

	offset := (page - 1) * pageSize
	var events []*types.ClawbackEvent
	err := k.dbHandler.Table(CLAWBACK_TABLE).
		Where("admin_addr = ?", request.Admin).
		Joins("JOIN transactions ON transactions.hash = asset_contract_clawback_events.tx_hash").
		Order("transactions.ledger DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&events).Error
	if err != nil {
		return &types.ClawbackEventsAdminResponse{}, err
	}

	return &types.ClawbackEventsAdminResponse{
		Events: events,
	}, nil
}

func (k Keeper) ClawbackEventsFrom(ctx context.Context, request *types.ClawbackEventsFromRequest) (*types.ClawbackEventsFromResponse, error) {
	page := int(request.Page)
	if request.Page < 1 {
		page = 1
	}
	pageSize := int(request.PageSize)
	if request.PageSize < 1 {
		pageSize = PAGE_SIZE
	}

	offset := (page - 1) * pageSize
	var events []*types.ClawbackEvent
	err := k.dbHandler.Table(CLAWBACK_TABLE).
		Where("from_addr = ?", request.From).
		Joins("JOIN transactions ON transactions.hash = asset_contract_clawback_events.tx_hash").
		Order("transactions.ledger DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&events).Error
	if err != nil {
		return &types.ClawbackEventsFromResponse{}, err
	}

	return &types.ClawbackEventsFromResponse{
		Events: events,
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
