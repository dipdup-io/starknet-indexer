package filters

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Deploy -
type Deploy struct {
	*pb.DeployFilters
	isEmpty bool
}

// NewDeploy -
func NewDeploy(req *pb.DeployFilters) Deploy {
	deploy := Deploy{
		isEmpty: true,
	}
	if req == nil {
		return deploy
	}
	deploy.isEmpty = false
	deploy.DeployFilters = req
	return deploy
}

// Filter -
func (f Deploy) Filter(data storage.Deploy) bool {
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

	if !validBytes(f.Class, data.Class.Hash) {
		return false
	}

	if !validMap(f.ParsedCalldata, data.ParsedCalldata) {
		return false
	}

	return false
}
