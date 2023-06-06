package filters

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Event -
type Event struct {
	fltrs   []*pb.EventFilter
	isEmpty bool

	contracts []ids
	from      []ids
}

// NewEvent -
func NewEvent(ctx context.Context, address storage.IAddress, req []*pb.EventFilter) (Event, error) {
	event := Event{
		isEmpty: true,
	}
	if req == nil {
		return event, nil
	}
	event.contracts = make([]ids, 0)
	event.from = make([]ids, 0)
	event.isEmpty = false
	event.fltrs = req

	for i := range event.fltrs {
		event.contracts = append(event.contracts, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, event.fltrs[i].Contract, event.contracts[i]); err != nil {
			return event, err
		}
		event.from = append(event.from, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, event.fltrs[i].From, event.from[i]); err != nil {
			return event, err
		}
	}
	return event, nil
}

// Filter -
func (f *Event) Filter(data storage.Event) bool {
	if f.isEmpty {
		return true
	}

	var ok bool
	for i := range f.fltrs {
		if !validInteger(f.fltrs[i].Id, data.ID) {
			continue
		}

		if !validInteger(f.fltrs[i].Height, data.Height) {
			continue
		}

		if !validTime(f.fltrs[i].Time, data.Time) {
			continue
		}

		if f.fltrs[i].Contract != nil {
			if !f.contracts[i].In(data.ContractID) {
				continue
			}
		}

		if f.fltrs[i].From != nil {
			if !f.from[i].In(data.FromID) {
				continue
			}
		}

		if !validString(f.fltrs[i].Name, data.Name) {
			continue
		}

		if !validMap(f.fltrs[i].ParsedData, data.ParsedData) {
			continue
		}

		ok = true
		break
	}

	return ok
}
