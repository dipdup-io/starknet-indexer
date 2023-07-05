package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// ProxyUpgrade -
type ProxyUpgrade struct {
	*postgres.Table[*storage.ProxyUpgrade]
}

// NewProxyUpgrade -
func NewProxyUpgrade(db *database.PgGo) *ProxyUpgrade {
	return &ProxyUpgrade{
		Table: postgres.NewTable[*storage.ProxyUpgrade](db),
	}
}

// LastBefore -
func (pu *ProxyUpgrade) LastBefore(ctx context.Context, height uint64) (upg storage.ProxyUpgrade, err error) {
	err = pu.DB().ModelContext(ctx, &upg).Where("height < ?", height).Last()
	return
}

// ListWithHeight -
func (pu *ProxyUpgrade) ListWithHeight(ctx context.Context, height uint64, limit, offset int) (upgrades []storage.ProxyUpgrade, err error) {
	if limit == 0 {
		limit = 10
	}

	err = pu.DB().ModelContext(ctx, &upgrades).
		Where("height > ?", height).
		Limit(limit).
		Offset(offset).
		Order("id desc").
		Select(&upgrades)

	return upgrades, err
}
