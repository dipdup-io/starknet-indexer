package filters

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Event -
type Event struct {
	*pb.EventFilter
	isEmpty bool

	contracts ids
	from      ids
}

// NewEvent -
func NewEvent(ctx context.Context, address storage.IAddress, req *pb.EventFilter) (Event, error) {
	event := Event{
		isEmpty:   true,
		contracts: make(ids),
		from:      make(ids),
	}
	if req == nil {
		return event, nil
	}
	event.isEmpty = false
	event.EventFilter = req
	if err := fillAddressMapFromBytesFilter(ctx, address, req.Contract, event.contracts); err != nil {
		return event, err
	}
	if err := fillAddressMapFromBytesFilter(ctx, address, req.From, event.from); err != nil {
		return event, err
	}
	return event, nil
}

// Filter -
func (f *Event) Filter(data storage.Event) bool {
	if f.isEmpty {
		return true
	}
	if f.EventFilter == nil {
		return false
	}

	if !validInteger(f.Id, data.ID) {
		return false
	}

	if !validInteger(f.Height, data.Height) {
		return false
	}

	if !validTime(f.Time, data.Time) {
		return false
	}

	if f.Contract != nil {
		if !f.contracts.In(data.ContractID) {
			return false
		}
	}

	if f.From != nil {
		if !f.from.In(data.FromID) {
			return false
		}
	}

	if !validString(f.Name, data.Name) {
		return false
	}

	if !validMap(f.ParsedData, data.ParsedData) {
		return false
	}

	return true
}
