package filters

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// DeployAccount -
type DeployAccount struct {
	*pb.DeployAccountFilters
	isEmpty bool
}

// NewDeployAccount -
func NewDeployAccount(req *pb.DeployAccountFilters) DeployAccount {
	deploy := DeployAccount{
		isEmpty: true,
	}
	if req == nil {
		return deploy
	}
	deploy.isEmpty = false
	deploy.DeployAccountFilters = req
	return deploy
}

// Filter -
func (f DeployAccount) Filter(data storage.DeployAccount) bool {
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
