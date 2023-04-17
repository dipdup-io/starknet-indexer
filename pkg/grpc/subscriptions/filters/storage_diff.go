package filters

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// StorageDiff -
type StorageDiff struct {
	*pb.StorageDiffFilter
	isEmpty bool

	contracts ids
}

// NewStorageDiff -
func NewStorageDiff(ctx context.Context, address storage.IAddress, req *pb.StorageDiffFilter) (StorageDiff, error) {
	sd := StorageDiff{
		isEmpty:   true,
		contracts: make(ids),
	}
	if req == nil {
		return sd, nil
	}
	sd.isEmpty = false
	sd.StorageDiffFilter = req
	if err := fillAddressMapFromBytesFilter(ctx, address, req.Contract, sd.contracts); err != nil {
		return sd, err
	}
	return sd, nil
}

// Filter -
func (f StorageDiff) Filter(data storage.StorageDiff) bool {
	if f.isEmpty {
		return true
	}
	if f.StorageDiffFilter == nil {
		return false
	}

	if !validInteger(f.Id, data.ID) {
		return false
	}

	if !validInteger(f.Height, data.Height) {
		return false
	}

	if f.Contract != nil {
		if !f.contracts.In(data.ContractID) {
			return false
		}
	}

	if !validEquality(f.Key, encoding.EncodeHex(data.Key)) {
		return false
	}

	return true
}
