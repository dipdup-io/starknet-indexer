package store

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (store *Store) saveInvokeV1(
	ctx context.Context,
	tx storage.Transaction,
	result parser.Result,
) error {
	if result.Block.InvokeV1Count == 0 {
		return nil
	}

	for i := range result.Block.InvokeV1 {
		ptr := &result.Block.InvokeV1[i]

		if ptr.SenderID == 0 {
			if address, ok := result.Addresses[encoding.EncodeHex(ptr.Sender.Hash)]; ok {
				ptr.SenderID = address.ID
			}
		}

		if ptr.ContractID == 0 {
			if address, ok := result.Addresses[encoding.EncodeHex(ptr.Contract.Hash)]; ok {
				ptr.ContractID = address.ID
			}
		}

		if err := tx.Add(ctx, ptr); err != nil {
			return err
		}

		if len(ptr.Internals) > 0 {
			for j := range ptr.Internals {
				ptrInt := &ptr.Internals[j]
				ptrInt.InvokeV1ID = &ptr.ID
			}

			if err := store.saveInternals(ctx, tx, result, ptr.Internals); err != nil {
				return err
			}
		}

		if len(ptr.Events) > 0 {
			for j := range ptr.Events {
				ptrEv := &ptr.Events[j]
				ptrEv.InvokeV1ID = &ptr.ID

				if err := store.saveEvent(ctx, tx, ptrEv); err != nil {
					return err
				}
			}
		}

		if len(ptr.Messages) > 0 {
			for j := range ptr.Messages {
				ptrMsg := &ptr.Messages[j]
				ptrMsg.InvokeV1ID = &ptr.ID

				if err := store.saveMessage(ctx, tx, ptrMsg); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
