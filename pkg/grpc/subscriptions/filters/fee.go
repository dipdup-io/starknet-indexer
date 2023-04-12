package filters

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Fee -
type Fee struct {
	*pb.FeeFilter
	isEmpty bool

	contracts ids
	callers   ids
	class     ids
}

// NewFee -
func NewFee(ctx context.Context, address storage.IAddress, class storage.IClass, req *pb.FeeFilter) (Fee, error) {
	fee := Fee{
		isEmpty:   true,
		contracts: make(ids),
		callers:   make(ids),
		class:     make(ids),
	}
	if req == nil {
		return fee, nil
	}
	fee.isEmpty = false
	fee.FeeFilter = req

	if err := fillAddressMapFromBytesFilter(ctx, address, req.Contract, fee.contracts); err != nil {
		return fee, err
	}
	if err := fillAddressMapFromBytesFilter(ctx, address, req.Caller, fee.callers); err != nil {
		return fee, err
	}
	if err := fillClassMapFromBytesFilter(ctx, class, req.Class, fee.class); err != nil {
		return fee, err
	}

	return fee, nil
}

// Filter -
func (f Fee) Filter(data storage.Fee) bool {
	if f.isEmpty {
		return true
	}
	if f.FeeFilter == nil {
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
