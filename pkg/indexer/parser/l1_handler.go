package parser

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/decode"
)

func (parser *Parser) getL1Handler(ctx context.Context, raw *data.L1Handler, block storage.Block, trace sequencer.Trace) (storage.L1Handler, error) {
	tx := storage.L1Handler{
		Height:             block.Height,
		Time:               block.Time,
		Status:             block.Status,
		Hash:               encoding.MustDecodeHex(trace.TransactionHash),
		EntrypointSelector: encoding.MustDecodeHex(raw.EntrypointSelector),
		CallData:           raw.Calldata,
		Nonce:              decimalFromHex(raw.Nonce),
	}

	if raw.ContractAddress != "" {
		tx.Contract = storage.Address{
			Hash: encoding.MustDecodeHex(raw.ContractAddress),
		}

		if err := parser.findAddress(ctx, &tx.Contract); err != nil {
			return tx, err
		}
		tx.ContractID = tx.Contract.ID
	}

	contractAbi, err := parser.cache.GetAbiByAddress(ctx, tx.Contract.Hash)
	if err != nil {
		return tx, err
	}

	if len(tx.CallData) > 0 {
		parsed, entrypoint, err := decode.CalldataForL1Handler(parser.cache, contractAbi, tx.EntrypointSelector, tx.CallData)
		if err != nil {
			return tx, err
		}
		tx.ParsedCalldata = parsed
		tx.Entrypoint = entrypoint
	}

	if trace.FunctionInvocation != nil {
		for i := range trace.FunctionInvocation.Events {
			event, err := parser.getEvent(block, contractAbi, trace.FunctionInvocation.Events[i])
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
