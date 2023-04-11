package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IAddress -
type IAddress interface {
	storage.Table[*Address]

	GetByHash(ctx context.Context, hash []byte) (Address, error)
	GetAddresses(ctx context.Context, ids ...uint64) ([]Address, error)
}

// Address -
type Address struct {
	// nolint
	tableName struct{} `pg:"address"`

	ID      uint64 `pg:"id,type:bigint,pk,notnull"`
	ClassID *uint64
	Height  uint64 `pg:",use_zero"`
	Hash    []byte `pg:",unique:address_hash"`

	Class Class `pg:"rel:has-one"`
}

// TableName -
func (Address) TableName() string {
	return "address"
}
