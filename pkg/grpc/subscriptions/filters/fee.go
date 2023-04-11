package filters

import (
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Fee -
type Fee struct {
	*pb.FeeFilter
	isEmpty bool
}

// NewFee -
func NewFee(req *pb.FeeFilter) Fee {
	fee := Fee{
		isEmpty: true,
	}
	if req == nil {
		return fee
	}
	fee.isEmpty = false
	fee.FeeFilter = req
	return fee
}

// Filter -
func (f Fee) Filter(data storage.Fee) bool {
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

	if !validBytes(f.Caller, data.Caller.Hash) {
		return false
	}

	if !validBytes(f.Class, data.Class.Hash) {
		return false
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