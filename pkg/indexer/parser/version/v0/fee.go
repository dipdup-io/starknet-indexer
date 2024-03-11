package v0

import (
	"bytes"
	"context"
	"errors"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	starknetData "github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/decode"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/interfaces"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/resolver"
	"github.com/shopspring/decimal"
)

const (
	actualFeeContractHash = "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
	actualFeeToHash       = "0x05dcd266a80b8a5f29f04d779c6b166b80150c24f2180a75e82427242dab20a9"
)

// FeeParser -
type FeeParser struct {
	cache    *cache.Cache
	resolver resolver.Resolver
	blocks   storage.IBlock

	EventParser      interfaces.EventParser
	MessageParser    interfaces.MessageParser
	TransferParser   interfaces.TransferParser
	InternalTxParser interfaces.InternalTxParser

	actualFeeContractId uint64
	actualFeeToId       uint64
}

// NewFeeParser -
func NewFeeParser(
	cache *cache.Cache,
	resolver resolver.Resolver,
	blocks storage.IBlock,
	eventParser interfaces.EventParser,
	messageParser interfaces.MessageParser,
	transferParser interfaces.TransferParser,
	internalTxParser interfaces.InternalTxParser,
) FeeParser {
	return FeeParser{
		cache:            cache,
		resolver:         resolver,
		blocks:           blocks,
		EventParser:      eventParser,
		MessageParser:    messageParser,
		TransferParser:   transferParser,
		InternalTxParser: internalTxParser,
	}
}

// ParseInvocation -
func (parser FeeParser) ParseActualFee(ctx context.Context, txCtx data.TxContext, actualFee starknetData.Felt) (*storage.Transfer, error) {
	fee := actualFee.Decimal()
	if fee.IsZero() {
		return nil, nil
	}

	if parser.actualFeeContractId == 0 {
		address := storage.Address{
			Hash: starknetData.Felt(actualFeeContractHash).Bytes(),
		}
		if err := parser.resolver.FindAddress(ctx, &address); err != nil {
			return nil, err
		}
		parser.actualFeeContractId = address.ID
	}

	if parser.actualFeeToId == 0 {
		address := storage.Address{
			Hash: starknetData.Felt(actualFeeToHash).Bytes(),
		}
		if err := parser.resolver.FindAddress(ctx, &address); err != nil {
			return nil, err
		}
		parser.actualFeeToId = address.ID
	}

	return &storage.Transfer{
		Height:     txCtx.Height,
		Time:       txCtx.Time,
		Amount:     fee,
		FromID:     txCtx.ContractId,
		ToID:       parser.actualFeeToId,
		ContractID: parser.actualFeeContractId,
		TokenID:    decimal.Zero,

		DeclareID:       txCtx.DeclareID,
		DeployID:        txCtx.DeployID,
		DeployAccountID: txCtx.DeployAccountID,
		InvokeID:        txCtx.InvokeID,
		L1HandlerID:     txCtx.L1HandlerID,
		Token: storage.Token{
			TokenId:     decimal.Zero,
			ContractId:  parser.actualFeeContractId,
			Type:        storage.TokenTypeERC20,
			FirstHeight: txCtx.Height,
		},
	}, nil
}

// ParseInvocation -
func (parser FeeParser) ParseInvocation(ctx context.Context, txCtx data.TxContext, feeInvocation sequencer.Invocation) (*storage.Fee, error) {
	tx := storage.Fee{
		ID:             parser.resolver.NextTxId(),
		Height:         txCtx.Height,
		Time:           txCtx.Time,
		Status:         txCtx.Status,
		CallType:       storage.NewCallType(feeInvocation.CallType),
		EntrypointType: storage.NewEntrypointType(feeInvocation.EntrypointType),
		Selector:       feeInvocation.Selector.Bytes(),
		Result:         make([]string, len(feeInvocation.Result)),
		Calldata:       make([]string, len(feeInvocation.Calldata)),

		DeclareID:       txCtx.DeclareID,
		DeployID:        txCtx.DeployID,
		DeployAccountID: txCtx.DeployAccountID,
		InvokeID:        txCtx.InvokeID,
		L1HandlerID:     txCtx.L1HandlerID,

		Events:    make([]storage.Event, 0),
		Messages:  make([]storage.Message, 0),
		Internals: make([]storage.Internal, 0),
	}
	for i := 0; i < len(feeInvocation.Calldata); i++ {
		tx.Calldata[i] = feeInvocation.Calldata[i].String()
	}
	for i := 0; i < len(feeInvocation.Result); i++ {
		tx.Result[i] = feeInvocation.Result[i].String()
	}

	if class, err := parser.resolver.FindClassByHash(ctx, feeInvocation.ClassHash, tx.Height); err != nil {
		return nil, err
	} else if class != nil {
		tx.Class = *class
		tx.ClassID = class.ID
	}

	if address, err := parser.resolver.FindAddressByHash(ctx, feeInvocation.ContractAddress); err != nil {
		return nil, err
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

	if address, err := parser.resolver.FindAddressByHash(ctx, feeInvocation.CallerAddress); err != nil {
		return nil, err
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
		contractAbi, err = parser.cache.GetAbiByClassHash(ctx, tx.Class.Hash)
	case len(tx.Contract.Hash) > 0:
		contractAbi, err = parser.cache.GetAbiByAddress(ctx, tx.Contract.Hash)
	}
	if err != nil {
		if !parser.blocks.IsNoRows(err) {
			return nil, err
		}

		if err := parser.resolver.ReceiveClass(ctx, &tx.Class); err != nil {
			return nil, err
		}
		tx.Class.Height = tx.Height
	}

	isExecute := bytes.Equal(tx.Selector, encoding.ExecuteEntrypointSelector)
	_, hasExecute := contractAbi.Functions[encoding.ExecuteEntrypoint]

	if len(tx.Selector) > 0 && !isExecute {
		if _, has := contractAbi.GetByTypeAndSelector(feeInvocation.EntrypointType, encoding.EncodeHex(tx.Selector)); !has {
			if tx.Class.ID == 0 {
				class, err := parser.cache.GetClassById(ctx, *tx.Contract.ClassID)
				if err != nil {
					return nil, err
				}
				tx.Class = *class
			}
			contractAbi, err = parser.resolver.Proxy(ctx, txCtx, tx.Class, tx.Contract, tx.Selector)
			if err != nil {
				return nil, err
			}
			if tx.Class.Type.Is(storage.ClassTypeProxy) {
				proxyId = tx.ContractID
			}
		}
	}

	if len(tx.Calldata) > 0 && len(tx.Selector) > 0 {
		if isExecute && !hasExecute {
			tx.Entrypoint = encoding.ExecuteEntrypoint
			tx.ParsedCalldata, err = abi.DecodeExecuteCallData(tx.Calldata)
		} else {
			tx.ParsedCalldata, tx.Entrypoint, err = decode.CalldataBySelector(contractAbi, tx.Selector, tx.Calldata)
		}
		if err != nil {
			if !errors.Is(err, abi.ErrNoLenField) {
				return nil, err
			}
		}
	}

	feeTxCtx := data.NewTxContextFromFee(tx, proxyId)

	tx.Events, err = parseEvents(ctx, parser.EventParser, feeTxCtx, contractAbi, feeInvocation.Events)
	if err != nil {
		return nil, err
	}

	tx.Messages, err = parseMessages(ctx, parser.MessageParser, feeTxCtx, feeInvocation.Messages)
	if err != nil {
		return nil, err
	}

	tx.Internals, err = parseInternals(ctx, parser.InternalTxParser, feeTxCtx, feeInvocation.InternalCalls)
	if err != nil {
		return nil, err
	}

	isNew := len(tx.Internals) > 0 && isInternalNotEqualParent(feeTxCtx, tx.Internals[0])

	var countInternalTransfers int
	for i := range tx.Internals {
		countInternalTransfers += len(tx.Internals[i].Transfers)
	}

	if isNew || countInternalTransfers == 0 {
		tx.Transfers, err = parser.TransferParser.ParseEvents(ctx, txCtx, tx.Contract, tx.Events)
		if err != nil {
			return nil, err
		}
		if len(tx.Transfers) == 0 {
			tx.Transfers, err = parser.TransferParser.ParseCalldata(ctx, feeTxCtx, tx.Entrypoint, tx.ParsedCalldata)
			if err != nil {
				return nil, err
			}
		}
	}

	return &tx, nil
}
