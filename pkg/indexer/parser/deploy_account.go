package parser

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/decode"
)

func (parser *Parser) getDeployAccount(ctx context.Context, raw *data.DeployAccount, block storage.Block, trace sequencer.Trace) (storage.DeployAccount, error) {
	tx := storage.DeployAccount{
		Height:              block.Height,
		Time:                block.Time,
		Status:              block.Status,
		Hash:                encoding.MustDecodeHex(trace.TransactionHash),
		ContractAddressSalt: encoding.MustDecodeHex(raw.ContractAddressSalt),
		ConstructorCalldata: raw.ConstructorCalldata,
		MaxFee:              decimalFromHex(raw.MaxFee),
		Nonce:               decimalFromHex(raw.Nonce),
		Signature:           raw.Signature,
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

	classAbi, err := parser.cache.GetAbiByAddress(ctx, tx.Contract.Hash)
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
