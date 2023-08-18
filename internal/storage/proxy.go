package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

// EntityType -
type EntityType int8

// entity types
const (
	EntityTypeClass = iota
	EntityTypeContract
)

// IProxy -
type IProxy interface {
	storage.Table[*Proxy]

	GetByHash(ctx context.Context, address, selector []byte) (Proxy, error)
}

// Proxy -
type Proxy struct {
	bun.BaseModel `bun:"proxy"`

	ID         uint64     `bun:",pk,autoincrement" comment:"Unique internal identity"`
	ContractID uint64     `comment:"Proxy contract id"`
	Hash       []byte     `comment:"Proxy contract hash"`
	Selector   []byte     `comment:"Proxy contract selector (for modules)"`
	EntityType EntityType `comment:"Entity type behind proxy (0 - class | 1 - contract)"`
	EntityID   uint64     `comment:"Entity id behind proxy"`
	EntityHash []byte     `comment:"Entity hash behind proxy"`

	Contract Address `bun:"rel:belongs-to,join:contract_id=id" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Entity   Address `bun:"rel:belongs-to,join:entity_id=id" hasura:"table:address,field:entity_id,remote_field:id,type:oto,name:entity"`
}

// TableName -
func (Proxy) TableName() string {
	return "proxy"
}
