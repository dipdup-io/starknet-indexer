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
	GetIdsByHash(ctx context.Context, hash [][]byte) (ids []uint64, err error)
}

// Address -
type Address struct {
	// nolint
	tableName struct{} `pg:"address,comment:Table with starknet and ethereum addresses."`

	ID      uint64  `pg:"id,type:bigint,pk,notnull,comment:Unique internal identity"`
	ClassID *uint64 `pg:",comment:Class identity. It is NULL for ethereum addresses."`
	Height  uint64  `pg:",use_zero,comment:Block number of the first address occurrence."`
	Hash    []byte  `pg:",unique:address_hash,comment:Address hash."`

	Class Class `pg:"rel:has-one" hasura:"table:class,field:class_id,remote_field:id,type:oto,name:class"`
}

// TableName -
func (Address) TableName() string {
	return "address"
}
