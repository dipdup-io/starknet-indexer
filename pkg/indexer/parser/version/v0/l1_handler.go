package v0

import (
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

// ParseL1Handler -
func (parser Parser) ParseL1Handler(ctx context.Context, raw *data.L1Handler, block storage.Block, trace sequencer.Trace, receipts sequencer.Receipt) (storage.L1Handler, error) {
	tx := storage.L1Handler{
		ID:                 parser.Resolver.NextTxId(),
		Height:             block.Height,
		Time:               block.Time,
		Status:             block.Status,
		Hash:               trace.TransactionHash.Bytes(),
		EntrypointSelector: raw.EntrypointSelector.Bytes(),
		CallData:           raw.Calldata,
		Nonce:              raw.Nonce.Decimal(),
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
	}

	if len(tx.CallData) > 0 {
		if _, ok := contractAbi.GetL1HandlerBySelector(encoding.EncodeHex(tx.EntrypointSelector)); !ok {
			class, err := parser.Cache.GetClassById(ctx, *tx.Contract.ClassID)
			if err != nil {
				return tx, err
			}
			contractAbi, err = parser.Resolver.Proxy(ctx, *class, tx.Contract)
			if err != nil {
				return tx, err
			}
			if class.Type.Is(storage.ClassTypeProxy) {
				proxyId = tx.ContractID
			}
		}

		tx.ParsedCalldata, tx.Entrypoint, err = decode.CalldataForL1Handler(parser.Cache, contractAbi, tx.EntrypointSelector, tx.CallData)
		if err != nil {
			if !errors.Is(err, abi.ErrNoLenField) {
				return tx, err
			}
		}
	}

	txCtx := parserData.NewTxContextFromL1Hadler(tx, proxyId)

	if trace.FunctionInvocation != nil {
		tx.Events, err = parseEvents(ctx, parser.EventParser, txCtx, contractAbi, trace.FunctionInvocation.Events)
		if err != nil {
			return tx, err
		}
		tx.Transfers, err = parser.TransferParser.ParseEvents(ctx, txCtx, tx.Contract, tx.Events)
		if err != nil {
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
		fee.L1HandlerID = &tx.ID
		tx.Fee = fee
	}

	return tx, nil
}
