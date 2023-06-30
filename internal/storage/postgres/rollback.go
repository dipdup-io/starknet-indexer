package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/dipdup-io/starknet-indexer/internal/starknet"
	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/go-pg/pg/v10"
)

const rollbackQuery = `WITH deleted AS (DELETE FROM ? WHERE height > ? RETURNING *) SELECT count(*) FROM deleted;`

// RollbackManager -
type RollbackManager struct {
	state        models.IState
	blocks       models.IBlock
	storageDiffs models.IStorageDiff
	transfers    models.ITransfer
	transactable storage.Transactable
}

// NewRollbackManager -
func NewRollbackManager(
	transactable storage.Transactable,
	state models.IState,
	blocks models.IBlock,
	storageDiffs models.IStorageDiff,
	transfers models.ITransfer,
) RollbackManager {
	return RollbackManager{
		transactable: transactable,
		state:        state,
		blocks:       blocks,
		storageDiffs: storageDiffs,
		transfers:    transfers,
	}
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

	if err := rm.rollbackProxy(ctx, height, tx); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := rm.rollbackTokenBalances(ctx, height, tx); err != nil {
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
			state.TxCount -= uint64(deletedCount)
		case *models.Declare:
			state.DeclaresCount -= uint64(deletedCount)
			state.TxCount -= uint64(deletedCount)
		case *models.Deploy:
			state.DeployAccountsCount -= uint64(deletedCount)
			state.TxCount -= uint64(deletedCount)
		case *models.DeployAccount:
			state.DeployAccountsCount -= uint64(deletedCount)
			state.TxCount -= uint64(deletedCount)
		case *models.L1Handler:
			state.L1HandlersCount -= uint64(deletedCount)
			state.TxCount -= uint64(deletedCount)
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

func (rm RollbackManager) rollbackProxy(ctx context.Context, height uint64, tx storage.Transaction) error {
	keys := make([][]byte, 0)
	for _, value := range starknet.ProxyStorageVars {
		keys = append(keys, value)
	}
	_, err := rm.storageDiffs.GetForKeys(ctx, keys, height)
	if err != nil {
		return err
	}
	// TODO: rollback proxy by diffs
	return nil
}

func (rm RollbackManager) rollbackTokenBalances(ctx context.Context, height uint64, tx storage.Transaction) error {
	var (
		offset = 0
		limit  = 100
		end    = false
	)

	updates := make(map[string]*models.TokenBalance, 0)
	for !end {
		transfers, err := rm.transfers.Filter(ctx,
			[]models.TransferFilter{
				{
					Height: models.IntegerFilter{
						Gt: height,
					},
				},
			},
			models.WithDescSortByIdFilter(),
			models.WithLimitFilter(limit),
			models.WithOffsetFilter(offset),
		)
		if err != nil {
			return err
		}

		offset += len(transfers)
		end = len(transfers) != limit

		for i := range transfers {
			fromKey := fmt.Sprintf("%d_%d_%s", transfers[i].ContractID, transfers[i].FromID, transfers[i].TokenID.String())
			if upd, ok := updates[fromKey]; ok {
				upd.Balance = upd.Balance.Add(transfers[i].Amount)
			} else {
				updates[fromKey] = &models.TokenBalance{
					OwnerID:    transfers[i].FromID,
					ContractID: transfers[i].ContractID,
					TokenID:    transfers[i].TokenID,
					Balance:    transfers[i].Amount.Copy(),
				}
			}
			toKey := fmt.Sprintf("%d_%d_%s", transfers[i].ContractID, transfers[i].ToID, transfers[i].TokenID.String())
			if upd, ok := updates[toKey]; ok {
				upd.Balance = upd.Balance.Sub(transfers[i].Amount)
			} else {
				updates[toKey] = &models.TokenBalance{
					OwnerID:    transfers[i].FromID,
					ContractID: transfers[i].ContractID,
					TokenID:    transfers[i].TokenID,
					Balance:    transfers[i].Amount.Copy().Neg(),
				}
			}
		}
	}

	values := make([]string, 0)
	for _, update := range updates {
		value := fmt.Sprintf(
			"(%d,%d,%s,%s)",
			update.OwnerID,
			update.ContractID,
			update.TokenID,
			update.Balance,
		)
		values = append(values, value)
	}
	_, err := tx.Exec(ctx,
		`INSERT INTO token_balance (owner_id, contract_id, token_id, balance)
		VALUES ? 
		ON CONFLICT (owner_id, contract_id, token_id)
		DO 
		UPDATE SET balance = token_balance.balance + excluded.balance`,
		pg.Safe(strings.Join(values, ",")),
	)
	return err
}
