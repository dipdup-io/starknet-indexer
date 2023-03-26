package store

import (
	"context"
	"fmt"
	"strings"

	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/go-pg/pg/v10"
)

type subModels struct {
	Transfers     []any
	TokenBalances []models.TokenBalance
	Events        []any
	Messages      []any
	Internals     []any
}

func newSubModels() *subModels {
	return &subModels{
		Transfers:     make([]any, 0),
		Events:        make([]any, 0),
		Messages:      make([]any, 0),
		Internals:     make([]any, 0),
		TokenBalances: make([]models.TokenBalance, 0),
	}
}

// Save -
func (sm *subModels) Save(ctx context.Context, tx storage.Transaction) error {
	if len(sm.Internals) > 0 {
		if err := tx.BulkSave(ctx, sm.Internals); err != nil {
			return err
		}
	}
	if len(sm.Events) > 0 {
		if err := tx.BulkSave(ctx, sm.Events); err != nil {
			return err
		}
	}
	if len(sm.Messages) > 0 {
		if err := tx.BulkSave(ctx, sm.Messages); err != nil {
			return err
		}
	}
	if len(sm.Transfers) > 0 {
		if err := tx.BulkSave(ctx, sm.Transfers); err != nil {
			return err
		}
	}
	if len(sm.TokenBalances) > 0 {
		if err := sm.saveTokenBalanceUpdates(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

func (sm *subModels) addEvents(events []models.Event) {
	for i := range events {
		sm.Events = append(sm.Events, &events[i])
	}
}

func (sm *subModels) addMessages(msgs []models.Message) {
	for i := range msgs {
		sm.Messages = append(sm.Messages, &msgs[i])
	}
}

func (sm *subModels) addTransfers(transfers []models.Transfer) {
	for i := range transfers {
		sm.Transfers = append(sm.Transfers, &transfers[i])
		sm.TokenBalances = append(sm.TokenBalances, transfers[i].TokenBalanceUpdates()...)
	}
}

func (sm *subModels) addInternals(internals []models.Internal) {
	for i := range internals {
		sm.Internals = append(sm.Internals, &internals[i])
	}
}

func (sm *subModels) saveTokenBalanceUpdates(ctx context.Context, tx storage.Transaction) error {
	updates := make(map[string]*models.TokenBalance)
	for i := range sm.TokenBalances {
		key := sm.TokenBalances[i].String()
		if update, ok := updates[key]; ok {
			update.Balance = update.Balance.Add(sm.TokenBalances[i].Balance)
		} else {
			updates[key] = &sm.TokenBalances[i]
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
