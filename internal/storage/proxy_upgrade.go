package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// ProxyAction -
type ProxyAction int

// default proxy actions
const (
	ProxyActionAdd ProxyAction = iota
	ProxyActionUpdate
	ProxyActionDelete
)

// IProxy -
type IProxyUpgrade interface {
	storage.Table[*ProxyUpgrade]

	LastBefore(ctx context.Context, height uint64) (ProxyUpgrade, error)
	ListWithHeight(ctx context.Context, height uint64, limit, offset int) ([]ProxyUpgrade, error)
}

// ProxyUpgrade -
type ProxyUpgrade struct {
	// nolint
	tableName struct{} `pg:"proxy_upgrade"`

	ID         uint64      `pg:",comment:Unique internal identity"`
	ContractID uint64      `pg:",comment:Proxy contract id"`
	Hash       []byte      `pg:",comment:Proxy contract hash"`
	Selector   []byte      `pg:",comment:Proxy contract selector (for modules)"`
	EntityType EntityType  `pg:",use_zero,comment:Entity type behind proxy (0 - class | 1 - contract)"`
	EntityID   uint64      `pg:",comment:Entity id behind proxy"`
	EntityHash []byte      `pg:",comment:Entity hash behind proxy"`
	Height     uint64      `pg:",comment:Height when event occured"`
	Action     ProxyAction `pg:",comment:Action which occured with proxy (0 - add | 1 - update | 2 - delete)"`
}

// NewUpgradeFromProxy -
func NewUpgradeFromProxy(p Proxy) ProxyUpgrade {
	return ProxyUpgrade{
		ContractID: p.ContractID,
		Hash:       p.Hash,
		Selector:   p.Selector,
		EntityType: p.EntityType,
		EntityID:   p.EntityID,
		EntityHash: p.EntityHash,
	}
}

// TableName -
func (ProxyUpgrade) TableName() string {
	return "proxy_upgrade"
}

// ToProxy -
func (pu ProxyUpgrade) ToProxy() Proxy {
	return Proxy{
		ContractID: pu.ContractID,
		Hash:       pu.Hash,
		Selector:   pu.Selector,
		EntityType: pu.EntityType,
		EntityID:   pu.EntityID,
		EntityHash: pu.EntityHash,
	}
}
