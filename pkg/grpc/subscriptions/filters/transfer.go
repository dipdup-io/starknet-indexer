package filters

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Transfer -
type Transfer struct {
	*pb.TransferFilter
	isEmpty bool
}

// NewTransfer -
func NewTransfer(req *pb.TransferFilter) Transfer {
	transfer := Transfer{
		isEmpty: true,
	}
	if req == nil {
		return transfer
	}
	transfer.isEmpty = false
	transfer.TransferFilter = req
	return transfer
}

// Filter -
func (f Transfer) Filter(data storage.Transfer) bool {
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

	if !validBytes(f.Contract, data.Contract.Hash) {
		return false
	}

	if !validBytes(f.From, data.From.Hash) {
		return false
	}

	if !validBytes(f.To, data.To.Hash) {
		return false
	}

	if !validString(f.TokenId, data.TokenID.String()) {
		return false
	}

	return true
}
