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
)

// ParseDeploy -
func (parser Parser) ParseDeploy(ctx context.Context, raw *data.Deploy, block storage.Block, trace sequencer.Trace, receipts sequencer.Receipt) (storage.Deploy, *storage.Fee, error) {
	tx := storage.Deploy{
		ID:                  parser.Resolver.NextTxId(),
		Height:              block.Height,
		Time:                block.Time,
		Status:              block.Status,
		Hash:                trace.TransactionHash.Bytes(),
		ContractAddressSalt: encoding.MustDecodeHex(raw.ContractAddressSalt),
		ConstructorCalldata: raw.ConstructorCalldata,
	}

	if trace.RevertedError != "" {
		tx.Status = storage.StatusReverted
		tx.Error = &trace.RevertedError
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
		tx.ContractID = address.ID
		tx.Contract = *address
		tx.Contract.Height = tx.Height
		if tx.ClassID != 0 {
			tx.Contract.ClassID = &tx.Class.ID
			tx.Contract.Class = tx.Class
		}
	}

	classAbi, err := parser.Cache.GetAbiByClassHash(ctx, tx.Class.Hash)
	if err != nil {
		return tx, nil, err
	}

	if tx.Status != storage.StatusReverted && len(tx.ConstructorCalldata) > 0 {
		tx.ParsedCalldata, err = decode.CalldataForConstructor(classAbi, tx.ConstructorCalldata)
		if err != nil {
			if !errors.Is(err, abi.ErrNoLenField) {
				return tx, nil, err
			}
		}
	}

	parser.Cache.SetAbiByAddress(tx.Class, tx.Contract.Hash)

	var proxyId uint64
	if tx.Class.Type.Is(storage.ClassTypeProxy) {
		proxyId = tx.ContractID
	}

	txCtx := parserData.NewTxContextFromDeploy(tx, proxyId)

	if trace.FunctionInvocation != nil {
		tx.Events, err = parseEvents(ctx, parser.EventParser, txCtx, classAbi, trace.FunctionInvocation.Events)
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
