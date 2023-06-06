package filters

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Message -
type Token struct {
	fltrs   []*pb.TokenFilter
	isEmpty bool

	contracts []ids
	owners    []ids
}

// NewToken -
func NewToken(ctx context.Context, address storage.IAddress, req []*pb.TokenFilter) (Token, error) {
	token := Token{
		isEmpty: true,
	}
	if req == nil {
		return token, nil
	}
	token.isEmpty = false
	token.fltrs = req
	token.contracts = make([]ids, 0)
	token.owners = make([]ids, 0)

	for i := range token.fltrs {
		token.contracts = append(token.contracts, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, token.fltrs[i].Contract, token.contracts[i]); err != nil {
			return token, err
		}
		token.owners = append(token.owners, make(ids))
		if err := fillAddressMapFromBytesFilter(ctx, address, token.fltrs[i].Owner, token.owners[i]); err != nil {
			return token, err
		}
	}
	return token, nil
}

// Filter -
func (f Token) Filter(data storage.Token) bool {
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

		if !validEnum(f.fltrs[i].Type, uint64(data.Type)) {
			continue
		}

		ok = true
		break
	}

	return ok
}
