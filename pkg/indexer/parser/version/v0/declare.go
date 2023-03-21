package v0

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
)

// ParseDeclare -
func (parser Parser) ParseDeclare(ctx context.Context, version data.Felt, raw *data.Declare, block storage.Block, trace sequencer.Trace, receipts sequencer.Receipt) (storage.Declare, error) {
	tx := storage.Declare{
		ID:        parser.Resolver.NextTxId(),
		Height:    block.Height,
		Time:      block.Time,
		Status:    block.Status,
		Hash:      trace.TransactionHash.Bytes(),
		Signature: raw.Signature,
		MaxFee:    raw.MaxFee.Decimal(),
		Nonce:     raw.Nonce.Decimal(),
	}

	var err error
	tx.Version, err = version.Uint64()
	if err != nil {
		return tx, err
	}

	if class, err := parser.Resolver.FindClassByHash(ctx, raw.ClassHash); err != nil {
		return tx, err
	} else if class != nil {
		tx.Class = *class
		tx.ClassID = class.ID
	}

	if address, err := parser.Resolver.FindAddressByHash(ctx, raw.ContractAddress); err != nil {
		return tx, err
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
		return tx, err
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
				return tx, err
			}
			tx.Events, err = parseEvents(ctx, parser.EventParser, txCtx, classAbi, trace.FunctionInvocation.Events)
			if err != nil {
				return tx, err
			}
			tx.Transfers, err = parser.TransferParser.ParseEvents(ctx, txCtx, tx.Contract, tx.Events)
			if err != nil {
				return tx, err
			}
		}

		tx.Messages, err = parseMessages(ctx, parser.MessageParser, txCtx, trace.FunctionInvocation.Messages)
		if err != nil {
			return tx, err
		}

		tx.Internals, err = parseInternals(ctx, parser.InternalTxParser, txCtx, trace.FunctionInvocation.InternalCalls)
		if err != nil {
			return tx, err
		}
	}

	var fee *storage.Fee
	if trace.FeeTransferInvocation != nil {
		fee, err = parser.FeeParser.ParseInvocation(ctx, txCtx, *trace.FeeTransferInvocation)
	} else {
		fee, err = parser.FeeParser.ParseActualFee(ctx, txCtx, receipts.ActualFee)
	}
	if err != nil {
		return tx, err
	}
	if fee != nil {
		fee.DeclareID = &tx.ID
		tx.Fee = fee
	}

	return tx, nil
}
