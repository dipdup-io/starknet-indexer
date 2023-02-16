package indexer

import (
	"context"

	models "github.com/dipdup-io/starknet-indexer/internal/storage"
)

func (indexer *Indexer) save(ctx context.Context, block models.Block) error {
	tx, err := indexer.transactable.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	if err := tx.Add(ctx, &block); err != nil {
		return tx.HandleError(ctx, err)
	}

	if block.TxCount == 0 {
		return nil
	}

	if block.DeclareCount > 0 {
		entities := make([]any, 0)
		for i := range block.Declare {
			ptr := &block.Declare[i]
			ptr.BlockID = block.ID
			entities = append(entities, ptr)
		}
		if err := tx.BulkSave(ctx, entities); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if block.DeployCount > 0 {
		entities := make([]any, 0)
		for i := range block.Deploy {
			ptr := &block.Deploy[i]
			ptr.BlockID = block.ID
			entities = append(entities, ptr)
		}
		if err := tx.BulkSave(ctx, entities); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if block.DeployAccountCount > 0 {
		entities := make([]any, 0)
		for i := range block.DeployAccount {
			ptr := &block.DeployAccount[i]
			ptr.BlockID = block.ID
			entities = append(entities, ptr)
		}
		if err := tx.BulkSave(ctx, entities); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if block.InvokeV0Count > 0 {
		entities := make([]any, 0)
		for i := range block.InvokeV0 {
			ptr := &block.InvokeV0[i]
			ptr.BlockID = block.ID
			entities = append(entities, ptr)
		}
		if err := tx.BulkSave(ctx, entities); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if block.InvokeV1Count > 0 {
		entities := make([]any, 0)
		for i := range block.InvokeV1 {
			ptr := &block.InvokeV1[i]
			ptr.BlockID = block.ID
			entities = append(entities, ptr)
		}
		if err := tx.BulkSave(ctx, entities); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if block.L1HandlerCount > 0 {
		entities := make([]any, 0)
		for i := range block.L1Handler {
			ptr := &block.L1Handler[i]
			ptr.BlockID = block.ID
			entities = append(entities, ptr)
		}
		if err := tx.BulkSave(ctx, entities); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	return tx.Flush(ctx)
}
