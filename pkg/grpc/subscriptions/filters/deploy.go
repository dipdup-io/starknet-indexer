package filters

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Deploy -
type Deploy struct {
	*pb.DeployFilters
	isEmpty bool

	classes ids
}

// NewDeploy -
func NewDeploy(ctx context.Context, class storage.IClass, req *pb.DeployFilters) (Deploy, error) {
	deploy := Deploy{
		isEmpty: true,
		classes: make(ids),
	}
	if req == nil {
		return deploy, nil
	}
	deploy.isEmpty = false
	deploy.DeployFilters = req
	if err := fillClassMapFromBytesFilter(ctx, class, req.Class, deploy.classes); err != nil {
		return deploy, err
	}
	return deploy, nil
}

// Filter -
func (f Deploy) Filter(data storage.Deploy) bool {
	if f.isEmpty {
		return true
	}
	if f.DeployFilters == nil {
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
