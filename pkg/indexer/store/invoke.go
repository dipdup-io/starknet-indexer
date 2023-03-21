package store

import (
	"context"

	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (store *Store) saveInvoke(
	ctx context.Context,
	tx storage.Transaction,
	result parserData.Result,
) error {
	if result.Block.InvokeCount == 0 {
		return nil
	}

	models := make([]any, result.Block.InvokeCount)
	for i := range result.Block.Invoke {
		models[i] = &result.Block.Invoke[i]

		allInternals, err := store.saveInternals(ctx, tx, result, result.Block.Invoke[i].Internals)
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

		if err := store.saveEvents(ctx, tx, result.Block.Invoke[i].Events); err != nil {
			return err
		}

		if err := store.saveMessages(ctx, tx, result.Block.Invoke[i].Messages); err != nil {
			return err
		}

		if err := store.saveFee(ctx, tx, result.Block.Invoke[i].Fee); err != nil {
			return err
		}

		if err := store.saveTransfers(ctx, tx, result.Block.Invoke[i].Transfers); err != nil {
			return err
		}
	}

	return tx.BulkSave(ctx, models)
}
