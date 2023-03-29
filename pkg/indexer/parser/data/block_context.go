package data

import "github.com/dipdup-io/starknet-indexer/internal/storage"

// BlockContext -
type BlockContext struct {
	block storage.Block

	addresses       map[string]*storage.Address
	classes         map[string]*storage.Class
	endBlockProxies ProxyMap[*ProxyWithAction]
	contextProxies  ProxyMap[*ProxyWithAction]
}

// NewBlockContext -
func NewBlockContext(block storage.Block) *BlockContext {
	return &BlockContext{
		block:           block,
		addresses:       make(map[string]*storage.Address),
		classes:         make(map[string]*storage.Class),
		endBlockProxies: NewProxyMap[*ProxyWithAction](),
		contextProxies:  NewProxyMap[*ProxyWithAction](),
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
func (blockCtx *BlockContext) Proxies() ProxyMap[*ProxyWithAction] {
	return blockCtx.endBlockProxies
}

// CurrentProxies -
func (blockCtx *BlockContext) CurrentProxies() ProxyMap[*ProxyWithAction] {
	return blockCtx.contextProxies
}
