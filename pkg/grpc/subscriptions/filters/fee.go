package filters

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Fee -
type Fee struct {
	fltrs   []*pb.FeeFilter
	isEmpty bool

	contracts []ids
	callers   []ids
	class     []ids
}

// NewFee -
func NewFee(ctx context.Context, address storage.IAddress, class storage.IClass, req []*pb.FeeFilter) (Fee, error) {
	fee := Fee{
		isEmpty: true,
	}
	if req == nil {
		return fee, nil
	}
	fee.callers = make([]ids, 0)
	fee.class = make([]ids, 0)
	fee.contracts = make([]ids, 0)
	fee.isEmpty = false
	fee.fltrs = req

	for i := range fee.fltrs {
		fee.contracts = append(fee.contracts, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, fee.fltrs[i].Contract, fee.contracts[i]); err != nil {
			return fee, err
		}
		fee.callers = append(fee.callers, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, fee.fltrs[i].Caller, fee.callers[i]); err != nil {
			return fee, err
		}
		fee.class = append(fee.class, make(ids))
		if err := fillClassMapFromBytesFilter(ctx, class, fee.fltrs[i].Class, fee.class[i]); err != nil {
			return fee, err
		}
	}

	return fee, nil
}

// Filter -
func (f Fee) Filter(data storage.Fee) bool {
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
