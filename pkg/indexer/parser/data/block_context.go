package data

import "github.com/dipdup-io/starknet-indexer/internal/storage"

// BlockContext -
type BlockContext struct {
	block storage.Block

	addresses       map[string]*storage.Address
	classes         map[string]*storage.Class
	tokens          map[string]*storage.Token
	endBlockProxies ProxyMap[*storage.ProxyUpgrade]
	contextProxies  ProxyMap[*storage.ProxyUpgrade]
	classReplaces   map[string]*storage.ClassReplace
}

// NewBlockContext -
func NewBlockContext(block storage.Block) *BlockContext {
	return &BlockContext{
		block:           block,
		addresses:       make(map[string]*storage.Address),
		classes:         make(map[string]*storage.Class),
		tokens:          make(map[string]*storage.Token),
		endBlockProxies: NewProxyMap[*storage.ProxyUpgrade](),
		contextProxies:  NewProxyMap[*storage.ProxyUpgrade](),
		classReplaces:   make(map[string]*storage.ClassReplace),
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

// Tokens -
func (blockCtx *BlockContext) Tokens() map[string]*storage.Token {
	return blockCtx.tokens
}

// Proxies -
func (blockCtx *BlockContext) Proxies() ProxyMap[*storage.ProxyUpgrade] {
	return blockCtx.endBlockProxies
}

// CurrentProxies -
func (blockCtx *BlockContext) CurrentProxies() ProxyMap[*storage.ProxyUpgrade] {
	return blockCtx.contextProxies
}

// ClassReplaces -
func (blockCtx *BlockContext) ClassReplaces() map[string]*storage.ClassReplace {
	return blockCtx.classReplaces
}
