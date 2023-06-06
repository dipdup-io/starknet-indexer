package filters

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// DeployAccount -
type DeployAccount struct {
	fltrs   []*pb.DeployAccountFilters
	isEmpty bool

	classes []ids
}

// NewDeployAccount -
func NewDeployAccount(ctx context.Context, class storage.IClass, req []*pb.DeployAccountFilters) (DeployAccount, error) {
	deploy := DeployAccount{
		isEmpty: true,
	}
	if req == nil {
		return deploy, nil
	}
	deploy.classes = make([]ids, 0)
	deploy.isEmpty = false
	deploy.fltrs = req

	for i := range deploy.fltrs {
		deploy.classes = append(deploy.classes, make(ids))
		if err := fillClassMapFromBytesFilter(ctx, class, deploy.fltrs[i].Class, deploy.classes[i]); err != nil {
			return deploy, err
		}
	}

	return deploy, nil
}

// Filter -
func (f DeployAccount) Filter(data storage.DeployAccount) bool {
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

		if f.fltrs[i].Class != nil {
			if !f.classes[i].In(data.ClassID) {
				continue
			}
		}

		if !validMap(f.fltrs[i].ParsedCalldata, data.ParsedCalldata) {
			continue
		}

		ok = true
		break
	}

	return ok
}
