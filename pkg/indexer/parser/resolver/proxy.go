package resolver

import (
	"context"

	jsoniter "github.com/json-iterator/go"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Proxy -
func (resolver *Resolver) Proxy(ctx context.Context, class storage.Class, contract storage.Address) (a abi.Abi, err error) {
	var (
		current = class
		hash    = contract.Hash
	)
	for current.Type.Is(storage.ClassTypeProxy) {
		c, address, err := resolver.resolveHash(ctx, hash)
		if err != nil {
			return a, err
		}
		current = *c
		hash = address
	}

	if len(current.Abi) > 0 {
		err = json.Unmarshal(current.Abi, &a)
		return a, err
	}
	return a, errors.Errorf("can't find contract abi under proxy: %x contract=%x", class.Hash, contract.Hash)
}

// UpgradeProxy -
func (resolver *Resolver) UpgradeProxy(ctx context.Context, contract storage.Address, impl []byte, height uint64) error {
	id, typ, err := resolver.findProxyEntity(ctx, impl, height)
	if err != nil {
		return err
	}

	proxy := storage.Proxy{
		ContractID: contract.ID,
		Hash:       contract.Hash,
		EntityID:   id,
		EntityHash: impl,
		EntityType: typ,
	}
	resolver.contextProxies[encoding.EncodeHex(contract.Hash)] = &proxy
	resolver.cache.SetProxy(proxy)
	return nil
}

func (resolver *Resolver) resolveHash(ctx context.Context, hash []byte) (*storage.Class, []byte, error) {
	proxy, err := resolver.findProxy(ctx, hash)
	if err != nil {
		return nil, nil, err
	}

	switch proxy.EntityType {
	case storage.EntityTypeClass:
		class := &storage.Class{
			ID:   proxy.EntityID,
			Hash: proxy.EntityHash,
		}
		err := resolver.FindClass(ctx, class)
		return class, proxy.EntityHash, err
	case storage.EntityTypeContract:
		address := &storage.Address{
			ID:   proxy.EntityID,
			Hash: proxy.EntityHash,
		}
		if err := resolver.FindAddress(ctx, address); err != nil {
			return nil, nil, err
		}
		if address.ClassID != nil {
			c, err := resolver.FindClassByID(ctx, *address.ClassID)
			return c, address.Hash, err
		}
		return nil, nil, errors.Errorf("unknown class id for contract: %x", address.Hash)
	default:
		return nil, nil, errors.Errorf("unknown proxy entity type: %d", proxy.EntityType)
	}
}

func (resolver *Resolver) findProxy(ctx context.Context, hash []byte) (storage.Proxy, error) {
	sHash := encoding.EncodeHex(hash)
	if proxy, ok := resolver.contextProxies[sHash]; ok {
		return *proxy, nil
	}

	proxy, err := resolver.cache.GetProxy(ctx, hash)
	switch {
	case err == nil:
		return proxy, err
	case resolver.blocks.IsNoRows(err):
		if proxy, ok := resolver.endBlockProxies[sHash]; ok {
			return *proxy, nil
		}
		return storage.Proxy{}, errors.Wrapf(ErrUnknownProxy, "%x", hash)
	default:
		return proxy, err
	}
}

func (resolver *Resolver) findProxyEntity(ctx context.Context, hash []byte, height uint64) (uint64, storage.EntityType, error) {
	sHash := encoding.EncodeHex(hash)
	if value, ok := resolver.addresses[sHash]; ok {
		return value.ID, storage.EntityTypeContract, nil
	}
	if value, ok := resolver.classes[sHash]; ok {
		return value.ID, storage.EntityTypeClass, nil
	}

	address, err := resolver.cache.GetAddress(ctx, hash)
	switch {
	case err == nil:
		return address.ID, storage.EntityTypeContract, nil

	case resolver.blocks.IsNoRows(err):
		class, err := resolver.cache.GetClassByHash(ctx, hash)
		if err != nil {
			if resolver.blocks.IsNoRows(err) {
				class, err := resolver.parseClass(ctx, hash, height)
				if err != nil {
					return 0, storage.EntityTypeClass, ErrUnknownProxy
				}
				return class.ID, storage.EntityTypeClass, nil
			}
			return 0, storage.EntityTypeClass, err
		}
		return class.ID, storage.EntityTypeClass, nil

	default:
		return 0, storage.EntityTypeClass, errors.Wrapf(err, "get address: %x", hash)
	}
}
