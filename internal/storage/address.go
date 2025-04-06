package storage

import (
	"context"
	"encoding/json"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IAddress interface {
	storage.Table[*Address]

	Filterable[Address, AddressFilter]

	GetByHash(ctx context.Context, hash []byte) (Address, error)
	GetAddresses(ctx context.Context, ids ...uint64) ([]Address, error)
	GetByHashes(ctx context.Context, hash [][]byte) ([]Address, error)
}

// AddressFilter -
type AddressFilter struct {
	ID           IntegerFilter
	Height       IntegerFilter
	OnlyStarknet bool
}

// Address -
type Address struct {
	bun.BaseModel `bun:"address" comment:"Table with starknet and ethereum addresses."`

	ID      uint64  `bun:"id,type:bigint,pk,notnull" json:"id" comment:"Unique internal identity"`
	ClassID *uint64 `bun:"class_id" json:"class_id" comment:"Class identity. It is NULL for ethereum addresses."`
	Height  uint64  `json:"height" comment:"Block number of the first address occurrence."`
	Hash    []byte  `bun:",unique:address_hash" json:"hash" comment:"Address hash."`

	Class Class `bun:"rel:belongs-to,join:class_id=id" json:"class" hasura:"table:class,field:class_id,remote_field:id,type:oto,name:class"`
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

func (address Address) MarshalJSON() ([]byte, error) {
	type Alias Address

	return json.Marshal(&struct {
		Alias
		Hash string `json:"hash"`
	}{
		Alias: Alias(address),
		Hash:  BytesToFormattedHex(address.Hash),
	})
}
