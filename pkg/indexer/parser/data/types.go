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

// Token -
type Token struct {
	ERC20   *storage.ERC20
	ERC721  *storage.ERC721
	ERC1155 *storage.ERC1155
}
