package parser

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/decode"
)

func (parser *Parser) getDeploy(ctx context.Context, raw *data.Deploy, block storage.Block, trace sequencer.Trace) (storage.Deploy, error) {
	tx := storage.Deploy{
		Height:              block.Height,
		Time:                block.Time,
		Status:              block.Status,
		Hash:                encoding.MustDecodeHex(trace.TransactionHash),
		ContractAddressSalt: encoding.MustDecodeHex(raw.ContractAddressSalt),
		ConstructorCalldata: raw.ConstructorCalldata,
	}

	if raw.ClassHash != "" {
		tx.Class = storage.Class{
			Hash: encoding.MustDecodeHex(raw.ClassHash),
		}

		if err := parser.findClass(ctx, &tx.Class); err != nil {
			return tx, err
		}
		tx.ClassID = tx.Class.ID
	}

	if raw.ContractAddress != "" {
		tx.Contract = storage.Address{
			Hash: encoding.MustDecodeHex(raw.ContractAddress),
		}

		if err := parser.findAddress(ctx, &tx.Contract); err != nil {
			return tx, err
		}
		tx.ContractID = tx.Contract.ID

		if tx.ClassID != 0 {
			tx.Contract.ClassID = &tx.ClassID
			tx.Contract.Class = tx.Class
		}
	}

	classAbi, err := parser.cache.GetAbiByClassHash(ctx, tx.Class.Hash)
	if err != nil {
		return tx, err
	}

	if len(tx.ConstructorCalldata) > 0 {
		parsedCalldata, err := decode.CalldataForConstructor(parser.cache, classAbi, tx.ConstructorCalldata)
		if err != nil {
			return tx, err
		}
		tx.ParsedCalldata = parsedCalldata
	}

	parser.cache.SetAbiByAddress(tx.Class, tx.Contract.Hash)

	if trace.FunctionInvocation != nil {
		for i := range trace.FunctionInvocation.Events {
			event, err := parser.getEvent(block, classAbi, trace.FunctionInvocation.Events[i])
			if err != nil {
				return tx, err
			}
			tx.Events = append(tx.Events, event)
		}
		for i := range trace.FunctionInvocation.Messages {
			msg, err := parser.getMessage(ctx, block, trace.FunctionInvocation.Messages[i])
			if err != nil {
				return tx, err
			}
			tx.Messages = append(tx.Messages, msg)
		}
		for i := range trace.FunctionInvocation.InternalCalls {
			internal, err := parser.getInternal(ctx, block, trace.FunctionInvocation.InternalCalls[i], tx.Hash, tx.Status)
			if err != nil {
				return tx, err
			}
			tx.Internals = append(tx.Internals, internal)
		}
	}

	return tx, nil
}
