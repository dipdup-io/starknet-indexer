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
func (parser Parser) ParseInvokeV1(ctx context.Context, raw *data.Invoke, block storage.Block, trace sequencer.Trace, receipts sequencer.Receipt) (storage.Invoke, *storage.Fee, error) {
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

	var contract data.Felt
	switch {
	case raw.ContractAddress != "":
		contract = raw.ContractAddress
	case raw.SenderAddress != "":
		contract = raw.SenderAddress
	}
	if contract != "" {
		if address, err := parser.Resolver.FindAddressByHash(ctx, contract); err != nil {
			return tx, nil, err
		} else if address != nil {
			tx.ContractID = address.ID
			tx.Contract = *address
			tx.Contract.Height = tx.Height
		}
	}

	if len(tx.CallData) > 0 {
		parsed, err := abi.DecodeExecuteCallData(tx.CallData)
		if err != nil {
			if !errors.Is(err, abi.ErrNoLenField) {
				return tx, nil, err
			}
		}
		tx.ParsedCalldata = parsed
	}

	var (
		proxyId uint64
		class   storage.Class
		err     error
	)

	if trace.FunctionInvocation != nil && trace.FunctionInvocation.ClassHash.Length() > 0 {
		class, err = parser.Cache.GetClassByHash(ctx, trace.FunctionInvocation.ClassHash.Bytes())
	} else {
		class, err = parser.Cache.GetClassForAddress(ctx, tx.Contract.Hash)
	}
	if err != nil {
		return tx, nil, err
	}

	if class.Type.Is(storage.ClassTypeProxy) {
		proxyId = tx.ContractID
	}

	txCtx := parserData.NewTxContextFromInvoke(tx, proxyId)

	tx.Transfers, err = parser.TransferParser.ParseCalldata(ctx, txCtx, tx.Entrypoint, tx.ParsedCalldata)
	if err != nil {
		return tx, nil, err
	}

	if trace.FunctionInvocation != nil {
		if len(trace.FunctionInvocation.Events) > 0 {
			contractAbi, err := parser.Cache.GetAbiByAddress(ctx, tx.Contract.Hash)
			if err != nil {
				return tx, nil, err
			}
			tx.Events, err = parseEvents(ctx, parser.EventParser, txCtx, contractAbi, trace.FunctionInvocation.Events)
			if err != nil {
				return tx, nil, err
			}
			tx.Transfers, err = parser.TransferParser.ParseEvents(ctx, txCtx, tx.Contract, tx.Events)
			if err != nil {
				return tx, nil, err
			}
			if err := parser.ProxyUpgrader.Parse(ctx, txCtx, tx.Contract, tx.Events, tx.Entrypoint, tx.ParsedCalldata); err != nil {
				return tx, nil, err
			}
		}

		var err error
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
		transfer, err := parser.FeeParser.ParseActualFee(ctx, txCtx, receipts.ActualFee)
		if err != nil {
			return tx, nil, nil
		}
		if transfer != nil {
			tx.Transfers = append(tx.Transfers, *transfer)
		}
	}

	return tx, nil, nil
}
