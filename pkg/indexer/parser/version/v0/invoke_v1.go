package v0

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/decode"
	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/pkg/errors"
)

// ParseInvokeV1 -
func (parser Parser) ParseInvokeV1(ctx context.Context, raw *data.Invoke, block storage.Block, trace sequencer.Trace, receipts sequencer.Receipt) (storage.Invoke, *storage.Fee, error) {
	tx := storage.Invoke{
		ID:                 parser.Resolver.NextTxId(),
		Height:             block.Height,
		Time:               block.Time,
		Status:             block.Status,
		Hash:               trace.TransactionHash.Bytes(),
		CallData:           raw.Calldata,
		MaxFee:             raw.MaxFee.Decimal(),
		Nonce:              raw.Nonce.Decimal(),
		EntrypointSelector: encoding.ExecuteEntrypointSelector,
		Entrypoint:         encoding.ExecuteEntrypoint,
		Version:            1,
	}

	if trace.RevertedError != "" {
		tx.Status = storage.StatusReverted
		tx.Error = &trace.RevertedError
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
			return tx, nil, errors.Wrap(err, "FindAddressByHash")
		} else if address != nil {
			tx.ContractID = address.ID
			tx.Contract = *address
			tx.Contract.Height = tx.Height
		}
	}

	if tx.Status != storage.StatusReverted && len(tx.CallData) > 0 {
		var found bool
		contractAbi, err := parser.Cache.GetAbiByAddress(ctx, tx.Contract.Hash)
		if err == nil {
			if _, found = contractAbi.GetFunctionBySelector(encoding.EncodeHex(encoding.ExecuteEntrypointSelector)); found {
				parsed, _, err := decode.CalldataBySelector(contractAbi, tx.EntrypointSelector, tx.CallData)
				if err != nil {
					if !errors.Is(err, abi.ErrNoLenField) && !errors.Is(err, abi.ErrTooShortCallData) {
						return tx, nil, errors.Wrap(err, "custom __execute__ function")
					}
				}
				tx.ParsedCalldata = parsed
			}
		}

		if !found {
			parsed, err := abi.DecodeExecuteCallData(tx.CallData)
			if err != nil {
				if !errors.Is(err, abi.ErrNoLenField) && !errors.Is(err, abi.ErrTooShortCallData) {
					return tx, nil, errors.Wrap(err, "default __execute__ function")
				}
			}
			tx.ParsedCalldata = parsed
		}
	}

	txCtx := parserData.NewTxContextFromInvoke(tx, 0)

	if tx.Status != storage.StatusReverted {
		var (
			class       storage.Class
			err         error
			contractAbi abi.Abi
		)

		if trace.FunctionInvocation != nil && trace.FunctionInvocation.ClassHash.Length() > 0 {
			class, err = parser.Cache.GetClassByHash(ctx, trace.FunctionInvocation.ClassHash.Bytes())
		} else {
			class, err = parser.Cache.GetClassForAddress(ctx, tx.Contract.Hash)
		}
		if err != nil {
			return tx, nil, errors.Wrap(err, "receive class hash")
		}

		if class.Type.Is(storage.ClassTypeProxy) {
			txCtx.ProxyId = tx.ContractID
		}

		tx.Transfers, err = parser.TransferParser.ParseCalldata(ctx, txCtx, tx.Entrypoint, tx.ParsedCalldata)
		if err != nil {
			return tx, nil, errors.Wrap(err, "transfer parse")
		}

		if trace.FunctionInvocation != nil {
			if len(trace.FunctionInvocation.Events) > 0 {
				contractAbi, err = parser.Cache.GetAbiByAddress(ctx, tx.Contract.Hash)
				if err != nil {
					return tx, nil, errors.Wrapf(err, "get abi: %x", tx.Contract.Hash)
				}
				tx.Events, err = parseEvents(ctx, parser.EventParser, txCtx, contractAbi, trace.FunctionInvocation.Events)
				if err != nil {
					return tx, nil, errors.Wrap(err, "parse events")
				}
				tx.Transfers, err = parser.TransferParser.ParseEvents(ctx, txCtx, tx.Contract, tx.Events)
				if err != nil {
					return tx, nil, errors.Wrap(err, "TransferParser.ParseEvents")
				}
				if err := parser.ProxyUpgrader.Parse(ctx, txCtx, tx.Contract, tx.Events, tx.Entrypoint, tx.ParsedCalldata); err != nil {
					return tx, nil, errors.Wrap(err, "ProxyUpgrader.Parse")
				}
			}

			var err error
			tx.Messages, err = parseMessages(ctx, parser.MessageParser, txCtx, trace.FunctionInvocation.Messages)
			if err != nil {
				return tx, nil, errors.Wrap(err, "parseMessages")
			}

			tx.Internals, err = parseInternals(ctx, parser.InternalTxParser, txCtx, trace.FunctionInvocation.InternalCalls)
			if err != nil {
				return tx, nil, err
			}
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
