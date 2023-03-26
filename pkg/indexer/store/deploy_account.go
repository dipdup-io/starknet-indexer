package store

import (
	"context"

	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (store *Store) saveDeployAccount(
	ctx context.Context,
	tx storage.Transaction,
	result parserData.Result,
	sm *subModels,
) error {
	if result.Block.DeployAccountCount == 0 {
		return nil
	}

	models := make([]any, result.Block.DeployAccountCount)
	for i := range result.Block.DeployAccount {
		models[i] = &result.Block.DeployAccount[i]

		sm.addInternals(result.Block.DeployAccount[i].Internals)
		sm.addEvents(result.Block.DeployAccount[i].Events)
		sm.addMessages(result.Block.DeployAccount[i].Messages)
		sm.addTransfers(result.Block.DeployAccount[i].Transfers)

		if err := store.saveInternals(ctx, tx, result.Block.DeployAccount[i].Internals, sm); err != nil {
			return err
		}
	}

	return tx.BulkSave(ctx, models)
}
