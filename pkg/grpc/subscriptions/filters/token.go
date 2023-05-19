package filters

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Message -
type Token struct {
	*pb.TokenFilter
	isEmpty bool

	contracts ids
	owners    ids
}

// NewToken -
func NewToken(ctx context.Context, address storage.IAddress, req *pb.TokenFilter) (Token, error) {
	msg := Token{
		isEmpty:   true,
		contracts: make(ids),
		owners:    make(ids),
	}
	if req == nil {
		return msg, nil
	}
	msg.isEmpty = false
	msg.TokenFilter = req
	if err := fillAddressMapFromBytesFilter(ctx, address, req.Contract, msg.contracts); err != nil {
		return msg, err
	}
	if err := fillAddressMapFromBytesFilter(ctx, address, req.Owner, msg.owners); err != nil {
		return msg, err
	}
	return msg, nil
}

// Filter -
func (f Token) Filter(data storage.Token) bool {
	if f.isEmpty {
		return true
	}
	if f.TokenFilter == nil {
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

	if !validEnum(f.Type, uint64(data.Type)) {
		return false
	}

	return true
}
