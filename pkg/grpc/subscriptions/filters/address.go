package filters

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Address -
type Address struct {
	fltrs   []*pb.AddressFilter
	isEmpty bool
}

// NewAddress -
func NewAddress(req []*pb.AddressFilter) Address {
	address := Address{
		isEmpty: true,
	}
	if req == nil {
		return address
	}
	address.isEmpty = false
	address.fltrs = req
	return address
}

// Filter -
func (f Address) Filter(data storage.Address) bool {
	if f.isEmpty {
		return true
	}

	var ok bool
	for i := range f.fltrs {
		if f.fltrs[i] == nil {
			continue
		}

		if !validInteger(f.fltrs[i].Id, data.ID) {
			continue
		}
		if !validInteger(f.fltrs[i].Height, data.Height) {
			continue
		}

		if f.fltrs[i].OnlyStarknet && data.ClassID == nil {
			continue
		}

		ok = true
		break
	}

	return ok
}
