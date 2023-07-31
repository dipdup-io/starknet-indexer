package indexer

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
)

// topics
const (
	OutputBlocks string = "blocks"
)

// IndexerMessage -
type IndexerMessage struct {
	Block     *storage.Block
	Addresses map[string]*storage.Address
	Tokens    []*storage.Token
}

// Input -
func (indexer *Indexer) Input(name string) (*modules.Input, error) {
	return nil, errors.Wrap(modules.ErrUnknownInput, name)
}

// Output -
func (indexer *Indexer) Output(name string) (*modules.Output, error) {
	output, ok := indexer.outputs[name]
	if !ok {
		return nil, errors.Wrap(modules.ErrUnknownInput, name)
	}
	return output, nil
}

// AttachTo -
func (indexer *Indexer) AttachTo(name string, input *modules.Input) error {
	output, err := indexer.Output(name)
	if err != nil {
		return err
	}
	output.Attach(input)
	return nil
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
	indexer.outputs[OutputBlocks].Push(&IndexerMessage{
		Block:     &blocks,
		Addresses: addresses,
		Tokens:    newTokens,
	})
}
