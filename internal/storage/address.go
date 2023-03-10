package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IAddress -
type IAddress interface {
	storage.Table[*Address]

	GetByHash(ctx context.Context, hash []byte) (Address, error)
}

// Address -
type Address struct {
	// nolint
	tableName struct{} `pg:"address"`

	ID      uint64
	ClassID *uint64
	Hash    []byte `pg:",unique:address_hash"`

	Class Class `pg:"rel:has-one"`
}

// TableName -
func (Address) TableName() string {
	return "address"
}
