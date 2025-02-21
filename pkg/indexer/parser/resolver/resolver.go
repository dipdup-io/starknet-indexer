package resolver

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
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
	blocks       storage.IBlock
	addresses    storage.IAddress
	proxies      storage.IProxy
	receiver     receiver.IReceiver
	cache        *cache.Cache
	idGenerator  *generator.IdGenerator
	blockContext *data.BlockContext
}

// NewResolver -
func NewResolver(
	receiver receiver.IReceiver,
	cache *cache.Cache,
	idGenerator *generator.IdGenerator,
	blocks storage.IBlock,
	addresses storage.IAddress,
	proxies storage.IProxy,
	blockContext *data.BlockContext,
) Resolver {
	return Resolver{
		receiver:     receiver,
		cache:        cache,
		idGenerator:  idGenerator,
		blocks:       blocks,
		addresses:    addresses,
		proxies:      proxies,
		blockContext: blockContext,
	}
}

// NextTxId -
func (resolver *Resolver) NextTxId() uint64 {
	return resolver.idGenerator.NextTxId()
}

// NextEventId -
func (resolver *Resolver) NextEventId() uint64 {
	return resolver.idGenerator.NextEventId()
}

// AddTokenToBlockContext -
func (resolver *Resolver) AddTokenToBlockContext(token *storage.Token) {
	if token == nil {
		return
	}
	key := token.String()
	tokens := resolver.blockContext.Tokens()
	if _, ok := tokens[key]; !ok {
		tokens[key] = token
	}
}
