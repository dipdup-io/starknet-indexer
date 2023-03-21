package v0

import (
	"context"
	"errors"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
)

// ParseInvokeV1 -
func (parser Parser) ParseInvokeV1(ctx context.Context, raw *data.Invoke, block storage.Block, trace sequencer.Trace, receipts sequencer.Receipt) (storage.Invoke, error) {
	tx := storage.Invoke{
		ID:                 parser.Resolver.NextTxId(),
		Height:             block.Height,
		Time:               block.Time,
		Status:             block.Status,
		Hash:               trace.TransactionHash.Bytes(),
		Signature:          raw.Signature,
		CallData:           raw.Calldata,
		MaxFee:             raw.MaxFee.Decimal(),
		Nonce:              raw.Nonce.Decimal(),
		EntrypointSelector: encoding.ExecuteEntrypointSelector,
		Entrypoint:         encoding.ExecuteEntrypoint,
		Version:            1,
	}

	if address, err := parser.Resolver.FindAddressByHash(ctx, raw.ContractAddress); err != nil {
		return tx, err
	} else if address != nil {
		tx.ContractID = address.ID
		tx.Contract = *address
		tx.Contract.Height = tx.Height
	}

	if len(tx.CallData) > 0 {
		parsed, err := abi.DecodeExecuteCallData(tx.CallData)
		if err != nil {
			if !errors.Is(err, abi.ErrNoLenField) {
				return tx, err
			}
		}
		tx.ParsedCalldata = parsed
	}

	class, err := parser.Cache.GetClassForAddress(ctx, tx.Contract.Hash)
	if err != nil {
		return tx, err
	}

	var proxyId uint64
	if class.Type.Is(storage.ClassTypeProxy) {
		proxyId = tx.ContractID
	}

	txCtx := parserData.NewTxContextFromInvoke(tx, proxyId)

	tx.Transfers, err = parser.TransferParser.ParseCalldata(ctx, txCtx, tx.Entrypoint, tx.ParsedCalldata)
	if err != nil {
		return tx, err
	}
	txCtx.TransfersCount = len(tx.Transfers)

	if trace.FunctionInvocation != nil {
		if len(trace.FunctionInvocation.Events) > 0 {
			contractAbi, err := parser.Cache.GetAbiByAddress(ctx, tx.Contract.Hash)
			if err != nil {
				return tx, err
			}
			tx.Events, err = parseEvents(ctx, parser.EventParser, txCtx, contractAbi, trace.FunctionInvocation.Events)
			if err != nil {
				return tx, err
			}
			tx.Transfers, err = parser.TransferParser.ParseEvents(ctx, txCtx, tx.Contract, tx.Events)
			if err != nil {
				return tx, err
			}
			if err := parser.ProxyUpgrader.Parse(ctx, tx.Contract, tx.Events); err != nil {
				return tx, err
			}
		}

		var err error
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
		fee.InvokeID = &tx.ID
		tx.Fee = fee
	}

	return tx, nil
}
