package v0

import (
	"bytes"
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/interfaces"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/resolver"
	"github.com/pkg/errors"
)

// Parser -
type Parser struct {
	Resolver resolver.Resolver
	Cache    *cache.Cache

	InternalTxParser interfaces.InternalTxParser
	EventParser      interfaces.EventParser
	MessageParser    interfaces.MessageParser
	FeeParser        interfaces.FeeParser
	TokenParser      interfaces.TokenParser
	TransferParser   interfaces.TransferParser
	ProxyUpgrader    interfaces.ProxyUpgrader
}

// NewWithParsers -
func NewWithParsers(
	resolver resolver.Resolver,
	cache *cache.Cache,
	blocks storage.IBlock,

	internalTxParser interfaces.InternalTxParser,
	eventParser interfaces.EventParser,
	messageParser interfaces.MessageParser,
	tokenParser interfaces.TokenParser,
	transferParser interfaces.TransferParser,
	proxyUpgrader interfaces.ProxyUpgrader,
) Parser {
	return Parser{
		Cache:    cache,
		Resolver: resolver,

		InternalTxParser: internalTxParser,
		EventParser:      eventParser,
		MessageParser:    messageParser,
		TokenParser:      tokenParser,
		TransferParser:   transferParser,
		ProxyUpgrader:    proxyUpgrader,
	}
}

// New -
func New(
	resolver resolver.Resolver,
	cache *cache.Cache,
	blocks storage.IBlock,
) Parser {
	event := NewEventParser(cache, resolver)
	message := NewMessageParser(cache, resolver)
	transferParser := NewTransferParser(resolver)
	proxyUpgrader := NewProxyUpgrader(resolver)
	tokenParser := NewTokenParser(cache, resolver)
	internal := NewInternalTxParser(resolver, cache, blocks, event, message, transferParser, tokenParser, proxyUpgrader)
	fee := NewFeeParser(cache, resolver, blocks, event, message, transferParser, internal)

	return Parser{
		Cache:            cache,
		Resolver:         resolver,
		InternalTxParser: internal,
		EventParser:      event,
		MessageParser:    message,
		FeeParser:        fee,
		TokenParser:      tokenParser,
		TransferParser:   transferParser,
		ProxyUpgrader:    proxyUpgrader,
	}
}

func parseEvents(ctx context.Context, eventParser interfaces.EventParser, txCtx parserData.TxContext, contractAbi abi.Abi, events []data.Event) ([]storage.Event, error) {
	result := make([]storage.Event, 0)
	for i := range events {
		event, err := eventParser.Parse(ctx, txCtx, contractAbi, events[i])
		if err != nil {
			return nil, errors.Wrap(err, "event")
		}
		result = append(result, event)
	}
	return result, nil
}

func parseMessages(ctx context.Context, msgParser interfaces.MessageParser, txCtx parserData.TxContext, msgs []data.Message) ([]storage.Message, error) {
	result := make([]storage.Message, 0)
	for i := range msgs {
		msg, err := msgParser.Parse(ctx, txCtx, msgs[i])
		if err != nil {
			return nil, errors.Wrap(err, "message")
		}
		result = append(result, msg)
	}
	return result, nil
}

func parseInternals(ctx context.Context, internalParser interfaces.InternalTxParser, txCtx parserData.TxContext, txs []sequencer.Invocation) ([]storage.Internal, error) {
	result := make([]storage.Internal, 0)
	for i := range txs {
		tx, err := internalParser.Parse(ctx, txCtx, txs[i])
		if err != nil {
			return nil, errors.Wrap(err, "internal")
		}
		result = append(result, tx)
	}
	return result, nil
}

func isInternalNotEqualParent(txCtx parserData.TxContext, tx storage.Internal) bool {
	switch {
	case txCtx.Internal != nil:
		if tx.ContractID != txCtx.Internal.ContractID {
			return true
		}
		if tx.CallerID != txCtx.Internal.CallerID {
			return true
		}
		if !bytes.Equal(tx.Selector, txCtx.Internal.Selector) {
			return true
		}
		if !stringArrayIsEqual(tx.Calldata, txCtx.Internal.Calldata) {
			return true
		}
		return false
	case txCtx.Fee != nil:
		if tx.ContractID != txCtx.Fee.ContractID {
			return true
		}
		if tx.CallerID != txCtx.Fee.CallerID {
			return true
		}
		if !bytes.Equal(tx.Selector, txCtx.Fee.Selector) {
			return true
		}
		if !stringArrayIsEqual(tx.Calldata, txCtx.Fee.Calldata) {
			return true
		}
		return false
	case txCtx.Invoke != nil:
		if tx.ContractID != txCtx.Invoke.ContractID {
			return true
		}
		if !bytes.Equal(tx.Selector, txCtx.Invoke.EntrypointSelector) {
			return true
		}
		if !stringArrayIsEqual(tx.Calldata, txCtx.Invoke.CallData) {
			return true
		}
		return false
	default:
		return true
	}
}

func stringArrayIsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
