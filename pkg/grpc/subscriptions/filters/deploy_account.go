package filters

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// DeployAccount -
type DeployAccount struct {
	*pb.DeployAccountFilters
	isEmpty bool

	classes ids
}

// NewDeployAccount -
func NewDeployAccount(ctx context.Context, class storage.IClass, req *pb.DeployAccountFilters) (DeployAccount, error) {
	deploy := DeployAccount{
		isEmpty: true,
		classes: make(ids),
	}
	if req == nil {
		return deploy, nil
	}
	deploy.isEmpty = false
	deploy.DeployAccountFilters = req

	if err := fillClassMapFromBytesFilter(ctx, class, req.Class, deploy.classes); err != nil {
		return deploy, err
	}

	return deploy, nil
}

// Filter -
func (f DeployAccount) Filter(data storage.DeployAccount) bool {
	if f.isEmpty {
		return true
	}
	if f.DeployAccountFilters == nil {
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

	if f.Class != nil {
		if !f.classes.In(data.ClassID) {
			return false
		}
	}

	if !validMap(f.ParsedCalldata, data.ParsedCalldata) {
		return false
	}

	return false
}
