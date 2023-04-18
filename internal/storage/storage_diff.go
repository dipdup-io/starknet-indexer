package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IStorageDiff -
type IStorageDiff interface {
	storage.Table[*StorageDiff]
	Copiable[StorageDiff]
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
	// nolint
	tableName struct{} `pg:"storage_diff"`

	ID         uint64 `pg:",comment:Unique internal identity"`
	Height     uint64 `pg:",use_zero,comment:Block height"`
	ContractID uint64 `pg:",comment:Contract id which storage was changed"`
	Key        []byte `pg:",comment:Storage key"`
	Value      []byte `pg:",comment:Data"`

	Contract Address `pg:"rel:has-one" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
}

// TableName -
func (StorageDiff) TableName() string {
	return "storage_diff"
}

// GetHeight -
func (sd StorageDiff) GetHeight() uint64 {
	return sd.Height
}
