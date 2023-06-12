package data

import (
	"sync"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
)

// ProxyAction -
type ProxyAction int

// default proxy actions
const (
	ProxyActionAdd ProxyAction = iota
	ProxyActionUpdate
	ProxyActionDelete
)

// ProxyWithAction -
type ProxyWithAction struct {
	storage.Proxy
	Action ProxyAction
}

// NewProxyWithAction -
func NewProxyWithAction(proxy storage.Proxy, action ProxyAction) *ProxyWithAction {
	return &ProxyWithAction{proxy, action}
}

// ProxyUpgrade -
type ProxyUpgrade struct {
	Address  []byte
	Selector []byte
	Action   ProxyAction
	IsModule bool
}

// ProxyKey -
type ProxyKey struct {
	Address  string
	Selector string
}

// NewProxyKey -
func NewProxyKey(address, selector []byte) ProxyKey {
	key := ProxyKey{
		Address: encoding.EncodeHex(address),
	}
	if len(selector) > 0 {
		key.Selector = encoding.EncodeHex(selector)
	}

	return key
}

// NewProxyKeyFromString -
func NewProxyKeyFromString(address, selector string) ProxyKey {
	return ProxyKey{
		Address:  address,
		Selector: selector,
	}
}

// ProxyMap -
type ProxyMap[V any] struct {
	m  map[ProxyKey]V
	mx *sync.RWMutex
}

// NewProxyMap -
func NewProxyMap[V any]() ProxyMap[V] {
	return ProxyMap[V]{
		m:  make(map[ProxyKey]V),
		mx: new(sync.RWMutex),
	}
}

// Get -
func (pm ProxyMap[V]) Get(key ProxyKey) (V, bool) {
	pm.mx.RLock()
	defer pm.mx.RUnlock()
	if value, ok := pm.m[key]; ok {
		return value, ok
	}

	nilKey := NewProxyKeyFromString(key.Address, "")
	value, ok := pm.m[nilKey]
	return value, ok
}

// GetByHash -
func (pm ProxyMap[V]) GetByHash(address, selector []byte) (V, bool) {
	key := NewProxyKey(address, selector)
	return pm.Get(key)
}

// Add -
func (pm ProxyMap[V]) Add(key ProxyKey, value V) {
	pm.mx.Lock()
	pm.m[key] = value
	pm.mx.Unlock()
}

// AddByHash -
func (pm ProxyMap[V]) AddByHash(address, selector []byte, value V) {
	key := NewProxyKey(address, selector)
	pm.Add(key, value)
}

// Range -
func (pm ProxyMap[V]) Range(handler func(key ProxyKey, value V) (bool, error)) error {
	pm.mx.RLock()
	defer pm.mx.RUnlock()

	for key, value := range pm.m {
		stop, err := handler(key, value)
		if err != nil {
			return err
		}
		if stop {
			return nil
		}
	}
	return nil
}

// Delete -
func (pm ProxyMap[V]) Delete(key ProxyKey) {
	pm.mx.Lock()
	delete(pm.m, key)
	pm.mx.Unlock()
}

// Clone -
func (pm ProxyMap[V]) Clone() (ProxyMap[V], error) {
	newMap := NewProxyMap[V]()
	err := pm.Range(func(key ProxyKey, value V) (bool, error) {
		newMap.Add(key, value)
		return false, nil
	})
	return newMap, err
}
