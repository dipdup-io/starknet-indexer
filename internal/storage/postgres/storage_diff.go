package postgres

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/go-pg/pg/v10/orm"
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

// InsertByCopy -
func (sd *StorageDiff) InsertByCopy(diffs []storage.StorageDiff) (io.Reader, string, error) {
	if len(diffs) == 0 {
		return nil, "", nil
	}
	builder := new(strings.Builder)

	for i := range diffs {
		if err := writeUint64(builder, diffs[i].Height); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64(builder, diffs[i].ContractID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeBytes(builder, diffs[i].Key); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeBytes(builder, diffs[i].Value); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte('\n'); err != nil {
			return nil, "", err
		}
	}

	query := fmt.Sprintf(`COPY %s (
		height, contract_id, key, value
	) FROM STDIN WITH (FORMAT csv, ESCAPE '\', QUOTE '"', DELIMITER ',')`, storage.StorageDiff{}.TableName())
	return strings.NewReader(builder.String()), query, nil
}

// Filter -
func (sd *StorageDiff) Filter(ctx context.Context, fltr []storage.StorageDiffFilter, opts ...storage.FilterOption) ([]storage.StorageDiff, error) {
	query := sd.DB().ModelContext(ctx, (*storage.StorageDiff)(nil))
	query = query.WhereGroup(func(q1 *orm.Query) (*orm.Query, error) {
		for i := range fltr {
			q1 = q1.WhereOrGroup(func(q *orm.Query) (*orm.Query, error) {
				q = integerFilter(q, "storage_diff.id", fltr[i].ID)
				q = integerFilter(q, "storage_diff.height", fltr[i].Height)
				q = addressFilter(q, "storage_diff.contract_id", fltr[i].Contract, "Contract")
				q = equalityFilter(q, "storage_diff.key", fltr[i].Key)
				return q, nil
			})
		}
		return q1, nil
	})
	query = optionsFilter(query, "storage_diff", opts...)
	query.Relation("Contract")

	var result []storage.StorageDiff
	err := query.Select(&result)
	return result, err
}
