package filters

import (
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Invoke -
type Invoke struct {
	*pb.InvokeFilters
	isEmpty bool
}

// NewInvoke -
func NewInvoke(req *pb.InvokeFilters) Invoke {
	invoke := Invoke{
		isEmpty: true,
	}
	if req == nil {
		return invoke
	}
	invoke.isEmpty = false
	invoke.InvokeFilters = req
	return invoke
}

// Filter -
func (f Invoke) Filter(data storage.Invoke) bool {
	if f.isEmpty {
		return true
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

	if !validBytes(f.Contract, data.Contract.Hash) {
		return false
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
