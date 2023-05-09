package data

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
)

// Result -
type Result struct {
	Block   storage.Block
	Context *BlockContext
	State   *storage.State
}

// entrypoint names
const (
	UpgradeEntrypoint = "upgrade"
)
