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
)

// InternalTxParser -
type InternalTxParser struct {
	Resolver resolver.Resolver
	Cache    *cache.Cache
	blocks   storage.IBlock

	EventParser    interfaces.EventParser
	MessageParser  interfaces.MessageParser
	TransferParser interfaces.TransferParser
	TokenParser    interfaces.TokenParser
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
	tokenParser interfaces.TokenParser,
	proxyUpgrader interfaces.ProxyUpgrader,
) InternalTxParser {
	return InternalTxParser{
		Resolver:       resolver,
		Cache:          cache,
		blocks:         blocks,
		EventParser:    eventParser,
		MessageParser:  messageParser,
		TransferParser: transferParser,
		TokenParser:    tokenParser,
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

	if class, err := parser.Resolver.FindClassByHash(ctx, internal.ClassHash); err != nil {
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

	isExecute := bytes.Equal(tx.Selector, encoding.ExecuteEntrypointSelector)
	_, hasExecute := contractAbi.Functions[encoding.ExecuteEntrypoint]

	isChangeModules := bytes.Equal(tx.Selector, encoding.ChangeModuleEntrypointSelector)
	_, hasChangeModules := contractAbi.Functions[encoding.ChangeModulesEntrypoint]

	if len(tx.Selector) > 0 && !isExecute && !isChangeModules {
		if _, has := contractAbi.GetByTypeAndSelector(internal.EntrypointType, encoding.EncodeHex(tx.Selector)); !has {
			if tx.Class.ID == 0 {
				class, err := parser.Cache.GetClassById(ctx, *tx.Contract.ClassID)
				if err != nil {
					return tx, err
				}
				tx.Class = *class
			}
			contractAbi, err = parser.Resolver.Proxy(ctx, txCtx, tx.Class, tx.Contract)
			if err != nil {
				return tx, err
			}
			if tx.Class.Type.Is(storage.ClassTypeProxy) {
				proxyId = tx.ContractID
			}
		}
	}

	if len(internal.Calldata) > 0 && len(tx.Selector) > 0 {
		switch {
		case isExecute && !hasExecute:
			tx.Entrypoint = encoding.ExecuteEntrypoint
			tx.ParsedCalldata, err = abi.DecodeExecuteCallData(internal.Calldata)
		case isChangeModules && !hasChangeModules:
			tx.Entrypoint = encoding.ChangeModulesEntrypoint
			tx.ParsedCalldata, err = abi.DecodeChangeModulesCallData(internal.Calldata)
		default:
			switch tx.EntrypointType {
			case storage.EntrypointTypeExternal:
				tx.ParsedCalldata, tx.Entrypoint, err = decode.CalldataBySelector(contractAbi, tx.Selector, tx.Calldata)
			case storage.EntrypointTypeConstructor:
				tx.ParsedCalldata, err = decode.CalldataForConstructor(contractAbi, tx.Calldata)
			case storage.EntrypointTypeL1Handler:
				tx.ParsedCalldata, tx.Entrypoint, err = decode.CalldataForL1Handler(contractAbi, tx.Selector, tx.Calldata)
			}
		}

		if err != nil {
			if !errors.Is(err, abi.ErrNoLenField) {
				return tx, err
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

	if tx.EntrypointType == storage.EntrypointTypeConstructor && tx.Class.Type.OneOf(storage.ClassTypeERC20, storage.ClassTypeERC721, storage.ClassTypeERC1155) {
		token, err := parser.TokenParser.Parse(ctx, txCtx, tx.Contract, tx.Class.Type, tx.ParsedCalldata)
		if err != nil {
			return tx, err
		}
		tx.ERC20 = token.ERC20
		tx.ERC721 = token.ERC721
		tx.ERC1155 = token.ERC1155
	}

	return tx, nil
}
