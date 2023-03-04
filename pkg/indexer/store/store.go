package store

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// Store -
type Store struct {
	classes      models.IClass
	address      models.IAddress
	cache        *cache.Cache
	transactable storage.Transactable
}

// New -
func New(
	cache *cache.Cache,
	classes models.IClass,
	address models.IAddress,
	transactable storage.Transactable,
) *Store {
	return &Store{
		cache:        cache,
		classes:      classes,
		address:      address,
		transactable: transactable,
	}
}

// Save -
func (store *Store) Save(
	ctx context.Context,
	result parser.Result,
) error {
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

	for _, address := range result.Addresses {
		if address.ClassID == nil || *address.ClassID == 0 {
			if class, ok := result.Classes[encoding.EncodeHex(address.Class.Hash)]; ok {
				address.ClassID = &class.ID
			}
		}
		if err := tx.Add(ctx, address); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if err := tx.Add(ctx, &result.Block); err != nil {
		return tx.HandleError(ctx, err)
	}

	if result.Block.StorageDiffCount > 0 {
		for i := range result.Block.StorageDiffs {
			if err := tx.Add(ctx, &result.Block.StorageDiffs[i]); err != nil {
				return tx.HandleError(ctx, err)
			}
		}
	}

	if result.Block.TxCount > 0 {
		if err := store.saveDeclare(ctx, tx, result); err != nil {
			return tx.HandleError(ctx, err)
		}

		if err := store.saveDeploy(ctx, tx, result); err != nil {
			return tx.HandleError(ctx, err)
		}

		if err := store.saveDeployAccount(ctx, tx, result); err != nil {
			return tx.HandleError(ctx, err)
		}

		if err := store.saveInvokeV0(ctx, tx, result); err != nil {
			return tx.HandleError(ctx, err)
		}

		if err := store.saveInvokeV1(ctx, tx, result); err != nil {
			return tx.HandleError(ctx, err)
		}

		if err := store.saveL1Handler(ctx, tx, result); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if err := tx.Flush(ctx); err != nil {
		return tx.HandleError(ctx, err)
	}
	return nil
}
