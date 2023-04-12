package filters

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Transfer -
type Transfer struct {
	*pb.TransferFilter
	isEmpty bool

	contracts ids
	from      ids
	to        ids
}

// NewTransfer -
func NewTransfer(ctx context.Context, address storage.IAddress, req *pb.TransferFilter) (Transfer, error) {
	transfer := Transfer{
		isEmpty: true,
	}
	if req == nil {
		return transfer, nil
	}
	transfer.isEmpty = false
	transfer.TransferFilter = req

	if err := fillAddressMapFromBytesFilter(ctx, address, req.Contract, transfer.contracts); err != nil {
		return transfer, err
	}
	if err := fillAddressMapFromBytesFilter(ctx, address, req.From, transfer.from); err != nil {
		return transfer, err
	}
	if err := fillAddressMapFromBytesFilter(ctx, address, req.To, transfer.to); err != nil {
		return transfer, err
	}
	return transfer, nil
}

// Filter -
func (f Transfer) Filter(data storage.Transfer) bool {
	if f.isEmpty {
		return true
	}
	if f.TransferFilter == nil {
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

	if !validString(f.TokenId, data.TokenID.String()) {
		return false
	}

	return true
}
