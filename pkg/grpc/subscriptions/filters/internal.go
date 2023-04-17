package filters

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Internal -
type Internal struct {
	*pb.InternalFilter
	isEmpty   bool
	contracts ids
	callers   ids
	class     ids
}

// NewInternal -
func NewInternal(ctx context.Context, address storage.IAddress, class storage.IClass, req *pb.InternalFilter) (Internal, error) {
	internal := Internal{
		isEmpty:   true,
		contracts: make(ids),
		callers:   make(ids),
		class:     make(ids),
	}
	if req == nil {
		return internal, nil
	}
	internal.isEmpty = false
	internal.InternalFilter = req

	if err := fillAddressMapFromBytesFilter(ctx, address, req.Contract, internal.contracts); err != nil {
		return internal, err
	}
	if err := fillAddressMapFromBytesFilter(ctx, address, req.Caller, internal.callers); err != nil {
		return internal, err
	}
	if err := fillClassMapFromBytesFilter(ctx, class, req.Class, internal.class); err != nil {
		return internal, err
	}
	return internal, nil
}

// Filter -
func (f Internal) Filter(data storage.Internal) bool {
	if f.isEmpty {
		return true
	}
	if f.InternalFilter == nil {
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

	if f.Contract != nil {
		if !f.contracts.In(data.ContractID) {
			return false
		}
	}

	if f.Caller != nil {
		if !f.callers.In(data.CallerID) {
			return false
		}
	}

	if f.Class != nil {
		if !f.class.In(data.ClassID) {
			return false
		}
	}

	if !validString(f.Entrypoint, data.Entrypoint) {
		return false
	}

	if !validEquality(f.Selector, encoding.EncodeHex(data.Selector)) {
		return false
	}

	if !validEnum(f.EntrypointType, uint64(data.EntrypointType)) {
		return false
	}

	if !validEnum(f.CallType, uint64(data.CallType)) {
		return false
	}

	if !validMap(f.ParsedCalldata, data.ParsedCalldata) {
		return false
	}

	return false
}
