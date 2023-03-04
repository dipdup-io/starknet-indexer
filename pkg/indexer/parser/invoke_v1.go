package parser

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
)

func (parser *Parser) getInvokeV1(ctx context.Context, raw *data.InvokeV1, block storage.Block, trace sequencer.Trace) (storage.InvokeV1, error) {
	tx := storage.InvokeV1{
		Height:    block.Height,
		Time:      block.Time,
		Status:    block.Status,
		Hash:      encoding.MustDecodeHex(trace.TransactionHash),
		Signature: raw.Signature,
		CallData:  raw.Calldata,
		MaxFee:    decimalFromHex(raw.MaxFee),
		Nonce:     decimalFromHex(raw.Nonce),
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

	if raw.SenderAddress != "" {
		tx.Sender = storage.Address{
			Hash: encoding.MustDecodeHex(raw.SenderAddress),
		}

		if err := parser.findAddress(ctx, &tx.Sender); err != nil {
			return tx, err
		}
		tx.SenderID = tx.Sender.ID
	}

	if len(tx.CallData) > 0 {
		parsed, err := abi.DecodeFunctionCallData(tx.CallData, abi.ExecuteFunction, map[string]*abi.StructItem{
			"CallArray": &abi.CallArray,
		})
		if err != nil {
			return tx, err
		}
		tx.ParsedCalldata = parsed
	}

	if trace.FunctionInvocation != nil {
		if len(trace.FunctionInvocation.Events) > 0 {
			contractAbi, err := parser.cache.GetAbiByAddress(ctx, tx.Contract.Hash)
			if err != nil {
				return tx, err
			}
			for i := range trace.FunctionInvocation.Events {
				event, err := parser.getEvent(block, contractAbi, trace.FunctionInvocation.Events[i])
				if err != nil {
					return tx, err
				}
				tx.Events = append(tx.Events, event)
			}
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
