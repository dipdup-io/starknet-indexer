package data

import "github.com/dipdup-io/starknet-indexer/internal/storage"

// Result -
type Result struct {
	Addresses map[string]*storage.Address
	Block     storage.Block
	Classes   map[string]*storage.Class
	Proxies   map[string]*storage.Proxy
	State     *storage.State
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
