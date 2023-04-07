package subscriptions

import "github.com/dipdup-io/starknet-indexer/internal/storage"

// Message -
type Message struct {
	Block *storage.Block

	EndOfBlock bool
}

// NewBlockMessage -
func NewBlockMessage(block *storage.Block) *Message {
	return &Message{
		Block: block,
	}
}

// NewEndMessage -
func NewEndMessage() *Message {
	return &Message{
		EndOfBlock: true,
	}
}
