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
	sm *subModels,
) error {
	if result.Block.InvokeCount == 0 {
		return nil
	}

	for i := range result.Block.Invoke {
		sm.addInternals(result.Block.Invoke[i].Internals)
		sm.addEvents(result.Block.Invoke[i].Events)
		sm.addMessages(result.Block.Invoke[i].Messages)
		sm.addTransfers(result.Block.Invoke[i].Transfers)

		if err := store.saveInternals(ctx, tx, result.Block.Invoke[i].Internals, sm); err != nil {
			return err
		}
	}

	return bulkSaveWithCopy(ctx, tx, result.Block.Invoke)
}
