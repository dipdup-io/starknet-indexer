package filters

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Event -
type Event struct {
	*pb.EventFilter
	isEmpty bool
}

// NewEvent -
func NewEvent(req *pb.EventFilter) Event {
	event := Event{
		isEmpty: true,
	}
	if req == nil {
		return event
	}
	event.isEmpty = false
	event.EventFilter = req
	return event
}

// Filter -
func (f *Event) Filter(data storage.Event) bool {
	if f.isEmpty {
		return true
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

	// TODO: think about passing contract address
	if !validBytes(f.Contract, data.Contract.Hash) {
		return false
	}

	if !validBytes(f.From, data.From.Hash) {
		return false
	}

	if !validString(f.Name, data.Name) {
		return false
	}

	if !validMap(f.ParsedData, data.ParsedData) {
		return false
	}

	return true
}
