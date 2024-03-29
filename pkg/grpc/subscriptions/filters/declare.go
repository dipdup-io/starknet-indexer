package filters

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Declare -
type Declare struct {
	fltrs   []*pb.DeclareFilters
	isEmpty bool
}

// NewDeclare -
func NewDeclare(req []*pb.DeclareFilters) Declare {
	declare := Declare{
		isEmpty: true,
	}
	if req == nil {
		return declare
	}
	declare.isEmpty = false
	declare.fltrs = req
	return declare
}

// Filter -
func (f Declare) Filter(data storage.Declare) bool {
	if f.isEmpty {
		return true
	}

	var ok bool
	for i := range f.fltrs {
		if f.fltrs[i] == nil {
			continue
		}

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

		ok = true
		break
	}

	return ok
}
