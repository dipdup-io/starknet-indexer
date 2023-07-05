package data

import "github.com/dipdup-io/starknet-indexer/internal/storage"

// BlockContext -
type BlockContext struct {
	block storage.Block

	addresses       map[string]*storage.Address
	classes         map[string]*storage.Class
	endBlockProxies ProxyMap[*storage.ProxyUpgrade]
	contextProxies  ProxyMap[*storage.ProxyUpgrade]
}

// NewBlockContext -
func NewBlockContext(block storage.Block) *BlockContext {
	return &BlockContext{
		block:           block,
		addresses:       make(map[string]*storage.Address),
		classes:         make(map[string]*storage.Class),
		endBlockProxies: NewProxyMap[*storage.ProxyUpgrade](),
		contextProxies:  NewProxyMap[*storage.ProxyUpgrade](),
	}
}

// Block -
func (blockCtx *BlockContext) Block() storage.Block {
	return blockCtx.block
}

// Addresses -
func (blockCtx *BlockContext) Addresses() map[string]*storage.Address {
	return blockCtx.addresses
}

// Classes -
func (blockCtx *BlockContext) Classes() map[string]*storage.Class {
	return blockCtx.classes
}

// Proxies -
func (blockCtx *BlockContext) Proxies() ProxyMap[*storage.ProxyUpgrade] {
	return blockCtx.endBlockProxies
}

// CurrentProxies -
func (blockCtx *BlockContext) CurrentProxies() ProxyMap[*storage.ProxyUpgrade] {
	return blockCtx.contextProxies
}
