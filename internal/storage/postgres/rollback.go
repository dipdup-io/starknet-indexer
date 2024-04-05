package postgres

import (
	"context"
	"fmt"
	"strings"

	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

const rollbackQuery = `DELETE FROM ? WHERE height > ? RETURNING id;`

// RollbackManager -
type RollbackManager struct {
	state         models.IState
	blocks        models.IBlock
	proxyUpgrades models.IProxyUpgrade
	transfers     models.ITransfer
	classReplaces models.IClassReplace
	transactable  storage.Transactable
}

// NewRollbackManager -
func NewRollbackManager(
	transactable storage.Transactable,
	state models.IState,
	blocks models.IBlock,
	proxyUpgrades models.IProxyUpgrade,
	classReplaces models.IClassReplace,
	transfers models.ITransfer,
) RollbackManager {
	return RollbackManager{
		transactable:  transactable,
		state:         state,
		blocks:        blocks,
		proxyUpgrades: proxyUpgrades,
		transfers:     transfers,
		classReplaces: classReplaces,
	}
}

// Rollback - rollback database to height
func (rm RollbackManager) Rollback(ctx context.Context, indexerName string, height uint64) error {
	log.Info().Uint64("new_height", height).Str("indexer", indexerName).Msg("rollback starting...")
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

	if err := rm.rollbackReplaceClass(ctx, height, tx); err != nil {
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
		&models.ProxyUpgrade{},
		&models.ClassReplace{},
	} {
		deletedCount, err := tx.Exec(ctx, rollbackQuery, bun.Ident(model.TableName()), height)
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
	log.Info().Uint64("new_height", height).Str("indexer", indexerName).Msg("rollback is finished")
	return nil
}

func (rm RollbackManager) rollbackProxy(ctx context.Context, height uint64, tx storage.Transaction) error {
	var (
		offset = 0
		limit  = 100
		end    = false
	)

	for !end {
		upgrades, err := rm.proxyUpgrades.ListWithHeight(ctx, height, limit, offset)
		if err != nil {
			return err
		}

		offset += len(upgrades)
		end = len(upgrades) != limit

		for i := range upgrades {
			switch upgrades[i].Action {
			case models.ProxyActionAdd:
				if _, err = tx.Exec(ctx, `delete from proxy where selector = ? and hash = ?`, upgrades[i].Selector, upgrades[i].Hash); err != nil {
					return err
				}
			case models.ProxyActionUpdate:
				last, err := rm.proxyUpgrades.LastBefore(ctx, upgrades[i].Hash, upgrades[i].Selector, upgrades[i].Height)
				if err != nil {
					if rm.proxyUpgrades.IsNoRows(err) {
						if _, err = tx.Exec(ctx, `delete from proxy where selector = ? and hash = ?`, upgrades[i].Selector, upgrades[i].Hash); err != nil {
							return err
						}
						continue
					}
					return err
				}

				if upgrades[i].Selector == nil {
					if _, err := tx.Exec(ctx,
						`update proxy set entity_type = ?, entity_id = ?, entity_hash = ? where selector is NULL and hash = ?`,
						last.EntityType, last.EntityID, last.EntityHash, last.Hash,
					); err != nil {
						return err
					}
				} else {
					if _, err := tx.Exec(ctx,
						`update proxy set entity_type = ?, entity_id = ?, entity_hash = ? where selector = ? and hash = ?`,
						last.EntityType, last.EntityID, last.EntityHash, last.Selector, last.Hash,
					); err != nil {
						return err
					}
				}
			case models.ProxyActionDelete:
				if _, err = tx.Exec(ctx,
					`insert into proxy (contract_id, hash, selector, entity_type, entity_id, entity_hash) values (?,?,?,?,?,?)`,
					upgrades[i].ContractID, upgrades[i].Hash, upgrades[i].Selector, upgrades[i].EntityType, upgrades[i].EntityID, upgrades[i].EntityHash,
				); err != nil {
					return err
				}
			}
		}
	}
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
			models.WithLimitFilter(limit),
			models.WithOffsetFilter(offset),
			models.WithMultiSort("time desc", "id desc"),
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
					OwnerID:    transfers[i].ToID,
					ContractID: transfers[i].ContractID,
					TokenID:    transfers[i].TokenID,
					Balance:    transfers[i].Amount.Copy().Neg(),
				}
			}
		}
	}
	if len(updates) == 0 {
		return nil
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
		bun.Safe(strings.Join(values, ",")),
	)
	return err
}

func (rm RollbackManager) rollbackReplaceClass(ctx context.Context, height uint64, tx storage.Transaction) error {
	replaces, err := rm.classReplaces.ByHeight(ctx, height)
	if err != nil {
		return err
	}
	if len(replaces) == 0 {
		return nil
	}

	for i := range replaces {
		_, err = tx.Tx().NewUpdate().Model((*models.Address)(nil)).
			Where("id = ?", replaces[i].ContractId).
			Set("class_id = ?", replaces[i].PrevClassId).
			Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
