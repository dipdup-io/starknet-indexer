package resolver

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/generator"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/receiver"
	"github.com/pkg/errors"
)

// errors
var (
	ErrUnknownProxy = errors.New("unknown proxy")
)

// Resolver -
type Resolver struct {
	blocks      storage.IBlock
	receiver    *receiver.Receiver
	cache       *cache.Cache
	idGenerator *generator.IdGenerator

	addresses       map[string]*storage.Address
	classes         map[string]*storage.Class
	endBlockProxies map[string]*storage.Proxy
	contextProxies  map[string]*storage.Proxy
}

// NewResolver -
func NewResolver(
	receiver *receiver.Receiver,
	cache *cache.Cache,
	idGenerator *generator.IdGenerator,
	blocks storage.IBlock,
) Resolver {
	return Resolver{
		receiver:    receiver,
		cache:       cache,
		idGenerator: idGenerator,
		blocks:      blocks,

		addresses:       make(map[string]*storage.Address),
		classes:         make(map[string]*storage.Class),
		endBlockProxies: make(map[string]*storage.Proxy),
		contextProxies:  make(map[string]*storage.Proxy),
	}
}

// Addresses -
func (resolver *Resolver) Addresses() map[string]*storage.Address {
	return resolver.addresses
}

// Classes -
func (resolver *Resolver) Classes() map[string]*storage.Class {
	return resolver.classes
}

// Proxies -
func (resolver *Resolver) Proxies() map[string]*storage.Proxy {
	return resolver.endBlockProxies
}

// NextTxId -
func (resolver *Resolver) NextTxId() uint64 {
	return resolver.idGenerator.NextTxId()
}

// NextEventId -
func (resolver *Resolver) NextEventId() uint64 {
	return resolver.idGenerator.NextEventId()
}
