package filters

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Message -
type Message struct {
	fltrs   []*pb.MessageFilter
	isEmpty bool

	contracts []ids
	from      []ids
	to        []ids
}

// NewMessage -
func NewMessage(ctx context.Context, address storage.IAddress, req []*pb.MessageFilter) (Message, error) {
	msg := Message{
		isEmpty: true,
	}
	if req == nil {
		return msg, nil
	}
	msg.isEmpty = false
	msg.fltrs = req
	msg.contracts = make([]ids, 0)
	msg.from = make([]ids, 0)
	msg.to = make([]ids, 0)
	for i := range msg.fltrs {
		msg.contracts = append(msg.contracts, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, msg.fltrs[i].Contract, msg.contracts[i]); err != nil {
			return msg, err
		}
		msg.from = append(msg.from, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, msg.fltrs[i].From, msg.from[i]); err != nil {
			return msg, err
		}
		msg.to = append(msg.to, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, msg.fltrs[i].To, msg.to[i]); err != nil {
			return msg, err
		}
	}
	return msg, nil
}

// Filter -
func (f Message) Filter(data storage.Message) bool {
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

		if f.fltrs[i].To != nil {
			if !f.to[i].In(data.ToID) {
				continue
			}
		}

		if !validEquality(f.fltrs[i].Selector, data.Selector) {
			continue
		}

		ok = true
		break
	}

	return ok
}
