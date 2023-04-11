package filters

import (
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// L1Handler -
type L1Handler struct {
	*pb.L1HandlerFilter
	isEmpty bool
}

// NewL1Handler -
func NewL1Handler(req *pb.L1HandlerFilter) L1Handler {
	l1Handler := L1Handler{
		isEmpty: true,
	}
	if req == nil {
		return l1Handler
	}
	l1Handler.isEmpty = false
	l1Handler.L1HandlerFilter = req
	return l1Handler
}

// Filter -
func (f L1Handler) Filter(data storage.L1Handler) bool {
	if f.isEmpty {
		return true
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

	if !validBytes(f.Contract, data.Contract.Hash) {
		return false
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
