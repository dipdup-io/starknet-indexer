package filters

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// TokenBalance -
type TokenBalance struct {
	fltrs []*pb.TokenBalanceFilter

	contracts []ids
	owners    []ids

	isEmpty bool
}

// NewTokenBalance -
func NewTokenBalance(ctx context.Context, address storage.IAddress, req []*pb.TokenBalanceFilter) (TokenBalance, error) {
	balance := TokenBalance{
		isEmpty: true,
	}
	if req == nil {
		return balance, nil
	}
	balance.isEmpty = false
	balance.fltrs = req
	balance.contracts = make([]ids, 0)
	balance.owners = make([]ids, 0)

	for i := range balance.fltrs {
		balance.contracts = append(balance.contracts, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, balance.fltrs[i].Contract, balance.contracts[i]); err != nil {
			return balance, err
		}
		balance.owners = append(balance.owners, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, balance.fltrs[i].Owner, balance.owners[i]); err != nil {
			return balance, err
		}
	}

	return balance, nil
}

// Filter -
func (f TokenBalance) Filter(data storage.TokenBalance) bool {
	if f.isEmpty {
		return true
	}

	var ok bool
	for i := range f.fltrs {
		if f.fltrs[i].Contract != nil {
			if !f.contracts[i].In(data.ContractID) {
				continue
			}
		}
		if f.fltrs[i].Owner != nil {
			if !f.owners[i].In(data.OwnerID) {
				continue
			}
		}

		if !validString(f.fltrs[i].TokenId, data.TokenID.String()) {
			continue
		}

		ok = true
		break
	}

	return ok
}
