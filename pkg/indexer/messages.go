package indexer

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
)

// topics
const (
	OutputBlocks string = "blocks"
	InputStopSubsquid
)

// IndexerMessage -
type IndexerMessage struct {
	Block     *storage.Block
	Addresses map[string]*storage.Address
	Tokens    []*storage.Token
}

func (indexer *Indexer) notifyAllAboutBlock(
	blocks storage.Block,
	addresses map[string]*storage.Address,
	tokens map[string]*storage.Token,
) {
	newTokens := make([]*storage.Token, 0)
	for _, token := range tokens {
		if token.ID > 0 {
			newTokens = append(newTokens, token)
		}
	}

	output := indexer.MustOutput(OutputBlocks)
	output.Push(&IndexerMessage{
		Block:     &blocks,
		Addresses: addresses,
		Tokens:    newTokens,
	})
}
