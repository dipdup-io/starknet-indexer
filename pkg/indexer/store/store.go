package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/go-pg/pg/v10"
)

// Store -
type Store struct {
	classes          models.IClass
	address          models.IAddress
	cache            *cache.Cache
	transactable     storage.Transactable
	partitionManager postgres.PartitionManager
}

// New -
func New(
	cache *cache.Cache,
	classes models.IClass,
	address models.IAddress,
	transactable storage.Transactable,
	partitionManager postgres.PartitionManager,
) *Store {
	return &Store{
		cache:            cache,
		classes:          classes,
		address:          address,
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

	for _, class := range result.Classes {
		if err := tx.Add(ctx, class); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if err := saveAddresses(ctx, tx, result); err != nil {
		return tx.HandleError(ctx, err)
	}

	for _, proxy := range result.Proxies {
		if _, err := tx.Exec(ctx, `
			INSERT INTO proxy (contract_id, hash, entity_type, entity_id, entity_hash)
			VALUES(?,?,?,?,?) 
			ON CONFLICT (hash)
			DO 
			UPDATE SET entity_type = excluded.entity_type, entity_id = excluded.entity_id, entity_hash = excluded.entity_hash`,
			proxy.ContractID, proxy.Hash, proxy.EntityType, proxy.EntityID, proxy.EntityHash); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if err := tx.Add(ctx, &result.Block); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := saveStorageDiff(ctx, tx, result); err != nil {
		return tx.HandleError(ctx, err)
	}

	if result.Block.TxCount > 0 {
		sm := newSubModels()

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
	if len(result.Addresses) == 0 {
		return nil
	}

	values := make([]string, 0)
	for _, address := range result.Addresses {
		if address.ClassID == nil || *address.ClassID == 0 {
			if class, ok := result.Classes[encoding.EncodeHex(address.Class.Hash)]; ok {
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

func saveStorageDiff(ctx context.Context, tx storage.Transaction, result parserData.Result) error {
	if result.Block.StorageDiffCount == 0 {
		return nil
	}
	models := make([]any, len(result.Block.StorageDiffs))
	for i := range result.Block.StorageDiffs {
		models[i] = &result.Block.StorageDiffs[i]
	}
	return tx.BulkSave(ctx, models)
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
