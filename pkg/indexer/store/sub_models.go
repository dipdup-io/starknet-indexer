package store

import (
	"context"

	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

const copyThreashold = 25

type subModels struct {
	Transfers     []models.Transfer
	TokenBalances []models.TokenBalance
	Events        []models.Event
	Messages      []any
	Internals     []models.Internal
}

func newSubModels(
	internalsStorage models.IInternal,
	transfersStorage models.ITransfer,
	eventsStorage models.IEvent,
) *subModels {
	return &subModels{
		Transfers:     make([]models.Transfer, 0),
		Events:        make([]models.Event, 0),
		Messages:      make([]any, 0),
		Internals:     make([]models.Internal, 0),
		TokenBalances: make([]models.TokenBalance, 0),
	}
}

// Save -
func (sm *subModels) Save(ctx context.Context, tx postgres.Transaction) error {
	if err := bulkSaveWithCopy(ctx, tx, sm.Internals); err != nil {
		return err
	}

	if err := bulkSaveWithCopy(ctx, tx, sm.Events); err != nil {
		return err
	}

	if len(sm.Messages) > 0 {
		if err := tx.BulkSave(ctx, sm.Messages); err != nil {
			return err
		}
	}

	if err := bulkSaveWithCopy(ctx, tx, sm.Transfers); err != nil {
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

func (sm *subModels) saveTokenBalanceUpdates(ctx context.Context, tx postgres.Transaction) error {
	updates := make(map[string]*models.TokenBalance)
	for i := range sm.TokenBalances {
		key := sm.TokenBalances[i].String()
		if update, ok := updates[key]; ok {
			update.Balance = update.Balance.Add(sm.TokenBalances[i].Balance)
		} else {
			updates[key] = &sm.TokenBalances[i]
		}
	}

	arr := make([]*models.TokenBalance, 0)
	for _, update := range updates {
		arr = append(arr, update)
	}

	return tx.SaveTokenBalanceUpdates(ctx, arr...)
}

func bulkSaveWithCopy[M models.CopiableModel](ctx context.Context, tx storage.Transaction, arr []M) error {
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
		tableName := arr[0].TableName()
		data := make([]storage.Copiable, len(arr))
		for i := range arr {
			data[i] = arr[i]
		}
		return tx.CopyFrom(ctx, tableName, data)
	}
}
