package filters

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// L1Handler -
type L1Handler struct {
	*pb.L1HandlerFilter
	isEmpty bool

	contracts ids
}

// NewL1Handler -
func NewL1Handler(ctx context.Context, address storage.IAddress, req *pb.L1HandlerFilter) (L1Handler, error) {
	l1Handler := L1Handler{
		isEmpty:   true,
		contracts: make(ids),
	}
	if req == nil {
		return l1Handler, nil
	}
	l1Handler.isEmpty = false
	l1Handler.L1HandlerFilter = req

	if err := fillAddressMapFromBytesFilter(ctx, address, req.Contract, l1Handler.contracts); err != nil {
		return l1Handler, err
	}
	return l1Handler, nil
}

// Filter -
func (f L1Handler) Filter(data storage.L1Handler) bool {
	if f.isEmpty {
		return true
	}
	if f.L1HandlerFilter == nil {
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

	if !validString(f.Entrypoint, data.Entrypoint) {
		return false
	}

	if !validEquality(f.Selector, encoding.EncodeHex(data.Selector)) {
		return false
	}

	if !validMap(f.ParsedCalldata, data.ParsedCalldata) {
		return false
	}

	return false
}
