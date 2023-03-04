package parser

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
)

func (parser *Parser) entitiesFromStateUpdate(ctx context.Context, block *storage.Block, upd data.StateUpdate) error {
	if err := parser.parseDeclaredContracts(ctx, upd.StateDiff.DeclaredContracts); err != nil {
		return err
	}

	if err := parser.parseDeployedContracts(ctx, upd.StateDiff.DeployedContracts); err != nil {
		return err
	}

	if err := parser.parseStorageDiffs(ctx, block, upd.StateDiff.StorageDiffs); err != nil {
		return err
	}

	return nil
}

func (parser *Parser) parseDeclaredContracts(ctx context.Context, declared []string) error {
	for i := range declared {
		hash := encoding.MustDecodeHex(declared[i])
		class := storage.Class{
			Hash: hash,
			ID:   parser.idGenerator.NextClassId(),
		}
		if err := parser.receiveClass(ctx, &class); err != nil {
			return err
		}
	}
	return nil
}

func (parser *Parser) parseDeployedContracts(ctx context.Context, contracts []data.DeployedContract) error {
	for i := range contracts {
		classHash := encoding.MustDecodeHex(contracts[i].ClassHash)
		class := storage.Class{
			Hash: classHash,
		}
		if err := parser.findClass(ctx, &class); err != nil {
			return err
		}
		if class.Abi == nil {
			if err := parser.receiveClass(ctx, &class); err != nil {
				return err
			}
		}

		addressHash := encoding.MustDecodeHex(contracts[i].Address)
		address := storage.Address{
			Hash:    addressHash,
			ID:      parser.idGenerator.NextAddressId(),
			ClassID: &class.ID,
		}
		parser.addAddress(&address)
	}

	return nil
}

func (parser *Parser) parseStorageDiffs(ctx context.Context, block *storage.Block, diffs map[string][]data.KeyValue) error {
	block.StorageDiffs = make([]storage.StorageDiff, 0)
	for hash, updates := range diffs {
		address := storage.Address{
			Hash: encoding.MustDecodeHex(hash),
		}
		if err := parser.findAddress(ctx, &address); err != nil {
			return err
		}

		for i := range updates {
			block.StorageDiffs = append(block.StorageDiffs, storage.StorageDiff{
				ContractID: address.ID,
				Height:     block.Height,
				Key:        encoding.MustDecodeHex(updates[i].Key),
				Value:      encoding.MustDecodeHex(updates[i].Value),
			})
		}
	}
	block.StorageDiffCount = len(block.StorageDiffs)
	return nil
}
