package store

import (
	"context"

	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (store *Store) saveDeploy(
	ctx context.Context,
	tx storage.Transaction,
	result parserData.Result,
) error {
	if result.Block.DeployCount == 0 {
		return nil
	}

	models := make([]any, result.Block.DeployCount)
	for i := range result.Block.Deploy {
		models[i] = &result.Block.Deploy[i]

		allInternals, err := store.saveInternals(ctx, tx, result, result.Block.Deploy[i].Internals)
		if err != nil {
			return err
		}
		internalModels := make([]any, len(allInternals))
		for i := range allInternals {
			internalModels[i] = &allInternals[i]
		}
		if len(allInternals) > 0 {
			if err := tx.BulkSave(ctx, internalModels); err != nil {
				return err
			}
		}

		if err := store.saveEvents(ctx, tx, result.Block.Deploy[i].Events); err != nil {
			return err
		}

		if err := store.saveMessages(ctx, tx, result.Block.Deploy[i].Messages); err != nil {
			return err
		}

		if err := store.saveTransfers(ctx, tx, result.Block.Deploy[i].Transfers); err != nil {
			return err
		}

		if err := store.saveFee(ctx, tx, result.Block.Deploy[i].Fee); err != nil {
			return err
		}

		switch {
		case result.Block.Deploy[i].ERC20 != nil:
			err = tx.Add(ctx, result.Block.Deploy[i].ERC20)
		case result.Block.Deploy[i].ERC721 != nil:
			err = tx.Add(ctx, result.Block.Deploy[i].ERC721)
		case result.Block.Deploy[i].ERC1155 != nil:
			err = tx.Add(ctx, result.Block.Deploy[i].ERC1155)
		}
		if err != nil {
			return err
		}
	}

	return tx.BulkSave(ctx, models)
}
