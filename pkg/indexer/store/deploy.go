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
	sm *subModels,
) error {
	if result.Block.DeployCount == 0 {
		return nil
	}

	var err error

	for i := range result.Block.Deploy {

		sm.addInternals(result.Block.Deploy[i].Internals)
		sm.addEvents(result.Block.Deploy[i].Events)
		sm.addMessages(result.Block.Deploy[i].Messages)
		sm.addTransfers(result.Block.Deploy[i].Transfers)

		if err := store.saveInternals(ctx, tx, result.Block.Deploy[i].Internals, sm); err != nil {
			return err
		}

		if result.Block.Deploy[i].Token != nil {
			if err = tx.Add(ctx, result.Block.Deploy[i].Token); err != nil {
				return err
			}
		}
	}

	return bulkSaveWithCopy(ctx, tx, result.Block.Deploy)
}
