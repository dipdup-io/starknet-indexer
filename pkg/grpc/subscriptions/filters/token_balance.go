package filters

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// TokenBalance -
type TokenBalance struct {
	*pb.TokenBalanceFilter
	isEmpty bool
}

// NewTokenBalance -
func NewTokenBalance(req *pb.TokenBalanceFilter) TokenBalance {
	sd := TokenBalance{
		isEmpty: true,
	}
	if req == nil {
		return sd
	}
	sd.isEmpty = false
	sd.TokenBalanceFilter = req
	return sd
}

// Filter -
func (f TokenBalance) Filter(data storage.TokenBalance) bool {
	if f.isEmpty {
		return true
	}

	// TODO: think about passing contract address
	if !validBytes(f.Contract, data.Contract.Hash) {
		return false
	}
	if !validBytes(f.Owner, data.Owner.Hash) {
		return false
	}
	if !validString(f.TokenId, data.TokenID.String()) {
		return false
	}

	return true
}
