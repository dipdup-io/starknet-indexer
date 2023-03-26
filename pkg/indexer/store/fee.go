package store

import (
	"context"

	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (store *Store) saveFees(
	ctx context.Context,
	tx storage.Transaction,
	result parserData.Result,
	sm *subModels,
) error {
	if len(result.Block.Fee) == 0 {
		return nil
	}

	models := make([]any, len(result.Block.Fee))
	for i := range result.Block.Fee {
		models[i] = &result.Block.Fee[i]

		sm.addInternals(result.Block.Fee[i].Internals)
		sm.addEvents(result.Block.Fee[i].Events)
		sm.addMessages(result.Block.Fee[i].Messages)
		sm.addTransfers(result.Block.Fee[i].Transfers)

		if err := store.saveInternals(ctx, tx, result.Block.Fee[i].Internals, sm); err != nil {
			return err
		}
	}

	return tx.BulkSave(ctx, models)
}
