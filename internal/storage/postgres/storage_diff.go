package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// StorageDiff -
type StorageDiff struct {
	*postgres.Table[*storage.StorageDiff]
}

// NewStorageDiff -
func NewStorageDiff(db *database.Bun) *StorageDiff {
	return &StorageDiff{
		Table: postgres.NewTable[*storage.StorageDiff](db),
	}
}

// GetOnBlock -
func (sd *StorageDiff) GetOnBlock(ctx context.Context, height, contractId uint64, key []byte) (diff storage.StorageDiff, err error) {
	err = sd.DB().NewSelect().Model(&diff).
		Where("contract_id = ?", contractId).
		Where("key = ?", key).
		Where("height >= ?", height).
		Order("id desc").
		Limit(1).
		Scan(ctx)
	return
}

// Filter -
func (sd *StorageDiff) Filter(ctx context.Context, fltr []storage.StorageDiffFilter, opts ...storage.FilterOption) (result []storage.StorageDiff, err error) {
	query := sd.DB().NewSelect().Model(&result)
	query = query.WhereGroup(" AND ", func(q1 *bun.SelectQuery) *bun.SelectQuery {
		for i := range fltr {
			q1 = q1.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
				q = integerFilter(q, "storage_diff.id", fltr[i].ID)
				q = integerFilter(q, "storage_diff.height", fltr[i].Height)
				q = addressFilter(q, "hash", fltr[i].Contract, "Contract")
				q = equalityFilter(q, "storage_diff.key", fltr[i].Key)
				return q
			})
		}
		return q1
	})
	query = optionsFilter(query, "storage_diff", opts...)
	query.Relation("Contract")

	err = query.Scan(ctx)
	return
}
