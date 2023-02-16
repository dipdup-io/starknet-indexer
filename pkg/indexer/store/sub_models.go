package store

import (
	"context"
	"fmt"
	"strings"

	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/go-pg/pg/v10"
)

const copyThreashold = 25

type subModels struct {
	Transfers     []models.Transfer
	TokenBalances []models.TokenBalance
	Events        []models.Event
	Messages      []any
	Internals     []models.Internal

	internalsStorage models.IInternal
	transfersStorage models.ITransfer
	eventsStorage    models.IEvent
}

func newSubModels(
	internalsStorage models.IInternal,
	transfersStorage models.ITransfer,
	eventsStorage models.IEvent,
) *subModels {
	return &subModels{
		Transfers:        make([]models.Transfer, 0),
		Events:           make([]models.Event, 0),
		Messages:         make([]any, 0),
		Internals:        make([]models.Internal, 0),
		TokenBalances:    make([]models.TokenBalance, 0),
		internalsStorage: internalsStorage,
		transfersStorage: transfersStorage,
		eventsStorage:    eventsStorage,
	}
}

// Save -
func (sm *subModels) Save(ctx context.Context, tx storage.Transaction) error {
	if err := bulkSaveWithCopy[models.Internal](ctx, tx, sm.internalsStorage, sm.Internals); err != nil {
		return err
	}

	if err := bulkSaveWithCopy[models.Event](ctx, tx, sm.eventsStorage, sm.Events); err != nil {
		return err
	}

	if len(sm.Messages) > 0 {
		if err := tx.BulkSave(ctx, sm.Messages); err != nil {
			return err
		}
	}

	if err := bulkSaveWithCopy[models.Transfer](ctx, tx, sm.transfersStorage, sm.Transfers); err != nil {
		return err
	}

	if len(sm.TokenBalances) > 0 {
		if err := sm.saveTokenBalanceUpdates(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

func (sm *subModels) addEvents(events []models.Event) {
	sm.Events = append(sm.Events, events...)
}

func (sm *subModels) addMessages(msgs []models.Message) {
	for i := range msgs {
		sm.Messages = append(sm.Messages, &msgs[i])
	}
}

func (sm *subModels) addTransfers(transfers []models.Transfer) {
	sm.Transfers = append(sm.Transfers, transfers...)
	for i := range transfers {
		sm.TokenBalances = append(sm.TokenBalances, transfers[i].TokenBalanceUpdates()...)
	}
}

func (sm *subModels) addInternals(internals []models.Internal) {
	sm.Internals = append(sm.Internals, internals...)
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

func bulkSaveWithCopy[M storage.Model](ctx context.Context, tx storage.Transaction, copiable models.Copiable[M], arr []M) error {
	switch {
	case len(arr) == 0:
		return nil
	case len(arr) < copyThreashold:
		data := make([]any, len(arr))
		for i := range arr {
			data[i] = &arr[i]
		}
		return tx.BulkSave(ctx, data)
	default:
		reader, query, err := copiable.InsertByCopy(arr)
		if err != nil {
			return err
		}
		// return tx.CopyFrom(io.TeeReader(reader, os.Stdout), query)
		return tx.CopyFrom(reader, query)
	}
}
