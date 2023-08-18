package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Blocks -
type Blocks struct {
	*postgres.Table[*storage.Block]
}

// NewBlocks -
func NewBlocks(db *database.Bun) *Blocks {
	return &Blocks{
		Table: postgres.NewTable[*storage.Block](db),
	}
}

// ByHeight -
func (b *Blocks) ByHeight(ctx context.Context, height uint64) (block storage.Block, err error) {
	err = b.DB().NewSelect().Model(&block).Where("height = ?", height).Limit(1).Scan(ctx)
	return
}

// Last -
func (b *Blocks) Last(ctx context.Context) (block storage.Block, err error) {
	err = b.DB().NewSelect().Model(&block).Order("height desc").Limit(1).Scan(ctx)
	return
}

// ByStatus -
func (b *Blocks) ByStatus(ctx context.Context, status storage.Status, limit, offset uint64, order sdk.SortOrder) (blocks []storage.Block, err error) {
	query := b.DB().NewSelect().Model(&blocks).
		Where("status = ?", status).
		Limit(int(limit)).
		Offset(int(offset))

	if order == sdk.SortOrderAsc {
		query.Order("id asc")
	} else {
		query.Order("id desc")
	}

	err = query.Scan(ctx)
	return
}
