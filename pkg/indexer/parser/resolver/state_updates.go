package resolver

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-indexer/internal/starknet"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// ResolveStateUpdates -
func (resolver *Resolver) ResolveStateUpdates(ctx context.Context, block *storage.Block, upd data.StateUpdate) error {
	if err := resolver.parseDeclaredContracts(ctx, block, upd.StateDiff.OldDeclaredContracts); err != nil {
		return err
	}

	if err := resolver.parseDeclaredClasses(ctx, block, upd.StateDiff.DeclaredClasses); err != nil {
		return err
	}

	if err := resolver.parseDeployedContracts(ctx, block, upd.StateDiff.DeployedContracts); err != nil {
		return err
	}

	if err := resolver.parseStorageDiffs(ctx, block, upd.StateDiff.StorageDiffs); err != nil {
		return err
	}

	if err := resolver.parseReplaceClasses(ctx, block, upd.StateDiff.ReplacedClasses); err != nil {
		return err
	}

	return nil
}

func (resolver *Resolver) parseReplaceClasses(ctx context.Context, block *storage.Block, replaced []data.ReplacedClass) error {
	addrs := make(map[data.Felt]storage.Class)
	for i := range replaced {
		class, err := resolver.cache.GetClassByHash(ctx, replaced[i].ClassHash.Bytes())
		if err != nil {
			class, err = resolver.parseClassFromFelt(ctx, replaced[i].ClassHash, block.Height)
			if err != nil {
				return errors.Wrap(err, replaced[i].ClassHash.String())
			}
		}

		key := data.NewFeltFromBytes(replaced[i].Address.Bytes())
		addrs[key] = class
	}

	hash := make([][]byte, 0)
	for felt := range addrs {
		hash = append(hash, felt.Bytes())
	}

	addresses, err := resolver.addresses.GetByHashes(ctx, hash)
	if err != nil {
		return err
	}

	for i := range addresses {
		h := data.NewFeltFromBytes(addresses[i].Hash)
		if class, ok := addrs[h]; ok {
			addresses[i].ClassID = &class.ID
			resolver.addAddress(&addresses[i])
			resolver.cache.SetAbiByAddress(class, addresses[i].Hash)
			delete(addrs, h)
		}
	}

	for h, class := range addrs {
		id := class.ID
		address := storage.Address{
			Hash:    h.Bytes(),
			ID:      resolver.idGenerator.NextAddressId(),
			ClassID: &id,
			Height:  block.Height,
		}
		resolver.addAddress(&address)
		resolver.cache.SetAddress(ctx, address)
	}

	return nil
}

func (resolver *Resolver) parseDeclaredClasses(ctx context.Context, block *storage.Block, declared []data.DeclaredClass) error {
	for i := range declared {
		if _, err := resolver.parseClassFromFelt(ctx, declared[i].ClassHash, block.Height); err != nil {
			return errors.Wrap(err, declared[i].ClassHash.String())
		}
	}
	return nil
}

func (resolver *Resolver) parseDeclaredContracts(ctx context.Context, block *storage.Block, declared []data.Felt) error {
	for i := range declared {
		if _, err := resolver.parseClassFromFelt(ctx, declared[i], block.Height); err != nil {
			return errors.Wrap(err, declared[i].String())
		}
	}
	return nil
}

func (resolver *Resolver) parseDeployedContracts(ctx context.Context, block *storage.Block, contracts []data.DeployedContract) error {
	addrs := make(map[data.Felt]storage.Class)
	for i := range contracts {
		class, err := resolver.cache.GetClassByHash(ctx, contracts[i].ClassHash.Bytes())
		if err != nil {
			class, err = resolver.parseClassFromFelt(ctx, contracts[i].ClassHash, block.Height)
			if err != nil {
				return errors.Wrap(err, contracts[i].ClassHash.String())
			}
		}

		key := data.NewFeltFromBytes(contracts[i].Address.Bytes())
		addrs[key] = class
	}

	hash := make([][]byte, 0)
	for felt := range addrs {
		hash = append(hash, felt.Bytes())
	}

	addresses, err := resolver.addresses.GetByHashes(ctx, hash)
	if err != nil {
		return err
	}

	for i := range addresses {
		h := data.NewFeltFromBytes(addresses[i].Hash)
		if class, ok := addrs[h]; ok && addresses[i].ClassID == nil {
			addresses[i].ClassID = &class.ID
			resolver.addAddress(&addresses[i])
			resolver.cache.SetAbiByAddress(class, addresses[i].Hash)
			delete(addrs, h)
		}
	}

	for h, class := range addrs {
		id := class.ID
		address := storage.Address{
			Hash:    h.Bytes(),
			ID:      resolver.idGenerator.NextAddressId(),
			ClassID: &id,
			Height:  block.Height,
		}
		resolver.addAddress(&address)
		resolver.cache.SetAddress(ctx, address)
	}

	return nil
}

func (resolver *Resolver) parseStorageDiffs(ctx context.Context, block *storage.Block, diffs map[data.Felt][]data.KeyValue) error {
	endBlockProxies := resolver.blockContext.Proxies()
	block.StorageDiffs = make([]storage.StorageDiff, 0)
	for hash, updates := range diffs {
		address, err := resolver.FindAddressByHash(ctx, hash)
		if err != nil {
			return err
		}
		if address == nil {
			continue
		}

		for i := range updates {
			diff := storage.StorageDiff{
				ContractID: address.ID,
				Height:     block.Height,
				Key:        updates[i].Key.Bytes(),
				Value:      updates[i].Value.Bytes(),
			}
			block.StorageDiffs = append(block.StorageDiffs, diff)

			sKey := updates[i].Key.String()
			if _, ok := starknet.ProxyStorageVars[sKey]; ok {
				proxyUpgrade := storage.ProxyUpgrade{
					Hash:       address.Hash,
					ContractID: address.ID,
					EntityHash: diff.Value,
					Action:     storage.ProxyActionUpdate,
					Height:     block.Height,
				}
				id, typ, err := resolver.findProxyEntity(ctx, diff.Value, block.Height)
				if err != nil {
					if errors.Is(err, ErrUnknownProxy) {
						log.Warn().Err(err).Msgf("%x", diff.Value)
						continue
					}
					return errors.Wrap(err, "find proxy entity")
				}
				proxyUpgrade.EntityID = id
				proxyUpgrade.EntityType = typ
				endBlockProxies.AddByHash(address.Hash, nil, &proxyUpgrade)
			}
		}
	}
	block.StorageDiffCount = len(block.StorageDiffs)
	return nil
}
