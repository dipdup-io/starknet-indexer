package store

import (
	"context"

	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (store *Store) saveDeclare(
	ctx context.Context,
	tx storage.Transaction,
	result parserData.Result,
	sm *subModels,
) error {
	if result.Block.DeclareCount == 0 {
		return nil
	}

	models := make([]any, result.Block.DeclareCount)
	for i := range result.Block.Declare {
		models[i] = &result.Block.Declare[i]

		sm.addInternals(result.Block.Declare[i].Internals)
		sm.addEvents(result.Block.Declare[i].Events)
		sm.addMessages(result.Block.Declare[i].Messages)
		sm.addTransfers(result.Block.Declare[i].Transfers)

		if err := store.saveInternals(ctx, tx, result.Block.Declare[i].Internals, sm); err != nil {
			return err
		}
	}

	return tx.BulkSave(ctx, models)
}
