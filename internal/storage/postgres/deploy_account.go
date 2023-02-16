package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// DeployAccount -
type DeployAccount struct {
	*postgres.Table[*storage.DeployAccount]
}

// NewDeployAccount -
func NewDeployAccount(db *database.PgGo) *DeployAccount {
	return &DeployAccount{
		Table: postgres.NewTable[*storage.DeployAccount](db),
	}
}

// ByHeight -
func (d *DeployAccount) ByHeight(ctx context.Context, height, limit, offset uint64) (response []storage.DeployAccount, err error) {
	err = d.DB().ModelContext(ctx, (*storage.DeployAccount)(nil)).
		Where("height = ?", height).
		Limit(int(limit)).
		Offset(int(offset)).
		Select(&response)
	return
}
