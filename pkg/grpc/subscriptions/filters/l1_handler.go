package filters

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// L1Handler -
type L1Handler struct {
	fltrs   []*pb.L1HandlerFilter
	isEmpty bool

	contracts []ids
}

// NewL1Handler -
func NewL1Handler(ctx context.Context, address storage.IAddress, req []*pb.L1HandlerFilter) (L1Handler, error) {
	l1Handler := L1Handler{
		isEmpty: true,
	}
	if req == nil {
		return l1Handler, nil
	}
	l1Handler.isEmpty = false
	l1Handler.fltrs = req
	l1Handler.contracts = make([]ids, 0)

	for i := range l1Handler.fltrs {
		l1Handler.contracts = append(l1Handler.contracts, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, l1Handler.fltrs[i].Contract, l1Handler.contracts[i]); err != nil {
			return l1Handler, err
		}
	}
	return l1Handler, nil
}

// Filter -
func (f L1Handler) Filter(data storage.L1Handler) bool {
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

		if !validString(f.fltrs[i].Entrypoint, data.Entrypoint) {
			continue
		}

		if !validEquality(f.fltrs[i].Selector, encoding.EncodeHex(data.Selector)) {
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
