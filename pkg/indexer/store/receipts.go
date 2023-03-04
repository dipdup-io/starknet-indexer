package store

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (store *Store) saveInternals(
	ctx context.Context,
	tx storage.Transaction,
	result parser.Result,
	internals []models.Internal,
) error {
	if len(internals) == 0 {
		return nil
	}

	for i := range internals {
		ptr := &internals[i]

		if ptr.ContractID == 0 {
			if address, ok := result.Addresses[encoding.EncodeHex(ptr.Contract.Hash)]; ok {
				ptr.ContractID = address.ID
			}
		}

		if ptr.CallerID == 0 {
			if address, ok := result.Addresses[encoding.EncodeHex(ptr.Contract.Hash)]; ok {
				ptr.ContractID = address.ID
			}
		}

		if ptr.ClassID == 0 {
			if class, ok := result.Classes[encoding.EncodeHex(ptr.Class.Hash)]; ok {
				ptr.ClassID = class.ID
				store.cache.SetClassByHash(ptr.Class)
			}
		}

		if err := tx.Add(ctx, ptr); err != nil {
			return err
		}

		if len(ptr.Internals) > 0 {
			for j := range ptr.Internals {
				ptrInt := &ptr.Internals[j]
				ptrInt.InternalID = &ptr.ID
			}

			if err := store.saveInternals(ctx, tx, result, ptr.Internals); err != nil {
				return err
			}
		}
	}

	return nil
}

func (store *Store) saveMessage(ctx context.Context, tx storage.Transaction, msg *models.Message) error {
	if len(msg.From.Hash) > 0 {
		msg.FromID = msg.From.ID
	}

	if len(msg.To.Hash) > 0 {
		msg.ToID = msg.To.ID
	}

	return tx.Add(ctx, msg)
}

func (store *Store) saveEvent(ctx context.Context, tx storage.Transaction, event *models.Event) error {
	if len(event.From.Hash) > 0 {
		event.FromID = event.From.ID
	}

	return tx.Add(ctx, event)
}
