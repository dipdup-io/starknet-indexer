package store

import (
	"context"

	models "github.com/dipdup-io/starknet-indexer/internal/storage"
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

	for i := range result.Block.L1Handler {
		sm.addInternals(result.Block.L1Handler[i].Internals)
		sm.addEvents(result.Block.L1Handler[i].Events)
		sm.addMessages(result.Block.L1Handler[i].Messages)
		sm.addTransfers(result.Block.L1Handler[i].Transfers)

		if err := store.saveInternals(ctx, tx, result.Block.L1Handler[i].Internals, sm); err != nil {
			return err
		}
	}

	return bulkSaveWithCopy[models.L1Handler](ctx, tx, store.l1Handlers, result.Block.L1Handler)
}
