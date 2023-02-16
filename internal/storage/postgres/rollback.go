package postgres

import (
	"context"

	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/go-pg/pg/v10"
)

const rollbackQuery = `WITH deleted AS (DELETE FROM ? WHERE height > ? RETURNING *) SELECT count(*) FROM deleted;`

// RollbackManager -
type RollbackManager struct {
	state        models.IState
	blocks       models.IBlock
	transactable storage.Transactable
}

// NewRollbackManager -
func NewRollbackManager(transactable storage.Transactable, state models.IState, blocks models.IBlock) RollbackManager {
	return RollbackManager{transactable: transactable, state: state, blocks: blocks}
}

// Rollback - rollback database to height
func (rm RollbackManager) Rollback(ctx context.Context, indexerName string, height uint64) error {
	tx, err := rm.transactable.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	state, err := rm.state.ByName(ctx, indexerName)
	if err != nil {
		return tx.HandleError(ctx, err)
	}

	for _, model := range []storage.Model{
		&models.Address{},
		&models.Class{},
		&models.StorageDiff{},
		&models.Block{},
		&models.Invoke{},
		&models.Declare{},
		&models.Deploy{},
		&models.DeployAccount{},
		&models.L1Handler{},
		&models.Internal{},
		&models.Event{},
		&models.Message{},
		&models.Fee{},
		&models.Transfer{},
	} {
		deletedCount, err := tx.Exec(ctx, rollbackQuery, pg.Ident(model.TableName()), height)
		if err != nil {
			return tx.HandleError(ctx, err)
		}

		switch model.(type) {
		case *models.Invoke:
			state.InvokesCount -= uint64(deletedCount)
		case *models.Declare:
			state.DeclaresCount -= uint64(deletedCount)
		case *models.Deploy:
			state.DeployAccountsCount -= uint64(deletedCount)
		case *models.DeployAccount:
			state.DeployAccountsCount -= uint64(deletedCount)
		case *models.L1Handler:
			state.L1HandlersCount -= uint64(deletedCount)
		}
	}

	lastBlock, err := rm.blocks.ByHeight(ctx, height)
	if err != nil {
		return tx.HandleError(ctx, err)
	}
	state.LastTime = lastBlock.Time
	state.LastHeight = lastBlock.Height

	if err := tx.Update(ctx, &state); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Flush(ctx); err != nil {
		return tx.HandleError(ctx, err)
	}
	return nil
}
