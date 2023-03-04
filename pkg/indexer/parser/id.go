package parser

import (
	"context"
	"sync/atomic"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
)

type identifiable interface {
	LastID(ctx context.Context) (uint64, error)
}

func initializeCounter(ctx context.Context, i identifiable, counter *atomic.Uint64) error {
	id, err := i.LastID(ctx)
	if err != nil {
		return err
	}
	counter.Store(id)
	return nil
}

// IdGenerator -
type IdGenerator struct {
	address storage.IAddress
	class   storage.IClass

	cache *cache.Cache

	addressId *atomic.Uint64
	classId   *atomic.Uint64
}

// NewIdGenerator -
func NewIdGenerator(address storage.IAddress, class storage.IClass, cache *cache.Cache) *IdGenerator {
	return &IdGenerator{
		address:   address,
		class:     class,
		cache:     cache,
		addressId: new(atomic.Uint64),
		classId:   new(atomic.Uint64),
	}
}

// Init -
func (gen *IdGenerator) Init(ctx context.Context) error {
	if err := initializeCounter(ctx, gen.address, gen.addressId); err != nil {
		return err
	}
	if err := initializeCounter(ctx, gen.class, gen.classId); err != nil {
		return err
	}
	return nil
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
	return gen.addressId.Add(1)
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
	return gen.classId.Add(1)
}
