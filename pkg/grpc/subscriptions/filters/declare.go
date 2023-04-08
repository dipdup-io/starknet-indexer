package filters

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Declare -
type Declare struct {
	*pb.DeclareFilters
	isEmpty bool
}

// NewDeclare -
func NewDeclare(req *pb.DeclareFilters) Declare {
	declare := Declare{
		isEmpty: true,
	}
	if req == nil {
		return declare
	}
	declare.isEmpty = false
	declare.DeclareFilters = req
	return declare
}

// Filter -
func (f Declare) Filter(data storage.Declare) bool {
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

	return false
}
