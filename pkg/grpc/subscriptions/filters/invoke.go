package filters

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Invoke -
type Invoke struct {
	fltrs     []*pb.InvokeFilters
	isEmpty   bool
	contracts []ids
}

// NewInvoke -
func NewInvoke(ctx context.Context, address storage.IAddress, req []*pb.InvokeFilters) (Invoke, error) {
	invoke := Invoke{
		isEmpty: true,
	}
	if req == nil {
		return invoke, nil
	}
	invoke.isEmpty = false
	invoke.fltrs = req
	invoke.contracts = make([]ids, 0)
	for i := range invoke.fltrs {
		invoke.contracts = append(invoke.contracts, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, invoke.fltrs[i].Contract, invoke.contracts[i]); err != nil {
			return invoke, err
		}
	}
	return invoke, nil
}

// Filter -
func (f Invoke) Filter(data storage.Invoke) bool {
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

		if !validEnum(f.fltrs[i].Version, data.Version) {
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

		if !validEquality(f.fltrs[i].Selector, encoding.EncodeHex(data.EntrypointSelector)) {
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
