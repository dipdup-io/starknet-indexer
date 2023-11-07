package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

// ProxyAction -
type ProxyAction int

// default proxy actions
const (
	ProxyActionAdd ProxyAction = iota
	ProxyActionUpdate
	ProxyActionDelete
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IProxyUpgrade interface {
	storage.Table[*ProxyUpgrade]

	LastBefore(ctx context.Context, hash, selector []byte, height uint64) (ProxyUpgrade, error)
	ListWithHeight(ctx context.Context, height uint64, limit, offset int) ([]ProxyUpgrade, error)
}

// ProxyUpgrade -
type ProxyUpgrade struct {
	bun.BaseModel `bun:"proxy_upgrade"`

	ID         uint64      `bun:",pk,autoincrement" comment:":Unique internal identity"`
	ContractID uint64      `comment:":Proxy contract id"`
	Hash       []byte      `comment:":Proxy contract hash"`
	Selector   []byte      `comment:":Proxy contract selector (for modules)"`
	EntityType EntityType  `comment:"Entity type behind proxy (0 - class | 1 - contract)"`
	EntityID   uint64      `comment:":Entity id behind proxy"`
	EntityHash []byte      `comment:":Entity hash behind proxy"`
	Height     uint64      `comment:":Height when event occurred"`
	Action     ProxyAction `comment:":Action which occurred with proxy (0 - add | 1 - update | 2 - delete)"`
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
