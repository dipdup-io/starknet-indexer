package filters

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Transfer -
type Transfer struct {
	fltrs   []*pb.TransferFilter
	isEmpty bool

	contracts []ids
	from      []ids
	to        []ids
}

// NewTransfer -
func NewTransfer(ctx context.Context, address storage.IAddress, req []*pb.TransferFilter) (Transfer, error) {
	transfer := Transfer{
		isEmpty: true,
	}
	if req == nil {
		return transfer, nil
	}
	transfer.isEmpty = false
	transfer.fltrs = req
	transfer.contracts = make([]ids, 0)
	transfer.from = make([]ids, 0)
	transfer.to = make([]ids, 0)

	for i := range transfer.fltrs {
		transfer.contracts = append(transfer.contracts, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, transfer.fltrs[i].Contract, transfer.contracts[i]); err != nil {
			return transfer, err
		}
		transfer.from = append(transfer.from, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, transfer.fltrs[i].From, transfer.from[i]); err != nil {
			return transfer, err
		}
		transfer.to = append(transfer.to, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, transfer.fltrs[i].To, transfer.to[i]); err != nil {
			return transfer, err
		}
	}
	return transfer, nil
}

// Filter -
func (f Transfer) Filter(data storage.Transfer) bool {
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

		if !validString(f.fltrs[i].TokenId, data.TokenID.String()) {
			continue
		}

		ok = true
		break
	}

	return ok
}
