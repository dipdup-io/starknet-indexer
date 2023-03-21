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
	query := sd.DB().ModelContext(ctx, &diff).
		Where("contract_id = ?", contractId).
		Where("key = ?", key)

	if height > 0 {
		query = query.Where("height >= ?", height)
	}

	err = query.Order("id desc").
		Limit(1).
		Select(&diff)
	return
}
