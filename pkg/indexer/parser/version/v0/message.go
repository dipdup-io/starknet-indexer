package v0

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/resolver"
)

// MessageParser -
type MessageParser struct {
	cache    *cache.Cache
	resolver resolver.Resolver
}

// NewMessageParser -
func NewMessageParser(
	cache *cache.Cache,
	resolver resolver.Resolver,
) MessageParser {
	return MessageParser{cache: cache, resolver: resolver}
}

// ParseMessage -
func (parser MessageParser) Parse(ctx context.Context, txCtx parserData.TxContext, msg data.Message) (storage.Message, error) {
	message := storage.Message{
		Height:          txCtx.Height,
		Time:            txCtx.Time,
		Order:           msg.Order,
		Selector:        msg.Selector.String(),
		Payload:         msg.Payload,
		Nonce:           msg.Nonce.Decimal(),
		ContractID:      txCtx.ContractId,
		DeclareID:       txCtx.DeclareID,
		DeployID:        txCtx.DeployID,
		DeployAccountID: txCtx.DeployAccountID,
		InvokeID:        txCtx.InvokeID,
		L1HandlerID:     txCtx.L1HandlerID,
		FeeID:           txCtx.FeeID,
		InternalID:      txCtx.InternalID,
	}
	if txCtx.ProxyId > 0 {
		message.ContractID = txCtx.ProxyId
	}

	if msg.FromAddress != "" {
		message.From = storage.Address{
			Hash:   data.Felt(msg.FromAddress).Bytes(),
			Height: message.Height,
		}

		if err := parser.resolver.FindAddress(ctx, &message.From); err != nil {
			return message, err
		}
		message.FromID = message.From.ID
	}

	if msg.ToAddress != "" {
		message.To = storage.Address{
			Hash:   data.Felt(msg.ToAddress).Bytes(),
			Height: message.Height,
		}

		if err := parser.resolver.FindAddress(ctx, &message.To); err != nil {
			return message, err
		}
		message.ToID = message.To.ID
	}

	return message, nil
}
