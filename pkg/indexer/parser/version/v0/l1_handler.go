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
func (parser Parser) ParseL1Handler(ctx context.Context, raw *data.L1Handler, block storage.Block, trace sequencer.Trace, receipts sequencer.Receipt) (storage.L1Handler, *storage.Fee, error) {
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
		return tx, nil, err
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
			return tx, nil, err
		}
	}

	if len(tx.CallData) > 0 {
		if _, ok := contractAbi.GetL1HandlerBySelector(encoding.EncodeHex(tx.EntrypointSelector)); !ok {
			class, err := parser.Cache.GetClassById(ctx, *tx.Contract.ClassID)
			if err != nil {
				return tx, nil, err
			}
			contractAbi, err = parser.Resolver.Proxy(ctx, parserData.NewEmptyTxContext(), *class, tx.Contract)
			if err != nil {
				return tx, nil, err
			}
			if class.Type.Is(storage.ClassTypeProxy) {
				proxyId = tx.ContractID
			}
		}

		tx.ParsedCalldata, tx.Entrypoint, err = decode.CalldataForL1Handler(contractAbi, tx.EntrypointSelector, tx.CallData)
		if err != nil {
			if !errors.Is(err, abi.ErrNoLenField) {
				return tx, nil, err
			}
		}
	}

	txCtx := parserData.NewTxContextFromL1Hadler(tx, proxyId)

	if trace.FunctionInvocation != nil {
		tx.Events, err = parseEvents(ctx, parser.EventParser, txCtx, contractAbi, trace.FunctionInvocation.Events)
		if err != nil {
			return tx, nil, err
		}
		tx.Transfers, err = parser.TransferParser.ParseEvents(ctx, txCtx, tx.Contract, tx.Events)
		if err != nil {
			return tx, nil, err
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
