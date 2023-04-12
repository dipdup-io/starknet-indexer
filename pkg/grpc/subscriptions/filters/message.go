package filters

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Message -
type Message struct {
	*pb.MessageFilter
	isEmpty bool

	contracts ids
	from      ids
	to        ids
}

// NewMessage -
func NewMessage(ctx context.Context, address storage.IAddress, req *pb.MessageFilter) (Message, error) {
	msg := Message{
		isEmpty:   true,
		contracts: make(ids),
		from:      make(ids),
		to:        make(ids),
	}
	if req == nil {
		return msg, nil
	}
	msg.isEmpty = false
	msg.MessageFilter = req
	if err := fillAddressMapFromBytesFilter(ctx, address, req.Contract, msg.contracts); err != nil {
		return msg, err
	}
	if err := fillAddressMapFromBytesFilter(ctx, address, req.From, msg.from); err != nil {
		return msg, err
	}
	if err := fillAddressMapFromBytesFilter(ctx, address, req.To, msg.to); err != nil {
		return msg, err
	}
	return msg, nil
}

// Filter -
func (f Message) Filter(data storage.Message) bool {
	if f.isEmpty {
		return true
	}
	if f.MessageFilter == nil {
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

	if f.To != nil {
		if !f.to.In(data.ToID) {
			return false
		}
	}

	if !validEquality(f.Selector, data.Selector) {
		return false
	}

	return true
}
