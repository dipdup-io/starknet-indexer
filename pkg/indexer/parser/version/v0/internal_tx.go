package v0

import (
	"bytes"
	"context"
	"errors"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/decode"
	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/interfaces"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/resolver"
	"github.com/rs/zerolog/log"
)

// InternalTxParser -
type InternalTxParser struct {
	Resolver resolver.Resolver
	Cache    *cache.Cache
	blocks   storage.IBlock

	EventParser    interfaces.EventParser
	MessageParser  interfaces.MessageParser
	TransferParser interfaces.TransferParser
	ProxyUpgrader  interfaces.ProxyUpgrader
}

// NewInternalTxParser -
func NewInternalTxParser(
	resolver resolver.Resolver,
	cache *cache.Cache,
	blocks storage.IBlock,
	eventParser interfaces.EventParser,
	messageParser interfaces.MessageParser,
	transferParser interfaces.TransferParser,
	proxyUpgrader interfaces.ProxyUpgrader,
) InternalTxParser {
	return InternalTxParser{
		Resolver:       resolver,
		Cache:          cache,
		blocks:         blocks,
		EventParser:    eventParser,
		MessageParser:  messageParser,
		TransferParser: transferParser,
		ProxyUpgrader:  proxyUpgrader,
	}
}

// Parse -
func (parser InternalTxParser) Parse(ctx context.Context, txCtx parserData.TxContext, internal sequencer.Invocation) (storage.Internal, error) {
	tx := storage.Internal{
		ID:             parser.Resolver.NextTxId(),
		Height:         txCtx.Height,
		Time:           txCtx.Time,
		Hash:           txCtx.Hash,
		Status:         txCtx.Status,
		CallType:       storage.NewCallType(internal.CallType),
		EntrypointType: storage.NewEntrypointType(internal.EntrypointType),
		Selector:       internal.Selector.Bytes(),
		Result:         internal.Result,
		Calldata:       internal.Calldata,

		DeclareID:       txCtx.DeclareID,
		DeployID:        txCtx.DeployID,
		DeployAccountID: txCtx.DeployAccountID,
		InvokeID:        txCtx.InvokeID,
		L1HandlerID:     txCtx.L1HandlerID,
		InternalID:      txCtx.InternalID,

		Events:    make([]storage.Event, 0),
		Messages:  make([]storage.Message, 0),
		Internals: make([]storage.Internal, 0),
	}

	if class, err := parser.Resolver.FindClassByHash(ctx, internal.ClassHash, tx.Height); err != nil {
		return tx, err
	} else if class != nil {
		tx.Class = *class
		tx.ClassID = class.ID
	}

	if address, err := parser.Resolver.FindAddressByHash(ctx, internal.ContractAddress); err != nil {
		return tx, err
	} else if address != nil {
		tx.ContractID = address.ID
		tx.Contract = *address
		tx.Contract.Height = tx.Height
		if tx.ClassID != 0 {
			tx.Contract.ClassID = &tx.Class.ID
			tx.Contract.Class = tx.Class
		} else if tx.Contract.ClassID != nil {
			tx.ClassID = *tx.Contract.ClassID
		}
	}

	if address, err := parser.Resolver.FindAddressByHash(ctx, internal.CallerAddress); err != nil {
		return tx, err
	} else if address != nil {
		tx.CallerID = address.ID
		tx.Caller = *address
		tx.Caller.Height = tx.Height
	}

	var (
		contractAbi abi.Abi
		err         error
		proxyId     uint64
	)

	switch {
	case len(tx.Class.Hash) > 0:
		contractAbi, err = parser.Cache.GetAbiByClassHash(ctx, tx.Class.Hash)
	case len(tx.Contract.Hash) > 0:
		contractAbi, err = parser.Cache.GetAbiByAddress(ctx, tx.Contract.Hash)
	}
	if err != nil {
		if !parser.blocks.IsNoRows(err) {
			return tx, err
		}

		if err := parser.Resolver.ReceiveClass(ctx, &tx.Class); err != nil {
			return tx, err
		}
		tx.Class.Height = tx.Height
	}

	if len(tx.Selector) > 0 {
		var (
			_, has = contractAbi.GetByTypeAndSelector(internal.EntrypointType, encoding.EncodeHex(tx.Selector))

			isExecute       = bytes.Equal(tx.Selector, encoding.ExecuteEntrypointSelector)
			isChangeModules = bytes.Equal(tx.Selector, encoding.ChangeModuleEntrypointSelector)
			isUnknownProxy  = false
		)

		if !(isExecute || isChangeModules) && !has {
			if tx.Class.ID == 0 {
				class, err := parser.Cache.GetClassById(ctx, *tx.Contract.ClassID)
				if err != nil {
					return tx, err
				}
				tx.Class = *class
			}
			if !has {
				contractAbi, err = parser.Resolver.Proxy(ctx, txCtx, tx.Class, tx.Contract, tx.Selector)
				if err != nil {
					isUnknownProxy = errors.Is(err, resolver.ErrUnknownProxy)
					if !isUnknownProxy {
						return tx, err
					}
					log.Warn().Hex("contract", tx.Contract.Hash).Msg("unknown proxy")
				}
				if tx.Class.Type.Is(storage.ClassTypeProxy) {
					proxyId = tx.ContractID
				}
			}
		}

		if len(internal.Calldata) > 0 && !isUnknownProxy {
			switch {
			case isExecute && !has:
				tx.Entrypoint = encoding.ExecuteEntrypoint
				tx.ParsedCalldata, err = abi.DecodeExecuteCallData(internal.Calldata)
			case isChangeModules && !has:
				tx.Entrypoint = encoding.ChangeModulesEntrypoint
				tx.ParsedCalldata, err = abi.DecodeChangeModulesCallData(internal.Calldata)
			default:
				tx.ParsedCalldata, tx.Entrypoint, err = decode.InternalCalldata(contractAbi, tx.Selector, internal.Calldata, tx.EntrypointType)
			}

			if err != nil {
				switch {
				case (errors.Is(err, abi.ErrNoLenField) || errors.Is(err, abi.ErrTooShortCallData)):
				case errors.Is(err, decode.ErrUnknownSelector):
					log.Err(err).Hex("tx_hash", tx.Hash).Msg("can't decode calldata")
				default:
					return tx, err
				}
			}
		}

		if len(internal.Result) > 0 && !isUnknownProxy {
			switch {
			case isExecute && !has:
			case isChangeModules && !has:
				tx.ParsedResult, err = abi.DecodeChangeModulesResult(internal.Result)
			default:
				tx.ParsedResult, err = decode.Result(contractAbi, internal.Result, tx.Selector, tx.EntrypointType)
			}
			if err != nil {
				switch {
				case (errors.Is(err, abi.ErrNoLenField) || errors.Is(err, abi.ErrTooShortCallData)):
				case errors.Is(err, decode.ErrUnknownSelector):
					log.Err(err).Hex("tx_hash", tx.Hash).Msg("can't decode result")
				default:
					return tx, err
				}
			}
		}
	}

	internalTxCtx := parserData.NewTxContextFromInternal(tx, txCtx.ProxyUpgrades, proxyId)

	tx.Events, err = parseEvents(ctx, parser.EventParser, internalTxCtx, contractAbi, internal.Events)
	if err != nil {
		return tx, err
	}

	if err := parser.ProxyUpgrader.Parse(ctx, internalTxCtx, tx.Contract, tx.Events, tx.Entrypoint, tx.ParsedCalldata); err != nil {
		return tx, err
	}

	tx.Messages, err = parseMessages(ctx, parser.MessageParser, internalTxCtx, internal.Messages)
	if err != nil {
		return tx, err
	}

	tx.Internals, err = parseInternals(ctx, parser, internalTxCtx, internal.InternalCalls)
	if err != nil {
		return tx, err
	}

	isNew := len(tx.Internals) > 0 && isInternalNotEqualParent(internalTxCtx, tx.Internals[0])

	var countInternalTransfers int
	for i := range tx.Internals {
		countInternalTransfers += len(tx.Internals[i].Transfers)
	}

	if isNew || countInternalTransfers == 0 {
		tx.Transfers, err = parser.TransferParser.ParseEvents(ctx, txCtx, tx.Contract, tx.Events)
		if err != nil {
			return tx, err
		}
		if len(tx.Transfers) == 0 {
			tx.Transfers, err = parser.TransferParser.ParseCalldata(ctx, internalTxCtx, tx.Entrypoint, tx.ParsedCalldata)
			if err != nil {
				return tx, err
			}
		}
	}

	return tx, nil
}
