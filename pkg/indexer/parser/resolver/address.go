package resolver

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
)

// FindAddress -
func (resolver *Resolver) FindAddress(ctx context.Context, address *storage.Address) error {
	if value, ok := resolver.addresses[encoding.EncodeHex(address.Hash)]; ok {
		address.ID = value.ID
		address.ClassID = value.ClassID
		return nil
	}
	generated, err := resolver.idGenerator.SetAddressId(ctx, address)
	if err != nil {
		return err
	}
	if generated {
		resolver.addAddress(address)
	}
	return nil
}

// FindAddressByHash -
func (resolver *Resolver) FindAddressByHash(ctx context.Context, hash data.Felt) (*storage.Address, error) {
	if hash == "" {
		return nil, nil
	}

	address := &storage.Address{
		Hash: hash.Bytes(),
	}
	err := resolver.FindAddress(ctx, address)
	return address, err
}

func (resolver *Resolver) addAddress(address *storage.Address) {
	if len(address.Hash) == 0 {
		return
	}
	key := encoding.EncodeHex(address.Hash)
	if _, ok := resolver.addresses[key]; !ok {
		resolver.addresses[key] = address
	}
}
