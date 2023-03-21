package store

import (
	"context"

	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (store *Store) saveInternals(
	ctx context.Context,
	tx storage.Transaction,
	result parserData.Result,
	internals []models.Internal,
) ([]models.Internal, error) {
	if len(internals) == 0 {
		return nil, nil
	}

	models := make([]models.Internal, 0)
	for i := range internals {
		models = append(models, internals[i])
		if err := store.saveEvents(ctx, tx, internals[i].Events); err != nil {
			return nil, err
		}

		if err := store.saveMessages(ctx, tx, internals[i].Messages); err != nil {
			return nil, err
		}

		if err := store.saveTransfers(ctx, tx, internals[i].Transfers); err != nil {
			return nil, err
		}

		received, err := store.saveInternals(ctx, tx, result, internals[i].Internals)
		if err != nil {
			return nil, err
		}
		models = append(models, received...)
	}

	return models, nil
}

func (store *Store) saveEvents(ctx context.Context, tx storage.Transaction, events []models.Event) error {
	models := make([]any, len(events))
	for i := range events {
		models[i] = &events[i]
	}

	return tx.BulkSave(ctx, models)
}

func (store *Store) saveMessages(ctx context.Context, tx storage.Transaction, msgs []models.Message) error {
	models := make([]any, len(msgs))
	for i := range msgs {
		models[i] = &msgs[i]
	}
	return tx.BulkSave(ctx, models)
}

func (store *Store) saveTransfers(ctx context.Context, tx storage.Transaction, transfers []models.Transfer) error {
	if len(transfers) == 0 {
		return nil
	}
	result := make([]any, 0)
	for j := range transfers {
		result = append(result, &transfers[j])
		if err := store.saveBalanceUpdates(ctx, tx, transfers[j].TokenBalanceUpdates()); err != nil {
			return err
		}
	}
	return tx.BulkSave(ctx, result)
}

func (store *Store) saveBalanceUpdates(ctx context.Context, tx storage.Transaction, updates []models.TokenBalance) error {
	for i := range updates {
		if _, err := tx.Exec(ctx, `
		INSERT INTO token_balance (owner_id, contract_id, token_id, balance)
		VALUES(?,?,?,?) 
		ON CONFLICT (owner_id, contract_id, token_id)
		DO 
		UPDATE SET balance = token_balance.balance + excluded.balance
		`, updates[i].OwnerID, updates[i].ContractID, updates[i].TokenID, updates[i].Balance); err != nil {
			return err
		}
	}
	return nil
}

func (store *Store) saveFee(ctx context.Context, tx storage.Transaction, fee *models.Fee) error {
	if fee == nil {
		return nil
	}

	if err := store.saveBalanceUpdates(ctx, tx, fee.TokenBalanceUpdates()); err != nil {
		return err
	}

	return tx.Add(ctx, fee)
}
