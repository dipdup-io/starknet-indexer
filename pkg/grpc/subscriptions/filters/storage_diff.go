package filters

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// StorageDiff -
type StorageDiff struct {
	fltrs   []*pb.StorageDiffFilter
	isEmpty bool

	contracts []ids
}

// NewStorageDiff -
func NewStorageDiff(ctx context.Context, address storage.IAddress, req []*pb.StorageDiffFilter) (StorageDiff, error) {
	sd := StorageDiff{
		isEmpty: true,
	}
	if req == nil {
		return sd, nil
	}
	sd.isEmpty = false
	sd.fltrs = req
	sd.contracts = make([]ids, 0)

	for i := range sd.fltrs {
		sd.contracts = append(sd.contracts, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, sd.fltrs[i].Contract, sd.contracts[i]); err != nil {
			return sd, err
		}
	}
	return sd, nil
}

// Filter -
func (f StorageDiff) Filter(data storage.StorageDiff) bool {
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

		if f.fltrs[i].Contract != nil {
			if !f.contracts[i].In(data.ContractID) {
				continue
			}
		}

		if !validEquality(f.fltrs[i].Key, encoding.EncodeHex(data.Key)) {
			continue
		}

		ok = true
		break
	}

	return ok
}
