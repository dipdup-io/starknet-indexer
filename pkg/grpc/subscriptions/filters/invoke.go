package filters

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Invoke -
type Invoke struct {
	*pb.InvokeFilters
	isEmpty   bool
	contracts ids
}

// NewInvoke -
func NewInvoke(ctx context.Context, address storage.IAddress, req *pb.InvokeFilters) (Invoke, error) {
	invoke := Invoke{
		isEmpty:   true,
		contracts: make(ids),
	}
	if req == nil {
		return invoke, nil
	}
	invoke.isEmpty = false
	invoke.InvokeFilters = req
	if err := fillAddressMapFromBytesFilter(ctx, address, req.Contract, invoke.contracts); err != nil {
		return invoke, err
	}
	return invoke, nil
}

// Filter -
func (f Invoke) Filter(data storage.Invoke) bool {
	if f.isEmpty {
		return true
	}
	if f.InvokeFilters == nil {
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

	if !validEnum(f.Status, uint64(data.Status)) {
		return false
	}

	if !validEnum(f.Version, data.Version) {
		return false
	}

	if f.Contract != nil {
		if !f.contracts.In(data.ContractID) {
			return false
		}
	}

	if !validString(f.Entrypoint, data.Entrypoint) {
		return false
	}

	if !validEquality(f.Selector, encoding.EncodeHex(data.EntrypointSelector)) {
		return false
	}

	if !validMap(f.ParsedCalldata, data.ParsedCalldata) {
		return false
	}

	return false
}
