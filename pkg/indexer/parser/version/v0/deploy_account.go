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

// ParseDeployAccount -
func (parser Parser) ParseDeployAccount(ctx context.Context, raw *data.DeployAccount, block storage.Block, trace sequencer.Trace, receipts sequencer.Receipt) (storage.DeployAccount, error) {
	tx := storage.DeployAccount{
		ID:                  parser.Resolver.NextTxId(),
		Height:              block.Height,
		Time:                block.Time,
		Status:              block.Status,
		Hash:                trace.TransactionHash.Bytes(),
		ContractAddressSalt: encoding.MustDecodeHex(raw.ContractAddressSalt),
		ConstructorCalldata: raw.ConstructorCalldata,
		MaxFee:              raw.MaxFee.Decimal(),
		Nonce:               raw.Nonce.Decimal(),
		Signature:           raw.Signature,
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
		tx.ContractID = address.ID
		tx.Contract = *address
		tx.Contract.Height = tx.Height
		if tx.ClassID != 0 {
			tx.Contract.ClassID = &tx.Class.ID
			tx.Contract.Class = tx.Class
		}
	}

	var (
		classAbi abi.Abi
		err      error
	)

	if err = json.Unmarshal(tx.Class.Abi, &classAbi); err != nil {
		return tx, err
	}

	if len(tx.ConstructorCalldata) > 0 {
		tx.ParsedCalldata, err = decode.CalldataForConstructor(parser.Cache, classAbi, tx.ConstructorCalldata)
		if err != nil {
			if !errors.Is(err, abi.ErrNoLenField) {
				return tx, err
			}
		}
	}

	var proxyId uint64
	if tx.Class.Type.Is(storage.ClassTypeProxy) {
		proxyId = tx.ContractID
	}

	txCtx := parserData.NewTxContextFromDeployAccount(tx, proxyId)

	if trace.FunctionInvocation != nil {
		tx.Events, err = parseEvents(ctx, parser.EventParser, txCtx, classAbi, trace.FunctionInvocation.Events)
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
		fee.DeployAccountID = &tx.ID
		tx.Fee = fee
	}

	return tx, nil
}
