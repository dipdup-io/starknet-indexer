package generator

import (
	"context"
	"sync/atomic"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
)

// IdGenerator -
type IdGenerator struct {
	address storage.IAddress
	class   storage.IClass

	cache *cache.Cache
	state *storage.State

	addressId *atomic.Uint64
	classId   *atomic.Uint64
	txId      *atomic.Uint64
	eventId   *atomic.Uint64
}

// NewIdGenerator -
func NewIdGenerator(address storage.IAddress, class storage.IClass, cache *cache.Cache, state *storage.State) *IdGenerator {
	return &IdGenerator{
		address:   address,
		class:     class,
		cache:     cache,
		state:     state,
		addressId: new(atomic.Uint64),
		classId:   new(atomic.Uint64),
		txId:      new(atomic.Uint64),
		eventId:   new(atomic.Uint64),
	}
}

// Init -
func (gen *IdGenerator) Init() {
	gen.addressId.Store(gen.state.LastAddressID)
	gen.classId.Store(gen.state.LastClassID)
	gen.txId.Store(gen.state.LastTxID)
	gen.eventId.Store(gen.state.LastEventID)
}

// SetAddressId -
func (gen *IdGenerator) SetAddressId(ctx context.Context, address *storage.Address) (bool, error) {
	stored, err := gen.cache.GetAddress(ctx, address.Hash)
	if err != nil {
		if gen.class.IsNoRows(err) {
			address.ID = gen.NextAddressId()
			return true, nil
		}
		return false, err
	}
	address.ID = stored.ID
	address.ClassID = stored.ClassID
	return false, nil
}

// NextAddressId -
func (gen *IdGenerator) NextAddressId() uint64 {
	id := gen.addressId.Add(1)
	gen.state.LastAddressID = id
	return id
}

// SetClassId -
func (gen *IdGenerator) SetClassId(ctx context.Context, class *storage.Class) (bool, error) {
	stored, err := gen.cache.GetClassByHash(ctx, class.Hash)
	if err != nil {
		if gen.class.IsNoRows(err) {
			class.ID = gen.NextClassId()
			return true, nil
		}
		return false, err
	}
	class.ID = stored.ID
	class.Abi = stored.Abi
	class.Type = stored.Type
	return false, nil
}

// NextClassId -
func (gen *IdGenerator) NextClassId() uint64 {
	id := gen.classId.Add(1)
	gen.state.LastClassID = id
	return id
}

// NextTxId -
func (gen *IdGenerator) NextTxId() uint64 {
	id := gen.txId.Add(1)
	gen.state.LastTxID = id
	return id
}

// NextEventId -
func (gen *IdGenerator) NextEventId() uint64 {
	id := gen.eventId.Add(1)
	gen.state.LastEventID = id
	return id
}
