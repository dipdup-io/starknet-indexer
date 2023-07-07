package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IAddress -
type IAddress interface {
	storage.Table[*Address]

	Filterable[Address, AddressFilter]

	GetByHash(ctx context.Context, hash []byte) (Address, error)
	GetAddresses(ctx context.Context, ids ...uint64) ([]Address, error)
	GetIdsByHash(ctx context.Context, hash [][]byte) (ids []uint64, err error)
}

// AddressFilter -
type AddressFilter struct {
	ID           IntegerFilter
	Height       IntegerFilter
	OnlyStarknet bool
}

// Address -
type Address struct {
	// nolint
	tableName struct{} `pg:"address" comment:"Table with starknet and ethereum addresses."`

	ID      uint64  `pg:"id,type:bigint,pk,notnull" comment:"Unique internal identity"`
	ClassID *uint64 `pg:"class_id" comment:"Class identity. It is NULL for ethereum addresses."`
	Height  uint64  `pg:",use_zero" comment:"Block number of the first address occurrence."`
	Hash    []byte  `pg:",unique:address_hash" comment:"Address hash."`

	Class Class `pg:"rel:has-one" hasura:"table:class,field:class_id,remote_field:id,type:oto,name:class"`
}

// TableName -
func (Address) TableName() string {
	return "address"
}

// GetHeight -
func (address Address) GetHeight() uint64 {
	return address.Height
}

// GetId -
func (address Address) GetId() uint64 {
	return address.ID
}
