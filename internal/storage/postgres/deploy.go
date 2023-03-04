package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Deploy -
type Deploy struct {
	*postgres.Table[*storage.Deploy]
}

// NewDeploy -
func NewDeploy(db *database.PgGo) *Deploy {
	return &Deploy{
		Table: postgres.NewTable[*storage.Deploy](db),
	}
}

// ByHeight -
func (d *Deploy) ByHeight(ctx context.Context, height, limit, offset uint64) (response []storage.Deploy, err error) {
	err = d.DB().ModelContext(ctx, (*storage.Deploy)(nil)).
		Where("height = ?", height).
		Limit(int(limit)).
		Offset(int(offset)).
		Select(&response)
	return
}
