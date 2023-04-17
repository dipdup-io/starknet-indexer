package filters

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// TokenBalance -
type TokenBalance struct {
	*pb.TokenBalanceFilter

	contracts ids
	owners    ids

	isEmpty bool
}

// NewTokenBalance -
func NewTokenBalance(ctx context.Context, address storage.IAddress, req *pb.TokenBalanceFilter) (TokenBalance, error) {
	balance := TokenBalance{
		isEmpty:   true,
		contracts: make(ids),
		owners:    make(ids),
	}
	if req == nil {
		return balance, nil
	}
	balance.isEmpty = false
	balance.TokenBalanceFilter = req

	if err := fillAddressMapFromBytesFilter(ctx, address, req.Contract, balance.contracts); err != nil {
		return balance, err
	}
	if err := fillAddressMapFromBytesFilter(ctx, address, req.Owner, balance.owners); err != nil {
		return balance, err
	}

	return balance, nil
}

// Filter -
func (f TokenBalance) Filter(data storage.TokenBalance) bool {
	if f.isEmpty {
		return true
	}
	if f.TokenBalanceFilter == nil {
		return false
	}

	if f.Contract != nil {
		if !f.contracts.In(data.ContractID) {
			return false
		}
	}
	if f.Owner != nil {
		if !f.owners.In(data.OwnerID) {
			return false
		}
	}

	if !validString(f.TokenId, data.TokenID.String()) {
		return false
	}

	return true
}
