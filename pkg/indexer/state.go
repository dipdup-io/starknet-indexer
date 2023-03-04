package indexer

import (
	"sync"

	models "github.com/dipdup-io/starknet-indexer/internal/storage"
)

type state struct {
	block *models.Block

	mx *sync.RWMutex
}

func newState(block *models.Block) *state {
	if block == nil {
		block = new(models.Block)
	}
	return &state{
		block: block,
		mx:    new(sync.RWMutex),
	}
}

// Height -
func (state *state) Height() uint64 {
	state.mx.RLock()
	defer state.mx.RUnlock()

	return state.block.Height
}

// Set -
func (state *state) Set(block models.Block) {
	state.mx.Lock()
	state.block = &block
	state.mx.Unlock()
}
