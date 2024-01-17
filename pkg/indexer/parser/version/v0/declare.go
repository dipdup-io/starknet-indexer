package v0

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/receiver"
)

// ParseDeclare -
func (parser Parser) ParseDeclare(ctx context.Context, version data.Felt, raw *data.Declare, block storage.Block, receiverTx receiver.Transaction, trace sequencer.Trace) (storage.Declare, *storage.Fee, error) {
	tx := storage.Declare{
		ID:     parser.Resolver.NextTxId(),
		Height: block.Height,
		Time:   block.Time,
		Status: block.Status,
		Hash:   trace.TransactionHash.Bytes(),
		MaxFee: raw.MaxFee.Decimal(),
		Nonce:  raw.Nonce.Decimal(),
	}

	if trace.RevertedError != "" {
		tx.Status = storage.StatusReverted
		tx.Error = &trace.RevertedError
	}

	var err error
	tx.Version, err = version.Uint64()
	if err != nil {
		return tx, nil, err
	}

	if class, err := parser.Resolver.FindClassByHash(ctx, raw.ClassHash, tx.Height); err != nil {
		return tx, nil, err
	} else if class != nil {
		tx.Class = *class
		tx.ClassID = class.ID
	}

	if address, err := parser.Resolver.FindAddressByHash(ctx, raw.ContractAddress); err != nil {
		return tx, nil, err
	} else if address != nil {
		tx.ContractID = &address.ID
		tx.Contract = *address
		tx.Contract.Height = tx.Height
		if tx.ClassID != 0 {
			tx.Contract.ClassID = &tx.Class.ID
			tx.Contract.Class = tx.Class
		}
	}

	if address, err := parser.Resolver.FindAddressByHash(ctx, raw.SenderAddress); err != nil {
		return tx, nil, err
	} else if address != nil {
		tx.SenderID = &address.ID
		tx.Sender = *address
		tx.Sender.Height = tx.Height
	}

	var proxyId uint64
	if tx.Class.Type.Is(storage.ClassTypeProxy) && tx.ContractID != nil {
		proxyId = *tx.ContractID
	}

	txCtx := parserData.NewTxContextFromDeclare(tx, proxyId)

	if trace.FunctionInvocation != nil {
		if len(trace.FunctionInvocation.Events) > 0 {
			classAbi, err := parser.Cache.GetAbiByClassHash(ctx, tx.Class.Hash)
			if err != nil {
				return tx, nil, err
			}
			tx.Events, err = parseEvents(ctx, parser.EventParser, txCtx, classAbi, trace.FunctionInvocation.Events)
			if err != nil {
				return tx, nil, err
			}
			tx.Transfers, err = parser.TransferParser.ParseEvents(ctx, txCtx, tx.Contract, tx.Events)
			if err != nil {
				return tx, nil, err
			}
		}

		tx.Messages, err = parseMessages(ctx, parser.MessageParser, txCtx, trace.FunctionInvocation.Messages)
		if err != nil {
			return tx, nil, err
		}

		tx.Internals, err = parseInternals(ctx, parser.InternalTxParser, txCtx, trace.FunctionInvocation.InternalCalls)
		if err != nil {
			return tx, nil, err
		}
	}

	if trace.FeeTransferInvocation != nil {
		fee, err := parser.FeeParser.ParseInvocation(ctx, txCtx, *trace.FeeTransferInvocation)
		if err != nil {
			return tx, nil, nil
		}
		if fee != nil {
			return tx, fee, nil
		}
	} else {
		transfer, err := parser.FeeParser.ParseActualFee(ctx, txCtx, receiverTx.ActualFee)
		if err != nil {
			return tx, nil, nil
		}
		if transfer != nil {
			tx.Transfers = append(tx.Transfers, *transfer)
		}
	}

	return tx, nil, nil
}
