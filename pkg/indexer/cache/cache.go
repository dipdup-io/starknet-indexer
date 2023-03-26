package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/karlseguin/ccache/v2"
	"github.com/pkg/errors"
)

const ttl = time.Hour

// Cache -
type Cache struct {
	*ccache.Cache

	address storage.IAddress
	class   storage.IClass
	proxy   storage.IProxy
}

// New -
func New(address storage.IAddress, class storage.IClass, proxy storage.IProxy) *Cache {
	return &Cache{
		Cache:   ccache.New(ccache.Configure().MaxSize(50000)),
		address: address,
		class:   class,
		proxy:   proxy,
	}
}

// GetAbiByAddress -
func (cache *Cache) GetAbiByAddress(ctx context.Context, hash []byte) (abi.Abi, error) {
	item, err := cache.Fetch(fmt.Sprintf("abi:address:%x", hash), ttl, func() (interface{}, error) {
		address, err := cache.GetAddress(ctx, hash)
		if err != nil {
			return nil, err
		}

		if address.ClassID == nil {
			return nil, errors.Errorf("unknown class id for address: %x", hash)
		}

		class, err := cache.class.GetByID(ctx, *address.ClassID)
		if err != nil {
			return nil, err
		}

		return class.Abi, nil
	})
	if err != nil {
		return abi.Abi{}, err
	}

	var a abi.Abi
	err = json.Unmarshal(item.Value().(storage.Bytes), &a)
	return a, err
}

// SetAbiByAddress -
func (cache *Cache) SetAbiByAddress(class storage.Class, hash []byte) {
	cache.Set(fmt.Sprintf("abi:address:%x", hash), class.Abi, ttl)
}

// GetAbiByClassHash -
func (cache *Cache) GetAbiByClassHash(ctx context.Context, hash []byte) (abi.Abi, error) {
	item, err := cache.Fetch(fmt.Sprintf("abi:class_hash:%x", hash), ttl, func() (interface{}, error) {
		class, err := cache.class.GetByHash(ctx, hash)
		if err != nil {
			return nil, err
		}

		var a abi.Abi
		err = json.Unmarshal(class.Abi, &a)
		return a, err
	})
	if err != nil {
		return abi.Abi{}, err
	}

	return item.Value().(abi.Abi), nil
}

// SetAbiByClassHash -
func (cache *Cache) SetAbiByClassHash(class storage.Class, a abi.Abi) {
	cache.Set(fmt.Sprintf("abi:class_hash:%x", class.Hash), a, ttl)
}

// GetClassByHash -
func (cache *Cache) GetClassByHash(ctx context.Context, hash []byte) (storage.Class, error) {
	item, err := cache.Fetch(fmt.Sprintf("class:hash:%x", hash), ttl, func() (interface{}, error) {
		return cache.class.GetByHash(ctx, hash)
	})
	if err != nil {
		return storage.Class{}, err
	}

	return item.Value().(storage.Class), nil
}

// SetClassByHash -
func (cache *Cache) SetClassByHash(class storage.Class) {
	cache.Set(fmt.Sprintf("class:hash:%x", class.Hash), class, ttl)
}

// GetAddress -
func (cache *Cache) GetAddress(ctx context.Context, hash []byte) (storage.Address, error) {
	item, err := cache.Fetch(fmt.Sprintf("address:hash:%x", hash), ttl, func() (interface{}, error) {
		return cache.address.GetByHash(ctx, hash)
	})
	if err != nil {
		return storage.Address{}, err
	}

	address := item.Value().(storage.Address)
	if address.ClassID == nil {
		address, err = cache.address.GetByHash(ctx, hash)
		if err != nil {
			return address, err
		}
		cache.Set(fmt.Sprintf("address:hash:%x", hash), address, ttl)
	}

	return address, nil
}

// GetAbiByAddressId -
func (cache *Cache) GetAbiByAddressId(ctx context.Context, id uint64) (abi.Abi, error) {
	item, err := cache.Fetch(fmt.Sprintf("abi:address_id:%d", id), ttl, func() (interface{}, error) {
		address, err := cache.address.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}

		if address.ClassID == nil {
			return nil, errors.Errorf("unknown class id for address id: %d", id)
		}

		class, err := cache.class.GetByID(ctx, *address.ClassID)
		if err != nil {
			return nil, err
		}

		return class.Abi, nil
	})
	if err != nil {
		return abi.Abi{}, err
	}

	var a abi.Abi
	err = json.Unmarshal(item.Value().(storage.Bytes), &a)
	return a, err
}

// GetClassById -
func (cache *Cache) GetClassById(ctx context.Context, id uint64) (*storage.Class, error) {
	item, err := cache.Fetch(fmt.Sprintf("class:id:%d", id), ttl, func() (interface{}, error) {
		return cache.class.GetByID(ctx, id)
	})
	if err != nil {
		return nil, err
	}

	return item.Value().(*storage.Class), nil
}

// GetClassForAddress -
func (cache *Cache) GetClassForAddress(ctx context.Context, hash []byte) (storage.Class, error) {
	item, err := cache.Fetch(fmt.Sprintf("class:address:%x", hash), ttl, func() (interface{}, error) {
		address, err := cache.address.GetByHash(ctx, hash)
		if err != nil {
			return nil, err
		}

		if address.ClassID != nil {
			return cache.GetClassById(ctx, *address.ClassID)
		}

		return nil, errors.Errorf("address is not a starknet contract: hash=%x", hash)
	})
	if err != nil {
		return storage.Class{}, err
	}
	class := item.Value().(*storage.Class)
	return *class, err
}
