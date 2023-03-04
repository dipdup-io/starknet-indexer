package parser

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/decode"
)

func (parser *Parser) getInternal(ctx context.Context, block storage.Block, internal sequencer.Invocation, hash []byte, status storage.Status) (storage.Internal, error) {
	tx := storage.Internal{
		Height:         block.Height,
		Time:           block.Time,
		Hash:           hash,
		Status:         status,
		CallType:       storage.NewCallType(internal.CallType),
		EntrypointType: storage.NewEntrypointType(internal.EntrypointType),
		Selector:       encoding.MustDecodeHex(internal.Selector),
		Result:         internal.Result,
		Calldata:       internal.Calldata,

		Events:    make([]storage.Event, 0),
		Messages:  make([]storage.Message, 0),
		Internals: make([]storage.Internal, 0),
	}

	if internal.ClassHash != "" {
		tx.Class = storage.Class{
			Hash: encoding.MustDecodeHex(internal.ClassHash),
		}

		if err := parser.findClass(ctx, &tx.Class); err != nil {
			return tx, err
		}
		tx.ClassID = tx.Class.ID
	}

	if internal.ContractAddress != "" {
		tx.Contract = storage.Address{
			Hash: encoding.MustDecodeHex(internal.ContractAddress),
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

	if internal.CallerAddress != "" {
		tx.Caller = storage.Address{
			Hash: encoding.MustDecodeHex(internal.CallerAddress),
		}

		if err := parser.findAddress(ctx, &tx.Caller); err != nil {
			return tx, err
		}
		tx.CallerID = tx.Caller.ID
	}

	var contractAbi abi.Abi
	var err error

	switch {
	case len(tx.Class.Hash) > 0:
		contractAbi, err = parser.cache.GetAbiByClassHash(ctx, tx.Class.Hash)
	case len(tx.Contract.Hash) > 0:
		contractAbi, err = parser.cache.GetAbiByAddress(ctx, tx.Contract.Hash)
	}
	if err != nil {
		return tx, err
	}

	if _, ok := contractAbi.GetFunctionBySelector(encoding.EncodeHex(tx.Selector)); !ok {
		diff, err := parser.storageDiffs.GetOnBlock(ctx, block.Height, tx.ContractID, encoding.MustDecodeHex("0xF920571B9F85BDD92A867CFDC73319D0F8836F0E69E06E4C5566B6203F75CC"))
		if err != nil {
			return tx, err
		}
		contractAbi, err = parser.cache.GetAbiByAddress(ctx, diff.Value)
		if err != nil {
			return tx, err
		}
	}

	if len(internal.Calldata) > 0 && len(contractAbi.Functions) > 0 && len(tx.Selector) > 0 {
		parsed, entrypoint, err := decode.CalldataBySelector(parser.cache, contractAbi, tx.Selector, tx.Calldata)
		if err != nil {
			return tx, err
		}
		tx.ParsedCalldata = parsed
		tx.Entrypoint = entrypoint
	}

	for i := range internal.Events {
		event, err := parser.getEvent(block, contractAbi, internal.Events[i])
		if err != nil {
			return tx, err
		}
		tx.Events = append(tx.Events, event)
	}

	for i := range internal.Messages {
		msg, err := parser.getMessage(ctx, block, internal.Messages[i])
		if err != nil {
			return tx, err
		}
		tx.Messages = append(tx.Messages, msg)
	}

	for i := range internal.InternalCalls {
		internal, err := parser.getInternal(ctx, block, internal.InternalCalls[i], tx.Hash, tx.Status)
		if err != nil {
			return tx, err
		}
		tx.Internals = append(tx.Internals, internal)
	}

	return tx, nil
}
