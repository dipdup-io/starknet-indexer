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

	models := make([]any, result.Block.DeployCount)
	for i := range result.Block.Deploy {
		models[i] = &result.Block.Deploy[i]

		sm.addInternals(result.Block.Deploy[i].Internals)
		sm.addEvents(result.Block.Deploy[i].Events)
		sm.addMessages(result.Block.Deploy[i].Messages)
		sm.addTransfers(result.Block.Deploy[i].Transfers)

		if err := store.saveInternals(ctx, tx, result.Block.Deploy[i].Internals, sm); err != nil {
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
