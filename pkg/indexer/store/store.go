package store

import (
	"context"

	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// Store -
type Store struct {
	classes          models.IClass
	address          models.IAddress
	internals        models.IInternal
	transfers        models.ITransfer
	events           models.IEvent
	cache            *cache.Cache
	transactable     storage.Transactable
	partitionManager database.RangePartitionManager
}

// New -
func New(
	cache *cache.Cache,
	classes models.IClass,
	address models.IAddress,
	internals models.IInternal,
	transfers models.ITransfer,
	eventsStorage models.IEvent,
	transactable storage.Transactable,
	partitionManager database.RangePartitionManager,
) *Store {
	return &Store{
		cache:            cache,
		classes:          classes,
		address:          address,
		internals:        internals,
		transfers:        transfers,
		events:           eventsStorage,
		transactable:     transactable,
		partitionManager: partitionManager,
	}
}

// Save -
func (store *Store) Save(
	ctx context.Context,
	result parserData.Result,
) error {
	if err := store.partitionManager.CreatePartitions(
		ctx,
		result.Block.Time,
		models.Internal{}.TableName(),
		models.Declare{}.TableName(),
		models.DeployAccount{}.TableName(),
		models.Deploy{}.TableName(),
		models.Event{}.TableName(),
		models.Invoke{}.TableName(),
		models.L1Handler{}.TableName(),
		models.Message{}.TableName(),
		models.Transfer{}.TableName(),
		models.Fee{}.TableName(),
	); err != nil {
		return err
	}

	tx, err := postgres.BeginTransaction(ctx, store.transactable)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	if err := saveClasses(ctx, tx, result.Context.Classes()); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := saveAddresses(ctx, tx, result); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := saveClassReplaces(ctx, tx, store.cache, result.Context.ClassReplaces()); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := store.saveProxies(ctx, tx, result.Context.Proxies()); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Add(ctx, &result.Block); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := saveTokens(ctx, tx, result.Context.Tokens()); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := saveStorageDiff(ctx, tx, result); err != nil {
		return tx.HandleError(ctx, err)
	}

	if result.Block.TxCount > 0 {
		sm := newSubModels(store.internals, store.transfers, store.events)

		if err := store.saveDeclare(ctx, tx, result, sm); err != nil {
			return tx.HandleError(ctx, err)
		}

		if err := store.saveDeploy(ctx, tx, result, sm); err != nil {
			return tx.HandleError(ctx, err)
		}

		if err := store.saveDeployAccount(ctx, tx, result, sm); err != nil {
			return tx.HandleError(ctx, err)
		}

		if err := store.saveInvoke(ctx, tx, result, sm); err != nil {
			return tx.HandleError(ctx, err)
		}

		if err := store.saveL1Handler(ctx, tx, result, sm); err != nil {
			return tx.HandleError(ctx, err)
		}

		if err := store.saveFees(ctx, tx, result, sm); err != nil {
			return tx.HandleError(ctx, err)
		}

		if err := sm.Save(ctx, tx); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if err := tx.Update(ctx, result.State); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Flush(ctx); err != nil {
		return tx.HandleError(ctx, err)
	}
	return nil
}

func saveAddresses(ctx context.Context, tx postgres.Transaction, result parserData.Result) error {
	addresses := result.Context.Addresses()
	if len(addresses) == 0 {
		return nil
	}

	values := make([]*models.Address, 0)
	for _, address := range addresses {
		values = append(values, address)
	}

	return tx.SaveAddresses(ctx, values...)
}

func saveStorageDiff(ctx context.Context, tx storage.Transaction, result parserData.Result) error {
	if result.Block.StorageDiffCount == 0 {
		return nil
	}
	return bulkSaveWithCopy(ctx, tx, result.Block.StorageDiffs)
}

func (store *Store) saveInternals(
	ctx context.Context,
	tx storage.Transaction,
	internals []models.Internal,
	sm *subModels,
) error {
	if len(internals) == 0 {
		return nil
	}

	for i := range internals {
		sm.addInternals(internals[i].Internals)
		sm.addEvents(internals[i].Events)
		sm.addMessages(internals[i].Messages)
		sm.addTransfers(internals[i].Transfers)

		if err := store.saveInternals(ctx, tx, internals[i].Internals, sm); err != nil {
			return err
		}
	}
	return nil
}

func (store *Store) saveProxies(ctx context.Context, tx postgres.Transaction, proxies data.ProxyMap[*models.ProxyUpgrade]) error {
	return proxies.Range(func(_ parserData.ProxyKey, value *models.ProxyUpgrade) (bool, error) {
		store.cache.SetProxy(value.Hash, value.Selector, value.ToProxy())
		if err := tx.Add(ctx, value); err != nil {
			return true, err
		}
		return saveProxy(ctx, tx, value)
	})
}

func saveProxy(ctx context.Context, tx postgres.Transaction, proxy *models.ProxyUpgrade) (bool, error) {
	switch proxy.Action {
	case models.ProxyActionAdd, models.ProxyActionUpdate:
		if err := tx.SaveProxy(ctx, proxy.ToProxy()); err != nil {
			return true, err
		}
	case models.ProxyActionDelete:
		if err := tx.DeleteProxy(ctx, proxy.ToProxy()); err != nil {
			return true, err
		}
	}
	return false, nil
}
func saveTokens(ctx context.Context, tx postgres.Transaction, tokens map[string]*models.Token) error {
	if len(tokens) == 0 {
		return nil
	}

	arr := make([]*models.Token, 0)
	for _, token := range tokens {
		arr = append(arr, token)
	}

	return tx.SaveTokens(ctx, arr...)
}

func saveClasses(ctx context.Context, tx postgres.Transaction, classes map[string]*models.Class) error {
	if len(classes) == 0 {
		return nil
	}

	arr := make([]*models.Class, 0)
	for _, class := range classes {
		arr = append(arr, class)
	}

	return tx.SaveClasses(ctx, arr...)
}

func saveClassReplaces(ctx context.Context, tx postgres.Transaction, cache *cache.Cache, replaces map[string]*models.ClassReplace) error {
	if len(replaces) == 0 {
		return nil
	}

	arr := make([]*models.ClassReplace, 0)
	for _, replace := range replaces {
		arr = append(arr, replace)
		cache.SetAbiByAddress(replace.NextClass, replace.Contract.Hash)
		cache.SetAddress(ctx, replace.Contract)
	}

	return tx.SaveClassReplaces(ctx, arr...)
}
