package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/go-pg/pg/v10"
)

// Store -
type Store struct {
	classes          models.IClass
	address          models.IAddress
	internals        models.IInternal
	transfers        models.ITransfer
	events           models.IEvent
	diffs            models.IStorageDiff
	invokes          models.IInvoke
	deploys          models.IDeploy
	l1Handlers       models.IL1Handler
	fees             models.IFee
	cache            *cache.Cache
	transactable     storage.Transactable
	partitionManager postgres.PartitionManager
}

// New -
func New(
	cache *cache.Cache,
	classes models.IClass,
	address models.IAddress,
	internals models.IInternal,
	transfers models.ITransfer,
	eventsStorage models.IEvent,
	diffs models.IStorageDiff,
	invokes models.IInvoke,
	deploys models.IDeploy,
	l1Handlers models.IL1Handler,
	fees models.IFee,
	transactable storage.Transactable,
	partitionManager postgres.PartitionManager,
) *Store {
	return &Store{
		cache:            cache,
		classes:          classes,
		address:          address,
		internals:        internals,
		transfers:        transfers,
		events:           eventsStorage,
		diffs:            diffs,
		invokes:          invokes,
		deploys:          deploys,
		l1Handlers:       l1Handlers,
		fees:             fees,
		transactable:     transactable,
		partitionManager: partitionManager,
	}
}

// Save -
func (store *Store) Save(
	ctx context.Context,
	result parserData.Result,
) error {
	if err := store.partitionManager.CreatePartitions(ctx, result.Block.Time); err != nil {
		return err
	}

	tx, err := store.transactable.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	for _, class := range result.Context.Classes() {
		if _, err := tx.Exec(ctx, `INSERT INTO class (id, type, hash, abi, height)
			VALUES (?,?,?,?,?)
			ON CONFLICT (id)
			DO 
			UPDATE SET abi = excluded.abi`, class.ID, class.Type, class.Hash, class.Abi, class.Height); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if err := saveAddresses(ctx, tx, result); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := store.saveProxies(ctx, tx, result.Context.Proxies()); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Add(ctx, &result.Block); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := saveStorageDiff(ctx, tx, store.diffs, result); err != nil {
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

func saveAddresses(ctx context.Context, tx storage.Transaction, result parserData.Result) error {
	addresses := result.Context.Addresses()
	if len(addresses) == 0 {
		return nil
	}

	values := make([]string, 0)
	for _, address := range addresses {
		if address.ClassID == nil || *address.ClassID == 0 {
			if class, ok := result.Context.Classes()[encoding.EncodeHex(address.Class.Hash)]; ok {
				address.ClassID = &class.ID
			}
		}
		if address.ClassID == nil {
			values = append(values, fmt.Sprintf("(%d,NULL,%d,'\\x%x')", address.ID, address.Height, address.Hash))
		} else {
			values = append(values, fmt.Sprintf("(%d,%d,%d,'\\x%x')", address.ID, *address.ClassID, address.Height, address.Hash))
		}
	}

	_, err := tx.Exec(ctx, `INSERT INTO address (id, class_id, height, hash)
	VALUES ?
	ON CONFLICT (hash)
	DO 
	UPDATE SET class_id = excluded.class_id, height = excluded.height`, pg.Safe(strings.Join(values, ",")))
	return err
}

func saveStorageDiff(ctx context.Context, tx storage.Transaction, copiable models.Copiable[models.StorageDiff], result parserData.Result) error {
	if result.Block.StorageDiffCount == 0 {
		return nil
	}
	return bulkSaveWithCopy(ctx, tx, copiable, result.Block.StorageDiffs)
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

		var err error
		switch {
		case internals[i].ERC20 != nil:
			err = tx.Add(ctx, internals[i].ERC20)
		case internals[i].ERC721 != nil:
			err = tx.Add(ctx, internals[i].ERC721)
		case internals[i].ERC1155 != nil:
			err = tx.Add(ctx, internals[i].ERC1155)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (store *Store) saveProxies(ctx context.Context, tx storage.Transaction, proxies data.ProxyMap[*data.ProxyWithAction]) error {
	return proxies.Range(func(_ parserData.ProxyKey, value *parserData.ProxyWithAction) (bool, error) {
		store.cache.SetProxy(value.Hash, value.Selector, value.Proxy)
		return saveProxy(ctx, tx, value)
	})
}

func saveProxy(ctx context.Context, tx storage.Transaction, proxy *data.ProxyWithAction) (bool, error) {
	switch proxy.Action {
	case parserData.ProxyActionAdd, parserData.ProxyActionUpdate:
		if _, err := tx.Exec(ctx, `
			INSERT INTO proxy (contract_id, hash, selector, entity_type, entity_id, entity_hash)
			VALUES(?,?,?,?,?,?) 
			ON CONFLICT (hash, selector)
			DO 
			UPDATE SET entity_type = excluded.entity_type, entity_id = excluded.entity_id, entity_hash = excluded.entity_hash, selector = excluded.selector`,
			proxy.ContractID, proxy.Hash, proxy.Selector, proxy.EntityType, proxy.EntityID, proxy.EntityHash); err != nil {
			return true, err
		}
	case parserData.ProxyActionDelete:
		if _, err := tx.Exec(ctx, `DELETE FROM proxy WHERE hash = ? AND selector = ?`,
			proxy.Hash, proxy.Selector); err != nil {
			return true, err
		}
	}
	return false, nil
}
