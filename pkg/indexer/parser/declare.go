package parser

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
)

func (parser *Parser) getDeclare(ctx context.Context, raw *data.Declare, block storage.Block, trace sequencer.Trace) (storage.Declare, error) {
	tx := storage.Declare{
		Height:    block.Height,
		Time:      block.Time,
		Status:    block.Status,
		Hash:      encoding.MustDecodeHex(trace.TransactionHash),
		Signature: raw.Signature,
		MaxFee:    decimalFromHex(raw.MaxFee),
		Nonce:     decimalFromHex(raw.Nonce),
	}

	if raw.ClassHash != "" {
		tx.Class = storage.Class{
			Hash: encoding.MustDecodeHex(raw.ClassHash),
		}

		if err := parser.findClass(ctx, &tx.Class); err != nil {
			return tx, err
		}
		tx.ClassID = tx.Class.ID
	}

	if raw.ContractAddress != "" {
		tx.Contract = storage.Address{
			Hash: encoding.MustDecodeHex(raw.ContractAddress),
		}

		if err := parser.findAddress(ctx, &tx.Contract); err != nil {
			return tx, err
		}
		tx.ContractID = &tx.Contract.ID

		if tx.ClassID != 0 {
			tx.Contract.ClassID = &tx.Class.ID
			tx.Contract.Class = tx.Class
		}
	}

	if raw.SenderAddress != "" {
		tx.Sender = storage.Address{
			Hash: encoding.MustDecodeHex(raw.SenderAddress),
		}

		if err := parser.findAddress(ctx, &tx.Sender); err != nil {
			return tx, err
		}
		tx.SenderID = &tx.Sender.ID
	}

	if trace.FunctionInvocation != nil {
		if len(trace.FunctionInvocation.Events) > 0 {
			classAbi, err := parser.cache.GetAbiByClassHash(ctx, tx.Class.Hash)
			if err != nil {
				return tx, err
			}

			for i := range trace.FunctionInvocation.Events {
				event, err := parser.getEvent(block, classAbi, trace.FunctionInvocation.Events[i])
				if err != nil {
					return tx, err
				}
				tx.Events = append(tx.Events, event)
			}
		}
		for i := range trace.FunctionInvocation.Messages {
			msg, err := parser.getMessage(ctx, block, trace.FunctionInvocation.Messages[i])
			if err != nil {
				return tx, err
			}
			tx.Messages = append(tx.Messages, msg)
		}
		for i := range trace.FunctionInvocation.InternalCalls {
			internal, err := parser.getInternal(ctx, block, trace.FunctionInvocation.InternalCalls[i], tx.Hash, tx.Status)
			if err != nil {
				return tx, err
			}
			tx.Internals = append(tx.Internals, internal)
		}
	}

	return tx, nil
}
