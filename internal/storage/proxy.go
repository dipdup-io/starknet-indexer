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

	ID         uint64
	ContractID uint64
	Hash       []byte
	Selector   []byte
	EntityType EntityType `pg:",use_zero"`
	EntityID   uint64
	EntityHash []byte

	Contract Address `pg:"rel:has-one,fk:contract_id"`
	Entity   Address `pg:"rel:has-one,fk:entity_id"`
}

// TableName -
func (Proxy) TableName() string {
	return "proxy"
}
