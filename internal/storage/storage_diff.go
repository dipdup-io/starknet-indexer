package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

// IStorageDiff -
type IStorageDiff interface {
	storage.Table[*StorageDiff]
	Filterable[StorageDiff, StorageDiffFilter]

	GetOnBlock(ctx context.Context, height, contractId uint64, key []byte) (StorageDiff, error)
}

// StorageDiffFilter -
type StorageDiffFilter struct {
	ID       IntegerFilter
	Height   IntegerFilter
	Contract BytesFilter
	Key      EqualityFilter
}

// StorageDiff -
type StorageDiff struct {
	bun.BaseModel `bun:"storage_diff"`

	ID         uint64 `bun:",pk,autoincrement" comment:"Unique internal identity"`
	Height     uint64 `comment:"Block height"`
	ContractID uint64 `comment:"Contract id which storage was changed"`
	Key        []byte `comment:"Storage key"`
	Value      []byte `comment:"Data"`

	Contract Address `bun:"rel:belongs-to" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
}

// TableName -
func (StorageDiff) TableName() string {
	return "storage_diff"
}

// GetHeight -
func (sd StorageDiff) GetHeight() uint64 {
	return sd.Height
}

// GetId -
func (sd StorageDiff) GetId() uint64 {
	return sd.ID
}

// Columns -
func (StorageDiff) Columns() []string {
	return []string{
		"height", "contract_id", "key", "value",
	}
}

// Flat -
func (sd StorageDiff) Flat() []any {
	return []any{
		sd.Height,
		sd.ContractID,
		sd.Key,
		sd.Value,
	}
}
