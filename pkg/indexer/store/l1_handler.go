package store

import (
	"context"

	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (store *Store) saveL1Handler(
	ctx context.Context,
	tx storage.Transaction,
	result parserData.Result,
	sm *subModels,
) error {
	if result.Block.L1HandlerCount == 0 {
		return nil
	}

	models := make([]any, result.Block.L1HandlerCount)
	for i := range result.Block.L1Handler {
		models[i] = &result.Block.L1Handler[i]

		sm.addInternals(result.Block.L1Handler[i].Internals)
		sm.addEvents(result.Block.L1Handler[i].Events)
		sm.addMessages(result.Block.L1Handler[i].Messages)
		sm.addTransfers(result.Block.L1Handler[i].Transfers)

		if err := store.saveInternals(ctx, tx, result.Block.L1Handler[i].Internals, sm); err != nil {
			return err
		}
	}

	return tx.BulkSave(ctx, models)
}
