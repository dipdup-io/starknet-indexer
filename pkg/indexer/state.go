package indexer

import (
	"sync"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
)

type state struct {
	state *storage.State

	mx *sync.RWMutex
}

func newState(s *storage.State) *state {
	if s == nil {
		s = new(storage.State)
	}
	return &state{
		state: s,
		mx:    new(sync.RWMutex),
	}
}

// Height -
func (state *state) Height() uint64 {
	state.mx.RLock()
	defer state.mx.RUnlock()

	return state.state.LastHeight
}

// Current -
func (state *state) Current() *storage.State {
	state.mx.RLock()
	defer state.mx.RUnlock()

	return state.state
}

// Set -
func (state *state) Set(s storage.State) {
	state.mx.Lock()
	defer state.mx.Unlock()

	state.state.DeclaresCount = s.DeclaresCount
	state.state.DeployAccountsCount = s.DeployAccountsCount
	state.state.DeploysCount = s.DeploysCount
	state.state.ID = s.ID
	state.state.InvokesCount = s.InvokesCount
	state.state.L1HandlersCount = s.L1HandlersCount
	state.state.LastAddressID = s.LastAddressID
	state.state.LastClassID = s.LastClassID
	state.state.LastEventID = s.LastEventID
	state.state.LastTxID = s.LastTxID
	state.state.LastHeight = s.LastHeight
	state.state.LastTime = s.LastTime
	state.state.Name = s.Name
	state.state.TxCount = s.TxCount
}
