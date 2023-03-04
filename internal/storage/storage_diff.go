package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IStorageDiff -
type IStorageDiff interface {
	storage.Table[*StorageDiff]

	GetOnBlock(ctx context.Context, height, contractId uint64, key []byte) (StorageDiff, error)
}

// StorageDiff -
type StorageDiff struct {
	// nolint
	tableName struct{} `pg:"storage_diff"`

	ID         uint64
	Height     uint64 `pg:",use_zero"`
	ContractID uint64
	Key        []byte
	Value      []byte
}

// TableName -
func (StorageDiff) TableName() string {
	return "storage_diff"
}
