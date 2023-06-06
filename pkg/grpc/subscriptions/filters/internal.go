package filters

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Internal -
type Internal struct {
	fltrs     []*pb.InternalFilter
	isEmpty   bool
	contracts []ids
	callers   []ids
	class     []ids
}

// NewInternal -
func NewInternal(ctx context.Context, address storage.IAddress, class storage.IClass, req []*pb.InternalFilter) (Internal, error) {
	internal := Internal{
		isEmpty: true,
	}
	if req == nil {
		return internal, nil
	}
	internal.callers = make([]ids, 0)
	internal.class = make([]ids, 0)
	internal.contracts = make([]ids, 0)
	internal.isEmpty = false
	internal.fltrs = req

	for i := range internal.fltrs {
		internal.contracts = append(internal.contracts, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, internal.fltrs[i].Contract, internal.contracts[i]); err != nil {
			return internal, err
		}
		internal.callers = append(internal.callers, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, internal.fltrs[i].Caller, internal.callers[i]); err != nil {
			return internal, err
		}
		internal.class = append(internal.class, make(ids))
		if err := fillClassMapFromBytesFilter(ctx, class, internal.fltrs[i].Class, internal.class[i]); err != nil {
			return internal, err
		}
	}
	return internal, nil
}

// Filter -
func (f Internal) Filter(data storage.Internal) bool {
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

		if !validEnum(f.fltrs[i].Status, uint64(data.Status)) {
			continue
		}

		if f.fltrs[i].Contract != nil {
			if !f.contracts[i].In(data.ContractID) {
				continue
			}
		}

		if f.fltrs[i].Caller != nil {
			if !f.callers[i].In(data.CallerID) {
				continue
			}
		}

		if f.fltrs[i].Class != nil {
			if !f.class[i].In(data.ClassID) {
				continue
			}
		}

		if !validString(f.fltrs[i].Entrypoint, data.Entrypoint) {
			continue
		}

		if !validEquality(f.fltrs[i].Selector, encoding.EncodeHex(data.Selector)) {
			continue
		}

		if !validEnum(f.fltrs[i].EntrypointType, uint64(data.EntrypointType)) {
			continue
		}

		if !validEnum(f.fltrs[i].CallType, uint64(data.CallType)) {
			continue
		}

		if !validMap(f.fltrs[i].ParsedCalldata, data.ParsedCalldata) {
			continue
		}

		ok = true
		break
	}

	return ok
}
