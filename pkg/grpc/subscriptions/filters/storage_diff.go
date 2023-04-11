package filters

import (
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// StorageDiff -
type StorageDiff struct {
	*pb.StorageDiffFilter
	isEmpty bool
}

// NewStorageDiff -
func NewStorageDiff(req *pb.StorageDiffFilter) StorageDiff {
	sd := StorageDiff{
		isEmpty: true,
	}
	if req == nil {
		return sd
	}
	sd.isEmpty = false
	sd.StorageDiffFilter = req
	return sd
}

// Filter -
func (f StorageDiff) Filter(data storage.StorageDiff) bool {
	if f.isEmpty {
		return true
	}

	if !validInteger(f.Id, data.ID) {
		return false
	}

	if !validInteger(f.Height, data.Height) {
		return false
	}

	// TODO: think about passing contract address
	if !validBytes(f.Contract, data.Contract.Hash) {
		return false
	}

	if !validEquality(f.Key, encoding.EncodeHex(data.Key)) {
		return false
	}

	return true
}
