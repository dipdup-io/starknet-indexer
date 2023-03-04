package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// StorageDiff -
type StorageDiff struct {
	*postgres.Table[*storage.StorageDiff]
}

// NewStorageDiff -
func NewStorageDiff(db *database.PgGo) *StorageDiff {
	return &StorageDiff{
		Table: postgres.NewTable[*storage.StorageDiff](db),
	}
}

// GetOnBlock -
func (sd *StorageDiff) GetOnBlock(ctx context.Context, height, contractId uint64, key []byte) (diff storage.StorageDiff, err error) {
	err = sd.DB().ModelContext(ctx, &diff).
		Where("height <= ?", height).
		Where("contract_id = ?", contractId).
		Where("key = ?", key).
		Order("id desc").
		Limit(1).
		Select(&diff)
	return
}
