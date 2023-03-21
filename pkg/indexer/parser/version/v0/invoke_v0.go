package v0

import (
	"bytes"
	"context"
	"errors"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/decode"
	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/helpers"
)

// ParseInvokeV0 -
func (parser Parser) ParseInvokeV0(ctx context.Context, raw *data.Invoke, block storage.Block, trace sequencer.Trace, receipts sequencer.Receipt) (storage.Invoke, error) {
	tx := storage.Invoke{
		ID:                 parser.Resolver.NextTxId(),
		Height:             block.Height,
		Time:               block.Time,
		Status:             block.Status,
		Hash:               trace.TransactionHash.Bytes(),
		EntrypointSelector: raw.EntrypointSelector.Bytes(),
		Signature:          raw.Signature,
		CallData:           raw.Calldata,
		MaxFee:             raw.MaxFee.Decimal(),
		Nonce:              raw.Nonce.Decimal(),
		Version:            0,

		Events:    make([]storage.Event, 0),
		Messages:  make([]storage.Message, 0),
		Internals: make([]storage.Internal, 0),
	}

	if address, err := parser.Resolver.FindAddressByHash(ctx, raw.ContractAddress); err != nil {
		return tx, err
	} else if address != nil {
		tx.ContractID = address.ID
		tx.Contract = *address
		tx.Contract.Height = tx.Height
	}

	var (
		contractAbi abi.Abi
		err         error
		proxyId     uint64
	)

	if helpers.NeedDecode(tx.CallData, trace.FunctionInvocation) {
		contractAbi, err = parser.Cache.GetAbiByAddress(ctx, tx.Contract.Hash)
		if err != nil {
			return tx, err
		}

		if _, ok := contractAbi.GetFunctionBySelector(encoding.EncodeHex(tx.EntrypointSelector)); !ok {
			var class storage.Class
			if tx.Contract.ClassID == nil {
				class, err = parser.Cache.GetClassForAddress(ctx, tx.Contract.Hash)
				if err != nil {
					return tx, err
				}
			} else {
				c, err := parser.Cache.GetClassById(ctx, *tx.Contract.ClassID)
				if err != nil {
					return tx, err
				}
				class = *c
			}

			contractAbi, err = parser.Resolver.Proxy(ctx, class, tx.Contract)
			if err != nil {
				return tx, err
			}

			if class.Type.Is(storage.ClassTypeProxy) {
				proxyId = tx.ContractID
			}
		}
	}

	isExecute := bytes.Equal(tx.EntrypointSelector, encoding.ExecuteEntrypointSelector)
	_, hasExecute := contractAbi.Functions[encoding.ExecuteEntrypoint]

	if len(tx.CallData) > 0 {
		if isExecute && !hasExecute {
			tx.Entrypoint = encoding.ExecuteEntrypoint
			tx.ParsedCalldata, err = abi.DecodeExecuteCallData(tx.CallData)
		} else {
			tx.ParsedCalldata, tx.Entrypoint, err = decode.CalldataBySelector(parser.Cache, contractAbi, tx.EntrypointSelector, tx.CallData)
		}
		if err != nil {
			if !errors.Is(err, abi.ErrNoLenField) {
				return tx, err
			}
		}
	}

	txCtx := parserData.NewTxContextFromInvoke(tx, proxyId)

	tx.Transfers, err = parser.TransferParser.ParseCalldata(ctx, txCtx, tx.Entrypoint, tx.ParsedCalldata)
	if err != nil {
		return tx, err
	}
	txCtx.TransfersCount = len(tx.Transfers)

	if trace.FunctionInvocation != nil {
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
