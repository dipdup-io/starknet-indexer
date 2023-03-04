package parser

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/decode"
)

func (parser *Parser) getEvent(block storage.Block, contractAbi abi.Abi, event data.Event) (storage.Event, error) {
	model := storage.Event{
		Height: block.Height,
		Time:   block.Time,
		Order:  event.Order,
		Data:   event.Data,
		Keys:   event.Keys,
	}

	if len(contractAbi.Events) > 0 {
		parsed, name, err := decode.Event(parser.cache, contractAbi, model.Keys, model.Data)
		if err != nil {
			return model, err
		}
		model.ParsedData = parsed
		model.Name = name
	}

	return model, nil
}
func (parser *Parser) getMessage(ctx context.Context, block storage.Block, msg data.Message) (storage.Message, error) {
	message := storage.Message{
		Height:   block.Height,
		Time:     block.Time,
		Order:    msg.Order,
		Selector: msg.Selector,
		Payload:  msg.Payload,
		Nonce:    decimalFromHex(msg.Nonce),
	}

	if msg.FromAddress != "" {
		message.From = storage.Address{
			Hash: encoding.MustDecodeHex(msg.FromAddress),
		}

		if err := parser.findAddress(ctx, &message.From); err != nil {
			return message, err
		}
		message.FromID = message.From.ID
	}

	if msg.ToAddress != "" {
		message.To = storage.Address{
			Hash: encoding.MustDecodeHex(msg.ToAddress),
		}

		if err := parser.findAddress(ctx, &message.To); err != nil {
			return message, err
		}
		message.ToID = message.To.ID
	}

	return message, nil
}
