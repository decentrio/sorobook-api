package app

import (
	"context"

	types "github.com/decentrio/sorobook-api/types/v1"
)

func (k Keeper) Event(ctx context.Context, request *types.EventRequest) (*types.EventResponse, error) {
	return nil, nil
}

func (k Keeper) ContractEvents(ctx context.Context, request *types.ContractEventsRequest) (*types.ContractEventsResponse, error) {
	return nil, nil
}

func (k Keeper) ContractEventCount(ctx context.Context, request *types.ContractEventCountRequest) (*types.ContractEventCountResponse, error) {
	return nil, nil
}
