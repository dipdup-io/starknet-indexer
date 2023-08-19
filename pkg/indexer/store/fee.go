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

	for i := range result.Block.Fee {
		sm.addInternals(result.Block.Fee[i].Internals)
		sm.addEvents(result.Block.Fee[i].Events)
		sm.addMessages(result.Block.Fee[i].Messages)
		sm.addTransfers(result.Block.Fee[i].Transfers)

		if err := store.saveInternals(ctx, tx, result.Block.Fee[i].Internals, sm); err != nil {
			return err
		}
	}

	return bulkSaveWithCopy(ctx, tx, result.Block.Fee)
}
