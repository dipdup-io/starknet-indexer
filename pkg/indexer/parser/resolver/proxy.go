package resolver

import (
	"bytes"
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Proxy -
func (resolver *Resolver) Proxy(ctx context.Context, txCtx data.TxContext, class storage.Class, contract storage.Address, selector []byte) (a abi.Abi, err error) {
	var (
		current = class
		hash    = contract.Hash
	)
	for current.Type.Is(storage.ClassTypeProxy) {
		c, address, err := resolver.resolveHash(ctx, txCtx, hash, selector)
		if err != nil {
			return a, err
		}
		current = *c
		if bytes.Equal(address, hash) {
			break
		}
		hash = address
	}

	if len(current.Abi) > 0 {
		return resolver.cache.GetAbiByClass(current)
	}
	return a, errors.Errorf("can't find contract abi under proxy: %x contract=%x", current.Hash, contract.Hash)
}

// UpgradeProxy -
func (resolver *Resolver) UpgradeProxy(ctx context.Context, contract storage.Address, upgrades []data.ProxyUpgrade, height uint64) error {
	for i := range upgrades {
		id, typ, err := resolver.findProxyEntity(ctx, upgrades[i].Address, height)
		if err != nil {
			return err
		}

		proxy := storage.Proxy{
			ContractID: contract.ID,
			Hash:       contract.Hash,
			Selector:   upgrades[i].Selector,
			EntityID:   id,
			EntityHash: upgrades[i].Address,
			EntityType: typ,
		}

		withAction := data.NewProxyWithAction(proxy, upgrades[i].Action)
		contextProxies := resolver.blockContext.CurrentProxies()

		key := data.NewProxyKey(contract.Hash, upgrades[i].Selector)
		contextProxies.Add(key, withAction)

		if upgrades[i].IsModule {
			endBlockProxies := resolver.blockContext.Proxies()
			endBlockProxies.Add(key, withAction)
		}
	}
	return nil
}

func (resolver *Resolver) resolveHash(ctx context.Context, txCtx data.TxContext, address, selector []byte) (*storage.Class, []byte, error) {
	proxy, err := resolver.findProxy(ctx, txCtx, address, selector)
	if err != nil {
		return nil, nil, err
	}

	switch proxy.EntityType {
	case storage.EntityTypeClass:
		class := &storage.Class{
			ID:     proxy.EntityID,
			Hash:   proxy.EntityHash,
			Height: txCtx.Height,
		}
		err := resolver.FindClass(ctx, class)
		return class, proxy.EntityHash, err
	case storage.EntityTypeContract:
		address := &storage.Address{
			ID:     proxy.EntityID,
			Hash:   proxy.EntityHash,
			Height: txCtx.Height,
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

func (resolver *Resolver) findProxy(ctx context.Context, txCtx data.TxContext, address, selector []byte) (storage.Proxy, error) {

	log.Info().
		Bytes("address", address).
		Bytes("selector", selector).
		Msg("find proxy")

	key := data.NewProxyKey(address, selector)
	if _, ok := txCtx.ProxyUpgrades.Get(key); !ok {
		log.Info().
			Bytes("address", address).
			Bytes("selector", selector).
			Msg("not found proxy")
		contextProxies := resolver.blockContext.CurrentProxies()
		if proxy, ok := contextProxies.Get(key); ok {
			log.Info().
				Bytes("address", address).
				Bytes("selector", selector).
				Msg("found in block context")
			return proxy.Proxy, nil
		}
	} else {
		log.Info().
			Bytes("address", address).
			Bytes("selector", selector).
			Msg("found proxy upgrade")
	}

	proxy, err := resolver.cache.GetProxy(ctx, address, selector)
	switch {
	case err == nil:
		log.Info().
			Bytes("address", address).
			Bytes("selector", selector).
			Msg("found in database")
		return proxy, err
	case resolver.blocks.IsNoRows(err):
		endBlockProxies := resolver.blockContext.Proxies()
		if proxy, ok := endBlockProxies.Get(key); ok {
			log.Info().
				Bytes("address", address).
				Bytes("selector", selector).
				Msg("found in end block context")
			return proxy.Proxy, nil
		}
		return storage.Proxy{}, errors.Wrapf(ErrUnknownProxy, "%x", address)
	default:
		return proxy, err
	}
}

func (resolver *Resolver) findProxyEntity(ctx context.Context, hash []byte, height uint64) (uint64, storage.EntityType, error) {
	sHash := encoding.EncodeHex(hash)
	if value, ok := resolver.blockContext.Addresses()[sHash]; ok {
		return value.ID, storage.EntityTypeContract, nil
	}
	if value, ok := resolver.blockContext.Classes()[sHash]; ok {
		return value.ID, storage.EntityTypeClass, nil
	}

	class, err := resolver.cache.GetClassByHash(ctx, hash)
	switch {
	case err == nil:
		return class.ID, storage.EntityTypeClass, nil

	case resolver.blocks.IsNoRows(err):
		address, err := resolver.cache.GetAddress(ctx, hash)
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
		return address.ID, storage.EntityTypeContract, nil

	default:
		return 0, storage.EntityTypeClass, errors.Wrapf(err, "get address: %x", hash)
	}
}
