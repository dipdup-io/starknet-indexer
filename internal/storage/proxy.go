package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
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
	// nolint
	tableName struct{} `pg:"proxy"`

	ID         uint64     `pg:",comment:Unique internal identity"`
	ContractID uint64     `pg:",comment:Proxy contract id"`
	Hash       []byte     `pg:",comment:Proxy contract hash"`
	Selector   []byte     `pg:",comment:Proxy contract selector (for modules)"`
	EntityType EntityType `pg:",use_zero,comment:Entity type behind proxy (0 - class | 1 - contract)"`
	EntityID   uint64     `pg:",comment:Entity id behind proxy"`
	EntityHash []byte     `pg:",comment:Entity hash behind proxy"`

	Contract Address `pg:"rel:has-one,fk:contract_id" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Entity   Address `pg:"rel:has-one,fk:entity_id" hasura:"table:address,field:entity_id,remote_field:id,type:oto,name:entity"`
}

// TableName -
func (Proxy) TableName() string {
	return "proxy"
}
