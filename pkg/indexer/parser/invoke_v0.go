package parser

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/decode"
	"github.com/pkg/errors"
)

func (parser *Parser) getInvokeV0(ctx context.Context, raw *data.InvokeV0, block storage.Block, trace sequencer.Trace) (storage.InvokeV0, error) {
	tx := storage.InvokeV0{
		Height:             block.Height,
		Time:               block.Time,
		Status:             block.Status,
		Hash:               encoding.MustDecodeHex(trace.TransactionHash),
		EntrypointSelector: encoding.MustDecodeHex(raw.EntrypointSelector),
		Signature:          raw.Signature,
		CallData:           raw.Calldata,
		MaxFee:             decimalFromHex(raw.MaxFee),
		Nonce:              decimalFromHex(raw.Nonce),

		Events:    make([]storage.Event, 0),
		Messages:  make([]storage.Message, 0),
		Internals: make([]storage.Internal, 0),
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

	var contractAbi abi.Abi
	var err error

	if len(tx.CallData) > 0 || len(trace.FunctionInvocation.Events) > 0 {
		contractAbi, err = parser.cache.GetAbiByAddress(ctx, tx.Contract.Hash)
		if err != nil {
			return tx, err
		}

		if _, ok := contractAbi.GetFunctionBySelector(encoding.EncodeHex(tx.EntrypointSelector)); !ok {
			diff, err := parser.storageDiffs.GetOnBlock(ctx, block.Height, tx.ContractID, encoding.MustDecodeHex("0xF920571B9F85BDD92A867CFDC73319D0F8836F0E69E06E4C5566B6203F75CC"))
			if err != nil {
				return tx, err
			}
			contractAbi, err = parser.cache.GetAbiByAddress(ctx, diff.Value)
			if err != nil {
				return tx, err
			}
		}
	}

	if len(tx.CallData) > 0 {
		decoded, entrypoint, err := decode.CalldataBySelector(parser.cache, contractAbi, tx.EntrypointSelector, tx.CallData)
		if err != nil {
			return tx, err
		}

		tx.ParsedCalldata = decoded
		tx.Entrypoint = entrypoint
	}

	if trace.FunctionInvocation != nil {
		for i := range trace.FunctionInvocation.Events {
			event, err := parser.getEvent(block, contractAbi, trace.FunctionInvocation.Events[i])
			if err != nil {
				return tx, errors.Wrap(err, "event")
			}
			tx.Events = append(tx.Events, event)
		}
		for i := range trace.FunctionInvocation.Messages {
			msg, err := parser.getMessage(ctx, block, trace.FunctionInvocation.Messages[i])
			if err != nil {
				return tx, errors.Wrap(err, "message")
			}
			tx.Messages = append(tx.Messages, msg)
		}
		for i := range trace.FunctionInvocation.InternalCalls {
			internal, err := parser.getInternal(ctx, block, trace.FunctionInvocation.InternalCalls[i], tx.Hash, tx.Status)
			if err != nil {
				return tx, errors.Wrap(err, "internal")
			}
			tx.Internals = append(tx.Internals, internal)
		}
	}

	return tx, nil
}
